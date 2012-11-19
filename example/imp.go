package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
)

type (
	Imp interface {
		entity.Entity
		entity.Positioner
		entity.Healther
		entity.Thinker
		spell.Caster
	}

	imp struct {
		entity.EntityID
		entity.Healther
		spell.SpellCaster

		x, y, z float64
		master  Magician
	}
)

func NewImp(master Magician, x, y, z float64) Imp {
	const (
		maxHealth = 10
	)
	i := &imp{
		x:      x,
		y:      y,
		z:      z,
		master: master,
	}

	i.Healther = entity.BaseHealth(i, maxHealth)

	entity.Spawn(i)

	return i
}

func (i *imp) Parent() entity.Entity {
	return i.master
}

func (i *imp) Position() (x, y, z float64) {
	return i.x, i.y, i.z
}

func (i *imp) Think(delta float64) {
	const (
		maxCastDistance = 100
		spellCastTime   = 1
		spellDamage     = 5
	)

	if i.Health() <= 0 {
		entity.Despawn(i)
		return
	}

	if i.CasterThink(delta) {
		// currently casting spell
		return
	}

	entity.ForOneNearby(i, maxCastDistance, func(e entity.Entity) bool {
		if o, ok := e.(Magician); ok {
			return o != i.master
		}
		if o, ok := e.(Imp); ok {
			return o.Parent() != i.master
		}
		return false
	}, func(e entity.Entity) {
		i.Cast(spell.DamageSpell(spellCastTime, spellDamage, i, e, false))
	})
}
