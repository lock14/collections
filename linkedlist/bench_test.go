package linkedlist_test

import (
	"testing"

	"github.com/lock14/collections/linkedlist"
)

const benchSize = 1000

var sinkInt int

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
	for i := 0; i < benchSize; i++ {
		l.AddBack(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.AddBack(i)
		r = l.RemoveFront()
	}
	sinkInt = r
}

func BenchmarkLinkedList_PeekFront(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	for i := 0; i < benchSize; i++ {
		l.AddBack(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = l.PeekFront()
	}
	sinkInt = r
}

func BenchmarkLinkedList_PeekBack(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	for i := 0; i < benchSize; i++ {
		l.AddBack(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = l.PeekBack()
	}
	sinkInt = r
}

func BenchmarkLinkedList_Get(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	for i := 0; i < benchSize; i++ {
		l.AddBack(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = l.Get(i % benchSize)
	}
	sinkInt = r
}

func BenchmarkLinkedList_All(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	for i := 0; i < benchSize; i++ {
		l.AddBack(i)
	}
	var sum int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := range l.All() {
			sum += v
		}
	}
	sinkInt = sum
}

func BenchmarkLinkedList_Backward(b *testing.B) {
	b.ReportAllocs()
	l := linkedlist.New[int]()
	for i := 0; i < benchSize; i++ {
		l.AddBack(i)
	}
	var sum int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := range l.Backward() {
			sum += v
		}
	}
	sinkInt = sum
}
