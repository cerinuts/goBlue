package main

import (
	//"github.com/ceriath/goBlue/network"
//	"github.com/ceriath/goBlue/archium"
//	"github.com/ceriath/goBlue/log"
	//	"log"
	"github.com/ceriath/goBlue/clockwork"
	"time"
	//"github.com/ceriath/goBlue/network/client"
)

func test(){
	print("tick")
}

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
//	log.PrintToStdout = true
//	log.CurrentLevel = log.LevelDebug
//	log.CurrentLogFileBehaviour = log.LogfileBehaviourDaily
//	adl := new(archium.ArchiumDebugListener)
//	a := archium.ArchiumCore
//	a.Register(adl)
//	ev := archium.CreateEvent()
//	ev.EventType = "chat.abc.d"
//	ev.EventSource = "MAIN"
//	ev.Data["test"] = "abc"
//	a.FireEvent(*ev)
	cw := clockwork.Clockwork
	interrupt := cw.RepeatEvery(10 * time.Second, test, true)
	cw.RunAfter(50 * time.Second, test)
	time.Sleep(30 * time.Second)
	close(interrupt)
	cw.WaitForFinish()
}
