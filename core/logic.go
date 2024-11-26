package core

import (
	"ddt-copilot/data"
	"ddt-copilot/defs"
	"ddt-copilot/utils"
	"github.com/lxn/win"
)

var (
	defaultAttackList = []string{"2", "3",
		"4", "4", "4", "4", "4", "4", "4", "4",
		"5", "5", "5", "5", "5", "5", "5", "5",
		"6", "6", "6", "6", "6", "6", "6", "6",
		"7", "7", "7", "7", "7", "7", "7", "7",
		"8", "8", "8", "8", "8", "8", "8", "8"} // 默认的攻击指令
)

func GetAttackCmd() []uintptr {
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
	cmds := GetAttackCmd()
	for _, cmd := range cmds {
		utils.UseSkill(hwnd, cmd)
	}
}

func Launch(hwnd win.HWND, needAngle, power int) {
	num, err := data.GDefsAngle.GetAngle(hwnd)
	if err != nil {
		data.Log().Error().Err(err).Int("hwnd", int(hwnd)).Msg("Launch, GetAngle failed")
	}
	utils.UpdateAngle(hwnd, needAngle-num)
	utils.Launch(hwnd, power)
}
