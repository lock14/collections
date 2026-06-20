package bitset

import (
	"bytes"
	"fmt"
	"slices"
	"testing"
)

// https://oeis.org/A000040
var first100Primes = []int{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
	31, 37, 41, 43, 47, 53, 59, 61, 67, 71,
	73, 79, 83, 89, 97, 101, 103, 107, 109, 113,
	127, 131, 137, 139, 149, 151, 157, 163, 167, 173,
	179, 181, 191, 193, 197, 199, 211, 223, 227, 229,
	233, 239, 241, 251, 257, 263, 269, 271, 277, 281,
	283, 293, 307, 311, 313, 317, 331, 337, 347, 349,
	353, 359, 367, 373, 379, 383, 389, 397, 401, 409,
	419, 421, 431, 433, 439, 443, 449, 457, 461, 463,
	467, 479, 487, 491, 499, 503, 509, 521, 523, 541,
}

func TestAllBitsInitializedToZero(t *testing.T) {
	t.Parallel()
	n := 128
	bitSet := New(NumBits(n))
	for i := 0; i < n; i++ {
		if bitSet.Get(i) {
			t.Errorf("excepted bit %d to be unset, but it was not", i)
		}
	}
}

func TestSetBit(t *testing.T) {
	t.Parallel()
	n := 128
	bitSet := New(NumBits(n))
	for i := 0; i < n; i++ {
		bitSet.Set(i)
		if !bitSet.Get(i) {
			t.Errorf("excepted bit %d to be set, but it was not", i)
		}
	}
}

func TestString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		bitSetInitFunc func() *BitSet
		want           string
	}{
		{
			name: "new_bit_set_all_zeros",
			bitSetInitFunc: func() *BitSet {
				return New()
			},
			want: "0000000000000000",
		},
		{
			name: "new_flip_all_f",
			bitSetInitFunc: func() *BitSet {
				b := New()
				b.Flip()
				return b
			},
			want: "FFFFFFFFFFFFFFFF",
		},
		{
			name: "two_words_bottom_word_1_top_word_2",
			bitSetInitFunc: func() *BitSet {
				b := New(NumBits(128))
				b.Set(0)
				b.Set(127)
				return b
			},
			want: "80000000000000000000000000000001",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := tc.bitSetInitFunc()
			got := b.String()
			if got != tc.want {
				t.Errorf("b.String() mismatch got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

func TestFromBytesToBytes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "empty_slice",
			input: []byte{},
			want:  []byte{},
		},
		{
			name:  "one_byte",
			input: []byte{0xFF},
			want:  []byte{0xFF},
		},
		{
			name:  "eight_bytes",
			input: []byte{0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF},
			want:  []byte{0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF},
		},
		{
			name:  "twelve_bytes",
			input: []byte{0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF},
			want:  []byte{0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF},
		},
		{
			name:  "sixteen_bytes",
			input: []byte{0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF},
			want:  []byte{0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := FromBytes(tc.input)
			got := b.ToBytes()
			if !bytes.Equal(got, tc.want) {
				t.Errorf("unexpected result got:\n%v\nwant:\n%v", got, tc.want)
			}
			for n := 0; n < b.Length(); n++ {
				gotSetBit := (got[n/8] & (1 << (n % 8))) != 0
				wantSetBit := b.Get(n)
				if gotSetBit != wantSetBit {
					t.Errorf("unexpected result for bit %d, got: %v, want: %v", n, gotSetBit, wantSetBit)
				}
			}
		})
	}
}

func TestFlipRange(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		start int
		end   int
		want  int
	}{
		{
			name:  "flip_entire_range_does_not_expand_size",
			start: 0,
			end:   64,
			want:  64,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := New()
			b.FlipRange(tc.start, tc.end)
			got := b.Size()
			if got != tc.want {
				t.Errorf("unexpected result got: %d, want: %d", got, tc.want)
			}
		})
	}
}

func TestBitSetPrimeGen(t *testing.T) {
	t.Parallel()
	// a prime sieve is a good gamut test of a BitSet
	cases := make([]struct {
		name     string
		lessThan int
		want     []int
	}, 0, len(first100Primes))
	for i := 0; i < len(first100Primes); i++ {
		lessThan := first100Primes[i] + 1
		cases = append(cases, struct {
			name     string
			lessThan int
			want     []int
		}{
			name:     fmt.Sprintf("primes_less_than_%d", lessThan),
			lessThan: lessThan,
			want:     first100Primes[:i+1],
		})
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := primesLessThan(tc.lessThan)
			primes := make([]int, 0, len(tc.want))
			for n := range b.SetBits() {
				primes = append(primes, n)
			}
			if !slices.Equal(primes, tc.want) {
				t.Errorf("unexpected result got: %v, want: %v", primes, tc.want)
			}
		})
	}
}

func primesLessThan(n int) *BitSet {
	b := New(NumBits(n))
	if n > 2 {
		b.Set(0)
		b.Set(1)
		for i := 4; i < n; i += 2 {
			b.Set(i)
		}
		for i := 3; (i*i) > i && (i*i) < n; i += 2 {
			if !b.Get(i) {
				// i is prime
				for j := i * i; j > i && j < n; j += i {
					b.Set(j)
				}
			}
		}
		b.FlipRange(0, n)
	}
	return b
}

func TestLength(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "empty_bitset",
			check: func(t *testing.T, b *BitSet) {
				if got := b.Length(); got != 0 {
					t.Errorf("expected 0, got %d", got)
				}
			},
		},
		{
			name: "set_bit_0",
			check: func(t *testing.T, b *BitSet) {
				b.Set(0)
				if got := b.Length(); got != 1 {
					t.Errorf("expected 1, got %d", got)
				}
			},
		},
		{
			name: "set_bit_5",
			check: func(t *testing.T, b *BitSet) {
				b.Set(5)
				if got := b.Length(); got != 6 {
					t.Errorf("expected 6, got %d", got)
				}
			},
		},
		{
			name: "set_bit_63",
			check: func(t *testing.T, b *BitSet) {
				b.Set(63)
				if got := b.Length(); got != 64 {
					t.Errorf("expected 64, got %d", got)
				}
			},
		},
		{
			name: "set_bit_64",
			check: func(t *testing.T, b *BitSet) {
				b.Set(64)
				if got := b.Length(); got != 65 {
					t.Errorf("expected 65, got %d", got)
				}
			},
		},
		{
			name: "set_multiple_bits_in_word_0",
			check: func(t *testing.T, b *BitSet) {
				b.Set(3)
				b.Set(10)
				b.Set(7)
				if got := b.Length(); got != 11 {
					t.Errorf("expected 11, got %d", got)
				}
			},
		},
		{
			name: "set_bits_across_words",
			check: func(t *testing.T, b *BitSet) {
				b.Set(5)
				b.Set(100)
				if got := b.Length(); got != 101 {
					t.Errorf("expected 101, got %d", got)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := New()
			tc.check(t, b)
		})
	}
}

func TestClear(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "clear_set_bit",
			check: func(t *testing.T, b *BitSet) {
				b.Set(5)
				b.Clear(5)
				if b.Get(5) {
					t.Errorf("expected bit 5 to be unset after Clear")
				}
			},
		},
		{
			name: "clear_updates_length",
			check: func(t *testing.T, b *BitSet) {
				b.Set(5)
				b.Set(10)
				b.Clear(10)
				if got := b.Length(); got != 6 {
					t.Errorf("expected Length 6 after clearing highest bit, got %d", got)
				}
			},
		},
		{
			name: "clear_only_bit_in_word_0",
			check: func(t *testing.T, b *BitSet) {
				b.Set(3)
				b.Clear(3)
				if got := b.Length(); got != 0 {
					t.Errorf("expected Length 0, got %d", got)
				}
			},
		},
		{
			name: "clear_out_of_range_does_not_grow",
			check: func(t *testing.T, b *BitSet) {
				sizeBefore := b.Size()
				b.Clear(1000)
				if got := b.Size(); got != sizeBefore {
					t.Errorf("Clear on out-of-range bit grew the BitSet from %d to %d", sizeBefore, got)
				}
			},
		},
		{
			name: "clear_unset_bit_is_noop",
			check: func(t *testing.T, b *BitSet) {
				b.Set(5)
				b.Clear(3)
				if !b.Get(5) {
					t.Errorf("clearing an unset bit should not affect other bits")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := New()
			tc.check(t, b)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "get_out_of_range_returns_false",
			check: func(t *testing.T, b *BitSet) {
				if b.Get(1000) {
					t.Errorf("expected false for out-of-range bit")
				}
			},
		},
		{
			name: "get_out_of_range_does_not_grow",
			check: func(t *testing.T, b *BitSet) {
				sizeBefore := b.Size()
				b.Get(1000)
				if got := b.Size(); got != sizeBefore {
					t.Errorf("Get on out-of-range bit grew the BitSet from %d to %d", sizeBefore, got)
				}
			},
		},
		{
			name: "get_negative_panics",
			check: func(t *testing.T, b *BitSet) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic for negative index")
					}
				}()
				b.Get(-1)
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := New()
			tc.check(t, b)
		})
	}
}

func TestSet_Negative_Panics(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for negative index")
		}
	}()
	b := New()
	b.Set(-1)
}

func TestSetBits(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "no_bits_set",
			check: func(t *testing.T, b *BitSet) {
				got := slices.Collect(b.SetBits())
				if len(got) != 0 {
					t.Errorf("expected no set bits, got: %v", got)
				}
			},
		},
		{
			name: "sparse_bits",
			check: func(t *testing.T, b *BitSet) {
				b.Set(0)
				b.Set(63)
				b.Set(64)
				b.Set(127)
				want := []int{0, 63, 64, 127}
				got := slices.Collect(b.SetBits())
				if !slices.Equal(got, want) {
					t.Errorf("wrong set bits, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "early_break",
			check: func(t *testing.T, b *BitSet) {
				b.Set(1)
				b.Set(5)
				b.Set(100)
				got := []int{}
				for n := range b.SetBits() {
					got = append(got, n)
					if n == 5 {
						break
					}
				}
				want := []int{1, 5}
				if !slices.Equal(got, want) {
					t.Errorf("wrong set bits, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := New(NumBits(128))
			tc.check(t, b)
		})
	}
}

func TestNew_ZeroBits(t *testing.T) {
	t.Parallel()
	b := New(NumBits(0))
	if got := b.Size(); got != 0 {
		t.Errorf("expected Size 0, got %d", got)
	}
	b.Set(5) // should grow without panic
	if !b.Get(5) {
		t.Errorf("expected bit 5 to be set after growing from zero")
	}
}

