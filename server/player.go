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

type Permission uint16

const (
	PermSecurity Permission = 1 << iota
	PermEngineer
	PermMedical
)

type Player struct {
	actor.Actor
	ID        uint64 // network ID (public)
	id        string // network ID (private)
	x, y      int
	flags     uint32
	perms     Permission
	send      chan<- *packet.Packet
	onmove    <-chan message.Message
	forcemove chan *packet.Packet
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

	canSee := make(map[uint64]bool)

	for {
		select {
		case msg, ok := <-msgIn:
			if !ok {
				return
			}
			switch m := msg.(type) {
			case SendLocation:
				select {
				case m.forcemove <- &packet.Packet{
					Location: &packet.Location{
						ID:    p.ID,
						Coord: layout.Coord{p.x, p.y},
						Flags: p.flags,
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
				// We missed a few ticks. No big deal.
				move = actor.Tick(time.Second / 2)
				continue
			}

			dx, dy := moveRequest.X, moveRequest.Y

			if dx == 0 && dy == 0 {
				move = nil
				continue
			}

			target := layout.Coord{p.x + dx, p.y + dy}
			canMove := layout.GetCoord(target).Passable()
			if canMove && dx != 0 && dy != 0 {
				canMove = canMove && (layout.Get(p.x+dx, p.y).Passable() ||
					layout.Get(p.x, p.y+dy).Passable())
			} else if !canMove && (dx == 0 || dy == 0) {
				if layout.GetCoord(target).Door() {
					go world.OpenDoor(p, target)
					continue
				}
			}

			if !canMove {
				moveRequest.X, moveRequest.Y = 0, 0
				move = nil
				continue
			}

			p.x += dx
			p.y += dy

			go func(m SetLocation) {
				world.Send <- m
			}(SetLocation{
				ID:    p.ID,
				Actor: &p.Actor,
				Flags: p.flags,
				Coord: layout.Coord{p.x, p.y},
			})

		case msg := <-p.forcemove:
			if layout.Visible(layout.Coord{p.x, p.y}, msg.Location.Coord) {
				canSee[msg.Location.ID] = true
				p.send <- msg
			}

		case msg := <-p.onmove:
			m := msg.(SetLocation)
			if layout.Visible(layout.Coord{p.x, p.y}, m.Coord) {
				canSee[m.ID] = true
				p.send <- &packet.Packet{
					Location: &packet.Location{
						ID:    m.ID,
						Coord: m.Coord,
						Flags: m.Flags,
					},
				}
			} else if canSee[m.ID] {
				p.send <- &packet.Packet{
					Despawn: &packet.Despawn{
						ID: m.ID,
					},
				}
				delete(canSee, m.ID)
			}
		}
	}
}

func (p *Player) HasPermissions(perm Permission) bool {
	return p.perms&perm == perm
}

func (p *Player) Disconnected() {
	world.Send <- packet.Despawn{ID: p.ID}
	world.Send <- actor.RemoveHeld{&p.Actor}
	close(p.Send)
}

func NewPlayer(id string, send chan<- *packet.Packet, flags uint32) (player *Player) {
	player = new(Player)
	player.id = id
	player.send = send
	player.flags = flags
	actor.Init("player:"+id, &player.Actor, player)

	send <- &packet.Packet{
		Handshake: &packet.Handshake{
			ID: player.ID,
		},
	}

	player.forcemove = make(chan *packet.Packet)

	var onmove actor.Subscriber
	onmove, player.onmove = actor.Subscribe(MsgSetLocation, 4)
	world.Send <- actor.AddSubscriber{onmove}

	world.Send <- actor.AddHeld{&player.Actor}
	world.onConnect <- player
	world.Send <- SetLocation{
		ID:    player.ID,
		Actor: &player.Actor,
		Flags: flags,
		Coord: layout.Coord{player.x, player.y},
	}

	return
}
