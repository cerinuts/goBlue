/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Clockwork runs scheduled tasks
package clockwork

import (
	"github.com/ceriath/goBlue/log"
	"sync"
	"time"
	"errors"
)

const AppName, VersionMajor, VersionMinor, VersionBuild string = "goBlue/Clockwork", "1", "0", "d"
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

var Clockwork _clockwork
var interrupts map[string]chan bool

type _clockwork struct {
	cw *clockwork
}

type clockwork struct {
	Waitgroup sync.WaitGroup
}

func init() {
	if Clockwork.cw == nil {
		Clockwork.cw = new(clockwork)
		interrupts = make(map[string]chan bool)
	}
}

//run fn every duration, if onStart is set, run first on function call
func (cw *_clockwork) RepeatEvery(d time.Duration, fn func(), onStart bool, id string){
	interrupt := make(chan bool)
	if onStart {
		cw.cw.Waitgroup.Add(1)
		go func() {
			defer cw.cw.Waitgroup.Done()
			fn()
		}()
	}
	cw.cw.Waitgroup.Add(1)
	go func() {
		defer cw.cw.Waitgroup.Done()
		defer delete(interrupts, id)
		for range time.Tick(d) {
			select {
			case <-interrupt:
				return
			default:
				go fn()
			}
		}
	}()
	log.I("Added repetitive task", id, "every", d)
	interrupts[id] = interrupt
}

//runs fn after duration
func (cw *_clockwork) RunAfter(d time.Duration, fn func(), id string){
	interrupt := make(chan bool)
	run := make(chan bool)
	go func() {
		<-time.After(d)
		run <- true
	}()
	cw.cw.Waitgroup.Add(1)
	go func() {
		defer cw.cw.Waitgroup.Done()
		defer delete(interrupts, id)
		select {
		case <-interrupt:
			return
		case <-run:
			go fn()
		}
	}()
	log.I("Added scheduled task", id, "after", d)
	interrupts[id] = interrupt
}

func (cw *_clockwork) InterruptTask(id string) error{
	interrupt := interrupts[id]
	if interrupt == nil{
		return errors.New("Task not found:" + id)
	}
	close(interrupt)
	return nil
}

//Waits for all tasks to be done.
func (cw *_clockwork) WaitForFinish() {
	cw.cw.Waitgroup.Wait()
}
