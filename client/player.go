package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"image"
	"time"
)

type Player struct {
	actor.Actor
	x, y int
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
				Invalidate(p.screenRect())
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
	return image.Rect(p.x<<TileSize, p.y<<TileSize,
		(p.x+1)<<TileSize, (p.y+1)<<TileSize)
}

var thePlayer = NewPlayer()

func NewPlayer() (player Player) {
	actor.TopLevel(player.Initialize())
	world.Send <- actor.AddHeld{&player.Actor}
	return
}
