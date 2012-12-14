package actor

import (
	"time"
)

const (
	MaxDroppedTicks = 2
)

type Ticker chan struct{}

func (t Ticker) tick(delay time.Duration) {
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
	c := make(Ticker, 1)

	go c.tick(delay)

	return c
}
