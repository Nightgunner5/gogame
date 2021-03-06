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
	ToolboxHeight = 16
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

	center := layout.Coord{ViewportWidth/2 - xOffset, ViewportHeight/2 - yOffset}
	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			res.Tile(space, res.Terrain, uint16(layout.GetSpace(x-xOffset, y-yOffset)), x, y)
			draw.Draw(scene, image.Rect(x<<res.TileSize, y<<res.TileSize, (x+1)<<res.TileSize, (y+1)<<res.TileSize), image.Transparent, image.ZP, draw.Src)
			if WireView() {
				tile := layout.Get(x-xOffset, y-yOffset)
				for i := len(tile) - 1; i > 0; i-- {
					res.Tile(scene, res.Terrain, uint16(tile[i]), x, y)
				}
			} else {
				for _, t := range layout.Get(x-xOffset, y-yOffset) {
					res.Tile(scene, res.Terrain, uint16(t), x, y)
				}
			}
			if VisibilityOn() && !layout.Visible(center, layout.Coord{x - xOffset, y - yOffset}) {
				res.Tile(scene, image.Black, 0, x, y)
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
		draw.DrawMask(viewport, rect, image.Black, image.ZP, light.Image(-xOffset, -yOffset), rect.Min.Add(light.Origin(-xOffset, -yOffset)), draw.Over)
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
			res.Tile(toolbox.Screen(), res.Terrain, uint16(i+j*ToolboxWidth), i, j)
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
					mouseTile = layout.Tile((e.Where.X >> res.TileSize) + (e.Where.Y>>res.TileSize)*ToolboxWidth)
					mouseValid = true
				case wde.RightButton:
					mouseValid = false
				}
				mouseLock.Unlock()
			case wde.ResizeEvent:
				toolbox.SetSize(ToolboxWidth<<res.TileSize, ToolboxHeight<<res.TileSize)
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
					for {
						old := layout.GetCoord(mouseCoord)
						var changed layout.MultiTile

						if WireView() {
							old_ := old
							if len(old_) == 0 {
								old_ = layout.MultiTile{layout.Plating}
							}
							changed = append(changed, old_[0])
							changed = append(changed, mouseTile)
							changed = append(changed, old_[1:]...)
						} else {
							changed = append(changed, old...)

							changed = append(changed, mouseTile)
						}

						if layout.SetCoord(mouseCoord, old, changed) {
							break
						}
					}
				} else {
					for {
						tile := layout.GetCoord(mouseCoord)
						if len(tile) == 0 {
							break
						}
						t, removed := tile[len(tile)-1], tile[:len(tile)-1]
						if WireView() {
							if len(tile) == 1 {
								break
							}
							t = tile[1]
							removed = append(layout.MultiTile{tile[0]}, tile[2:]...)
						}
						if layout.SetCoord(mouseCoord, tile, removed) {
							mouseTile = t
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
			w.SetSize(ViewportWidth<<res.TileSize, ViewportHeight<<res.TileSize)

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
