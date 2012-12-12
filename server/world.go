package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"time"
)

var (
	SendToAll = make(chan packet.Packet)
)

type World struct {
	actor.Holder

	onConnect chan<- chan<- packet.Packet

	idToActor map[uint64]*actor.Actor
	location  map[*actor.Actor]layout.Coord
}

func (w *World) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := w.Holder.Initialize()

	onConnect := make(chan chan<- packet.Packet)
	w.onConnect = onConnect

	w.idToActor = make(map[uint64]*actor.Actor)
	w.location = make(map[*actor.Actor]layout.Coord)

	messages := make(chan message.Message)

	go w.dispatch(msgIn, messages, onConnect)

	return messages, broadcast
}

func (w *World) dispatch(msgIn messages.Receiver, messages message.Sender, onConnect <-chan chan<- packet.Packet) {
	for {
		select {
		case msg, ok := <-msgIn:
			if !ok {
				close(messages)
				return
			}

			switch m := msg.(type) {
			case SetLocation:
				w.idToActor[m.ID] = m.Actor
				w.location[m.Actor] = m.Coord

				SendToAll <- packet.Packet{
					Location: &packet.Location{
						ID:    m.ID,
						Coord: m.Coord,
					},
				}

			case packet.Despawn:
				a := w.idToActor[m.ID]
				delete(w.idToActor, m.ID)
				delete(w.location, a)
				SendToAll <- packet.Packet{
					Despawn: &m,
				}

			default:
				messages <- m
			}

		case c := <-onConnect:
			go func(c SendLocation) {
				for _, a := range w.GetHeld() {
					select {
					case a.Send <- c:
					case <-time.After(10 * time.Millisecond):
					}
				}
			}(SendLocation(c))
		}
	}
}

var world = NewWorld()

func NewWorld() (world *World) {
	world = new(World)
	actor.Init("world", &world.Actor, world)
	return
}
