package network

import (
	//	"fmt"
	"github.com/ceriath/goBlue/log"
	"time"
)

type Ratelimiter struct {
	Name            string
	limit           int
	tokens          int
	resetAfter      time.Duration
	resetSignal     chan int
	queueSize       int
	queueEnd        chan int
	burstLimit      int
	burstTokens     int
	burstResetAfter time.Duration
}

func (rl *Ratelimiter) Init(name string, limit int, reset time.Duration) {
	rl.limit = limit
	rl.Name = name
	rl.tokens = limit
	rl.resetAfter = reset
	rl.resetSignal = make(chan int)
	rl.queueEnd = make(chan int)
	rl.queueSize = 0
}

func (rl *Ratelimiter) InitBurst(limit int, reset time.Duration) {
	rl.burstLimit = limit - rl.limit
	rl.burstTokens = limit - rl.limit
	rl.burstResetAfter = reset
}

func (rl *Ratelimiter) Request() {
	if rl.tokens > 0 && rl.burstTokens == rl.burstLimit {
		//		fmt.Printf("not using burstmode\n")
		rl.tokens--
		//		fmt.Printf("token used, remaining: %d\n", rl.tokens)
		go func(rl *Ratelimiter) {
			time.Sleep(rl.resetAfter)
			if rl.tokens < rl.limit {
				rl.tokens++
				//				fmt.Printf("token added, remaining: %d\n", rl.tokens)
				rl.resetSignal <- 1
			}
		}(rl)
		return
	} else if rl.burstTokens > 0 {
		//		fmt.Printf("using burstmode\n")
		rl.burstTokens--
		//		fmt.Printf("bursttoken used, remaining: %d\n", rl.burstTokens)
		go func(rl *Ratelimiter) {
			time.Sleep(rl.burstResetAfter)
			if rl.burstTokens < rl.burstLimit {
				rl.burstTokens++
				//				fmt.Printf("bursttoken added, remaining: %d\n", rl.burstTokens)
				rl.resetSignal <- 1
			}
		}(rl)
		return
	} else {
		rl.queueSize++
		if rl.queueSize > rl.limit*2 {
			log.I("Warning: Ratelimiter", rl.Name, "'s queue is", rl.queueSize)
		}
		//		fmt.Printf("queued %s %d\n", rl.Name, rl.queueSize)
		<-rl.resetSignal
		rl.Request()
		rl.queueSize--
		//		fmt.Printf("dequeue %d\n", rl.queueSize)
		if rl.queueSize == 0 {
			rl.queueEnd <- 1
		}
		return
	}
}

func (rl *Ratelimiter) GetQueuesize() int {
	return rl.queueSize
}

func (rl *Ratelimiter) WaitForQueue() {
	if rl.queueSize <= 0 {
		return
	}
	<-rl.queueEnd
	return
}
