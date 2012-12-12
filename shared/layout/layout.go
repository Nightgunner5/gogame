package layout

type Coord struct{ X, Y int }

var baseLayout = map[Coord]MultiTile{
	Coord{-2, -2}: {Wall1},
	Coord{-2, -1}: {Wall1},
	Coord{-2, 0}:  {Wall1},
	Coord{-2, 1}:  {Wall1},
	Coord{-2, 2}:  {Wall1},
	Coord{-1, -2}: {Wall1},
	Coord{-1, -1}: {WhiteTile},
	Coord{-1, 0}:  {WhiteTile},
	Coord{-1, 1}:  {WhiteTile},
	Coord{-1, 2}:  {Wall1},
	Coord{0, -2}:  {Wall1},
	Coord{0, -1}:  {WhiteTile},
	Coord{0, 0}:   {WhiteTile},
	Coord{0, 1}:   {WhiteTile},
	Coord{0, 2}:   {Wall1},
	Coord{1, -2}:  {Wall1},
	Coord{1, -1}:  {WhiteTile},
	Coord{1, 0}:   {WhiteTile},
	Coord{1, 1}:   {WhiteTile},
	Coord{1, 2}:   {Wall1},
	Coord{2, -2}:  {Wall1},
	Coord{2, -1}:  {Wall1},
	Coord{2, 0}:   {Wall1},
	Coord{2, 1}:   {Wall1},
	Coord{2, 2}:   {Wall1},
}

var currentLayout map[Coord]MultiTile

func init() {
	currentLayout = make(map[Coord]MultiTile, len(baseLayout))
	for k, v := range baseLayout {
		currentLayout[k] = v
	}
}

var space = [...]Tile{Space1, Space2}

func Get(x, y int) MultiTile {
	return currentLayout[Coord{x, y}]
}

func GetSpace(x, y int) Tile {
	return space[uint(x^y)%uint(len(space))]
}
