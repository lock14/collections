package linkedlist_test

import (
	"github.com/lock14/collections/linkedlist"
	"testing"
)

func BenchmarkLinkedList_AddFront(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.AddFront(i)
	}
}

func BenchmarkLinkedList_AddBack(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.AddBack(i)
	}
}

func BenchmarkLinkedList_AddRemove(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.AddBack(i)
		l.RemoveFront()
	}
}

func BenchmarkLinkedList_Get(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	for i := 0; i < 1000; i++ {
		l.AddBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Get(i % 1000)
	}
}

func BenchmarkLinkedList_IterateAll(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	for i := 0; i < 1000; i++ {
		l.AddBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range l.All() {
		}
	}
}
