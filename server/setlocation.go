package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
)

var (
	MsgSetLocation = message.NewKind("SetLocation")
)

type SetLocation struct {
	ID    uint64
	Actor *actor.Actor
	Flags uint32
	Coord layout.Coord
}

func (SetLocation) Kind() message.Kind {
	return MsgSetLocation
}
