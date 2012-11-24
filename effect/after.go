package effect

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"sync"
)

type after struct {
	add      *Effect
	timeLeft float64
	sync.Mutex
}

func After(effect *Effect, delay float64) Primitive {
	return &after{
		add:      effect,
		timeLeft: delay,
	}
}

func (a *after) String() string {
	a.Lock()
	defer a.Unlock()

	if a.timeLeft <= 0 {
		return ""
	}

	t := int(a.timeLeft)
	if t <= 1 {
		return fmt.Sprintf("After 1 second:\n\n%s\n", a.add)
	}
	return fmt.Sprintf("After %d seconds:\n\n%s\n", t, a.add)
}

func (*after) isSameType(Primitive) bool {
	return false
}

func (*after) combine(Primitive) Primitive {
	panic("combine on effect.after is unimplemented by design")
}

func (a *after) effectThink(delta float64, ent entity.Entity) {
	a.Lock()
	defer a.Unlock()

	if a.timeLeft <= 0 {
		return
	}

	a.timeLeft -= delta
	if a.timeLeft <= 0 {
		// This has to be done in a goroutine since the effect list is locked during effectThink
		go ent.(Effected).AddEffect(a.add)
	}
}
