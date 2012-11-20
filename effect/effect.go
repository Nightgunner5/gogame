package effect

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
	"sync"
)

type Effect interface {
	// Returning an empty string means "this effect is no longer active
	// and may be deleted"
	String() string

	effect() // marker
}

type EffectAdder interface {
	AddEffect(effect Effect, duration float64)
	Effects() []Effect
	EffectDescription() string
}

type basicEffect struct {
	Effect
	duration float64
}

func (b *basicEffect) String() string {
	s := b.Effect.String()
	if s == "" || b.duration <= 0 {
		return s
	}
	if int(b.duration) <= 1 {
		return fmt.Sprintf("%s (1 second remaining)", s)
	}
	return fmt.Sprintf("%s (%d seconds remaining)", s, int(b.duration))
}

type BasicEffectAdder struct {
	entity.Listeners
	e []*basicEffect
	m sync.RWMutex
}

func (b *BasicEffectAdder) AddEffect(effect Effect, duration float64) {
	b.m.Lock()
	defer b.m.Unlock()

	if duration <= 0 {
		duration = 0
	}
	b.e = append(b.e, &basicEffect{Effect: effect, duration: duration})

	b.AddAll(effect)
}

func (b *BasicEffectAdder) EffectThink(delta float64) {
	b.m.Lock()
	defer b.m.Unlock()

	var removed []*basicEffect

	for i, e := range b.e {
		if e.duration <= 0 && e.Effect.String() != "" {
			if removed != nil {
				removed = append(removed, e)
			}
			continue
		}
		if e.duration <= delta || e.Effect.String() == "" {
			if removed == nil {
				removed = make([]*basicEffect, i, len(b.e))
				copy(removed, b.e)
			}
			b.RemoveAll(e.Effect)
			continue
		}
		e.duration -= delta
		if removed != nil {
			removed = append(removed, e)
		}
	}

	if removed != nil {
		b.e = removed
	}
}

func (b *BasicEffectAdder) Effects() []Effect {
	b.m.RLock()
	defer b.m.RUnlock()

	l := make([]Effect, len(b.e))
	for i, e := range b.e {
		l[i] = e.Effect
	}
	return l
}

func (b *BasicEffectAdder) EffectDescription() string {
	b.m.RLock()
	defer b.m.RUnlock()

	var s []byte
	for _, e := range b.e {
		s = append(s, e.String()...)
		s = append(s, '\n')
	}

	return string(s)
}
