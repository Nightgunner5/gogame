package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/power"
	"time"
)

type Door struct {
	actor.Actor
	coord      layout.Coord
	open       bool
	permission Permission
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
				if !d.open && power.Powered(d.coord.X, d.coord.Y) && m.Opener.HasPermissions(d.permission) {
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
			closeDoor = nil
			if d.open && power.Powered(d.coord.X, d.coord.Y) {
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

func init() {
	layout.AllTiles(func(coord layout.Coord, tile layout.MultiTile) {
		if tile.Door() {
			door := NewDoor(coord)
			world.Send <- AddDoor{coord, &door.Actor}
		}
	})
}

func NewDoor(coord layout.Coord) *Door {
	door := new(Door)
	door.coord = coord
	tile := layout.Get(coord.X, coord.Y)
	if !tile.Door() {
		panic("NewDoor on non-door coordinate")
	}
	door.open = tile.Passable()
	for _, t := range tile {
		switch t {
		case layout.DoorGeneralClosed, layout.DoorGeneralOpen:
			// No extra permission
		case layout.DoorSecurityClosed, layout.DoorSecurityOpen:
			door.permission |= PermSecurity
		case layout.DoorEngineerClosed, layout.DoorEngineerOpen:
			door.permission |= PermEngineer
		case layout.DoorMedicalClosed, layout.DoorMedicalOpen:
			door.permission |= PermMedical
		}
	}

	actor.Init("door:"+coord.String(), &door.Actor, door)

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
