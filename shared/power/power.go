package power

import (
	"github.com/Nightgunner5/gogame/shared/layout"
	"sort"
	"sync"
)

type sortableCoord []layout.Coord

func (s sortableCoord) Less(i, j int) bool {
	return s[i].X < s[j].X || (s[i].X == s[j].X && s[i].Y < s[j].Y)
}

func (s sortableCoord) Len() int {
	return len(s)
}

func (s sortableCoord) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var (
	powerNetwork *network
	powerVersion uint64
	powerLock    sync.RWMutex
)

func recomputeAll() {
	powerNetwork = new(network)
	powerVersion = layout.Version()

	var roots sortableCoord

	layout.AllTiles(func(c layout.Coord, tile layout.MultiTile) {
		for _, t := range tile {
			if t == layout.Generator {
				roots = append(roots, c)
			}
		}
	})

	sort.Sort(roots)

	for _, root := range roots {
		powerNetwork.construct(root)
	}

	powerNetwork.compute()
}

func Powered(x, y int, tile layout.Tile) bool {
	powerLock.RLock()
	if powerVersion == layout.Version() && powerNetwork != nil {
		result := powerNetwork.get(x, y, tile)
		powerLock.RUnlock()
		return result
	}
	powerLock.RUnlock()

	powerLock.Lock()
	if powerVersion == layout.Version() && powerNetwork != nil {
		result := powerNetwork.get(x, y, tile)
		powerLock.Unlock()
		return result
	}

	recomputeAll()
	result := powerNetwork.get(x, y, tile)
	powerLock.Unlock()
	return result
}
