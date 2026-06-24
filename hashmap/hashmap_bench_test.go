package hashmap_test

import (
	"github.com/lock14/collections/hashmap"
	"testing"
)

func BenchmarkHashMap_Put(b *testing.B) {
	hm := hashmap.New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hm.Put(i, i)
	}
}

func BenchmarkHashMap_Put_Preallocated(b *testing.B) {
	hm := hashmap.New[int, int](hashmap.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hm.Put(i, i)
	}
}

func BenchmarkHashMap_Get(b *testing.B) {
	hm := hashmap.New[int, int]()
	for i := 0; i < 1000; i++ {
		hm.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hm.Get(i % 1000)
	}
}

func BenchmarkHashMap_Remove(b *testing.B) {
	hm := hashmap.New[int, int]()
	for i := 0; i < b.N; i++ {
		hm.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hm.Remove(i)
	}
}

func BenchmarkHashMap_IterateAll(b *testing.B) {
	hm := hashmap.New[int, int]()
	for i := 0; i < 1000; i++ {
		hm.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range hm.All() {
		}
	}
}
