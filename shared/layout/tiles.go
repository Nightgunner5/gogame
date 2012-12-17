package layout

import (
	"strconv"
)

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
	GrueTile      Tile = 9
	TurquoiseTile Tile = 10
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

	Door2Open   Tile = 28
	Door2Closed Tile = 29

	Computer Tile = 30
	Safe     Tile = 31

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
	return !t.Passable() && t != Window1 && t != Computer && t != Safe
}

func (t Tile) Door() bool {
	return t == Door1Open || t == Door1Closed ||
		t == Door2Open || t == Door2Closed
}

func (t Tile) LightLevel() byte {
	if t >= Light1W && t <= Light1S {
		return 80
	}
	if t == Computer {
		return 45
	}
	return 0
}

func (t Tile) String() string {
	switch t {
	case WhiteTile:
		return "WhiteTile"
	case GrayTile:
		return "GrayTile"
	case BlackTile:
		return "BlackTile"
	case RedTile:
		return "RedTile"
	case OrangeTile:
		return "OrangeTile"
	case RellowTile:
		return "RellowTile"
	case YellowTile:
		return "YellowTile"
	case GrellowTile:
		return "GrellowTile"
	case GreenTile:
		return "GreenTile"
	case GrueTile:
		return "GrueTile"
	case TurquoiseTile:
		return "TurquoiseTile"
	case CyanTile:
		return "CyanTile"
	case BlueTile:
		return "BlueTile"
	case IndigoTile:
		return "IndigoTile"
	case PurpleTile:
		return "PurpleTile"
	case PinkTile:
		return "PinkTile"

	case Wall1:
		return "Wall1"
	case Wall1NE:
		return "Wall1NE"
	case Wall1NW:
		return "Wall1NW"
	case Wall1SE:
		return "Wall1SE"
	case Wall1SW:
		return "Wall1SW"

	case Window1:
		return "Window1"

	case Door1Open:
		return "Door1Open"
	case Door1Closed:
		return "Door1Closed"

	case Light1W:
		return "Light1W"
	case Light1N:
		return "Light1N"
	case Light1E:
		return "Light1E"
	case Light1S:
		return "Light1S"

	case Door2Open:
		return "Door2Open"
	case Door2Closed:
		return "Door2Closed"

	case Computer:
		return "Computer"
	case Safe:
		return "Safe"

	case Space1:
		return "Space1"
	case Space2:
		return "Space2"
	}
	return strconv.FormatUint(uint64(t), 10)
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
	if m.Space() {
		return 30
	}
	var light byte
	for _, t := range m {
		light += t.LightLevel()
	}
	return light
}

func (a MultiTile) String() string {
	var s []byte

	s = append(s, '{')

	for i, t := range a {
		if i != 0 {
			s = append(s, ", "...)
		}
		s = append(s, t.String()...)
	}

	s = append(s, '}')
	return string(s)
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
