package core

import (
	"ddt-copilot/defs"
	"fmt"
)

type beforeStartCB func(ctrl *ScriptCtrl)
type fightCB func(ctrl *ScriptCtrl)
type firstRoundInitCB func(ctrl *ScriptCtrl)

type handleNode struct {
	id               defs.FunctionID
	cbBeforeStart    beforeStartCB
	cbFight          fightCB
	cbFirstRoundInit firstRoundInitCB
}

var handles map[defs.FunctionID]*handleNode

func RegisterHandle(id defs.FunctionID, cbBeforeStart beforeStartCB, cbFirstRoundInit firstRoundInitCB, cbFight fightCB) {
	if handles == nil {
		handles = make(map[defs.FunctionID]*handleNode)
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

func findHandle(id defs.FunctionID) (*handleNode, error) {
	if handles == nil {
		return nil, fmt.Errorf("handles nil")
	}
	v := handles[id]
	if v == nil {
		return nil, fmt.Errorf("not found id:%d", v.id)
	}
	return v, nil
}

func DoHandleBeforeStart(ctrl *ScriptCtrl) {
	handler, err := findHandle(ctrl.id)
	if err != nil {
		panic(err)
	}
	if handler.cbBeforeStart != nil {
		handler.cbBeforeStart(ctrl)
	}
}

func DoHandleFight(ctrl *ScriptCtrl) {
	handler, err := findHandle(ctrl.id)
	if err != nil {
		panic(err)
	}
	if handler.cbFight != nil {
		handler.cbFight(ctrl)
	}
}

func DoHandleFirstRoundInit(ctrl *ScriptCtrl) {
	handler, err := findHandle(ctrl.id)
	if err != nil {
		panic(err)
	}
	if handler.cbFirstRoundInit != nil {
		handler.cbFirstRoundInit(ctrl)
	}
}

func Init() {
	RegisterHandle(defs.FunctionIDCustomFuben, nil, nil, OnBattleCustom)
	RegisterHandle(defs.FunctionIDMaYiGeneral, nil, OnRoundInitMaYiGeneral, OnBattleMaYiGeneral)
}
