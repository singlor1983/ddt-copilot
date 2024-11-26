package core

import (
	"ddt-copilot/data"
	"ddt-copilot/defs"
	"ddt-copilot/utils"
	"fmt"
	"github.com/lxn/win"
	"sync"
	"time"
)

var (
	monitorListFightProcess = []defs.DDTElementRectType{
		defs.DDTElementRectTypeIsYourTurn,
		defs.DDTElementRectTypeWinOrFail,
	}
	monitorListFubenFanCard = []defs.DDTElementRectType{
		defs.DDTElementRectTypeFanCardSmall,
		defs.DDTElementRectTypeFanCardBoss,
	}
	monitorListFubenRoom = []defs.DDTElementRectType{
		defs.DDTElementRectTypeFubenSelectText,
		defs.DDTElementRectTypeFubenInviteAndChangeTeam,
	}
)

var (
	instanceMgr *ScriptCtrlMgr
)

// InstanceMgr 单例
func InstanceMgr() *ScriptCtrlMgr {
	if instanceMgr == nil {
		instanceMgr = NewInstanceMgr()
	}
	return instanceMgr
}

func NewInstanceMgr() *ScriptCtrlMgr {
	return &ScriptCtrlMgr{items: make(map[win.HWND]*ScriptCtrl)}
}

type ChildNode struct {
	hwnd       win.HWND
	readyState defs.ReadyState // 是否已准备，副本结束时需要重置状态
}

type ScriptCtrl struct {
	hwnd win.HWND // 窗口句柄
	// 以下选项启动后不能修改，界面置灰
	scriptIndex defs.FubenID    // 脚本功能索引
	lvIndex     defs.FubenLevel // 副本难度-不是所有功能都和难度有关
	isChild     bool            // 是否为小号
	isBossFight bool            // 是否为boss战
	masterHWND  win.HWND        // 所属主号句柄

	mu sync.RWMutex

	// 内存数据
	fightCount          int                     // 已挑战了多少次关卡
	fixFightCount       int                     // 修正的挑战次数，战斗失败会把这个置为0
	winCount, failCount int                     // 关卡胜利失败次数
	roundCount          int                     // 第几次自己的回合出手
	initPosition        defs.FightInnerPosition // 初始位置
	childHWNDs          map[win.HWND]*ChildNode // 小号列表，小号开始运行时加入主号的列表

}

type ScriptCtrlMgr struct {
	items map[win.HWND]*ScriptCtrl
}

func (mgr *ScriptCtrlMgr) AddCtrl(ctrl *ScriptCtrl) {
	if mgr.items == nil || ctrl == nil {
		return
	}
	mgr.items[ctrl.hwnd] = ctrl
}

func (mgr *ScriptCtrlMgr) GetCtrl(hwnd win.HWND) *ScriptCtrl {
	if mgr.items == nil {
		return nil
	}
	v := mgr.items[hwnd]
	return v
}

func (mgr *ScriptCtrlMgr) GetChild(master, child win.HWND) (*ChildNode, error) {
	if mgr.items == nil {
		return nil, fmt.Errorf("mgr.items == nil")
	}
	m := mgr.items[master]
	if m == nil {
		return nil, fmt.Errorf("m == nil. master:%d", master)
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	v := m.childHWNDs[child]
	return v, nil
}

func (mgr *ScriptCtrlMgr) AddChild(master, child win.HWND) error {
	if mgr.items == nil {
		return fmt.Errorf("mgr.items == nil")
	}
	m := mgr.items[master]
	if m == nil {
		return fmt.Errorf("m == nil. master:%d", master)
	}
	if len(m.childHWNDs) >= 3 { // 最多3个小号
		return fmt.Errorf("len(m.childHWNDs) >= 3. list:%v", m.childHWNDs)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.childHWNDs == nil {
		m.childHWNDs = make(map[win.HWND]*ChildNode)
	}
	if m.childHWNDs[child] != nil {
		return fmt.Errorf("repeat child:%d", child)
	}
	m.childHWNDs[child] = &ChildNode{
		hwnd:       child,
		readyState: defs.ReadyStateNo,
	}
	return nil
}

func (mgr *ScriptCtrlMgr) DelChild(master, child win.HWND) error {
	if mgr.items == nil {
		return fmt.Errorf("mgr.items == nil")
	}
	m := mgr.items[master]
	if m == nil {
		return fmt.Errorf("m == nil. master:%d", master)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.childHWNDs, child)
	return nil
}

func (mgr *ScriptCtrlMgr) SetChildReady(master, child win.HWND, state defs.ReadyState) error {
	v, _ := mgr.GetChild(master, child)
	if v == nil {
		return fmt.Errorf("not found child:%d", child)
	}
	v.readyState = state
	return nil
}

func NewScriptCtrl(hwnd, master win.HWND, scriptIndex defs.FubenID, lvIndex defs.FubenLevel, isChild, isBossFight bool) *ScriptCtrl {
	return &ScriptCtrl{hwnd: hwnd, scriptIndex: scriptIndex, lvIndex: lvIndex, isChild: isChild, isBossFight: isBossFight, masterHWND: master}
}

func (ctrl *ScriptCtrl) SelectFubenMap() {
	if err := utils.SelectFubenMap(ctrl.hwnd, ctrl.scriptIndex, ctrl.lvIndex, ctrl.isBossFight, data.GConfig.GetFubenPosition()); err != nil {
		panic(err.Error())
	}
}

func (ctrl *ScriptCtrl) OnBeforeStartFuben() {
	// Todo zhangzhihi 还有个通用before，比如喂宠物粮
	DoHandleBeforeStart(ctrl)
}

// waitChildReady 开始战斗之前检查小号是否都已经准备
func (ctrl *ScriptCtrl) waitChildReady() {
	if ctrl.isChild { // 本身就是小号无需监听
		return
	}
	for {
		fmt.Printf("time:%s check child ready\n", time.Now().Format(time.TimeOnly))
		allReady := true
		for _, node := range ctrl.childHWNDs {
			if node == nil {
				continue
			}
			if node.readyState == defs.ReadyStateNo {
				allReady = false
				break
			}
		}
		if allReady {
			fmt.Printf("time:%s child all ready\n", time.Now().Format(time.TimeOnly))
			return
		}
		time.Sleep(time.Second) // 一秒检查一次
	}
}

func (ctrl *ScriptCtrl) waitMasterStartFuben() {
	if !ctrl.isChild {
		return
	}
	for {
		master := instanceMgr.GetCtrl(ctrl.masterHWND)
		if master == nil {
			panic("master nil")
		}
		if master.fightCount > ctrl.fightCount {
			ctrl.fightCount = master.fightCount // Todo 用这个来判断是否开启了副本，有待商榷，什么清楚？
			return
		}
	}
}

func (ctrl *ScriptCtrl) StartFuben() {
	ctrl.waitChildReady()
	ctrl.fightCount++ // 重新启动脚本自然会重置【重启脚本需要重新初始化ScriptCtrl】
	ctrl.fixFightCount++
	ctrl.roundCount = 0 // 每次新开战斗要重置回合次数
	utils.ClickElementAfterMid(ctrl.hwnd, defs.ElementFightStart)
	utils.ClickElementAfterMid(ctrl.hwnd, defs.ElementFubenFightStartAck) // 多人战斗的时候没有确认框，但是点击一下也不影响
	fmt.Printf("time:%s start fuben: %v\n", time.Now().Format(time.TimeOnly), ctrl)
}

// EnterFubenOnce 在房间内执行一次进入副本的操作
func (ctrl *ScriptCtrl) EnterFubenOnce(needSelectFuben bool) {
	if needSelectFuben { // 多关卡的是不需要再次选择副本的,实际上点击一次也没关系，没有选项
		ctrl.SelectFubenMap()
	}
	ctrl.OnBeforeStartFuben()
	ctrl.StartFuben()
}

func (ctrl *ScriptCtrl) monitorFightProcess() {
	ctrl.Monitor(monitorListFightProcess)
}

func (ctrl *ScriptCtrl) monitorFubenRoom() {
	ctrl.Monitor(monitorListFubenRoom)
}

func (ctrl *ScriptCtrl) monitorFubenFanCard() {
	ctrl.Monitor(monitorListFubenFanCard)
}

// Monitor 每秒只截屏一次，然后从这张图里面再截取对应的部分，而不是每次重新截屏
func (ctrl *ScriptCtrl) Monitor(nextMonitorList []defs.DDTElementRectType) {
	time.Sleep(time.Millisecond * 500)
	img, err := utils.CaptureWindowLight(ctrl.hwnd, nil, true)
	if err != nil {
		panic(err.Error())
	}
	//// 创建文件并保存为 PNG 格式
	//file, err := os.Create("fubenleveltmp.png")
	//if err != nil {
	//	log.Fatalf("无法创建文件: %v", err)
	//}
	//defer file.Close()
	//
	//// 将图像编码为 PNG 格式并写入文件
	////err = png.Encode(file, img)
	//err = png.Encode(file, img)
	//if err != nil {
	//	log.Fatalf("保存图片失败: %v", err)
	//}
	//log.Println("截图成功，已保存为 fubenleveltmp.png")
	for _, rectType := range nextMonitorList { // case的先后顺序有严格要求
		switch rectType {
		case defs.DDTElementRectTypeIsYourTurn: // 以下每个分支下的return必不可少，不然可能造成混乱，多个监听
			if utils.InSenseIsYourTurn(ctrl.hwnd, img) { // 监听：轮到你出手了
				fmt.Printf("time:%s in your turn\n", time.Now().Format(time.TimeOnly))
				ctrl.roundCount++
				if ctrl.roundCount == 1 {
					fmt.Printf("time:%s is first round\n", time.Now().Format(time.TimeOnly))
					DoHandleFirstRoundInit(ctrl)
				}
				fmt.Printf("time:%s do fight handle\n", time.Now().Format(time.TimeOnly))
				DoHandleFight(ctrl)
				ctrl.monitorFightProcess()
				return
			}
		case defs.DDTElementRectTypeWinOrFail:
			if utils.InSenseFightWin(ctrl.hwnd, img) { // 战斗结算页面-胜利，最后一关打完才有，小关是直接翻牌了
				fmt.Printf("time:%s fight win\n", time.Now().Format(time.TimeOnly))
				ctrl.winCount++
				ctrl.monitorFubenFanCard()
				return
			} else if utils.InSenseFighFail(ctrl.hwnd, img) { // 战斗结算页面-失败，每一小关都有
				fmt.Printf("time:%s fight fail\n", time.Now().Format(time.TimeOnly))
				ctrl.failCount++
				ctrl.fixFightCount = 0
				ctrl.monitorFubenFanCard()
				return
			}
		case defs.DDTElementRectTypeFanCardSmall, defs.DDTElementRectTypeFanCardBoss: // 小关翻牌、boss翻牌
			if utils.InSenseFubenFanCardSmall(ctrl.hwnd, img) {
				fmt.Printf("time:%s fan card small\n", time.Now().Format(time.TimeOnly))
				// Todo add 翻牌逻辑
				ctrl.monitorFubenRoom()
				return
			} else if utils.InSenseFubenFanCardBoss(ctrl.hwnd, img) {
				fmt.Printf("time:%s fan card boss\n", time.Now().Format(time.TimeOnly))
				// Todo add 翻牌逻辑
				ctrl.monitorFubenRoom()
				return
			}
		case defs.DDTElementRectTypeFubenSelectText: // 有这个select必然重新选择副本并开始
			if utils.InSenseFubenSelectText(ctrl.hwnd, img) {
				fmt.Printf("time:%s select fuben\n", time.Now().Format(time.TimeOnly))
				if ctrl.isChild {
					if !utils.IsFubenRoomReady(ctrl.hwnd) { // 如果未准备则准备，准备后通知主号
						utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightStart) // 这个也是准备按钮的位置
						if err := instanceMgr.SetChildReady(ctrl.masterHWND, ctrl.hwnd, defs.ReadyStateOK); err != nil {
							panic(err.Error())
						}
					}
				} else {
					ctrl.EnterFubenOnce(true)
					ctrl.monitorFightProcess()
				}
				return
			}
		case defs.DDTElementRectTypeFubenInviteAndChangeTeam: // 先判断select再判断，副本房间特征元素，这样能保证需要选择副本的时候必然可以选择副本 Todo 可能不能这么快就开始，要等小号都准备。小号也不能点击开始和选择副本
			if utils.InSenseFubenInviteAndChangeTeam(ctrl.hwnd, img) {
				fmt.Printf("time:%s direct enter fuben\n", time.Now().Format(time.TimeOnly))
				if ctrl.isChild {
					if !utils.IsFubenRoomReady(ctrl.hwnd) { // 如果未准备则准备，准备后通知主号
						utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightStart) // 这个也是准备按钮的位置
						if err := instanceMgr.SetChildReady(ctrl.masterHWND, ctrl.hwnd, defs.ReadyStateOK); err != nil {
							panic(err.Error())
						}
					}
				} else {
					ctrl.EnterFubenOnce(false)
					ctrl.monitorFightProcess()
				}
				return
			}
		}
	}
	ctrl.Monitor(nextMonitorList)
}

func (ctrl *ScriptCtrl) runFuben() {
	if ctrl.isChild { // 小号的话只允许在副本房间内开启运行
		sense := utils.GetSense(ctrl.hwnd)
		if sense != defs.SenseTypeFubenRoom {
			panic(fmt.Sprintf("小号请先进入副本房间内再开启脚本. sense:%d hwnd:%d", sense, ctrl.hwnd))
		}
		if err := InstanceMgr().AddChild(ctrl.masterHWND, ctrl.hwnd); err != nil {
			panic(err.Error())
		}
		if !utils.IsFubenRoomReady(ctrl.hwnd) { // 如果未准备则准备，准备后通知主号
			utils.ClickElementAfterShort(ctrl.hwnd, defs.ElementFightStart) // 这个也是准备按钮的位置
		}
		// 走到这里来之后肯定是准备了
		if err := instanceMgr.SetChildReady(ctrl.masterHWND, ctrl.hwnd, defs.ReadyStateOK); err != nil {
			panic(err.Error())
		}
		ctrl.waitMasterStartFuben()
		ctrl.monitorFightProcess() // 不能在我的线程去调用其他线程的方法
	} else {
		sense := utils.GetSense(ctrl.hwnd)
		switch sense {
		case defs.SenseTypeOther, defs.SenseTypeJinjiRoom, defs.SenseTypeJinjiHall:
			utils.BackToIndexPage(ctrl.hwnd)
			utils.ClickElementAfterLong(ctrl.hwnd, defs.ElementFubenHall)
			utils.ClickElementAfterLong(ctrl.hwnd, defs.ElementFubenEnter)
		case defs.SenseTypeFubenRoom:
		case defs.SenseTypeFubenHall:
			utils.ClickElementAfterLong(ctrl.hwnd, defs.ElementFubenEnter)
		case defs.SenseTypeInFight, defs.SenseTypeFubenFightLoading: // 在战斗内的话直接开始监听战斗
			ctrl.monitorFightProcess()
		//case : Todo zhangzhihui more，增加更多场景的判断，比如在结算胜负判断时，在翻牌时可以监控其他的
		default:
			panic("请不要在当前场景下开启脚本")
		}
		ctrl.EnterFubenOnce(true)
		ctrl.monitorFightProcess()
	}
}

func (ctrl *ScriptCtrl) Run() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// 这里退出的时候通知客户端显示退出的原因
				fmt.Printf("hwnd:%d exit. reason:%v\n", ctrl.hwnd, r)
			}
		}()
		if ctrl.scriptIndex > defs.FubenIDBegin && ctrl.scriptIndex < defs.FubenIDEnd {
			ctrl.runFuben()
		} else if ctrl.scriptIndex > defs.FubenIDJinjiBegin && ctrl.scriptIndex < defs.FubenIDJinjiEnd {
			//ctrl.runJinji()
		} else if ctrl.scriptIndex > defs.FuBenIDOtherBegin && ctrl.scriptIndex < defs.FuBenIDOtherEnd {
			//ctrl.runOther()
		}
	}()
}
