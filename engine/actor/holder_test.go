package actor

import (
	"testing"
)

func TestHolder(t *testing.T) {
	t.Parallel()

	var a Actor
	TopLevel(a.Initialize())

	var h Holder
	TopLevel(h.Initialize())

	if l := len(h.GetHeld()); l != 0 {
		t.Errorf("Expected no held actors, but there are %d", l)
	}

	sub, s := Subscribe(MsgAddHeld, 1)
	h.Send <- AddSubscriber{sub}
	h.Send <- AddHeld{&a}
	if a2 := (<-s).(AddHeld).Actor; a2 != &a {
		t.Errorf("Expected added actor %v but got %v", &a, a2)
	}
	h.Send <- RemoveSubscriber{sub}

	held := h.GetHeld()
	if l := len(held); l != 1 {
		t.Errorf("Expected one held actor, but there are %d", l)
	}

	if held[0] != &a {
		t.Errorf("Expected held actor to be %v, but got %v", &a, held[0])
	}

	sub, s = Subscribe(MsgRemoveHeld, 1)
	h.Send <- AddSubscriber{sub}
	h.Send <- RemoveHeld{&a}
	if a2 := (<-s).(RemoveHeld).Actor; a2 != &a {
		t.Errorf("Expected removed actor %v but got %v", &a, a2)
	}
	h.Send <- RemoveSubscriber{sub}

	held = h.GetHeld()
	if l := len(held); l != 0 {
		t.Errorf("Expected no held actors, but there are %d", l)
	}
}
