/*
Copyright (c) 2017 ceriath
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

type Tokenbucket struct {
	Limiter *rate.Limiter
}

func NewTokenbucket(refill time.Duration, tokens int) *Tokenbucket {
	lim := new(Tokenbucket)
	lim.Limiter = rate.NewLimiter(rate.Every(refill), tokens)
	return lim
}

func (t *Tokenbucket) WaitUntil(ctx context.Context) error{
	return t.Limiter.Wait(ctx)
}

func (t *Tokenbucket) Wait() error{
	ctx := context.Background()
	return t.Limiter.Wait(ctx)
}
