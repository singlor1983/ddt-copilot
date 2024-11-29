package data

import "ddt-copilot/defs"

type SettingGeneral struct {
	IsBossFightEnable bool     // 是否启用boss战
	AttackCMD         []string // 战斗攻击指令
}

type SettingFubenPosition struct {
	Position map[defs.FunctionID]defs.FubenPosition
}

type SettingFubenCustom struct {
	name      string
	Angle     int
	Power     int
	Direction defs.Direction // 攻击方向
	AttackCMD []string       // 战斗攻击指令

	defs.FubenSetting
}

type GameSetting struct {
	SettingGeneral
	SettingFubenPosition
	SettingFubenCustom
}

func (self *GameSetting) Init() {
}

func (self *GameSetting) SetSettingFubenPosition() {
	self.SettingFubenPosition.Position = map[defs.FunctionID]defs.FubenPosition{
		defs.FunctionIDMaYiGeneral: {Type: defs.FubenTypeNormal, Page: 1, Index: 1},
	}
}

func (self *GameSetting) SetSettingFubenCustom(st SettingFubenCustom) {
	self.SettingFubenCustom = st
}
