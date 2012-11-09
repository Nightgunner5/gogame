package entity

import "sync"

type EntityID uint64

var nextID EntityID
var idLock sync.Mutex

func (id *EntityID) ID() EntityID {
	if *id == 0 {
		idLock.Lock()
		defer idLock.Unlock()
		if *id != 0 {
			return *id
		}
		nextID++
		*id = nextID
	}
	return *id
}

type Entity interface {
	// Entities have unique IDs. This can most easily be accomplished
	// by embedding EntityID and always using the ID method instead of
	// accessing the EntityID directly. EntityID.ID() handles locking
	// and initializing on its own.
	ID() EntityID

	// The parent of this entity. Top-level entities have a parent of
	// the World entity. Only the World entity has a nil parent.
	Parent() Entity
}
