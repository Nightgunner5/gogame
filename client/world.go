package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/Nightgunner5/netchan"
	"image/draw"
	"net"
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

type World struct {
	actor.Holder

	in  <-chan packet.Packet
	out chan<- packet.Packet

	idToActor map[uint64]*actor.Actor
	actorToID map[*actor.Actor]uint64
}

func (w *World) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := w.Holder.Initialize()

	w.idToActor = make(map[uint64]*actor.Actor)
	w.actorToID = make(map[*actor.Actor]uint64)

	messages := make(chan message.Message)

	go func() {
		for {
			select {
			case msg := <-msgIn:
				switch m := msg.(type) {
				default:
					messages <- m
				}
			case pkt := <-w.in:
				switch {
				case pkt.HandshakeServer != nil:
					id := pkt.HandshakeServer.ID
					w.idToActor[id] = &thePlayer.Actor
					w.actorToID[&thePlayer.Actor] = id
					w.Send <- actor.AddHeld{&thePlayer.Actor}

				case pkt.PlayerLocation != nil:
					id, coord := pkt.PlayerLocation.ID, pkt.PlayerLocation.Coord
					if _, ok := w.idToActor[id]; !ok {
						a := &NewPlayer(false).Actor
						w.idToActor[id] = a
						w.actorToID[a] = id
						w.Send <- actor.AddHeld{a}
					}
					w.idToActor[id].Send <- SetLocation{coord}
				}
			}
		}
	}()

	return messages, broadcast
}

var world = NewWorld()

func NewWorld() (world World) {
	conn, err := net.Dial("tcp", "nightgunner5.is-a-geek.net:7031")
	if err != nil {
		panic(err)
	}
	c := netchan.New(conn, packet.Type, MaxQueue)
	world.in, world.out = c.ChanRecv().(<-chan packet.Packet), c.ChanSend().(chan<- packet.Packet)
	actor.TopLevel(world.Initialize())
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
