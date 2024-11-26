package defs

import "github.com/lxn/win"

type Rect struct {
	X, Y int32
	W, H int32
}

var (
	RectPassBtn = &Rect{
		X: 477,
		Y: 159,
		W: 47,
		H: 16,
	}
	RectAngle = &Rect{
		X: 32,
		Y: 558,
		W: 35,
		H: 15,
	}
	RectFubenSelectText = &Rect{
		X: 540,
		Y: 110,
		W: 100,
		H: 20,
	}
	RectFubenInviteAndChangeTeam = &Rect{
		X: 761,
		Y: 453,
		W: 125,
		H: 65,
	}
	RectJinjiInviteAndChangeArea = &Rect{
		X: 761,
		Y: 453,
		W: 125,
		H: 65,
	}
	RectFubenHall = &Rect{
		X: 30,
		Y: 83,
		W: 58,
		H: 16,
	}
	RectJinjiHall = &Rect{
		X: 25,
		Y: 60,
		W: 188,
		H: 25,
	}
	RectFightRightTop = &Rect{
		X: 946,
		Y: 4,
		W: 50,
		H: 16,
	}
	RectFightResult = &Rect{
		X: 678,
		Y: 34,
		W: 125,
		H: 30,
	}
	RectFightLoading = &Rect{
		X: 396,
		Y: 301,
		W: 190,
		H: 62,
	}
	RectFubenFightSettle = &Rect{
		X: 159,
		Y: 17,
		W: 149,
		H: 38,
	}
	RectJinjiFightSettle = &Rect{
		X: 749,
		Y: 32,
		W: 149,
		H: 38,
	}
	RectBackAndExit = &Rect{
		X: 945,
		Y: 580,
		W: 38,
		H: 18,
	}
	RectMiniMap = &Rect{
		X: 788,
		Y: 24,
		W: 211,
		H: 97,
	}
	RectIsYourTurn = &Rect{
		X: 460,
		Y: 225,
		W: 100,
		H: 25,
	}
	RectWinOrFail = &Rect{
		X: 821,
		Y: 13,
		W: 155,
		H: 104,
	}
	RectFightReady = &Rect{
		X: 893,
		Y: 507,
		W: 90,
		H: 35,
	}
	RectFubenBtn1 = &Rect{
		X: 219,
		Y: 332,
		W: 130,
		H: 49,
	}
	RectFubenBtn2 = &Rect{
		X: RectFubenBtn1.X + 130 + 6,
		Y: 332,
		W: 130,
		H: 49,
	}
	RectFubenBtn3 = &Rect{
		X: RectFubenBtn2.X + 130 + 6,
		Y: 332,
		W: 130,
		H: 49,
	}
	RectFubenBtn4 = &Rect{
		X: RectFubenBtn3.X + 130 + 6,
		Y: 332,
		W: 130,
		H: 49,
	}
	RectFubenBtn5 = &Rect{
		X: 219,
		Y: 332 + 49 + 15,
		W: 130,
		H: 49,
	}
	RectFubenBtn6 = &Rect{
		X: RectFubenBtn5.X + 130 + 6,
		Y: RectFubenBtn5.Y,
		W: 130,
		H: 49,
	}
	RectFubenBtn7 = &Rect{
		X: RectFubenBtn6.X + 130 + 6,
		Y: RectFubenBtn5.Y,
		W: 130,
		H: 49,
	}
	RectFubenBtn8 = &Rect{
		X: RectFubenBtn7.X + 130 + 6,
		Y: RectFubenBtn5.Y,
		W: 130,
		H: 49,
	}
	RectFubenLevelBtn1 = &Rect{
		X: 437,
		Y: 489,
		W: 110,
		H: 30,
	}
	RectFubenLevelBtn2 = &Rect{
		X: 325,
		Y: 489,
		W: 110,
		H: 30,
	}
	RectFubenLevelBtn3 = &Rect{
		X: 260,
		Y: 489,
		W: 110,
		H: 30,
	}
	RectFubenLevelBtn4 = &Rect{
		X: 232,
		Y: 489,
		W: 110,
		H: 30,
	}
	RectFubenLevelBtn5 = &Rect{
		X: 565,
		Y: 489,
		W: 110,
		H: 30,
	}
	RectFubenLevelBtn6 = &Rect{
		X: 440,
		Y: 489,
		W: 110,
		H: 30,
	}
	RectFubenLevelBtn7 = &Rect{
		X: 366,
		Y: 489,
		W: 110,
		H: 30,
	}
	RectFubenLevelBtn8 = &Rect{
		X: 620,
		Y: 489,
		W: 110,
		H: 30,
	}
	RectFubenLevelBtn9 = &Rect{
		X: 496,
		Y: 489,
		W: 110,
		H: 30,
	}
	RectFubenLevelBtn10 = &Rect{
		X: 632,
		Y: 489,
		W: 110,
		H: 30,
	}
)

// fubenLevelRect 副本关卡难度元素所在的可能位置
var fubenLevelRect = []*Rect{
	RectFubenLevelBtn1,
	RectFubenLevelBtn2,
	RectFubenLevelBtn3,
	RectFubenLevelBtn4,
	RectFubenLevelBtn5,
	RectFubenLevelBtn6,
	RectFubenLevelBtn7,
	RectFubenLevelBtn8,
	RectFubenLevelBtn9,
	RectFubenLevelBtn10,
}

func RangeFubenLevelRect(cb func(*Rect) bool) {
	for _, rect := range fubenLevelRect {
		if rect == nil {
			continue
		}
		if cb(rect) {
			break
		}
	}
}

type RectType int

const (
	RectTypeFubenBtn1 RectType = 1
	RectTypeFubenBtn2 RectType = 2
	RectTypeFubenBtn3 RectType = 3
	RectTypeFubenBtn4 RectType = 4
	RectTypeFubenBtn5 RectType = 5
	RectTypeFubenBtn6 RectType = 6
	RectTypeFubenBtn7 RectType = 7
	RectTypeFubenBtn8 RectType = 8

	RectTypePassBtn         RectType = 101 // 战斗内PASS按钮
	RectTypeAngle           RectType = 102 // 角度
	RectTypeFubenSelectText RectType = 103 // 副本选择文字
	RectTypeBack            RectType = 104 // 返回
	RectTypeExit            RectType = 105 // 退出
	RectTypeMiniMap         RectType = 106 // 小地图
	RectTypeIsYourTurn      RectType = 107 // 轮到你出手了
	RectTypeWinOrFail       RectType = 108 // 结算画面 胜利-失败
	RectTypeFanCardSmall    RectType = 109 // 翻牌画面 小关翻牌
	RectTypeFanCardBoss     RectType = 110 // 翻牌画面 boss翻牌
	RectTypeFightReady      RectType = 111 // 房间内【已准备-显示取消，未准备-显示准备】

	RectTypeSenseMin                  RectType = 200
	RectTypeFubenInviteAndChangeTeam  RectType = 201 // 副本房间的特征元素【邀请&换队】
	RectTypeJinjiInviteAndChangeArea1 RectType = 202 // 竞技房间的特征元素【邀请&换区】
	RectTypeJinjiInviteAndChangeArea2 RectType = 203 // 竞技房间的特征元素【邀请&本区】
	RectTypeFubenHall                 RectType = 204 // 副本大厅的特征元素【筛选副本】
	RectTypeJinjiHall                 RectType = 205 // 竞技大厅的特征元素【房间列表-所有模式】
	RectTypeFightRightTop             RectType = 206 // 在竞技战斗或副本内的特征元素【右上角的设置和退出按钮】
	RectTypeFightResult               RectType = 207 // 战斗结算特征元素【我的成绩，竞技战和副本战都用到、胜利和失败都有】
	RectTypeFubenFightLoading         RectType = 208 // 副本战斗加载特征元素【副本战】
	RectTypeJinjiFightLoading         RectType = 209 // 竞技战斗加载特征元素【自由战】
	RectTypeFubenFightSettle          RectType = 210 // 副本战斗结算特征元素【游戏结算，左上角】也是boss关翻牌画面
	RectTypeJinjiFightSettle          RectType = 211 // 竞技战斗结算特征元素【游戏结算，右上角】也是小关翻牌画面
	RectTypeSenseMan                  RectType = 300
)

func RangeSenseRect(cb func(RectType) bool) {
	for i := RectTypeSenseMin + 1; i < RectTypeSenseMan; i++ {
		if cb(i) {
			break
		}
	}
}

func RangeFubenBtn(cb func(RectType) bool) {
	for i := RectTypeFubenBtn1; i <= RectTypeFubenBtn8; i++ {
		if cb(i) {
			break
		}
	}
}

var elementRect = map[RectType]*Rect{
	RectTypePassBtn:                   RectPassBtn,
	RectTypeFubenBtn1:                 RectFubenBtn1,
	RectTypeFubenBtn2:                 RectFubenBtn2,
	RectTypeFubenBtn3:                 RectFubenBtn3,
	RectTypeFubenBtn4:                 RectFubenBtn4,
	RectTypeFubenBtn5:                 RectFubenBtn5,
	RectTypeFubenBtn6:                 RectFubenBtn6,
	RectTypeFubenBtn7:                 RectFubenBtn7,
	RectTypeFubenBtn8:                 RectFubenBtn8,
	RectTypeAngle:                     RectAngle,
	RectTypeFubenSelectText:           RectFubenSelectText,
	RectTypeBack:                      RectBackAndExit,
	RectTypeExit:                      RectBackAndExit,
	RectTypeMiniMap:                   RectMiniMap,
	RectTypeIsYourTurn:                RectIsYourTurn,
	RectTypeWinOrFail:                 RectWinOrFail,
	RectTypeFightReady:                RectFightReady,
	RectTypeFubenInviteAndChangeTeam:  RectFubenInviteAndChangeTeam,
	RectTypeJinjiInviteAndChangeArea1: RectJinjiInviteAndChangeArea,
	RectTypeJinjiInviteAndChangeArea2: RectJinjiInviteAndChangeArea,
	RectTypeFubenHall:                 RectFubenHall,
	RectTypeJinjiHall:                 RectJinjiHall,
	RectTypeFightRightTop:             RectFightRightTop,
	RectTypeFightResult:               RectFightResult,
	RectTypeFubenFightLoading:         RectFightLoading,
	RectTypeJinjiFightLoading:         RectFightLoading,
	RectTypeFubenFightSettle:          RectFubenFightSettle,
	RectTypeJinjiFightSettle:          RectJinjiFightSettle,
}

func GetWinRect(tp RectType) *win.RECT {
	rectTmp, ok := elementRect[tp]
	if !ok {
		return nil
	}
	return ToWinRect(rectTmp)
}

func ToWinRect(rect *Rect) *win.RECT {
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

func ToRect(rect *win.RECT) *Rect {
	if rect == nil {
		return nil
	}
	return &Rect{
		X: rect.Left,
		Y: rect.Top,
		W: rect.Right - rect.Left,
		H: rect.Bottom - rect.Top,
	}
}
