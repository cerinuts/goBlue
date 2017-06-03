package main

import (
	//"github.com/ceriath/goBlue/network"
	"github.com/ceriath/goBlue/archium"
	"github.com/ceriath/goBlue/log"
	//	"log"
	"time"
	//"github.com/ceriath/goBlue/network/client"
)

func main() {
	//	s := new(network.Server)
	//	go s.Create("localhost", "12345")
	//	time.Sleep(2 * time.Second)
	//	c := new(network.Client)
	//	err := c.Connect("localhost", "12345")
	//	if(err != nil){
	//		println("Error connecting client")
	//	}
	//	c.Sendln("hi")
	//	println("sent.")
	//	time.Sleep(2 * time.Second)
	log.PrintToStdout = true
	log.CurrentLevel = log.LevelDebug
	log.CurrentLogFileBehaviour = log.LogfileBehaviourDaily
	adl := new(archium.ArchiumDebugListener)
	a := archium.ArchiumCore
	a.Register(adl)
	ev := archium.CreateEvent(1)
	ev.EventType = "chat.abc.d"
	ev.EventSource = "MAIN"
	ev.Data["test"] = "abc"
	a.FireEvent(*ev)
	time.Sleep(2 * time.Minute)
}
