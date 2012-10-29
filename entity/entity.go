package entity

import (
	"encoding/gob"
	"sync/atomic"
	"github.com/Nightgunner5/gogame/log"
)

type EntityID uint64

type Entity interface {
	ID() EntityID

	Position() (x, y, z float64)
	SetPosition(x, y, z float64)

	AcceptEvent(data EventData)
}

type Base struct {
	id      EntityID
	x, y, z float64
}

func init() {
	RegisterEntityType((*Base)(nil))
}

const (
	EventKill = "Kill"
)

func (b *Base) ID() EntityID {
	return b.id
}

func (b *Base) Position() (x, y, z float64) {
	x, y, z = b.x, b.y, b.z
	return
}

func (b *Base) SetPosition(x, y, z float64) {
	b.x, b.y, b.z = x, y, z
}

func (b *Base) AcceptEvent(data EventData) {
	switch data.Type {
	case EventKill:
		removeFromEntList(b.ID())

	default:
		log.Panic("Entity %d: Unknown event type %q", b.ID(), data.Type)
	}
}

var lastID uint64

// Return a new Base with a unque Entity.ID chosen sequentially in a thread-safe manner. The Entity.Position is 0, 0, 0
func BaseEntity() (ent *Base) {
	ent = &Base{
		id: EntityID(atomic.AddUint64(&lastID, 1)),
	}
	return
}

func RegisterEntityType(instance Entity) {
	gob.Register(instance)
}
