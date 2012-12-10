package main

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"image/draw"
)

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

func (p PaintRequest) Reply(spriteID, x, y int) {
	p <- PaintContext{spriteID, x, y}
}

func (p PaintRequest) Kind() message.Kind {
	return MsgPaintRequest
}

type PaintContext struct {
	spriteID, x, y int
}

func (p PaintContext) Paint(viewport draw.Image) {
	Tile(viewport, Actors, p.spriteID, p.x, p.y)
}
