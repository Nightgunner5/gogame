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
	entities := make([]nullEntity, b.N)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if i%10000 == 0 {
			b.StopTimer()
			nukeForTesting()
			b.StartTimer()
		}
		Spawn(&entities[i])
	}
}

func BenchmarkDespawn(b *testing.B) {
	b.StopTimer()
	nukeForTesting()
	entities := make([]nullEntity, b.N)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if i%10000 == 0 {
			b.StopTimer()
			nukeForTesting()
			for j := i; j < i+10000 && j < b.N; j++ {
				Spawn(&entities[j])
			}
			b.StartTimer()
		}
		Despawn(&entities[i])
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
		ForEach(func(e Entity) {
			// empty function
		})
	}
}
