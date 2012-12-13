package packet

import (
	"github.com/Nightgunner5/gogame/shared/layout"
)

type MapChange struct {
	Coord layout.Coord
	Tile  layout.MultiTile
}

type MapOverride struct {
	NewMap map[layout.Coord]layout.MultiTile
}
