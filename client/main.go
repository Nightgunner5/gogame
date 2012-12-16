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
	"os"
	"sync"
	"time"
)

var (
	Terrain image.Image
	Actors  image.Image
)

func init() {
	var err error

	Terrain, err = png.Decode(bytes.NewReader(res.TerrainPng))
	if err != nil {
		panic(err)
	}
	res.TerrainPng = nil

	Actors, err = png.Decode(bytes.NewReader(res.ActorsPng))
	if err != nil {
		panic(err)
	}
	res.ActorsPng = nil
}

func Main() {
	layout.OnChange = func(c layout.Coord, t layout.MultiTile) {
		if t.Door() {
			Invalidate(image.Rect(0, 0, ViewportWidth<<TileSize, ViewportHeight<<TileSize))
		} else {
			xOffset, yOffset := GetTopLeft()
			x, y := c.X+xOffset, c.Y+yOffset
			Invalidate(image.Rect(x<<TileSize, y<<TileSize, (x+1)<<TileSize, (y+1)<<TileSize))
		}
	}

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

func TileFloat(viewport draw.Image, base image.Image, index uint16, x1, y1, x2, y2 int, interp float32) {
	x1, y1 = x1<<TileSize, y1<<TileSize
	x2, y2 = x2<<TileSize, y2<<TileSize
	x, y := x1+int(float32(x2-x1)*interp), y1+int(float32(y2-y1)*interp)
	draw.Draw(viewport, image.Rect(x, y, x+1<<TileSize, y+1<<TileSize), base, tileCoord(index), draw.Over)
}

var (
	viewport = image.NewRGBA(image.Rect(0, 0, ViewportWidth<<TileSize, ViewportHeight<<TileSize))
	space    = image.NewRGBA(viewport.Bounds())
	scene    = image.NewRGBA(viewport.Bounds())
	light    = new(lighting)
)

func Paint(w wde.Window, rect image.Rectangle) {
	xOffset, yOffset := GetTopLeft()
	center := layout.Coord{ViewportWidth/2 - xOffset, ViewportHeight/2 - yOffset}

	minX, maxX := rect.Min.X>>TileSize, (rect.Max.X-1)>>TileSize+1
	minY, maxY := rect.Min.Y>>TileSize, (rect.Max.Y-1)>>TileSize+1

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			if layout.Visible(center, layout.Coord{x - xOffset, y - yOffset}) {
				Tile(space, Terrain, uint16(layout.GetSpace(x-xOffset, y-yOffset)), x, y)
				draw.Draw(scene, image.Rect(x<<TileSize, y<<TileSize, (x+1)<<TileSize, (y+1)<<TileSize), image.Transparent, image.ZP, draw.Src)
				for _, t := range layout.Get(x-xOffset, y-yOffset) {
					Tile(scene, Terrain, uint16(t), x, y)
				}
			} else {
				Tile(scene, image.Black, 0, x, y)
			}
		}
	}

	hasAnimation := false
	paintLock.Lock()
	for _, p := range paintContexts {
		if layout.Visible(center, p.To) {
			x1, y1 := p.From.X+xOffset, p.From.Y+yOffset
			x2, y2 := p.To.X+xOffset, p.To.Y+yOffset

			if minX <= x2 && x2 <= maxX && minY <= y2 && y2 <= maxY {
				interp := float32(time.Since(p.Changed)*2) / float32(time.Second)
				if interp > 1 {
					Tile(scene, Actors, p.Sprite, x2, y2)
				} else {
					hasAnimation = true
					TileFloat(scene, Actors, p.Sprite, x1, y1, x2, y2, interp)
				}
			}
		}
	}
	paintLock.Unlock()

	draw.Draw(viewport, rect, space, rect.Min, draw.Src)
	draw.Draw(viewport, rect, scene, rect.Min, draw.Over)
	drawLightOverlay(viewport, rect, light.Image(-xOffset, -yOffset), rect.Min.Add(light.Origin(-xOffset, -yOffset)), scene, rect.Min)

	if hasAnimation {
		Invalidate(viewport.Bounds())
	}

	w.Screen().CopyRGBA(viewport, viewport.Bounds())

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

func init() {
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

	w.SetTitle("Stace Spation 2½") // Yes, that is the correct spelling.

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
			Disconnected()
			return
		}
	}
}

var Disconnected = func() {
	wde.Stop()
	os.Exit(0)
}

func Handle(msg *packet.Packet) {
	switch {
	case msg.Handshake != nil:
		world.Send <- *msg.Handshake

	case msg.Location != nil:
		world.Send <- *msg.Location

	case msg.Despawn != nil:
		world.Send <- *msg.Despawn

	case msg.MapOverride != nil:
		layout.SetChanges(msg.MapOverride.NewMap)

	case msg.MapChange != nil:
		for !layout.SetCoord(msg.MapChange.Coord, layout.GetCoord(msg.MapChange.Coord), msg.MapChange.Tile) {
		}

	default:
		log.Fatalf("unknown packet: %#v", msg)
	}
}

var (
	paintContexts = make(map[*actor.Actor]*PaintContext)
	paintLock     sync.RWMutex
)

type PaintContext struct {
	From    layout.Coord
	To      layout.Coord
	Changed time.Time
	Sprite  uint16
}
