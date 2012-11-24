package effect

import "github.com/Nightgunner5/gogame/entity"

func AbsorbDamage(amount float64) Primitive {
	return &absorbDamageFlat{flatPrimitive(amount)}
}

type absorbDamageFlat struct {
	flatPrimitive
}
var _ entity.DamageListener = new(absorbDamageFlat)

func (*absorbDamageFlat) isSameType(other Primitive) bool {
	_, ok := other.(*absorbDamageFlat)
	return ok
}

func (a *absorbDamageFlat) combine(b Primitive) Primitive {
	return &absorbDamageFlat{a.flatPrimitive.combine(b.(*absorbDamageFlat).flatPrimitive)}
}

func (absorb *absorbDamageFlat) String() string {
	s := absorb.flatPrimitive.String()
	if s == "" {
		return ""
	}
	return "Absorbs " + s + " damage."
}

func (absorb *absorbDamageFlat) OnTakeDamage(amount *float64, attacker, victim entity.Entity) bool {
	c := absorb.consume(*amount)
	*amount -= c
	return c != 0
}
