package entity

type DamageListener interface {
	// Optionally modifies *amount. Returns true if the state of the DamageListener has changed.
	OnTakeDamage(amount *float64, attacker, victim Entity) (changed bool)
}

type DoDamageListener interface {
	// Optionally modifies *amount. Returns true if the state of the DoDamageListener has changed.
	OnDoDamage(amount *float64, attacker, victim Entity) (changed bool)
}

type AllListeners interface {
	DamageListener
	DoDamageListener
}

type ListenerAdder interface {
	AddAll(interface{})
	AddDamageListener(DamageListener)
	AddDoDamageListener(DoDamageListener)

	RemoveAll(interface{})
	RemoveDamageListener(DamageListener)
	RemoveDoDamageListener(DoDamageListener)
}

type Listeners struct {
	damage   []DamageListener
	doDamage []DoDamageListener
}
var _ AllListeners = new(Listeners)

func (l *Listeners) AddAll(listener interface{}) {
	if dmg, ok := listener.(DamageListener); ok {
		l.AddDamageListener(dmg)
	}
	if dmg, ok := listener.(DoDamageListener); ok {
		l.AddDoDamageListener(dmg)
	}
}

func (l *Listeners) AddDamageListener(listener DamageListener) {
	l.damage = append(l.damage, listener)
}

func (l *Listeners) AddDoDamageListener(listener DoDamageListener) {
	l.doDamage = append(l.doDamage, listener)
}

func (l *Listeners) RemoveAll(listener interface{}) {
	if dmg, ok := listener.(DamageListener); ok {
		l.RemoveDamageListener(dmg)
	}
	if dmg, ok := listener.(DoDamageListener); ok {
		l.RemoveDoDamageListener(dmg)
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

func (l *Listeners) RemoveDoDamageListener(listener DoDamageListener) {
	for i := range l.doDamage {
		if l.doDamage[i] == listener {
			l.doDamage = append(l.doDamage[:i], l.doDamage[i+1:]...)
			return
		}
	}
}

func (l *Listeners) OnTakeDamage(amount *float64, attacker, victim Entity) (changed bool) {
	for _, listener := range l.damage {
		changed = listener.OnTakeDamage(amount, attacker, victim) || changed
	}
	return
}

func (l *Listeners) OnDoDamage(amount *float64, attacker, victim Entity) (changed bool) {
	for _, listener := range l.doDamage {
		changed = listener.OnDoDamage(amount, attacker, victim) || changed
	}
	return
}
