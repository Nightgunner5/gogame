package layout

import (
	"sync"
)

var (
	visLock sync.RWMutex
	visible map[[2]Coord]bool
)

func visInvalidateAll() {
	visLock.Lock()
	defer visLock.Unlock()

	visible = make(map[[2]Coord]bool)
}

func init() {
	visInvalidateAll()
}

func visInvalidate(coord Coord) {
	visLock.Lock()
	defer visLock.Unlock()

	for link := range visible {
		if link[0] == coord || link[1] == coord {
			delete(visible, link)
		}
	}
}

func visInvalidateRecursive(coord Coord) {
	visLock.Lock()
	defer visLock.Unlock()

	for link := range visible {
		if link[0] == coord && link[1] == coord {
			continue
		} else if link[0] == coord {
			for link2 := range visible {
				if link2[0] == link[1] || link2[1] == link[1] {
					delete(visible, link2)
				}
			}
			delete(visible, link)
		} else if link[1] == coord {
			for link2 := range visible {
				if link2[0] == link[0] || link2[1] == link[0] {
					delete(visible, link2)
				}
			}
			delete(visible, link)
		}
	}
}

func Visible(a, b Coord) bool {
	visLock.RLock()
	if seen, ok := visible[[2]Coord{a, b}]; ok {
		visLock.RUnlock()
		return seen
	}
	visLock.RUnlock()

	visLock.Lock()
	defer visLock.Unlock()
	// double-check
	if seen, ok := visible[[2]Coord{a, b}]; ok {
		return seen
	}

	seen := visTrace(a.X, a.Y, b.X, b.Y, 1)

	visible[[2]Coord{a, b}] = seen
	return seen
}

func visTrace(ax, ay, bx, by, off int) bool {
	if ax == bx && ay == by {
		return true
	}

	var vis bool
	dx, dy := bx-ax, by-ay

	if Get(ax, ay).BlocksVision() && (dx*dx+dy*dy != 1 || !Get(bx, by).BlocksVision()) {
		return false
	}

	if dx*dx < dy*dy {
		if dy > 0 {
			vis = visTrace(ax, ay+1, bx, by, off)
		} else {
			vis = visTrace(ax, ay-1, bx, by, off)
		}
		if !vis && off > 0 {
			if dx > 0 {
				vis = visTrace(ax+1, ay, bx, by, off-1)
			} else {
				vis = visTrace(ax-1, ay, bx, by, off-1)
			}
		}
		return vis
	} else {
		if dx > 0 {
			vis = visTrace(ax+1, ay, bx, by, off)
		} else {
			vis = visTrace(ax-1, ay, bx, by, off)
		}
		if !vis && off > 0 {
			if dy > 0 {
				vis = visTrace(ax, ay+1, bx, by, off-1)
			} else {
				vis = visTrace(ax, ay-1, bx, by, off-1)
			}
		}
		return vis
	}
	panic("unreachable")
}
