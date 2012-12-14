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
	"image/color"
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

type lightStep struct {
	layout.Coord
	bright byte
}

func spreadLight(draw map[layout.Coord]byte, queue *[]lightStep, step lightStep) {
	if 255-draw[step.Coord] < step.bright {
		draw[step.Coord] = 255
	} else {
		draw[step.Coord] += step.bright
	}

	var lightLoss = step.bright>>1 + 1
	if step.bright <= lightLoss {
		return
	}
	step.bright -= lightLoss

	blocksVision := layout.GetCoord(step.Coord).BlocksVision()

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if !blocksVision || layout.Get(step.X+dx, step.Y+dy).Space() {
				*queue = append(*queue, lightStep{layout.Coord{step.X + dx, step.Y + dy}, step.bright})
			}
		}
	}

	for len(*queue) > 0 {
		step := (*queue)[0]
		*queue = (*queue)[1:]
		spreadLight(draw, queue, step)
	}
}

var (
	brightness        map[layout.Coord]byte
	lastBrightnessX   int
	lastBrightnessY   int
	brightnessVersion uint64
	brightnessReady   = make(chan struct{}, 1)
)

func Paint(w wde.Window, rect image.Rectangle) {
	viewport := w.Screen()

	xOffset, yOffset := GetTopLeft()

	minX, maxX := rect.Min.X>>TileSize, (rect.Max.X-1)>>TileSize+1
	minY, maxY := rect.Min.Y>>TileSize, (rect.Max.Y-1)>>TileSize+1

	if lastBrightnessX != xOffset || lastBrightnessY != yOffset {
		brightness = nil
	}

	if brightness == nil {
		go func() {
			brightness = make(map[layout.Coord]byte)
			spreadLight(brightness, new([]lightStep), lightStep{layout.Coord{ViewportWidth/2 - xOffset, ViewportHeight/2 - yOffset}, 100})
			brightnessReady <- struct{}{}
		}()
	} else {
		brightnessReady <- struct{}{}
	}

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			Tile(viewport, Terrain, uint16(layout.GetSpace(x-xOffset, y-yOffset)), x, y)
			for _, t := range layout.Get(x-xOffset, y-yOffset) {
				Tile(viewport, Terrain, uint16(t), x, y)
			}
		}
	}

	paintLock.Lock()
	for _, p := range paintContexts {
		x, y := p.Coord.X+xOffset, p.Coord.Y+yOffset
		if minX <= x && x <= maxX && minY <= y && y <= maxY {
			Tile(viewport, Actors, p.Sprite, x, y)
		}
	}
	paintLock.Unlock()

	<-brightnessReady
	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			bright := brightness[layout.Coord{x - xOffset, y - yOffset}]
			if bright < 255 && minX <= x && x <= maxX && minY <= y && y <= maxY {
				draw.Draw(viewport,
					image.Rect(x<<TileSize, y<<TileSize, (x+1)<<TileSize, (y+1)<<TileSize),
					image.NewUniform(color.RGBA{0, 0, 0, 255 - bright}), image.Pt(0, 0), draw.Over)
			}
		}
	}

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
	Coord  layout.Coord
	Sprite uint16
}
