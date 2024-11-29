package core

import (
	"ddt-copilot/data"
	"ddt-copilot/defs"
	"ddt-copilot/utils"
	"fmt"
	"github.com/lxn/win"
	"runtime/debug"
	"sync"
	"time"
)

var (
	monitorListFightProcess = []defs.RectType{
		defs.RectTypeIsYourTurn,
		defs.RectTypeFightWin,
		defs.RectTypeFightFail,
	}
	monitorListFubenFanCard = []defs.RectType{
		defs.RectTypeFanCardSmall,
		defs.RectTypeFanCardBoss,
	}
	monitorListFubenRoom = []defs.RectType{
		defs.RectTypeFubenSelectText,
		defs.RectTypeFubenInviteAndChangeTeam,
	}
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
	data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Msg("set ready ok")
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

func (self *ScriptCtrl) OnBeforeStartFuben() {
	// Todo zhangzhihi 还有个通用before，比如喂宠物粮
	DoHandleBeforeStart(self)
}

func (self *ScriptCtrl) waitChildReady() {
	if self.isChild { // 本身就是小号无需监听
		return
	}
	for {
		data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Msg("check child ready")
		allReady := true
		for child := range self.childs {
			ctrl := InstanceMgr().GetCtrl(child)
			if ctrl == nil {
				continue
			}
			if ctrl.readyState != defs.ReadyStateOK {
				allReady = false
				data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Int("child", int(ctrl.hwnd)).Msg("child not ready")
				break
			}
		}
		if allReady {
			data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Msg("all child ready")
			return
		}
		time.Sleep(time.Second) // 一秒检查一次
	}
}

func (self *ScriptCtrl) StartFuben() {
	self.waitChildReady()
	self.fightCount++   // 重新启动脚本自然会重置【重启脚本需要重新初始化ScriptCtrl】
	self.roundCount = 0 // 每次新开战斗要重置回合次数
	ClickFubenStart(self.hwnd)
	ClickFubenStartAck(self.hwnd) // 多人战斗的时候没有确认框，但是点击一下也不影响
}

// EnterFubenOnce 在房间内执行一次进入副本的操作
func (self *ScriptCtrl) EnterFubenOnce(needSelectFuben bool) {
	if needSelectFuben { // 多关卡的是不需要再次选择副本的,实际上点击一次也没关系，没有选项
		self.SelectFubenMap()
	}
	self.OnBeforeStartFuben()
	self.StartFuben()
}

func (self *ScriptCtrl) Monitor(nextMonitorList []defs.RectType) {
	time.Sleep(time.Second)
	data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Interface("nextMonitorList", nextMonitorList).Msg("monitor")
	for _, rectType := range nextMonitorList { // case的先后顺序有严格要求
		standard := data.GDefsOther.GetStandard(rectType)
		switch rectType {
		case defs.RectTypeIsYourTurn: // case分子下是否需要return视功能而定
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.roundCount++
				data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Int("roundCount", self.roundCount).Msg("is your turn")
				if self.roundCount == 1 {
					DoHandleFirstRoundInit(self)
				}
				DoHandleFight(self)
				self.monitorFightProcess()
				return
			}
		case defs.RectTypeFightWin:
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.fightWinCount++
				data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Int("fightWinCount", self.fightWinCount).Msg("fight win")
				return
			}
		case defs.RectTypeFightFail:
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.fightFailCount++
				data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Int("fightFailCount", self.fightFailCount).Msg("fight fail")
				return
			}
		case defs.RectTypeSettleWin:
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.settleWinCount++
				data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Int("settleWinCount", self.settleWinCount).Msg("final settle win")
				return
			}
		case defs.RectTypeSettleFail:
			if utils.IsSimilarity(self.hwnd, standard, rectType, 0.8, true) {
				self.settleFailCount++
				data.Log().Info().Timestamp().Int("hwnd", int(self.hwnd)).Int("settleFailCount", self.settleFailCount).Msg("final settle fail")
				return
			}
			//case defs.RectTypeWinOrFail:
			//	if utils.InSenseFightWin(ctrl.hwnd, img) { // 战斗结算页面-胜利，最后一关打完才有，小关是直接翻牌了
			//		fmt.Printf("time:%s fight win\n", time.Now().Format(time.TimeOnly))
			//		ctrl.winCount++
			//		ctrl.monitorFubenFanCard()
			//		return
			//	} else if utils.InSenseFighFail(ctrl.hwnd, img) { // 战斗结算页面-失败，每一小关都有
			//		fmt.Printf("time:%s fight fail\n", time.Now().Format(time.TimeOnly))
			//		ctrl.failCount++
			//		ctrl.fixFightCount = 0
			//		ctrl.monitorFubenFanCard()
			//		return
			//	}
			//case defs.DDTElementRectTypeFanCardSmall, defs.DDTElementRectTypeFanCardBoss: // 小关翻牌、boss翻牌
			//	if utils.InSenseFubenFanCardSmall(ctrl.hwnd, img) {
			//		fmt.Printf("time:%s fan card small\n", time.Now().Format(time.TimeOnly))
			//		// Todo add 翻牌逻辑
			//		ctrl.monitorFubenRoom()
			//		return
			//	} else if utils.InSenseFubenFanCardBoss(ctrl.hwnd, img) {
			//		fmt.Printf("time:%s fan card boss\n", time.Now().Format(time.TimeOnly))
			//		// Todo add 翻牌逻辑
			//		ctrl.monitorFubenRoom()
			//		return
			//	}
			//case defs.DDTElementRectTypeFubenSelectText: // 有这个select必然重新选择副本并开始
			//	if utils.InSenseFubenSelectText(ctrl.hwnd, img) {
			//		fmt.Printf("time:%s select fuben\n", time.Now().Format(time.TimeOnly))
			//		if ctrl.isChild {
			//			if !utils.IsFubenRoomReady(ctrl.hwnd) { // 如果未准备则准备，准备后通知主号
			//				utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightStart) // 这个也是准备按钮的位置
			//				if err := instanceMgr.SetChildReady(ctrl.masterHWND, ctrl.hwnd, defs.ReadyStateOK); err != nil {
			//					panic(err.Error())
			//				}
			//			}
			//		} else {
			//			ctrl.EnterFubenOnce(true)
			//			ctrl.monitorFightProcess()
			//		}
			//		return
			//	}
			//case defs.DDTElementRectTypeFubenInviteAndChangeTeam: // 先判断select再判断，副本房间特征元素，这样能保证需要选择副本的时候必然可以选择副本 Todo 可能不能这么快就开始，要等小号都准备。小号也不能点击开始和选择副本
			//	if utils.InSenseFubenInviteAndChangeTeam(ctrl.hwnd, img) {
			//		fmt.Printf("time:%s direct enter fuben\n", time.Now().Format(time.TimeOnly))
			//		if ctrl.isChild {
			//			if !utils.IsFubenRoomReady(ctrl.hwnd) { // 如果未准备则准备，准备后通知主号
			//				utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightStart) // 这个也是准备按钮的位置
			//				if err := instanceMgr.SetChildReady(ctrl.masterHWND, ctrl.hwnd, defs.ReadyStateOK); err != nil {
			//					panic(err.Error())
			//				}
			//			}
			//		} else {
			//			ctrl.EnterFubenOnce(false)
			//			ctrl.monitorFightProcess()
			//		}
			//		return
			//	}
		}
	}
	self.Monitor(nextMonitorList)
}

func (self *ScriptCtrl) monitorFightProcess() {
	self.Monitor(monitorListFightProcess)
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
		switch true {
		case InJinRoom(self.hwnd) || InJinjiHall(self.hwnd):
			BackToIndexPage(self.hwnd)
			EnterFubenHall(self.hwnd)
			EnterFubenRoom(self.hwnd)
		case InIndexPage(self.hwnd):
			EnterFubenHall(self.hwnd)
			EnterFubenRoom(self.hwnd)
		case InFubenHall(self.hwnd):
			EnterFubenRoom(self.hwnd)
		case InFubenRoom(self.hwnd):
		//case IsFubenFight(): Todo zhangzhihui 开启脚本时正在战斗内
		default:
			panic(fmt.Sprintf("不支持在当前场景下开启脚本. hwnd:%d", self.hwnd))
		}
		// 经过上面的步骤已经进入副本房间内了
		self.EnterFubenOnce(true)
		self.monitorFightProcess()
	}
}

func (self *ScriptCtrl) runJinji() {
}

func (self *ScriptCtrl) runOther() {
}

//
//func (ctrl *ScriptCtrl) SelectFubenMap() {
//	if err := utils.SelectFubenMap(ctrl.hwnd, self.id, ctrl.lvIndex, ctrl.isBossFight, data.GConfig.GetFubenPosition()); err != nil {
//		panic(err.Error())
//	}
//}
//
//func (ctrl *ScriptCtrl) OnBeforeStartFuben() {
//	// Todo zhangzhihi 还有个通用before，比如喂宠物粮
//	DoHandleBeforeStart(ctrl)
//}
//
//// waitChildReady 开始战斗之前检查小号是否都已经准备
//func (ctrl *ScriptCtrl) waitChildReady() {
//	if ctrl.isChild { // 本身就是小号无需监听
//		return
//	}
//	for {
//		fmt.Printf("time:%s check child ready\n", time.Now().Format(time.TimeOnly))
//		allReady := true
//		for _, node := range ctrl.childHWNDs {
//			if node == nil {
//				continue
//			}
//			if node.readyState == defs.ReadyStateNo {
//				allReady = false
//				break
//			}
//		}
//		if allReady {
//			fmt.Printf("time:%s child all ready\n", time.Now().Format(time.TimeOnly))
//			return
//		}
//		time.Sleep(time.Second) // 一秒检查一次
//	}
//}
//
//func (ctrl *ScriptCtrl) waitMasterStartFuben() {
//	if !ctrl.isChild {
//		return
//	}
//	for {
//		master := instanceMgr.GetCtrl(ctrl.masterHWND)
//		if master == nil {
//			panic("master nil")
//		}
//		if master.fightCount > ctrl.fightCount {
//			ctrl.fightCount = master.fightCount // Todo 用这个来判断是否开启了副本，有待商榷，什么清楚？
//			return
//		}
//	}
//}
//
//func (ctrl *ScriptCtrl) StartFuben() {
//	ctrl.waitChildReady()
//	ctrl.fightCount++ // 重新启动脚本自然会重置【重启脚本需要重新初始化ScriptCtrl】
//	ctrl.fixFightCount++
//	ctrl.roundCount = 0 // 每次新开战斗要重置回合次数
//	utils.ClickElementAfterMid(ctrl.hwnd, defs.ElementFightStart)
//	utils.ClickElementAfterMid(ctrl.hwnd, defs.ElementFubenFightStartAck) // 多人战斗的时候没有确认框，但是点击一下也不影响
//	fmt.Printf("time:%s start fuben: %v\n", time.Now().Format(time.TimeOnly), ctrl)
//}
//
//// EnterFubenOnce 在房间内执行一次进入副本的操作
//func (ctrl *ScriptCtrl) EnterFubenOnce(needSelectFuben bool) {
//	if needSelectFuben { // 多关卡的是不需要再次选择副本的,实际上点击一次也没关系，没有选项
//		ctrl.SelectFubenMap()
//	}
//	ctrl.OnBeforeStartFuben()
//	ctrl.StartFuben()
//}
//
//func (ctrl *ScriptCtrl) monitorFightProcess() {
//	ctrl.Monitor(monitorListFightProcess)
//}
//
//func (ctrl *ScriptCtrl) monitorFubenRoom() {
//	ctrl.Monitor(monitorListFubenRoom)
//}
//
//func (ctrl *ScriptCtrl) monitorFubenFanCard() {
//	ctrl.Monitor(monitorListFubenFanCard)
//}
//
//// Monitor 每秒只截屏一次，然后从这张图里面再截取对应的部分，而不是每次重新截屏
//func (ctrl *ScriptCtrl) Monitor(nextMonitorList []defs.RectType) {
//	time.Sleep(time.Millisecond * 500)
//	img, err := utils.CaptureWindowLight(ctrl.hwnd, nil, true)
//	if err != nil {
//		panic(err.Error())
//	}
//	//// 创建文件并保存为 PNG 格式
//	//file, err := os.Create("fubenleveltmp.png")
//	//if err != nil {
//	//	log.Fatalf("无法创建文件: %v", err)
//	//}
//	//defer file.Close()
//	//
//	//// 将图像编码为 PNG 格式并写入文件
//	////err = png.Encode(file, img)
//	//err = png.Encode(file, img)
//	//if err != nil {
//	//	log.Fatalf("保存图片失败: %v", err)
//	//}
//	//log.Println("截图成功，已保存为 fubenleveltmp.png")
//	for _, rectType := range nextMonitorList { // case的先后顺序有严格要求
//		switch rectType {
//		case defs.RectTypeIsYourTurn: // 以下每个分支下的return必不可少，不然可能造成混乱，多个监听
//			if utils.InSenseIsYourTurn(ctrl.hwnd, img) { // 监听：轮到你出手了
//				fmt.Printf("time:%s in your turn\n", time.Now().Format(time.TimeOnly))
//				ctrl.roundCount++
//				if ctrl.roundCount == 1 {
//					fmt.Printf("time:%s is first round\n", time.Now().Format(time.TimeOnly))
//					DoHandleFirstRoundInit(ctrl)
//				}
//				fmt.Printf("time:%s do fight handle\n", time.Now().Format(time.TimeOnly))
//				DoHandleFight(ctrl)
//				ctrl.monitorFightProcess()
//				return
//			}
//		case defs.RectTypeWinOrFail:
//			if utils.InSenseFightWin(ctrl.hwnd, img) { // 战斗结算页面-胜利，最后一关打完才有，小关是直接翻牌了
//				fmt.Printf("time:%s fight win\n", time.Now().Format(time.TimeOnly))
//				ctrl.winCount++
//				ctrl.monitorFubenFanCard()
//				return
//			} else if utils.InSenseFighFail(ctrl.hwnd, img) { // 战斗结算页面-失败，每一小关都有
//				fmt.Printf("time:%s fight fail\n", time.Now().Format(time.TimeOnly))
//				ctrl.failCount++
//				ctrl.fixFightCount = 0
//				ctrl.monitorFubenFanCard()
//				return
//			}
//		case defs.RectTypeFanCardSmall, defs.RectTypeFanCardBoss: // 小关翻牌、boss翻牌
//			if utils.InSenseFubenFanCardSmall(ctrl.hwnd, img) {
//				fmt.Printf("time:%s fan card small\n", time.Now().Format(time.TimeOnly))
//				// Todo add 翻牌逻辑
//				ctrl.monitorFubenRoom()
//				return
//			} else if utils.InSenseFubenFanCardBoss(ctrl.hwnd, img) {
//				fmt.Printf("time:%s fan card boss\n", time.Now().Format(time.TimeOnly))
//				// Todo add 翻牌逻辑
//				ctrl.monitorFubenRoom()
//				return
//			}
//		case defs.RectTypeFubenSelectText: // 有这个select必然重新选择副本并开始
//			if utils.InSenseFubenSelectText(ctrl.hwnd, img) {
//				fmt.Printf("time:%s select fuben\n", time.Now().Format(time.TimeOnly))
//				if ctrl.isChild {
//					if !utils.IsFubenRoomReady(ctrl.hwnd) { // 如果未准备则准备，准备后通知主号
//						utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightStart) // 这个也是准备按钮的位置
//						if err := instanceMgr.SetChildReady(ctrl.masterHWND, ctrl.hwnd, defs.ReadyStateOK); err != nil {
//							panic(err.Error())
//						}
//					}
//				} else {
//					ctrl.EnterFubenOnce(true)
//					ctrl.monitorFightProcess()
//				}
//				return
//			}
//		case defs.RectTypeFubenInviteAndChangeTeam: // 先判断select再判断，副本房间特征元素，这样能保证需要选择副本的时候必然可以选择副本 Todo 可能不能这么快就开始，要等小号都准备。小号也不能点击开始和选择副本
//			if utils.InSenseFubenInviteAndChangeTeam(ctrl.hwnd, img) {
//				fmt.Printf("time:%s direct enter fuben\n", time.Now().Format(time.TimeOnly))
//				if ctrl.isChild {
//					if !utils.IsFubenRoomReady(ctrl.hwnd) { // 如果未准备则准备，准备后通知主号
//						utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightStart) // 这个也是准备按钮的位置
//						if err := instanceMgr.SetChildReady(ctrl.masterHWND, ctrl.hwnd, defs.ReadyStateOK); err != nil {
//							panic(err.Error())
//						}
//					}
//				} else {
//					ctrl.EnterFubenOnce(false)
//					ctrl.monitorFightProcess()
//				}
//				return
//			}
//		}
//	}
//	ctrl.Monitor(nextMonitorList)
//}
//
//func (ctrl *ScriptCtrl) runFuben() {
//	if ctrl.isChild { // 小号的话只允许在副本房间内开启运行
//		sense := utils.GetSense(ctrl.hwnd)
//		if sense != defs.SenseTypeFubenRoom {
//			panic(fmt.Sprintf("小号请先进入副本房间内再开启脚本. sense:%d hwnd:%d", sense, ctrl.hwnd))
//		}
//		if err := InstanceMgr().AddChild(ctrl.masterHWND, ctrl.hwnd); err != nil {
//			panic(err.Error())
//		}
//		if !utils.IsFubenRoomReady(ctrl.hwnd) { // 如果未准备则准备，准备后通知主号
//			utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightStart) // 这个也是准备按钮的位置
//		}
//		// 走到这里来之后肯定是准备了
//		if err := instanceMgr.SetChildReady(ctrl.masterHWND, ctrl.hwnd, defs.ReadyStateOK); err != nil {
//			panic(err.Error())
//		}
//		ctrl.waitMasterStartFuben()
//		ctrl.monitorFightProcess() // 不能在我的线程去调用其他线程的方法
//	} else {
//		sense := utils.GetSense(ctrl.hwnd)
//		switch sense {
//		case defs.SenseTypeOther, defs.SenseTypeJinjiRoom, defs.SenseTypeJinjiHall:
//			utils.BackToIndexPage(ctrl.hwnd)
//			utils.ClickElementAfterLong(ctrl.hwnd, defs.ElementFubenHall)
//			utils.ClickElementAfterLong(ctrl.hwnd, defs.ElementFubenEnter)
//		case defs.SenseTypeFubenRoom:
//		case defs.SenseTypeFubenHall:
//			utils.ClickElementAfterLong(ctrl.hwnd, defs.ElementFubenEnter)
//		case defs.SenseTypeInFight, defs.SenseTypeFubenFightLoading: // 在战斗内的话直接开始监听战斗
//			ctrl.monitorFightProcess()
//		//case : Todo zhangzhihui more，增加更多场景的判断，比如在结算胜负判断时，在翻牌时可以监控其他的
//		default:
//			panic("请不要在当前场景下开启脚本")
//		}
//		ctrl.EnterFubenOnce(true)
//		ctrl.monitorFightProcess()
//	}
//}
//
//func (ctrl *ScriptCtrl) Run() {
//	go func() {
//		defer func() {
//			if r := recover(); r != nil {
//				// 这里退出的时候通知客户端显示退出的原因
//				fmt.Printf("hwnd:%d exit. reason:%v\n", ctrl.hwnd, r)
//			}
//		}()
//		if self.id > defs.FubenIDBegin && self.id < defs.FubenIDEnd {
//			ctrl.runFuben()
//		} else if self.id > defs.FubenIDJinjiBegin && self.id < defs.FubenIDJinjiEnd {
//			//ctrl.runJinji()
//		} else if self.id > defs.FuBenIDOtherBegin && self.id < defs.FuBenIDOtherEnd {
//			//ctrl.runOther()
//		}
//	}()
//}
