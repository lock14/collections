package hashmap_test

import (
	"testing"

	"github.com/lock14/collections/hashmap"
)

const benchSize = 1000

var (
	sinkInt  int
	sinkBool bool
)

func BenchmarkHashMap_Put(b *testing.B) {
	b.ReportAllocs()
	m := hashmap.New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func BenchmarkHashMap_Put_Preallocated(b *testing.B) {
	b.ReportAllocs()
	m := hashmap.New[int, int](hashmap.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func BenchmarkHashMap_PutRemove(b *testing.B) {
	b.ReportAllocs()
	m := hashmap.New[int, int](hashmap.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := i % benchSize
		m.Remove(k)
		m.Put(k, i)
	}
}

func BenchmarkHashMap_Get(b *testing.B) {
	b.ReportAllocs()
	m := hashmap.New[int, int](hashmap.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt, sinkBool = m.Get(i % benchSize)
	}
}

func BenchmarkHashMap_ContainsKey(b *testing.B) {
	b.ReportAllocs()
	m := hashmap.New[int, int](hashmap.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBool = m.ContainsKey(i % benchSize)
	}
}

func BenchmarkHashMap_All(b *testing.B) {
	b.ReportAllocs()
	m := hashmap.New[int, int](hashmap.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for k, v := range m.All() {
			sum += k + v
		}
		sinkInt = sum
	}
}
