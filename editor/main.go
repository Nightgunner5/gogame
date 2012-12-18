package main

import (
	"bytes"
	"github.com/Nightgunner5/gogame/client/res"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"image"
	"image/draw"
	"image/png"
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

func main() {
	layout.OnChange = func(c layout.Coord, t layout.MultiTile) {
		xOffset, yOffset := GetTopLeft()
		x, y := c.X+xOffset, c.Y+yOffset
		Invalidate(image.Rect(x<<TileSize, y<<TileSize, (x+1)<<TileSize, (y+1)<<TileSize))
	}

	go UI()

	wde.Run()
}

var (
	mouseValid bool
	mouseTile  layout.Tile
	mouseCoord layout.Coord
	mouseLock  sync.Mutex
)

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

var (
	viewport = image.NewRGBA(image.Rect(0, 0, ViewportWidth<<TileSize, ViewportHeight<<TileSize))
	space    = image.NewRGBA(viewport.Bounds())
	scene    = image.NewRGBA(viewport.Bounds())
	light    = new(lighting)
)

func Paint(w wde.Window, rect image.Rectangle) {
	xOffset, yOffset := GetTopLeft()

	minX, maxX := rect.Min.X>>TileSize, (rect.Max.X-1)>>TileSize+1
	minY, maxY := rect.Min.Y>>TileSize, (rect.Max.Y-1)>>TileSize+1

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			Tile(space, Terrain, uint16(layout.GetSpace(x-xOffset, y-yOffset)), x, y)
			draw.Draw(scene, image.Rect(x<<TileSize, y<<TileSize, (x+1)<<TileSize, (y+1)<<TileSize), image.Transparent, image.ZP, draw.Src)
			for _, t := range layout.Get(x-xOffset, y-yOffset) {
				Tile(scene, Terrain, uint16(t), x, y)
			}
		}
	}

	mouseLock.Lock()
	if mouseValid {
		Tile(scene, Terrain, uint16(mouseTile), mouseCoord.X+xOffset, mouseCoord.Y+yOffset)
	}
	mouseLock.Unlock()

	draw.Draw(viewport, rect, space, rect.Min, draw.Src)
	draw.Draw(viewport, rect, scene, rect.Min, draw.Over)
	if LightsOn() {
		drawLightOverlay(viewport, rect, light.Image(-xOffset, -yOffset), rect.Min.Add(light.Origin(-xOffset, -yOffset)), scene, rect.Min)
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

	w.SetTitle("Stace Spation 2½ Map Editor") // Yes, that is the correct spelling.

	toolbox, err := wde.NewWindow(8<<TileSize, 5<<TileSize)
	if err != nil {
		panic(err)
	}
	defer toolbox.Close()

	toolbox.SetTitle("Tiles")

	draw.Draw(toolbox.Screen(), toolbox.Screen().Bounds(), image.White, image.ZP, draw.Src)
	for i := 0; i < 8; i++ {
		Tile(toolbox.Screen(), Terrain, uint16(i), i, 0)
		Tile(toolbox.Screen(), Terrain, uint16(i+8), i, 1)
		Tile(toolbox.Screen(), Terrain, uint16(i+16), i, 2)
		Tile(toolbox.Screen(), Terrain, uint16(i+24), i, 3)
		Tile(toolbox.Screen(), Terrain, uint16(i+32), i, 4)
	}
	toolbox.FlushImage(toolbox.Screen().Bounds())

	w.Show()
	toolbox.Show()

	go paintHandler(w)

	Invalidate(w.Screen().Bounds())

	keys := make(map[string]bool)

	go func() {
		for event := range toolbox.EventChan() {
			switch e := event.(type) {
			case wde.MouseDownEvent:
				mouseLock.Lock()
				switch e.Which {
				case wde.LeftButton:
					mouseTile = layout.Tile((e.Where.X >> TileSize) + (e.Where.Y>>TileSize)*8)
					mouseValid = true
				case wde.RightButton:
					mouseValid = false
				}
				mouseLock.Unlock()
			}
		}
	}()

	for event := range w.EventChan() {
		switch e := event.(type) {
		case wde.MouseMovedEvent:
			xOffset, yOffset := GetTopLeft()

			mouseLock.Lock()
			mouseCoord.X, mouseCoord.Y = e.Where.X>>TileSize-xOffset, e.Where.Y>>TileSize-yOffset
			mouseLock.Unlock()

			Invalidate(viewport.Bounds())

		case wde.MouseDownEvent:
			switch e.Which {
			case wde.LeftButton:
				mouseLock.Lock()
				if mouseValid {
					for !layout.SetCoord(mouseCoord, layout.GetCoord(mouseCoord), append(append(layout.MultiTile{}, layout.GetCoord(mouseCoord)...), mouseTile)) {
					}
				} else {
					for {
						tile := layout.GetCoord(mouseCoord)
						if len(tile) == 0 {
							break
						}
						if layout.SetCoord(mouseCoord, tile, tile[:len(tile)-1]) {
							mouseTile = tile[len(tile)-1]
							mouseValid = true
							break
						}
					}
				}
				mouseLock.Unlock()

			case wde.RightButton:
				mouseLock.Lock()
				mouseValid = false
				mouseLock.Unlock()
			}

			Invalidate(viewport.Bounds())

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
	writeMapLayout()
	writeMapTest()
	wde.Stop()
	os.Exit(0)
}
