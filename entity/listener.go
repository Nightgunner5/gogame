package entity

type DamageListener interface {
	OnTakeDamage(amount *float64, attacker, victim Entity)
}

type ListenerAdder interface {
	AddDamageListener(listener DamageListener)
}

type Listeners struct {
	damage []DamageListener
}

func (l *Listeners) AddDamageListener(listener DamageListener) {
	l.damage = append(l.damage, listener)
}

func (l *Listeners) OnTakeDamage(amount *float64, attacker, victim Entity) {
	for _, listener := range l.damage {
		listener.OnTakeDamage(amount, attacker, victim)
	}
}
