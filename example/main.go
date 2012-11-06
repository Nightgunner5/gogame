package main

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
	//"log"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"
	"sync"
)

const (
	maxHealth = 100
	maxMana   = 20

	initialHealth = 100
	initialMana   = 10

	manaForDamageSpell = 10
	damageCastTime     = 1
	spellDamage        = 75

	manaForHealingSpell = 3
	healCastTime        = 0.7
	spellHealing        = 25
)

func variance() float64 {
	return rand.Float64()/2 + 0.75
}

type (
	Person interface {
		entity.Entity
		entity.Positioner
		entity.Healther
		entity.TakeDamager
		entity.Thinker
		Mana() float64
		Action() string
	}

	person struct {
		entity.EntityID

		spell.SpellCaster

		x, y, z float64
		health  float64
		mana    float64

		lock sync.Mutex

		action string
	}
)

func damageSpell(target, caster entity.Entity) {
	damage := spellDamage * variance()
	//log.Printf("%d: Hit %d for %0.1f damage", caster.ID(), target.ID(), damage)
	target.(entity.TakeDamager).TakeDamage(damage, caster)
}

func healingSpell(target, caster entity.Entity) {
	healing := spellHealing * variance()
	//log.Printf("%d: Healed %d for %0.1f", caster.ID(), target.ID(), healing)
	target.(entity.TakeDamager).TakeDamage(-healing, caster)
}

func (p *person) Health() float64 {
	if p.health > maxHealth {
		return maxHealth
	}
	if p.health < 0 {
		return 0
	}
	return p.health
}

func (p *person) Mana() float64 {
	if p.mana > maxMana {
		return maxMana
	}
	if p.mana < 0 {
		return 0
	}
	return p.mana
}

func (p *person) Parent() entity.Entity {
	return entity.World
}

func (p *person) Position() (x, y, z float64) {
	return p.x, p.y, p.z
}

func (p *person) TakeDamage(amount float64, attacker entity.Entity) {
	if attacker != p {
		p.lock.Lock()
		defer p.lock.Unlock()
	}
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
		//log.Printf("%d: Killed by %d", p.ID(), attacker.ID())
		p.health = 0
		entity.Despawn(p)
	}
}

func (p *person) Action() string {
	return p.action
}

func (p *person) Think(delta float64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.mana += delta
	if p.mana > maxMana {
		p.mana = maxMana
	}

	if p.CasterThink(delta) {
		// do nothing; spell is casting
	} else if p.Health() < maxHealth-spellHealing {
		if p.mana >= manaForHealingSpell {
			//log.Printf("%d: H%d M%d", p.ID(), int(p.Health()), int(p.mana))

			p.mana -= manaForHealingSpell

			p.Cast(&spell.BasicSpell{
				CastTime:      healCastTime,
				Interruptable: true,
				Caster_:       p.ID(),
				Target_:       p.ID(),
				Action:        healingSpell,
			})
			//log.Printf("%d: Started healing self", p.ID())

			p.action = "Healing"
		} else {
			p.action = ""
		}
	} else {
		if p.mana >= manaForDamageSpell {
			var target entity.Entity
			entity.ForOneNearby(p, 1, func(e entity.Entity) bool {
				_, ok := e.(Person)
				return ok
			}, func(e entity.Entity) {
				target = e
			})

			if target == nil {
				//log.Printf("%d: Spell failed: No target", p.ID())
				notarget <- true
				return
			}

			p.mana -= manaForDamageSpell

			p.Cast(&spell.BasicSpell{
				CastTime:      damageCastTime,
				Interruptable: true,
				Target_:       target.ID(),
				Caster_:       p.ID(),
				Action:        damageSpell,
			})
			//log.Printf("%d: Started casting spell on %d", p.ID(), target.ID())
			p.action = fmt.Sprintf("Spell:%d", target.ID())
		} else {
			p.action = ""
		}
	}
}

var notarget = make(chan bool)

func main() {
	maxPersonID := runtime.GOMAXPROCS(0)
	for i := 0; i < maxPersonID; i++ {
		entity.Spawn(&person{
			health: initialHealth * variance(),
			mana:   initialMana * variance(),
		})
	}

	format := fmt.Sprintf("%%%dd: [%%-%ds] [%%-%ds] %%s\n", int(math.Ceil(math.Log10(float64(maxPersonID)))), maxHealth/2, maxMana)

	for {
		select {
		case <-notarget:
			return
		default:
			fmt.Println()
			for i := 1; i <= maxPersonID; i++ {
				if e := entity.Get(entity.EntityID(i)); e != nil {
					p := e.(Person)
					fmt.Printf(format, i, strings.Repeat("H", int(p.Health()/2)), strings.Repeat("M", int(p.Mana())), p.Action())
				} else {
					fmt.Printf(format, i, "", "", "Dead")
				}
			}
			time.Sleep(time.Second / 2)
		}
	}
}
