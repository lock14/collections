package arraydeque_test

import (
	"testing"

	"github.com/lock14/collections/arraydeque"
)

const benchSize = 1000

var sinkInt int

func BenchmarkArrayDeque_AddFront(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddFront(i)
	}
}

func BenchmarkArrayDeque_AddBack(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddBack(i)
	}
}

func BenchmarkArrayDeque_AddFront_Preallocated(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int](arraydeque.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddFront(i)
	}
}

func BenchmarkArrayDeque_AddBack_Preallocated(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int](arraydeque.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddBack(i)
	}
}

func BenchmarkArrayDeque_SteadyStatePushPop(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int](arraydeque.WithCapacity(benchSize * 2))
	for i := 0; i < benchSize; i++ {
		d.Push(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Push(i)
		r = d.Pop()
	}
	sinkInt = r
}

func BenchmarkArrayDeque_AddRemove(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int](arraydeque.WithCapacity(benchSize * 2))
	for i := 0; i < benchSize; i++ {
		d.AddBack(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddBack(i)
		r = d.RemoveFront()
	}
	sinkInt = r
}

func BenchmarkArrayDeque_PeekFront(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int]()
	for i := 0; i < benchSize; i++ {
		d.AddBack(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = d.PeekFront()
	}
	sinkInt = r
}

func BenchmarkArrayDeque_PeekBack(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int]()
	for i := 0; i < benchSize; i++ {
		d.AddBack(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = d.PeekBack()
	}
	sinkInt = r
}

func BenchmarkArrayDeque_All(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int]()
	for i := 0; i < benchSize; i++ {
		d.AddBack(i)
	}
	var sum int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := range d.All() {
			sum += v
		}
	}
	sinkInt = sum
}

func BenchmarkArrayDeque_Backward(b *testing.B) {
	b.ReportAllocs()
	d := arraydeque.New[int]()
	for i := 0; i < benchSize; i++ {
		d.AddBack(i)
	}
	var sum int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := range d.Backward() {
			sum += v
		}
	}
	sinkInt = sum
}
