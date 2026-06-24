package arraylist_test

import (
	"github.com/lock14/collections/arraylist"
	"testing"
)

func BenchmarkSliceWrapper_Add(b *testing.B) {
	l := arraylist.Wrap([]int{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Add(i)
	}
}

func BenchmarkSliceWrapper_Remove(b *testing.B) {
	l := arraylist.Wrap(make([]int, 0, b.N))
	for i := 0; i < b.N; i++ {
		l.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Remove()
	}
}

func BenchmarkSliceWrapper_Get(b *testing.B) {
	l := arraylist.Wrap(make([]int, 1000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Get(i % 1000)
	}
}

func BenchmarkSliceWrapper_Set(b *testing.B) {
	l := arraylist.Wrap(make([]int, 1000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Set(i%1000, i)
	}
}

func BenchmarkSliceWrapper_IterateAll(b *testing.B) {
	l := arraylist.Wrap(make([]int, 1000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range l.All() {
		}
	}
}
