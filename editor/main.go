package main

import (
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/lighting"
	"github.com/Nightgunner5/gogame/shared/res"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"image"
	"image/draw"
	"os"
	"sync"
	"time"
)

func main() {
	layout.OnChange = func(c layout.Coord, t layout.MultiTile) {
		xOffset, yOffset := GetTopLeft()
		x, y := c.X+xOffset, c.Y+yOffset
		Invalidate(image.Rect(x<<res.TileSize, y<<res.TileSize, (x+1)<<res.TileSize, (y+1)<<res.TileSize))
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
	ViewportWidth  = 20
	ViewportHeight = 15

	ToolboxWidth  = 8
	ToolboxHeight = 5
)

var (
	viewport = image.NewRGBA(image.Rect(0, 0, ViewportWidth<<res.TileSize, ViewportHeight<<res.TileSize))
	space    = image.NewRGBA(viewport.Bounds())
	scene    = image.NewRGBA(viewport.Bounds())
	light    = new(lighting.Lighting)
)

func Paint(w wde.Window, rect image.Rectangle) {
	xOffset, yOffset := GetTopLeft()

	minX, maxX := rect.Min.X>>res.TileSize, (rect.Max.X-1)>>res.TileSize+1
	minY, maxY := rect.Min.Y>>res.TileSize, (rect.Max.Y-1)>>res.TileSize+1

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			res.Tile(space, res.Terrain, uint16(layout.GetSpace(x-xOffset, y-yOffset)), x, y)
			draw.Draw(scene, image.Rect(x<<res.TileSize, y<<res.TileSize, (x+1)<<res.TileSize, (y+1)<<res.TileSize), image.Transparent, image.ZP, draw.Src)
			for _, t := range layout.Get(x-xOffset, y-yOffset) {
				res.Tile(scene, res.Terrain, uint16(t), x, y)
			}
		}
	}

	mouseLock.Lock()
	if mouseValid {
		res.Tile(scene, res.Terrain, uint16(mouseTile), mouseCoord.X+xOffset, mouseCoord.Y+yOffset)
	}
	mouseLock.Unlock()

	draw.Draw(viewport, rect, space, rect.Min, draw.Src)
	draw.Draw(viewport, rect, scene, rect.Min, draw.Over)
	if LightsOn() {
		lighting.DrawLightOverlay(viewport, rect, light.Image(-xOffset, -yOffset), rect.Min.Add(light.Origin(-xOffset, -yOffset)), scene, rect.Min)
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

	w, err := wde.NewWindow(ViewportWidth<<res.TileSize, ViewportHeight<<res.TileSize)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	w.SetTitle("Stace Spation 2½ Map Editor") // Yes, that is the correct spelling.

	toolbox, err := wde.NewWindow(ToolboxWidth<<res.TileSize, ToolboxHeight<<res.TileSize)
	if err != nil {
		panic(err)
	}
	defer toolbox.Close()

	toolbox.SetTitle("Tiles")

	draw.Draw(toolbox.Screen(), toolbox.Screen().Bounds(), image.White, image.ZP, draw.Src)
	for i := 0; i < ToolboxWidth; i++ {
		for j := 0; j < ToolboxHeight; j++ {
			res.Tile(toolbox.Screen(), res.Terrain, uint16(i+j<<3), i, j)
		}
	}
	toolbox.FlushImage(toolbox.Screen().Bounds())

	toolbox.LockSize(true)
	w.LockSize(true)

	toolbox.Show()
	w.Show()

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
					mouseTile = layout.Tile((e.Where.X >> res.TileSize) + (e.Where.Y>>res.TileSize)*8)
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
			mouseCoord.X, mouseCoord.Y = e.Where.X>>res.TileSize-xOffset, e.Where.Y>>res.TileSize-yOffset
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
