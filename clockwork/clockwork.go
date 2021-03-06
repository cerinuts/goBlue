/*
Copyright (c) 2018 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package clockwork runs scheduled tasks
package clockwork

import (
	"errors"
	"sync"
	"time"

	"code.cerinuts.io/libs/goBlue/log"
)

//AppName is the name of the application
const AppName string = "goBlue/clockwork"

//VersionMajor 0 means in development, >1 ensures compatibility with each minor version, but breakes with new major version
const VersionMajor string = "0"

//VersionMinor introduces changes that require a new version number. If the major version is 0, they are likely to break compatibility
const VersionMinor string = "1"

//VersionBuild is the type of this release. s(table), b(eta), d(evelopment), n(ightly)
const VersionBuild string = "s"

//FullVersion contains the full name and version of this package in a printable string
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild

// Clockwork singleton
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

//RepeatEvery runs fn every duration, if onStart is set, run first on function call
func (cw *_clockwork) RepeatEvery(d time.Duration, fn func(), onStart bool, id string) {
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
		for range time.NewTicker(d).C {
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

//RunAfter runs fn after duration
func (cw *_clockwork) RunAfter(d time.Duration, fn func(), id string) {
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

//InterruptTask interrupts stops task with id and prevents it from running (again)
func (cw *_clockwork) InterruptTask(id string) error {
	interrupt := interrupts[id]
	if interrupt == nil {
		return errors.New("Task not found:" + id)
	}
	close(interrupt)
	return nil
}

//Waits for all tasks to be done.
func (cw *_clockwork) WaitForFinish() {
	cw.cw.Waitgroup.Wait()
}
