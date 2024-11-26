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
	items map[int]*image.Gray // key->角度
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
	items map[defs.FubenLv]*image.Gray
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

type DefsConstants struct {
	items map[defs.RectType]*image.Gray
}

func (self *DefsConstants) load(tp defs.RectType, name string, items map[defs.RectType]*image.Gray) {
	var err error
	defer func() {
		if err != nil {
			Log().Error().Err(err).Int("tp", int(tp)).Str("name", name).Msg("DefsConstants.load failed")
		} else {
			Log().Info().Err(err).Int("tp", int(tp)).Str("name", name).Msg("DefsConstants.load success")
		}
	}()

	img, err := utils.LoadPngToImage(fmt.Sprintf("%s%s", imgPath, name))
	if err != nil {
		return
	}
	v, ok := img.(*image.Gray)
	if !ok {
		err = fmt.Errorf("not gray image")
		return
	}

	items[tp] = v
}

func (self *DefsConstants) Init() {
	items := make(map[defs.RectType]*image.Gray)
	self.load(defs.RectTypePassBtn, "pbtn", items)
	self.load(defs.RectTypeFubenSelectText, "fst", items)
	self.load(defs.RectTypeFubenInviteAndChangeTeam, "fiact", items)
	self.load(defs.RectTypeJinjiInviteAndChangeArea1, "jiaca1", items)
	self.load(defs.RectTypeJinjiInviteAndChangeArea2, "jiaca2", items)
	self.load(defs.RectTypeFubenHall, "fh", items)
	self.load(defs.RectTypeJinjiHall, "jh", items)
	self.load(defs.RectTypeFightRightTop, "frt", items)
	self.load(defs.RectTypeFightResult, "fr", items)
	self.load(defs.RectTypeFubenFightLoading, "fl1", items)
	self.load(defs.RectTypeJinjiFightLoading, "fl2", items)
	self.load(defs.RectTypeFubenFightSettle, "ffs", items)
	self.load(defs.RectTypeJinjiFightSettle, "jfs", items)
	self.load(defs.RectTypeBack, "back", items)
	self.load(defs.RectTypeExit, "exit", items)

	self.items = items

	self.BackToIndexPage(31065808)
}

func (self *DefsConstants) BackToIndexPage(hwnd win.HWND) {
	standard := self.items[defs.RectTypeBack]
	if standard == nil {
		return
	}
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
		utils.LeftClickPoint(hwnd, defs.PointBackAndExit, defs.ClickWaitLong)
	}
}

func Init() {

}
