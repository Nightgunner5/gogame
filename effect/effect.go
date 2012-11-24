package effect

import (
	"fmt"
	"github.com/Nightgunner5/gogame/entity"
)

type Effect struct {
	name        string
	duration    float64
	currentTime float64
	primitives  []Primitive
}

var _ entity.AllListeners = new(Effect)

func NewEffect(name string, duration float64) *Effect {
	return &Effect{
		name:     name,
		duration: duration,
	}
}

func (e *Effect) Add(p Primitive) *Effect {
	for i := range e.primitives {
		if e.primitives[i].isSameType(p) {
			e.primitives[i] = e.primitives[i].combine(p)
			return e
		}
	}

	e.primitives = append(e.primitives, p)

	return e
}

func (e *Effect) TimeLeft() float64 {
	if e.duration <= 0 {
		return -1
	}
	if e.duration > e.currentTime {
		return e.duration - e.currentTime
	}
	return 0
}

func (e *Effect) String() string {
	s := e.name
	for _, p := range e.primitives {
		s = fmt.Sprintf("%s\n%v", s, p)
	}
	if e.duration > 0 {
		timeLeft := int(e.TimeLeft())
		if timeLeft <= 1 {
			s = fmt.Sprintf("%s\n(1 second remaining)", s)
		} else {
			s = fmt.Sprintf("%s\n(%d seconds remaining)", s, timeLeft)
		}
	}
	return s
}

func (e *Effect) effectThink(delta float64, ent entity.Entity) (changed bool) {
	changed = int(e.currentTime) != int(e.currentTime + delta)
	e.currentTime += delta

	for _, p := range e.primitives {
		changed = p.effectThink(delta, ent) || changed
	}
	return
}

func (e *Effect) OnTakeDamage(amount *float64, attacker, victim entity.Entity) (changed bool) {
	for _, p := range e.primitives {
		if l, ok := p.(entity.DamageListener); ok {
			changed = l.OnTakeDamage(amount, attacker, victim) || changed
		}
	}
	return
}

func (e *Effect) OnDoDamage(amount *float64, attacker, victim entity.Entity) (changed bool) {
	for _, p := range e.primitives {
		if l, ok := p.(entity.DoDamageListener); ok {
			changed = l.OnDoDamage(amount, attacker, victim) || changed
		}
	}
	return
}
