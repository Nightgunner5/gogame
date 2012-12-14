package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"time"
)

type Door struct {
	actor.Actor
	coord layout.Coord
	open  bool
}

func (d *Door) Initialize() (message.Receiver, func(message.Message)) {
	msgIn, broadcast := d.Actor.Initialize()

	messages := make(chan message.Message)

	go d.dispatch(msgIn, messages)

	return messages, broadcast
}

func (d *Door) dispatch(msgIn message.Receiver, messages message.Sender) {
	defer close(messages)
	var closeDoor <-chan time.Time
	for {
		select {
		case msg, ok := <-msgIn:
			if !ok {
				return
			}

			switch m := msg.(type) {
			case OpenDoor:
				if !d.open {
					d.open = true
					for {
						orig := layout.GetCoord(d.coord)
						tile := make(layout.MultiTile, len(orig))
						copy(tile, orig)

						for i := range tile {
							if tile[i].Door() {
								tile[i] &^= 1
							}
						}

						if layout.SetCoord(d.coord, orig, tile) {
							break
						}
					}
					closeDoor = time.After(10 * time.Second)
				}

			default:
				messages <- m
			}

		case <-closeDoor:
			if d.open {
				d.open = false
				for {
					orig := layout.GetCoord(d.coord)
					tile := make(layout.MultiTile, len(orig))
					copy(tile, orig)

					for i := range tile {
						if tile[i].Door() {
							tile[i] |= 1
						}
					}

					if layout.SetCoord(d.coord, orig, tile) {
						break
					}
				}
			}
		}
	}
}

func NewDoor(coord layout.Coord) *Door {
	door := new(Door)
	door.coord = coord
	tile := layout.Get(coord.X, coord.Y)
	if !tile.Door() {
		panic("NewDoor on non-door coordinate")
	}
	door.open = tile.Passable()

	actor.Init("door:"+coord.String(), &door.Actor, door)

	world.Send <- AddDoor{coord, &door.Actor}
	return door
}

var (
	MsgAddDoor  = message.NewKind("AddDoor")
	MsgOpenDoor = message.NewKind("OpenDoor")
)

type AddDoor struct {
	Coord layout.Coord
	Actor *actor.Actor
}

func (AddDoor) Kind() message.Kind {
	return MsgAddDoor
}

type OpenDoor struct {
	Opener *Player
}

func (OpenDoor) Kind() message.Kind {
	return MsgOpenDoor
}
