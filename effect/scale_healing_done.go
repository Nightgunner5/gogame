package effect

import "github.com/Nightgunner5/gogame/entity"

func ScaleHealingDone(amount float64) Primitive {
	return scaleHealingDone{scalePrimitive(amount)}
}

type scaleHealingDone struct {
	scalePrimitive
}

var _ Primitive = new(scaleHealingDone)
var _ entity.DoDamageListener = new(scaleHealingDone)

func (scaleHealingDone) isSameType(other Primitive) bool {
	_, ok := other.(scaleHealingDone)
	return ok
}

func (a scaleHealingDone) combine(b Primitive) Primitive {
	return scaleHealingDone{a.scalePrimitive.combine(b.(scaleHealingDone).scalePrimitive)}
}

func (scale scaleHealingDone) String() string {
	s := scale.scalePrimitive.String()
	if s == "" {
		return ""
	}
	if s[0] == '+' {
		return "Increases healing done by " + s[1:] + "."
	}
	return "Reduces healing done by " + s[1:] + "."
}

func (scale scaleHealingDone) OnDoDamage(amount *float64, attacker, victim entity.Entity) {
	*amount = -scale.scale(-*amount)
}
