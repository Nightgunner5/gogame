package main

import (
	"github.com/Nightgunner5/gogame/effect"
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/network"
	"github.com/Nightgunner5/gogame/spell"
	"math"
	"math/rand"
	"sync"
)

type (
	Magician interface {
		entity.Entity
		entity.Positioner
		entity.Healther
		entity.Resourcer
		entity.Thinker
		spell.SpellCaster
		effect.Effected

		Name() string
		SetMotion(x, y, z float64)

		magician()
	}

	magician struct {
		entity.EntityID
		entity.Positioner
		entity.Healther
		entity.Resourcer
		spell.SpellCaster
		effect.Effected

		name   string
		motion [3]float64

		mtx sync.Mutex
	}
)

func NewMagician(x, y, z float64, name string) Magician {
	const (
		maxHealth = 100
		maxMana   = 160
	)

	m := &magician{name: name}

	m.Positioner = entity.BasePosition(m, x, y, z)
	m.Healther = entity.BaseHealth(m, maxHealth)
	m.Resourcer = entity.BaseResource(m, maxMana)
	m.SpellCaster = spell.BaseSpellCaster(m)
	m.Effected = effect.BaseEffected(m)

	entity.Spawn(m)

	network.Broadcast(network.NewPacket(EntityName).
		Set(network.EntityID, m.ID()).
		Set(EntityName, m.name), false)

	return m
}

func (m *magician) Name() string {
	return m.name
}

func (m *magician) Parent() entity.Entity {
	return entity.World
}

func (m *magician) Tag() string {
	return "magician"
}

func (m *magician) SetMotion(x, y, z float64) {
	if x != 0 || y != 0 || z != 0 {
		magnitude := math.Sqrt(x*x + y*y + z*z)
		x, y, z = x/magnitude, y/magnitude, z/magnitude
	}
	m.mtx.Lock()
	m.motion[0], m.motion[1], m.motion[2] = x, y, z
	m.mtx.Unlock()
}

func (m *magician) Think(delta float64) {
	const (
		manaPerSecond = 10
	)

	if m.Health() <= 0 {
		entity.Despawn(m)
		return
	}

	m.mtx.Lock()
	m.Move(m.motion[0]*delta, m.motion[1]*delta, 0)
	if m.motion[0] != 0 || m.motion[1] != 0 {
		m.Interrupt()
	}
	m.mtx.Unlock()

	m.EffectThink(delta)

	if m.CasterThink(delta) {
		// currently casting spell
		return
	}

	m.UseResource(-delta * manaPerSecond)
}

func summonImp(target, caster entity.Entity) {
	m := caster.(Magician)
	x, y, z := m.Position()

	x += rand.Float64()*2 - 1
	y += rand.Float64()*2 - 1
	z += rand.Float64()*2 - 1

	imp := NewImp(m, x, y, z)
	imp.AddEffect(effect.NewEffect("Impending Sickness", 10.1).
		Add(effect.After(effect.NewEffect("Summoning Sickness", 0).
		Add(effect.DamageOverTime(2)), 10)))
}

func summonShield(target, caster entity.Entity) {
	m := caster.(Magician)

	m.AddEffect(effect.NewEffect("Shield", 5).
		Add(effect.AbsorbDamage(200)))
}

func (magician) magician() {}
