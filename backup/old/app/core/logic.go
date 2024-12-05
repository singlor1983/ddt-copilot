package core

import (
	"ddt-copilot/data"
	"ddt-copilot/defs"
	"ddt-copilot/utils"
	"github.com/lxn/win"
)

var (
	defaultAttackList = []string{"E", "2", "3",
		"4", "4", "4", "4", "4", "4", "4", "4",
		"5", "5", "5", "5", "5", "5", "5", "5",
		"6", "6", "6", "6", "6", "6", "6", "6",
		"7", "7", "7", "7", "7", "7", "7", "7",
		"8", "8", "8", "8", "8", "8", "8", "8"} // 默认的攻击指令
)

func GetAttackCmd() []uintptr {
	var cmds []uintptr
	attackCmd := data.GConfig.GetAttackCmd()
	if len(attackCmd) == 0 {
		attackCmd = defaultAttackList
	}
	for _, s := range attackCmd {
		vk := defs.GetVkFromStr(s)
		cmds = append(cmds, vk)
	}
	return cmds
}

func UseSkillByConfig(hwnd win.HWND) {
	cmds := GetAttackCmd()
	for _, cmd := range cmds {
		utils.UseSkill(hwnd, cmd)
	}
}

func Launch(hwnd win.HWND, needAngle, power int) {
	num := utils.GetAngle(hwnd)
	utils.UpdateAngle(hwnd, needAngle-num)
	utils.Launch(hwnd, power)
}

// 战前道具选择，蚂蚁无需选择，只是示例
func OnBeforeStartMY(ctrl *ScriptCtrl) {
	//utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightEquipItem1)
	//utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightEquipItem2)
	//utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightEquipItem3)
	//utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightSelectItem4)
	//utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightSelectItem4)
	//utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightSelectItem4)
}

// isMYBossFight 第二关或者boss战则为直接打boss
func isMYBossFight(ctrl *ScriptCtrl) bool {
	if ctrl.isBossFight {
		return true
	}
	if ctrl.fixFightCount%2 == 0 { // 对于蚂蚁副本来说，非boss战情况下第二关为boss战，战斗失败会重新从第一关开始打，所以需要重置
		return true
	}
	return false
}

func OnFirstRoundInitMY(ctrl *ScriptCtrl) {
	if isMYBossFight(ctrl) {
		positions := map[int]defs.FightInnerPosition{
			6:  defs.FightInnerPosition1,
			10: defs.FightInnerPosition2,
			41: defs.FightInnerPosition3,
			-5: defs.FightInnerPosition4,
		}
		ctrl.initPosition = positions[utils.GetAngle(ctrl.hwnd)]
		if ctrl.initPosition == defs.FightInnerPosition3 { // 蚂蚁3号位置角度要调整位置
			utils.Move(ctrl.hwnd, defs.DirectionLeft, 5)
			utils.ConfirmDirection(ctrl.hwnd, defs.DirectionRight)
		}
	} else {
		positions := map[int]defs.FightInnerPosition{
			27: defs.FightInnerPosition1,
			16: defs.FightInnerPosition2,
			10: defs.FightInnerPosition3,
			15: defs.FightInnerPosition4,
		}
		ctrl.initPosition = positions[utils.GetAngle(ctrl.hwnd)]
		switch ctrl.initPosition { // 全部移动到同一个位置去，避免打到队友 Todo 待验证
		case defs.FightInnerPosition1:
			utils.UseSkill(ctrl.hwnd, defs.VK_P)
		case defs.FightInnerPosition2:
			utils.Move(ctrl.hwnd, defs.DirectionLeft, 5)
			utils.UseSkill(ctrl.hwnd, defs.VK_P)
		case defs.FightInnerPosition3:
			utils.Move(ctrl.hwnd, defs.DirectionLeft, 10)
			utils.UseSkill(ctrl.hwnd, defs.VK_P)
		case defs.FightInnerPosition4:
			utils.Move(ctrl.hwnd, defs.DirectionLeft, 15)
			utils.UseSkill(ctrl.hwnd, defs.VK_P)
		}
	}
}

func OnFightMY(ctrl *ScriptCtrl) {
	if ctrl.isChild {
		utils.UseSkill(ctrl.hwnd, defs.VK_P)
		return
	}
	if isMYBossFight(ctrl) {
		UseSkillByConfig(ctrl.hwnd)
		Launch(ctrl.hwnd, 30, 75)
	} else { // 这就是蚂蚁的第一关小关
		if ctrl.roundCount%2 == 0 {
			utils.ConfirmDirection(ctrl.hwnd, defs.DirectionRight)
		} else {
			utils.ConfirmDirection(ctrl.hwnd, defs.DirectionLeft)
		}
		Launch(ctrl.hwnd, 20, 10)
	}
}

func OnFightXSZY(ctrl *ScriptCtrl) {
	if ctrl.isChild {
		utils.UseSkill(ctrl.hwnd, defs.VK_P)
		return
	}
	if isMYBossFight(ctrl) {
		// 第二关的位置初始力度 15 16 25 8
		var distance int
		num := utils.GetAngle(ctrl.hwnd) // Todo zhangzhihi 不能在这里获取角度，只有在第一回合取一次就行了，然后保存起来，不然移动之后第二回合的角度已经改变了
		if num == 15 {
			distance = 15
		} else if num == 16 {
			distance = 20
		} else if num == 25 {
			distance = 30
		} else if num == 25 {
			distance = 40
		}
		if ctrl.roundCount == 1 {
			utils.Move(ctrl.hwnd, defs.DirectionLeft, distance) // Todo zhangzhihui 看来这里要根据不同的位置计算不同的移动距离了.如何判断不同的位置？根据角度还是取红点来？
			utils.UseSkill(ctrl.hwnd, defs.VK_P)
			return
		} else if ctrl.roundCount == 2 {
			utils.Move(ctrl.hwnd, defs.DirectionLeft, distance)
			utils.ConfirmDirection(ctrl.hwnd, defs.DirectionRight)
			utils.UseSkill(ctrl.hwnd, defs.VK_F) // F
			Launch(ctrl.hwnd, 70, 62)
		} else if ctrl.roundCount == 3 {
			utils.Move(ctrl.hwnd, defs.DirectionLeft, 25)
			utils.ConfirmDirection(ctrl.hwnd, defs.DirectionRight)
			utils.UseSkill(ctrl.hwnd, defs.VK_P)
			return
		} else {
			utils.UseSkill(ctrl.hwnd, defs.VK_2) // 2
			utils.UseSkill(ctrl.hwnd, defs.VK_3) // 3
			utils.UseSkill(ctrl.hwnd, defs.VK_4) // 4
			utils.UseSkill(ctrl.hwnd, defs.VK_4) // 4
			utils.UseSkill(ctrl.hwnd, defs.VK_4) // 4
			utils.UseSkill(ctrl.hwnd, defs.VK_4) // 4
			utils.UseSkill(ctrl.hwnd, defs.VK_5) // 5
			utils.UseSkill(ctrl.hwnd, defs.VK_5) // 5
			utils.UseSkill(ctrl.hwnd, defs.VK_5) // 5
			utils.UseSkill(ctrl.hwnd, defs.VK_5) // 5
			utils.UseSkill(ctrl.hwnd, defs.VK_6) // 6
			utils.UseSkill(ctrl.hwnd, defs.VK_6) // 6
			utils.UseSkill(ctrl.hwnd, defs.VK_6) // 6
			utils.UseSkill(ctrl.hwnd, defs.VK_6) // 6
			utils.UseSkill(ctrl.hwnd, defs.VK_7) // 7
			utils.UseSkill(ctrl.hwnd, defs.VK_7) // 7
			utils.UseSkill(ctrl.hwnd, defs.VK_7) // 7
			utils.UseSkill(ctrl.hwnd, defs.VK_7) // 7
			Launch(ctrl.hwnd, 45, 40)
		}
	} else {
		// 所有人的初始位置与初始力度一样，力度均为16
		if ctrl.roundCount == 1 { // 第一回合飞过去
			utils.ConfirmDirection(ctrl.hwnd, defs.DirectionRight)
			utils.UseSkill(ctrl.hwnd, defs.VK_F) // F
			Launch(ctrl.hwnd, 35, 96)
		} else if ctrl.roundCount == 2 {
			utils.Move(ctrl.hwnd, defs.DirectionLeft, 26)
			utils.ConfirmDirection(ctrl.hwnd, defs.DirectionRight)
			utils.UseSkill(ctrl.hwnd, defs.VK_P)
		} else {
			utils.UseSkill(ctrl.hwnd, defs.VK_2) // 2
			utils.UseSkill(ctrl.hwnd, defs.VK_3) // 3
			utils.UseSkill(ctrl.hwnd, defs.VK_4) // 4
			utils.UseSkill(ctrl.hwnd, defs.VK_4) // 4
			utils.UseSkill(ctrl.hwnd, defs.VK_4) // 4
			utils.UseSkill(ctrl.hwnd, defs.VK_4) // 4
			utils.UseSkill(ctrl.hwnd, defs.VK_5) // 5
			utils.UseSkill(ctrl.hwnd, defs.VK_5) // 5
			utils.UseSkill(ctrl.hwnd, defs.VK_5) // 5
			utils.UseSkill(ctrl.hwnd, defs.VK_5) // 5
			utils.UseSkill(ctrl.hwnd, defs.VK_6) // 6
			utils.UseSkill(ctrl.hwnd, defs.VK_6) // 6
			utils.UseSkill(ctrl.hwnd, defs.VK_6) // 6
			utils.UseSkill(ctrl.hwnd, defs.VK_6) // 6
			utils.UseSkill(ctrl.hwnd, defs.VK_7) // 7
			utils.UseSkill(ctrl.hwnd, defs.VK_7) // 7
			utils.UseSkill(ctrl.hwnd, defs.VK_7) // 7
			utils.UseSkill(ctrl.hwnd, defs.VK_7) // 7
			Launch(ctrl.hwnd, 60, 40)
		}
	}
}

func OnFirstRoundInitBG(ctrl *ScriptCtrl) {
	positions := map[int]defs.FightInnerPosition{
		17: defs.FightInnerPosition1,
		34: defs.FightInnerPosition2,
		29: defs.FightInnerPosition3,
		43: defs.FightInnerPosition4,
	}
	ctrl.initPosition = positions[utils.GetAngle(ctrl.hwnd)]
	UseSkillByConfig(ctrl.hwnd)
	switch ctrl.initPosition {
	case defs.FightInnerPosition1:
		Launch(ctrl.hwnd, 30, 55)
	case defs.FightInnerPosition2:
		Launch(ctrl.hwnd, 34, 45)
	case defs.FightInnerPosition3:
		Launch(ctrl.hwnd, 29, 45)
	case defs.FightInnerPosition4:
		Launch(ctrl.hwnd, 43, 35)
	}
}

func OnFightBG(ctrl *ScriptCtrl) {
	UseSkillByConfig(ctrl.hwnd)
	switch ctrl.initPosition {
	case defs.FightInnerPosition1:
		Launch(ctrl.hwnd, 30, 55)
	case defs.FightInnerPosition2:
		Launch(ctrl.hwnd, 34, 45)
	case defs.FightInnerPosition3:
		Launch(ctrl.hwnd, 29, 45)
	case defs.FightInnerPosition4:
		Launch(ctrl.hwnd, 43, 35)
	}
}
