package effect

import "github.com/Nightgunner5/gogame/entity"

type Primitive interface {
	isSameType(Primitive) bool
	// Only to be called if isSameType returns true for the same argument.
	combine(Primitive) Primitive

	effectThink(delta float64, ent entity.Entity)
	String() string
}
