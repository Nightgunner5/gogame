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
		entity.TakeDamager
		entity.Thinker
		spell.Caster
		Mana() float64
	}

	mage struct {
		entity.EntityID

		spell.SpellCaster

		x, y, z float64
		health  float64
		mana    float64

		sync.RWMutex
	}
)

func (p *mage) Health() float64 {
	p.RLock()
	defer p.RUnlock()

	if p.health > maxHealth {
		return maxHealth
	}
	if p.health < 0 {
		return 0
	}
	return p.health
}

func (p *mage) Mana() float64 {
	p.RLock()
	defer p.RUnlock()

	if p.mana > maxMana {
		return maxMana
	}
	if p.mana < 0 {
		return 0
	}
	return p.mana
}

func (p *mage) Parent() entity.Entity {
	return entity.World
}

func (p *mage) Position() (x, y, z float64) {
	return p.x, p.y, p.z
}

func (p *mage) TakeDamage(amount float64, attacker entity.Entity) {
	p.Lock()
	defer p.Unlock()

	if p.health <= 0 {
		return
	}

	p.health -= amount
	if p.health > maxHealth {
		p.health = maxHealth
	}
	if amount > 0 {
		p.Interrupt()
	}
	if p.health <= 0 {
		// killed
		p.health = 0
		entity.Despawn(p)
	}
}

func (p *mage) Think(Δtime float64) {
	if p.CasterThink(Δtime) {
		// do nothing; spell is casting
	} else {
		p.Lock()
		p.mana += Δtime * manaPerSecond
		if p.mana > maxMana {
			p.mana = maxMana
		}
		p.Unlock()

		if p.Health() < maxHealth-spellHealing {
			if p.mana >= manaForHealingSpell {
				p.Lock()
				defer p.Unlock()

				p.mana -= manaForHealingSpell

				p.Cast(spell.HealingOverTimeSpell(healCastTime, spellHealing, p, p, true))
			}
		} else {
			if p.mana >= manaForDamageSpell {
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

				p.Lock()
				defer p.Unlock()

				p.mana -= manaForDamageSpell

				p.Cast(spell.DamageSpell(damageCastTime, spellDamage*variance(), p, target, true))
			}
		}
	}
}
