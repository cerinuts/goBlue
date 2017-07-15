/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Clockwork runs scheduled tasks
package clockwork

import (
	"time"
	"sync"
)

var Clockwork _clockwork

type _clockwork struct{
	cw *clockwork
}

type clockwork struct{
	Waitgroup sync.WaitGroup
}

func init(){
		if Clockwork.cw == nil {
		Clockwork.cw = new(clockwork)
	}
}

//run fn every duration, if onStart is set, run first on function call
func (cw *_clockwork) RepeatEvery(d time.Duration, fn func(), onStart bool) chan bool {
	interrupt := make(chan bool)
	if onStart {
		cw.cw.Waitgroup.Add(1)
		go func(){
			defer cw.cw.Waitgroup.Done()
			fn()
		}()
	}
	cw.cw.Waitgroup.Add(1)
	go func() {
		defer cw.cw.Waitgroup.Done()
		for range time.Tick(d) {
			select {
			case <-interrupt:
				return
			default:
				go fn()
			}
		}
	}()
	return interrupt
}

//runs fn after duration
func (cw *_clockwork) RunAfter(d time.Duration, fn func()) chan bool {
	interrupt := make(chan bool)
	run := make(chan bool)
	go func(){
		<- time.After(d)
		run <- true
	}()
	cw.cw.Waitgroup.Add(1)
	go func() {
		defer cw.cw.Waitgroup.Done()
		select {
		case <-interrupt:
			return
		case <-run:
			go fn()
		}
	}()
	return interrupt
}

//Waits for all tasks to be done.
func (cw *_clockwork) WaitForFinish(){
	cw.cw.Waitgroup.Wait()
}
