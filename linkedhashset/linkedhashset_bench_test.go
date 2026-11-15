package linkedhashset_test

import (
	"testing"

	"github.com/lock14/collections/linkedhashset"
)

const benchSize = 1000

var (
	sinkInt  int
	sinkBool bool
)

func BenchmarkLinkedHashSet_Add(b *testing.B) {
	b.ReportAllocs()
	s := linkedhashset.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkLinkedHashSet_Add_Preallocated(b *testing.B) {
	b.ReportAllocs()
	s := linkedhashset.New[int](linkedhashset.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkLinkedHashSet_AddRemoveElement(b *testing.B) {
	b.ReportAllocs()
	s := linkedhashset.New[int](linkedhashset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := i % benchSize
		s.RemoveElement(k)
		s.Add(k)
	}
}

func BenchmarkLinkedHashSet_Contains(b *testing.B) {
	b.ReportAllocs()
	s := linkedhashset.New[int](linkedhashset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBool = s.Contains(i % benchSize)
	}
}

func BenchmarkLinkedHashSet_First(b *testing.B) {
	b.ReportAllocs()
	s := linkedhashset.New[int](linkedhashset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt = s.First()
	}
}

func BenchmarkLinkedHashSet_All(b *testing.B) {
	b.ReportAllocs()
	s := linkedhashset.New[int](linkedhashset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for v := range s.All() {
			sum += v
		}
		sinkInt = sum
	}
}

func BenchmarkLinkedHashSet_Backward(b *testing.B) {
	b.ReportAllocs()
	s := linkedhashset.New[int](linkedhashset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for v := range s.Backward() {
			sum += v
		}
		sinkInt = sum
	}
}

func BenchmarkLinkedHashSet_RetainAll(b *testing.B) {
	b.ReportAllocs()
	other := linkedhashset.New[int](linkedhashset.WithCapacity(benchSize / 2))
	for i := 0; i < benchSize/2; i++ {
		other.Add(i)
	}
	s := linkedhashset.New[int](linkedhashset.WithCapacity(benchSize))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for j := 0; j < benchSize; j++ {
			s.Add(j)
		}
		b.StartTimer()
		s.RetainAll(other)
	}
}
