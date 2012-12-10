package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"net"
)

var nextPlayerID = make(chan uint64)

func init() {
	go func() {
		var id uint64
		for {
			id++
			nextPlayerID <- id
		}
	}()
}

type Player struct {
	actor.Actor
	ID   uint64
	x, y int
	addr net.Addr
}

func (p *Player) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := p.Actor.Initialize()

	p.ID = <-nextPlayerID

	messages := make(chan message.Message)

	go func() {
		for {
			select {
			case msg := <-msgIn:
				switch m := msg.(type) {
				case MoveRequest:
					if m.Dx*m.Dx+m.Dy*m.Dy == 1 {
						// TODO: space logic
						if !layout.Get(p.x+m.Dx, p.y+m.Dy).Passable() {
							return
						}
						p.x += m.Dx
						p.y += m.Dy
						world.Send <- SetLocation{
							ID:    p.ID,
							Actor: &p.Actor,
							Coord: layout.Coord{p.x, p.y},
						}
					}
				default:
					messages <- m
				}
			}
		}
	}()

	return messages, broadcast
}

func NewPlayer(addr net.Addr) (player Player) {
	player.addr = addr
	actor.TopLevel(player.Initialize())
	world.Send <- actor.AddHeld{&player.Actor}
	world.Send <- SetLocation{
		Actor: &player.Actor,
		Coord: layout.Coord{},
	}
	return
}
