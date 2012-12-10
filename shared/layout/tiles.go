package layout

type Tile uint16

// YES THESE ARE REAL COLORS NOW SHUT UP
const (
	WhiteTile     Tile = 0
	GrayTile      Tile = 1
	BlackTile     Tile = 2
	RedTile       Tile = 3
	OrangeTile    Tile = 4
	RellowTile    Tile = 5
	YellowTile    Tile = 6
	GrellowTile   Tile = 7
	GreenTile     Tile = 8
	TurquoiseTile Tile = 9
	GrueTile      Tile = 10
	CyanTile      Tile = 11
	BlueTile      Tile = 12
	IndigoTile    Tile = 13
	PurpleTile    Tile = 14
	PinkTile      Tile = 15

	Wall1 Tile = 16

	Space1 Tile = 1022
	Space2 Tile = 1023
)

func (t Tile) Space() bool {
	return t >= Space1
}

func (t Tile) Passable() bool {
	return (t >= WhiteTile && t <= PinkTile) ||
		(t >= Space1)
}
