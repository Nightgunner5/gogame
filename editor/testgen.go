package main

import (
	"fmt"
	"github.com/Nightgunner5/gogame/shared/layout"
	"os"
)

func writeMapTest() {
	minX, minY, maxX, maxY := 0, 0, 0, 0
	layout.AllTiles(func(c layout.Coord, t layout.MultiTile) {
		if t.Space() {
			return
		}
		if c.X < minX {
			minX = c.X
		}
		if c.Y < minY {
			minY = c.Y
		}
		if c.X > maxX {
			maxX = c.X
		}
		if c.Y > maxY {
			maxY = c.Y
		}
	})
	minX--
	minY--
	maxX++
	maxY++

	dx := maxX - minX
	dy := maxY - minY
	dx++
	dy++

	buf := make([]byte, dx*dy+dy-1)
	i := 0
	for y := minY; y <= maxY; y++ {
		if i != 0 {
			buf[i] = '\n'
			i++
		}
		for x := minX; x <= maxX; x++ {
			t := layout.Get(x, y)
			if t.Space() {
				buf[i] = ' '
			} else if t.Door() {
				if t.Passable() {
					buf[i] = 'd'
				} else {
					buf[i] = 'D'
				}
			} else if t.Passable() {
				buf[i] = '_'
			} else if t.BlocksVision() {
				buf[i] = 'W'
			} else {
				buf[i] = 'G'
			}
			i++
		}
	}

	f, err := os.Create("map_test.go")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(f, beforeTest)
	fmt.Fprint(f, "\n\ttop   = ", minY)
	fmt.Fprint(f, "\n\tleft  = ", minX)
	fmt.Fprint(f, "\n\tcheck = `", string(buf), "`\n")
	fmt.Fprint(f, afterTest)
	f.Close()
}

const (
	beforeTest = `package layout

import (
	"strings"
	"testing"
)

const (`

	afterTest = `)

func getCheck() [][]rune {
	lines := strings.Split(check, "\n")
	runes := make([][]rune, len(lines))
	for i := range lines {
		runes[i] = []rune(lines[i])
	}
	return runes
}

func TestCheck(t *testing.T) {
	c := getCheck()
	lengths := make([]int, len(c))
	for i, line := range c {
		lengths[i] = len(line)
		for _, r := range line {
			switch r {
			case ' ': // Space
			case 'W': // Wall
			case '_': // Floor
			case 'D': // Closed door
			case 'd': // Open door
			case 'G': // Window / gadget
			default:
				t.Errorf("Unexpected rune %c", r)
			}
		}
	}

	firstLength := lengths[0]
	for i, length := range lengths {
		if length != firstLength {
			t.Errorf("Length of line %d (%d) does not match length of line 0 (%d)", i, length, firstLength)
		}
	}
}

func TestLayout(t *testing.T) {
	assert := func(name string, tile MultiTile, x, y int, ok bool) {
		if ok {
			return
		}

		t.Errorf("Tile at (%d, %d) (%v) failed check %q", x, y, tile, name)
	}

	for y_, row := range getCheck() {
		for x_, r := range row {
			x, y := x_+left, y_+top
			tile := Get(x, y)
			switch r {
			case ' ':
				assert("space->space", tile, x, y, tile.Space())
				assert("space->passable", tile, x, y, tile.Passable())
				assert("space->!door", tile, x, y, !tile.Door())
				assert("space->!visblock", tile, x, y, !tile.BlocksVision())
			case 'W':
				assert("wall->!space", tile, x, y, !tile.Space())
				assert("wall->!passable", tile, x, y, !tile.Passable())
				assert("wall->!door", tile, x, y, !tile.Door())
				assert("wall->visblock", tile, x, y, tile.BlocksVision())
			case 'G':
				assert("window->!space", tile, x, y, !tile.Space())
				assert("window->!passable", tile, x, y, !tile.Passable())
				assert("window->!door", tile, x, y, !tile.Door())
				assert("window->!visblock", tile, x, y, !tile.BlocksVision())
			case '_':
				assert("floor->!space", tile, x, y, !tile.Space())
				assert("floor->passable", tile, x, y, tile.Passable())
				assert("floor->!door", tile, x, y, !tile.Door())
				assert("floor->!visblock", tile, x, y, !tile.BlocksVision())
			case 'D':
				assert("door->!space", tile, x, y, !tile.Space())
				assert("door->!passable", tile, x, y, !tile.Passable())
				assert("door->door", tile, x, y, tile.Door())
				assert("door->visblock", tile, x, y, tile.BlocksVision())
			case 'd':
				assert("door->!space", tile, x, y, !tile.Space())
				assert("door->passable", tile, x, y, tile.Passable())
				assert("door->door", tile, x, y, tile.Door())
				assert("door->!visblock", tile, x, y, !tile.BlocksVision())
			}
		}
	}
}`
)
