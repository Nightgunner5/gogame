package lighting

import (
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/power"
	"github.com/Nightgunner5/gogame/shared/res"
	"image"
	"math"
)

type Lighting struct {
	lightmapVersion  uint64
	lightmap         map[layout.Coord]uint8
	cachedImage      *image.Alpha
	cachedInvalid    bool
	cachedX, cachedY int
}

const (
	lightShift = 5
)

func (l *Lighting) Image(x, y int) *image.Alpha {
	if l.lightmapVersion != layout.Version() || l.lightmap == nil {
		l.lightmapVersion = layout.Version()
		l.lightmap = make(map[layout.Coord]uint8)
		l.recalculateLightmap()
		l.cachedInvalid = true
	}

	if !l.cachedInvalid && l.cachedImage != nil && l.cachedX == x>>lightShift && l.cachedY == y>>lightShift {
		return l.cachedImage
	}

	if l.cachedImage == nil {
		l.cachedImage = image.NewAlpha(image.Rect(0, 0, 2<<lightShift<<res.TileSize, 2<<lightShift<<res.TileSize))
	}

	l.cachedInvalid = false
	l.cachedX, l.cachedY = x>>lightShift, y>>lightShift
	l.initializeImage()

	return l.cachedImage
}

func (l *Lighting) Origin(x, y int) image.Point {
	return image.Pt((x&(1<<lightShift-1))<<res.TileSize,
		(y&(1<<lightShift-1))<<res.TileSize)
}

func (l *Lighting) lightAt(x, y int) uint16 {
	return uint16(220 - l.lightmap[layout.Coord{x, y}])
}

func (l *Lighting) initializeImage() {
	const (
		ws  = 1 + lightShift + res.TileSize
		w   = 1 << ws
		ts  = 1 << res.TileSize
		tsh = res.TileSize - 1
		tsm = 1<<tsh - 1
	)
	cx := l.cachedX << lightShift
	cy := l.cachedY << lightShift

	for y := 0; y < w; y += ts {
		for x := 0; x < w; x += ts {
			var input [3][3]uint16
			tx, ty := x>>res.TileSize+cx, y>>res.TileSize+cy
			input[1][1] = l.lightAt(tx, ty)

			input[0][1] = l.lightAt(tx-1, ty)
			input[2][1] = l.lightAt(tx+1, ty)
			input[1][0] = l.lightAt(tx, ty-1)
			input[1][2] = l.lightAt(tx, ty+1)

			input[0][0] = input[0][1] + input[1][0] + l.lightAt(tx-1, ty-1)
			input[2][0] = input[2][1] + input[1][0] + l.lightAt(tx+1, ty-1)
			input[0][2] = input[0][1] + input[1][2] + l.lightAt(tx-1, ty+1)
			input[2][2] = input[2][1] + input[1][2] + l.lightAt(tx+1, ty+1)

			input[0][1] = (input[0][1] + input[1][1]) / 2
			input[2][1] = (input[2][1] + input[1][1]) / 2
			input[1][0] = (input[1][0] + input[1][1]) / 2
			input[1][2] = (input[1][2] + input[1][1]) / 2

			input[0][0] = (input[0][0] + input[1][1]) / 4
			input[2][0] = (input[2][0] + input[1][1]) / 4
			input[0][2] = (input[0][2] + input[1][1]) / 4
			input[2][2] = (input[2][2] + input[1][1]) / 4

			for i := uint16(0); i < ts; i++ {
				for j := uint16(0); j < ts; j++ {
					a11 := input[i>>tsh][j>>tsh]
					a12 := input[i>>tsh][j>>tsh+1]
					a21 := input[i>>tsh+1][j>>tsh]
					a22 := input[i>>tsh+1][j>>tsh+1]

					b1 := (a11*(^i&tsm) + a21*(i&tsm)) >> tsh
					b2 := (a12*(^i&tsm) + a22*(i&tsm)) >> tsh

					c := (b1*(^j&tsm) + b2*(j&tsm)) >> tsh

					l.cachedImage.Pix[x|int(i)|(y|int(j))<<ws] = uint8(c)
				}
			}
		}
	}
}

func (l *Lighting) recalculateLightmap() {
	layout.AllTiles(func(c layout.Coord, t layout.MultiTile) {
		if brightness := t.LightLevel(); brightness != 0 && power.Powered(c.X, c.Y) {
			l.spread(c, brightness)
		}
	})
}

var brightnessMap [256]uint8

func init() {
	for i := range brightnessMap {
		brightnessMap[i] = uint8(math.Sqrt(float64(i))*2.5) + 2
	}
}

func (l *Lighting) spread(c layout.Coord, brightness uint8) {
	if l.lightmap[c] < 220 {
		if 220-l.lightmap[c] < brightness {
			l.lightmap[c] = 220
		} else {
			l.lightmap[c] += brightness
		}
	}

	if layout.GetCoord(c).BlocksVision() {
		return
	}

	var lightLoss = brightnessMap[brightness]
	if brightness <= lightLoss {
		return
	}
	brightness -= lightLoss
	l.spread(layout.Coord{c.X - 1, c.Y}, brightness)
	l.spread(layout.Coord{c.X + 1, c.Y}, brightness)
	l.spread(layout.Coord{c.X, c.Y - 1}, brightness)
	l.spread(layout.Coord{c.X, c.Y + 1}, brightness)
}
