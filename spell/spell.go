package spell

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/network"
	"sync"
)

type Spell interface {
	// Returns true if the spell has ended, false if it has not.
	Tick(float64) bool

	Caster() entity.Entity
	CasterID() entity.EntityID
	Target() entity.Entity
	TargetID() entity.EntityID

	TotalTime() float64
	TimeLeft() float64

	// Interrupts the spell and returns true. If the spell is not
	// interruptable, returns false and does nothing.
	Interrupt() bool

	Tag() string
}

type SpellCaster interface {
	Interrupt() bool
	CurrentSpell() Spell
	Cast(Spell)
	CasterThink(float64) bool
}

type baseSpellCaster struct {
	ent   entity.EntityID
	spell Spell
	sync.Mutex
}

func BaseSpellCaster(ent entity.Entity) SpellCaster {
	return &baseSpellCaster{
		ent: ent.ID(),
	}
}

func (b *baseSpellCaster) CasterThink(delta float64) bool {
	b.Lock()
	defer b.Unlock()

	if b.spell == nil {
		return false
	}
	if b.spell.Tick(delta) {
		b.spell = nil
		return false
	}
	return true
}

func (b *baseSpellCaster) CurrentSpell() Spell {
	b.Lock()
	defer b.Unlock()

	return b.spell
}

func (b *baseSpellCaster) Cast(spell Spell) {
	b.Lock()
	defer b.Unlock()

	b.spell = spell
	network.Broadcast(network.NewPacket(network.CastSpell).
		Set(network.EntityID, b.ent).
		Set(network.OtherEntID, spell.TargetID()).
		Set(network.Tag, spell.Tag()).
		Set(network.TimeLeft, spell.TimeLeft()).
		Set(network.TotalTime, spell.TotalTime()), false)
}

func (b *baseSpellCaster) Interrupt() bool {
	b.Lock()
	defer b.Unlock()

	if b.spell != nil && b.spell.Interrupt() {
		b.spell = nil
		network.Broadcast(network.NewPacket(network.CastSpell).
			Set(network.EntityID, b.ent).
			Set(network.OtherEntID, 0).
			Set(network.Tag, "interrupted").
			Set(network.TimeLeft, 0).
			Set(network.TotalTime, 0), false)
		return true
	}
	return false
}
