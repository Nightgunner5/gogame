package entity

import "testing"

// nullEntity defined in entity_test.go

func TestEntityListInsertRemove(t *testing.T) {
	nukeForTesting()

	var ent nullEntity

	list := NewEntityList(1)

	if !list.Add(&ent) {
		t.Error("Could not add entity to list (once)")
	}

	if list.Add(&ent) {
		t.Error("Added entity to list (twice)")
	}

	if e := list.Get(ent.ID()); e != &ent {
		t.Errorf("Get: Expected entity %v but found %v", &ent, e)
	}

	if list.Count() != 1 {
		t.Errorf("%d elements in single-element list", list.Count())
	}

	if e := list.Remove(ent.ID()); e != &ent {
		t.Errorf("Remove: Expected entity %v but found %v", &ent, e)
	}

	if e := list.Get(ent.ID()); e != nil {
		t.Errorf("Get: Expected nil but found %v", e)
	}

	if e := list.Remove(ent.ID()); e != nil {
		t.Errorf("Remove: Expected nil but found %v", e)
	}

	if list.Count() != 0 {
		t.Errorf("%d elements in empty list", list.Count())
	}

	if !list.Add(&ent) {
		t.Error("Could not re-add entity to list after removing it")
	}

	if list.Count() != 1 {
		t.Errorf("%d elements in single-element list", list.Count())
	}
}

func TestEntityListInsertionOrder(t *testing.T) {
	nextID = 0

	var entities [250]nullEntity

	list := NewEntityList(250)

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
