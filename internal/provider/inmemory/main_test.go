package inmemory

import (
	"runtime"
	"sync"
	"testing"
)

func BenchmarkCache_GetShard(b *testing.B) {
	cache := New(128)

	for i := 0; b.Loop(); i++ {
		cache.GetShard(i)
	}
}

func BenchmarkCache_Increment_Sequential(b *testing.B) {
	cache := New(128)

	for i := 0; b.Loop(); i++ {
		sh := cache.GetShard(i)
		sh.Mu.Lock()
		sh.Data[i]++
		sh.Mu.Unlock()
	}
}

func BenchmarkCache_Increment_Parallel(b *testing.B) {
	cache := New(runtime.NumCPU() * 2)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := 0
		for pb.Next() {
			sh := cache.GetShard(id)
			sh.Mu.Lock()
			sh.Data[id]++
			sh.Mu.Unlock()
			id++
		}
	})
}

func BenchmarkCache_Increment_Contention(b *testing.B) {
	cache := New(1) // Один шард = максимальная конкуренция

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sh := cache.GetShard(1) // Все в один шард
			sh.Mu.Lock()
			sh.Data[1]++
			sh.Mu.Unlock()
		}
	})
}

func BenchmarkCache_Distribution(b *testing.B) {
	shards := []int{1, 2, 4, 8, 16, 32, 64, 128}

	for _, shardCount := range shards {
		b.Run(string(rune('0'+shardCount/10))+string(rune('0'+shardCount%10)), func(b *testing.B) {
			cache := New(shardCount)
			var wg sync.WaitGroup

			b.ResetTimer()
			for i := 0; i < runtime.NumCPU(); i++ {
				wg.Add(1)
				go func(start int) {
					defer wg.Done()
					for j := 0; j < b.N/runtime.NumCPU(); j++ {
						id := start*1000 + j
						sh := cache.GetShard(id)
						sh.Mu.Lock()
						sh.Data[id]++
						sh.Mu.Unlock()
					}
				}(i)
			}
			wg.Wait()
		})
	}
}
