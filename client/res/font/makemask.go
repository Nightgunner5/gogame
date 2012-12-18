package main

import (
	"encoding/gob"
	"image"
	"image/draw"
	_ "image/png"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	runes = [][]rune{
		[]rune(`ABCDEFGHIJKLMN012345`),
		[]rune(`OPQRSTUVWXYZ.,-_6789`),
		[]rune(`abcdefghijklmn!@#$%^&`),
		[]rune(`opqrstuvwxyz*()[]+=/?<>'`),
	}

	coord = [3][3]int{
		{0, 0, 60},
		{0, 256, 31},
		{320, 256, 20},
	}
)

func main() {
	f, err := os.Open("text.png")
	check(err)

	src, _, err := image.Decode(f)
	check(err)
	f.Close()

	dest := image.NewAlpha(src.Bounds())
	draw.Draw(dest, dest.Bounds(), src, image.ZP, draw.Src)

	glyphs := make(map[rune][3]image.Rectangle)
	y := [3]int{coord[0][1], coord[1][1], coord[2][1]}
	for _, line := range runes {
		x := [3]int{coord[0][0], coord[1][0], coord[2][0]}
		for _, r := range line {
			println(string(r))
			var rect [3]image.Rectangle

			// Find left edges
			for i := range x {
				x[i]--
			leftEdge:
				for {
					x[i]++
					for j := y[i]; j < y[i]+coord[i][2]; j++ {
						if _, _, _, a := dest.At(x[i], j).RGBA(); a != 0 {
							break leftEdge
						}
					}
				}
			}

			for i := range rect {
				rect[i].Min.X, rect[i].Min.Y = x[i], y[i]
				rect[i].Max.Y = y[i] + coord[i][2]
			}

			// Find right edges
			for i := range x {
				x[i]--
			rightEdge:
				for {
					x[i]++
					for j := y[i]; j < y[i]+coord[i][2]; j++ {
						if _, _, _, a := dest.At(x[i], j).RGBA(); a != 0 {
							continue rightEdge
						}
					}
					break
				}
			}

			for i := range rect {
				rect[i].Max.X = x[i]
			}

			glyphs[r] = rect
		}
		for i := range y {
			y[i] += coord[i][2]
		}
	}

	f, err = os.Create("text.bin")
	check(err)

	encoder := gob.NewEncoder(f)
	err = encoder.Encode(dest)
	check(err)
	err = encoder.Encode(glyphs)
	check(err)
	f.Close()
}
