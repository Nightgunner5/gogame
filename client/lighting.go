package client

import (
	"github.com/Nightgunner5/gogame/shared/layout"
	"image"
	"image/color"
	"image/draw"
)

type lighting struct {
	lightmapVersion  uint64
	lightmap         map[layout.Coord]byte
	cachedImage      *image.RGBA
	cachedX, cachedY int
}

const LightShift = 5

func (l *lighting) Image(x, y int) *image.RGBA {
	if l.lightmapVersion != layout.Version() || l.lightmap == nil {
		l.lightmap = make(map[layout.Coord]byte)
		l.recalculateLightmap()
		l.cachedImage = nil
		l.lightmapVersion = layout.Version()
	}

	if l.cachedImage != nil && l.cachedX == x>>LightShift && l.cachedY == y>>LightShift {
		return l.cachedImage
	}

	l.cachedX, l.cachedY = x>>LightShift, y>>LightShift
	l.cachedImage = image.NewRGBA(image.Rect(0, 0, 2<<LightShift<<TileSize, 2<<LightShift<<TileSize))
	l.initializeImage()

	return l.cachedImage
}

func (l *lighting) Origin(x, y int) image.Point {
	return image.Pt((x&(1<<LightShift-1))<<TileSize,
		(y&(1<<LightShift-1))<<TileSize)
}

func (l *lighting) initializeImage() {
	for x := 0; x < 2<<LightShift; x++ {
		for y := 0; y < 2<<LightShift; y++ {
			draw.Draw(l.cachedImage, image.Rect(x<<TileSize, y<<TileSize, (x+1)<<TileSize, (y+1)<<TileSize),
				image.NewUniform(color.RGBA{A: 200 - l.lightmap[layout.Coord{x + l.cachedX<<LightShift, y + l.cachedY<<LightShift}]}), image.ZP, draw.Src)
		}
	}
}

func (l *lighting) recalculateLightmap() {
	layout.AllTiles(func(c layout.Coord, t layout.MultiTile) {
		if brightness := t.LightLevel(); brightness != 0 {
			l.spread(c, brightness)
		}
	})
}

func (l *lighting) spread(c layout.Coord, brightness byte) {
	if l.lightmap[c] < 200 {
		if 200-l.lightmap[c] < brightness {
			l.lightmap[c] = 200
		} else {
			l.lightmap[c] += brightness
		}
	}

	if layout.GetCoord(c).BlocksVision() {
		return
	}

	const LightLoss = 15
	if brightness <= LightLoss {
		return
	}
	brightness -= LightLoss
	l.spread(layout.Coord{c.X - 1, c.Y}, brightness)
	l.spread(layout.Coord{c.X + 1, c.Y}, brightness)
	l.spread(layout.Coord{c.X, c.Y - 1}, brightness)
	l.spread(layout.Coord{c.X, c.Y + 1}, brightness)
}
