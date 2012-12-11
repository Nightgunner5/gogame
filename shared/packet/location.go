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
}

func (Location) Kind() message.Kind {
	return MsgLocation
}
