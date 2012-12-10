package layout

type Coord struct{ X, Y int }

var BaseLayout = map[Coord]Tile{
	Coord{0, 0}: WhiteTile,
}

var CurrentLayout map[Coord]Tile

func init() {
	CurrentLayout = make(map[Coord]Tile, len(BaseLayout))
	for k, v := range BaseLayout {
		CurrentLayout[k] = v
	}
}

var Space = [...]Tile{Space1, Space2}
