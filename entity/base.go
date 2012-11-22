package entity

import (
	"fmt"
	"github.com/Nightgunner5/gogame/network"
	"sync"
)

type (
	baseCounter struct {
		mtx sync.Mutex
		cnt float64
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

	basePosition struct {
		xyz [3]float64
		ent EntityID
		mtx sync.Mutex
	}
)

var _ Healther = new(baseHealth)
var _ Resourcer = new(baseResource)
var _ Positioner = new(basePosition)

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

func BasePosition(e Entity, x, y, z float64) Positioner {
	return &basePosition{
		xyz: [3]float64{x, y, z},
		ent: e.ID(),
	}
}

func (b *baseCounter) get(max float64) float64 {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.set {
		if b.cnt > max {
			return max
		}
		if b.cnt < 0 {
			return 0
		}
		return b.cnt
	}
	return max
}

func (b *baseCounter) sub(amount, max float64, force bool) (changed bool) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if !b.set {
		b.cnt = max
		b.set = true
	}

	if b.cnt <= 0 {
		return
	}

	if b.cnt >= amount || force {
		b.cnt -= amount
		changed = true
	}

	if b.cnt > max {
		b.cnt = max
	}

	if b.cnt < 0 {
		b.cnt = 0
	}
	return
}

func (b *baseHealth) Health() float64 {
	return b.get(b.max)
}
func (b *baseResource) Resource() float64 {
	return b.get(b.max)
}
func (b *basePosition) Position() (x, y, z float64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return b.xyz[0], b.xyz[1], b.xyz[2]
}
func (b *basePosition) positionArray() []float64 {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return b.xyz[:]
}

func (b *baseHealth) TakeDamage(amount float64, attacker Entity) {
	ent := Get(b.ent)
	if l, ok := attacker.(DoDamageListener); ok {
		l.OnDoDamage(&amount, attacker, ent)
	}
	if l, ok := ent.(DamageListener); ok {
		l.OnTakeDamage(&amount, attacker, ent)
	}
	if amount != 0 {
		b.sub(amount, b.max, true)
		network.Broadcast(network.NewPacket(network.HealthChange).
			Set(network.AttackerID, attacker.ID()).
			Set(network.VictimID, b.ent).
			Set(network.Amount, b.get(b.max)), false)
	}
}
func (b *baseResource) UseResource(amount float64) bool {
	return b.sub(amount, b.max, false)
}
func (b *basePosition) Move(dx, dy, dz float64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if dx != 0 || dy != 0 || dz != 0 {
		b.xyz[0] += dx
		b.xyz[1] += dy
		b.xyz[2] += dz

		network.Broadcast(network.NewPacket(network.EntityPosition).
			Set(network.EntityID, b.ent).
			Set(network.EntityPosition, b.xyz[:]), false)
	}
}
