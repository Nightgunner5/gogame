package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"image"
	"sync/atomic"
	"time"
)

type Player struct {
	actor.Actor
	x, y          int
	isLocalPlayer bool
}

func (p *Player) Initialize() (message.Receiver, message.Sender) {
	msgIn, broadcast := p.Actor.Initialize()

	messages := make(chan message.Message)

	tick := actor.Tick(time.Second)

	go func() {
		for {
			select {
			case <-tick:
				Invalidate(p.screenRect())
				p.x++
				if p.x == ViewportWidth {
					p.x = 0
					p.y++
				}
				if p.y == ViewportHeight {
					p.y = 0
				}
				if p.isLocalPlayer {
					atomic.StoreInt64(&topLeftX, ViewportWidth/2-int64(p.x))
					atomic.StoreInt64(&topLeftY, ViewportHeight/2-int64(p.y))
				} else {
					Invalidate(p.screenRect())
				}
			case msg := <-msgIn:
				switch m := msg.(type) {
				case PaintRequest:
					m.Reply(0, p.x, p.y)
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

func NewPlayer(isLocalPlayer bool) (player Player) {
	player.isLocalPlayer = isLocalPlayer
	actor.TopLevel(player.Initialize())
	world.Send <- actor.AddHeld{&player.Actor}
	return
}
