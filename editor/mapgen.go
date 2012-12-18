package main

import (
	"fmt"
	"github.com/Nightgunner5/gogame/shared/layout"
	"os"
	"sort"
)

type sortMap []struct {
	C layout.Coord
	T layout.MultiTile
}

func (l sortMap) Len() int {
	return len(l)
}

func (l sortMap) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l sortMap) Less(i, j int) bool {
	a, b := l[i], l[j]
	if a.C.X < b.C.X {
		return true
	}

	if a.C.X > b.C.X {
		return false
	}

	if a.C.Y < b.C.Y {
		return true
	}

	return false
}

func writeMapLayout() {
	var l sortMap

	layout.AllTiles(func(c layout.Coord, t layout.MultiTile) {
		if t.Space() {
			return
		}
		l = append(l, struct {
			C layout.Coord
			T layout.MultiTile
		}{c, t})
	})

	sort.Sort(l)

	f, err := os.Create("map.go")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(f, beforeLayout)
	for _, t := range l {
		fmt.Fprintf(f, "\n\tCoord{%d, %d}: %s,", t.C.X, t.C.Y, t.T)
	}
	fmt.Fprint(f, afterLayout)
	f.Close()
}

const (
	beforeLayout = "package layout\n\nvar baseLayout = map[Coord]MultiTile{"
	afterLayout  = "\n}\n"
)
