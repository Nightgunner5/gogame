package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/network"
	"github.com/Nightgunner5/gogame/spell"
	"math"
)

type (
	Imp interface {
		entity.Entity
		entity.Positioner
		entity.Healther
		entity.Thinker
		spell.Caster

		Name() string

		imp()
	}

	imp struct {
		entity.EntityID
		entity.Positioner
		entity.Healther
		spell.SpellCaster

		name       string
		master     Magician
		nextSearch float64
	}
)

func NewImp(master Magician, x, y, z float64) Imp {
	const (
		maxHealth = 10
	)
	i := &imp{master: master, name: impName()}

	i.Positioner = entity.BasePosition(i, x, y, z)
	i.Healther = entity.BaseHealth(i, maxHealth)

	entity.Spawn(i)

	network.Broadcast(network.NewPacket(EntityName).
		Set(network.EntityID, i.ID()).
		Set(EntityName, i.name), false)

	return i
}

func (i *imp) Name() string {
	return i.name
}

func (i *imp) Parent() entity.Entity {
	return i.master
}

func (m *imp) Tag() string {
	return "imp"
}

func (i *imp) Think(delta float64) {
	const (
		maxCastDistance     = 10
		spellCastTime       = 1
		spellDamage         = 5
		moveSpeed           = 0.75
		personalSpaceBuffer = 5
	)

	if i.Health() <= 0 {
		entity.Despawn(i)
		return
	}

	i.nextSearch -= delta

	if i.CasterThink(delta) {
		// currently casting spell
		return
	}

	foundTarget := false
	if i.nextSearch <= 0 {
		i.nextSearch = spellCastTime
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
			foundTarget = true
		})
	}

	if !foundTarget {
		px, py, pz := i.master.Position()
		x, y, z := i.Position()
		x, y, z = px-x, py-y, pz-z
		if x != 0 || y != 0 || z != 0 {
			m := math.Sqrt(x*x + y*y + z*z)
			if m < personalSpaceBuffer {
				return
			}
			m /= moveSpeed * delta
			x, y, z = x/m, y/m, z/m
			i.Move(x, y, z)
		}
	}
}

func (imp) imp() {}
