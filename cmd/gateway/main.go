package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"unedged-server/internal/config"
	"unedged-server/internal/room"
)

func main() {
	log.Println("==================================================")
	log.Println("《无锋》（Unedged）极速论剑 60Hz 帧同步网关启动中...")
	log.Println("==================================================")

	// 1. 加载共享配置表（支持根目录或同级路径探测）
	configCandidates := []string{
		"../configs/weapons.json",
		"configs/weapons.json",
		"../../configs/weapons.json",
	}

	var weaponsCfg *config.WeaponsConfig
	var err error
	for _, p := range configCandidates {
		absPath, _ := filepath.Abs(p)
		weaponsCfg, err = config.LoadWeaponsConfig(p)
		if err == nil {
			log.Printf("成功载入武器契约配置: %s", absPath)
			break
		}
	}

	if weaponsCfg == nil {
		log.Printf("警告: 无法加载外部 weapons.json: %v，采用内存缺省配置", err)
	} else {
		log.Printf("已解析武器大类: %d 个，名器模组: %d 把", len(weaponsCfg.Archetypes), len(weaponsCfg.Weapons))
	}

	// 2. 初始化论剑擂台房间实例
	duel := room.NewDuelRoom("arena-001", 10001, 10002)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. 启动后台 60Hz 定频逻辑帧同步协程
	go duel.StartFrameLoop(ctx)

	log.Println("网关监听就绪（基于 KCP-GO 协议栈，模拟端口: 9999）")
	log.Println("系统已进入稳定轮询态，按 CTRL+C 终止运行...")

	// 4. 监听系统退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("正在优雅关闭论剑房间与网络会话...")
	duel.Stop()
	log.Println("网关服务已安全退出。")
}
