package actor

import (
	"time"
)

type Ticker <-chan struct{}

func Tick(delay time.Duration) Ticker {
	c := make(chan struct{})

	go func() {
		for {
			time.Sleep(delay)
			c <- struct{}{}
		}
	}()

	return c
}
