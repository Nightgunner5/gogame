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

	Wall1   Tile = 16
	Wall1NE Tile = 17
	Wall1NW Tile = 18
	Wall1SE Tile = 19
	Wall1SW Tile = 20

	Window1 Tile = 21

	Door1Open   Tile = 22
	Door1Closed Tile = 23

	Light1W Tile = 24
	Light1N Tile = 25
	Light1E Tile = 26
	Light1S Tile = 27

	Space1 Tile = 1022
	Space2 Tile = 1023
)

func (t Tile) Space() bool {
	return t >= Space1
}

func (t Tile) Passable() bool {
	return (t >= WhiteTile && t <= PinkTile) ||
		(t.Door() && t&1 == 0) ||
		(t >= Light1W && t <= Light1S) ||
		t.Space()
}

func (t Tile) BlocksVision() bool {
	return !t.Passable() && t != Window1
}

func (t Tile) Door() bool {
	return (t >= Door1Open && t <= Door1Closed)
}

func (t Tile) LightLevel() byte {
	if t >= Light1W && t <= Light1S {
		return 80
	}
	return 0
}

type MultiTile []Tile

func (m MultiTile) Space() bool {
	for _, t := range m {
		if !t.Space() {
			return false
		}
	}
	return true
}

func (m MultiTile) Passable() bool {
	for _, t := range m {
		if !t.Passable() {
			return false
		}
	}
	return true
}

func (m MultiTile) Door() bool {
	for _, t := range m {
		if t.Door() {
			return true
		}
	}
	return false
}

func (m MultiTile) BlocksVision() bool {
	for _, t := range m {
		if t.BlocksVision() {
			return true
		}
	}
	return false
}

func (m MultiTile) LightLevel() byte {
	var light byte
	for _, t := range m {
		light += t.LightLevel()
	}
	return light
}

func (a MultiTile) equal(b MultiTile) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
