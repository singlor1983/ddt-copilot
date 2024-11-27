package main

import (
	"ddt-copilot/core"
	"ddt-copilot/data"
	"ddt-copilot/defs"
	"ddt-copilot/utils"
)

func main() {
	data.InitGlobal()

	core.Init()

	hwnds := utils.GetDDTHwnds()
	if len(hwnds) == 0 {
		return
	}

	ctrl1 := core.NewScriptCtrl(hwnds[0], 0, defs.FunctionIDCustomFuben, defs.FubenLvNormal, false, data.GGameSetting.SettingGeneral.IsBossFightEnable)
	core.InstanceMgr().AddCtrl(ctrl1)

	ctrl1.Run()

	for {
	}
}
