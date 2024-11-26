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
