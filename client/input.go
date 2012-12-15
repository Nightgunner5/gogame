package client

import (
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/skelterjohn/go.wde"
)

var (
	keyDown = make(map[string]func(map[string]bool))
	keyUp   = make(map[string]func(map[string]bool))
)

func init() {
	keyDown[wde.KeyUpArrow] = updateMotion
	keyDown[wde.KeyDownArrow] = updateMotion
	keyDown[wde.KeyLeftArrow] = updateMotion
	keyDown[wde.KeyRightArrow] = updateMotion

	keyUp[wde.KeyUpArrow] = updateMotion
	keyUp[wde.KeyDownArrow] = updateMotion
	keyUp[wde.KeyLeftArrow] = updateMotion
	keyUp[wde.KeyRightArrow] = updateMotion
}

func updateMotion(keys map[string]bool) {
	dx, dy := 0, 0

	if keys[wde.KeyUpArrow] {
		,
	}
	if keys[wde.KeyDownArrow] {
		o
	}
	if keys[wde.KeyLeftArrow] {
		a
	}
	if keys[wde.KeyRightArrow] {
		e
	}

	Network <- &packet.Packet{
		Location: &packet.Location{
			Coord: layout.Coord{dx, dy},
		},
	}
}
