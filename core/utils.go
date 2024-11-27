package core

import (
	"ddt-copilot/data"
	"ddt-copilot/defs"
	"ddt-copilot/utils"
	"github.com/lxn/win"
)

func InFubenRoom(hwnd win.HWND) bool {
	rectType := defs.RectTypeFubenInviteAndChangeTeam
	standard := data.GDefsOther.GetStandard(rectType)
	return utils.IsSimilarity(hwnd, standard, rectType, 0.8, true)
}

func InJinRoom(hwnd win.HWND) bool {
	for _, rectType := range []defs.RectType{defs.RectTypeJinjiInviteAndChangeArea1, defs.RectTypeJinjiInviteAndChangeArea2} {
		standard := data.GDefsOther.GetStandard(rectType)
		if utils.IsSimilarity(hwnd, standard, rectType, 0.8, true) {
			return true
		}
	}
	return false
}

func InFubenHall(hwnd win.HWND) bool {
	rectType := defs.RectTypeFubenHall
	standard := data.GDefsOther.GetStandard(rectType)
	return utils.IsSimilarity(hwnd, standard, rectType, 0.8, true)
}

func InJinjiHall(hwnd win.HWND) bool {
	rectType := defs.RectTypeJinjiHall
	standard := data.GDefsOther.GetStandard(rectType)
	return utils.IsSimilarity(hwnd, standard, rectType, 0.8, true)
}

func InIndexPage(hwnd win.HWND) bool {
	rectType := defs.RectTypeIndexPage
	standard := data.GDefsOther.GetStandard(rectType)
	return utils.IsSimilarity(hwnd, standard, rectType, 0.8, true)
}

func IsReady(hwnd win.HWND) bool {
	rectType := defs.RectTypeFightReady
	standard := data.GDefsOther.GetStandard(rectType)
	return utils.IsSimilarity(hwnd, standard, rectType, 0.8, true)
}

func BackToIndexPage(hwnd win.HWND) {
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("back to index page")
	standard := data.GDefsOther.GetStandard(defs.RectTypeBack)
	for {
		rect := defs.GetWinRect(defs.RectTypeBack)
		if rect == nil {
			return
		}
		gray, err := utils.CaptureWindowLightWithGray(hwnd, rect, true)
		if err != nil {
			return
		}
		if !utils.IsImageSimilarity(standard, gray, 0.9) {
			break
		}
		utils.ClickPointByType(hwnd, defs.PointBackAndExit, defs.TimeWaitLong)
	}
}

func EnterFubenHall(hwnd win.HWND) {
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("enter fuben hall")
	utils.ClickPointByType(hwnd, defs.PointFubenHall, defs.TimeWaitLong)
}

func EnterFubenRoom(hwnd win.HWND) {
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("enter fuben room")
	utils.ClickPointByType(hwnd, defs.PointFubenEnter, defs.TimeWaitLong)
}

func ClickFubenStart(hwnd win.HWND) {
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("start fuben")
	utils.ClickPointByType(hwnd, defs.PointFightStart, defs.TimeWaitLong)
}

func ClickFubenStartAck(hwnd win.HWND) {
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("start fuben ack")
	utils.ClickPointByType(hwnd, defs.PointFubenFightStartAck, defs.TimeWaitLong)
}

func SelectFubenMap(hwnd win.HWND, lv defs.FubenLv, isBossFight bool, position defs.FubenPosition) error {
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("click fuben select")
	utils.ClickPointByType(hwnd, defs.PointFubenSelect, defs.TimeWaitMid)
	// 先选择副本类型
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("click fuben type")
	switch position.Type {
	case defs.FubenTypeSpecial:
		utils.ClickPointByType(hwnd, defs.PointFubenTypeSpecial, defs.TimeWaitMid)
	}
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("jump to page")
	// 跳转到副本指定页面
	for i := 1; i < position.Page; i++ {
		nextPageClickTimes := 8
		if i%2 == 0 {
			nextPageClickTimes = 9
		}
		for j := 0; j < nextPageClickTimes; j++ {
			utils.ClickPointByType(hwnd, defs.PointFubenPageDown, 0)
		}
	}
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("select fuben name")
	// 选择副本
	utils.ClickRectByType(hwnd, position.Index, defs.TimeWaitMid)
	data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("select fuben lv")
	// 选择难度
	utils.SelectFubenLv(hwnd, data.GDefsFubenLv.GetStandard(lv))
	// 确认选择
	if isBossFight {
		data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("click boss fight")
		utils.ClickPointByType(hwnd, defs.PointFubenBossFight, defs.TimeWaitMid)
		data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("click select ack")
		utils.ClickPointByType(hwnd, defs.PointFubenSelectAck, defs.TimeWaitMid)
		data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("click boss fight ack")
		utils.ClickPointByType(hwnd, defs.PointFubenBossFightAck, defs.TimeWaitMid)
	} else {
		data.Log().Info().Timestamp().Int("hwnd", int(hwnd)).Msg("click select ack")
		utils.ClickPointByType(hwnd, defs.PointFubenSelectAck, defs.TimeWaitMid)
	}
	return nil
}
