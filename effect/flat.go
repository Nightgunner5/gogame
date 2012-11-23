package effect

import (
	"github.com/Nightgunner5/gogame/entity"
	"strconv"
)

type flatPrimitive float64

func (*flatPrimitive) effectThink(float64, entity.Entity) {}

func (flat *flatPrimitive) String() string {
	if *flat <= 0 {
		return ""
	}
	return strconv.FormatUint(uint64(*flat), 10)
}

func (a flatPrimitive) combine(b flatPrimitive) flatPrimitive {
	return a + b
}

func (flat *flatPrimitive) consume(max float64) (available float64) {
	if *flat <= 0 || max <= 0 {
		return
	}
	if float64(*flat) < max {
		available = float64(*flat)
		*flat = 0
	} else {
		available = max
		*flat -= flatPrimitive(max)
	}
	return
}
