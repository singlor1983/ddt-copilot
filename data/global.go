package data

import (
	"ddt-copilot/defs"
	"ddt-copilot/logger"
)

var (
	GLog         *logger.CustomLogger // 线程安全的全局logger
	GDefsAngle   *DefsAngle           // 识别-角度
	GDefsFubenLv *DefsFubenLv         // 识别-副本难度
	GDefsOther   *DefsOther           // 识别-杂项
	GGameSetting *GameSetting         // 游戏配置
)

func Log() *logger.CustomLogger {
	return GLog
}

func InitGlobal() {
	GLog = logger.NewConsoleLogger() // 不写入文件，只在控制台打印

	GDefsAngle = &DefsAngle{}
	GDefsAngle.Init()

	GDefsFubenLv = &DefsFubenLv{}
	GDefsFubenLv.Init()

	GDefsOther = &DefsOther{}
	GDefsOther.Init()

	GGameSetting = &GameSetting{}
	GGameSetting.Init()
	GGameSetting.SetSettingGeneral(SettingGeneral{
		IsBossFightEnable:      true,
		AttackCMD:              nil,
		UsePetFoodByFightCount: 10,
	})
	GGameSetting.SetSettingFubenPosition()
	GGameSetting.SetSettingFubenCustom(SettingFubenCustom{
		name:      "custom",
		Angle:     55,
		Power:     80,
		Direction: defs.DirectionRight,
		AttackCMD: nil,
		FubenSetting: defs.FubenSetting{
			Lv:                defs.FubenLvDifficulty,
			IsBossFightEnable: false,
			FubenPosition: defs.FubenPosition{
				Type:  defs.FubenTypeNormal,
				Page:  2,
				Index: 5,
			},
		},
	})
}
