package layout

import (
	"fmt"
	"sync"
)

type Coord struct{ X, Y int }

func (c Coord) String() string {
	return fmt.Sprintf("[%d,%d]", c.X, c.Y)
}

var (
	currentLayout = make(map[Coord]MultiTile)
	layoutLock    sync.RWMutex
	version       uint64
)

var space = [...]Tile{Space1, Space2}

func Version() uint64 {
	layoutLock.RLock()
	defer layoutLock.RUnlock()

	return version
}

func Get(x, y int) MultiTile {
	return GetCoord(Coord{x, y})
}

func GetCoord(coord Coord) MultiTile {
	layoutLock.RLock()
	if t, ok := currentLayout[coord]; ok {
		layoutLock.RUnlock()
		return t
	}
	layoutLock.RUnlock()
	return baseLayout[coord]
}

func GetChanges() map[Coord]MultiTile {
	layoutLock.RLock()
	defer layoutLock.RUnlock()

	clone := make(map[Coord]MultiTile, len(currentLayout))
	for k, v := range currentLayout {
		clone[k] = v
	}
	return clone
}

func SetChanges(m map[Coord]MultiTile) {
	layoutLock.Lock()
	defer layoutLock.Unlock()

	currentLayout = m
	version++
	visInvalidateAll()
}

func SetCoord(coord Coord, check, t MultiTile) bool {
	layoutLock.Lock()
	defer layoutLock.Unlock()

	old := currentLayout[coord]
	if old == nil {
		old = baseLayout[coord]
	}

	if old.equal(check) {
		if old.BlocksVision() != t.BlocksVision() {
			if old.Door() || t.Door() {
				visInvalidateRecursive(coord)
			} else {
				visInvalidate(coord)
			}
		}
		currentLayout[coord] = t
		version++
		OnChange(coord, t)
		return true
	}
	return false
}

func GetSpace(x, y int) Tile {
	return space[uint(x^y)%uint(len(space))]
}

func AllTiles(f func(Coord, MultiTile)) {
	layoutLock.RLock()
	defer layoutLock.RUnlock()

	for coord := range baseLayout {
		if t, ok := currentLayout[coord]; ok {
			if t == nil {
				continue
			}

			f(coord, t)
		} else {
			f(coord, baseLayout[coord])
		}
	}

	for coord, tile := range currentLayout {
		if _, ok := baseLayout[coord]; ok || tile == nil {
			continue
		}

		f(coord, tile)
	}
}

var OnChange = func(Coord, MultiTile) {}
