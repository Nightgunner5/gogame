package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
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
	in   <-chan packet.Packet
	out  chan<- packet.Packet
}

func (p *Player) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := p.Actor.Initialize()

	p.ID = <-nextPlayerID
	p.out <- packet.Packet{
		HandshakeServer: &packet.HandshakeServer{
			ID: p.ID,
		},
	}

	messages := make(chan message.Message)

	go func() {
		for {
			select {
			case pkt := <-p.in:
				switch {
				case pkt.MoveRequest != nil:
					dx, dy := pkt.MoveRequest.Dx, pkt.MoveRequest.Dy
					if dx*dx+dy*dy == 1 {
						// TODO: space logic
						if !layout.Get(p.x+dx, p.y+dy).Passable() {
							return
						}
						p.x += dx
						p.y += dy
						go func() {
							world.Send <- SetLocation{
								ID:    p.ID,
								Actor: &p.Actor,
								Coord: layout.Coord{p.x, p.y},
							}
						}()
					}
				}
			case msg := <-msgIn:
				switch m := msg.(type) {
				case SendLocation:
					m <- packet.Packet{
						PlayerLocation: &packet.PlayerLocation{
							ID:    p.ID,
							Coord: layout.Coord{p.x, p.y},
						},
					}
				default:
					messages <- m
				}
			}
		}
	}()

	return messages, broadcast
}

func NewPlayer(addr net.Addr, in <-chan packet.Packet, out chan<- packet.Packet) (player Player) {
	player.addr = addr
	player.in, player.out = in, out
	actor.TopLevel(player.Initialize())
	world.Send <- actor.AddHeld{&player.Actor}
	world.Send <- SetLocation{
		Actor: &player.Actor,
		Coord: layout.Coord{},
	}
	return
}
