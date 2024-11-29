package core

import (
	"ddt-copilot/data"
	"ddt-copilot/defs"
	"ddt-copilot/utils"
	"fmt"
	"github.com/lxn/win"
	"github.com/rs/zerolog"
	"runtime/debug"
	"sync"
	"time"
)

var (
	monitorListFightProcess = []defs.RectType{
		defs.RectTypePassBtn,
		defs.RectTypeSettleWin,
		defs.RectTypeSettleFail,
		defs.RectTypeFanCardSmall,
		defs.RectTypeFanCardBoss,
		defs.RectTypeFubenSelectText,
		defs.RectTypeFubenInviteAndChangeTeam,
	}
	//monitorListFubenFanCard = []defs.RectType{
	//	defs.RectTypeFanCardSmall,
	//	defs.RectTypeFanCardBoss,
	//}
	//monitorListFubenRoom = []defs.RectType{
	//	defs.RectTypeFubenSelectText,
	//	defs.RectTypeFubenInviteAndChangeTeam,
	//}
)

type ScriptCtrl struct {
	hwnd win.HWND // 窗口句柄
	// 以下选项启动后不能修改，界面置灰
	id          defs.FunctionID // 脚本功能索引
	lv          defs.FubenLv    // 副本难度-不是所有功能都和难度有关
	isChild     bool            // 是否为小号
	isBossFight bool            // 是否为boss战
	masterHWND  win.HWND        // 所属主号句柄

	mu        sync.RWMutex
	isrunning bool // 是否正在运行

	// 内存数据
	fightCount                      int                 // 总关卡挑战次数
	settleWinCount, settleFailCount int                 // 最终结算的胜利失败次数
	fightWinCount, fightFailCount   int                 // 小关的胜利失败次数
	roundCount                      int                 // 第几次自己的回合出手
	initPosition                    defs.InitPosition   // 初始位置
	childs                          utils.Set[win.HWND] // 小号列表，小号开始运行时加入主号的列表
	readyState                      defs.ReadyState     // 是否已准备，副本结束时需要重置状态
}

func NewScriptCtrl(hwnd, master win.HWND, id defs.FunctionID, lv defs.FubenLv, isChild, isBossFight bool) *ScriptCtrl {
	return &ScriptCtrl{hwnd: hwnd, id: id, lv: lv, isChild: isChild, isBossFight: isBossFight, masterHWND: master, childs: make(map[win.HWND]utils.Empty)}
}

func (self *ScriptCtrl) AddChild(child win.HWND) error {
	if self == nil {
		return fmt.Errorf("self nil")
	}
	if len(self.childs) >= 3 {
		return fmt.Errorf("child list already full")
	}
	if self.childs.Contain(child) {
		return fmt.Errorf("add repeat child:%d", child)
	}
	self.mu.Lock()
	defer self.mu.Unlock()
	self.childs.Add(child)
	return nil
}

func (self *ScriptCtrl) DelChild(child win.HWND) {
	if self == nil {
		return
	}
	self.mu.Lock()
	defer self.mu.Unlock()
	self.childs.Del(child)
}

func (self *ScriptCtrl) SetReadyState(state defs.ReadyState) {
	if self == nil {
		return
	}
	self.mu.Lock()
	defer self.mu.Unlock()
	self.readyState = state
}

func (self *ScriptCtrl) GetReadyState() defs.ReadyState {
	if self == nil {
		return defs.ReadyStateNo
	}
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.readyState
}

func (self *ScriptCtrl) Run() {
	if self.isrunning {
		// Todo 通知客户端，脚本正在运行
		data.Log().Error().Timestamp().Int("hwnd", int(self.hwnd)).Msg("script ctrl in running")
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				// Todo 这里退出的时候通知客户端显示退出的原因
				data.Log().Error().Timestamp().Msg(fmt.Sprintf("script ctrl panic! hwnd：%d\n%v|%s", self.hwnd, err, string(debug.Stack())))
			}
		}()
		if self.id > defs.FunctionIDFubenBegin && self.id < defs.FunctionIDFubenEnd {
			self.runFuben()
		} else if self.id > defs.FunctionIDJinjiBegin && self.id < defs.FunctionIDJinjiEnd {
			self.runJinji()
		} else if self.id > defs.FunctionIDOtherBegin && self.id < defs.FunctionIDOtherEnd {
			self.runOther()
		}
	}()
}

func (self *ScriptCtrl) tryReady() {
	if IsReady(self.hwnd) {
		return
	}
	// 点击准备
	utils.ClickPoint(self.hwnd, defs.GetPoint(defs.PointFightStart), defs.TimeWaitShort)
	// 设置状态
	self.SetReadyState(defs.ReadyStateOK)
	self.logger().Msg("set ready ok")
}

func (self *ScriptCtrl) GetFubenPosition() defs.FubenSetting {
	lv := self.lv
	isBossFight := self.isBossFight
	position := data.GGameSetting.SettingFubenPosition.Position[self.id]
	if self.id == defs.FunctionIDCustomFuben {
		lv = data.GGameSetting.SettingFubenCustom.Lv
		isBossFight = data.GGameSetting.SettingFubenCustom.IsBossFightEnable
		position = defs.FubenPosition{
			Type:  data.GGameSetting.SettingFubenCustom.Type,
			Page:  data.GGameSetting.SettingFubenCustom.Page,
			Index: data.GGameSetting.SettingFubenCustom.Index,
		}
	}
	return defs.FubenSetting{
		Lv:                lv,
		IsBossFightEnable: isBossFight,
		FubenPosition:     position,
	}
}

func (self *ScriptCtrl) SelectFubenMap() {
	setting := self.GetFubenPosition()
	if err := SelectFubenMap(self.hwnd, setting.Lv, setting.IsBossFightEnable, setting.FubenPosition); err != nil {
		panic(err)
	}
}

func (self *ScriptCtrl) tryUsePetFood() {
	count := data.GGameSetting.UsePetFoodByFightCount
	if count <= 0 {
		return
	}
	if self.fightCount > 0 && self.fightCount%count == 0 { // 每count局使用一次宠物粮食
		self.logger().Int("fightCount", self.fightCount).Msg("eat pet food")
		utils.FocusDDTWindow(self.hwnd, true)
		utils.ClickPoint(self.hwnd, defs.GetPoint(defs.PointBackpack), defs.TimeWaitMid)
		utils.ClickPoint(self.hwnd, defs.GetPoint(defs.PointBackpackPet), defs.TimeWaitMid)
		utils.ClickPoint(self.hwnd, defs.GetPoint(defs.PointBackpackItem), defs.TimeWaitMid)
		utils.ClickPoint(self.hwnd, defs.GetPoint(defs.PointBackpackIndex1), defs.TimeWaitMid)
		utils.ClickPoint(self.hwnd, defs.GetPoint(defs.PointBackpackPetFoodPutDown), defs.TimeWaitMid)
		utils.ClickPoint(self.hwnd, defs.GetPoint(defs.PointBackpackPetFoodPutDown), defs.TimeWaitMid)
		utils.ClickPoint(self.hwnd, defs.GetPoint(defs.PointBackpackPetFoodEatAck), defs.TimeWaitMid)
		utils.ClickPoint(self.hwnd, defs.GetPoint(defs.PointBackpackPetFoodEat), defs.TimeWaitMid)
		utils.FocusDDTWindow(self.hwnd, true)
	}
}

func (self *ScriptCtrl) OnBeforeStartFuben() {
	self.tryUsePetFood()
	DoHandleBeforeStart(self)
}

func (self *ScriptCtrl) waitChildReady() {
	if self.isChild { // 本身就是小号无需监听
		return
	}
	for {
		self.logger().Msg("check child ready")
		allReady := true
		for child := range self.childs {
			ctrl := InstanceMgr().GetCtrl(child)
			if ctrl == nil {
				continue
			}
			if ctrl.readyState != defs.ReadyStateOK {
				allReady = false
				self.logger().Int("child", int(ctrl.hwnd)).Msg("child not ready")
				break
			}
		}
		if allReady {
			self.logger().Msg("all child ready")
			return
		}
		time.Sleep(time.Second) // 一秒检查一次
	}
}

func (self *ScriptCtrl) StartFuben() {
	self.OnBeforeStartFuben()
	self.waitChildReady()
	self.fightCount++   // 重新启动脚本自然会重置【重启脚本需要重新初始化ScriptCtrl】
	self.roundCount = 0 // 每次新开战斗要重置回合次数
	ClickFubenStart(self.hwnd)
	ClickFubenStartAck(self.hwnd) // 多人战斗的时候没有确认框，但是点击一下也不影响
	self.logger().Int("fightCount", self.fightCount).Msg("start fuben enter fight")
}

// EnterFubenOnce 在房间内执行一次进入副本的操作
func (self *ScriptCtrl) EnterFubenOnce(needSelectFuben bool) {
	if needSelectFuben { // 多关卡的是不需要再次选择副本的,实际上点击一次也没关系，没有选项
		self.SelectFubenMap()
	}
	self.StartFuben()
}

func (self *ScriptCtrl) logger() *zerolog.Event {
	return data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Int("id", int(self.id)).
		Int("lv", int(self.lv)).Bool("isChild", self.isChild).Bool("isBossFight", self.isBossFight).
		Int("masterHWND", int(self.masterHWND))
}

func (self *ScriptCtrl) monitor() {
	self.logger().Msg("monitor")
	for _, rectType := range monitorListFightProcess { // case的先后顺序有严格要求
		standard := data.GDefsOther.GetStandard(rectType)
		switch rectType {
		case defs.RectTypePassBtn: // case分子下是否需要return视功能而定，Todo 第一回合总是识别不到？
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.roundCount++
				self.logger().Int("roundCount", self.roundCount).Msg("is your turn")
				if self.roundCount == 1 {
					DoHandleFirstRoundInit(self)
				}
				DoHandleFight(self)
				time.Sleep(time.Second)
			}
		case defs.RectTypeSettleWin: // 最终关胜利-小关没有
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.settleWinCount++
				self.logger().Int("settleWinCount", self.settleWinCount).Msg("final settle win")
				time.Sleep(time.Second * 3)
			}
		case defs.RectTypeSettleFail: // 最终关失败-小关没有
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.settleFailCount++
				self.logger().Int("settleFailCount", self.settleFailCount).Msg("final settle fail")
				time.Sleep(time.Second * 3)
			}
		case defs.RectTypeFanCardSmall: // Todo 执行翻牌操作
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.logger().Msg("fan card small")
				time.Sleep(time.Second * 18)
			}
		case defs.RectTypeFanCardBoss: // Todo 执行翻牌操作
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.logger().Msg("fan card boss")
				time.Sleep(time.Second * 18)
			}
		case defs.RectTypeFubenSelectText:
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.logger().Msg("next select fuben and start")
				self.EnterFubenOnce(true)
				time.Sleep(time.Second)
			}
		case defs.RectTypeFubenInviteAndChangeTeam:
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.logger().Msg("next direct start fuben")
				self.EnterFubenOnce(false)
				time.Sleep(time.Second)
			}
		}
	}
	time.Sleep(time.Second)
	self.monitor()
}

func (self *ScriptCtrl) runFuben() {
	if self.isChild {
		if !InFubenRoom(self.hwnd) { // 小号的话只允许在副本房间内开启运行
			// Todo 通知客户端显示，做成一个接口把
			panic(fmt.Sprintf("小号需要在副本房间内再开启脚本. hwnd:%d", self.hwnd))
		}
		if err := InstanceMgr().GetCtrl(self.masterHWND).AddChild(self.hwnd); err != nil {
			panic(err)
		}
		self.tryReady()
		// 等待主号开启副本战斗之后就行了
	} else {
		// 主号只支持在特定场景下开启脚本，没必要考虑到所有的场景情况
		if InJinRoom(self.hwnd) || InJinjiHall(self.hwnd) {
			BackToIndexPage(self.hwnd)
		}
		if InIndexPage(self.hwnd) {
			EnterFubenHall(self.hwnd)
		}
		if InFubenHall(self.hwnd) {
			EnterFubenRoom(self.hwnd)
		}
		self.monitor()
	}
}

func (self *ScriptCtrl) runJinji() {
}

func (self *ScriptCtrl) runOther() {
}
