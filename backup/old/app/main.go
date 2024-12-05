package main

import (
	"ddt-copilot/app/core"
	"ddt-copilot/data"
	"ddt-copilot/defs"
	"ddt-copilot/utils"
)

func main() {
	data.InitGlobal()
	data.GConfig.SetBossFightEnable(true)
	data.GConfig.SetFubenPosition(data.SettingFubenPostion{Position: map[defs.FubenID]defs.FubenPosition{
		defs.FubenIDMY:     {Type: defs.FubenTypeNormal, Page: 1, Count: 1},
		defs.FubenIDXJ:     {Type: defs.FubenTypeNormal, Page: 1, Count: 2},
		defs.FubenIDBG:     {Type: defs.FubenTypeNormal, Page: 1, Count: 3},
		defs.FubenIDXS:     {Type: defs.FubenTypeNormal, Page: 1, Count: 4},
		defs.FubenIDBL:     {Type: defs.FubenTypeNormal, Page: 1, Count: 5},
		defs.FubenIDLC:     {Type: defs.FubenTypeNormal, Page: 1, Count: 6},
		defs.FubenIDYDH:    {Type: defs.FubenTypeNormal, Page: 1, Count: 7},
		defs.FubenIDJJC:    {Type: defs.FubenTypeNormal, Page: 1, Count: 8},
		defs.FubenIDHDMD:   {Type: defs.FubenTypeNormal, Page: 2, Count: 1},
		defs.FubenIDMDSC:   {Type: defs.FubenTypeNormal, Page: 2, Count: 2},
		defs.FubenIDXSZY:   {Type: defs.FubenTypeNormal, Page: 2, Count: 3},
		defs.FubenIDFCDK:   {Type: defs.FubenTypeNormal, Page: 2, Count: 4},
		defs.FubenIDBFXY:   {Type: defs.FubenTypeNormal, Page: 2, Count: 5},
		defs.FubenIDYZSL2:  {Type: defs.FubenTypeNormal, Page: 2, Count: 6},
		defs.FubenIDYZSL3:  {Type: defs.FubenTypeNormal, Page: 2, Count: 7},
		defs.FubenIDYZSL4:  {Type: defs.FubenTypeNormal, Page: 2, Count: 8},
		defs.FubenIDYZSL5:  {Type: defs.FubenTypeNormal, Page: 3, Count: 1},
		defs.FubenIDYZSL6:  {Type: defs.FubenTypeNormal, Page: 3, Count: 2},
		defs.FubenIDWSMY:   {Type: defs.FubenTypeNormal, Page: 3, Count: 3},
		defs.FubenIDSHMD:   {Type: defs.FubenTypeNormal, Page: 3, Count: 4},
		defs.FubenIDJSWC:   {Type: defs.FubenTypeNormal, Page: 3, Count: 5},
		defs.FubenIDXSZYTB: {Type: defs.FubenTypeNormal, Page: 3, Count: 6},
	}})
	hwnds := utils.GetDDTWindowsHWND()
	if len(hwnds) == 0 {
		return
	}

	core.Init()

	ctrl1 := core.NewScriptCtrl(hwnds[0], 0, defs.FubenIDBG, defs.FubenLevelHero, false, data.GConfig.IsBossFightEnable())
	//ctrl2 := core.NewScriptCtrl(hwnds[1], hwnds[0], defs.FubenIDMY, defs.FubenLevelEasy, true, data.GConfig.IsBossFightEnable())
	//ctrl3 := core.NewScriptCtrl(hwnds[2], hwnds[0], defs.FubenIDMY, defs.FubenLevelEasy, true, data.GConfig.IsBossFightEnable())
	//ctrl4 := core.NewScriptCtrl(hwnds[3], hwnds[0], defs.FubenIDMY, defs.FubenLevelEasy, true, data.GConfig.IsBossFightEnable())

	core.InstanceMgr().AddCtrl(ctrl1)
	//core.InstanceMgr().AddCtrl(ctrl2)
	//core.InstanceMgr().AddCtrl(ctrl3)
	//core.InstanceMgr().AddCtrl(ctrl4)

	//ctrl2.Run()
	//time.Sleep(time.Second) // Todo 并发执行问题？
	//ctrl3.Run()
	//time.Sleep(time.Second)
	//ctrl4.Run()
	//time.Sleep(time.Second)
	ctrl1.Run()

	for {
	}
}
