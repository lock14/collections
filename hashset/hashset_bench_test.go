package hashset_test

import (
	"github.com/lock14/collections/hashset"
	"testing"
)

func BenchmarkHashSet_Add(b *testing.B) {
	s := hashset.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkHashSet_Add_Preallocated(b *testing.B) {
	s := hashset.New[int](hashset.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkHashSet_Contains(b *testing.B) {
	s := hashset.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains(i % 1000)
	}
}

func BenchmarkHashSet_RemoveElement(b *testing.B) {
	s := hashset.New[int]()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.RemoveElement(i)
	}
}

func BenchmarkHashSet_RetainAll(b *testing.B) {
	s := hashset.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	other := hashset.New[int]()
	for i := 0; i < 500; i++ {
		other.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.RetainAll(other)
	}
}
