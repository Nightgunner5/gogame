package main

import (
	"github.com/skelterjohn/go.wde"
	"sync/atomic"
	"time"
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

	keyDown[wde.KeyL] = toggleLights
	keyDown[wde.KeyW] = toggleWires
	keyDown[wde.KeyV] = toggleVisibility
}

func updateMotion(keys map[string]bool) {
	dx, dy := 0, 0

	if keys[wde.KeyUpArrow] {
		dy--
	}
	if keys[wde.KeyDownArrow] {
		dy++
	}
	if keys[wde.KeyLeftArrow] {
		dx--
	}
	if keys[wde.KeyRightArrow] {
		dx++
	}

	atomic.StoreInt64(&topLeftDx, int64(dx))
	atomic.StoreInt64(&topLeftDy, int64(dy))
	select {
	case updateLocationImmediately <- struct{}{}:
	default:
	}
}

func toggleLights(keys map[string]bool) {
	for {
		old := atomic.LoadUint32(&viewFlags)
		if atomic.CompareAndSwapUint32(&viewFlags, old, old^1) {
			break
		}
	}
	Invalidate(viewport.Bounds())
}

func toggleWires(keys map[string]bool) {
	for {
		old := atomic.LoadUint32(&viewFlags)
		if atomic.CompareAndSwapUint32(&viewFlags, old, old^2) {
			break
		}
	}
	Invalidate(viewport.Bounds())
}

func toggleVisibility(keys map[string]bool) {
	for {
		old := atomic.LoadUint32(&viewFlags)
		if atomic.CompareAndSwapUint32(&viewFlags, old, old^4) {
			break
		}
	}
	Invalidate(viewport.Bounds())
}

var (
	topLeftX, topLeftY        int64 = ViewportWidth / 2, ViewportHeight / 2
	topLeftDx, topLeftDy      int64 = 0, 0
	updateLocationImmediately       = make(chan struct{}, 1)

	viewFlags uint32
)

func LightsOn() bool {
	return atomic.LoadUint32(&viewFlags)&1 != 0
}

func WireView() bool {
	return atomic.LoadUint32(&viewFlags)&2 != 0
}

func VisibilityOn() bool {
	return atomic.LoadUint32(&viewFlags)&4 != 0
}

func GetTopLeft() (x, y int) {
	x = int(atomic.LoadInt64(&topLeftX))
	y = int(atomic.LoadInt64(&topLeftY))
	return
}

func init() {
	go func() {
		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ticker:
			case <-updateLocationImmediately:
			}

			atomic.AddInt64(&topLeftX, -atomic.LoadInt64(&topLeftDx))
			atomic.AddInt64(&topLeftY, -atomic.LoadInt64(&topLeftDy))
			Invalidate(viewport.Bounds())
		}
	}()
}
