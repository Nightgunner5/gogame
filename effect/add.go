package effect

import (
	"github.com/Nightgunner5/gogame/entity"
	"strconv"
)

type addPrimitive float64

func (addPrimitive) effectThink(float64, entity.Entity) {}

func (add addPrimitive) String() string {
	if add <= 0 {
		return ""
	}
	return strconv.FormatUint(uint64(add), 10)
}

func (a addPrimitive) combine(b addPrimitive) addPrimitive {
	return a + b
}

func (add addPrimitive) get() float64 {
	return float64(add)
}
