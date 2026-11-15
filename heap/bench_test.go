package heap_test

import (
	"testing"

	"github.com/lock14/collections/comparator"
	"github.com/lock14/collections/heap"
)

const benchSize = 1000

var sinkInt int

func BenchmarkHeap_Add(b *testing.B) {
	b.ReportAllocs()
	h := heap.Min[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(i)
	}
}

func BenchmarkHeap_Add_Preallocated(b *testing.B) {
	b.ReportAllocs()
	h := heap.New[int](heap.WithCapacity[int](b.N), heap.WithComparator(comparator.NaturalOrder[int]()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(i)
	}
}

func BenchmarkHeap_AddReverse(b *testing.B) {
	b.ReportAllocs()
	h := heap.Min[int]()
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		h.Add(i)
	}
}

func BenchmarkHeap_SteadyStateAddRemove(b *testing.B) {
	b.ReportAllocs()
	h := heap.Min[int](heap.WithCapacity[int](benchSize * 2))
	for i := 0; i < benchSize; i++ {
		h.Add(i * 2)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(benchSize + (i % benchSize))
		r = h.Remove()
	}
	sinkInt = r
}

func BenchmarkHeap_Peek(b *testing.B) {
	b.ReportAllocs()
	h := heap.Min[int](heap.WithCapacity[int](benchSize))
	for i := 0; i < benchSize; i++ {
		h.Add(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = h.Peek()
	}
	sinkInt = r
}

func BenchmarkHeap_All(b *testing.B) {
	b.ReportAllocs()
	h := heap.Min[int](heap.WithCapacity[int](benchSize))
	for i := 0; i < benchSize; i++ {
		h.Add(i)
	}
	var sum int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := range h.All() {
			sum += v
		}
	}
	sinkInt = sum
}
