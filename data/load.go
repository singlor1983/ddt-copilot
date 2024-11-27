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
			Log().Error().Timestamp().Err(err).Int("angle", angle).Msg("load DefsAngle failed, load png failed")
			continue
		}
		v, ok := img.(*image.Gray)
		if !ok {
			Log().Error().Timestamp().Err(err).Int("angle", angle).Msg("load DefsAngle failed, img is not *image.Gray")
			continue
		}
		items[angle] = v
		Log().Info().Timestamp().Int("angle", angle).Msg("load DefsAngle success")
	}
	self.items = items
}

func (self *DefsAngle) GetAngle(hwnd win.HWND) (int, error) {
	gray, _ := utils.CaptureWindowLightWithNormalization(hwnd, defs.GetWinRect(defs.RectTypeAngle), true)
	for angle, standard := range self.items {
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
			Log().Error().Timestamp().Err(err).Int("lv", int(lv)).Msg("load DefsFubenLv failed, load png failed")
			continue
		}
		v, ok := img.(*image.Gray)
		if !ok {
			Log().Error().Timestamp().Err(err).Int("lv", int(lv)).Msg("load DefsFubenLv failed, img is not *image.Gray")
			continue
		}
		items[lv] = v
		Log().Info().Timestamp().Int("lv", int(lv)).Msg("load DefsFubenLv success")
	}
	self.items = items
}

func (self *DefsFubenLv) GetStandard(lv defs.FubenLv) *image.Gray {
	standard, _ := self.items[lv]
	return standard
}

type DefsOther struct {
	items map[defs.RectType]*image.Gray
}

func (self *DefsOther) load(tp defs.RectType, name string, items map[defs.RectType]*image.Gray) {
	var err error
	defer func() {
		if err != nil {
			Log().Error().Timestamp().Err(err).Int("type", int(tp)).Str("name", name).Msg("load DefsOther failed")
		} else {
			Log().Info().Timestamp().Err(err).Int("type", int(tp)).Str("name", name).Msg("load DefsOther success")
		}
	}()

	img, err := utils.LoadPngToImage(fmt.Sprintf("%s%s", imgPath, name))
	if err != nil {
		return
	}
	v, ok := img.(*image.Gray)
	if !ok {
		err = fmt.Errorf("img is not *image.Gray")
		return
	}

	items[tp] = v
}

func (self *DefsOther) Init() {
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
	self.load(defs.RectTypeIndexPage, "index", items)
	self.load(defs.RectTypeIsYourTurn, "iyt", items)
	self.load(defs.RectTypeSettleWin, "win", items)
	self.load(defs.RectTypeSettleFail, "fail", items)
	self.load(defs.RectTypeFightWin, "fwin", items)
	self.load(defs.RectTypeFightFail, "ffail", items)

	self.items = items
}

func (self *DefsOther) GetStandard(tp defs.RectType) *image.Gray {
	standard := self.items[tp]
	return standard
}

func Init() {

}
