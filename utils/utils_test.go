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

func TestSaveAngle(t *testing.T) {
	hwnd := GetFirstDDTHwnds()

	i := 0
	for {
		img, _ := CaptureWindowLight(hwnd, defs.GetWinRect(defs.RectTypeAngle), true)
		grayImg := ConvertToGrayWithNormalization(img)
		_ = SaveImageToPng(grayImg, fmt.Sprintf("%d", i))
		imgLoad, _ := LoadPngToImage(fmt.Sprintf("%d", i))
		imgLoadNew := imgLoad.(*image.Gray)
		println(imgLoadNew)
		i++
		time.Sleep(time.Second * 3)
	}
}
