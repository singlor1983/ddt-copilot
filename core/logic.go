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

func GetAttackCMD() []uintptr {
	var cmds []uintptr
	attackCmd := data.GGameSetting.SettingGeneral.AttackCMD
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
	cmds := GetAttackCMD()
	for _, cmd := range cmds {
		utils.UseSkill(hwnd, cmd)
	}
}

func Launch(hwnd win.HWND, needAngle, power int) {
	num, err := data.GDefsAngle.GetAngle(hwnd)
	if err != nil {
		data.Log().Error().Timestamp().Timestamp().Err(err).Int("hwnd", int(hwnd)).Msg("launch failed")
	}
	diff := needAngle - num
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Int("captureNum", num).Int("diff", diff).
		Int("needAngle", needAngle).Int("power", power).Msg("launch")
	utils.UpdateAngle(hwnd, diff)
	utils.Launch(hwnd, power)
}

func OnBattleCustom(ctrl *ScriptCtrl) {
	utils.ConfirmDirection(ctrl.hwnd, data.GGameSetting.SettingFubenCustom.Direction)

	attackCmd := data.GGameSetting.SettingFubenCustom.AttackCMD
	if len(attackCmd) == 0 {
		attackCmd = defaultAttackList
	}
	for _, s := range attackCmd {
		vk := defs.GetVkFromStr(s)
		utils.UseSkill(ctrl.hwnd, vk)
	}

	Launch(ctrl.hwnd, data.GGameSetting.SettingFubenCustom.Angle, data.GGameSetting.SettingFubenCustom.Power)
}

func OnRoundInitMaYiGeneral(ctrl *ScriptCtrl) {

}

func OnBattleMaYiGeneral(ctrl *ScriptCtrl) {

}
