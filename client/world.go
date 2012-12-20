package client

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
)

var Network chan<- *packet.Packet

type World struct {
	actor.Holder

	idToActor map[uint64]*actor.Actor
}

func (w *World) Initialize() (message.Receiver, func(message.Message)) {
	msgIn, broadcast := w.Holder.Initialize()

	w.idToActor = make(map[uint64]*actor.Actor)

	messages := make(chan message.Message)

	go w.dispatch(msgIn, messages)

	return messages, broadcast
}

func (w *World) dispatch(msgIn message.Receiver, messages message.Sender) {
	for msg := range msgIn {
		switch m := msg.(type) {
		case packet.Handshake:
			a := &thePlayer.Actor
			w.idToActor[m.ID] = a
			go w.addHeld(a)

		case packet.Location:
			id, coord := m.ID, m.Coord
			if _, ok := w.idToActor[id]; !ok {
				a := &NewPlayer(false, m.Flags&packet.FlagSpriteMask == packet.FlagMonkey).Actor
				w.idToActor[id] = a
				go w.addHeld(a)
			}
			w.idToActor[id].Send <- SetLocation{coord, m.Flags}

		case packet.Despawn:
			if a, ok := w.idToActor[m.ID]; ok {
				delete(w.idToActor, m.ID)
				go w.removeHeld(a)
				a.Send <- m
				close(a.Send)
			}

		case MoveRequest:
			if m.X == 0 && m.Y == 0 {
				continue
			}
			var dx, dy int
			if m.X*m.X > m.Y*m.Y {
				if m.X > 0 {
					dx = 1
				} else {
					dx = -1
				}
			} else {
				if m.Y > 0 {
					dy = 1
				} else {
					dy = -1
				}
			}
			Network <- &packet.Packet{
				Location: &packet.Location{
					Coord: layout.Coord{dx, dy},
				},
			}

		default:
			messages <- m
		}
	}
	close(messages)
}

func (w *World) addHeld(a *actor.Actor) {
	w.Send <- actor.AddHeld{a}
}

func (w *World) removeHeld(a *actor.Actor) {
	w.Send <- actor.RemoveHeld{a}
}

var world = NewWorld()

func NewWorld() *World {
	world := new(World)
	actor.Init("client:world", &world.Actor, world)
	return world
}

var (
	MsgMoveRequest = message.NewKind("MoveRequest")
)

type MoveRequest struct {
	X, Y int
}

func (MoveRequest) Kind() message.Kind {
	return MsgMoveRequest
}
