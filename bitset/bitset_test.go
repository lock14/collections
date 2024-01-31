package bitset

import (
	"github.com/lock14/collections/util"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// https://oeis.org/A000040
var first100Primes = []uint{
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
	bitSet := New(NumBits(uint(n)))
	for i := 0; i < n; i++ {
		if bitSet.Get(uint(i)) {
			t.Fatalf("excepted bit %d to be unset, but it was not", i)
		}
	}
}

func TestSetBit(t *testing.T) {
	t.Parallel()
	n := 128
	bitSet := New(NumBits(uint(n)))
	for i := 0; i < n; i++ {
		bitSet.Set(uint(i))
		if !bitSet.Get(uint(i)) {
			t.Fatalf("excepted bit %d to be set, but it was not", i)
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

func TestBitSetPrimeGen(t *testing.T) {
	t.Parallel()
	// a prime sieve is a good gamut test of a BitSet
	b := primesLessThan(first100Primes[len(first100Primes)-1] + 1)
	primes := make([]uint, 0, 100)
	for i := b.SetBits(); !i.Empty(); util.MustDo(i.PopFront()) {
		n := *util.Must(i.Front())
		primes = append(primes, n)
	}
	if diff := cmp.Diff(primes, first100Primes); diff != "" {
		t.Fatalf("wrong values: %s", diff)
	}
}

func primesLessThan(n uint) *BitSet {
	b := New(NumBits(n))
	if n > 2 {
		b.Set(0)
		b.Set(1)
		for i := uint(4); i < n; i += 2 {
			b.Set(i)
		}
		for i := uint(3); (i*i) > i && (i*i) < n; i += 2 {
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
