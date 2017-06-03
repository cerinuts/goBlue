/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

package archium

import (
	"github.com/ceriath/goBlue/log"
	"strings"
)

//Types

type ArchiumEvent struct {
	EventType, EventSource string
	Data                   map[string]string
}

type ArchiumEventListener interface {
	Trigger(ae ArchiumEvent)
	GetTypes() []string
}

type _ArchiumCore struct {
	ac *archiumCore
}

type archiumCore struct {
	listener []ArchiumEventListener
}

// Core

var ArchiumCore _ArchiumCore

func init() {
	ArchiumCore.ac = new(archiumCore)
}

func (core *_ArchiumCore) FireEvent(ev ArchiumEvent) {
	for _, el := range core.ac.listener {
		for t := range el.GetTypes() {
			if checkTypes(t, ev.EventType) {
				el.Trigger(ev)
			}
		}
	}
}

//Core Util

func checkTypes(lType, eType string) bool {
	lType = strings.ToLower(lType)
	eType = strings.ToLower(eType)
	lTypeSplit := strings.Split(lType, ".")
	if lTypeSplit[0] == "*" {
		return true
	}
	eTypeSplit := strings.Split(eType, ".")
	index := 0
	for len(lTypeSplit) > index && len(eTypeSplit) > index && (lTypeSplit[index] == eTypeSplit[index] || lTypeSplit[index] == "*") {
		//everything matches
		if index == len(lTypeSplit)-1 && index == len(eTypeSplit)-1 {
			return true
		}
		//next is * and its the last one
		if len(lTypeSplit)-1 > index && lTypeSplit[index+1] == "*" && len(lTypeSplit)-1 == index+1 {
			return true
		}
		//next is * and not the last one
		if len(lTypeSplit)-1 > index && lTypeSplit[index+1] == "*" {
			index++
			continue
		}
		//otherwise
		if len(lTypeSplit)-1 == index {
			break
		}
		index++
	}
	return false
}

func (core *_ArchiumCore) Register(al ArchiumEventListener) {
	core.ac.listener = append(core.ac.listener, al)
}

// Events

func CreateEvent(mapsize int) *ArchiumEvent {
	ev := new(ArchiumEvent)
	ev.Data = make(map[string]string, mapsize)
	return ev
}

// Debug Listener

type ArchiumDebugListener struct {
}

func (adl *ArchiumDebugListener) Trigger(ae ArchiumEvent) {
	mapstr := ""
	for k, v := range ae.Data {
		mapstr = mapstr + " --- " + k + ":" + v
	}
	log.D(ae.EventType, ae.EventSource, mapstr)
}

func (adl *ArchiumDebugListener) GetTypes() string {
	return []string{"*"}
}
