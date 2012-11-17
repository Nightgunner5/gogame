package main

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"time"
)

func renderer() {
	const (
		tickLength = time.Second
	)

	for _ = range time.Tick(tickLength) {
		fmt.Println()
		entity.ForEach(func(e entity.Entity) {
			if h, ok := e.(entity.Healther); ok {
				fmt.Printf("%d %T %f\n", e.ID(), e, h.Health())
			}
		})
	}
}
