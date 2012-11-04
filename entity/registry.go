package entity

import (
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
