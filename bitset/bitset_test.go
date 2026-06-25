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

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *BitSet)
	}{
		{
			name: "default_capacity",
			check: func(t *testing.T, b *BitSet) {
				if got := b.Capacity(); got != 64 {
					t.Errorf("expected Capacity 64, got %d", got)
				}
				if got := b.Length(); got != 0 {
					t.Errorf("expected Length 0, got %d", got)
				}
				for i := 0; i < 64; i++ {
					if b.GetBit(i) {
						t.Errorf("expected bit %d to be unset", i)
					}
				}
			},
		},
		{
			name: "zero_bits",
			opts: []Option{NumBits(0)},
			check: func(t *testing.T, b *BitSet) {
				if got := b.Capacity(); got != 0 {
					t.Errorf("Capacity() = %v, want %v", got, 0)
				}
				b.SetBit(5) // should grow without panic
				if !b.GetBit(5) {
					t.Errorf("expected bit 5 to be set after growing from zero")
				}
			},
		},
		{
			name: "exact_word_size",
			opts: []Option{NumBits(64)},
			check: func(t *testing.T, b *BitSet) {
				if got := b.Capacity(); got != 64 {
					t.Errorf("expected Capacity 64, got %d", got)
				}
			},
		},
		{
			name: "non_exact_word_size",
			opts: []Option{NumBits(100)},
			check: func(t *testing.T, b *BitSet) {
				if got := b.Capacity(); got != 128 { // 100 bits requires 2 words
					t.Errorf("Capacity() = %v, want %v", got, 128)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := New(tc.opts...)
			tc.check(t, b)
		})
	}
}

func TestGetBit(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "get_out_of_range_returns_false",
			check: func(t *testing.T, b *BitSet) {
				if b.GetBit(1000) {
					t.Errorf("expected false for out-of-range bit")
				}
			},
		},
		{
			name: "get_out_of_range_does_not_grow",
			check: func(t *testing.T, b *BitSet) {
				sizeBefore := b.Capacity()
				b.GetBit(1000)
				if got := b.Capacity(); got != sizeBefore {
					t.Errorf("Capacity() = %v, want %v", got, sizeBefore)
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
				b.GetBit(-1)
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

func TestSetBit(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "set_single_bit",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(5)
				if !b.GetBit(5) {
					t.Errorf("expected bit 5 to be set")
				}
			},
		},
		{
			name: "set_many_bits",
			check: func(t *testing.T, b *BitSet) {
				for i := 0; i < 128; i++ {
					b.SetBit(i)
					if !b.GetBit(i) {
						t.Errorf("expected bit %d to be set", i)
					}
				}
			},
		},
		{
			name: "set_negative_panics",
			check: func(t *testing.T, b *BitSet) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic for negative index")
					}
				}()
				b.SetBit(-1)
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

func TestClearBit(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "clear_set_bit",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(5)
				b.ClearBit(5)
				if b.GetBit(5) {
					t.Errorf("expected bit 5 to be unset after Clear")
				}
			},
		},
		{
			name: "clear_updates_length",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(5)
				b.SetBit(10)
				b.ClearBit(10)
				if got := b.Length(); got != 6 {
					t.Errorf("expected Length 6 after clearing highest bit, got %d", got)
				}
			},
		},
		{
			name: "clear_only_bit_in_word_0",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(3)
				b.ClearBit(3)
				if got := b.Length(); got != 0 {
					t.Errorf("expected Length 0, got %d", got)
				}
			},
		},
		{
			name: "clear_out_of_range_does_not_grow",
			check: func(t *testing.T, b *BitSet) {
				sizeBefore := b.Capacity()
				b.ClearBit(1000)
				if got := b.Capacity(); got != sizeBefore {
					t.Errorf("Capacity() = %v, want %v", got, sizeBefore)
				}
			},
		},
		{
			name: "clear_unset_bit_is_noop",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(5)
				b.ClearBit(3)
				if !b.GetBit(5) {
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
				b.SetBit(0)
				if got := b.Length(); got != 1 {
					t.Errorf("expected 1, got %d", got)
				}
			},
		},
		{
			name: "set_bit_5",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(5)
				if got := b.Length(); got != 6 {
					t.Errorf("expected 6, got %d", got)
				}
			},
		},
		{
			name: "set_bit_63",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(63)
				if got := b.Length(); got != 64 {
					t.Errorf("expected 64, got %d", got)
				}
			},
		},
		{
			name: "set_bit_64",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(64)
				if got := b.Length(); got != 65 {
					t.Errorf("expected 65, got %d", got)
				}
			},
		},
		{
			name: "set_multiple_bits_in_word_0",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(3)
				b.SetBit(10)
				b.SetBit(7)
				if got := b.Length(); got != 11 {
					t.Errorf("expected 11, got %d", got)
				}
			},
		},
		{
			name: "set_bits_across_words",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(5)
				b.SetBit(100)
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

func TestFlip(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "flip_empty",
			check: func(t *testing.T, b *BitSet) {
				b.Flip()
				if !b.GetBit(0) || !b.GetBit(63) {
					t.Errorf("expected all bits in default capacity to be set")
				}
				if b.Length() != 64 {
					t.Errorf("expected Length 64 after flipping default size, got %d", b.Length())
				}
			},
		},
		{
			name: "flip_twice_restores",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(5)
				b.Flip()
				b.Flip()
				if got := slices.Collect(b.SetBits()); !slices.Equal(got, []int{5}) {
					t.Errorf("expected [5], got %v", got)
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

func TestFlipRange(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "flip_entire_range_does_not_expand_size",
			check: func(t *testing.T, b *BitSet) {
				b.FlipRange(0, 64)
				if got := b.Capacity(); got != 64 {
					t.Errorf("Capacity() = %v, want %v", got, 64)
				}
				if got := b.Length(); got != 64 {
					t.Errorf("unexpected length got: %d, want: %d", got, 64)
				}
			},
		},
		{
			name: "flip_partial_word",
			check: func(t *testing.T, b *BitSet) {
				b.FlipRange(2, 5) // flips 2, 3, 4
				want := []int{2, 3, 4}
				if got := slices.Collect(b.SetBits()); !slices.Equal(got, want) {
					t.Errorf("expected %v, got %v", want, got)
				}
			},
		},
		{
			name: "flip_across_words",
			check: func(t *testing.T, b *BitSet) {
				b.FlipRange(60, 68) // flips 60-63 (word 0), 64-67 (word 1)
				want := []int{60, 61, 62, 63, 64, 65, 66, 67}
				if got := slices.Collect(b.SetBits()); !slices.Equal(got, want) {
					t.Errorf("expected %v, got %v", want, got)
				}
			},
		},
		{
			name: "flip_full_middle_words",
			check: func(t *testing.T, b *BitSet) {
				b.FlipRange(60, 130) // covers word 0 (partial), word 1 (full), word 2 (partial)
				if !b.GetBit(60) || !b.GetBit(127) || !b.GetBit(129) || b.GetBit(59) || b.GetBit(130) {
					t.Errorf("flip range failed for multi-word span")
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

func TestString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "new_bit_set_all_zeros",
			check: func(t *testing.T, b *BitSet) {
				want := "0000000000000000"
				if got := b.String(); got != want {
					t.Errorf("b.String() mismatch got:\n%s\nwant:\n%s", got, want)
				}
			},
		},
		{
			name: "new_flip_all_f",
			check: func(t *testing.T, b *BitSet) {
				b.Flip()
				want := "FFFFFFFFFFFFFFFF"
				if got := b.String(); got != want {
					t.Errorf("b.String() mismatch got:\n%s\nwant:\n%s", got, want)
				}
			},
		},
		{
			name: "two_words_bottom_word_1_top_word_2",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(0)
				b.SetBit(127) // expands to 2 words
				want := "80000000000000000000000000000001"
				if got := b.String(); got != want {
					t.Errorf("b.String() mismatch got:\n%s\nwant:\n%s", got, want)
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
				wantSetBit := b.GetBit(n)
				if gotSetBit != wantSetBit {
					t.Errorf("unexpected result for bit %d, got: %v, want: %v", n, gotSetBit, wantSetBit)
				}
			}
		})
	}
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
				b.SetBit(0)
				b.SetBit(63)
				b.SetBit(64)
				b.SetBit(127)
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
				b.SetBit(1)
				b.SetBit(5)
				b.SetBit(100)
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
		b.SetBit(0)
		b.SetBit(1)
		for i := 4; i < n; i += 2 {
			b.SetBit(i)
		}
		for i := 3; (i*i) > i && (i*i) < n; i += 2 {
			if !b.GetBit(i) {
				// i is prime
				for j := i * i; j > i && j < n; j += i {
					b.SetBit(j)
				}
			}
		}
		b.FlipRange(0, n)
	}
	return b
}

func assertKV(t *testing.T, k int, ok bool, expectedK int, expectedOk bool, name string) {
	t.Helper()
	if ok != expectedOk {
		t.Errorf("%s: expected ok=%v, got %v", name, expectedOk, ok)
	} else if ok && k != expectedK {
		t.Errorf("%s: expected %v, got %v", name, expectedK, k)
	}
}

func TestBitSet_Navigable(t *testing.T) {
	t.Parallel()
	b := New()

	_, ok := b.Lower(10)
	if ok {
		t.Errorf("expected empty set to return false for Lower")
	}

	b.Add(10)
	b.Add(20)
	b.Add(30)
	b.Add(40)
	b.Add(50)

	k, ok := b.Lower(30)
	assertKV(t, k, ok, 20, true, "Lower(30)")
	k, ok = b.Lower(10)
	assertKV(t, k, ok, 0, false, "Lower(10)")

	k, ok = b.Floor(30)
	assertKV(t, k, ok, 30, true, "Floor(30)")
	k, ok = b.Floor(25)
	assertKV(t, k, ok, 20, true, "Floor(25)")
	k, ok = b.Floor(5)
	assertKV(t, k, ok, 0, false, "Floor(5)")

	k, ok = b.Ceiling(30)
	assertKV(t, k, ok, 30, true, "Ceiling(30)")
	k, ok = b.Ceiling(35)
	assertKV(t, k, ok, 40, true, "Ceiling(35)")
	k, ok = b.Ceiling(60)
	assertKV(t, k, ok, 0, false, "Ceiling(60)")

	k, ok = b.Higher(30)
	assertKV(t, k, ok, 40, true, "Higher(30)")
	k, ok = b.Higher(50)
	assertKV(t, k, ok, 0, false, "Higher(50)")
}

func TestBitSet_Sequenced(t *testing.T) {
	t.Parallel()

	assertPanic := func(name string, f func()) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s did not panic", name)
			}
		}()
		f()
	}

	b := New()
	assertPanic("First", func() { b.First() })
	assertPanic("Last", func() { b.Last() })
	assertPanic("PollFirst", func() { b.PollFirst() })
	assertPanic("PollLast", func() { b.PollLast() })

	b.Add(10)
	b.Add(20)
	b.Add(30)

	v := b.First()
	if v != 10 {
		t.Errorf("First: expected 10, got %v", v)
	}

	v = b.Last()
	if v != 30 {
		t.Errorf("Last: expected 30, got %v", v)
	}

	v = b.PollFirst()
	if v != 10 {
		t.Errorf("PollFirst: expected 10, got %v", v)
	}
	if b.Contains(10) {
		t.Errorf("PollFirst didn't remove")
	}

	v = b.PollLast()
	if v != 30 {
		t.Errorf("PollLast: expected 30, got %v", v)
	}
	if b.Contains(30) {
		t.Errorf("PollLast didn't remove")
	}

	assertPanic("AddFirst", func() { b.AddFirst(1) })
	assertPanic("AddLast", func() { b.AddLast(1) })
}

func TestBitSet_Reversed(t *testing.T) {
	t.Parallel()
	b := New()
	for i := 1; i <= 5; i++ {
		b.Add(i * 10)
	}

	var keys []int
	for k := range b.ReversedAll() {
		keys = append(keys, k)
	}
	if !slices.Equal(keys, []int{50, 40, 30, 20, 10}) {
		t.Errorf("ReversedAll: %v", keys)
	}

	// coverage for ReversedAll early exit
	for k := range b.ReversedAll() {
		_ = k
		break
	}
}

func TestBitSet_Iterators_Bounds(t *testing.T) {
	t.Parallel()
	b := New()
	for i := 1; i <= 10; i++ {
		b.Add(i)
	}

	if got := slices.Collect(b.AllFrom(5)); !slices.Equal(got, []int{5, 6, 7, 8, 9, 10}) {
		t.Errorf("AllFrom(5): %v", got)
	}
	if got := slices.Collect(b.AllTo(5)); !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Errorf("AllTo(5): %v", got)
	}
	if got := slices.Collect(b.AllBetween(3, 7)); !slices.Equal(got, []int{3, 4, 5, 6}) {
		t.Errorf("AllBetween(3, 7): %v", got)
	}

	// Early exit coverage
	for k := range b.AllBetween(1, 10) {
		_ = k
		break
	}
}
