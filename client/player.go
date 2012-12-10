package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"image"
	"sync/atomic"
)

type Player struct {
	actor.Actor
	x, y          int
	isLocalPlayer bool
}

func (p *Player) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := p.Actor.Initialize()

	messages := make(chan message.Message)

	go func() {
		for {
			select {
			case msg := <-msgIn:
				switch m := msg.(type) {
				case PaintRequest:
					m.Reply(0, p.x, p.y)

				case SetLocation:
					Invalidate(p.screenRect())
					p.x, p.y = m.X, m.Y
					if p.isLocalPlayer {
						atomic.StoreInt64(&topLeftX, ViewportWidth/2-int64(p.x))
						atomic.StoreInt64(&topLeftY, ViewportHeight/2-int64(p.y))
					} else {
						Invalidate(p.screenRect())
					}

				default:
					messages <- m
				}
			}
		}
	}()

	return messages, broadcast
}

func (p *Player) screenRect() image.Rectangle {
	if p.isLocalPlayer {
		return image.Rect(0, 0, ViewportWidth<<TileSize, ViewportHeight<<TileSize)
	}
	x, y := GetTopLeft()
	return image.Rect((p.x+x)<<TileSize, (p.y+y)<<TileSize,
		(p.x+1+x)<<TileSize, (p.y+1+y)<<TileSize)
}

var thePlayer = NewPlayer(true)

func NewPlayer(isLocalPlayer bool) (player *Player) {
	player = new(Player)
	player.isLocalPlayer = isLocalPlayer
	actor.TopLevel(player.Initialize())
	return
}

var (
	MsgSetLocation = message.NewKind("SetLocation")
)

type SetLocation struct {
	layout.Coord
}

func (SetLocation) Kind() message.Kind {
	return MsgSetLocation
}

