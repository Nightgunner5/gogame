package bench

import (
	"github.com/Nightgunner5/gogame/entity"
	"sync"
	"testing"
)

type nullEntity struct {
	entity.EntityID
}

func (nullEntity) Parent() entity.Entity {
	return entity.World
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
	entity.NukeForTesting()
	for i := 0; i < 10000; i++ {
		entity.Spawn(new(nullEntity))
	}

	concurrentBench(b, func(i int) {
		_ = entity.Get(1)
	})
}

func BenchmarkSpawn(b *testing.B) {
	b.StopTimer()
	entity.NukeForTesting()

	concurrentBench(b, func(i int) {
		entity.Spawn(new(nullEntity))
	})
}

func BenchmarkDespawn(b *testing.B) {
	b.StopTimer()
	entity.NukeForTesting()
	entities := make([]nullEntity, b.N)
	for i := 0; i < b.N; i++ {
		entity.Spawn(&entities[i])
	}

	concurrentBench(b, func(i int) {
		entity.Despawn(&entities[i])
	})
}
