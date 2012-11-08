package spell

import (
	"github.com/Nightgunner5/gogame/entity"
	"sync"
)

type Spell interface {
	// Returns true if the spell has ended, false if it has not.
	Tick(Δtime float64) bool

	Caster() entity.Entity
	Target() entity.Entity

	TotalTime() float64
	TimeLeft() float64

	// Interrupts the spell and returns true. If the spell is not
	// interruptable, returns false and does nothing.
	Interrupt() bool
}

type Caster interface {
	Interrupt() bool
	CurrentSpell() Spell
}

type SpellCaster struct {
	current Spell
	m       sync.Mutex
}

func (c *SpellCaster) CasterThink(Δtime float64) bool {
	c.m.Lock()
	defer c.m.Unlock()

	spell := c.current
	if spell == nil {
		return false
	}
	if spell.Tick(Δtime) {
		c.current = nil
		return false
	}
	return true
}

func (c *SpellCaster) CurrentSpell() Spell {
	c.m.Lock()
	defer c.m.Unlock()

	return c.current
}

func (c *SpellCaster) Cast(spell Spell) {
	c.m.Lock()
	defer c.m.Unlock()

	c.current = spell
}

func (c *SpellCaster) Interrupt() bool {
	c.m.Lock()
	defer c.m.Unlock()

	spell := c.current
	if spell != nil && spell.Interrupt() {
		c.current = nil
		return true
	}
	return false
}
