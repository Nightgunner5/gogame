package power

import (
	"github.com/Nightgunner5/gogame/shared/layout"
)

type network struct {
	roots   []*graph
	powered map[node]bool
}

type node struct {
	layout.Coord
	layout.Tile
}

type graph struct {
	node
	neighbors []*graph
}

func (n *network) construct(root layout.Coord) {
	g := &graph{
		node: node{root, layout.Generator},
	}
	n.roots = append(n.roots, g)

	visit(g, map[node]bool{g.node: true})
}

func visit(root *graph, visited map[node]bool) {
	switch root.Tile {
	case layout.Generator, layout.GeneratorOff:
		next(root, visited, -1, 0, layout.WireE)
		next(root, visited, 1, 0, layout.WireW)
		next(root, visited, 0, -1, layout.WireS)
		next(root, visited, 0, 1, layout.WireN)

	case layout.WireN:
		next(root, visited, 0, 0, layout.Light1NOn)
		next(root, visited, 0, -1, layout.WireS)
		next(root, visited, 0, -1, layout.Generator)
		next(root, visited, 0, -1, layout.GeneratorOff)

	case layout.WireS:
		next(root, visited, 0, 0, layout.Light1SOn)
		next(root, visited, 0, 1, layout.WireN)
		next(root, visited, 0, 1, layout.Generator)
		next(root, visited, 0, 1, layout.GeneratorOff)

	case layout.WireW:
		next(root, visited, 0, 0, layout.Light1WOn)
		next(root, visited, -1, 0, layout.WireE)
		next(root, visited, -1, 0, layout.Generator)
		next(root, visited, -1, 0, layout.GeneratorOff)

	case layout.WireE:
		next(root, visited, 0, 0, layout.Light1EOn)
		next(root, visited, 1, 0, layout.WireW)
		next(root, visited, 1, 0, layout.Generator)
		next(root, visited, 1, 0, layout.GeneratorOff)
	}

	if root.Tile == layout.WireW || root.Tile == layout.WireE ||
		root.Tile == layout.WireN || root.Tile == layout.WireS {
		next(root, visited, 0, 0, layout.DoorGeneralClosed)
		next(root, visited, 0, 0, layout.DoorGeneralOpen)
		next(root, visited, 0, 0, layout.DoorSecurityClosed)
		next(root, visited, 0, 0, layout.DoorSecurityOpen)
		next(root, visited, 0, 0, layout.DoorEngineerClosed)
		next(root, visited, 0, 0, layout.DoorEngineerOpen)
		next(root, visited, 0, 0, layout.DoorMedicalClosed)
		next(root, visited, 0, 0, layout.DoorMedicalOpen)

		next(root, visited, 0, 0, layout.WireW)
		next(root, visited, 0, 0, layout.WireE)
		next(root, visited, 0, 0, layout.WireN)
		next(root, visited, 0, 0, layout.WireS)
	}
}

func next(root *graph, visited map[node]bool, dx, dy int, tile layout.Tile) {
	n := node{layout.Coord{root.X + dx, root.Y + dy}, tile}
	if visited[n] {
		return
	}
	visited[n] = true

	tiles := layout.GetCoord(n.Coord)
	for _, t := range tiles {
		if t == tile {
			g := &graph{node: n}
			root.neighbors = append(root.neighbors, g)
			visit(g, visited)
			return
		}
	}
}

func (n *network) compute() {
	n.powered = make(map[node]bool)

	// TODO: limited power from generators

	for _, root := range n.roots {
		n.walk(root)
	}

	// Time to feed the GC
	n.roots = nil
}

func (n *network) walk(g *graph) {
	if !n.powered[g.node] {
		n.powered[g.node] = true
		for _, g2 := range g.neighbors {
			n.walk(g2)
		}
	}
}

func (n *network) get(x, y int, tile layout.Tile) bool {
	return n.powered[node{layout.Coord{x, y}, tile}]
}
