package entity

import "testing"

// nullEntity defined in entity_test.go

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	nukeForTesting()
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
	nukeForTesting()
	entities := make([]nullEntity, b.N)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Spawn(&entities[i])
	}
}

func BenchmarkDespawn(b *testing.B) {
	b.StopTimer()
	nukeForTesting()
	entities := make([]nullEntity, b.N)
	for i := 0; i < b.N; i++ {
		Spawn(&entities[i])
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Despawn(&entities[i])
	}
}

func BenchmarkForAll(b *testing.B) {
	b.StopTimer()
	nukeForTesting()
	for i := 0; i < 10000; i++ {
		Spawn(new(nullEntity))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ForAll(func(e Entity){
			// empty function
		})
	}
}

func BenchmarkForEach(b *testing.B) {
	b.StopTimer()
	nukeForTesting()
	for i := 0; i < 10000; i++ {
		Spawn(new(nullEntity))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		GlobalEntityList().Each(func(e Entity){
			// empty function
		})
	}
}
