package lighting

import (
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/res"
	"image"
	"image/color"
	"image/draw"
)

type Lighting struct {
	lightmapVersion  uint64
	lightmap         map[layout.Coord]byte
	cachedImage      *image.RGBA
	cachedInvalid    bool
	cachedX, cachedY int
}

const (
	lightShift = 5
)

func (l *Lighting) Image(x, y int) *image.RGBA {
	if l.lightmapVersion != layout.Version() || l.lightmap == nil {
		l.lightmapVersion = layout.Version()
		l.lightmap = make(map[layout.Coord]byte)
		l.recalculateLightmap()
		l.cachedInvalid = true
	}

	if !l.cachedInvalid && l.cachedImage != nil && l.cachedX == x>>lightShift && l.cachedY == y>>lightShift {
		return l.cachedImage
	}

	if l.cachedImage == nil {
		l.cachedImage = image.NewRGBA(image.Rect(0, 0, 2<<lightShift<<res.TileSize, 2<<lightShift<<res.TileSize))
	}

	l.cachedX, l.cachedY = x>>lightShift, y>>lightShift
	l.initializeImage()

	return l.cachedImage
}

func (l *Lighting) Origin(x, y int) image.Point {
	return image.Pt((x&(1<<lightShift-1))<<res.TileSize,
		(y&(1<<lightShift-1))<<res.TileSize)
}

func (l *Lighting) initializeImage() {
	for x := 0; x < 2<<lightShift; x++ {
		for y := 0; y < 2<<lightShift; y++ {
			draw.Draw(l.cachedImage, image.Rect(x<<res.TileSize, y<<res.TileSize, (x+1)<<res.TileSize, (y+1)<<res.TileSize),
				image.NewUniform(color.RGBA{A: 250 - l.lightmap[layout.Coord{x + l.cachedX<<lightShift, y + l.cachedY<<lightShift}]}), image.ZP, draw.Src)
		}
	}
}

func (l *Lighting) recalculateLightmap() {
	layout.AllTiles(func(c layout.Coord, t layout.MultiTile) {
		if brightness := t.LightLevel(); brightness != 0 {
			l.spread(c, brightness)
		}
	})
}

func (l *Lighting) spread(c layout.Coord, brightness byte) {
	if l.lightmap[c] < 250 {
		if 250-l.lightmap[c] < brightness {
			l.lightmap[c] = 250
		} else {
			l.lightmap[c] += brightness
		}
	}

	if layout.GetCoord(c).BlocksVision() {
		return
	}

	var lightLoss = brightness>>2 + 3
	if brightness <= lightLoss {
		return
	}
	brightness -= lightLoss
	l.spread(layout.Coord{c.X - 1, c.Y}, brightness)
	l.spread(layout.Coord{c.X + 1, c.Y}, brightness)
	l.spread(layout.Coord{c.X, c.Y - 1}, brightness)
	l.spread(layout.Coord{c.X, c.Y + 1}, brightness)
}
