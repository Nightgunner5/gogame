package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
)

var (
	MsgMoveRequest = message.NewKind("MoveRequest")
)

type MoveRequest struct {
	Actor  *actor.Actor
	Dx, Dy int
}

func (MoveRequest) Kind() message.Kind {
	return MsgMoveRequest
}
