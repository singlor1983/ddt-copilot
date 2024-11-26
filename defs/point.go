package defs

type Point struct {
	X int
	Y int
}

var (
	IllegalPoint = Point{-1, -1}
	EmptyPoint   = Point{1, 1}
)

type ElementPoint int

const (
	PointIllegal            ElementPoint = 0  // 非法坐标点
	PointEmpty              ElementPoint = 1  // 空白处，用于激活窗口
	PointFubenEnter         ElementPoint = 2  // 从副本大厅进入副本房间
	PointFubenSelect        ElementPoint = 3  // 副本选择按钮
	PointFubenPageDown      ElementPoint = 4  // 副本下拉按钮
	PointFubenTypeNormal    ElementPoint = 5  // 普通副本
	PointFubenTypeSpecial   ElementPoint = 6  // 特殊副本
	PointFubenBossFight     ElementPoint = 7  // 副本BOSS战按钮
	PointFubenSelectAck     ElementPoint = 8  // 副本选择确认按钮
	PointFubenBossFightAck  ElementPoint = 9  // 副本BOSS战确认按钮
	PointFightEquipItem1    ElementPoint = 10 // 房间内战斗已装备道具1
	PointFightEquipItem2    ElementPoint = 11 // 房间内战斗已装备道具2
	PointFightEquipItem3    ElementPoint = 12 // 房间内战斗已装备道具3
	PointFightSelectItem1   ElementPoint = 13 // 房间内供选择道具1
	PointFightSelectItem2   ElementPoint = 14 // 房间内供选择道具2
	PointFightSelectItem3   ElementPoint = 15 // 房间内供选择道具3
	PointFightSelectItem4   ElementPoint = 16 // 房间内供选择道具4
	PointFightSelectItem5   ElementPoint = 17 // 房间内供选择道具5
	PointFightSelectItem6   ElementPoint = 18 // 房间内供选择道具6
	PointFightSelectItem7   ElementPoint = 19 // 房间内供选择道具7
	PointFightSelectItem8   ElementPoint = 20 // 房间内供选择道具8
	PointFightStart         ElementPoint = 21 // 战斗开始按钮，副本和竞技是同一个位置的按钮
	PointFubenFightStartAck ElementPoint = 22 // 副本战斗开始确认按钮
	PointBackAndExit        ElementPoint = 23 // 返回&&退出按钮
	PointFubenHall          ElementPoint = 24 // 进入副本大厅的按钮
)

// pointMap 可点击元素中心相对于左上角顶点的坐标
var pointMap = map[ElementPoint]Point{
	PointIllegal:            IllegalPoint,
	PointEmpty:              EmptyPoint,
	PointFubenEnter:         {722, 479},
	PointFubenSelect:        {597, 234},
	PointFubenPageDown:      {768, 440},
	PointFubenTypeNormal:    {287, 305},
	PointFubenTypeSpecial:   {427, 305},
	PointFubenBossFight:     {320, 565},
	PointFubenSelectAck:     {500, 565},
	PointFubenBossFightAck:  {433, 340},
	PointFightEquipItem1:    {810, 140},
	PointFightEquipItem2:    {880, 140},
	PointFightEquipItem3:    {950, 140},
	PointFightSelectItem1:   {807, 234},
	PointFightSelectItem2:   {857, 234},
	PointFightSelectItem3:   {907, 234},
	PointFightSelectItem4:   {957, 234},
	PointFightSelectItem5:   {807, 289},
	PointFightSelectItem6:   {857, 289},
	PointFightSelectItem7:   {907, 289},
	PointFightSelectItem8:   {957, 289},
	PointFightStart:         {940, 500},
	PointFubenFightStartAck: {413, 339},
	PointBackAndExit:        {965, 570},
	PointFubenHall:          {869, 501},
}

func GetPoint(element ElementPoint) Point {
	point, ok := pointMap[element]
	if !ok {
		return IllegalPoint
	}
	return point
}
