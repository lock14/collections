package linkedhashmap_test

import (
	"github.com/lock14/collections/linkedhashmap"
	"testing"
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
	m := linkedhashmap.New[int, int](linkedhashmap.WithMaxElements(1000))
	for i := 0; i < 1000; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 1000; i < b.N+1000; i++ {
		m.Put(i, i)
	}
}

func BenchmarkLinkedHashMap_Get_InsertionOrder(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int]()
	for i := 0; i < 1000; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(i % 1000)
	}
}

func BenchmarkLinkedHashMap_Get_AccessOrder(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int](linkedhashmap.WithAccessOrder())
	for i := 0; i < 1000; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(i % 1000)
	}
}

func BenchmarkLinkedHashMap_AddRemove(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
		m.Remove(i)
	}
}

func BenchmarkLinkedHashMap_IterateAll(b *testing.B) {
	b.ReportAllocs()
	m := linkedhashmap.New[int, int]()
	for i := 0; i < 1000; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range m.All() {
		}
	}
}
