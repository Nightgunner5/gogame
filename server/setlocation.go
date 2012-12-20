package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/packet"
)

var (
	MsgSetLocation = message.NewKind("SetLocation")
)

type SetLocation struct {
	Actor *actor.Actor
	*packet.Packet
}

func (SetLocation) Kind() message.Kind {
	return MsgSetLocation
}
