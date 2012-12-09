package engine

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
)

var (
	// Spawn Message - see Spawn.
	MsgSpawn = message.NewKind("Spawn")
)

// A Message which, when given to World, is forwarded to all Subscribers to
// the MsgSpawn Message.
type Spawn struct {
	Actor *actor.Actor
}

// Returns MsgSpawn.
func (Spawn) Kind() message.Kind {
	return MsgSpawn
}

type world struct {
	actor.Actor
}

func (w world) Initialize() (message.Receiver, message.Sender) {
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
