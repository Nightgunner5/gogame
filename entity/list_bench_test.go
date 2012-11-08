package entity

import (
	"testing"
)

type nullEntity struct {
	EntityID
}

func (nullEntity) Parent() Entity {
	return World
}

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	NukeForTesting()
	for i := 0; i < 10000; i++ {
		Spawn(new(nullEntity))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = Get(EntityID(i % 10000))
	}
}

func BenchmarkSpawn(b *testing.B) {
	b.StopTimer()
	NukeForTesting()
	entities := make([]nullEntity, b.N)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Spawn(&entities[i])
	}
}

func BenchmarkDespawn(b *testing.B) {
	b.StopTimer()
	NukeForTesting()
	entities := make([]nullEntity, b.N)
	for i := 0; i < b.N; i++ {
		Spawn(&entities[i])
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Despawn(&entities[i])
	}
}
