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

	keyDown[wde.KeyW] = updateMotion
	keyDown[wde.KeyS] = updateMotion
	keyDown[wde.KeyA] = updateMotion
	keyDown[wde.KeyD] = updateMotion

	keyUp[wde.KeyUpArrow] = updateMotion
	keyUp[wde.KeyDownArrow] = updateMotion
	keyUp[wde.KeyLeftArrow] = updateMotion
	keyUp[wde.KeyRightArrow] = updateMotion

	keyUp[wde.KeyW] = updateMotion
	keyUp[wde.KeyS] = updateMotion
	keyUp[wde.KeyA] = updateMotion
	keyUp[wde.KeyD] = updateMotion
}

func updateMotion(keys map[string]bool) {
	dx, dy := 0, 0

	if keys[wde.KeyUpArrow] || keys[wde.KeyW] || keys[wde.KeyComma] {
		dy--
	}
	if keys[wde.KeyDownArrow] || keys[wde.KeyS] || keys[wde.KeyO] {
		dy++
	}
	if keys[wde.KeyLeftArrow] || keys[wde.KeyA] {
		dx--
	}
	if keys[wde.KeyRightArrow] || keys[wde.KeyD] || keys[wde.KeyE] {
		dx++
	}

	Network <- &packet.Packet{
		Location: &packet.Location{
			Coord: layout.Coord{dx, dy},
		},
	}
}
