package bitset

import (
	"bytes"
	"fmt"
	"slices"
	"testing"

	"github.com/lock14/collections/arraylist"
	"github.com/lock14/collections/hashset"
	linkedlist "github.com/lock14/collections/linkedlist"
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
			opts: []Option{WithCapacity(0)},
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
			opts: []Option{WithCapacity(64)},
			check: func(t *testing.T, b *BitSet) {
				if got := b.Capacity(); got != 64 {
					t.Errorf("expected Capacity 64, got %d", got)
				}
			},
		},
		{
			name: "non_exact_word_size",
			opts: []Option{WithCapacity(100)},
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
			name: "new_bit_set_empty",
			check: func(t *testing.T, b *BitSet) {
				want := "[]"
				if got := b.String(); got != want {
					t.Errorf("b.String() mismatch got:\n%s\nwant:\n%s", got, want)
				}
			},
		},
		{
			name: "single_and_multiple_elements",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(0)
				b.SetBit(4)
				want := "[0 4]"
				if got := b.String(); got != want {
					t.Errorf("b.String() mismatch got:\n%s\nwant:\n%s", got, want)
				}
			},
		},
		{
			name: "two_words_elements",
			check: func(t *testing.T, b *BitSet) {
				b.SetBit(0)
				b.SetBit(127) // expands to 2 words
				want := "[0 127]"
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
			b := New(WithCapacity(128))
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
	b := New(WithCapacity(n))
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

func TestBitSet_Backward(t *testing.T) {
	t.Parallel()
	b := New()
	for i := 1; i <= 5; i++ {
		b.Add(i * 10)
	}

	var keys []int
	for k := range b.Backward() {
		keys = append(keys, k)
	}
	if !slices.Equal(keys, []int{50, 40, 30, 20, 10}) {
		t.Errorf("Backward: %v", keys)
	}

	// coverage for Backward early exit
	for k := range b.Backward() {
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

	if got := slices.Collect(b.From(5)); !slices.Equal(got, []int{5, 6, 7, 8, 9, 10}) {
		t.Errorf("From(5): %v", got)
	}
	if got := slices.Collect(b.To(5)); !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Errorf("To(5): %v", got)
	}
	if got := slices.Collect(b.Between(3, 7)); !slices.Equal(got, []int{3, 4, 5, 6}) {
		t.Errorf("Between(3, 7): %v", got)
	}

	// Early exit coverage
	for k := range b.Between(1, 10) {
		_ = k
		break
	}
}

func TestBitSet_MutableSetTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_MutableSet",
			run: func(t *testing.T) {

				t.Parallel()
				b := New(WithCapacity(100))
				b.Add(10)
				b.Add(20)

				if !b.Contains(10) {
					t.Errorf("expected to contain 10")
				}
				if !b.Contains(20) {
					t.Errorf("expected to contain 20")
				}

				b.RemoveElement(10)
				if b.Contains(10) {
					t.Errorf("expected 10 to be removed")
				}

				if b.Empty() {
					t.Errorf("expected not to be empty")
				}

				b.RemoveElement(20)
				if !b.Empty() {
					t.Errorf("expected to be empty")
				}

				b.Add(5)
				b.Add(15)
				val := b.Remove()
				if val != 5 && val != 15 {
					t.Errorf("Remove returned unexpected value: %v", val)
				}

				b.Clear()
				if !b.Empty() {
					t.Errorf("expected empty after Clear")
				}

				b.AddAll(slices.Values([]int{1, 2, 3}))
				if !b.Contains(1) || !b.Contains(2) || !b.Contains(3) {
					t.Errorf("AddAll failed")
				}

				other := New(WithCapacity(100))
				other.Add(2)
				other.Add(3)
				other.Add(4)

				b.RemoveAll(other)
				if b.Contains(2) || b.Contains(3) || !b.Contains(1) {
					t.Errorf("RemoveAll failed")
				}

				b.Clear()
				b.AddAll(slices.Values([]int{1, 2, 3}))
				b.RetainAll(other) // retains 2, 3
				if b.Contains(1) || !b.Contains(2) || !b.Contains(3) {
					t.Errorf("RetainAll failed")
				}

				// test RetainAll when b has more words than other
				b.Clear()
				b.Add(100) // maxWordInUse = 2
				other.Clear()
				other.Add(2) // maxWordInUse = 1
				b.RetainAll(other)
				if b.Contains(100) {
					t.Errorf("RetainAll did not clear upper words")
				}

				b.Clear()
				b.AddAll(slices.Values([]int{2, 3}))
				other.Clear()
				other.AddAll(slices.Values([]int{2, 3, 4}))
				if b.ContainsAll(other) {
					t.Errorf("expected not to contain all")
				}

				// test ContainsAll when other has more words than b
				other.Clear()
				other.AddAll(slices.Values([]int{2, 3, 100}))
				if b.ContainsAll(other) {
					t.Errorf("expected not to contain all")
				}
				// test ContainsAll when other has more words than b but empty upper words
				other.Clear()
				other.Add(2)
				other.ensureSize(5)
				other.maxWordInUse = 6
				if !b.ContainsAll(other) {
					t.Errorf("expected to contain all even with empty upper words")
				}

				other.Clear()
				other.AddAll(slices.Values([]int{2, 3}))
				if !other.ContainsAll(b) {
					t.Errorf("ContainsAll failed")
				}

				collected := slices.Collect(b.All())
				if !slices.Equal(collected, []int{2, 3}) {
					t.Errorf("All iterator failed")
				}

				// Test generic Collection branch (using arraylist.Wrap)
				list := arraylist.Wrap([]int{2, 3, 4})
				b.Clear()
				b.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
				b.RemoveAll(list)
				if b.Contains(2) || b.Contains(3) || b.Contains(4) || !b.Contains(1) || !b.Contains(5) {
					t.Errorf("RemoveAll with generic collection failed")
				}

				b.Clear()
				b.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
				b.RetainAll(list)
				if b.Contains(1) || b.Contains(5) || !b.Contains(2) || !b.Contains(3) || !b.Contains(4) {
					t.Errorf("RetainAll with generic collection failed")
				}

				if !b.ContainsAll(arraylist.Wrap([]int{2, 3, 4})) {
					t.Errorf("ContainsAll with generic collection should be true")
				}
				if b.ContainsAll(arraylist.Wrap([]int{2, 3, 4, 99})) {
					t.Errorf("ContainsAll with generic collection should be false")
				}

				// Test generic Set branch (using hashset.New)
				hset := hashset.New[int]()
				hset.Add(2)
				hset.Add(3)
				hset.Add(4)

				b.Clear()
				b.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
				b.RemoveAll(hset)
				if b.Contains(2) || b.Contains(3) || b.Contains(4) || !b.Contains(1) || !b.Contains(5) {
					t.Errorf("RemoveAll with Set failed")
				}

				b.Clear()
				b.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
				b.RetainAll(hset)
				if b.Contains(1) || b.Contains(5) || !b.Contains(2) || !b.Contains(3) || !b.Contains(4) {
					t.Errorf("RetainAll with Set failed")
				}

				// Test Size() lazy evaluation
				b.Clear()
				b.AddAll(slices.Values([]int{1, 2, 3}))
				b.FlipRange(0, 5)  // sets size to -1
				if b.Size() != 2 { // triggers recomputeSize
					t.Errorf("Lazy size evaluation failed: expected 2, got %v", b.Size())
				}

				// Test Remove() panic
				b.Clear()
				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Errorf("Remove() did not panic on empty set")
						}
					}()
					b.Remove()
				}()

				// Test Remove() panic when size is broken
				b.Clear()
				b.size = 1 // break size
				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Errorf("Remove() did not panic on broken size")
						}
					}()
					b.Remove()
				}()

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestBitSet_SizeTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_Size",
			run: func(t *testing.T) {

				b := New()
				if b.Size() != 0 {
					t.Fatalf("expected 0 got %d", b.Size())
				}
				b.Add(5)
				b.Add(10)
				if b.Size() != 2 {
					t.Fatalf("expected 2 got %d", b.Size())
				}

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestBitSet_RemoveTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_Remove",
			run: func(t *testing.T) {

				b := New()
				b.Add(5)
				val := b.Remove()
				if val != 5 {
					t.Fatalf("expected 5, got %d", val)
				}

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestBitSet_RetainAllTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_RetainAll",
			run: func(t *testing.T) {

				b := New()
				b.Add(1)
				b.Add(2)
				b.Add(3)

				other := linkedlist.New[int]()
				other.Add(2)
				other.Add(4)

				b.RetainAll(other)
				if b.Size() != 1 || !b.Contains(2) {
					t.Fatalf("expected only 2 to be retained, got size %d", b.Size())
				}

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestBitSet_ContainsAllTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_ContainsAll",
			run: func(t *testing.T) {

				b := New()
				b.Add(1)
				b.Add(2)

				other := linkedlist.New[int]()
				other.Add(1)
				other.Add(2)

				if !b.ContainsAll(other) {
					t.Fatal("expected contains all to be true")
				}
				other.Add(3)
				if b.ContainsAll(other) {
					t.Fatal("expected contains all to be false")
				}

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestBitSet_FirstLastTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_FirstLast",
			run: func(t *testing.T) {

				b := New()
				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Fatal("expected panic on First() when empty")
						}
					}()
					b.First()
				}()
				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Fatal("expected panic on Last() when empty")
						}
					}()
					b.Last()
				}()
				b.Add(5)
				b.Add(10)
				if v := b.First(); v != 5 {
					t.Fatalf("expected 5, got %d", v)
				}
				if v := b.Last(); v != 10 {
					t.Fatalf("expected 10, got %d", v)
				}

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestBitSet_PollFirstLastTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_PollFirstLast",
			run: func(t *testing.T) {

				b := New()
				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Fatal("expected panic on PollFirst() when empty")
						}
					}()
					b.PollFirst()
				}()
				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Fatal("expected panic on PollLast() when empty")
						}
					}()
					b.PollLast()
				}()
				b.Add(5)
				b.Add(10)
				if v := b.PollFirst(); v != 5 {
					t.Fatalf("expected 5, got %d", v)
				}
				if b.Contains(5) {
					t.Fatal("expected 5 to be removed")
				}
				if v := b.PollLast(); v != 10 {
					t.Fatalf("expected 10, got %d", v)
				}
				if b.Contains(10) {
					t.Fatal("expected 10 to be removed")
				}

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestBitSet_AddFirstLastTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_AddFirstLast",
			run: func(t *testing.T) {

				b := New()
				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Fatal("expected panic on AddFirst")
						}
					}()
					b.AddFirst(1)
				}()
				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Fatal("expected panic on AddLast")
						}
					}()
					b.AddLast(1)
				}()

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestBitSet_NavigableQueriesTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_NavigableQueries",
			run: func(t *testing.T) {

				b := New()

				// Test empty conditions
				if _, ok := b.Lower(100); ok {
					t.Fatal("expected lower to be false on empty")
				}
				if _, ok := b.Floor(100); ok {
					t.Fatal("expected floor to be false on empty")
				}
				if _, ok := b.Ceiling(0); ok {
					t.Fatal("expected ceiling to be false on empty")
				}
				if _, ok := b.Higher(0); ok {
					t.Fatal("expected higher to be false on empty")
				}

				b.Add(5)
				b.Add(10)
				b.Add(15)

				if v, ok := b.Lower(10); !ok || v != 5 {
					t.Fatalf("expected 5, got %d", v)
				}
				if _, ok := b.Lower(5); ok {
					t.Fatal("expected no lower than 5")
				}
				if v, ok := b.Lower(100); !ok || v != 15 {
					t.Fatalf("expected 15, got %d", v)
				}

				if v, ok := b.Floor(10); !ok || v != 10 {
					t.Fatalf("expected 10, got %d", v)
				}
				if v, ok := b.Floor(12); !ok || v != 10 {
					t.Fatalf("expected 10, got %d", v)
				}
				if v, ok := b.Floor(100); !ok || v != 15 {
					t.Fatalf("expected 15, got %d", v)
				}

				if v, ok := b.Ceiling(10); !ok || v != 10 {
					t.Fatalf("expected 10, got %d", v)
				}
				if v, ok := b.Ceiling(12); !ok || v != 15 {
					t.Fatalf("expected 15, got %d", v)
				}
				if v, ok := b.Ceiling(-10); !ok || v != 5 {
					t.Fatalf("expected 5, got %d", v)
				}

				if v, ok := b.Higher(10); !ok || v != 15 {
					t.Fatalf("expected 15, got %d", v)
				}
				if _, ok := b.Higher(15); ok {
					t.Fatal("expected no higher than 15")
				}
				if v, ok := b.Higher(-10); !ok || v != 5 {
					t.Fatalf("expected 5, got %d", v)
				}

				b.Add(63)
				b.Add(64)
				b.Add(127)
				b.Add(128)
				b.Lower(64)
				b.Lower(128)
				b.Floor(63)
				b.Floor(127)
				b.Ceiling(64)
				b.Ceiling(128)
				b.Higher(63)
				b.Higher(127)

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestBitSet_IteratorsTableDriven(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "TestBitSet_Iterators",
			run: func(t *testing.T) {

				b := New()

				// Test empty iterators
				if len(slices.Collect(b.Backward())) != 0 {
					t.Fatal("expected empty backward")
				}
				if len(slices.Collect(b.From(0))) != 0 {
					t.Fatal("expected empty from")
				}
				if len(slices.Collect(b.To(100))) != 0 {
					t.Fatal("expected empty to")
				}
				if len(slices.Collect(b.Between(0, 100))) != 0 {
					t.Fatal("expected empty between")
				}

				b.Add(5)
				b.Add(10)
				b.Add(15)

				back := slices.Collect(b.Backward())
				if !slices.Equal(back, []int{15, 10, 5}) {
					t.Fatalf("expected [15 10 5], got %v", back)
				}

				from := slices.Collect(b.From(10))
				if !slices.Equal(from, []int{10, 15}) {
					t.Fatalf("expected [10 15], got %v", from)
				}

				to := slices.Collect(b.To(15))
				if !slices.Equal(to, []int{5, 10}) {
					t.Fatalf("expected [5 10], got %v", to)
				}

				between := slices.Collect(b.Between(5, 15))
				if !slices.Equal(between, []int{5, 10}) {
					t.Fatalf("expected [5 10], got %v", between)
				}

				b.Add(63)
				b.Add(64)
				if got := slices.Collect(b.From(64)); !slices.Equal(got, []int{64}) {
					t.Fatalf("expected [64], got %v", got)
				}
				if got := slices.Collect(b.To(64)); !slices.Equal(got, []int{5, 10, 15, 63}) {
					t.Fatalf("expected [5 10 15 63], got %v", got)
				}
				if got := slices.Collect(b.Between(63, 65)); !slices.Equal(got, []int{63, 64}) {
					t.Fatalf("expected [63 64], got %v", got)
				}

			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}
