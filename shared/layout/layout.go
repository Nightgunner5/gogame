package layout

import (
	"fmt"
	"sync"
)

type Coord struct{ X, Y int }

func (c Coord) String() string {
	return fmt.Sprintf("[%d,%d]", c.X, c.Y)
}

var baseLayout = map[Coord]MultiTile{
	// Arrival Ship
	Coord{-6, -3}: {Wall1SE},
	Coord{-6, -2}: {Wall1},
	Coord{-6, -1}: {Wall1},
	Coord{-6, 0}:  {Wall1},
	Coord{-6, 1}:  {Wall1},
	Coord{-6, 2}:  {Wall1},
	Coord{-6, 3}:  {Wall1NE},
	Coord{-5, -3}: {Wall1},
	Coord{-5, -2}: {GrayTile, Wall1NW},
	Coord{-5, -1}: {GrayTile},
	Coord{-5, 0}:  {GrayTile},
	Coord{-5, 1}:  {GrayTile},
	Coord{-5, 2}:  {GrayTile, Wall1SW},
	Coord{-5, 3}:  {Wall1},
	Coord{-4, -3}: {Wall1},
	Coord{-4, -2}: {GrayTile},
	Coord{-4, -1}: {WhiteTile},
	Coord{-4, 0}:  {WhiteTile},
	Coord{-4, 1}:  {WhiteTile},
	Coord{-4, 2}:  {GrayTile},
	Coord{-4, 3}:  {Wall1},
	Coord{-3, -3}: {Wall1},
	Coord{-3, -2}: {GrayTile},
	Coord{-3, -1}: {WhiteTile},
	Coord{-3, 0}:  {WhiteTile},
	Coord{-3, 1}:  {WhiteTile},
	Coord{-3, 2}:  {GrayTile},
	Coord{-3, 3}:  {GrayTile, Door1Closed},
	Coord{-2, -3}: {Wall1},
	Coord{-2, -2}: {GrayTile},
	Coord{-2, -1}: {WhiteTile},
	Coord{-2, 0}:  {WhiteTile},
	Coord{-2, 1}:  {WhiteTile},
	Coord{-2, 2}:  {GrayTile},
	Coord{-2, 3}:  {Wall1},
	Coord{-1, -4}: {Wall1SE},
	Coord{-1, -3}: {Wall1},
	Coord{-1, -2}: {GrayTile, Wall1},
	Coord{-1, -1}: {WhiteTile, Door1Closed},
	Coord{-1, 0}:  {WhiteTile, Wall1},
	Coord{-1, 1}:  {WhiteTile, Window1},
	Coord{-1, 2}:  {GrayTile, Wall1},
	Coord{-1, 3}:  {Wall1},
	Coord{-1, 4}:  {Wall1NE},
	Coord{0, -4}:  {Wall1},
	Coord{0, -3}:  {GrayTile, Wall1NW},
	Coord{0, -2}:  {GrayTile},
	Coord{0, -1}:  {WhiteTile},
	Coord{0, 0}:   {WhiteTile},
	Coord{0, 1}:   {WhiteTile},
	Coord{0, 2}:   {GrayTile},
	Coord{0, 3}:   {GrayTile, Wall1SW},
	Coord{0, 4}:   {Wall1},
	Coord{1, -4}:  {Wall1},
	Coord{1, -3}:  {GrayTile},
	Coord{1, -2}:  {WhiteTile},
	Coord{1, -1}:  {WhiteTile},
	Coord{1, 0}:   {WhiteTile},
	Coord{1, 1}:   {WhiteTile},
	Coord{1, 2}:   {WhiteTile},
	Coord{1, 3}:   {GrayTile},
	Coord{1, 4}:   {Wall1},
	Coord{2, -4}:  {Wall1},
	Coord{2, -3}:  {GrayTile},
	Coord{2, -2}:  {WhiteTile},
	Coord{2, -1}:  {WhiteTile},
	Coord{2, 0}:   {WhiteTile},
	Coord{2, 1}:   {WhiteTile},
	Coord{2, 2}:   {WhiteTile},
	Coord{2, 3}:   {GrayTile},
	Coord{2, 4}:   {Wall1},
	Coord{3, -4}:  {Wall1},
	Coord{3, -3}:  {GrayTile},
	Coord{3, -2}:  {WhiteTile},
	Coord{3, -1}:  {WhiteTile},
	Coord{3, 0}:   {WhiteTile},
	Coord{3, 1}:   {WhiteTile},
	Coord{3, 2}:   {WhiteTile},
	Coord{3, 3}:   {GrayTile},
	Coord{3, 4}:   {GrayTile, Door1Closed},
	Coord{4, -4}:  {Wall1},
	Coord{4, -3}:  {GrayTile},
	Coord{4, -2}:  {WhiteTile},
	Coord{4, -1}:  {WhiteTile},
	Coord{4, 0}:   {WhiteTile},
	Coord{4, 1}:   {WhiteTile},
	Coord{4, 2}:   {WhiteTile},
	Coord{4, 3}:   {GrayTile},
	Coord{4, 4}:   {Wall1},
	Coord{5, -4}:  {Wall1},
	Coord{5, -3}:  {GrayTile, Wall1NE},
	Coord{5, -2}:  {GrayTile},
	Coord{5, -1}:  {WhiteTile},
	Coord{5, 0}:   {WhiteTile},
	Coord{5, 1}:   {WhiteTile},
	Coord{5, 2}:   {GrayTile},
	Coord{5, 3}:   {GrayTile, Wall1SE},
	Coord{5, 4}:   {Wall1},
	Coord{6, -4}:  {Wall1SW},
	Coord{6, -3}:  {Wall1},
	Coord{6, -2}:  {GrayTile, Wall1NE},
	Coord{6, -1}:  {GrayTile},
	Coord{6, 0}:   {GrayTile},
	Coord{6, 1}:   {GrayTile},
	Coord{6, 2}:   {GrayTile, Wall1SE},
	Coord{6, 3}:   {Wall1},
	Coord{6, 4}:   {Wall1NW},
	Coord{7, -3}:  {Wall1SW},
	Coord{7, -2}:  {Wall1},
	Coord{7, -1}:  {Window1},
	Coord{7, 0}:   {Window1},
	Coord{7, 1}:   {Window1},
	Coord{7, 2}:   {Wall1},
	Coord{7, 3}:   {Wall1NW},

	// Arrival Wing
	Coord{-1, 7}:  {Wall1SE},
	Coord{-1, 8}:  {Wall1},
	Coord{-1, 9}:  {Wall1},
	Coord{-1, 10}: {Wall1},
	Coord{-1, 11}: {Wall1NE},
	Coord{0, 7}:   {Window1},
	Coord{0, 8}:   {OrangeTile},
	Coord{0, 9}:   {OrangeTile},
	Coord{0, 10}:  {OrangeTile, Wall1SW},
	Coord{0, 11}:  {Wall1},
	Coord{0, 12}:  {Wall1},
	Coord{0, 13}:  {Wall1},
	Coord{0, 14}:  {Window1},
	Coord{0, 15}:  {Window1},
	Coord{1, 7}:   {Window1},
	Coord{1, 8}:   {OrangeTile},
	Coord{1, 9}:   {WhiteTile},
	Coord{1, 10}:  {OrangeTile},
	Coord{1, 11}:  {OrangeTile},
	Coord{1, 12}:  {OrangeTile},
	Coord{1, 13}:  {OrangeTile},
	Coord{1, 14}:  {OrangeTile},
	Coord{1, 15}:  {Window1},
	Coord{2, 5}:   {Wall1},
	Coord{2, 6}:   {Wall1},
	Coord{2, 7}:   {Wall1},
	Coord{2, 8}:   {OrangeTile},
	Coord{2, 9}:   {WhiteTile},
	Coord{2, 10}:  {WhiteTile},
	Coord{2, 11}:  {WhiteTile},
	Coord{2, 12}:  {WhiteTile},
	Coord{2, 13}:  {WhiteTile},
	Coord{2, 14}:  {OrangeTile},
	Coord{2, 15}:  {Window1},
	Coord{3, 5}:   {BlackTile},
	Coord{3, 6}:   {BlackTile},
	Coord{3, 7}:   {BlackTile, Door1Closed},
	Coord{3, 8}:   {OrangeTile},
	Coord{3, 9}:   {WhiteTile},
	Coord{3, 10}:  {WhiteTile},
	Coord{3, 11}:  {WhiteTile},
	Coord{3, 12}:  {WhiteTile},
	Coord{3, 13}:  {WhiteTile},
	Coord{3, 14}:  {OrangeTile},
	Coord{3, 15}:  {Window1},
	Coord{4, 5}:   {Wall1},
	Coord{4, 6}:   {Wall1},
	Coord{4, 7}:   {Wall1},
	Coord{4, 8}:   {OrangeTile},
	Coord{4, 9}:   {WhiteTile},
	Coord{4, 10}:  {WhiteTile},
	Coord{4, 11}:  {WhiteTile},
	Coord{4, 12}:  {WhiteTile},
	Coord{4, 13}:  {WhiteTile},
	Coord{4, 14}:  {OrangeTile},
	Coord{4, 15}:  {Window1},
	Coord{5, 7}:   {Window1},
	Coord{5, 8}:   {OrangeTile},
	Coord{5, 9}:   {WhiteTile},
	Coord{5, 10}:  {WhiteTile},
	Coord{5, 11}:  {WhiteTile},
	Coord{5, 12}:  {WhiteTile},
	Coord{5, 13}:  {WhiteTile},
	Coord{5, 14}:  {OrangeTile},
	Coord{5, 15}:  {Window1},
	Coord{6, 7}:   {Window1},
	Coord{6, 8}:   {OrangeTile},
	Coord{6, 9}:   {WhiteTile},
	Coord{6, 10}:  {WhiteTile},
	Coord{6, 11}:  {OrangeTile},
	Coord{6, 12}:  {OrangeTile},
	Coord{6, 13}:  {OrangeTile},
	Coord{6, 14}:  {OrangeTile},
	Coord{6, 15}:  {Wall1},
	Coord{6, 16}:  {Wall1NE},
	Coord{7, 7}:   {Wall1},
	Coord{7, 8}:   {OrangeTile, Wall1NE},
	Coord{7, 9}:   {OrangeTile},
	Coord{7, 10}:  {OrangeTile},
	Coord{7, 11}:  {OrangeTile, Wall1SE},
	Coord{7, 12}:  {OrangeTile, Wall1},
	Coord{7, 13}:  {OrangeTile, Door1Closed},
	Coord{7, 14}:  {OrangeTile, Wall1},
	Coord{7, 15}:  {Wall1},
	Coord{7, 16}:  {Wall1},
	Coord{8, 7}:   {Wall1SW},
	Coord{8, 8}:   {Wall1},
	Coord{8, 9}:   {Wall1},
	Coord{8, 10}:  {Wall1},
	Coord{8, 11}:  {Wall1},
	Coord{8, 12}:  {WhiteTile, Wall1NW},
	Coord{8, 13}:  {WhiteTile},
	Coord{8, 14}:  {WhiteTile},
	Coord{8, 15}:  {WhiteTile},
	Coord{8, 16}:  {WhiteTile},
	Coord{9, 11}:  {Window1},
	Coord{9, 12}:  {WhiteTile},
	Coord{9, 13}:  {WhiteTile},
	Coord{9, 14}:  {WhiteTile},
	Coord{9, 15}:  {WhiteTile},
	Coord{9, 16}:  {WhiteTile},
	Coord{10, 11}: {Window1},
	Coord{10, 12}: {WhiteTile},
	Coord{10, 13}: {WhiteTile},
	Coord{10, 14}: {WhiteTile},
	Coord{10, 15}: {WhiteTile},
	Coord{10, 16}: {WhiteTile},
	Coord{11, 11}: {Wall1},
	Coord{11, 12}: {WhiteTile},
	Coord{11, 13}: {WhiteTile},
	Coord{11, 14}: {WhiteTile},
	Coord{11, 15}: {WhiteTile},
	Coord{11, 16}: {WhiteTile},
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
	defer layoutLock.RUnlock()

	if t, ok := currentLayout[coord]; ok {
		return t
	}
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
