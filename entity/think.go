package entity

import (
	"log"
	"time"
)

var TimeScale float64 = 1

// The entity list is passed as a parameter to prevent benchmarks from being fouled up.
func think(global EntityList) {
	const (
		seconds = 1 / float64(time.Second)
		delay   = time.Second / 10
	)

	then := time.Now()
	for now := range time.Tick(delay) {
		Δtime := float64(now.Sub(then)) * seconds * TimeScale
		global.All(func(e Entity) {
			if t, ok := e.(Thinker); ok {
				start := time.Now()
				t.Think(Δtime)
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
	go think(globalEntityList)
}
