package utils

import (
	"ddt-copilot/defs"
	"fmt"
	"github.com/lxn/win"
	"image"
	"math"
	"time"
)

// 各弹弹堂版本通用的工具函数

// GetDDTWindowsHWND 获取所有的弹弹堂游戏窗口
func GetDDTWindowsHWND() []win.HWND {
	var ret []win.HWND
	pids, _ := GetProcessID(string(defs.ProcessTgWeb))
	for _, pid := range pids {
		hwnd, err := GetFirstWindowByPID(pid)
		if err != nil {
			continue
		}
		wds, _ := GetAllChildWindows(hwnd)
		if len(wds) != 5 {
			continue
		}
		lastWd := win.HWND(wds[len(wds)-1])
		ret = append(ret, lastWd)
		threadID := win.GetWindowThreadProcessId(lastWd, &pid)
		fmt.Printf("hwnd:%d pid:%d threadID:%d\n", lastWd, pid, threadID)
	}
	return ret
}

func FocusDDTWindow(hwnd win.HWND, dropBlock bool) {
	//win.SetFocus(hwnd) //选取这个窗口
	win.NotifyWinEvent(win.EVENT_OBJECT_FOCUS, hwnd, win.OBJID_CLIENT, 0)
	win.NotifyWinEvent(win.EVENT_OBJECT_HIDE, hwnd, win.OBJID_CARET, 0)
	win.NotifyWinEvent(win.EVENT_OBJECT_IME_SHOW, hwnd, win.OBJID_CARET, 0)
	point := defs.GetPointByElement(defs.ElementEmpty)
	Click(hwnd, point[0], point[1], 0) // 激活窗口，便于截图为亮色的
	if dropBlock {                     // 副本选择页面不能删除遮挡，因为这本身就是遮挡
		for i := 0; i < 5; i++ { // 截屏之前用ESC把其他遮挡界面关闭，避免影响判断，5次是因为可能有多层折叠
			KeyBoard(hwnd, defs.VK_ESCAPE, 0)
		}
	}
	time.Sleep(time.Millisecond * 100) // 等待那个继续操作图标消失
}

func CaptureWindowLight(hwnd win.HWND, captureRect *win.RECT, dropBlock bool) (*image.RGBA, error) {
	FocusDDTWindow(hwnd, dropBlock)
	return CaptureWindow(hwnd, captureRect)
}

func ClickElement(hwnd win.HWND, element defs.DDTClickElement, duration time.Duration) {
	point := defs.GetPointByElement(element)
	if defs.IsEmptyPoint(point) {
		return
	}
	Click(hwnd, point[0], point[1], duration)
}

func ClickElementAfterLongEx(hwnd win.HWND, element defs.DDTClickElement, duration time.Duration) {
	ClickElement(hwnd, element, duration)
	time.Sleep(defs.ClickWaitLong)
}

func ClickElementAfterMidEx(hwnd win.HWND, element defs.DDTClickElement, duration time.Duration) {
	ClickElement(hwnd, element, duration)
	time.Sleep(defs.ClickWaitMid)
}

func ClickElementAfterShortEx(hwnd win.HWND, element defs.DDTClickElement, duration time.Duration) {
	ClickElement(hwnd, element, duration)
	time.Sleep(defs.ClickWaitShort)
}

func ClickElementAfterLong(hwnd win.HWND, element defs.DDTClickElement) {
	ClickElementAfterLongEx(hwnd, element, 0)
}

func ClickElementAfterMid(hwnd win.HWND, element defs.DDTClickElement) {
	ClickElementAfterMidEx(hwnd, element, 0)
}

func ClickElementAfterShort(hwnd win.HWND, element defs.DDTClickElement) {
	ClickElementAfterShortEx(hwnd, element, 0)
}

func ClickByRect(hwnd win.HWND, rect win.RECT, duration time.Duration) {
	point := defs.NewPointWithRectCenter(rect)
	if defs.IsEmptyPoint(point) {
		return
	}
	Click(hwnd, point[0], point[1], duration)
}

func ClickByRectAfterLongEx(hwnd win.HWND, rect win.RECT, duration time.Duration) {
	ClickByRect(hwnd, rect, duration)
	time.Sleep(defs.ClickWaitLong)
}

func ClickByRectAfterMidEx(hwnd win.HWND, rect win.RECT, duration time.Duration) {
	ClickByRect(hwnd, rect, duration)
	time.Sleep(defs.ClickWaitMid)
}

func ClickByRectAfterShortEx(hwnd win.HWND, rect win.RECT, duration time.Duration) {
	ClickByRect(hwnd, rect, duration)
	time.Sleep(defs.ClickWaitShort)
}

func ClickByRectAfterLong(hwnd win.HWND, rect win.RECT) {
	ClickByRectAfterLongEx(hwnd, rect, 0)
}

func ClickByRectAfterMid(hwnd win.HWND, rect win.RECT) {
	ClickByRectAfterMidEx(hwnd, rect, 0)
}

func ClickByRectAfterShort(hwnd win.HWND, rect win.RECT) {
	ClickByRectAfterShortEx(hwnd, rect, 0)
}

// SelectFubenLevel 根据难度选择副本
func SelectFubenLevel(hwnd win.HWND, level defs.FubenLevel) {
	standerImg := defs.GetFubenLevelElementImg(level)
	defs.RangeFubenLevelElementRect(func(rect *defs.DDTElementRect) bool {
		if rect == nil {
			return false
		}
		winRect := defs.NewElementRectWithDDTRect(rect)
		img, err := CaptureWindowLight(hwnd, winRect, false)
		if err != nil {
			return false
		}
		grayImg := ConvertToGray(img)
		diff, err := CompareGrayImages(grayImg, standerImg)
		if err != nil {
			return false
		}
		if !IsSimilarity(diff, defs.ImgSimilarityThresholdFubenLevel) {
			return false
		}
		ClickByRectAfterMid(hwnd, *winRect)
		return true
	})
}

// GetAngle 获得当前的角度
func GetAngle(hwnd win.HWND) int {
	img, err := CaptureWindowLight(hwnd, defs.GetElementRect(defs.DDTElementRectTypeAngle), true)
	if err != nil {
		return 0
	}
	grayImg := ConvertToGrayWithNormalization(img)
	var curAngle int
	defs.RangeAngleElementImg(func(num int, gray *image.Gray) bool {
		if gray == nil {
			return false
		}
		diff, inErr := CompareGrayImages(grayImg, gray)
		if inErr != nil {
			return false
		}
		if !IsSimilarity(diff, defs.ImgSimilarityThresholdAngle) {
			return false
		}
		curAngle = num
		return true
	})
	return curAngle
}

// UpdateAngle 改变力度
func UpdateAngle(hwnd win.HWND, angle int) {
	direction := defs.DirectionUp
	if angle < 0 {
		direction = defs.DirectionDown
		angle = int(math.Abs(float64(angle)))
	}
	switch direction {
	case defs.DirectionUp:
		for i := 0; i < angle; i++ {
			KeyBoard(hwnd, defs.VK_UP, 0)
		}
	case defs.DirectionDown:
		for i := 0; i < angle; i++ {
			KeyBoard(hwnd, defs.VK_DOWN, 0)
		}
	}
}

// ConfirmDirection 转向、确认方向
func ConfirmDirection(hwnd win.HWND, direction defs.Direction) {
	switch direction {
	case defs.DirectionLeft:
		KeyBoard(hwnd, defs.VK_RIGHT, 0)
		KeyBoard(hwnd, defs.VK_LEFT, 0)
	case defs.DirectionRight:
		KeyBoard(hwnd, defs.VK_LEFT, 0)
		KeyBoard(hwnd, defs.VK_RIGHT, 0)
	}
}

// Move 移动。1格距离100毫秒
func Move(hwnd win.HWND, direction defs.Direction, distance int) {
	ConfirmDirection(hwnd, direction)
	if distance < 0 {
		distance = 0
	}
	ts := time.Duration(distance*100) * time.Millisecond
	switch direction {
	case defs.DirectionLeft:
		KeyBoard(hwnd, defs.VK_LEFT, ts)
	case defs.DirectionRight:
		KeyBoard(hwnd, defs.VK_RIGHT, ts)
	}
}

func UseSkill(hwnd win.HWND, skill uintptr) {
	KeyBoard(hwnd, skill, 0)
}

// Launch 攻击、发射，根据力度计算出需要按压持续的时间。1度40毫秒
func Launch(hwnd win.HWND, power int) {
	if power < 0 {
		power = 0
	}
	if power > 100 {
		power = 100
	}
	ts := time.Duration(power*40) * time.Millisecond
	KeyBoard(hwnd, defs.VK_SPACE, ts)
}

func IsImgSimilarity(checker, stander *image.Gray, threshold float64) bool {
	if checker == nil || stander == nil {
		return false
	}
	diff, err := CompareGrayImages(checker, stander)
	if err != nil {
		return false
	}
	return IsSimilarity(diff, threshold)
}

// GetSense 获取当前所处的场景，后续根据当前所处的场景不同执行不同的操作
func GetSense(hwnd win.HWND) defs.SenseType {
	st := defs.SenseTypeOther
	defs.RangeSenseRect(func(rectType defs.DDTElementRectType) bool {
		rect := defs.GetElementRect(rectType)
		if rect == nil {
			return true
		}
		img, err := CaptureWindowLight(hwnd, rect, true)
		if err != nil {
			return false
		}
		grayImg := ConvertToGray(img)
		switch rectType {
		case defs.DDTElementRectTypeFubenInviteAndChangeTeam:
			if !IsImgSimilarity(grayImg, defs.ElementImgFubenInviteAndChangeTeam, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeFubenRoom
		case defs.DDTElementRectTypeJinjiInviteAndChangeArea1:
			if !IsImgSimilarity(grayImg, defs.ElementImgJinjiInviteAndChangeArea1, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeJinjiRoom
		case defs.DDTElementRectTypeJinjiInviteAndChangeArea2:
			if !IsImgSimilarity(grayImg, defs.ElementImgJinjiInviteAndChangeArea2, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeJinjiRoom
		case defs.DDTElementRectTypeFubenHall:
			if !IsImgSimilarity(grayImg, defs.ElementImgFubenHall, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeFubenHall
		case defs.DDTElementRectTypeJinjiHall:
			if !IsImgSimilarity(grayImg, defs.ElementImgJinjiHall, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeJinjiHall
		case defs.DDTElementRectTypeFightRightTop:
			if !IsImgSimilarity(grayImg, defs.ElementImgFightRightTop, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeInFight
		case defs.DDTElementRectTypeFightResult:
			if !IsImgSimilarity(grayImg, defs.ElementImgFightResult, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeFightResult
		case defs.DDTElementRectTypeFubenFightLoading:
			if !IsImgSimilarity(grayImg, defs.ElementImgFubenFightLoading, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeFubenFightLoading
		case defs.DDTElementRectTypeJinjiFightLoading:
			if !IsImgSimilarity(grayImg, defs.ElementImgJinjiFightLoading, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeJinjiFightLoading
		case defs.DDTElementRectTypeFubenFightSettle:
			if !IsImgSimilarity(grayImg, defs.ElementImgFubenFightSettle, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeFubenFightSettle
		case defs.DDTElementRectTypeJinjiFightSettle:
			if !IsImgSimilarity(grayImg, defs.ElementImgJinjiFightSettle, defs.ImgSimilarityThresholdSence) {
				return false
			}
			st = defs.SenseTypeJinjiFightSettle
		default:
			return false
		}
		return true
	})
	return st
}

func BackToIndexPage(hwnd win.HWND) {
	for {
		rect := defs.GetElementRect(defs.DDTElementRectTypeBack)
		if rect == nil {
			return
		}
		img, err := CaptureWindowLight(hwnd, rect, true)
		if err != nil {
			return
		}
		grayImg := ConvertToGray(img)
		if !IsImgSimilarity(grayImg, defs.ElementImgBack, defs.ImgSimilarityThresholdBack) {
			break
		}
		ClickElementAfterLong(hwnd, defs.ElementBackAndExit)
	}
}

func SelectFubenMap(hwnd win.HWND, t defs.FubenID, l defs.FubenLevel, isBossFight bool, fubenPosition map[defs.FubenID]defs.FubenPosition) error {
	position, ok := fubenPosition[t]
	if !ok {
		return fmt.Errorf("not found fuben position:%d", t)
	}
	ClickElementAfterMid(hwnd, defs.ElementFubenSelect)
	// 先选择副本类型
	switch position.Type {
	case defs.FubenTypeSpecial:
		ClickElementAfterMid(hwnd, defs.ElementFubenTypeSpecial)
	}
	// 跳转到副本指定页面
	for i := 1; i < position.Page; i++ {
		nextPageClickTimes := 8
		if i%2 == 0 {
			nextPageClickTimes = 9
		}
		for j := 0; j < nextPageClickTimes; j++ {
			ClickElement(hwnd, defs.ElementFubenPageDown, 0)
		}
	}
	// 选择副本
	ClickByRectAfterMid(hwnd, *defs.GetElementRect(defs.DDTElementRectType(position.Count)))
	// 选择难度
	SelectFubenLevel(hwnd, l)
	// 确认选择
	if isBossFight {
		ClickElementAfterMid(hwnd, defs.ElementFubenBossFight)
		ClickElementAfterMid(hwnd, defs.ElementFubenAck)
		ClickElementAfterMid(hwnd, defs.ElementFubenBossFightAck)
	} else {
		ClickElementAfterMid(hwnd, defs.ElementFubenAck)
	}
	return nil
}

func GetRectImgFromWindowImg(hwnd win.HWND, windowImg *image.RGBA, captureRect *win.RECT) *image.RGBA {
	var err error
	wdImgTmp := windowImg
	if wdImgTmp == nil {
		wdImgTmp, err = CaptureWindowLight(hwnd, nil, true)
		if err != nil {
			panic(err.Error())
		}
	}
	return ExtractRegion(wdImgTmp, captureRect)
}

func InSenseIsYourTurn(hwnd win.HWND, windowImg *image.RGBA) bool {
	img := GetRectImgFromWindowImg(hwnd, windowImg, defs.GetElementRect(defs.DDTElementRectTypeIsYourTurn))
	diff, err := CompareGrayImages(ConvertToGray(img), defs.ElementImgIsYourTurn)
	if err != nil {
		panic(err.Error())
	}
	return IsSimilarity(diff, defs.ImgSimilarityThresholdIsYourTurn)
}

func InSensePassBtn(hwnd win.HWND, windowImg *image.RGBA) bool {
	img := GetRectImgFromWindowImg(hwnd, windowImg, defs.GetElementRect(defs.DDTElementRectTypePassBtn))
	diff, err := CompareGrayImages(ConvertToGray(img), defs.ElementImgPassBtn)
	if err != nil {
		panic(err.Error())
	}
	return IsSimilarity(diff, defs.ImgSimilarityThresholdPassBtn)
}

func InSenseFightWin(hwnd win.HWND, windowImg *image.RGBA) bool {
	img := GetRectImgFromWindowImg(hwnd, windowImg, defs.GetElementRect(defs.DDTElementRectTypeWinOrFail))
	diff, err := CompareGrayImages(ConvertToGray(img), defs.ElementImgFightWin)
	if err != nil {
		panic(err.Error())
	}
	return IsSimilarity(diff, defs.ImgSimilarityThresholdWinOrFail)
}

func InSenseFighFail(hwnd win.HWND, windowImg *image.RGBA) bool {
	img := GetRectImgFromWindowImg(hwnd, windowImg, defs.GetElementRect(defs.DDTElementRectTypeWinOrFail))
	diff, err := CompareGrayImages(ConvertToGray(img), defs.ElementImgFighFail)
	if err != nil {
		panic(err.Error())
	}
	return IsSimilarity(diff, defs.ImgSimilarityThresholdWinOrFail)
}

func InSenseFubenFanCardSmall(hwnd win.HWND, windowImg *image.RGBA) bool {
	img := GetRectImgFromWindowImg(hwnd, windowImg, defs.GetElementRect(defs.DDTElementRectTypeJinjiFightSettle))
	diff, err := CompareGrayImages(ConvertToGray(img), defs.ElementImgJinjiFightSettle)
	if err != nil {
		panic(err.Error())
	}
	return IsSimilarity(diff, defs.ImgSimilarityThresholdSence)
}

func InSenseFubenFanCardBoss(hwnd win.HWND, windowImg *image.RGBA) bool {
	img := GetRectImgFromWindowImg(hwnd, windowImg, defs.GetElementRect(defs.DDTElementRectTypeFubenFightSettle))
	diff, err := CompareGrayImages(ConvertToGray(img), defs.ElementImgFubenFightSettle)
	if err != nil {
		panic(err.Error())
	}
	return IsSimilarity(diff, defs.ImgSimilarityThresholdSence)
}

func InSenseFubenSelectText(hwnd win.HWND, windowImg *image.RGBA) bool {
	img := GetRectImgFromWindowImg(hwnd, windowImg, defs.GetElementRect(defs.DDTElementRectTypeFubenSelectText))
	diff, err := CompareGrayImages(ConvertToGray(img), defs.ElementImgFubenSelectText)
	if err != nil {
		panic(err.Error())
	}
	return IsSimilarity(diff, defs.ImgSimilarityThresholdFubenSelectText)
}

func InSenseFubenInviteAndChangeTeam(hwnd win.HWND, windowImg *image.RGBA) bool {
	img := GetRectImgFromWindowImg(hwnd, windowImg, defs.GetElementRect(defs.DDTElementRectTypeFubenInviteAndChangeTeam))
	diff, err := CompareGrayImages(ConvertToGray(img), defs.ElementImgFubenInviteAndChangeTeam)
	if err != nil {
		panic(err.Error())
	}
	return IsSimilarity(diff, defs.ImgSimilarityThresholdSence)
}

func IsFubenRoomReady(hwnd win.HWND) bool {
	img, err := CaptureWindowLight(hwnd, defs.GetElementRect(defs.DDTElementRectTypeFightReady), true)
	if err != nil {
		panic(err.Error())
	}
	diff, err := CompareGrayImages(ConvertToGray(img), defs.ElementImgFightRoomCancel)
	if err != nil {
		panic(err.Error())
	}
	return IsSimilarity(diff, defs.ImgSimilarityThresholdSence)
}
