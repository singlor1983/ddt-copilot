package core

import (
	"ddt-copilot/defs"
	"fmt"
)

type beforeStartCB func(ctrl *ScriptCtrl)
type fightCB func(ctrl *ScriptCtrl)
type firstRoundInitCB func(ctrl *ScriptCtrl)

type handleNode struct {
	id               defs.FubenID
	cbBeforeStart    beforeStartCB
	cbFight          fightCB
	cbFirstRoundInit firstRoundInitCB
}

var handles map[defs.FubenID]*handleNode

func RegisterHandle(id defs.FubenID, cbBeforeStart beforeStartCB, cbFirstRoundInit firstRoundInitCB, cbFight fightCB) {
	if handles == nil {
		handles = make(map[defs.FubenID]*handleNode)
	}
	if _, ok := handles[id]; ok {
		panic(fmt.Sprintf("handle:%d already exist", id))
	}
	handles[id] = &handleNode{
		id:               id,
		cbBeforeStart:    cbBeforeStart,
		cbFight:          cbFight,
		cbFirstRoundInit: cbFirstRoundInit,
	}
}

func DoHandleBeforeStart(ctrl *ScriptCtrl) {
	if handles == nil {
		panic("handles nil")
	}
	v := handles[ctrl.scriptIndex]
	if v == nil {
		panic(fmt.Sprintf("not found id:%d", v.id))
	}
	if v.cbBeforeStart != nil {
		v.cbBeforeStart(ctrl)
	}
}

func DoHandleFight(ctrl *ScriptCtrl) {
	if handles == nil {
		panic("handles nil")
	}
	v := handles[ctrl.scriptIndex]
	if v == nil {
		panic(fmt.Sprintf("not found id:%d", v.id))
	}
	if v.cbFight != nil {
		v.cbFight(ctrl)
	}
}

func DoHandleFirstRoundInit(ctrl *ScriptCtrl) {
	if handles == nil {
		panic("handles nil")
	}
	v := handles[ctrl.scriptIndex]
	if v == nil {
		panic(fmt.Sprintf("not found id:%d", v.id))
	}
	if v.cbFirstRoundInit != nil {
		v.cbFirstRoundInit(ctrl)
	}
}

func Init() {
	RegisterHandle(defs.FubenIDMY, OnBeforeStartMY, OnFirstRoundInitMY, OnFightMY)
	RegisterHandle(defs.FubenIDXSZY, nil, nil, OnFightXSZY)
	RegisterHandle(defs.FubenIDBG, nil, OnFirstRoundInitBG, OnFightBG)
}
