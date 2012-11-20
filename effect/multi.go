package effect

import "github.com/Nightgunner5/gogame/entity"

type multiEffect []Effect

func Multi(effects ...Effect) Effect {
	return multiEffect(effects)
}

func (m multiEffect) OnTakeDamage(amount *float64, attacker, victim entity.Entity) {
	for _, e := range m {
		if effect, ok := e.(entity.DamageListener); ok {
			effect.OnTakeDamage(amount, attacker, victim)
		}
	}
}

func (m multiEffect) OnDoDamage(amount *float64, attacker, victim entity.Entity) {
	for _, e := range m {
		if effect, ok := e.(entity.DoDamageListener); ok {
			effect.OnDoDamage(amount, attacker, victim)
		}
	}
}

func (m multiEffect) String() string {
	var s []byte
	addSpace := false
	for _, e := range m {
		if addSpace {
			s = append(s, ' ')
		}

		description := e.String()

		addSpace = description != ""

		s = append(s, description...)
	}
	return string(s)
}

func (m multiEffect) effect() {}
