package client

import (
	"bytes"
	"github.com/Nightgunner5/gogame/client/res"
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"image"
	"image/draw"
	"image/png"
	"log"
	"sync"
	"time"
)

var (
	Terrain image.Image
	Actors  image.Image
)

func init() {
	var err error

	Terrain, err = png.Decode(bytes.NewReader(res.TerrainPng[:]))
	if err != nil {
		panic(err)
	}

	Actors, err = png.Decode(bytes.NewReader(res.ActorsPng[:]))
	if err != nil {
		panic(err)
	}
}

func Main() {
	go UI()

	wde.Run()
}

const (
	TileSize  = 5 // 2**5 pixels by 2**5 pixels
	TileShift = 5 // 2**5 tiles by 2**5 tiles
	TileMask  = 1<<TileShift - 1

	ViewportWidth  = 20
	ViewportHeight = 15
)

func tileCoord(index uint16) (p image.Point) {
	return image.Pt(int((index&TileMask)<<TileSize),
		int((index>>TileShift&TileMask)<<TileSize))
}

func Tile(viewport draw.Image, base image.Image, index uint16, x, y int) {
	x, y = x<<TileSize, y<<TileSize
	draw.Draw(viewport, image.Rect(x, y, x+1<<TileSize, y+1<<TileSize), base, tileCoord(index), draw.Over)
}

func Paint(w wde.Window, rect image.Rectangle) {
	viewport := w.Screen()

	xOffset, yOffset := GetTopLeft()

	for x := rect.Min.X >> TileSize; x < (rect.Max.X-1)>>TileSize+1; x++ {
		for y := rect.Min.Y >> TileSize; y < (rect.Max.Y-1)>>TileSize+1; y++ {
			Tile(viewport, Terrain, uint16(layout.GetSpace(x-xOffset, y-yOffset)), x, y)
			for _, t := range layout.Get(x-xOffset, y-yOffset) {
				Tile(viewport, Terrain, uint16(t), x, y)
			}
		}
	}

	paintLock.Lock()
	for _, p := range paintContexts {
		Tile(viewport, Actors, p.Sprite, p.Coord.X+xOffset, p.Coord.Y+yOffset)
	}
	paintLock.Unlock()

	w.FlushImage(rect)
}

var (
	shouldPaint = make(chan image.Rectangle, 1)
)

func Invalidate(rect image.Rectangle) {
	for {
		select {
		case shouldPaint <- rect:
			return
		case r2 := <-shouldPaint:
			rect = rect.Union(r2)
		}
	}
}

func paintHandler(w wde.Window) {
	for rect := range shouldPaint {
		Paint(w, rect)
		time.Sleep(16 * time.Millisecond)
	}
}

func UI() {
	defer wde.Stop()

	w, err := wde.NewWindow(ViewportWidth<<TileSize, ViewportHeight<<TileSize)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	w.SetTitle("GoGame")

	w.Show()

	go paintHandler(w)

	Invalidate(w.Screen().Bounds())

	keys := make(map[string]bool)

	for event := range w.EventChan() {
		switch e := event.(type) {
		case wde.MouseMovedEvent:
		case wde.MouseDownEvent:
		case wde.MouseUpEvent:
		case wde.MouseDraggedEvent:
		case wde.MouseEnteredEvent:
		case wde.MouseExitedEvent:
		case wde.KeyDownEvent:
			if !keys[e.Key] {
				keys[e.Key] = true
				if f := keyDown[e.Key]; f != nil {
					f(keys)
				}
			}

		case wde.KeyUpEvent:
			if keys[e.Key] {
				delete(keys, e.Key)
				if f := keyUp[e.Key]; f != nil {
					f(keys)
				}
			}

		case wde.KeyTypedEvent:
		case wde.ResizeEvent:
			w.SetSize(ViewportWidth<<TileSize, ViewportHeight<<TileSize)
			Invalidate(w.Screen().Bounds())

		case wde.CloseEvent:
			return
		}
	}
}

func Disconnected() {
	log.Fatal("Disconnected")
}

func Handle(msg packet.Packet) {
	switch {
	case msg.Handshake != nil:
		world.Send <- *msg.Handshake

	case msg.Location != nil:
		world.Send <- *msg.Location

	case msg.Despawn != nil:
		world.Send <- *msg.Despawn

	default:
		log.Fatalf("unknown packet: %#v", msg)
	}
}

var (
	paintContexts = make(map[*actor.Actor]*PaintContext)
	paintLock     sync.RWMutex
)

type PaintContext struct {
	Coord  layout.Coord
	Sprite uint16
}
