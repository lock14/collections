package bitset

import (
	"testing"
)

func BenchmarkBitSet_SetBit(b *testing.B) {
	b.ReportAllocs()
	bs := New(WithCapacity(10000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.SetBit(i % 10000)
	}
}

func BenchmarkBitSet_GetBit(b *testing.B) {
	b.ReportAllocs()
	bs := New(WithCapacity(10000))
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			bs.SetBit(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.GetBit(i % 10000)
	}
}

func BenchmarkBitSet_ClearBit(b *testing.B) {
	b.ReportAllocs()
	bs := New(WithCapacity(10000))
	for i := 0; i < 10000; i++ {
		bs.SetBit(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.ClearBit(i % 10000)
	}
}

func BenchmarkBitSet_FlipRange(b *testing.B) {
	b.ReportAllocs()
	bs := New(WithCapacity(100_000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs.FlipRange(0, 100_000)
	}
}

func BenchmarkBitSet_SetBits(b *testing.B) {
	b.ReportAllocs()
	bs := New(WithCapacity(100_000))
	for i := 0; i < 100_000; i++ {
		if i%2 == 0 {
			bs.SetBit(i)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range bs.SetBits() {
		}
	}
}

func BenchmarkBitSet_PrimesLessThan(b *testing.B) {
	b.ReportAllocs()
	n := 10_000
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		primesLessThan(n)
	}
}

func BenchmarkBitSet_SetBits_Sparse(b *testing.B) {
	b.ReportAllocs()
	// Allocate 1 million bits
	bs := New(WithCapacity(1_000_000))
	// Only set a few bits at the very beginning
	bs.SetBit(0)
	bs.SetBit(5)
	bs.SetBit(63)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range bs.SetBits() {
		}
	}
}
