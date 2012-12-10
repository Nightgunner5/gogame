package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
)

var (
	MsgDespawn = message.NewKind("Despawn")
)

type Despawn struct {
	ID    uint64
	Actor *actor.Actor
}

func (Despawn) Kind() message.Kind {
	return MsgDespawn
}
