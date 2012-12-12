package actor

import (
	"time"
)

const (
	MaxDroppedTicks = 10
)

type Ticker <-chan struct{}

func (Ticker) tick(delay time.Duration, t chan struct{}) {
	skipped := 0
	for {
		time.Sleep(delay)
		select {
		case t <- struct{}{}:
			skipped = 0
		default:
			skipped++
		}
		if skipped > MaxDroppedTicks {
			close(t)
			return
		}
	}
}

func Tick(delay time.Duration) Ticker {
	c := make(chan struct{}, 1)

	go Ticker(c).tick(delay, c)

	return c
}
