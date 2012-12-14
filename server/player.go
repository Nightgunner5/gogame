package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"time"
)

var nextPlayerID = make(chan uint64)

func generatePlayerIDs() {
	var id uint64
	for {
		id++
		nextPlayerID <- id
	}
}

func init() {
	go generatePlayerIDs()
}

type Player struct {
	actor.Actor
	ID   uint64 // network ID (public)
	id   string // network ID (private)
	Name string
	x, y int
	out  chan<- packet.Packet
}

func (p *Player) Initialize() (message.Receiver, func(message.Message)) {
	msgIn, broadcast := p.Actor.Initialize()

	p.ID = <-nextPlayerID

	messages := make(chan message.Message)

	go p.dispatch(msgIn, messages)

	return messages, broadcast
}

func (p *Player) dispatch(msgIn message.Receiver, messages message.Sender) {
	defer close(messages)

	var moveRequest layout.Coord
	var move actor.Ticker

	for {
		select {
		case msg, ok := <-msgIn:
			if !ok {
				return
			}
			switch m := msg.(type) {
			case SendLocation:
				select {
				case m <- packet.Packet{
					Location: &packet.Location{
						ID:    p.ID,
						Coord: layout.Coord{p.x, p.y},
					},
				}:
				case <-time.After(20 * time.Millisecond):
				}

			case packet.Location:
				moveRequest = m.Coord
				if move == nil {
					move = actor.Tick(time.Second / 2)
					move <- struct{}{}
				}

			default:
				messages <- m
			}

		case _, ok := <-move:
			if !ok {
				// We missed at least 10 ticks, which means the
				// server is lagging quite heavily. Restart the
				// ticker as we aren't dead just yet.
				move = actor.Tick(time.Second / 2)
				continue
			}

			dx, dy := moveRequest.X, moveRequest.Y

			if dx == 0 && dy == 0 {
				move = nil
				continue
			}

			canMove := layout.Get(p.x+dx, p.y+dy).Passable()
			if canMove && dx != 0 && dy != 0 {
				canMove = canMove && (layout.Get(p.x+dx, p.y).Passable() ||
					layout.Get(p.x, p.y+dy).Passable())
			}

			// TODO: space logic
			if canMove {
				p.x += dx
				p.y += dy

				go func(m SetLocation) {
					world.Send <- m
				}(SetLocation{
					ID:    p.ID,
					Actor: &p.Actor,
					Coord: layout.Coord{p.x, p.y},
				})
			}
		}
	}
}

func (p *Player) Disconnected() {
	world.Send <- packet.Despawn{ID: p.ID}
	world.Send <- actor.RemoveHeld{&p.Actor}
	close(p.Send)
}

func NewPlayer(id string, name string, out chan<- packet.Packet) (player *Player) {
	player = new(Player)
	player.id = id
	player.Name = name
	player.out = out
	actor.Init("player:"+id, &player.Actor, player)

	out <- packet.Packet{
		Handshake: &packet.Handshake{
			ID: player.ID,
		},
	}

	world.Send <- actor.AddHeld{&player.Actor}
	world.onConnect <- out
	return
}
