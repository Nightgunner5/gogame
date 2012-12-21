package packet

import (
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
)

var (
	MsgInteract = message.NewKind("Interact")
)

type Interact struct {
	layout.Coord
}

func (Interact) Kind() message.Kind {
	return MsgInteract
}
