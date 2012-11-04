package entity

import (
	"log"
	"sort"
	"sync"
)

type EntityList interface {
	Get(EntityID) Entity
	Add(Entity) bool
	Remove(EntityID) Entity

	Count() int
	Each(func(Entity))
	All(func(Entity))
}

type entityList []Entity

func (list entityList) search(id EntityID) (i int, found bool) {
	i = sort.Search(len(list), func(i int) bool {
		return list[i].ID() >= id
	})
	found = i < len(list) && list[i].ID() == id
	return
}

func (list entityList) Get(id EntityID) Entity {
	i, found := list.search(id)
	if found {
		return list[i]
	}
	return nil
}

func (list *entityList) Add(entity Entity) bool {
	i, found := list.search(entity.ID())
	if found {
		return false
	}

	*list = append((*list)[:i], append(entityList{entity}, (*list)[i:]...)...)

	return true
}

func (list *entityList) Remove(id EntityID) Entity {
	i, found := list.search(id)

	var ret Entity
	if found {
		ret = (*list)[i]
		*list = append((*list)[:i], (*list)[i+1:]...)
	}
	return ret
}

func (list entityList) Count() int {
	return len(list)
}

func (list entityList) Each(f func(Entity)) {
	for i := range list {
		f(list[i])
	}
}

func (list entityList) All(f func(Entity)) {
	var wg sync.WaitGroup

	wg.Add(list.Count())

	list.Each(func(e Entity) {
		go func() {
			f(e)
			wg.Done()
		}()
	})

	wg.Wait()
}

type concurrentEntityList struct {
	l entityList
	m sync.RWMutex
}

func (c *concurrentEntityList) Get(id EntityID) Entity {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.l.Get(id)
}

func (c *concurrentEntityList) Add(entity Entity) bool {
	c.m.Lock()
	defer c.m.Unlock()

	return c.l.Add(entity)
}

func (c *concurrentEntityList) Remove(id EntityID) Entity {
	c.m.Lock()
	defer c.m.Unlock()

	return c.l.Remove(id)
}

func (c *concurrentEntityList) Count() int {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.l.Count()
}

func (c *concurrentEntityList) Each(f func(Entity)) {
	c.m.RLock()
	defer c.m.RUnlock()

	c.l.Each(f)
}

func (c *concurrentEntityList) All(f func(Entity)) {
	c.m.RLock()

	var wg sync.WaitGroup

	wg.Add(c.l.Count())

	c.l.Each(func(e Entity) {
		go func() {
			f(e)
			wg.Done()
		}()
	})

	c.m.RUnlock()

	wg.Wait()
}

func NewEntityList(capacity int) EntityList {
	list := make(entityList, 0, capacity)
	return &list
}

func ConcurrentEntityList(capacity int) EntityList {
	list := make(entityList, 0, capacity)
	return &concurrentEntityList{l: list}
}

type readOnlyEntityList struct {
	EntityList
}

func (readOnlyEntityList) Add(Entity) bool {
	return false
}

func (readOnlyEntityList) Remove(EntityID) Entity {
	return nil
}

var globalEntityList EntityList = ConcurrentEntityList(1)

func GlobalEntityList() EntityList {
	return readOnlyEntityList{globalEntityList}
}

func Spawn(entity Entity) {
	if !globalEntityList.Add(entity) {
		log.Printf("SpawnEntity called twice on entity %v", entity)
	}
}

func Despawn(entity Entity) {
	if e := globalEntityList.Remove(entity.ID()); e != entity {
		if e == nil {
			log.Printf("DespawnEntity called on non-spawned entity %v", entity)
		} else {
			log.Panicf("DespawnEntity: Entity ID %d is used multiple times: %v %v", entity.ID(), e, entity)
		}
	}
}

func Get(id EntityID) Entity {
	return globalEntityList.Get(id)
}

func ForAll(f func(Entity)) {
	globalEntityList.All(f)
}

func ForAllNearby(target Positioner, distance float64, f func(Entity)) {
	sX, sY, sZ := target.Position()
	d2 := distance * distance
	globalEntityList.All(func(e Entity) {
		if p, ok := e.(Positioner); ok {
			x, y, z := p.Position()
			x, y, z = sX-x, sY-y, sZ-z
			x, y, z = x*x, y*y, z*z

			if x+y+z <= d2 {
				f(e)
			}
		}
	})
}
