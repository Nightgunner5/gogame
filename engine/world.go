package engine

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
)

var (
	MsgSpawn = message.NewKind("Spawn")
)

type Spawn struct {
	Actor *actor.Actor
}

func (Spawn) Kind() message.Kind {
	return MsgSpawn
}

type world struct {
	actor.Actor
}

func (w world) Initialize() (message.Reciever, message.Sender) {
	messages, broadcast := w.Actor.Initialize()
	messages_ := make(chan message.Message)

	go func() {
		for {
			msg := <-messages
			switch m := msg.(type) {
			case Spawn:
				broadcast <- m
			default:
				messages_ <- m
			}
		}
	}()

	return messages_, broadcast
}

var World world

func init() {
	actor.TopLevel(World.Initialize())
}
