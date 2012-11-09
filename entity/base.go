package entity

import "sync"

type (
	baseCounter struct {
		m   sync.RWMutex
		c   float64
		set bool
	}

	BaseHealth struct {
		Max float64
		baseCounter
	}

	BaseResource struct {
		Max float64
		baseCounter
	}
)
var _ Healther = new(BaseHealth)
var _ Resourcer = new(BaseResource)

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

	if b.c > amount || force {
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

func (b *BaseHealth) Health() float64 {
	return b.get(b.Max)
}
func (b *BaseResource) Resource() float64 {
	return b.get(b.Max)
}

func (b *BaseHealth) TakeDamage(amount float64, attacker Entity) {
	b.sub(amount, b.Max, true)
}
func (b *BaseResource) UseResource(amount float64) bool {
	return b.sub(amount, b.Max, false)
}
