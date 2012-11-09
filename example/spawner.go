package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
	"math/rand"
)

type (
	Spawner interface {
		entity.Entity
		entity.Thinker
		entity.Resourcer
		spell.Caster
	}

	spawner struct {
		entity.EntityID
		entity.BaseResource
		spell.SpellCaster
	}
)

var _ Spawner = new(spawner)

func (s *spawner) Parent() entity.Entity {
	return entity.World
}

func (s *spawner) Think(Δtime float64) {
	if !s.CasterThink(Δtime) {
		s.UseResource(-Δtime)

		if s.UseResource(1) {
			s.Cast(&spell.BasicSpell{
				CastTime: 0.5,
				Caster_:  s.ID(),
				Target_:  s.ID(),
				Action: func(target, caster entity.Entity) {
					entity.Spawn(&mage{
						BaseHealth:   entity.BaseHealth{Max: maxHealth},
						BaseResource: entity.BaseResource{Max: maxMana},
						x:            rand.Float64()*20 - 10,
						y:            rand.Float64()*20 - 10,
					})
				},
			})
		}
	}
}
