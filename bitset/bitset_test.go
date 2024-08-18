package bitset

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("b.String() mismatch (-got, +want):\n%s", diff)
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
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want):\n%s", diff)
			}
			for n := 0; n < b.Length(); n++ {
				gotSetBit := (got[n/8] & (1 << (n % 8))) != 0
				wantSetBit := b.Get(n)
				if diff := cmp.Diff(gotSetBit, wantSetBit); diff != "" {
					t.Errorf("unexpected result (-got, +want):\n%s", diff)
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
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want):\n%s", diff)
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
			if diff := cmp.Diff(primes, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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
