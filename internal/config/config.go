package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// WeaponsConfig 定义共享武器配置表结构
type WeaponsConfig struct {
	Archetypes map[string]Archetype `json:"archetypes"`
	Weapons    []WeaponItem         `json:"weapons"`
}

// Archetype 武器大类基础属性
type Archetype struct {
	Name          string            `json:"name"`
	BaseStats     BaseStats         `json:"base_stats"`
	Stances       map[string]Stance `json:"stances"`
	DefaultSkill1 string            `json:"default_skill_1"`
}

// BaseStats 基础判定帧数据与属性
type BaseStats struct {
	MaxHP                    int `json:"max_hp"`
	MaxStamina               int `json:"max_stamina"`
	MaxPoise                 int `json:"max_poise"`
	LightAttackStartupFrames int `json:"light_attack_startup_frames"`
	LightAttackActiveFrames  int `json:"light_attack_active_frames"`
	LightAttackRecoveryFrames int `json:"light_attack_recovery_frames"`
	ParryWindowFrames        int `json:"parry_window_frames"`
}

// Stance 姿态修正倍率
type Stance struct {
	MoveSpeedMult            float64 `json:"move_speed_mult"`
	DamageMult               float64 `json:"damage_mult"`
	GuardStaminaCostMult     float64 `json:"guard_stamina_cost_mult"`
	HasHyperArmorOnTopHeavy  bool    `json:"has_hyper_armor_on_top_heavy"`
}

// WeaponItem 具体名器定义与代偿机制
type WeaponItem struct {
	WeaponID        string            `json:"weapon_id"`
	WeaponName      string            `json:"weapon_name"`
	Archetype       string            `json:"archetype"`
	SkillOverrides  map[string]string `json:"skill_overrides"`
	BalanceTradeoff string            `json:"balance_tradeoff"`
}

// LoadWeaponsConfig 从指定路径加载武器配置文件
func LoadWeaponsConfig(path string) (*WeaponsConfig, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg WeaponsConfig
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("反序列化武器配置失败: %w", err)
	}

	return &cfg, nil
}
