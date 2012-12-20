package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"log"
	"math/rand"
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
	ID         uint64 // network ID (public)
	id         string // network ID (private)
	coord      layout.Coord
	flags      uint32
	hasSetRole bool
	perms      Permission
	send       chan<- *packet.Packet
	onmove     <-chan message.Message
	forcemove  chan *packet.Packet
}

func (p *Player) Initialize() (message.Receiver, func(message.Message)) {
	msgIn, broadcast := p.Actor.Initialize()

	p.ID = <-nextPlayerID

	messages := make(chan message.Message)

	go p.dispatch(msgIn, messages)

	return messages, broadcast
}

func (p *Player) dispatch(msgIn message.Receiver, messages message.Sender) {
	// If the player manages to disconnect while this function is processing
	// one of their packets, responding to them will panic with "send on
	// closed channel". This is a race condition that would require more
	// work than it's worth - since the connection is already closed,
	// just drop the player.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in player:%s: %v", p.id, r)
			p.Disconnected()
		}
	}()

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
						Coord: p.coord,
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

			case LayoutChanged:
				if layout.Visible(p.coord, m.Coord) {
					go func(m SetLocation) {
						world.Send <- m
					}(SetLocation{
						Actor: &p.Actor,
						Packet: &packet.Packet{
							Location: &packet.Location{
								ID:    p.ID,
								Coord: p.coord,
								Flags: p.flags,
							},
						},
					})
				}

			case MoveIntoView:
				if layout.Visible(p.coord, m.Location.Coord) {
					p.send <- m.Packet
					canSee[m.Location.ID] = true
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

			x, y := p.coord.X, p.coord.Y
			target := layout.Coord{x + dx, y + dy}
			tile := layout.GetCoord(target)
			canMove := tile.Passable()
			if canMove && dx != 0 && dy != 0 {
				canMove = canMove && (layout.Get(x+dx, y).Passable() ||
					layout.Get(x, y+dy).Passable())
			} else if !canMove && (dx == 0 || dy == 0) {
				if tile.Door() {
					go world.OpenDoor(p, target)
					continue
				}
			}

			if !canMove {
				moveRequest.X, moveRequest.Y = 0, 0
				move = nil
				continue
			}

			for _, t := range tile {
				switch t {
				case layout.TriggerSelectRole:
					if !p.hasSetRole {
						var flags uint32
						switch rand.Intn(3) {
						case 0:
							flags = packet.FlagSecurity
							p.perms = PermSecurity
						case 1:
							flags = packet.FlagEngineer
							p.perms = PermEngineer
						case 2:
							flags = packet.FlagMedic
							p.perms = PermMedical
						}
						if p.flags&packet.FlagSpriteMask != packet.FlagMonkey {
							p.flags &^= packet.FlagSpriteMask
							p.flags |= flags
						}
						p.hasSetRole = true
					}
				}
			}

			p.coord.X += dx
			p.coord.Y += dy

			go func(m SetLocation) {
				world.Send <- m
			}(SetLocation{
				Actor: &p.Actor,
				Packet: &packet.Packet{
					Location: &packet.Location{
						ID:    p.ID,
						Coord: p.coord,
						Flags: p.flags,
					},
				},
			})

		case msg := <-p.forcemove:
			if layout.Visible(p.coord, msg.Location.Coord) {
				canSee[msg.Location.ID] = true
				p.send <- msg
			} else if canSee[msg.Location.ID] {
				p.send <- &packet.Packet{
					Despawn: &packet.Despawn{
						ID: msg.Location.ID,
					},
				}
				delete(canSee, msg.Location.ID)
			}

		case msg := <-p.onmove:
			m := msg.(SetLocation)
			if layout.Visible(p.coord, m.Location.Coord) {
				if !canSee[m.Location.ID] {
					go func(a *actor.Actor, m MoveIntoView) {
						defer func() {
							recover()
						}()

						a.Send <- m
					}(m.Actor, MoveIntoView{
						Packet: m.Packet,
					})
					canSee[m.Location.ID] = true
				}
				p.send <- m.Packet
			} else if canSee[m.Location.ID] {
				p.send <- &packet.Packet{
					Despawn: &packet.Despawn{
						ID: m.Location.ID,
					},
				}
				delete(canSee, m.Location.ID)
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
		Actor: &player.Actor,
		Packet: &packet.Packet{
			Location: &packet.Location{
				ID:    player.ID,
				Coord: player.coord,
				Flags: flags,
			},
		},
	}

	return
}
