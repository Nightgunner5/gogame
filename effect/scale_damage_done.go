package effect

import "github.com/Nightgunner5/gogame/entity"

func ScaleDamageDone(amount float64) Primitive {
	return &scaleDamageDone{scalePrimitive(amount)}
}

type scaleDamageDone struct {
	scalePrimitive
}

var _ Primitive = new(scaleDamageDone)
var _ entity.DoDamageListener = new(scaleDamageDone)

func (scaleDamageDone) isSameType(other Primitive) bool {
	_, ok := other.(scaleDamageDone)
	return ok
}

func (a scaleDamageDone) combine(b Primitive) Primitive {
	return scaleDamageDone{a.scalePrimitive.combine(b.(scaleDamageDone).scalePrimitive)}
}

func (scale scaleDamageDone) String() string {
	s := scale.scalePrimitive.String()
	if s == "" {
		return ""
	}
	if s[0] == '+' {
		return "Increases damage dealt by " + s[1:] + "."
	}
	return "Reduces damage dealt by " + s[1:] + "."
}

func (scale scaleDamageDone) OnDoDamage(amount *float64, attacker, victim entity.Entity) {
	*amount = scale.scale(*amount)
}
