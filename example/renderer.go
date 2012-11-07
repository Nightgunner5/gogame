package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/go-gl/gl"
)

type renderMage struct {
	x, y, z      float64
	mana, health float64
}

func render() {
	mages := make([]*renderMage, *mageCount)
	entity.ForAll(func(e entity.Entity) {
		if p, ok := e.(Mage); ok {
			var r renderMage

			r.x, r.y, r.z = p.Position()
			r.mana, r.health = p.Mana(), p.Health()

			mages[int(p.ID()-1)] = &r
		}
	})

	gl.Begin(gl.QUADS)
	for _, r := range mages {
		if r != nil {
			v := func(x, y, z float64) { gl.Vertex3d(r.x+x, r.y+y, r.z+z) }
			gl.Color3d(0, 0, 0)
			v(0, 0, 0)
			v(0, 1, 0)
			v(1, 1, 0)
			v(1, 0, 0)

			gl.Color3d(0, 1, 0)
			v(0.05, 0.55, 0.1)
			v(0.05, 0.95, 0.1)
			v(r.health/maxHealth*0.9+0.05, 0.95, 0.1)
			v(r.health/maxHealth*0.9+0.05, 0.55, 0.1)

			gl.Color3d(0, 0, 1)
			v(0.05, 0.05, 0.1)
			v(0.05, 0.45, 0.1)
			v(r.mana/maxMana*0.9+0.05, 0.45, 0.1)
			v(r.mana/maxMana*0.9+0.05, 0.05, 0.1)
		}
	}
	gl.End()
	gl.Flush()
}
