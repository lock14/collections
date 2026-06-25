package linkedhashset_test

import (
	"github.com/lock14/collections/linkedhashset"
	"testing"
)

func BenchmarkLinkedHashSet_Add(b *testing.B) {
	s := linkedhashset.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkLinkedHashSet_Add_Preallocated(b *testing.B) {
	s := linkedhashset.New[int](linkedhashset.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkLinkedHashSet_Contains(b *testing.B) {
	s := linkedhashset.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains(i % 1000)
	}
}

func BenchmarkLinkedHashSet_RemoveElement(b *testing.B) {
	s := linkedhashset.New[int]()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.RemoveElement(i)
	}
}

func BenchmarkLinkedHashSet_RetainAll(b *testing.B) {
	s := linkedhashset.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	other := linkedhashset.New[int]()
	for i := 0; i < 500; i++ {
		other.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.RetainAll(other)
	}
}
