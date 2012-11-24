package effect

import (
	"github.com/Nightgunner5/gogame/entity"
)

type damageOverTime struct {
	addPrimitive
}

func DamageOverTime(perSecond float64) Primitive {
	return damageOverTime{addPrimitive(perSecond)}
}

func (dot damageOverTime) effectThink(delta float64, ent entity.Entity) {
	if h, ok := ent.(entity.Healther); ok {
		// This has to be done in a goroutine since the effect list gets locked when damage is dealt or taken.
		go h.TakeDamage(dot.get()*delta, entity.World)
	}
}

func (damageOverTime) isSameType(other Primitive) bool {
	_, ok := other.(damageOverTime)
	return ok
}

func (a damageOverTime) combine(b Primitive) Primitive {
	return damageOverTime{a.addPrimitive.combine(b.(damageOverTime).addPrimitive)}
}

func (dot damageOverTime) String() string {
	s := dot.addPrimitive.String()
	if s == "" {
		return ""
	}
	return "Deals " + s + " damage per second."
}
