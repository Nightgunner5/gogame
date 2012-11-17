package main

import "math/rand"

func main() {
	const (
		initialMagicianCount = 10
		areaSize         = 10
	)
	for i := 0; i < initialMagicianCount; i++ {
		x := rand.Float64()*2*areaSize - areaSize
		y := rand.Float64()*2*areaSize - areaSize

		NewMagician(x, y, 0)
	}

	go renderer()

	select {}
}
