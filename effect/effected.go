package effect

import (
	"github.com/Nightgunner5/gogame/entity"
	"sync"
)

type Effected interface {
	Effects() []*Effect
	AddEffect(*Effect)
	EffectThink(float64)
	entity.AllListeners
}

func BaseEffected(ent entity.Entity) Effected {
	return &baseEffected{
		ent: ent.ID(),
	}
}

type baseEffected struct {
	ent     entity.EntityID
	effects []*Effect
	sync.Mutex
}

var _ Effected = new(baseEffected)

func (b *baseEffected) Effects() []*Effect {
	b.Lock()
	defer b.Unlock()

	effects := make([]*Effect, len(b.effects))
	copy(effects, b.effects)

	return effects
}

func (b *baseEffected) AddEffect(effect *Effect) {
	b.Lock()
	defer b.Unlock()

	if effect.TimeLeft() != 0 {
		b.effects = append(b.effects, effect)
	}

	broadcastPacket(b.ent, b.effects)
}

func (b *baseEffected) EffectThink(delta float64) {
	b.Lock()
	defer b.Unlock()

	ent := entity.Get(b.ent)

	changed := false
	for i, e := range b.effects {
		changed = e.effectThink(delta, ent) || changed
		if e.TimeLeft() == 0 {
			b.effects[i] = nil
			changed = true
		}
	}

	if changed {
		for i := 0; i < len(b.effects); i++ {
			if b.effects[i] == nil {
				b.effects = append(b.effects[:i], b.effects[i+1:]...)
				i--
			}
		}

		broadcastPacket(b.ent, b.effects)
	}
}

func (b *baseEffected) OnDoDamage(amount *float64, attacker, victim entity.Entity) (changed bool) {
	b.Lock()
	defer b.Unlock()

	for _, e := range b.effects {
		changed = e.OnDoDamage(amount, attacker, victim) || changed
	}
	if changed {
		broadcastPacket(b.ent, b.effects)
	}
	return
}

func (b *baseEffected) OnTakeDamage(amount *float64, attacker, victim entity.Entity) (changed bool) {
	b.Lock()
	defer b.Unlock()

	for _, e := range b.effects {
		changed = e.OnTakeDamage(amount, attacker, victim) || changed
	}
	if changed {
		broadcastPacket(b.ent, b.effects)
	}
	return
}
