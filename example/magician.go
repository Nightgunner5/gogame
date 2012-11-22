package main

import (
	"github.com/Nightgunner5/gogame/effect"
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

		effect.EffectAdder
		Cast(spell spell.Spell)

		magician()
	}

	magician struct {
		entity.EntityID
		entity.Positioner
		entity.Healther
		entity.Resourcer
		spell.SpellCaster

		effect.BasicEffectAdder
	}
)

func NewMagician(x, y, z float64) Magician {
	const (
		maxHealth = 100
		maxMana   = 160
	)

	m := new(magician)

	m.Positioner = entity.BasePosition(m, x, y, z)
	m.Healther = entity.BaseHealth(m, maxHealth)
	m.Resourcer = entity.BaseResource(m, maxMana)

	entity.Spawn(m)

	return m
}

func (m *magician) Parent() entity.Entity {
	return entity.World
}

func (m *magician) Tag() string {
	return "magician"
}

func (m *magician) Think(delta float64) {
	const (
		manaPerSecond = 10
	)

	if m.Health() <= 0 {
		entity.Despawn(m)
		return
	}

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

	NewImp(m, x, y, z)
}

func summonShield(target, caster entity.Entity) {
	m := caster.(Magician)

	m.AddEffect(effect.AbsorbDamage(20), 5)
}

func (magician) magician() {}
