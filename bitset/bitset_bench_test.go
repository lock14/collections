package bitset

import (
	"testing"
)

func BenchmarkSet(b *testing.B) {
	bs := New(NumBits(10000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.Set(i % 10000)
	}
}

func BenchmarkGet(b *testing.B) {
	bs := New(NumBits(10000))
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			bs.Set(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.Get(i % 10000)
	}
}

func BenchmarkClear(b *testing.B) {
	bs := New(NumBits(10000))
	for i := 0; i < 10000; i++ {
		bs.Set(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.Clear(i % 10000)
	}
}

func BenchmarkFlipRange(b *testing.B) {
	bs := New(NumBits(100_000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.FlipRange(0, 100_000)
	}
}

func BenchmarkSetBits(b *testing.B) {
	bs := New(NumBits(100_000))
	for i := 0; i < 100_000; i++ {
		if i%2 == 0 {
			bs.Set(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range bs.SetBits() {
		}
	}
}

func BenchmarkPrimesLessThan(b *testing.B) {
	n := 10_000
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		primesLessThan(n)
	}
}
