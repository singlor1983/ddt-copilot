package defs

import (
	"github.com/lxn/win"
	"image"
	"time"
)

type ProcessName string

const (
	ProcessTg3   ProcessName = "Tango3.exe"
	ProcessTgWeb ProcessName = "TangoWeb.exe"
)

const (
	DdtWindowStdWidth  = 1000 // 弹弹堂游戏窗口标准宽度，窗口宽度低于这个时不进行操作
	DdtWindowStdHeight = 600  // 弹弹堂游戏窗口标准高度，窗口高度低于这个时不进行操作
)

type Point [2]int32

var EmptyPoint Point = [2]int32{-1, -1}

func IsEmptyPoint(point Point) bool {
	return point[0] == EmptyPoint[0] && point[1] == EmptyPoint[1]
}

func NewPointWithRectCenter(rect win.RECT) Point {
	return Point{int32(rect.Right+rect.Left) / 2, int32(rect.Bottom+rect.Top) / 2}
}

type DDTClickElement int

const (
	ElementEmpty              DDTClickElement = 1  // 空白处，用于激活窗口
	ElementFubenEnter         DDTClickElement = 2  // 副本开始按钮
	ElementFubenSelect        DDTClickElement = 3  // 副本选择按钮
	ElementFubenPageDown      DDTClickElement = 4  // 副本下拉按钮
	ElementFubenTypeNormal    DDTClickElement = 5  // 普通副本
	ElementFubenTypeSpecial   DDTClickElement = 6  // 特殊副本
	ElementFubenBossFight     DDTClickElement = 7  // 副本BOSS战按钮
	ElementFubenAck           DDTClickElement = 8  // 副本确认按钮
	ElementFubenBossFightAck  DDTClickElement = 9  // 副本BOSS战确认按钮
	ElementFightEquipItem1    DDTClickElement = 10 // 房间内战斗已装备道具1
	ElementFightEquipItem2    DDTClickElement = 11 // 房间内战斗已装备道具2
	ElementFightEquipItem3    DDTClickElement = 12 // 房间内战斗已装备道具3
	ElementFightSelectItem1   DDTClickElement = 13 // 房间内供选择道具1
	ElementFightSelectItem2   DDTClickElement = 14 // 房间内供选择道具2
	ElementFightSelectItem3   DDTClickElement = 15 // 房间内供选择道具3
	ElementFightSelectItem4   DDTClickElement = 16 // 房间内供选择道具4
	ElementFightSelectItem5   DDTClickElement = 17 // 房间内供选择道具5
	ElementFightSelectItem6   DDTClickElement = 18 // 房间内供选择道具6
	ElementFightSelectItem7   DDTClickElement = 19 // 房间内供选择道具7
	ElementFightSelectItem8   DDTClickElement = 20 // 房间内供选择道具8
	ElementFightStart         DDTClickElement = 21 // 战斗开始按钮
	ElementFubenFightStartAck DDTClickElement = 22 // 副本战斗开始确认按钮
	ElementBackAndExit        DDTClickElement = 23 // 返回&&推出按钮
	ElementFubenHall          DDTClickElement = 24 // 进入副本大厅的按钮
)

// elementRelativePosition 可点击元素中心相对于左上角顶点的坐标
var elementRelativePosition = map[DDTClickElement]Point{
	ElementEmpty:              {1, 1},     // 空白处肯定不是屏幕中心
	ElementFubenEnter:         {722, 479}, // 从副本大厅进入副本房间
	ElementFubenSelect:        {597, 234},
	ElementFubenPageDown:      {768, 440},
	ElementFubenTypeNormal:    {287, 305},
	ElementFubenTypeSpecial:   {427, 305},
	ElementFubenBossFight:     {320, 565},
	ElementFubenAck:           {500, 565},
	ElementFubenBossFightAck:  {433, 340},
	ElementFightEquipItem1:    {810, 140},
	ElementFightEquipItem2:    {880, 140},
	ElementFightEquipItem3:    {950, 140},
	ElementFightSelectItem1:   {807, 234},
	ElementFightSelectItem2:   {857, 234},
	ElementFightSelectItem3:   {907, 234},
	ElementFightSelectItem4:   {957, 234},
	ElementFightSelectItem5:   {807, 289},
	ElementFightSelectItem6:   {857, 289},
	ElementFightSelectItem7:   {907, 289},
	ElementFightSelectItem8:   {957, 289},
	ElementFightStart:         {940, 500}, // 战斗开始按钮，副本和竞技是同一个位置的按钮
	ElementFubenFightStartAck: {413, 339},
	ElementBackAndExit:        {965, 570},
	ElementFubenHall:          {869, 501}, // 进入副本大厅的按钮
}

func GetPointByElement(element DDTClickElement) Point {
	point, ok := elementRelativePosition[element]
	if !ok {
		return EmptyPoint
	}
	return point
}

type DDTElementRect struct {
	X, Y int32
	W, H int32
}

var (
	elementRectPassBtn = &DDTElementRect{
		X: 477,
		Y: 159,
		W: 47,
		H: 16,
	}
	elementRectAngle = &DDTElementRect{
		X: 32,
		Y: 558,
		W: 35,
		H: 15,
	}
	elementRectFubenSelectText = &DDTElementRect{
		X: 540,
		Y: 110,
		W: 100,
		H: 20,
	}
	elementRectFubenInviteAndChangeTeam = &DDTElementRect{
		X: 761,
		Y: 453,
		W: 125,
		H: 65,
	}
	elementRectJinjiInviteAndChangeArea = &DDTElementRect{
		X: 761,
		Y: 453,
		W: 125,
		H: 65,
	}
	elementRectFubenHall = &DDTElementRect{
		X: 30,
		Y: 83,
		W: 58,
		H: 16,
	}
	elementRectJinjiHall = &DDTElementRect{
		X: 25,
		Y: 60,
		W: 188,
		H: 25,
	}
	elementRectFightRightTop = &DDTElementRect{
		X: 946,
		Y: 4,
		W: 50,
		H: 16,
	}
	elementRectFightResult = &DDTElementRect{
		X: 678,
		Y: 34,
		W: 125,
		H: 30,
	}
	elementRectFightLoading = &DDTElementRect{
		X: 396,
		Y: 301,
		W: 190,
		H: 62,
	}
	elementRectFubenFightSettle = &DDTElementRect{
		X: 159,
		Y: 17,
		W: 149,
		H: 38,
	}
	elementRectJinjiFightSettle = &DDTElementRect{
		X: 749,
		Y: 32,
		W: 149,
		H: 38,
	}
	elementRectBackAndExit = &DDTElementRect{
		X: 945,
		Y: 580,
		W: 38,
		H: 18,
	}
	elementRectMiniMap = &DDTElementRect{
		X: 788,
		Y: 24,
		W: 211,
		H: 97,
	}
	elementRectIsYourTurn = &DDTElementRect{
		X: 460,
		Y: 225,
		W: 100,
		H: 25,
	}
	elementRectWinOrFail = &DDTElementRect{
		X: 821,
		Y: 13,
		W: 155,
		H: 104,
	}
	elementRectFightReady = &DDTElementRect{
		X: 893,
		Y: 507,
		W: 90,
		H: 35,
	}
	elementRectFubenBtn1 = &DDTElementRect{
		X: 219,
		Y: 332,
		W: 130,
		H: 49,
	}
	elementRectFubenBtn2 = &DDTElementRect{
		X: elementRectFubenBtn1.X + 130 + 6,
		Y: 332,
		W: 130,
		H: 49,
	}
	elementRectFubenBtn3 = &DDTElementRect{
		X: elementRectFubenBtn2.X + 130 + 6,
		Y: 332,
		W: 130,
		H: 49,
	}
	elementRectFubenBtn4 = &DDTElementRect{
		X: elementRectFubenBtn3.X + 130 + 6,
		Y: 332,
		W: 130,
		H: 49,
	}
	elementRectFubenBtn5 = &DDTElementRect{
		X: 219,
		Y: 332 + 49 + 15,
		W: 130,
		H: 49,
	}
	elementRectFubenBtn6 = &DDTElementRect{
		X: elementRectFubenBtn5.X + 130 + 6,
		Y: elementRectFubenBtn5.Y,
		W: 130,
		H: 49,
	}
	elementRectFubenBtn7 = &DDTElementRect{
		X: elementRectFubenBtn6.X + 130 + 6,
		Y: elementRectFubenBtn5.Y,
		W: 130,
		H: 49,
	}
	elementRectFubenBtn8 = &DDTElementRect{
		X: elementRectFubenBtn7.X + 130 + 6,
		Y: elementRectFubenBtn5.Y,
		W: 130,
		H: 49,
	}
	elementRectFubenLevelBtn1 = &DDTElementRect{
		X: 437,
		Y: 489,
		W: 110,
		H: 30,
	}
	elementRectFubenLevelBtn2 = &DDTElementRect{
		X: 325,
		Y: 489,
		W: 110,
		H: 30,
	}
	elementRectFubenLevelBtn3 = &DDTElementRect{
		X: 260,
		Y: 489,
		W: 110,
		H: 30,
	}
	elementRectFubenLevelBtn4 = &DDTElementRect{
		X: 232,
		Y: 489,
		W: 110,
		H: 30,
	}
	elementRectFubenLevelBtn5 = &DDTElementRect{
		X: 565,
		Y: 489,
		W: 110,
		H: 30,
	}
	elementRectFubenLevelBtn6 = &DDTElementRect{
		X: 440,
		Y: 489,
		W: 110,
		H: 30,
	}
	elementRectFubenLevelBtn7 = &DDTElementRect{
		X: 366,
		Y: 489,
		W: 110,
		H: 30,
	}
	elementRectFubenLevelBtn8 = &DDTElementRect{
		X: 620,
		Y: 489,
		W: 110,
		H: 30,
	}
	elementRectFubenLevelBtn9 = &DDTElementRect{
		X: 496,
		Y: 489,
		W: 110,
		H: 30,
	}
	elementRectFubenLevelBtn10 = &DDTElementRect{
		X: 632,
		Y: 489,
		W: 110,
		H: 30,
	}
)

// fubenLevelElementRect 副本关卡难度元素所在的可能位置
var fubenLevelElementRect = []*DDTElementRect{
	elementRectFubenLevelBtn1,
	elementRectFubenLevelBtn2,
	elementRectFubenLevelBtn3,
	elementRectFubenLevelBtn4,
	elementRectFubenLevelBtn5,
	elementRectFubenLevelBtn6,
	elementRectFubenLevelBtn7,
	elementRectFubenLevelBtn8,
	elementRectFubenLevelBtn9,
	elementRectFubenLevelBtn10,
}

func RangeFubenLevelElementRect(cb func(*DDTElementRect) bool) {
	for _, rect := range fubenLevelElementRect {
		if rect == nil {
			continue
		}
		if cb(rect) {
			break
		}
	}
}

type DDTElementRectType int

const (
	DDTElementRectTypeFubenBtn1 DDTElementRectType = 1
	DDTElementRectTypeFubenBtn2 DDTElementRectType = 2
	DDTElementRectTypeFubenBtn3 DDTElementRectType = 3
	DDTElementRectTypeFubenBtn4 DDTElementRectType = 4
	DDTElementRectTypeFubenBtn5 DDTElementRectType = 5
	DDTElementRectTypeFubenBtn6 DDTElementRectType = 6
	DDTElementRectTypeFubenBtn7 DDTElementRectType = 7
	DDTElementRectTypeFubenBtn8 DDTElementRectType = 8

	DDTElementRectTypePassBtn         DDTElementRectType = 101 // 战斗内PASS按钮
	DDTElementRectTypeAngle           DDTElementRectType = 102 // 角度
	DDTElementRectTypeFubenSelectText DDTElementRectType = 103 // 副本选择文字
	DDTElementRectTypeBack            DDTElementRectType = 104 // 返回
	DDTElementRectTypeExit            DDTElementRectType = 105 // 退出
	DDTElementRectTypeMiniMap         DDTElementRectType = 106 // 小地图
	DDTElementRectTypeIsYourTurn      DDTElementRectType = 107 // 轮到你出手了
	DDTElementRectTypeWinOrFail       DDTElementRectType = 108 // 结算画面 胜利-失败
	DDTElementRectTypeFanCardSmall    DDTElementRectType = 109 // 翻牌画面 小关翻牌
	DDTElementRectTypeFanCardBoss     DDTElementRectType = 110 // 翻牌画面 boss翻牌
	DDTElementRectTypeFightReady      DDTElementRectType = 111 // 房间内【已准备-显示取消，未准备-显示准备】

	DDTElementRectTypeSenseMin                  DDTElementRectType = 200
	DDTElementRectTypeFubenInviteAndChangeTeam  DDTElementRectType = 201 // 副本房间的特征元素【邀请&换队】
	DDTElementRectTypeJinjiInviteAndChangeArea1 DDTElementRectType = 202 // 竞技房间的特征元素【邀请&换区】
	DDTElementRectTypeJinjiInviteAndChangeArea2 DDTElementRectType = 203 // 竞技房间的特征元素【邀请&本区】
	DDTElementRectTypeFubenHall                 DDTElementRectType = 204 // 副本大厅的特征元素【筛选副本】
	DDTElementRectTypeJinjiHall                 DDTElementRectType = 205 // 竞技大厅的特征元素【房间列表-所有模式】
	DDTElementRectTypeFightRightTop             DDTElementRectType = 206 // 在竞技战斗或副本内的特征元素【右上角的设置和退出按钮】
	DDTElementRectTypeFightResult               DDTElementRectType = 207 // 战斗结算特征元素【我的成绩，竞技战和副本战都用到、胜利和失败都有】
	DDTElementRectTypeFubenFightLoading         DDTElementRectType = 208 // 副本战斗加载特征元素【副本战】
	DDTElementRectTypeJinjiFightLoading         DDTElementRectType = 209 // 竞技战斗加载特征元素【自由战】
	DDTElementRectTypeFubenFightSettle          DDTElementRectType = 210 // 副本战斗结算特征元素【游戏结算，左上角】也是boss关翻牌画面
	DDTElementRectTypeJinjiFightSettle          DDTElementRectType = 211 // 竞技战斗结算特征元素【游戏结算，右上角】也是小关翻牌画面
	DDTElementRectTypeSenseMan                  DDTElementRectType = 300
)

func RangeSenseRect(cb func(DDTElementRectType) bool) {
	for i := DDTElementRectTypeSenseMin + 1; i < DDTElementRectTypeSenseMan; i++ {
		if cb(i) {
			break
		}
	}
}

func RangeFubenBtn(cb func(DDTElementRectType) bool) {
	for i := DDTElementRectTypeFubenBtn1; i <= DDTElementRectTypeFubenBtn8; i++ {
		if cb(i) {
			break
		}
	}
}

var elementRect = map[DDTElementRectType]*DDTElementRect{
	DDTElementRectTypePassBtn:                   elementRectPassBtn,
	DDTElementRectTypeFubenBtn1:                 elementRectFubenBtn1,
	DDTElementRectTypeFubenBtn2:                 elementRectFubenBtn2,
	DDTElementRectTypeFubenBtn3:                 elementRectFubenBtn3,
	DDTElementRectTypeFubenBtn4:                 elementRectFubenBtn4,
	DDTElementRectTypeFubenBtn5:                 elementRectFubenBtn5,
	DDTElementRectTypeFubenBtn6:                 elementRectFubenBtn6,
	DDTElementRectTypeFubenBtn7:                 elementRectFubenBtn7,
	DDTElementRectTypeFubenBtn8:                 elementRectFubenBtn8,
	DDTElementRectTypeAngle:                     elementRectAngle,
	DDTElementRectTypeFubenSelectText:           elementRectFubenSelectText,
	DDTElementRectTypeBack:                      elementRectBackAndExit,
	DDTElementRectTypeExit:                      elementRectBackAndExit,
	DDTElementRectTypeMiniMap:                   elementRectMiniMap,
	DDTElementRectTypeIsYourTurn:                elementRectIsYourTurn,
	DDTElementRectTypeWinOrFail:                 elementRectWinOrFail,
	DDTElementRectTypeFightReady:                elementRectFightReady,
	DDTElementRectTypeFubenInviteAndChangeTeam:  elementRectFubenInviteAndChangeTeam,
	DDTElementRectTypeJinjiInviteAndChangeArea1: elementRectJinjiInviteAndChangeArea,
	DDTElementRectTypeJinjiInviteAndChangeArea2: elementRectJinjiInviteAndChangeArea,
	DDTElementRectTypeFubenHall:                 elementRectFubenHall,
	DDTElementRectTypeJinjiHall:                 elementRectJinjiHall,
	DDTElementRectTypeFightRightTop:             elementRectFightRightTop,
	DDTElementRectTypeFightResult:               elementRectFightResult,
	DDTElementRectTypeFubenFightLoading:         elementRectFightLoading,
	DDTElementRectTypeJinjiFightLoading:         elementRectFightLoading,
	DDTElementRectTypeFubenFightSettle:          elementRectFubenFightSettle,
	DDTElementRectTypeJinjiFightSettle:          elementRectJinjiFightSettle,
}

func GetElementRect(tp DDTElementRectType) *win.RECT {
	rectTmp, ok := elementRect[tp]
	if !ok {
		return nil
	}
	return NewElementRectWithDDTRect(rectTmp)
}

func NewElementRectWithDDTRect(rect *DDTElementRect) *win.RECT {
	if rect == nil {
		return nil
	}
	return &win.RECT{
		Left:   rect.X,
		Top:    rect.Y,
		Right:  rect.X + rect.W,
		Bottom: rect.Y + rect.H,
	}
}

const (
	ImgSimilarityThresholdPassBtn         = 1000 // 相似性阈值 PASS
	ImgSimilarityThresholdIsYourTurn      = 1000 // 相似性阈值 你的回合了
	ImgSimilarityThresholdFubenLevel      = 1000 // 相似性阈值 FubenLevel
	ImgSimilarityThresholdAngle           = 500  // 相似性阈值 角度 这个不能太大，不然会导致识别错误
	ImgSimilarityThresholdFubenSelectText = 1000 // 相似性阈值 副本选择文字【选择副本】
	ImgSimilarityThresholdBack            = 1000 // 相似性阈值 返回按钮
	ImgSimilarityThresholdSence           = 1000 // 相似性阈值 场景判断
	ImgSimilarityThresholdWinOrFail       = 1000 // 相似性阈值 胜利-失败
)

type FubenID int

const (
	FubenIDBegin  FubenID = 0
	FubenIDMY     FubenID = 1  // 蚂蚁
	FubenIDXJ     FubenID = 2  // 小鸡
	FubenIDBG     FubenID = 3  // 波谷
	FubenIDXS     FubenID = 4  // 邪神
	FubenIDBL     FubenID = 5  // 堡垒
	FubenIDLC     FubenID = 6  // 龙巢
	FubenIDYDH    FubenID = 7  // 运动会
	FubenIDJJC    FubenID = 8  // 竞技场
	FubenIDHDMD   FubenID = 9  // 海盗迷岛
	FubenIDMDSC   FubenID = 10 // 迷岛深处
	FubenIDXSZY   FubenID = 11 // 血色庄园
	FubenIDFCDK   FubenID = 12 // 飞出洞窟
	FubenIDBFXY   FubenID = 13 // 冰峰雪域
	FubenIDYZSL2  FubenID = 14 // 勇者试炼2
	FubenIDYZSL3  FubenID = 15 // 勇者试炼3
	FubenIDYZSL4  FubenID = 16 // 勇者试炼4
	FubenIDYZSL5  FubenID = 17 // 勇者试炼5
	FubenIDYZSL6  FubenID = 18 // 勇者试炼6
	FubenIDWSMY   FubenID = 19 // 万圣墓园
	FubenIDSHMD   FubenID = 20 // 守护魔豆
	FubenIDJSWC   FubenID = 21 // 僵尸围城
	FubenIDXSZYTB FubenID = 22 // 万圣墓园自定义的团本？这个应该是私服自定义的
	FubenIDEnd    FubenID = 1000

	FubenIDJinjiBegin FubenID = 1000 // 竞技类战斗ID开始
	FubenIDJinjiEnd   FubenID = 2000 // 竞技类战斗ID结束

	FuBenIDOtherBegin FubenID = 2000 // 其他功能类脚本ID开始
	FuBenIDOtherEnd   FubenID = 3000 // 其他功能类脚本ID开始
)

type FubenType int

const (
	FubenTypeNormal  FubenType = 0 // 普通副本
	FubenTypeSpecial FubenType = 1 // 特殊副本
)

type FubenPosition struct {
	Type  FubenType
	Page  int
	Count int
}

type FubenLevel int

const (
	FubenLevelEasy       FubenLevel = 1 // 简单
	FubenLevelNormal     FubenLevel = 2 // 普通
	FubenLevelDifficulty FubenLevel = 3 // 困难
	FubenLevelHero       FubenLevel = 4 // 英雄
)

func GetFubenLevelElementImg(tp FubenLevel) *image.Gray {
	switch tp {
	case FubenLevelEasy:
		return ElementImgFubenLevelEasy
	case FubenLevelNormal:
		return ElementImgFubenLevelNormal
	case FubenLevelDifficulty:
		return ElementImgFubenLevelDifficulty
	case FubenLevelHero:
		return ElementImgFubenLevelHero
	}
	return nil
}

const (
	ClickWaitLong  = time.Millisecond * 2000 // 长久等待，比如切换场景，从大厅->副本房间内
	ClickWaitMid   = time.Millisecond * 700  // 中等等待，点击后的停顿，这个时间刚好够人类反应了
	ClickWaitShort = time.Millisecond * 150  // 短暂等待，比如准备副本道具
)

type Direction int

const (
	DirectionLeft  Direction = 1
	DirectionRight Direction = 2
	DirectionUp    Direction = 3
	DirectionDown  Direction = 4
)

type SenseType int

const (
	SenseTypeOther             SenseType = 0  // 未归类【战斗外的其他场景】
	SenseTypeFubenRoom         SenseType = 1  // 副本房间场景
	SenseTypeJinjiRoom         SenseType = 2  // 竞技房间场景
	SenseTypeFubenHall         SenseType = 3  // 副本大厅场景
	SenseTypeJinjiHall         SenseType = 4  // 竞技大厅场景
	SenseTypeInFight           SenseType = 5  // 在竞技战斗或副本内
	SenseTypeFightResult       SenseType = 6  // 在竞技战斗或副本的胜负判断界面
	SenseTypeFubenFightLoading SenseType = 7  // 副本加载场景
	SenseTypeJinjiFightLoading SenseType = 8  // 竞技加载场景
	SenseTypeFubenFightSettle  SenseType = 9  // 副本翻牌结算场景
	SenseTypeJinjiFightSettle  SenseType = 10 // 竞技翻牌结算场景
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

type FightInnerPosition int

const (
	FightInnerPosition1 FightInnerPosition = 1
	FightInnerPosition2 FightInnerPosition = 2
	FightInnerPosition3 FightInnerPosition = 3
	FightInnerPosition4 FightInnerPosition = 4
)

type ReadyState int

const (
	ReadyStateNo ReadyState = 0 // 未准备
	ReadyStateOK ReadyState = 1 // 已准备
)
