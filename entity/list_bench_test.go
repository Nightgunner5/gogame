package entity

import (
	"sync"
	"testing"
)

type nullEntity struct {
	EntityID
}

func (nullEntity) Parent() Entity {
	return World
}

func concurrentBench(b *testing.B, f func(int)) {
	var wg sync.WaitGroup

	wg.Add(b.N)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		go func(i int) {
			f(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	NukeForTesting()
	for i := 0; i < 10000; i++ {
		Spawn(new(nullEntity))
	}

	concurrentBench(b, func(i int) {
		_ = Get(1)
	})
}

func BenchmarkSpawn(b *testing.B) {
	b.StopTimer()
	NukeForTesting()
	entities := make([]nullEntity, b.N)

	concurrentBench(b, func(i int) {
		Spawn(&entities[i])
	})
}

func BenchmarkDespawn(b *testing.B) {
	b.StopTimer()
	NukeForTesting()
	entities := make([]nullEntity, b.N)
	for i := 0; i < b.N; i++ {
		Spawn(&entities[i])
	}

	concurrentBench(b, func(i int) {
		Despawn(&entities[i])
	})
}
