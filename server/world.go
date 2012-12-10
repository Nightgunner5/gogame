package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
)

type World struct {
	actor.Holder
	idToActor map[uint64]*actor.Actor
	actorToID map[*actor.Actor]uint64
	location  map[*actor.Actor]layout.Coord
}

func (w *World) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := w.Holder.Initialize()

	messages := make(chan message.Message)
	w.idToActor = make(map[uint64]*actor.Actor)
	w.actorToID = make(map[*actor.Actor]uint64)
	w.location = make(map[*actor.Actor]layout.Coord)

	go func() {
		for {
			select {
			case msg := <-msgIn:
				switch m := msg.(type) {
				case SetLocation:
					w.actorToID[m.Actor] = m.ID
					w.idToActor[m.ID] = m.Actor
					w.location[m.Actor] = m.Coord
				}
			}
		}
	}()

	return messages, broadcast
}

var world = NewWorld()

func NewWorld() (world World) {
	actor.TopLevel(world.Initialize())
	return
}
