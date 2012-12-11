package client

import (
	"bytes"
	"github.com/Nightgunner5/gogame/client/res"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"image"
	"image/draw"
	"image/png"
	"log"
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

	actors := world.GetHeld()
	count, paint := len(actors), make(chan PaintContext, len(actors))
	for _, actor := range actors {
		select {
		case actor.Send <- PaintRequest(paint):
		default:
			count--
		}
	}

	for x := rect.Min.X >> TileSize; x < (rect.Max.X-1)>>TileSize+1; x++ {
		for y := rect.Min.Y >> TileSize; y < (rect.Max.Y-1)>>TileSize+1; y++ {
			Tile(viewport, Terrain, uint16(layout.Get(x-xOffset, y-yOffset)), x, y)
		}
	}

	for i := 0; i < count; i++ {
		(<-paint).Paint(viewport, xOffset, yOffset)
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
			switch e.Key {
			case wde.KeyUpArrow:
				thePlayer.Send <- packet.Location{Coord: layout.Coord{0, -1}}
			case wde.KeyDownArrow:
				thePlayer.Send <- packet.Location{Coord: layout.Coord{0, 1}}
			case wde.KeyLeftArrow:
				thePlayer.Send <- packet.Location{Coord: layout.Coord{-1, 0}}
			case wde.KeyRightArrow:
				thePlayer.Send <- packet.Location{Coord: layout.Coord{1, 0}}
			}

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
	case msg.Chat != nil:
		log.Printf("<%s> %s", msg.Chat.User, msg.Chat.Message)
	case msg.Handshake != nil:
		world.Send <- *msg.Handshake
	case msg.Location != nil:
		world.Send <- *msg.Location
	default:
		log.Fatalf("unknown packet: %#v", msg)
	}
}
