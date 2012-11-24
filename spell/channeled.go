package spell

import (
	"github.com/Nightgunner5/gogame/entity"
	"sync"
)

type ChanneledSpell struct {
	CastTime      float64
	currentTime   float64
	Interruptable bool
	Caster_       entity.EntityID
	Target_       entity.EntityID
	Action        func(target, caster entity.Entity, Δtime float64)
	Tag_          string

	m sync.Mutex
}

func (s *ChanneledSpell) Caster() entity.Entity {
	return entity.Get(s.Caster_)
}

func (s *ChanneledSpell) CasterID() entity.EntityID {
	return s.Caster_
}

func (s *ChanneledSpell) Target() entity.Entity {
	return entity.Get(s.Target_)
}

func (s *ChanneledSpell) TargetID() entity.EntityID {
	return s.Target_
}

func (s *ChanneledSpell) Tag() string {
	return s.Tag_
}

func (s *ChanneledSpell) Interrupt() bool {
	s.m.Lock()
	defer s.m.Unlock()

	if !s.Interruptable || s.currentTime >= s.CastTime {
		return false
	}
	s.currentTime = s.CastTime
	return true
}

func (s *ChanneledSpell) TimeLeft() float64 {
	s.m.Lock()
	defer s.m.Unlock()

	return s.CastTime - s.currentTime
}

func (s *ChanneledSpell) TotalTime() float64 {
	return s.CastTime
}

func (s *ChanneledSpell) Tick(Δtime float64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	if s.currentTime >= s.CastTime {
		return true
	}
	prevTime := s.currentTime
	s.currentTime += Δtime
	if s.currentTime >= s.CastTime {
		s.currentTime = s.CastTime
	}

	target, caster := s.Target(), s.Caster()
	if target == nil || caster == nil {
		return true
	}

	s.Action(target, caster, s.currentTime-prevTime)

	if s.currentTime >= s.CastTime {
		return true
	}
	return false
}

var _ Spell = new(BasicSpell)
