package main

import (
	"bytes"
	"github.com/Nightgunner5/gogame/client/res"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"image"
	"image/draw"
	"image/png"
)

var (
	Terrain image.Image
	Actors  image.Image
)

func init() {
	var err error

	Terrain, err = png.Decode(bytes.NewReader(res.Terrain))
	if err != nil {
		panic(err)
	}

	Actors, err = png.Decode(bytes.NewReader(res.Actors))
	if err != nil {
		panic(err)
	}
}

func main() {
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

func tileCoord(index int) (p image.Point) {
	return image.Pt((index&TileMask)<<TileSize,
		(index>>TileShift&TileMask)<<TileSize)
}

func Tile(viewport draw.Image, base image.Image, index, x, y int) {
	x, y = x<<TileSize, y<<TileSize
	draw.Draw(viewport, image.Rect(x, y, x+1<<TileSize, y+1<<TileSize), base, tileCoord(index), draw.Over)
}

func getTile(x, y int) int {
	if x == 0 || y == 0 || x == ViewportWidth-1 || y == ViewportHeight-1 {
		return 16
	}
	return (x ^ y) & 1
}

func Paint(w wde.Window, rect image.Rectangle) {
	viewport := w.Screen()

	for x := rect.Min.X >> TileSize; x < (rect.Max.X-1)>>TileSize+1; x++ {
		for y := rect.Min.Y >> TileSize; y < (rect.Max.Y-1)>>TileSize+1; y++ {
			Tile(viewport, Terrain, getTile(x, y), x, y)
		}
	}
	actors := world.GetHeld()
	count, paint := len(actors), make(chan PaintContext, len(actors))
	for _, actor := range actors {
		actor.Send <- PaintRequest(paint)
	}
	for i := 0; i < count; i++ {
		(<-paint).Paint(viewport)
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

func UI() {
	defer wde.Stop()

	w, err := wde.NewWindow(ViewportWidth<<TileSize, ViewportHeight<<TileSize)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	w.SetTitle("GoGame")

	w.Show()

	go func() {
		for rect := range shouldPaint {
			Paint(w, rect)
		}
	}()

	Invalidate(w.Screen().Bounds())

	for event := range w.EventChan() {
		switch e := event.(type) {
		case wde.MouseMovedEvent:
		case wde.MouseDownEvent:
		case wde.MouseUpEvent:
		case wde.MouseDraggedEvent:
		case wde.MouseEnteredEvent:
		case wde.MouseExitedEvent:
		case wde.KeyDownEvent:
		case wde.KeyUpEvent:
		case wde.KeyTypedEvent:
		case wde.ResizeEvent:
			w.SetSize(ViewportWidth<<TileSize, ViewportHeight<<TileSize)
			Invalidate(w.Screen().Bounds())
		case wde.CloseEvent:
			_ = e
			return
		}
	}
}
