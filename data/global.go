package data

import (
	"ddt-copilot/logger"
)

var (
	GLog       *logger.CustomLogger // 线程安全的全局logger
	GDefsAngle *DefsAngle           // 角度识别
)

func Log() *logger.CustomLogger {
	return GLog
}

func InitGlobal() {
	GLog = logger.NewConsoleLogger() // 不写入文件，只在控制台打印

	GDefsAngle = &DefsAngle{}
	GDefsAngle.Init()
}
