/*
Copyright (c) 2018 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package network offers various network tools
package network

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

//Tokenbucket implements a simple Tokenbucket ratelimiter
type Tokenbucket struct {
	Limiter *rate.Limiter
}

//NewTokenbucket creates a new Tokenbucket
func NewTokenbucket(refill time.Duration, tokens int) *Tokenbucket {
	lim := new(Tokenbucket)
	lim.Limiter = rate.NewLimiter(rate.Every(refill), tokens)
	return lim
}

//WaitUntil waits until the context is canceled or a token can be used
func (t *Tokenbucket) WaitUntil(ctx context.Context) error {
	return t.Limiter.Wait(ctx)
}

//Wait waits until a token can be used
func (t *Tokenbucket) Wait() error {
	ctx := context.Background()
	return t.Limiter.Wait(ctx)
}
