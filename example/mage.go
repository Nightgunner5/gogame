package main

import (
	"fmt"
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
		Action() string
	}

	mage struct {
		entity.EntityID

		spell.SpellCaster

		x, y, z float64
		health  float64
		mana    float64

		sync.Mutex

		action string
		target entity.Entity
	}
)

func (p *mage) Health() float64 {
	if p.health > maxHealth {
		return maxHealth
	}
	if p.health < 0 {
		return 0
	}
	return p.health
}

func (p *mage) Mana() float64 {
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
	if p.health <= 0 {
		return
	}
	p.Lock()
	defer p.Unlock()

	p.health -= amount
	if p.health > maxHealth {
		p.health = maxHealth
	}
	if amount > 0 {
		p.Interrupt()
	}
	if p.health <= 0 {
		//log.Printf("%d: Killed by %d", p.ID(), attacker.ID())
		p.health = 0
		entity.Despawn(p)
	}
}

func (p *mage) Action() string {
	return p.action
}

func (p *mage) Think(Δtime float64) {
	if p.CasterThink(Δtime) {
		// do nothing; spell is casting
	} else {
		p.Lock()
		defer p.Unlock()

		p.mana += Δtime * manaPerSecond
		if p.mana > maxMana {
			p.mana = maxMana
		}

		if p.health < maxHealth-spellHealing {
			if p.mana >= manaForHealingSpell {
				//log.Printf("%d: H%d M%d", p.ID(), int(p.Health()), int(p.mana))

				p.target = p

				p.mana -= manaForHealingSpell

				p.Cast(spell.HealingOverTimeSpell(healCastTime, spellHealing, p, p, true))
				//log.Printf("%d: Started healing self", p.ID())

				p.action = "Healing"
			} else {
				p.action = ""
				p.target = nil
			}
		} else {
			if p.mana >= manaForDamageSpell {
				entity.ForOneNearby(p, 100, func(e entity.Entity) bool {
					_, ok := e.(Mage)
					return ok
				}, func(e entity.Entity) {
					p.target = e
				})

				if p.target == nil {
					//log.Printf("%d: Spell failed: No target", p.ID())
					return
				}

				p.mana -= manaForDamageSpell

				p.Cast(spell.DamageSpell(damageCastTime, spellDamage*variance(), p, p.target, true))
				//log.Printf("%d: Started casting spell on %d", p.ID(), target.ID())
				p.action = fmt.Sprintf("Spell:%d", p.target.ID())
			} else {
				p.action = ""
				p.target = nil
			}
		}
	}
}
