package client

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"image/draw"
	"sync/atomic"
)

const (
	MaxQueue = 8
)

var (
	topLeftX, topLeftY int64 = ViewportWidth / 2, ViewportHeight / 2
)

func GetTopLeft() (x, y int) {
	x = int(atomic.LoadInt64(&topLeftX))
	y = int(atomic.LoadInt64(&topLeftY))
	return
}

var Network chan<- packet.Packet

type World struct {
	actor.Holder

	idToActor map[uint64]*actor.Actor
}

func (w *World) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := w.Holder.Initialize()

	w.idToActor = make(map[uint64]*actor.Actor)

	messages := make(chan message.Message)

	go func() {
		for msg := range msgIn {
			switch m := msg.(type) {
			case packet.Handshake:
				a := &thePlayer.Actor
				w.idToActor[m.ID] = a
				go func(a *actor.Actor) {
					w.Send <- actor.AddHeld{a}
				}(a)

			case packet.Location:
				id, coord := m.ID, m.Coord
				if _, ok := w.idToActor[id]; !ok {
					a := &NewPlayer(false).Actor
					w.idToActor[id] = a
					go func(a *actor.Actor) {
						w.Send <- actor.AddHeld{a}
					}(a)
				}
				w.idToActor[id].Send <- SetLocation{coord}

			case MoveRequest:
				if m.X == 0 && m.Y == 0 {
					continue
				}
				var dx, dy int
				if m.X*m.X > m.Y*m.Y {
					if m.X > 0 {
						dx = 1
					} else {
						dx = -1
					}
				} else {
					if m.Y > 0 {
						dy = 1
					} else {
						dy = -1
					}
				}
				Network <- packet.Packet{
					Location: &packet.Location{
						Coord: layout.Coord{dx, dy},
					},
				}

			default:
				messages <- m
			}
		}
		close(messages)
	}()

	return messages, broadcast
}

var world = NewWorld()

func NewWorld() (world *World) {
	world = new(World)
	actor.Init("client:world", &world.Actor, world)
	return
}

var (
	MsgPaintRequest = message.NewKind("PaintRequest")
)

type PaintRequest chan<- PaintContext

func (p PaintRequest) Reply(spriteID uint16, x, y int) {
	p <- PaintContext{spriteID, x, y}
}

func (p PaintRequest) Kind() message.Kind {
	return MsgPaintRequest
}

type PaintContext struct {
	spriteID uint16
	x, y     int
}

func (p PaintContext) Paint(viewport draw.Image, xOffset, yOffset int) {
	Tile(viewport, Actors, p.spriteID, p.x+xOffset, p.y+yOffset)
}

var (
	MsgMoveRequest = message.NewKind("MoveRequest")
)

type MoveRequest struct {
	X, Y int
}

func (MoveRequest) Kind() message.Kind {
	return MsgMoveRequest
}
