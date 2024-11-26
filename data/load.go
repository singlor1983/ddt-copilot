package data

import (
	"ddt-copilot/defs"
	"ddt-copilot/utils"
	"fmt"
	"github.com/lxn/win"
	"image"
)

const (
	imgPath = "./assets/"
)

type DefsAngle struct {
	items map[int]*image.Gray // key->数值
}

func (self *DefsAngle) Init() {
	items := make(map[int]*image.Gray)
	for angle := -40; angle <= 90; angle++ {
		img, err := utils.LoadPngToImage(fmt.Sprintf("%s%d", imgPath, angle))
		if err != nil || img == nil {
			Log().Error().Err(err).Int("angle", angle).Msg("DefsAngle.Init, LoadPngToImage failed")
			continue
		}
		v, ok := img.(*image.Gray)
		if !ok {
			Log().Error().Err(err).Int("angle", angle).Msg("DefsAngle.Init, img is not *image.Gray")
			continue
		}
		items[angle] = v
		Log().Info().Int("angle", angle).Msg("DefsAngle.Init, load success")
	}
	self.items = items
}

func (self *DefsAngle) GetAngle(hwnd win.HWND) int {
	img, _ := utils.CaptureWindowLight(hwnd, defs.GetWinRect(defs.RectTypeAngle), true)
	standard := utils.ConvertToGrayWithNormalization(img)
	for angle, gray := range self.items {
		if gray == nil {
			continue
		}
		if utils.IsImageSimilarity(standard, gray, 0.9) {
			return angle
		}
	}
	return 0
}

func Init() {

}
