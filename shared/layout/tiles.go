package layout

import (
	"strconv"
)

type Tile uint16

// YES THESE ARE REAL COLORS NOW SHUT UP
const (
	TileWhite     Tile = 0
	TileGray      Tile = 1
	TileBlack     Tile = 2
	TileRed       Tile = 3
	TileOrange    Tile = 4
	TileRellow    Tile = 5
	TileYellow    Tile = 6
	TileGrellow   Tile = 7
	TileGreen     Tile = 8
	TileGrue      Tile = 9
	TileTurquoise Tile = 10
	TileCyan      Tile = 11
	TileBlue      Tile = 12
	TileIndigo    Tile = 13
	TilePurple    Tile = 14
	TilePink      Tile = 15

	Wall1   Tile = 16
	Wall1NE Tile = 17
	Wall1NW Tile = 18
	Wall1SE Tile = 19
	Wall1SW Tile = 20

	Window Tile = 21

	DoorGeneralOpen    Tile = 22
	DoorGeneralClosed  Tile = 23
	DoorSecurityOpen   Tile = 24
	DoorSecurityClosed Tile = 25
	DoorEngineerOpen   Tile = 26
	DoorEngineerClosed Tile = 27
	DoorMedicalOpen    Tile = 28
	DoorMedicalClosed  Tile = 29

	Computer Tile = 30
	Safe     Tile = 31

	Light1WOff Tile = 32
	Light1WOn  Tile = 33
	Light1NOff Tile = 34
	Light1NOn  Tile = 35
	Light1EOff Tile = 36
	Light1EOn  Tile = 37
	Light1SOff Tile = 38
	Light1SOn  Tile = 39

	Space1 Tile = 1022
	Space2 Tile = 1023
)

func (t Tile) Space() bool {
	return t >= Space1
}

func (t Tile) Passable() bool {
	return (t >= TileWhite && t <= TilePink) ||
		(t.Door() && t&1 == 0) ||
		(t >= Light1WOff && t <= Light1SOn) ||
		t.Space()
}

func (t Tile) BlocksVision() bool {
	return !t.Passable() && t != Window && t != Computer && t != Safe
}

func (t Tile) Door() bool {
	return t >= DoorGeneralOpen && t <= DoorMedicalClosed
}

func (t Tile) LightLevel() byte {
	if t >= Light1WOn && t <= Light1SOn && t&1 == 1 {
		return 80
	}
	if t == Computer {
		return 35
	}
	return 0
}

func (t Tile) String() string {
	switch t {
	case TileWhite:
		return "TileWhite"
	case TileGray:
		return "TileGray"
	case TileBlack:
		return "TileBlack"
	case TileRed:
		return "TileRed"
	case TileOrange:
		return "TileOrange"
	case TileRellow:
		return "TileRellow"
	case TileYellow:
		return "TileYellow"
	case TileGrellow:
		return "TileGrellow"
	case TileGreen:
		return "TileGreen"
	case TileGrue:
		return "TileGrue"
	case TileTurquoise:
		return "TileTurquoise"
	case TileCyan:
		return "TileCyan"
	case TileBlue:
		return "TileBlue"
	case TileIndigo:
		return "TileIndigo"
	case TilePurple:
		return "TilePurple"
	case TilePink:
		return "TilePink"

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

	case Window:
		return "Window"

	case DoorGeneralOpen:
		return "DoorGeneralOpen"
	case DoorGeneralClosed:
		return "DoorGeneralClosed"
	case DoorSecurityOpen:
		return "DoorSecurityOpen"
	case DoorSecurityClosed:
		return "DoorSecurityClosed"
	case DoorEngineerOpen:
		return "DoorEngineerOpen"
	case DoorEngineerClosed:
		return "DoorEngineerClosed"
	case DoorMedicalOpen:
		return "DoorMedicalOpen"
	case DoorMedicalClosed:
		return "DoorMedicalClosed"

	case Computer:
		return "Computer"
	case Safe:
		return "Safe"

	case Light1WOff:
		return "Light1WOff"
	case Light1WOn:
		return "Light1WOn"
	case Light1NOff:
		return "Light1NOff"
	case Light1NOn:
		return "Light1NOn"
	case Light1EOff:
		return "Light1EOff"
	case Light1EOn:
		return "Light1EOn"
	case Light1SOff:
		return "Light1SOff"
	case Light1SOn:
		return "Light1SOn"

	case Space1:
		return "Space1"
	case Space2:
		return "Space2"
	}
	return strconv.FormatUint(uint64(t), 10)
}

func (t Tile) describe() (string, bool) {
	switch t {
	case TileWhite, TileGray, TileBlack, TileRed, TileOrange, TileRellow, TileYellow, TileGrellow, TileGreen, TileGrue, TileTurquoise, TileCyan, TileBlue, TileIndigo, TilePurple, TilePink:
		return "floor", false

	case Wall1, Wall1NE, Wall1NW, Wall1SE, Wall1SW:
		return "wall", true

	case Window:
		return "window", true

	case DoorGeneralOpen:
		return "door", false
	case DoorGeneralClosed:
		return "door", true
	case DoorSecurityOpen:
		return "security door", false
	case DoorSecurityClosed:
		return "security door", true
	case DoorEngineerOpen:
		return "engineering door", false
	case DoorEngineerClosed:
		return "engineering door", true
	case DoorMedicalOpen:
		return "medbay door", false
	case DoorMedicalClosed:
		return "medbay door", true

	case Computer:
		return "computer", false
	case Safe:
		return "safe", false

	case Light1WOff, Light1NOff, Light1EOff, Light1SOff:
		return "light socket", false
	case Light1WOn, Light1NOn, Light1EOn, Light1SOn:
		return "flourescent light", false
	}
	return "ERROR", false
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

func (m MultiTile) String() string {
	var s []byte

	s = append(s, '{')

	for i, t := range m {
		if i != 0 {
			s = append(s, ", "...)
		}
		s = append(s, t.String()...)
	}

	s = append(s, '}')
	return string(s)
}

func (m MultiTile) Describe() []string {
	if m.Space() {
		return []string{"space"}
	}

	var description []string
	for _, t := range m {
		d, erase := t.describe()
		if erase {
			description = nil
		}
		description = append(description, d)
	}

	for i, j := 0, len(description)-1; i < j; i, j = i+1, j-1 {
		description[i], description[j] = description[j], description[i]
	}

	return description
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
