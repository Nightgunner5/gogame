package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/go-gl/gl"
)

type renderMage struct {
	x, y, z      float64
	mana, health float64
}

type renderSpell struct {
	x, y, z float64
}

func render() {
	mages := make([]*renderMage, *mageCount)
	spells := make([]*renderSpell, *mageCount)
	entity.ForAll(func(e entity.Entity) {
		if p, ok := e.(Mage); ok {
			var r renderMage

			r.x, r.y, r.z = p.Position()
			r.mana, r.health = p.Mana(), p.Health()

			mages[int(p.ID()-1)] = &r

			if spell := p.CurrentSpell(); spell != nil {
				if caster, target := spell.Caster(), spell.Target(); caster != nil && target != nil {
					var s renderSpell

					x1, y1, z1 := target.(Mage).Position()
					x2, y2, z2 := caster.(Mage).Position()

					interpolation := spell.TimeLeft() / spell.TotalTime()
					lerp := func(a, b float64) float64 {
						return a + (b-a)*interpolation
					}
					s.x, s.y, s.z = lerp(x1, x2), lerp(y1, y2), lerp(z1, z2)

					spells[int(p.ID()-1)] = &s
				}
			}
		}
	})

	gl.Begin(gl.QUADS)
	for _, r := range mages {
		if r != nil {
			v := func(x, y, z float64) { gl.Vertex3d(r.x+x-0.5, r.y+y-0.175, r.z+z) }
			gl.Color3d(0, 0, 0)
			v(0, 0, 0)
			v(0, 0.35, 0)
			v(1, 0.35, 0)
			v(1, 0, 0)

			gl.Color3d(0, 1, 0)
			v(0.05, 0.2, 0.1)
			v(0.05, 0.3, 0.1)
			v(r.health/maxHealth*0.9+0.05, 0.3, 0.1)
			v(r.health/maxHealth*0.9+0.05, 0.2, 0.1)

			gl.Color3d(0, 0, 1)
			v(0.05, 0.05, 0.1)
			v(0.05, 0.15, 0.1)
			v(r.mana/maxMana*0.9+0.05, 0.15, 0.1)
			v(r.mana/maxMana*0.9+0.05, 0.05, 0.1)
		}
	}
	for _, r := range spells {
		if r != nil {
			gl.Color3d(1, 0, 0)
			gl.Vertex3d(r.x-0.05, r.y-0.05, r.z)
			gl.Vertex3d(r.x-0.05, r.y+0.05, r.z)
			gl.Vertex3d(r.x+0.05, r.y+0.05, r.z)
			gl.Vertex3d(r.x+0.05, r.y-0.05, r.z)
		}
	}
	gl.End()
	gl.Flush()
}
