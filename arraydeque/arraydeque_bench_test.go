package arraydeque

import (
	"testing"
)

func BenchmarkArrayDeque_AddFront(b *testing.B) {
	d := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddFront(i)
	}
}

func BenchmarkArrayDeque_AddBack(b *testing.B) {
	d := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddBack(i)
	}
}

func BenchmarkArrayDeque_PushPop(b *testing.B) {
	d := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Push(i)
		d.Pop()
	}
}

func BenchmarkArrayDeque_AddRemove(b *testing.B) {
	d := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddBack(i)
		d.RemoveFront()
	}
}

func BenchmarkArrayDeque_AddFront_Preallocated(b *testing.B) {
	d := New[int](func(c *config) { c.capacity = b.N })
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddFront(i)
	}
}

func BenchmarkArrayDeque_AddBack_Preallocated(b *testing.B) {
	d := New[int](func(c *config) { c.capacity = b.N })
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddBack(i)
	}
}

func BenchmarkArrayDeque_IterateAll(b *testing.B) {
	d := New[int]()
	for i := 0; i < 1000; i++ {
		d.AddBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range d.All() {
		}
	}
}

func BenchmarkArrayDeque_IterateBackward(b *testing.B) {
	d := New[int]()
	for i := 0; i < 1000; i++ {
		d.AddBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range d.Backward() {
		}
	}
}

type largeStruct struct {
	a, b, c, d, e, f, g, h int64
}

func BenchmarkArrayDeque_StructAddBack(b *testing.B) {
	d := New[largeStruct]()
	item := largeStruct{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.AddBack(item)
	}
}

func BenchmarkArrayDeque_ClearAndReuse_NewEachTime(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := New[int]()
		for j := 0; j < 100; j++ {
			d.AddBack(j)
		}
		d.Clear() // Simulate clearing out the old one before tossing it
	}
}

func BenchmarkArrayDeque_ClearAndReuse_Existing(b *testing.B) {
	d := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			d.AddBack(j)
		}
		d.Clear() // Actually reuse it next iteration
	}
}
