package defs

import "time"

type ProcessName string

const (
	ProcessTgWeb ProcessName = "TangoWeb.exe"
)

const (
	VK_0      uintptr = 48
	VK_1      uintptr = 49
	VK_2      uintptr = 50
	VK_3      uintptr = 51
	VK_4      uintptr = 52
	VK_5      uintptr = 53
	VK_6      uintptr = 54
	VK_7      uintptr = 55
	VK_8      uintptr = 56
	VK_9      uintptr = 57
	VK_B      uintptr = 66
	VK_Q      uintptr = 81
	VK_E      uintptr = 69
	VK_T      uintptr = 84
	VK_Y      uintptr = 89
	VK_U      uintptr = 85
	VK_P      uintptr = 80
	VK_F      uintptr = 70
	VK_SPACE  uintptr = 32
	VK_LEFT   uintptr = 37
	VK_UP     uintptr = 38
	VK_RIGHT  uintptr = 39
	VK_DOWN   uintptr = 40
	VK_ESCAPE uintptr = 27
)

var vkMap = map[string]uintptr{
	"0":      VK_0,
	"1":      VK_1,
	"2":      VK_2,
	"3":      VK_3,
	"4":      VK_4,
	"5":      VK_5,
	"6":      VK_6,
	"7":      VK_7,
	"8":      VK_8,
	"9":      VK_9,
	"B":      VK_B,
	"Q":      VK_Q,
	"E":      VK_E,
	"T":      VK_T,
	"Y":      VK_Y,
	"U":      VK_U,
	"P":      VK_P,
	"F":      VK_F,
	"SPACE":  VK_SPACE,
	"LEFT":   VK_LEFT,
	"UP":     VK_UP,
	"RIGHT":  VK_RIGHT,
	"DOWN":   VK_DOWN,
	"ESCAPE": VK_ESCAPE,
}

func GetVkFromStr(key string) uintptr {
	v, ok := vkMap[key]
	if !ok {
		return 0
	}
	return v
}

const (
	Colorthreshold int = 30 // 图片相似度阈值
)

type FubenLv int

const (
	FubenLvEasy       FubenLv = 1 // 简单
	FubenLvNormal     FubenLv = 2 // 普通
	FubenLvDifficulty FubenLv = 3 // 困难
	FubenLvHero       FubenLv = 4 // 英雄
	FubenLvNightmare  FubenLv = 5 // 噩梦
)

type Direction int

const (
	DirectionLeft  Direction = 1
	DirectionRight Direction = 2
	DirectionUp    Direction = 3
	DirectionDown  Direction = 4
)

type FunctionID int

const (
	FunctionIDFubenBegin  FunctionID = 0    // 副本类战斗ID开始
	FunctionIDCustomFuben FunctionID = 1    // 自定义副本
	FunctionIDMaYiGeneral FunctionID = 2    // 蚂蚁-经典
	FunctionIDFubenEnd    FunctionID = 1000 // 副本类战斗ID结束

	FunctionIDJinjiBegin FunctionID = 1000 // 竞技类战斗ID开始
	FunctionIDJinjiEnd   FunctionID = 2000 // 竞技类战斗ID结束

	FunctionIDOtherBegin FunctionID = 2000 // 其他功能类脚本ID开始
	FunctionIDOtherEnd   FunctionID = 3000 // 其他功能类脚本ID开始
)

type ReadyState int

const (
	ReadyStateNo ReadyState = 0 // 未准备
	ReadyStateOK ReadyState = 1 // 已准备
)

type InitPosition int

const (
	InitPosition1 InitPosition = 1
	InitPosition2 InitPosition = 2
	InitPosition3 InitPosition = 3
	InitPosition4 InitPosition = 4
)

const (
	TimeWaitLong  = time.Millisecond * 2000 // 长久等待，比如切换场景，从大厅->副本房间内
	TimeWaitMid   = time.Millisecond * 700  // 中等等待，点击后的停顿，这个时间刚好够人类反应了
	TimeWaitShort = time.Millisecond * 150  // 短暂等待，比如准备副本道具
)

type FubenType int

const (
	FubenTypeNormal  FubenType = 0 // 普通副本
	FubenTypeSpecial FubenType = 1 // 特殊副本
)

type FubenPosition struct {
	Type  FubenType
	Page  int
	Index RectType // 1-8
}

type FubenSetting struct {
	Lv                FubenLv // 副本难度
	IsBossFightEnable bool    // 是否启用boss战
	FubenPosition
}
