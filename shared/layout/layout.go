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
	Coord{-1, -2}: {GrayTile},
	Coord{-1, -1}: {WhiteTile},
	Coord{-1, 0}:  {WhiteTile},
	Coord{-1, 1}:  {WhiteTile},
	Coord{-1, 2}:  {GrayTile},
	Coord{-1, 3}:  {Wall1},
	Coord{-1, 4}:  {Wall1NE},
	Coord{0, -4}:  {Wall1},
	Coord{0, -3}:  {Wall1},
	Coord{0, -2}:  {GrayTile},
	Coord{0, -1}:  {WhiteTile},
	Coord{0, 0}:   {WhiteTile},
	Coord{0, 0}:   {WhiteTile},
	Coord{0, 1}:   {WhiteTile},
	Coord{0, 2}:   {GrayTile},
	Coord{0, 3}:   {Wall1},
	Coord{0, 4}:   {Wall1},
	Coord{1, -4}:  {Wall1},
	Coord{1, -3}:  {GrayTile, Wall1NW},
	Coord{1, -2}:  {GrayTile},
	Coord{1, -1}:  {WhiteTile},
	Coord{1, 0}:   {WhiteTile},
	Coord{1, 1}:   {WhiteTile},
	Coord{1, 2}:   {GrayTile},
	Coord{1, 3}:   {GrayTile, Wall1SW},
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
	Coord{3, 4}:   {GrayTile, Door1Open},
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
}

var (
	currentLayout = make(map[Coord]MultiTile)
	layoutLock    sync.RWMutex
)

var space = [...]Tile{Space1, Space2}

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
}

func SetCoord(coord Coord, t MultiTile) {
	layoutLock.Lock()
	defer layoutLock.Unlock()

	currentLayout[coord] = t
}

func GetSpace(x, y int) Tile {
	return space[uint(x^y)%uint(len(space))]
}
