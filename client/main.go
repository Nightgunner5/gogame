package client

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/lighting"
	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/Nightgunner5/gogame/shared/res"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func Main() {
	layout.OnChange = func(c layout.Coord, t layout.MultiTile) {
		if t.Door() {
			Invalidate(image.Rect(0, 0, ViewportWidth<<res.TileSize, ViewportHeight<<res.TileSize))
		} else {
			xOffset, yOffset := GetTopLeft()
			x, y := c.X+xOffset, c.Y+yOffset
			Invalidate(image.Rect(x<<res.TileSize, y<<res.TileSize, (x+1)<<res.TileSize, (y+1)<<res.TileSize))
		}
	}

	go UI()

	wde.Run()
}

var (
	mouseTile       layout.Coord
	mouseTileString string
	mouseTileLock   sync.Mutex
)

const (
	ViewportWidth  = 20
	ViewportHeight = 15
)

var (
	viewport = image.NewRGBA(image.Rect(0, 0, ViewportWidth<<res.TileSize, ViewportHeight<<res.TileSize))
	space    = image.NewRGBA(viewport.Bounds())
	scene    = image.NewRGBA(viewport.Bounds())
	light    = new(lighting.Lighting)
)

func Paint(w wde.Window, rect image.Rectangle) {
	xOffset, yOffset := GetTopLeft()
	center := layout.Coord{ViewportWidth/2 - xOffset, ViewportHeight/2 - yOffset}

	minX, maxX := rect.Min.X>>res.TileSize, (rect.Max.X-1)>>res.TileSize+1
	minY, maxY := rect.Min.Y>>res.TileSize, (rect.Max.Y-1)>>res.TileSize+1

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			if layout.Visible(center, layout.Coord{x - xOffset, y - yOffset}) {
				res.Tile(space, res.Terrain, uint16(layout.GetSpace(x-xOffset, y-yOffset)), x, y)
				draw.Draw(scene, image.Rect(x<<res.TileSize, y<<res.TileSize, (x+1)<<res.TileSize, (y+1)<<res.TileSize), image.Transparent, image.ZP, draw.Src)
				for _, t := range layout.Get(x-xOffset, y-yOffset) {
					if !t.NoClient() {
						res.Tile(scene, res.Terrain, uint16(t), x, y)
					}
				}
			} else {
				res.Tile(scene, image.Black, 0, x, y)
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
				if interp >= 1 {
					res.Tile(scene, res.Actors, p.Sprite, x2, y2)
				} else {
					hasAnimation = true
					res.TileFloat(scene, res.Actors, p.Sprite, x1, y1, x2, y2, interp)
				}
			}
		}
	}
	paintLock.Unlock()

	draw.Draw(viewport, rect, space, rect.Min, draw.Src)
	draw.Draw(viewport, rect, scene, rect.Min, draw.Over)
	draw.DrawMask(viewport, rect, image.Black, image.ZP, light.Image(-xOffset, -yOffset), rect.Min.Add(light.Origin(-xOffset, -yOffset)), draw.Over)

	mouseTileLock.Lock()
	res.DrawString(viewport, mouseTileString, color.White, res.FontSmall, 1, 1)
	mouseTileLock.Unlock()

	if hasAnimation {
		Invalidate(rect)
	}

	w.Screen().CopyRGBA(viewport.SubImage(rect).(*image.RGBA), rect)

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
		time.Sleep(20 * time.Millisecond)
	}
}

func UI() {
	defer wde.Stop()

	w, err := wde.NewWindow(ViewportWidth<<res.TileSize, ViewportHeight<<res.TileSize)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	w.SetTitle("Stace Spation 2½") // Yes, that is the correct spelling.

	w.LockSize(true)

	w.Show()

	go paintHandler(w)

	Invalidate(w.Screen().Bounds())

	keys := make(map[string]bool)

	for event := range w.EventChan() {
		switch e := event.(type) {
		case wde.MouseMovedEvent:
			xOffset, yOffset := GetTopLeft()

			mouseTile.X, mouseTile.Y = e.Where.X>>res.TileSize-xOffset, e.Where.Y>>res.TileSize-yOffset
			if layout.Visible(layout.Coord{ViewportWidth/2 - xOffset, ViewportHeight/2 - yOffset}, mouseTile) {
				tooltip := strings.Join(layout.GetCoord(mouseTile).Describe(), ", ")

				mouseTileLock.Lock()
				mouseTileString = tooltip
				mouseTileLock.Unlock()
			} else {
				mouseTileLock.Lock()
				mouseTileString = ""
				mouseTileLock.Unlock()
			}
			Invalidate(image.Rect(0, 0, ViewportWidth<<res.TileSize, 1<<res.TileSize))

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
			w.SetSize(ViewportWidth<<res.TileSize, ViewportHeight<<res.TileSize)

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
