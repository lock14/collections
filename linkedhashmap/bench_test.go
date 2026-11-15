package linkedhashmap_test

import (
	"testing"

	"github.com/lock14/collections/linkedhashmap"
)

const benchSize = 1000

var (
	sinkInt  int
	sinkBool bool
)

func BenchmarkLinkedHashMap_Put(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func BenchmarkLinkedHashMap_Put_Preallocated(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int](linkedhashmap.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func BenchmarkLinkedHashMap_Put_Evict(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int](linkedhashmap.WithMaxElements(benchSize))
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := benchSize; i < b.N+benchSize; i++ {
		m.Put(i, i)
	}
}

func BenchmarkLinkedHashMap_PutRemove(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int](linkedhashmap.WithCapacity(benchSize))
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

func BenchmarkLinkedHashMap_Get_InsertionOrder(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int](linkedhashmap.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt, sinkBool = m.Get(i % benchSize)
	}
}

func BenchmarkLinkedHashMap_Get_AccessOrder(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int](
		linkedhashmap.WithCapacity(benchSize),
		linkedhashmap.WithAccessOrder(),
	)
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt, sinkBool = m.Get(i % benchSize)
	}
}

func BenchmarkLinkedHashMap_First(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int](linkedhashmap.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := m.First()
		sinkInt = k + v
	}
}

func BenchmarkLinkedHashMap_All(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int](linkedhashmap.WithCapacity(benchSize))
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

func BenchmarkLinkedHashMap_Backward(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int](linkedhashmap.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for k, v := range m.Backward() {
			sum += k + v
		}
		sinkInt = sum
	}
}
