package entity

import (
	"log"
	"runtime"
	"time"
)

var TimeScale float64 = 1

const (
	thinkDelay = time.Second / 10
)

type thinkTask struct {
	t Thinker
	d float64
}

func init() {
	toThink := make(chan thinkTask)
	go thinkDispatcher(toThink, globalEntityList)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go thinker(toThink)
	}
}

func thinkDispatcher(c chan<- thinkTask, global EntityList) {
	then := time.Now()
	for now := range time.Tick(thinkDelay) {
		delta := float64(now.Sub(then)) / float64(time.Second) * TimeScale

		global.Each(func(e Entity) {
			if t, ok := e.(Thinker); ok {
				c <- thinkTask{
					t: t,
					d: delta,
				}
			}
		})

		toSpawn.commit()

		then = now
	}
}

func thinker(c <-chan thinkTask) {
	for t := range c {
		start := time.Now()
		t.t.Think(t.d)
		if diff := time.Since(start); diff > thinkDelay/2 {
			log.Printf("Entity %d (%T) thinking for %v", t.t.(Entity).ID(), t.t, diff)
		}
	}
}
