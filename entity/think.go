package entity

import (
	"log"
	"time"
)

var TimeScale float64 = 1

func think() {
	const (
		seconds = 1 / float64(time.Second)
		delay   = time.Second / 10
	)

	then := time.Now()
	for now := range time.Tick(delay) {
		Δ := float64(now.Sub(then)) * seconds * TimeScale // I felt like being cool and using a Greek letter today.
		globalEntityList.All(func(e Entity) {
			if t, ok := e.(Thinker); ok {
				start := time.Now()
				t.Think(Δ)
				end := time.Now()
				if diff := end.Sub(start); diff > delay/2 {
					log.Printf("Entity %d (%T) thinking for %v", e.ID(), e, diff)
				}
			}
		})
		then = now
	}
}

func init() {
	go think()
}
