package entity

import (
	"fmt"
	"sync"
)

type (
	baseCounter struct {
		m   sync.RWMutex
		c   float64
		set bool
	}

	baseHealth struct {
		max float64
		ent EntityID
		baseCounter
	}

	baseResource struct {
		max float64
		ent EntityID
		baseCounter
	}
)

var _ Healther = new(baseHealth)
var _ Resourcer = new(baseResource)

func BaseHealth(e Entity, max float64) Healther {
	if max <= 0 {
		panic(fmt.Sprintf("max health (%v) must be positive", max))
	}
	return &baseHealth{
		max: max,
		ent: e.ID(),
	}
}

func BaseResource(e Entity, max float64) Resourcer {
	if max <= 0 {
		panic(fmt.Sprintf("max resource (%v) must be positive", max))
	}
	return &baseResource{
		max: max,
		ent: e.ID(),
	}
}

func (b *baseCounter) get(max float64) float64 {
	b.m.RLock()
	defer b.m.RUnlock()

	if b.set {
		if b.c > max {
			return max
		}
		if b.c < 0 {
			return 0
		}
		return b.c
	}
	return max
}

func (b *baseCounter) sub(amount, max float64, force bool) (changed bool) {
	b.m.Lock()
	defer b.m.Unlock()

	if !b.set {
		b.c = max
		b.set = true
	}

	if b.c <= 0 {
		return
	}

	if b.c >= amount || force {
		b.c -= amount
		changed = true
	}

	if b.c > max {
		b.c = max
	}

	if b.c < 0 {
		b.c = 0
	}
	return
}

func (b *baseHealth) Health() float64 {
	return b.get(b.max)
}
func (b *baseResource) Resource() float64 {
	return b.get(b.max)
}

func (b *baseHealth) TakeDamage(amount float64, attacker Entity) {
	ent := Get(b.ent)
	if l, ok := attacker.(DoDamageListener); ok {
		l.OnDoDamage(&amount, attacker, ent)
	}
	if l, ok := ent.(DamageListener); ok {
		l.OnTakeDamage(&amount, attacker, ent)
	}
	b.sub(amount, b.max, true)
}
func (b *baseResource) UseResource(amount float64) bool {
	return b.sub(amount, b.max, false)
}
