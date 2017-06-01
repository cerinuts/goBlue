package archium

import (
	"strings"
)

//Types

type ArchiumEvent struct{
	EventType, EventSource string
	Data map[string]string
}

type ArchiumEventListener interface{
	Trigger(ae ArchiumEvent)
	GetType() string
}

type _ArchiumCore struct{
	ac *archiumCore
}



type archiumCore struct{
	listener []ArchiumEventListener
}

// Core

var ArchiumCore _ArchiumCore

func init(){
	ArchiumCore.ac = new(archiumCore)
}

func (core *_ArchiumCore) FireEvent(ev ArchiumEvent){
	for _,el := range core.ac.listener{
			if (checkTypes(el.GetType(), ev.EventType)){
				el.Trigger(ev)
			}
	}
}

//Core Util

func checkTypes(lType, eType string) bool{
//	lType = strings.ToLower(lType)
//	eType = strings.ToLower(eType)
	lTypeSplit := strings.Split(lType, ".")
	if(lTypeSplit[0] == "*"){
		return true
	}
	eTypeSplit := strings.Split(eType, ".")
	index := 0
	for(lTypeSplit[index] == eTypeSplit[index]){
			if(index == len(lTypeSplit)-1 && index == len(eTypeSplit)-1){
				return true
			}
			if(len(lTypeSplit)-1 > index && lTypeSplit[index+1] == "*"){
				return true
			}
			if(len(lTypeSplit)-1 == index){
				break
			}
			index++
		}
	return false
}

func (core *_ArchiumCore) Register(al ArchiumEventListener){
	core.ac.listener = append(core.ac.listener, al)
}

// Events

func CreateEvent(mapsize int) *ArchiumEvent{
	ev := new(ArchiumEvent)
	ev.Data = make(map[string]string, mapsize)
	return ev
}

// Debug Listener

type ArchiumDebugListener struct{
	
}

func (adl *ArchiumDebugListener) Trigger(ae ArchiumEvent){
	mapstr := ""
	//create string of whole map
	println(ae.EventType, ae.EventSource, mapstr)
}

func (adl *ArchiumDebugListener) GetType() string{
	return "*"
}

