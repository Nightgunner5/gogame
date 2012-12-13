package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
)

type Door struct {
	actor.Actor
	coord layout.Coord
	open  bool
}

func (d *Door) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := d.Actor.Initialize()

	messages := make(chan message.Message)

	go d.dispatch(msgIn, messages)

	return messages, broadcast
}

func (d *Door) dispatch(msgIn message.Receiver, messages message.Sender) {
	defer close(messages)
	for msg := range msgIn {
		switch m := msg.(type) {
		default:
			messages <- m
		}
	}
}

func NewDoor(coord layout.Coord) (door *Door) {
	door = new(Door)
	door.coord = coord
	tile := layout.Get(coord.X, coord.Y)
	if !tile.Door() {
		panic("NewDoor on non-door coordinate")
	}
	door.open = tile.Passable()

	actor.Init("door:"+coord.String(), &door.Actor, door)

	world.Send <- actor.AddHeld{&door.Actor}
	return
}
