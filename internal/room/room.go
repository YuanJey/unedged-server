package room

import (
	"context"
	"log"
	"sync"
	"time"
)

// CombatDirection 三向攻防枚举
type CombatDirection uint8

const (
	DirNone  CombatDirection = 0
	DirTop   CombatDirection = 1
	DirLeft  CombatDirection = 2
	DirRight CombatDirection = 3
)

// GripStance 握持姿态枚举
type GripStance uint8

const (
	GripOneHandR GripStance = 0
	GripOneHandL GripStance = 1
	GripTwoHand  GripStance = 2
)

// PlayerInputFrame 单玩家单逻辑帧输入数据包
type PlayerInputFrame struct {
	PlayerID      uint64          `json:"player_id"`
	FrameID       uint32          `json:"frame_id"`
	MoveX         int32           `json:"move_x"` // 定点数 (-1000 ~ 1000)
	MoveY         int32           `json:"move_y"`
	GuardDir      CombatDirection `json:"guard_dir"`
	ActionBitmask uint32          `json:"action_bitmask"`
}

// DuelRoom 1v1 论剑房间实例
type DuelRoom struct {
	RoomID       string
	PlayerIDs    [2]uint64
	CurrentFrame uint32
	frameInputs  map[uint32][2]*PlayerInputFrame
	mu           sync.RWMutex
	cancelFunc   context.CancelFunc
}

// NewDuelRoom 初始化一个新的论剑房间
func NewDuelRoom(roomID string, p1 uint64, p2 uint64) *DuelRoom {
	return &DuelRoom{
		RoomID:      roomID,
		PlayerIDs:   [2]uint64{p1, p2},
		frameInputs: make(map[uint32][2]*PlayerInputFrame),
	}
}

// PushInput 接收玩家在特定帧的输入指令
func (r *DuelRoom) PushInput(input *PlayerInputFrame) {
	r.mu.Lock()
	defer r.mu.Unlock()

	frames := r.frameInputs[input.FrameID]
	if input.PlayerID == r.PlayerIDs[0] {
		frames[0] = input
	} else if input.PlayerID == r.PlayerIDs[1] {
		frames[1] = input
	}
	r.frameInputs[input.FrameID] = frames
}

// StartFrameLoop 启动 60Hz 逻辑帧同步循环 (约 16.66ms per tick)
func (r *DuelRoom) StartFrameLoop(ctx context.Context) {
	ctx, r.cancelFunc = context.WithCancel(ctx)
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	log.Printf("[Room %s] 60Hz 逻辑帧循环已启动，玩家: %d vs %d", r.RoomID, r.PlayerIDs[0], r.PlayerIDs[1])

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Room %s] 帧循环停止", r.RoomID)
			return
		case <-ticker.C:
			r.tick()
		}
	}
}

// tick 推进单逻辑帧并准备广播全量输入
func (r *DuelRoom) tick() {
	r.mu.Lock()
	r.CurrentFrame++
	current := r.CurrentFrame
	inputs := r.frameInputs[current]
	// 清理旧帧缓存（保留最近 120 帧用于丢包重传）
	if current > 120 {
		delete(r.frameInputs, current-120)
	}
	r.mu.Unlock()

	// 实际生产中在此通过 KCP session 将 inputs 广播给双方客户端
	_ = inputs
}

// Stop 停止房间帧循环
func (r *DuelRoom) Stop() {
	if r.cancelFunc != nil {
		r.cancelFunc()
	}
}
