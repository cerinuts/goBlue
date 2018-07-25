/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package archium is a simple PubSub-Service
package archium

import (
	"fmt"
	"strings"

	"code.cerinuts.io/libs/goBlue/log"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goBlue/archium", "0", "2", "s"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

//Types

//Event is an event that is fired by something. Data is a map<string, anything>,
//EventType is the type of the Event, often referred to as "topic"
type Event struct {
	EventType, EventSource string
	Data                   map[string]interface{}
}

//EventListener provides an interface for custom listeners. Trigger is activated if type of an occured event matches any type returned
//by GetTypes. GetTypes also can return static strings.
type EventListener interface {
	Trigger(ae Event)
	GetTypes() []string
}

//internal core wrapper for singleton design
type _ArchiumCore struct {
	ac *archiumCore
}

//the actual single core
type archiumCore struct {
	listener []EventListener
}

// ArchiumCore - Singleton
var ArchiumCore _ArchiumCore

//creates new singleton on init
func init() {
	if ArchiumCore.ac == nil {
		ArchiumCore.ac = new(archiumCore)
	}
}

//FireEvent fires an event, checks all listeners for their type and fires if necessary
func (core *_ArchiumCore) FireEvent(ev Event) {
	for _, el := range core.ac.listener {
		for _, t := range el.GetTypes() {
			if checkTypes(t, ev.EventType) {
				el.Trigger(ev)
			}
		}
	}
}

//Core Util

//checkTypes checks if two types match, considers wildcards
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

//Register adds a new Listener
func (core *_ArchiumCore) Register(al EventListener) {
	core.ac.listener = append(core.ac.listener, al)
}

//Deregister removes a listener
func (core *_ArchiumCore) Deregister(al EventListener) {
	var idx int
	for i, el := range core.ac.listener {
		if el == al {
			idx = i
		}
	}
	core.ac.listener = append(core.ac.listener[:idx], core.ac.listener[idx+1:]...)
}

// Events

//CreateEvent creates a new event - does NOT fire it!
func CreateEvent() *Event {
	ev := new(Event)
	ev.Data = make(map[string]interface{})
	return ev
}

// Debug Listener

//DebugListener is a basic listener which listens to everything and logs it to debug
type DebugListener struct {
}

//Trigger of the DebugListener logs everything on LogLevel debug
func (adl *DebugListener) Trigger(ae Event) {
	mapstr := ""
	for k, v := range ae.Data {
		mapstr = mapstr + " --- " + k + ":" + fmt.Sprintf("%b", v)
	}
	log.D(ae.EventType, ae.EventSource, mapstr)
}

//GetTypes returns a single wildcard for everything to catch every single Event
func (adl *DebugListener) GetTypes() []string {
	return []string{"*"}
}
