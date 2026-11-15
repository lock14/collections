package arraylist_test

import (
	"testing"

	"github.com/lock14/collections/arraylist"
)

const benchSize = 1000

var sinkInt int

func BenchmarkSliceWrapper_Add(b *testing.B) {
	b.ReportAllocs()
	l := arraylist.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Add(i)
	}
}

func BenchmarkSliceWrapper_Add_Preallocated(b *testing.B) {
	b.ReportAllocs()
	l := arraylist.New[int](arraylist.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Add(i)
	}
}

func BenchmarkSliceWrapper_AddRemove(b *testing.B) {
	b.ReportAllocs()
	l := arraylist.New[int](arraylist.WithCapacity(benchSize * 2))
	for i := 0; i < benchSize; i++ {
		l.Add(i)
	}
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Add(i)
		r = l.Remove()
	}
	sinkInt = r
}

func BenchmarkSliceWrapper_Get(b *testing.B) {
	b.ReportAllocs()
	l := arraylist.Wrap(make([]int, benchSize))
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = l.Get(i % benchSize)
	}
	sinkInt = r
}

func BenchmarkSliceWrapper_Peek(b *testing.B) {
	b.ReportAllocs()
	l := arraylist.Wrap(make([]int, benchSize))
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = l.Peek()
	}
	sinkInt = r
}

func BenchmarkSliceWrapper_Set(b *testing.B) {
	b.ReportAllocs()
	l := arraylist.Wrap(make([]int, benchSize))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Set(i%benchSize, i)
	}
}

func BenchmarkSliceWrapper_All(b *testing.B) {
	b.ReportAllocs()
	l := arraylist.Wrap(make([]int, benchSize))
	var sum int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for v := range l.All() {
			sum += v
		}
	}
	sinkInt = sum
}
