package effect

import (
	"github.com/Nightgunner5/gogame/entity"
	"strconv"
)

type scalePrimitive float64

func (scalePrimitive) effectThink(float64, entity.Entity) bool { return false }

func (a scalePrimitive) combine(b scalePrimitive) scalePrimitive {
	return a * b
}

func (scale scalePrimitive) String() string {
	if scale == 1 {
		return ""
	}
	if scale < 1 {
		if scale < 0 {
			scale = 0
		}
		return "-" + strconv.FormatUint(uint64(100-scale*100), 10) + "%"
	}
	return "+" + strconv.FormatUint(uint64(scale*100-100), 10) + "%"
}

func (scale scalePrimitive) scale(amount float64) float64 {
	if amount < 0 {
		return amount
	}
	if scale < 0 {
		return 0
	}
	return float64(scale) * amount
}
