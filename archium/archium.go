/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Archium is a simple PubSub-Service
package archium

import (
	"fmt"
	"gitlab.ceriath.net/libs/goBlue/log"
	"strings"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goBlue/archium", "0", "1", "s"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

//Types

//An ArchiumEvent is an event that is fired by something. Data is a map<string, anything>,
//EventType is the type of the Event, often referred to as "topic"
type ArchiumEvent struct {
	EventType, EventSource string
	Data                   map[string]interface{}
}

//Interface for an Event-Listener. Trigger is activated if type of an occured event matches any type returned
//by GetTypes. GetTypes also can return static strings.
type ArchiumEventListener interface {
	Trigger(ae ArchiumEvent)
	GetTypes() []string
}

//internal core wrapper for singleton design
type _ArchiumCore struct {
	ac *archiumCore
}

//the actual single core
type archiumCore struct {
	listener []ArchiumEventListener
}

// The Archium Core - Singleton
var ArchiumCore _ArchiumCore

//creates new singleton on init
func init() {
	if ArchiumCore.ac == nil {
		ArchiumCore.ac = new(archiumCore)
	}
}

//Fires an event, checks all listeners for their type and fires if necessary
func (core *_ArchiumCore) FireEvent(ev ArchiumEvent) {
	for _, el := range core.ac.listener {
		for _, t := range el.GetTypes() {
			if checkTypes(t, ev.EventType) {
				el.Trigger(ev)
			}
		}
	}
}

//Core Util

//checks if two types match, considers wildcards
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

//Add a new Listener
func (core *_ArchiumCore) Register(al ArchiumEventListener) {
	core.ac.listener = append(core.ac.listener, al)
}

//Add a new Listener
func (core *_ArchiumCore) Deregister(al ArchiumEventListener) {
	var idx int
	for i, el := range core.ac.listener{
		if(el == al){
			idx = i
		}
	}
	core.ac.listener  = append(core.ac.listener[:idx], core.ac.listener[idx+1:]...) 
}

// Events

//Create a new event - does NOT fire it!
func CreateEvent() *ArchiumEvent {
	ev := new(ArchiumEvent)
	ev.Data = make(map[string]interface{})
	return ev
}

// Debug Listener

//The debugListener is a basic listener which listens to everything and logs it to debug
type ArchiumDebugListener struct {
}

func (adl *ArchiumDebugListener) Trigger(ae ArchiumEvent) {
	mapstr := ""
	for k, v := range ae.Data {
		mapstr = mapstr + " --- " + k + ":" + fmt.Sprintf("%b", v)
	}
	log.D(ae.EventType, ae.EventSource, mapstr)
}

func (adl *ArchiumDebugListener) GetTypes() []string {
	return []string{"*"}
}
