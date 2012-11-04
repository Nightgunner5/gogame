package entity

import (
	"sync"
)

type EntityID uint64

var nextID EntityID
var idLock sync.Mutex

func (id *EntityID) ID() EntityID {
	if *id == 0 {
		idLock.Lock()
		nextID++
		*id = nextID
		idLock.Unlock()
	}
	return *id
}

type Entity interface {
	ID() EntityID

	AcceptEvent(Event)
}
