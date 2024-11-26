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

func (self *DefsAngle) GetAngle(hwnd win.HWND) (int, error) {
	standard, _ := utils.CaptureWindowLightWithNormalization(hwnd, defs.GetWinRect(defs.RectTypeAngle), true)
	for angle, gray := range self.items {
		if gray == nil {
			continue
		}
		if utils.IsImageSimilarity(standard, gray, 0.9) {
			return angle, nil
		}
	}
	return 0, fmt.Errorf("not found angle")
}

type DefsFubenLv struct {
	items map[defs.FubenLv]*image.Gray // key->数值
}

func (self *DefsFubenLv) Init() {
	items := make(map[defs.FubenLv]*image.Gray)
	for lv := defs.FubenLvEasy; lv <= defs.FubenLvHero; lv++ {
		img, err := utils.LoadPngToImage(fmt.Sprintf("%sf%d", imgPath, lv))
		if err != nil || img == nil {
			Log().Error().Err(err).Int("lv", int(lv)).Msg("DefsFubenLv.Init, LoadPngToImage failed")
			continue
		}
		v, ok := img.(*image.Gray)
		if !ok {
			Log().Error().Err(err).Int("lv", int(lv)).Msg("DefsFubenLv.Init, img is not *image.Gray")
			continue
		}
		items[lv] = v
		Log().Info().Int("lv", int(lv)).Msg("DefsFubenLv.Init, load success")
	}
	self.items = items
}

func (self *DefsFubenLv) GetPoint(hwnd win.HWND, lv defs.FubenLv) (defs.Point, error) {
	standard, ok := self.items[lv]
	if !ok {
		return defs.EmptyPoint, fmt.Errorf("not found lv：%d", lv)
	}
	var point defs.Point
	defs.RangeFubenLevelRect(func(rect *defs.Rect) bool {
		gray, err := utils.CaptureWindowLightWithGray(hwnd, defs.ToWinRect(rect), false)
		if err != nil {
			Log().Error().Err(err).Int("lv", int(lv)).Msg("DefsFubenLv.GetPoint, CaptureWindowLightWithGray failed")
			return false
		}
		if !utils.IsImageSimilarity(standard, gray, 0.9) {
			return false
		}
		point = defs.RectToPoint(rect)
		return true
	})
	return point, nil
}

func Init() {

}
