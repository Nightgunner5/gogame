package entity

import (
	"github.com/Nightgunner5/gogame/log"
	"sort"
	"sync"
)

type EntityList []Entity

var (
	entityList     = EntityList{}
	entityListLock sync.RWMutex
)

func (list EntityList) Less(i, j int) bool {
	return list[i].ID() < list[j].ID()
}

func (list EntityList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list EntityList) Len() int {
	return len(list)
}

func GetEntityByID(id EntityID) Entity {
	entityListLock.RLock()
	defer entityListLock.RUnlock()

	i := sort.Search(len(entityList), func(i int) bool {
		return entityList[i].ID() >= id
	})

	if i < len(entityList) && entityList[i].ID() == id {
		return entityList[i]
	}
	return nil
}

func RegisterEntity(entity Entity) {
	entityListLock.Lock()
	defer entityListLock.Unlock()

	i := sort.Search(len(entityList), func(i int) bool {
		return entityList[i].ID() >= entity.ID()
	})

	if i < len(entityList) {
		entityList = append(entityList, entity)
		return
	}

	log.Warning("Entity registration for ID %d, even though %d is the highest", entity.ID(), entityList[len(entityList)-1].ID())

	if entityList[i].ID() == entity.ID() {
		if entityList[i] == entity {
			log.Warning("Duplicate registration for entity %d", entity.ID())
		} else {
			log.Panic("Duplicate entity ID %d for two separate entities", entity.ID())
		}
		return
	}

	entityList = append(entityList[:i], append(EntityList{entity}, entityList[i:]...)...)
}

func removeFromEntList(id EntityID) {
	entityListLock.Lock()
	defer entityListLock.Unlock()

	i := sort.Search(len(entityList), func(i int) bool {
		return entityList[i].ID() <= id
	})

	if i < len(entityList) && entityList[i].ID() == id {
		entityList = append(entityList[:i], entityList[i+1:]...)
	}
}
