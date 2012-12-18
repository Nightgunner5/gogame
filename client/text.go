package client

import (
	"bytes"
	"encoding/gob"
	"github.com/Nightgunner5/gogame/client/res"
	"image"
	"image/color"
	"image/draw"
)

const (
	FontSmall  = 2
	FontMedium = 1
	FontLarge  = 0
)

var (
	Text      *image.Alpha
	glyphRect map[rune][3]image.Rectangle
)

func init() {
	dec := gob.NewDecoder(bytes.NewReader(res.TextBin))
	if err := dec.Decode(&Text); err != nil {
		panic(err)
	}
	if err := dec.Decode(&glyphRect); err != nil {
		panic(err)
	}
}

func drawString(dst draw.Image, s string, color color.Color, size, x, y int) int {
	c := image.NewUniform(color)
	for _, r := range s {
		x = drawRune(dst, r, c, size, x, y)
	}
	return x
}

func drawRune(dst draw.Image, r rune, color *image.Uniform, size, x, y int) int {
	if r == ' ' {
		return x + [3]int{24, 12, 8}[size]
	}
	rect := glyphRect[r][size]
	draw.DrawMask(dst, rect.Sub(rect.Min).Add(image.Pt(x, y)), color, image.ZP, Text, rect.Min, draw.Over)
	return x + rect.Dx()
}
