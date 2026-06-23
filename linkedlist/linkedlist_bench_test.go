package linked_list_test

import (
	"github.com/lock14/collections/linkedlist"
	"testing"
)

func BenchmarkLinkedList_AddFront(b *testing.B) {
	l := linked_list.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.AddFront(i)
	}
}

func BenchmarkLinkedList_AddBack(b *testing.B) {
	l := linked_list.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.AddBack(i)
	}
}

func BenchmarkLinkedList_RemoveFront(b *testing.B) {
	l := linked_list.New[int]()
	for i := 0; i < b.N; i++ {
		l.AddBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.RemoveFront()
	}
}

func BenchmarkLinkedList_RemoveBack(b *testing.B) {
	l := linked_list.New[int]()
	for i := 0; i < b.N; i++ {
		l.AddBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.RemoveBack()
	}
}

func BenchmarkLinkedList_Get(b *testing.B) {
	l := linked_list.New[int]()
	for i := 0; i < 1000; i++ {
		l.AddBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Get(i % 1000)
	}
}

func BenchmarkLinkedList_IterateAll(b *testing.B) {
	l := linked_list.New[int]()
	for i := 0; i < 1000; i++ {
		l.AddBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _ = range l.All() {
		}
	}
}
