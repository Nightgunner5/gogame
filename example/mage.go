package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
	"math/rand"
	"sync"
)

func variance() float64 {
	return rand.Float64()/2 + 0.75
}

type (
	Mage interface {
		entity.Entity
		entity.Positioner
		entity.Healther
		entity.Resourcer
		entity.Thinker
		spell.Caster
	}

	mage struct {
		entity.EntityID

		spell.SpellCaster

		x, y, z float64
		entity.BaseHealth
		entity.BaseResource

		sync.RWMutex
	}
)

var _ Mage = new(mage)

func (p *mage) Parent() entity.Entity {
	return entity.World
}

func (p *mage) Position() (x, y, z float64) {
	return p.x, p.y, p.z
}

func (p *mage) Think(Δtime float64) {
	if p.Health() <= 0 {
		// we are dead
		entity.Despawn(p)
		return
	}
	if p.CasterThink(Δtime) {
		// do nothing; spell is casting
		return
	}

	p.UseResource(-Δtime * manaPerSecond)

	if p.Health() < maxHealth-spellHealing {
		if p.UseResource(manaForHealingSpell) {
			p.Cast(spell.HealingOverTimeSpell(healCastTime, spellHealing, p, p, true))
		}
	} else {
		var target entity.Entity
		entity.ForOneNearby(p, 100, func(e entity.Entity) bool {
			_, ok := e.(Mage)
			return ok
		}, func(e entity.Entity) {
			target = e
		})

		if target == nil {
			return
		}

		if p.UseResource(manaForDamageSpell) {
			p.Cast(spell.DamageSpell(damageCastTime, spellDamage*variance(), p, target, true))
		}
	}
}
