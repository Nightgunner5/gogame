package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
	"math/rand"
)

type (
	Magician interface {
		entity.Entity
		entity.Positioner
		entity.Healther
		entity.Resourcer
		entity.Thinker
		spell.Caster
	}

	magician struct {
		entity.EntityID
		entity.Healther
		entity.Resourcer
		spell.SpellCaster

		x, y, z float64
	}
)

func NewMagician(x, y, z float64) Magician {
	const (
		maxHealth = 100
		maxMana   = 1000
	)

	m := &magician{
		x: x,
		y: y,
		z: z,
	}

	m.Healther = entity.BaseHealth(m, maxHealth)
	m.Resourcer = entity.BaseResource(m, maxMana)

	entity.Spawn(m)

	return m
}

func (m *magician) Parent() entity.Entity {
	return entity.World
}

func (m *magician) Position() (x, y, z float64) {
	return m.x, m.y, m.z
}

func (m *magician) Think(delta float64) {
	const (
		manaPerSecond  = 10
		summonCost     = 50
		summonCastTime = 0.5
	)

	if m.Health() <= 0 {
		entity.Despawn(m)
		return
	}

	if m.CasterThink(delta) {
		// currently casting spell
		return
	}

	m.UseResource(-delta * manaPerSecond)

	if m.UseResource(summonCost) {
		m.Cast(&spell.BasicSpell{
			CastTime: summonCastTime,
			Caster_:  m.ID(),
			Target_:  m.ID(),
			Action:   summonImp,

			Tag: "summonimp",
		})
	}
}

func summonImp(target, caster entity.Entity) {
	m := caster.(Magician)
	x, y, z := m.Position()

	x += rand.Float64()*2 - 1
	y += rand.Float64()*2 - 1
	z += rand.Float64()*2 - 1

	NewImp(m, x, y, z)
}
