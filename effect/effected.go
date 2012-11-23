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
}

func (b *baseEffected) EffectThink(delta float64) {
	b.Lock()
	defer b.Unlock()

	ent := entity.Get(b.ent)

	for i, e := range b.effects {
		e.effectThink(delta, ent)
		if e.TimeLeft() == 0 {
			b.effects[i] = nil
		}
	}

	for i := 0; i < len(b.effects); i++ {
		if b.effects[i] == nil {
			b.effects = append(b.effects[:i], b.effects[i+1:]...)
			i--
		}
	}
}

func (b *baseEffected) OnDoDamage(amount *float64, attacker, victim entity.Entity) {
	b.Lock()
	defer b.Unlock()

	for _, e := range b.effects {
		e.OnDoDamage(amount, attacker, victim)
	}
}

func (b *baseEffected) OnTakeDamage(amount *float64, attacker, victim entity.Entity) {
	b.Lock()
	defer b.Unlock()

	for _, e := range b.effects {
		e.OnTakeDamage(amount, attacker, victim)
	}
}
