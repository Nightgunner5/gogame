package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"image/draw"
	"sync/atomic"
)

var (
	topLeftX, topLeftY int64 = ViewportWidth/2, ViewportHeight/2
)

func GetTopLeft() (x, y int) {
	x = int(atomic.LoadInt64(&topLeftX))
	y = int(atomic.LoadInt64(&topLeftY))
	return
}

type World struct {
	actor.Holder
}

var world = NewWorld()

func NewWorld() (world World) {
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
