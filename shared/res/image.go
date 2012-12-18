package res

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
)

const (
	TileSize  = 5 // 2**5 pixels by 2**5 pixels
	TileShift = 5 // 2**5 tiles by 2**5 tiles
	TileMask  = 1<<TileShift - 1
)

var (
	Terrain image.Image
	Actors  image.Image
)

func init() {
	var err error

	Terrain, err = png.Decode(bytes.NewReader(TerrainPng))
	if err != nil {
		panic(err)
	}
	TerrainPng = nil

	Actors, err = png.Decode(bytes.NewReader(ActorsPng))
	if err != nil {
		panic(err)
	}
	ActorsPng = nil
}

func tileCoord(index uint16) (p image.Point) {
	return image.Pt(int((index&TileMask)<<TileSize),
		int((index>>TileShift&TileMask)<<TileSize))
}

func Tile(viewport draw.Image, base image.Image, index uint16, x, y int) {
	x, y = x<<TileSize, y<<TileSize
	draw.Draw(viewport, image.Rect(x, y, x+1<<TileSize, y+1<<TileSize), base, tileCoord(index), draw.Over)
}

func TileFloat(viewport draw.Image, base image.Image, index uint16, x1, y1, x2, y2 int, interp float32) {
	x1, y1 = x1<<TileSize, y1<<TileSize
	x2, y2 = x2<<TileSize, y2<<TileSize
	x, y := x1+int(float32(x2-x1)*interp), y1+int(float32(y2-y1)*interp)
	draw.Draw(viewport, image.Rect(x, y, x+1<<TileSize, y+1<<TileSize), base, tileCoord(index), draw.Over)
}
