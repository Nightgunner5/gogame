package effect

import "github.com/Nightgunner5/gogame/entity"

func ScaleHealingTaken(amount float64) Primitive {
	return scaleHealingTaken{scalePrimitive(amount)}
}

type scaleHealingTaken struct {
	scalePrimitive
}

var _ Primitive = new(scaleHealingTaken)
var _ entity.DamageListener = new(scaleHealingTaken)

func (scaleHealingTaken) isSameType(other Primitive) bool {
	_, ok := other.(scaleHealingTaken)
	return ok
}

func (a scaleHealingTaken) combine(b Primitive) Primitive {
	return scaleHealingTaken{a.scalePrimitive.combine(b.(scaleHealingTaken).scalePrimitive)}
}

func (scale scaleHealingTaken) String() string {
	s := scale.scalePrimitive.String()
	if s == "" {
		return ""
	}
	if s[0] == '+' {
		return "Increases healing taken by " + s[1:] + "."
	}
	return "Reduces healing taken by " + s[1:] + "."
}

func (scale scaleHealingTaken) OnTakeDamage(amount *float64, attacker, victim entity.Entity) {
	*amount = -scale.scale(-*amount)
}
