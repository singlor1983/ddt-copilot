package utils

import (
	"ddt-copilot/defs"
	"fmt"
	"github.com/lxn/win"
	"github.com/stretchr/testify/assert"
	"image"
	"testing"
	"time"
)

func GetFirstDDTHwnds() win.HWND {
	hwnds := GetDDTHwnds()
	return hwnds[0]
}

func TestCaptureWindow(t *testing.T) {
	hwnd := GetFirstDDTHwnds()
	imgSave, _ := CaptureWindow(hwnd, nil)
	_ = SaveImageToPng(imgSave, "1")
	imgLoad, _ := LoadPngToImage("1")
	imgLoadNew := imgLoad.(*image.RGBA)
	assert.Equal(t, imgSave.Stride, imgLoadNew.Stride)
	assert.Equal(t, imgSave.Rect, imgLoadNew.Rect)
	assert.Equal(t, imgSave.Pix, imgLoadNew.Pix)
}

func TestGenAngle(t *testing.T) {
	hwnd := GetFirstDDTHwnds()

	i := 91
	for {
		fmt.Printf("save: %d\n", i)
		grayImg, _ := CaptureWindowLightWithNormalization(hwnd, defs.GetWinRect(defs.RectTypeAngle), true)
		_ = SaveImageToPng(grayImg, fmt.Sprintf("%d", i))
		i++
		time.Sleep(time.Second * 2)
	}
}

func TestGenFubenLv(t *testing.T) {
	hwnd := GetFirstDDTHwnds()

	i := 0
	defs.RangeFubenLevelRect(func(rect *defs.Rect) bool {
		grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(rect), false)
		_ = SaveImageToPng(grayImg, fmt.Sprintf("f%d", i))
		i++
		return false
	})
}

func TestGenOther(t *testing.T) {
	//hwnd := GetFirstDDTHwnds()

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectPassBtn), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("pbtn"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectFubenSelectText), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("fst"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectFubenInviteAndChangeTeam), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("fiact"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectJinjiInviteAndChangeArea), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("jiaca1"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectJinjiInviteAndChangeArea), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("jiaca2"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectFubenHall), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("fh"))
	//
	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectJinjiHall), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("jh"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectFightRightTop), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("frt"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectFightResult), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("fr"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectFightLoading), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("fl1"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectFightLoading), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("fl2"))

	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectFubenFightSettle), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("ffs")) // 就是翻牌界面，小关需要手动翻牌，竞技用的就是小关
	//
	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectJinjiFightSettle), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("jfs")) // 就是翻牌界面，小关需要手动翻牌，竞技用的就是小关
	//
	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectBackAndExit), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("back"))
	//
	//grayImg, _ := CaptureWindowLightWithGray(hwnd, defs.ToWinRect(defs.RectBackAndExit), false)
	//_ = SaveImageToPng(grayImg, fmt.Sprintf("exit"))
	//
}
