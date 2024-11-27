package core

import (
	"github.com/lxn/win"
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

type ScriptCtrlMgr struct {
	items map[win.HWND]*ScriptCtrl
}

func (mgr *ScriptCtrlMgr) AddCtrl(ctrl *ScriptCtrl) {
	if mgr.items == nil || ctrl == nil {
		return
	}
	mgr.items[ctrl.hwnd] = ctrl
}

func (mgr *ScriptCtrlMgr) DelCtrl(ctrl *ScriptCtrl) {
	if mgr.items == nil || ctrl == nil {
		return
	}
	delete(mgr.items, ctrl.hwnd)
}

func (mgr *ScriptCtrlMgr) GetCtrl(hwnd win.HWND) *ScriptCtrl {
	if mgr.items == nil {
		return nil
	}
	v := mgr.items[hwnd]
	return v
}
