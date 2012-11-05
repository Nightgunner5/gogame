package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"log"
	"math/rand"
	"runtime"
)

const (
	maxHealth    = 100
	manaForSpell = 10

	castTime    = 1
	spellDamage = 75

	healCastTime = 0.9
	spellHealing = 25
)

func variance() float64 {
	return rand.Float64()/2 + 0.5
}

type (
	spellCast struct {
		damage   float64
		timeLeft float64
		target   entity.EntityID
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

		x, y, z float64
		health  float64
		mana    float64
		spell   *spellCast
	}
)

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
	if p.health <= 0 {
		log.Printf("%d: Killed by %d", p.ID(), attacker.ID())
		p.health = 0
		entity.Despawn(p)
	}
}

func (p *person) Think(delta float64) {
	p.mana += delta

	//log.Printf("%d H:%f M:%f", p.ID(), p.health, p.mana)

	if p.spell != nil {
		p.spell.timeLeft -= delta
		if p.spell.timeLeft <= 0 {
			target := entity.Get(p.spell.target)
			if target == nil {
				log.Printf("%d: Spell failed: Target is dead", p.ID())
			} else {
				log.Printf("%d: Hit %d with a spell", p.ID(), target.ID())
				target.(entity.TakeDamager).TakeDamage(p.spell.damage, p)
			}
			p.spell = nil
		}
	} else if p.mana >= manaForSpell {
		var target entity.Entity
		// Yay nondeterminism
		entity.ForAllNearby(p, 1, func(e entity.Entity) {
			if _, ok := e.(entity.TakeDamager); ok {
				if _, ok := e.(entity.Healther); ok {
					target = e
				}
			}
		})

		if target == nil {
			log.Printf("%d: Spell failed: No target", p.ID())
			notarget <- true
			return
		}

		p.mana -= manaForSpell

		if target.(entity.Healther).Health() <= p.Health() {
			p.spell = &spellCast{
				damage:   spellDamage * variance(),
				target:   target.ID(),
				timeLeft: castTime,
			}
		} else {
			p.spell = &spellCast{
				damage:   -spellHealing * variance(),
				target:   p.ID(),
				timeLeft: healCastTime,
			}
		}
		log.Printf("%d: Started casting spell on %d\n", p.ID(), target.ID())
	} else {
		//log.Printf("%d: Spell failed: Not enough mana", p.ID())
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
