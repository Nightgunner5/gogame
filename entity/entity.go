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
	ID() EntityID

	AcceptEvent(Event)
}
