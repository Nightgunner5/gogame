package message

import "testing"

func TestKindEquality(t *testing.T) {
	k1, k2, k3 := NewKind("foo"), NewKind("foo"), NewKind("bar")
	k4 := k1

	if k1 == k2 {
		t.Errorf("Unexpected equality: %v %v", k1, k2)
	}

	if k2 == k3 {
		t.Errorf("Unexpected equality: %v %v", k2, k3)
	}

	if k3 == k4 {
		t.Errorf("Unexpected equality: %v %v", k3, k4)
	}

	if k1 != k4 {
		t.Errorf("Unexpected inequality: %v %v", k1, k4)
	}
}
