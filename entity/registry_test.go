package entity

import "testing"

// TestEntity defined in entity_test.go

func TestEntityListInsertionOrder(t *testing.T) {
	nextID = 0

	var entities [250]TestEntity

	var list EntityList = new(entityList)

	for i := range entities {
		if !list.Add(&entities[i]) {
			t.Errorf("Entity %d was not added to the list", entities[i].ID())
		}
	}

	if c := list.Count(); c != 250 {
		t.Errorf("Entity count was %d, but expected 250", c)
	}

	var prev EntityID
	list.Each(func(ent Entity) {
		if prev >= ent.ID() {
			t.Errorf("Entity %d came before entity %d", prev, ent.ID())
		}
		prev = ent.ID()
	})
}
