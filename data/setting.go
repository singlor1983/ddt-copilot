package data

type SettingGeneral struct {
	IsBossFightEnable bool     // 是否启用boss战
	AttackCMD         []string // 战斗攻击指令
}

type GameSetting struct {
	SettingGeneral
}

func (self *GameSetting) Init() {
}
