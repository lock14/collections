package bitset

import (
	"bytes"
	"fmt"
	"iter"
	"math"
	"slices"
	"testing"

	"github.com/lock14/collections"
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
	353, 359, 367, 373, 379, 383, 389, 397, 397 - 0, // keep exact list below
}

func init() {
	first100Primes = []int{
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
}

// testSliceCollection is a minimal collections.Collection[int] (non-Set) for testing generic fallbacks.
type testSliceCollection []int

func (c testSliceCollection) All() iter.Seq[int] { return slices.Values(c) }
func (c testSliceCollection) Size() int          { return len(c) }
func (c testSliceCollection) Empty() bool        { return len(c) == 0 }

// testIntSet is a minimal collections.Set[int] (non-*BitSet) for testing Set-specialized fallbacks.
type testIntSet map[int]struct{}

func newTestIntSet(vals ...int) testIntSet {
	m := make(testIntSet, len(vals))
	for _, v := range vals {
		m[v] = struct{}{}
	}
	return m
}

func (m testIntSet) All() iter.Seq[int] {
	return func(yield func(int) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}
func (m testIntSet) Size() int           { return len(m) }
func (m testIntSet) Empty() bool         { return len(m) == 0 }
func (m testIntSet) Contains(v int) bool { _, ok := m[v]; return ok }
func (m testIntSet) ContainsAll(col collections.Collection[int]) bool {
	for v := range col.All() {
		if !m.Contains(v) {
			return false
		}
	}
	return true
}

func assertPanics(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: expected panic, did not panic", name)
		}
	}()
	fn()
}

func primesLessThan(n int) *BitSet {
	s := New(WithCapacity(n))
	if n > 2 {
		s.SetBit(0)
		s.SetBit(1)
		for i := 4; i < n; i += 2 {
			s.SetBit(i)
		}
		for i := 3; (i*i) > i && (i*i) < n; i += 2 {
			if !s.GetBit(i) {
				for j := i * i; j > i && j < n; j += i {
					s.SetBit(j)
				}
			}
		}
		s.FlipRange(0, n)
	}
	return s
}

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		opts      []Option
		wantPanic bool
		check     func(*testing.T, *BitSet)
	}{
		{
			name: "default_capacity",
			check: func(t *testing.T, s *BitSet) {
				if got := s.Capacity(); got != 64 {
					t.Errorf("expected Capacity 64, got %d", got)
				}
				if got := s.Length(); got != 0 {
					t.Errorf("expected Length 0, got %d", got)
				}
				if !s.Empty() || s.Size() != 0 {
					t.Errorf("expected Empty()=true, Size()=0")
				}
				for i := 0; i < 64; i++ {
					if s.GetBit(i) {
						t.Errorf("expected bit %d to be unset", i)
					}
				}
			},
		},
		{
			name: "zero_bits",
			opts: []Option{WithCapacity(0)},
			check: func(t *testing.T, s *BitSet) {
				if got := s.Capacity(); got != 0 {
					t.Errorf("Capacity() = %v, want 0", got)
				}
				s.SetBit(5)
				if !s.GetBit(5) {
					t.Errorf("expected bit 5 to be set after growing from zero")
				}
			},
		},
		{
			name: "exact_word_size",
			opts: []Option{WithCapacity(64)},
			check: func(t *testing.T, s *BitSet) {
				if got := s.Capacity(); got != 64 {
					t.Errorf("expected Capacity 64, got %d", got)
				}
			},
		},
		{
			name: "non_exact_word_size",
			opts: []Option{WithCapacity(100)},
			check: func(t *testing.T, s *BitSet) {
				if got := s.Capacity(); got != 128 {
					t.Errorf("Capacity() = %v, want 128", got)
				}
			},
		},
		{
			name:      "negative_capacity_panics",
			opts:      []Option{WithCapacity(-1)},
			wantPanic: true,
		},
		{
			name: "zero_value_struct_usable",
			check: func(t *testing.T, _ *BitSet) {
				var zero BitSet
				if !zero.Empty() || zero.Size() != 0 || zero.Capacity() != 0 || zero.String() != "[]" {
					t.Errorf("unexpected zero-value state")
				}
				zero.Add(70)
				if zero.Empty() || zero.Size() != 1 || !zero.Contains(70) {
					t.Errorf("zero-value BitSet failed after Add(70)")
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.wantPanic {
				assertPanics(t, tc.name, func() { New(tc.opts...) })
				return
			}
			s := New(tc.opts...)
			tc.check(t, s)
		})
	}
}

func TestBitOperations(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *BitSet)
	}{
		{
			name: "get_out_of_range_returns_false_without_growing",
			check: func(t *testing.T, s *BitSet) {
				capBefore := s.Capacity()
				if s.GetBit(1000) {
					t.Errorf("expected false for out-of-range bit")
				}
				if got := s.Capacity(); got != capBefore {
					t.Errorf("Capacity() = %d, want %d", got, capBefore)
				}
			},
		},
		{
			name: "negative_index_panics_on_get_set_clear",
			check: func(t *testing.T, s *BitSet) {
				assertPanics(t, "GetBit(-1)", func() { s.GetBit(-1) })
				assertPanics(t, "SetBit(-1)", func() { s.SetBit(-1) })
				assertPanics(t, "ClearBit(-1)", func() { s.ClearBit(-1) })
				assertPanics(t, "Add(-1)", func() { s.Add(-1) })
				assertPanics(t, "Contains(-1)", func() { s.Contains(-1) })
				assertPanics(t, "RemoveElement(-1)", func() { s.RemoveElement(-1) })
			},
		},
		{
			name: "set_and_idempotent_set",
			check: func(t *testing.T, s *BitSet) {
				s.SetBit(5)
				s.SetBit(5)
				if !s.GetBit(5) || s.Size() != 1 {
					t.Errorf("expected bit 5 set once, size=%d", s.Size())
				}
			},
		},
		{
			name: "geometric_growth_across_many_words",
			check: func(t *testing.T, s *BitSet) {
				for i := 0; i < 512; i++ {
					s.SetBit(i)
					if !s.GetBit(i) {
						t.Fatalf("expected bit %d to be set", i)
					}
				}
				if s.Size() != 512 || s.Length() != 512 || s.Capacity() != 512 {
					t.Errorf("got size=%d len=%d cap=%d", s.Size(), s.Length(), s.Capacity())
				}
			},
		},
		{
			name: "clear_bit_variants_and_after_invalidated_size",
			check: func(t *testing.T, s *BitSet) {
				s.SetBit(5)
				s.SetBit(10)
				s.ClearBit(10)
				if s.GetBit(10) || s.Length() != 6 || s.Size() != 1 {
					t.Errorf("unexpected state after ClearBit(10)")
				}
				// Clear unset bit & out-of-range bit are no-ops
				s.ClearBit(3)
				s.ClearBit(1000)
				if !s.GetBit(5) || s.Capacity() != 64 {
					t.Errorf("no-op ClearBit mutated state")
				}
				s.ClearBit(5)
				if !s.Empty() || s.Length() != 0 || s.Size() != 0 {
					t.Errorf("expected empty after clearing last bit")
				}
				// Clear last bit while size == -1 (invalidated by FlipRange)
				s.FlipRange(0, 2)
				s.ClearBit(0)
				s.ClearBit(1)
				if !s.Empty() || s.Size() != 0 {
					t.Errorf("expected Empty()=true and Size()=0 after clearing flipped bits")
				}
			},
		},
		{
			name: "length_across_word_boundaries",
			check: func(t *testing.T, s *BitSet) {
				if s.Length() != 0 {
					t.Errorf("expected 0, got %d", s.Length())
				}
				for _, pair := range [][2]int{{0, 1}, {5, 6}, {63, 64}, {64, 65}, {100, 101}} {
					s.SetBit(pair[0])
					if got := s.Length(); got != pair[1] {
						t.Errorf("after SetBit(%d): Length() = %d, want %d", pair[0], got, pair[1])
					}
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := New()
			tc.check(t, s)
		})
	}
}

func TestFlipAndFlipRange(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *BitSet)
	}{
		{
			name: "flip_empty_and_back_to_empty_preserves_empty_invariants",
			check: func(t *testing.T, s *BitSet) {
				s.Flip()
				if !s.GetBit(0) || !s.GetBit(63) || s.Length() != 64 || s.Empty() {
					t.Errorf("expected 64 set bits after Flip()")
				}
				s.Flip()
				// Regression test for Bug #1: Empty/First/Last/Remove/Higher(-1) after Flip clears all bits
				if !s.Empty() {
					t.Fatalf("expected Empty()=true after second Flip()")
				}
				if s.Size() != 0 {
					t.Fatalf("expected Size()=0, got %d", s.Size())
				}
				if _, ok := s.Higher(-1); ok {
					t.Errorf("expected Higher(-1) == false on empty flipped set")
				}
				assertPanics(t, "First on flipped-empty", func() { s.First() })
				assertPanics(t, "Last on flipped-empty", func() { s.Last() })
				assertPanics(t, "Remove on flipped-empty", func() { s.Remove() })
			},
		},
		{
			name: "flip_twice_restores_bits",
			check: func(t *testing.T, s *BitSet) {
				s.SetBit(5)
				s.Flip()
				s.Flip()
				if got := slices.Collect(s.SetBits()); !slices.Equal(got, []int{5}) {
					t.Errorf("expected [5], got %v", got)
				}
			},
		},
		{
			name: "flip_range_zero_capacity_noop_and_growth",
			opts: []Option{WithCapacity(0)},
			check: func(t *testing.T, s *BitSet) {
				// Regression test for Bug #2: FlipRange(0, 0) on capacity 0 must not panic
				s.FlipRange(0, 0)
				if !s.Empty() || s.Capacity() != 0 {
					t.Errorf("expected noop for FlipRange(0,0) on cap 0")
				}
				s.FlipRange(0, 64)
				if s.Capacity() != 64 || s.Size() != 64 {
					t.Errorf("got cap=%d size=%d, want 64, 64", s.Capacity(), s.Size())
				}
			},
		},
		{
			name: "flip_range_at_exact_capacity_boundary_is_noop",
			opts: []Option{WithCapacity(64)},
			check: func(t *testing.T, s *BitSet) {
				// Regression test for Bug #2: FlipRange(64, 64) on capacity 64 must not panic
				s.FlipRange(64, 64)
				if !s.Empty() || s.Capacity() != 64 {
					t.Errorf("expected noop at capacity boundary")
				}
			},
		},
		{
			name: "flip_range_exact_word_multiple_does_not_overallocate_words",
			opts: []Option{WithCapacity(64)},
			check: func(t *testing.T, s *BitSet) {
				// Regression test for Bug #2: FlipRange(0, 128) needs 2 words (128 bits), not 3 (192 bits)
				s.FlipRange(0, 128)
				if got := s.Capacity(); got != 128 {
					t.Errorf("Capacity() = %d, want 128", got)
				}
				if got := s.Size(); got != 128 {
					t.Errorf("Size() = %d, want 128", got)
				}
			},
		},
		{
			name: "flip_range_panics_on_negative_or_inverted_bounds",
			check: func(t *testing.T, s *BitSet) {
				assertPanics(t, "FlipRange(-1, 5)", func() { s.FlipRange(-1, 5) })
				assertPanics(t, "FlipRange(10, 5)", func() { s.FlipRange(10, 5) })
				assertPanics(t, "FlipRange(70, 10)", func() { s.FlipRange(70, 10) })
			},
		},
		{
			name: "flip_range_partial_cross_and_multi_word",
			check: func(t *testing.T, s *BitSet) {
				s.FlipRange(2, 5)
				if got := slices.Collect(s.SetBits()); !slices.Equal(got, []int{2, 3, 4}) {
					t.Errorf("got %v, want [2 3 4]", got)
				}
				// Flip same range back to empty -> test Bug #1 invariant
				s.FlipRange(2, 5)
				if !s.Empty() || s.Size() != 0 {
					t.Errorf("expected empty after unflipping [2, 5)")
				}

				s.FlipRange(60, 68)
				if got := slices.Collect(s.SetBits()); !slices.Equal(got, []int{60, 61, 62, 63, 64, 65, 66, 67}) {
					t.Errorf("got %v", got)
				}
				s.Clear()

				s.FlipRange(60, 130)
				if !s.GetBit(60) || !s.GetBit(127) || !s.GetBit(129) || s.GetBit(59) || s.GetBit(130) {
					t.Errorf("flip range failed for multi-word span")
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := New(tc.opts...)
			tc.check(t, s)
		})
	}
}

func TestMutableSetOperations(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T)
	}{
		{
			name: "add_remove_element_remove_and_clear",
			check: func(t *testing.T) {
				s := New(WithCapacity(100))
				s.Add(10)
				s.Add(20)
				if !s.Contains(10) || !s.Contains(20) || s.Size() != 2 {
					t.Fatalf("expected {10, 20}")
				}
				s.RemoveElement(10)
				if s.Contains(10) || s.Empty() {
					t.Fatalf("expected {20}")
				}
				if got := s.Remove(); got != 20 || !s.Empty() {
					t.Fatalf("expected Remove()=20 and empty set, got %d", got)
				}
				assertPanics(t, "Remove on empty", func() { s.Remove() })

				s.AddAll(slices.Values([]int{3, 7, 15}))
				s.Clear()
				if !s.Empty() || s.Size() != 0 {
					t.Fatalf("expected empty after Clear()")
				}
			},
		},
		{
			name: "remove_all_bitset_set_and_generic_collection",
			check: func(t *testing.T) {
				// BitSet fast path clearing all bits -> Empty() must be true (Bug #1 regression)
				s := New()
				s.AddAll(slices.Values([]int{1, 2, 3}))
				other := New()
				other.AddAll(slices.Values([]int{1, 2, 3, 4}))
				s.RemoveAll(other)
				if !s.Empty() || s.Size() != 0 {
					t.Fatalf("expected Empty()=true after RemoveAll cleared all bits")
				}
				assertPanics(t, "Remove after RemoveAll emptied set", func() { s.Remove() })

				// Partial RemoveAll with BitSet
				s.AddAll(slices.Values([]int{1, 2, 3, 100}))
				other.Clear()
				other.AddAll(slices.Values([]int{2, 3}))
				s.RemoveAll(other)
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{1, 100}) {
					t.Errorf("got %v, want [1 100]", got)
				}

				// RemoveAll with generic Set[int]
				s.RemoveAll(newTestIntSet(1))
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{100}) {
					t.Errorf("got %v, want [100]", got)
				}

				// RemoveAll with generic non-Set Collection[int]
				s.AddAll(slices.Values([]int{5, 10}))
				s.RemoveAll(testSliceCollection{5, 100})
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{10}) {
					t.Errorf("got %v, want [10]", got)
				}
			},
		},
		{
			name: "retain_all_bitset_set_and_generic_collection",
			check: func(t *testing.T) {
				// Disjoint RetainAll(*BitSet) -> Empty() must be true (Bug #1 regression)
				s := New()
				s.AddAll(slices.Values([]int{1, 2, 100}))
				other := New()
				other.AddAll(slices.Values([]int{50, 60}))
				s.RetainAll(other)
				if !s.Empty() || s.Size() != 0 {
					t.Fatalf("expected Empty()=true after disjoint RetainAll")
				}

				// RetainAll(*BitSet) where receiver has more words than other
				s.AddAll(slices.Values([]int{1, 2, 3, 200}))
				other.Clear()
				other.AddAll(slices.Values([]int{2, 3, 4}))
				s.RetainAll(other)
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{2, 3}) {
					t.Errorf("got %v, want [2 3]", got)
				}

				// RetainAll(Set[int])
				s.AddAll(slices.Values([]int{1, 2, 3, 4, 150}))
				s.RetainAll(newTestIntSet(2, 4, 999))
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{2, 4}) {
					t.Errorf("RetainAll(Set): got %v, want [2 4]", got)
				}

				// RetainAll(Collection[int])
				s.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
				s.RetainAll(testSliceCollection{3, 5, 99})
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{3, 5}) {
					t.Errorf("RetainAll(Collection): got %v, want [3 5]", got)
				}
			},
		},
		{
			name: "contains_all_bitset_and_generic_collection",
			check: func(t *testing.T) {
				s := New()
				s.AddAll(slices.Values([]int{2, 3}))

				sub := New()
				sub.AddAll(slices.Values([]int{2, 3}))
				if !s.ContainsAll(sub) {
					t.Errorf("expected ContainsAll(equal)=true")
				}

				sup := New()
				sup.AddAll(slices.Values([]int{2, 3, 4}))
				if s.ContainsAll(sup) {
					t.Errorf("expected ContainsAll(superset)=false")
				}

				wider := New()
				wider.AddAll(slices.Values([]int{2, 3, 150}))
				if s.ContainsAll(wider) {
					t.Errorf("expected ContainsAll(wider)=false")
				}

				// Defensive check when other has artificially inflated maxWordInUse with zero upper words
				padded := New()
				padded.Add(2)
				padded.ensureSize(4)
				padded.maxWordInUse = 5
				if !s.ContainsAll(padded) {
					t.Errorf("expected ContainsAll(padded)=true")
				}

				if !s.ContainsAll(testSliceCollection{2, 3}) {
					t.Errorf("expected ContainsAll(slice)=true")
				}
				if s.ContainsAll(testSliceCollection{2, 3, 99}) {
					t.Errorf("expected ContainsAll(slice with missing)=false")
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t)
		})
	}
}

func TestSequencedAndNavigable(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		elems []int
		check func(*testing.T, *BitSet)
	}{
		{
			name:  "empty_set_panics_and_false_queries",
			elems: nil,
			check: func(t *testing.T, s *BitSet) {
				assertPanics(t, "First", func() { s.First() })
				assertPanics(t, "Last", func() { s.Last() })
				assertPanics(t, "PollFirst", func() { s.PollFirst() })
				assertPanics(t, "PollLast", func() { s.PollLast() })

				for _, q := range []struct {
					name string
					fn   func(int) (int, bool)
				}{
					{"Lower", s.Lower},
					{"Floor", s.Floor},
					{"Ceiling", s.Ceiling},
					{"Higher", s.Higher},
				} {
					if _, ok := q.fn(-5); ok {
						t.Errorf("%s(-5) expected false on empty", q.name)
					}
					if _, ok := q.fn(0); ok {
						t.Errorf("%s(0) expected false on empty", q.name)
					}
					if _, ok := q.fn(100); ok {
						t.Errorf("%s(100) expected false on empty", q.name)
					}
				}
			},
		},
		{
			name:  "multi_word_first_last_poll",
			elems: []int{70, 130, 200},
			check: func(t *testing.T, s *BitSet) {
				// First must skip empty word 0 cleanly
				if got := s.First(); got != 70 {
					t.Errorf("First() = %d, want 70", got)
				}
				if got := s.Last(); got != 200 {
					t.Errorf("Last() = %d, want 200", got)
				}
				if got := s.PollFirst(); got != 70 || s.Contains(70) {
					t.Errorf("PollFirst() = %d", got)
				}
				if got := s.PollLast(); got != 200 || s.Contains(200) {
					t.Errorf("PollLast() = %d", got)
				}
				if got := s.PollLast(); got != 130 || !s.Empty() {
					t.Errorf("final PollLast() = %d", got)
				}
			},
		},
		{
			name:  "navigable_queries_across_word_boundaries_and_maxint",
			elems: []int{5, 10, 15, 63, 64, 127, 195},
			check: func(t *testing.T, s *BitSet) {
				navCases := []struct {
					label  string
					fn     func(int) (int, bool)
					arg    int
					want   int
					wantOk bool
				}{
					{"Lower(-5)", s.Lower, -5, 0, false},
					{"Lower(0)", s.Lower, 0, 0, false},
					{"Lower(5)", s.Lower, 5, 0, false},
					{"Lower(10)", s.Lower, 10, 5, true},
					{"Lower(64)", s.Lower, 64, 63, true},
					{"Lower(65)", s.Lower, 65, 64, true},
					{"Lower(190)", s.Lower, 190, 127, true},
					{"Lower(500)", s.Lower, 500, 195, true},

					{"Floor(-1)", s.Floor, -1, 0, false},
					{"Floor(4)", s.Floor, 4, 0, false},
					{"Floor(5)", s.Floor, 5, 5, true},
					{"Floor(12)", s.Floor, 12, 10, true},
					{"Floor(63)", s.Floor, 63, 63, true},
					{"Floor(190)", s.Floor, 190, 127, true},
					{"Floor(500)", s.Floor, 500, 195, true},

					{"Ceiling(-10)", s.Ceiling, -10, 5, true},
					{"Ceiling(0)", s.Ceiling, 0, 5, true},
					{"Ceiling(10)", s.Ceiling, 10, 10, true},
					{"Ceiling(12)", s.Ceiling, 12, 15, true},
					{"Ceiling(64)", s.Ceiling, 64, 64, true},
					{"Ceiling(130)", s.Ceiling, 130, 195, true},
					{"Ceiling(196)", s.Ceiling, 196, 0, false},
					{"Ceiling(500)", s.Ceiling, 500, 0, false},

					{"Higher(-10)", s.Higher, -10, 5, true},
					{"Higher(10)", s.Higher, 10, 15, true},
					{"Higher(63)", s.Higher, 63, 64, true},
					{"Higher(127)", s.Higher, 127, 195, true},
					{"Higher(195)", s.Higher, 195, 0, false},
					// Regression test for Bug #2: Higher(math.MaxInt) must not overflow to MinInt and panic
					{"Higher(MaxInt)", s.Higher, math.MaxInt, 0, false},
				}

				for _, nc := range navCases {
					got, ok := nc.fn(nc.arg)
					if ok != nc.wantOk || (ok && got != nc.want) {
						t.Errorf("%s = (%d, %v), want (%d, %v)", nc.label, got, ok, nc.want, nc.wantOk)
					}
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := New()
			s.AddAll(slices.Values(tc.elems))
			tc.check(t, s)
		})
	}
}

func TestIterators(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		elems []int
		check func(*testing.T, *BitSet)
	}{
		{
			name:  "empty_set_iterators",
			elems: nil,
			check: func(t *testing.T, s *BitSet) {
				if len(slices.Collect(s.All())) != 0 {
					t.Errorf("expected empty All")
				}
				if len(slices.Collect(s.Backward())) != 0 {
					t.Errorf("expected empty Backward")
				}
				if len(slices.Collect(s.From(0))) != 0 {
					t.Errorf("expected empty From")
				}
				if len(slices.Collect(s.To(100))) != 0 {
					t.Errorf("expected empty To")
				}
				if len(slices.Collect(s.Between(0, 100))) != 0 {
					t.Errorf("expected empty Between")
				}
			},
		},
		{
			name:  "range_iterators_and_word_boundaries",
			elems: []int{0, 5, 10, 15, 63, 64, 127, 128},
			check: func(t *testing.T, s *BitSet) {
				if got := slices.Collect(s.Backward()); !slices.Equal(got, []int{128, 127, 64, 63, 15, 10, 5, 0}) {
					t.Errorf("Backward() = %v", got)
				}
				if got := slices.Collect(s.From(-5)); !slices.Equal(got, []int{0, 5, 10, 15, 63, 64, 127, 128}) {
					t.Errorf("From(-5) = %v", got)
				}
				if got := slices.Collect(s.From(64)); !slices.Equal(got, []int{64, 127, 128}) {
					t.Errorf("From(64) = %v", got)
				}
				if got := slices.Collect(s.From(500)); len(got) != 0 {
					t.Errorf("From(500) = %v", got)
				}
				if got := slices.Collect(s.To(0)); len(got) != 0 {
					t.Errorf("To(0) = %v", got)
				}
				if got := slices.Collect(s.To(-5)); len(got) != 0 {
					t.Errorf("To(-5) = %v", got)
				}
				if got := slices.Collect(s.To(64)); !slices.Equal(got, []int{0, 5, 10, 15, 63}) {
					t.Errorf("To(64) = %v", got)
				}
				if got := slices.Collect(s.Between(-10, 12)); !slices.Equal(got, []int{0, 5, 10}) {
					t.Errorf("Between(-10, 12) = %v", got)
				}
				if got := slices.Collect(s.Between(63, 65)); !slices.Equal(got, []int{63, 64}) {
					t.Errorf("Between(63, 65) = %v", got)
				}
				if got := slices.Collect(s.Between(10, 10)); len(got) != 0 {
					t.Errorf("Between(10, 10) = %v", got)
				}
				if got := slices.Collect(s.Between(20, 10)); len(got) != 0 {
					t.Errorf("Between(20, 10) = %v", got)
				}
				if got := slices.Collect(s.Between(-10, 0)); len(got) != 0 {
					t.Errorf("Between(-10, 0) = %v", got)
				}
				if got := slices.Collect(s.Between(500, 600)); len(got) != 0 {
					t.Errorf("Between(500, 600) = %v", got)
				}
			},
		},
		{
			name:  "early_break_all_iterators",
			elems: []int{1, 5, 64, 70, 130},
			check: func(t *testing.T, s *BitSet) {
				take2 := func(seq iter.Seq[int]) []int {
					var out []int
					for v := range seq {
						out = append(out, v)
						if len(out) == 2 {
							break
						}
					}
					return out
				}
				if got := take2(s.All()); !slices.Equal(got, []int{1, 5}) {
					t.Errorf("All break1 = %v", got)
				}
				// Also test breaking in the second word loop of From
				take3 := func(seq iter.Seq[int]) []int {
					var out []int
					for v := range seq {
						out = append(out, v)
						if len(out) == 3 {
							break
						}
					}
					return out
				}
				if got := take3(s.From(1)); !slices.Equal(got, []int{1, 5, 64}) {
					t.Errorf("From break3 = %v", got)
				}
				if got := take2(s.From(1)); !slices.Equal(got, []int{1, 5}) {
					t.Errorf("From break2 = %v", got)
				}
				if got := take2(s.Backward()); !slices.Equal(got, []int{130, 70}) {
					t.Errorf("Backward break2 = %v", got)
				}
				if got := take2(s.To(100)); !slices.Equal(got, []int{1, 5}) {
					t.Errorf("To break2 = %v", got)
				}
				if got := take2(s.Between(1, 100)); !slices.Equal(got, []int{1, 5}) {
					t.Errorf("Between break2 = %v", got)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := New()
			s.AddAll(slices.Values(tc.elems))
			tc.check(t, s)
		})
	}
}

func TestStringAndBytes(t *testing.T) {
	t.Parallel()
	strCases := []struct {
		name  string
		elems []int
		want  string
	}{
		{name: "empty", elems: nil, want: "[]"},
		{name: "single", elems: []int{4}, want: "[4]"},
		{name: "multi_same_word", elems: []int{0, 4, 63}, want: "[0 4 63]"},
		{name: "multi_words", elems: []int{0, 64, 127}, want: "[0 64 127]"},
	}

	for _, tc := range strCases {
		tc := tc
		t.Run("string_"+tc.name, func(t *testing.T) {
			t.Parallel()
			s := New()
			s.AddAll(slices.Values(tc.elems))
			if got := s.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}

	byteCases := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{name: "nil_slice", input: nil, want: []byte{}},
		{name: "empty_slice", input: []byte{}, want: []byte{}},
		{name: "trailing_zero_bytes_trimmed", input: []byte{0x05, 0x00, 0x00}, want: []byte{0x05}},
		{name: "one_byte", input: []byte{0xFF}, want: []byte{0xFF}},
		{name: "eight_bytes", input: []byte{0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF}, want: []byte{0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF}},
		{name: "twelve_bytes", input: []byte{0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF}, want: []byte{0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF}},
		{name: "sixteen_bytes", input: []byte{0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF}, want: []byte{0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF}},
	}

	for _, tc := range byteCases {
		tc := tc
		t.Run("bytes_"+tc.name, func(t *testing.T) {
			t.Parallel()
			s := FromBytes(tc.input)
			got := s.ToBytes()
			if !bytes.Equal(got, tc.want) {
				t.Errorf("ToBytes() = %v, want %v", got, tc.want)
			}
			for n := 0; n < s.Length(); n++ {
				gotSetBit := (got[n/8] & (1 << (n % 8))) != 0
				if gotSetBit != s.GetBit(n) {
					t.Errorf("bit %d mismatch", n)
				}
			}
		})
	}
}

func TestBitSetPrimeGen(t *testing.T) {
	t.Parallel()
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
			s := primesLessThan(tc.lessThan)
			primes := slices.Collect(s.SetBits())
			if !slices.Equal(primes, tc.want) {
				t.Errorf("got %v, want %v", primes, tc.want)
			}
		})
	}
}
