package bitset_test

import (
	"testing"

	"github.com/lock14/collections/bitset"
)

const benchSize = 10_000

var (
	sinkBool   bool
	sinkInt    int
	sinkBitSet *bitset.BitSet
)

func sievePrimes(n int) *bitset.BitSet {
	s := bitset.New(bitset.WithCapacity(n))
	if n > 2 {
		s.SetBit(0)
		s.SetBit(1)
		for i := 4; i < n; i += 2 {
			s.SetBit(i)
		}
		for i := 3; i*i < n; i += 2 {
			if !s.GetBit(i) {
				for j := i * i; j < n; j += i {
					s.SetBit(j)
				}
			}
		}
		s.FlipRange(0, n)
	}
	return s
}

func BenchmarkBitSet_SetClearBit(b *testing.B) {
	b.ReportAllocs()
	s := bitset.New(bitset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i += 2 {
		s.SetBit(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % benchSize
		s.SetBit(idx)
		s.ClearBit(idx)
	}
}

func BenchmarkBitSet_Grow(b *testing.B) {
	b.ReportAllocs()
	var s *bitset.BitSet
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = bitset.New(bitset.WithCapacity(0))
		for bit := 0; bit < benchSize; bit += 64 {
			s.SetBit(bit)
		}
	}
	sinkBitSet = s
}

func BenchmarkBitSet_GetBit(b *testing.B) {
	b.ReportAllocs()
	s := bitset.New(bitset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i += 2 {
		s.SetBit(i)
	}
	var r bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = s.GetBit(i % benchSize)
	}
	sinkBool = r
}

func BenchmarkBitSet_Floor(b *testing.B) {
	b.ReportAllocs()
	s := bitset.New(bitset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i += 10 {
		s.SetBit(i)
	}
	var v int
	var ok bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, ok = s.Floor(i % benchSize)
	}
	sinkInt = v
	sinkBool = ok
}

func BenchmarkBitSet_Higher(b *testing.B) {
	b.ReportAllocs()
	s := bitset.New(bitset.WithCapacity(benchSize))
	for i := 0; i < benchSize; i += 10 {
		s.SetBit(i)
	}
	var v int
	var ok bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, ok = s.Higher(i % benchSize)
	}
	sinkInt = v
	sinkBool = ok
}

func BenchmarkBitSet_FlipRange(b *testing.B) {
	b.ReportAllocs()
	s := bitset.New(bitset.WithCapacity(100_000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.FlipRange(0, 100_000)
	}
}

func BenchmarkBitSet_All(b *testing.B) {
	b.ReportAllocs()
	s := bitset.New(bitset.WithCapacity(100_000))
	for i := 0; i < 100_000; i += 2 {
		s.SetBit(i)
	}
	var sum int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for bit := range s.All() {
			sum += bit
		}
	}
	sinkInt = sum
}

func BenchmarkBitSet_SetBits_Sparse(b *testing.B) {
	b.ReportAllocs()
	s := bitset.New(bitset.WithCapacity(1_000_000))
	s.SetBit(0)
	s.SetBit(5)
	s.SetBit(63)

	var sum int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for bit := range s.SetBits() {
			sum += bit
		}
	}
	sinkInt = sum
}

func BenchmarkBitSet_Between(b *testing.B) {
	b.ReportAllocs()
	s := bitset.New(bitset.WithCapacity(100_000))
	for i := 0; i < 100_000; i += 4 {
		s.SetBit(i)
	}
	var sum int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for bit := range s.Between(10_000, 10_500) {
			sum += bit
		}
	}
	sinkInt = sum
}

func BenchmarkBitSet_PrimesLessThan(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var s *bitset.BitSet
	for i := 0; i < b.N; i++ {
		s = sievePrimes(benchSize)
	}
	sinkBitSet = s
}
