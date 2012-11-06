package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
	"log"
	"math/rand"
	"runtime"
)

const (
	maxHealth = 100
	maxMana   = 1000

	manaForDamageSpell = 10
	castTime           = 1
	spellDamage        = 75

	manaForHealingSpell = 3
	healCastTime        = 0.7
	spellHealing        = 25
)

func variance() float64 {
	return rand.Float64()/2 + 0.75
}

type (
	allPurposeSpell struct {
		damage   float64
		timeLeft float64
		target   entity.EntityID
		caster   entity.EntityID
	}

	Person interface {
		entity.Entity
		entity.Positioner
		entity.Healther
		entity.TakeDamager
		entity.Thinker
	}

	person struct {
		entity.EntityID

		spell.SpellCaster

		x, y, z float64
		health  float64
		mana    float64
	}
)

func (s *allPurposeSpell) Tick(delta float64) bool {
	if s.timeLeft <= 0 {
		return true
	}
	s.timeLeft -= delta
	if s.timeLeft <= 0 {
		target, caster := s.Target(), s.Caster()
		if caster == nil {
			log.Printf("%d: Spell failed: You are dead", s.caster)
			return true
		}
		if target == nil {
			log.Printf("%d: Spell failed: Target is dead", s.caster)
			return true
		}
		if s.damage == 0 {
			return true
		}
		if s.damage > 0 {
			log.Printf("%d: Hit %d for %0.1f damage", s.caster, s.target, s.damage)
		} else {
			log.Printf("%d: Healed %d for %0.1f", s.caster, s.target, -s.damage)
		}
		target.(entity.TakeDamager).TakeDamage(s.damage, caster)
		return true
	}
	return false
}

func (s *allPurposeSpell) Target() entity.Entity {
	return entity.Get(s.target)
}

func (s *allPurposeSpell) Caster() entity.Entity {
	return entity.Get(s.caster)
}

func (s *allPurposeSpell) Interrupt() bool {
	if s.timeLeft > 0 {
		s.timeLeft = -1
		return true
	}
	return false
}

func (s *allPurposeSpell) TimeLeft() float64 {
	return s.timeLeft
}

func (p *person) Health() float64 {
	return p.health
}

func (p *person) Parent() entity.Entity {
	return entity.World
}

func (p *person) Position() (x, y, z float64) {
	return p.x, p.y, p.z
}

func (p *person) TakeDamage(amount float64, attacker entity.Entity) {
	if p.health <= 0 {
		return
	}
	p.health -= amount
	if p.health > maxHealth {
		p.health = maxHealth
	}
	p.Interrupt()
	if p.health <= 0 {
		log.Printf("%d: Killed by %d", p.ID(), attacker.ID())
		p.health = 0
		entity.Despawn(p)
	}
}

func (p *person) Think(delta float64) {
	p.mana += delta
	if p.mana > maxMana {
		p.mana = maxMana
	}

	if p.SpellCaster.Tick(delta) {
		// do nothing; spell is casting
	} else if p.Health() < maxHealth-spellHealing {
		if p.mana >= manaForHealingSpell {
			log.Printf("%d: H%d M%d", p.ID(), int(p.Health()), int(p.mana))

			p.mana -= manaForHealingSpell

			p.Cast(&allPurposeSpell{
				damage:   -spellHealing * variance(),
				target:   p.ID(),
				caster:   p.ID(),
				timeLeft: healCastTime,
			})
			log.Printf("%d: Started healing self", p.ID())
		}
	} else {
		if p.mana >= manaForDamageSpell {
			var target entity.Entity
			// Yay nondeterminism
			entity.ForAllNearby(p, 1, func(e entity.Entity) {
				if _, ok := e.(Person); ok {
					target = e
				}
			})

			if target == nil {
				log.Printf("%d: Spell failed: No target", p.ID())
				notarget <- true
				return
			}

			p.mana -= manaForDamageSpell

			p.Cast(&allPurposeSpell{
				damage:   spellDamage * variance(),
				target:   target.ID(),
				caster:   p.ID(),
				timeLeft: castTime,
			})
			log.Printf("%d: Started casting spell on %d", p.ID(), target.ID())
		}
	}
}

var notarget = make(chan bool)

func main() {
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		entity.Spawn(&person{
			health: maxHealth,
			mana:   10 * variance(),
		})
	}

	<-notarget
}
