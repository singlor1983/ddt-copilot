package data

import (
	"ddt-copilot/logger"
)

var (
	GLog         *logger.CustomLogger // 线程安全的全局logger
	GDefsAngle   *DefsAngle           // 识别-角度
	GDefsFubenLv *DefsFubenLv         // 识别-副本难度
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

	GGameSetting = &GameSetting{}
	GGameSetting.Init()
}
