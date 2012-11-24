package effect

import "github.com/Nightgunner5/gogame/entity"

func ScaleDamageTaken(amount float64) Primitive {
	return scaleDamageTaken{scalePrimitive(amount)}
}

type scaleDamageTaken struct {
	scalePrimitive
}

var _ entity.DamageListener = new(scaleDamageTaken)

func (scaleDamageTaken) isSameType(other Primitive) bool {
	_, ok := other.(scaleDamageTaken)
	return ok
}

func (a scaleDamageTaken) combine(b Primitive) Primitive {
	return scaleDamageTaken{a.scalePrimitive.combine(b.(scaleDamageTaken).scalePrimitive)}
}

func (scale scaleDamageTaken) String() string {
	s := scale.scalePrimitive.String()
	if s == "" {
		return ""
	}
	if s[0] == '+' {
		return "Increases damage taken by " + s[1:] + "."
	}
	return "Reduces damage taken by " + s[1:] + "."
}

func (scale scaleDamageTaken) OnTakeDamage(amount *float64, attacker, victim entity.Entity) bool {
	*amount = scale.scale(*amount)
	return false
}
