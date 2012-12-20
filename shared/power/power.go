package power

import (
	"github.com/Nightgunner5/gogame/shared/layout"
	"sync"
)

type graph struct {
	location   layout.Coord
	neighbors  []*graph
	energydiff int16
	powered    bool
}

func (g *graph) isPowered() bool {
	return g != nil && g.powered
}

var (
	powered      map[layout.Coord]*graph
	powerVersion uint64
	powerRoots   []*graph
	powerLock    sync.RWMutex
)

func getTile(x, y int) *graph {
	if g, ok := powered[layout.Coord{x, y}]; ok {
		return g
	}
	g := new(graph)
	g.location = layout.Coord{x, y}
	powered[g.location] = g
	return g
}

func recomputeAll() {
	powered = make(map[layout.Coord]*graph)
	powerVersion = layout.Version()
	powerRoots = nil

	layout.AllTiles(func(c layout.Coord, tile layout.MultiTile) {
		var g *graph
		for _, t := range tile {
			switch t {
			case layout.Generator:
				g = getTile(c.X, c.Y)
				g.energydiff += 10000
				powerRoots = append(powerRoots, g)

			case layout.Computer:
				g = getTile(c.X, c.Y)
				g.energydiff -= 200

			case layout.DoorGeneralClosed, layout.DoorGeneralOpen,
				layout.DoorSecurityClosed, layout.DoorSecurityOpen,
				layout.DoorEngineerClosed, layout.DoorEngineerOpen,
				layout.DoorMedicalClosed, layout.DoorMedicalOpen:
				g = getTile(c.X, c.Y)
				g.energydiff -= 500

			case layout.Light1WOn, layout.Light1NOn, layout.Light1EOn, layout.Light1SOn:
				g = getTile(c.X, c.Y)
				g.energydiff -= 20

			case layout.WireW:
				g = getTile(c.X, c.Y)
				g.neighbors = append(g.neighbors, getTile(c.X-1, c.Y))
				g.energydiff--

			case layout.WireN:
				g = getTile(c.X, c.Y)
				g.neighbors = append(g.neighbors, getTile(c.X, c.Y-1))
				g.energydiff--

			case layout.WireE:
				g = getTile(c.X, c.Y)
				g.neighbors = append(g.neighbors, getTile(c.X+1, c.Y))
				g.energydiff--

			case layout.WireS:
				g = getTile(c.X, c.Y)
				g.neighbors = append(g.neighbors, getTile(c.X, c.Y+1))
				g.energydiff--

			}
		}
	})

	for _, g := range powerRoots {
		visit(g, make(map[*graph]bool), 0)
	}
}

func visit(root *graph, visited map[*graph]bool, current int16) int16 {
	if visited[root] {
		return current
	}
	visited[root] = true

	if !root.powered {
		if current+root.energydiff < 0 {
			return current
		}
		current += root.energydiff
		root.powered = true
	}
	if current == 0 {
		return 0
	}

	for _, g := range root.neighbors {
		current = visit(g, visited, current)
	}
	return current
}

func Powered(x, y int) bool {
	powerLock.RLock()
	if powerVersion != layout.Version() || powered == nil {
		powerLock.RUnlock()
		powerLock.Lock()
		if powerVersion != layout.Version() || powered == nil {
			recomputeAll()
		}
		isPowered := powered[layout.Coord{x, y}].isPowered()
		powerLock.Unlock()
		return isPowered
	}
	isPowered := powered[layout.Coord{x, y}].isPowered()
	powerLock.RUnlock()
	return isPowered
}
