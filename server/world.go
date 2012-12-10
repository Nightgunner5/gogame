package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
)

type World struct {
	actor.Holder

	onConnect chan<- chan<- packet.Packet

	idToActor map[uint64]*actor.Actor
	actorToID map[*actor.Actor]uint64
	location  map[*actor.Actor]layout.Coord
}

func (w *World) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := w.Holder.Initialize()

	onConnect := make(chan chan<- packet.Packet)
	w.onConnect = onConnect

	w.idToActor = make(map[uint64]*actor.Actor)
	w.actorToID = make(map[*actor.Actor]uint64)
	w.location = make(map[*actor.Actor]layout.Coord)

	messages := make(chan message.Message)

	go func() {
		for {
			select {
			case msg := <-msgIn:
				switch m := msg.(type) {
				case SetLocation:
					w.actorToID[m.Actor] = m.ID
					w.idToActor[m.ID] = m.Actor
					w.location[m.Actor] = m.Coord

					sendToAll <- packet.Packet{
						PlayerLocation: &packet.PlayerLocation{
							ID:    m.ID,
							Coord: m.Coord,
						},
					}

				case Despawn:
					delete(w.actorToID, m.Actor)
					delete(w.idToActor, m.ID)
					delete(w.location, m.Actor)
					// TODO: despawn on client side

				default:
					messages <- m
				}

			case c := <-onConnect:
				for a := range w.actorToID {
					a.Send <- SendLocation(c)
				}
			}
		}
	}()

	return messages, broadcast
}

var world = NewWorld()

func NewWorld() (world *World) {
	world = new(World)
	actor.TopLevel(world.Initialize())
	return
}
