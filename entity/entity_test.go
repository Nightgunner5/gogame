package entity

import "testing"

type TestEntity struct {
	EntityID
}

func TestEntityID(t *testing.T) {
	nextID = 0

	var ents [250]TestEntity

	// Assign IDs in a nondeterministic order.
	for i := range ents {
		go ents[i].ID()
	}

	var entIDCounts [250]int
	// We can start looping immediately because of the synchronization
	for i := range ents {
		id := ents[i].ID()
		if id < 1 || id > 250 {
			t.Errorf("Entity ID out of range [1,250]: %d", id)
		}
		entIDCounts[int(id-1)]++
	}

	for i := range entIDCounts {
		if entIDCounts[i] != 1 {
			t.Errorf("Entity ID %d has %d uses, but expected 1", i+1, entIDCounts[i])
		}
	}
}
