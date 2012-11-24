package effect

import "github.com/Nightgunner5/gogame/entity"

func AbsorbHealing(amount float64) Primitive {
	return &absorbHealingFlat{flatPrimitive(amount)}
}

type absorbHealingFlat struct {
	flatPrimitive
}

var _ entity.DamageListener = new(absorbHealingFlat)

func (*absorbHealingFlat) isSameType(other Primitive) bool {
	_, ok := other.(*absorbHealingFlat)
	return ok
}

func (a *absorbHealingFlat) combine(b Primitive) Primitive {
	return &absorbHealingFlat{a.flatPrimitive.combine(b.(*absorbHealingFlat).flatPrimitive)}
}

func (absorb *absorbHealingFlat) String() string {
	s := absorb.flatPrimitive.String()
	if s == "" {
		return ""
	}
	return "Absorbs " + s + " healing."
}

func (absorb *absorbHealingFlat) OnTakeDamage(amount *float64, attacker, victim entity.Entity) bool {
	c := absorb.consume(-*amount)
	*amount += c

	return c != 0
}
