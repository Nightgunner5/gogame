package layout

import (
	"strings"
	"testing"
)

const (
	top   = -5
	left  = -7
	check = `                          
      WWWWWWWW            
 WWWWWWW____WWW           
 WW___W______WW           
 W____D_______G           
 W____W_______G           
 W____G_______G           
 WW___W______WW           
 WWWDWWW____WWW           
      WWWWDWWW            
         G_G              
         G_G              
      GGGWDWGGWWWWWW      
      G_______WWGGGWWWWW  
      W________G___WW_WW  
      WW_______G___D___W  
      WW______WWGGWW__GW  
       W______WW__WW___W  
       W______D____W__GW  
       G______W____W___W  
       GGGGGGWW____WW_WW  
             WW____WWDWWW 
              W________WW 
              W_________W 
              W_________W 
              WW_______WW 
              WWWWWWWWWWW 
                          `
)

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
}
