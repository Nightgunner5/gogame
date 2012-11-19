package entity

type DamageListener interface {
	OnTakeDamage(amount *float64, attacker, victim Entity)
}

type ListenerAdder interface {
	AddAll(interface{})
	AddDamageListener(DamageListener)

	RemoveAll(interface{})
	RemoveDamageListener(DamageListener)
}

type Listeners struct {
	damage []DamageListener
}

func (l *Listeners) AddAll(listener interface{}) {
	if dmg, ok := listener.(DamageListener); ok {
		l.AddDamageListener(dmg)
	}
}

func (l *Listeners) AddDamageListener(listener DamageListener) {
	l.damage = append(l.damage, listener)
}

func (l *Listeners) RemoveAll(listener interface{}) {
	if dmg, ok := listener.(DamageListener); ok {
		l.RemoveDamageListener(dmg)
	}
}

func (l *Listeners) RemoveDamageListener(listener DamageListener) {
	for i := range l.damage {
		if l.damage[i] == listener {
			l.damage = append(l.damage[:i], l.damage[i+1:]...)
			return
		}
	}
}

func (l *Listeners) OnTakeDamage(amount *float64, attacker, victim Entity) {
	for _, listener := range l.damage {
		listener.OnTakeDamage(amount, attacker, victim)
	}
}
