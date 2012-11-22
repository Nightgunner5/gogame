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
	EffectThink(delta float64)
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

func BaseEffectAdder(ent entity.Entity) EffectAdder {
	return &basicEffectAdder{ent: ent.ID()}
}

type basicEffectAdder struct {
	entity.Listeners
	effects []*basicEffect
	mtx     sync.RWMutex
	ent     entity.EntityID
}

func (b *basicEffectAdder) AddEffect(effect Effect, duration float64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if duration <= 0 {
		duration = 0
	}
	b.effects = append(b.effects, &basicEffect{Effect: effect, duration: duration})

	b.AddAll(effect)
}

func (b *basicEffectAdder) EffectThink(delta float64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	var removed []*basicEffect

	for i, effect := range b.effects {
		if effect.duration <= 0 && effect.Effect.String() != "" {
			if removed != nil {
				removed = append(removed, effect)
			}
			continue
		}
		if effect.duration <= delta || effect.Effect.String() == "" {
			if removed == nil {
				removed = make([]*basicEffect, i, len(b.effects))
				copy(removed, b.effects)
			}
			b.RemoveAll(effect.Effect)
			continue
		}
		effect.duration -= delta
		if removed != nil {
			removed = append(removed, effect)
		}
	}

	if removed != nil {
		b.effects = removed
	}
}

func (b *basicEffectAdder) Effects() []Effect {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	l := make([]Effect, len(b.effects))
	for i, effect := range b.effects {
		l[i] = effect.Effect
	}
	return l
}

func (b *basicEffectAdder) EffectDescription() string {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	var s []byte
	for _, effect := range b.effects {
		s = append(s, effect.String()...)
		s = append(s, '\n')
	}

	return string(s)
}
