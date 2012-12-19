package packet

import (
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
)

var (
	MsgLocation = message.NewKind("Location")
)

type Location struct {
	ID    uint64
	Coord layout.Coord
	Flags uint32
}

const (
	FlagNone       uint32 = 0
	FlagSpriteMask uint32 = 7
	FlagSuit       uint32 = 1
	FlagSuitHelm   uint32 = 2
	FlagMonkey     uint32 = 3
	FlagSecurity   uint32 = 4
	FlagEngineer   uint32 = 5
	FlagMedic      uint32 = 6
	_              uint32 = 7
)

func (Location) Kind() message.Kind {
	return MsgLocation
}
