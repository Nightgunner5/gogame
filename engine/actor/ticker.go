package actor

import (
	"time"
)

type Ticker <-chan struct{}

func (Ticker) tick(delay time.Duration, t chan struct{}) {
	for {
		time.Sleep(delay)
		t <- struct{}{}
	}
}

func Tick(delay time.Duration) Ticker {
	c := make(chan struct{}, 1)

	go Ticker(c).tick(delay, c)

	return c
}
