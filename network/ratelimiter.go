/*
Copyright (c) 2017 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package network offers various network tools
package network

import (
	"gitlab.ceriath.net/libs/goBlue/log"
	"sync"
	"time"
)

var mutex sync.Mutex
var tokenMutex sync.Mutex

//Ratelimiter implements a tokenbucket that supports bursts.
//A burst is an additional number of tokens that can be issued in a short time increasing the cooldown.
//Usually the ratelimiter will issue tokens on a constant rate, e.g. 2 per second. 
//when burst is enabled you can request e.g. 30 in 2 seconds but wait 15 seconds until issuing new tokens
type Ratelimiter struct {
	Name            string
	limit           int
	tokens          int
	resetAfter      time.Duration
	queueEnd        chan int
	resetSignal		chan int
	burstLimit      int
	burstTokens     int
	burstResetAfter time.Duration
	queue           []chan int
}

//Init intializes the ratelimiter with an identifier name, and a limit of tokens that are issued during a duration reset
func (rl *Ratelimiter) Init(name string, limit int, reset time.Duration) {
	rl.limit = limit
	rl.Name = name
	rl.tokens = limit
	rl.resetAfter = reset
	rl.queueEnd = make(chan int, 1)
	rl.resetSignal = make(chan int, 1)
}

//InitBurst inits a burst limit of tokens that are issued during duration reset
func (rl *Ratelimiter) InitBurst(limit int, reset time.Duration) {
	rl.burstLimit = limit - rl.limit
	rl.burstTokens = limit - rl.limit
	rl.burstResetAfter = reset
}

//return a channel, put that channel in the queue and on reset send a 1 to the next channel in queue
func (rl *Ratelimiter) Request(promote bool) chan int {
//	defer func() { fmt.Printf("%s\n", rl.queue) }()
	tokenMutex.Lock()
	if rl.tokens > 0 && rl.burstTokens == rl.burstLimit {
		//		fmt.Printf("not using burstmode\n")
		rl.tokens--
		tokenMutex.Unlock()
//		fmt.Printf("token used, remaining: %d\n", rl.tokens)
		go func(rl *Ratelimiter) {
			time.Sleep(rl.resetAfter)
			if rl.tokens < rl.limit {
				mutex.Lock()
				rl.tokens++
//				fmt.Printf("token added, remaining: %d\n", rl.tokens)
				if len(rl.queue) > 0 {
					rl.queue[0] <- 1
//					fmt.Printf("pop %s\n", rl.queue[0])
					rl.queue = rl.queue[1:]
				}
				rl.resetSignal <- 1
				mutex.Unlock()
			}
		}(rl)
		grant := make(chan int, 1)
		grant <- 1
		return grant
	} else if rl.burstTokens > 0 {
//		fmt.Printf("using burstmode\n")
		rl.burstTokens--
		tokenMutex.Unlock()
//		fmt.Printf("bursttoken used, remaining: %d\n", rl.burstTokens)
		go func(rl *Ratelimiter) {
			time.Sleep(rl.burstResetAfter)
			if rl.burstTokens < rl.burstLimit {
				mutex.Lock()
				rl.burstTokens++
//				fmt.Printf("bursttoken added, remaining: %d\n", rl.burstTokens)
				if len(rl.queue) > 0 {
					rl.queue[0] <- 1
//					fmt.Printf("pop %s\n", rl.queue[0])
					rl.queue = rl.queue[1:]
				}
				rl.resetSignal <- 1
				mutex.Unlock()
			}
		}(rl)
		grant := make(chan int, 1)
		grant <- 1
		return grant
	} else {
		grant := make(chan int, 1)
		mutex.Lock()
		if !promote {
			rl.queue = append(rl.queue, grant)
		} else {
			var tmpQ []chan int
			tmpQ = append(tmpQ, grant)
			tmpQ = append(tmpQ, rl.queue...)
			rl.queue = tmpQ
		}
		mutex.Unlock()
		tokenMutex.Unlock()

		if len(rl.queue) > rl.limit*3 {
			log.I("Warning: Ratelimiter", rl.Name, "'s queue is", len(rl.queue))
		}

		if len(rl.queue) == 0 {
			rl.queueEnd <- 1
		}
		<-rl.resetSignal
		<-rl.Request(promote)
		return grant
	}
}

func (rl *Ratelimiter) GetQueuesize() int {
	return len(rl.queue)
}

func (rl *Ratelimiter) WaitForQueue() {
	if len(rl.queue) <= 0 {
		return
	}
	<-rl.queueEnd
	return
}
