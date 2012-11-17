package entity

import "testing"

func TestBaseCounterExpend(t *testing.T) {
	t.Parallel()

	var c baseCounter

	if !c.sub(-100, 100, true) {
		t.Error("could not add to counter")
	}
	if n := c.get(100); n != 100 {
		t.Errorf("expected 100, got %v", n)
	}

	if !c.sub(40, 100, false) {
		t.Error("could not subtract from counter")
	}
	if n := c.get(100); n != 60 {
		t.Errorf("expected 60 (1), got %v", n)
	}

	if c.sub(100, 100, false) {
		t.Error("subtracted too much from counter")
	}
	if n := c.get(100); n != 60 {
		t.Errorf("expected 60 (2), got %v", n)
	}

	if !c.sub(60, 100, false) {
		t.Error("could not expend all from counter")
	}
	if n := c.get(100); n != 0 {
		t.Errorf("expected 0, got %v", n)
	}

	if c.sub(1, 100, false) {
		t.Error("expended from empty counter")
	}
}
