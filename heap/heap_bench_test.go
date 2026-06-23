package heap_test

import (
	"github.com/lock14/collections/heap"
	"testing"
)

func BenchmarkHeap_Add(b *testing.B) {
	h := heap.Min[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(i) // Adding sequential elements is best case for min heap (O(1) amortized)
	}
}

func BenchmarkHeap_Add_Preallocated(b *testing.B) {
	h := heap.New[int](heap.Capacity[int](b.N), heap.WithComparator(heap.NaturalOrder[int]()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(i)
	}
}

func BenchmarkHeap_AddReverse(b *testing.B) {
	h := heap.Min[int]()
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		h.Add(i) // Adding in reverse order forces more siftUps (O(log n))
	}
}

func BenchmarkHeap_Remove(b *testing.B) {
	h := heap.Min[int]()
	for i := 0; i < b.N; i++ {
		h.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Remove()
	}
}

func BenchmarkHeap_AddRemove(b *testing.B) {
	h := heap.Min[int]()
	h.Add(0) // Initialize with one element to prevent empty panic
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(i)
		h.Remove()
	}
}

func BenchmarkHeap_IterateAll(b *testing.B) {
	h := heap.Min[int]()
	for i := 0; i < 1000; i++ {
		h.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _ = range h.All() {
		}
	}
}

type largeStruct struct {
	a, b, c, d, e, f, g, h int64
}

func BenchmarkHeap_StructAdd(b *testing.B) {
	cmp := func(t1, t2 largeStruct) int {
		if t1.a < t2.a {
			return -1
		} else if t1.a > t2.a {
			return 1
		}
		return 0
	}
	h := heap.New[largeStruct](heap.WithComparator(cmp))
	item := largeStruct{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item.a = int64(i)
		h.Add(item)
	}
}
