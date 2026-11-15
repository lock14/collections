package hashset_test

import (
	"testing"

	"github.com/lock14/collections/hashset"
)

const benchSize = 1000

var (
	sinkInt  int
	sinkBool bool
)

func BenchmarkHashSet_Add(b *testing.B) {
	b.ReportAllocs()
	s := hashset.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkHashSet_Add_Preallocated(b *testing.B) {
	b.ReportAllocs()
	s := hashset.New[int](hashset.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkHashSet_AddRemove(b *testing.B) {
	b.ReportAllocs()
	s := hashset.New[int](hashset.WithCapacity(benchSize))
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

func BenchmarkHashSet_Contains(b *testing.B) {
	b.ReportAllocs()
	s := hashset.New[int](hashset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBool = s.Contains(i % benchSize)
	}
}

func BenchmarkHashSet_All(b *testing.B) {
	b.ReportAllocs()
	s := hashset.New[int](hashset.WithCapacity(benchSize))
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

func BenchmarkHashSet_RetainAll(b *testing.B) {
	b.ReportAllocs()
	other := hashset.New[int](hashset.WithCapacity(benchSize / 2))
	for i := 0; i < benchSize/2; i++ {
		other.Add(i)
	}
	s := hashset.New[int](hashset.WithCapacity(benchSize))
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
