package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
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
	ID   uint64 // network ID (public)
	id   string // network ID (private)
	Name string
	x, y int
	out  chan<- packet.Packet
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
				case SendLocation:
					m <- packet.Packet{
						Location: &packet.Location{
							ID:    p.ID,
							Coord: layout.Coord{p.x, p.y},
						},
					}
				case packet.Location:
					dx, dy := m.Coord.X, m.Coord.Y
					// TODO: space logic
					if !layout.Get(p.x+dx, p.y+dy).Passable() {
						continue
					}
					p.x += dx
					p.y += dy

					go func(m SetLocation) {
						world.Send <- m
					}(SetLocation{
						ID:    p.ID,
						Actor: &p.Actor,
						Coord: layout.Coord{p.x, p.y},
					})

				default:
					messages <- m
				}
			}
		}
	}()

	return messages, broadcast
}

func (p *Player) Disconnected() {
	world.Send <- actor.RemoveHeld{&p.Actor}
	world.Send <- packet.Despawn{
		ID: p.ID,
	}
}

func NewPlayer(id string, name string, out chan<- packet.Packet) (player *Player) {
	player = new(Player)
	player.id = id
	player.Name = name
	player.out = out
	actor.TopLevel(player.Initialize())
	out <- packet.Packet{
		Handshake: &packet.Handshake{
			ID: player.ID,
		},
	}

	go func() {
		world.Send <- actor.AddHeld{&player.Actor}
		world.onConnect <- out
	}()
	return
}
