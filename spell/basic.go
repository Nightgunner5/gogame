package spell

import (
	"github.com/Nightgunner5/gogame/entity"
	"sync"
)

type BasicSpell struct {
	CastTime      float64
	currentTime   float64
	Interruptable bool
	Caster_       entity.EntityID
	Target_       entity.EntityID
	Action        func(target, caster entity.Entity)
	Tag_          string

	m sync.Mutex
}

func (s *BasicSpell) Caster() entity.Entity {
	return entity.Get(s.Caster_)
}

func (s *BasicSpell) CasterID() entity.EntityID {
	return s.Caster_
}

func (s *BasicSpell) Target() entity.Entity {
	return entity.Get(s.Target_)
}

func (s *BasicSpell) TargetID() entity.EntityID {
	return s.Target_
}

func (s *BasicSpell) Tag() string {
	return s.Tag_
}

func (s *BasicSpell) Interrupt() bool {
	s.m.Lock()
	defer s.m.Unlock()

	if !s.Interruptable || s.currentTime >= s.CastTime {
		return false
	}
	s.currentTime = s.CastTime
	return true
}

func (s *BasicSpell) TotalTime() float64 {
	return s.CastTime
}

func (s *BasicSpell) TimeLeft() float64 {
	s.m.Lock()
	defer s.m.Unlock()

	return s.CastTime - s.currentTime
}

func (s *BasicSpell) Tick(Δtime float64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	if s.currentTime >= s.CastTime {
		return true
	}
	s.currentTime += Δtime
	if s.currentTime >= s.CastTime {
		s.currentTime = s.CastTime
		target, caster := s.Target(), s.Caster()
		if target == nil || caster == nil {
			return true
		}

		s.Action(target, caster)

		return true
	}
	return false
}

var _ Spell = new(BasicSpell)
