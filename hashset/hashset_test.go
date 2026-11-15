package hashset

import (
	"iter"
	"slices"
	"testing"
)

// testSliceCollection implements collections.Collection[T] without implementing collections.Set[T],
// allowing tests to exercise generic non-Set bulk paths in RetainAll/RemoveAll/ContainsAll.
type testSliceCollection[T any] []T

func (c testSliceCollection[T]) All() iter.Seq[T] {
	return slices.Values(c)
}

func (c testSliceCollection[T]) Size() int {
	return len(c)
}

func (c testSliceCollection[T]) Empty() bool {
	return len(c) == 0
}

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		opts      []Option
		wantPanic bool
		check     func(*testing.T, *HashSet[int])
	}{
		{
			name: "default",
			check: func(t *testing.T, s *HashSet[int]) {
				if s.Size() != 0 || !s.Empty() {
					t.Errorf("expected empty set, got size %d", s.Size())
				}
			},
		},
		{
			name: "zero_capacity",
			opts: []Option{WithCapacity(0)},
			check: func(t *testing.T, s *HashSet[int]) {
				if !s.Empty() {
					t.Errorf("expected empty set")
				}
			},
		},
		{
			name: "positive_capacity",
			opts: []Option{WithCapacity(100)},
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(1)
				if s.Size() != 1 {
					t.Errorf("expected size 1, got %d", s.Size())
				}
			},
		},
		{
			name:      "negative_capacity_panics",
			opts:      []Option{WithCapacity(-1)},
			wantPanic: true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				r := recover()
				if tc.wantPanic && r == nil {
					t.Errorf("expected panic, got none")
				}
				if !tc.wantPanic && r != nil {
					t.Errorf("unexpected panic: %v", r)
				}
			}()
			s := New[int](tc.opts...)
			if tc.check != nil {
				tc.check(t, s)
			}
		})
	}
}

func TestHashSet_AddRemoveContains(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *HashSet[int])
	}{
		{
			name: "add_and_contains",
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(2)
				if s.Size() != 2 || s.Empty() {
					t.Errorf("expected size 2, got %d", s.Size())
				}
				if !s.Contains(1) || !s.Contains(2) || s.Contains(3) {
					t.Errorf("contains check failed")
				}
			},
		},
		{
			name: "remove_element_present_and_absent",
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.RemoveElement(1)
				s.RemoveElement(99)
				if s.Contains(1) || !s.Contains(2) || s.Size() != 1 {
					t.Errorf("remove element failed, size=%d", s.Size())
				}
			},
		},
		{
			name: "remove_arbitrary_drains_set",
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(10)
				s.Add(20)
				first := s.Remove()
				second := s.Remove()
				if (first != 10 && first != 20) || (second != 10 && second != 20) || first == second {
					t.Errorf("unexpected drained items: %d, %d", first, second)
				}
				if !s.Empty() {
					t.Errorf("expected empty set after draining")
				}
			},
		},
		{
			name: "remove_empty_panics",
			check: func(t *testing.T, s *HashSet[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic on empty Remove()")
					}
				}()
				s.Remove()
			},
		},
		{
			name: "clear",
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Clear()
				if !s.Empty() || s.Size() != 0 {
					t.Errorf("expected empty set after Clear")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := New[int]()
			tc.check(t, s)
		})
	}
}

func TestHashSet_BulkAndIteration(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *HashSet[int])
	}{
		{
			name: "add_all",
			check: func(t *testing.T, s *HashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3, 2}))
				if s.Size() != 3 {
					t.Errorf("expected size 3, got %d", s.Size())
				}
			},
		},
		{
			name: "contains_all_including_self_alias",
			check: func(t *testing.T, s *HashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2}))
				if !s.ContainsAll(s) {
					t.Errorf("expected self ContainsAll to return true")
				}

				other := New[int]()
				other.Add(1)
				if !s.ContainsAll(other) {
					t.Errorf("expected subset ContainsAll to return true")
				}

				other.Add(3)
				if s.ContainsAll(other) {
					t.Errorf("expected non-subset ContainsAll to return false")
				}
			},
		},
		{
			name: "remove_all_empty_normal_and_self_alias",
			check: func(t *testing.T, s *HashSet[int]) {
				other := New[int]()
				other.Add(1)
				// Empty receiver early-return branch
				s.RemoveAll(other)
				if !s.Empty() {
					t.Errorf("expected empty set to remain empty")
				}

				s.AddAll(slices.Values([]int{1, 2, 3}))
				other.Add(2)
				s.RemoveAll(other)
				if s.Size() != 1 || !s.Contains(3) {
					t.Errorf("expected {3}, got %s", s.String())
				}

				// Self-aliasing branch
				s.Add(4)
				s.RemoveAll(s)
				if !s.Empty() {
					t.Errorf("expected empty set after RemoveAll(self), got %s", s.String())
				}
			},
		},
		{
			name: "retain_all_empty_set_collection_and_self_alias",
			check: func(t *testing.T, s *HashSet[int]) {
				otherSet := New[int]()
				otherSet.Add(2)
				// Empty receiver early-return branch
				s.RetainAll(otherSet)
				if !s.Empty() {
					t.Errorf("expected empty set to remain empty")
				}

				// Self-aliasing branch
				s.AddAll(slices.Values([]int{1, 2, 3}))
				s.RetainAll(s)
				if s.Size() != 3 {
					t.Errorf("expected self RetainAll to preserve size 3, got %d", s.Size())
				}

				// Set fast-path branch
				otherSet.Add(4)
				s.RetainAll(otherSet)
				if s.Size() != 1 || !s.Contains(2) {
					t.Errorf("expected {2} after Set RetainAll, got %s", s.String())
				}

				// Non-Set Collection fallback branch (with duplicates in other)
				s.AddAll(slices.Values([]int{10, 20, 30}))
				s.RetainAll(testSliceCollection[int]{20, 20, 30, 99})
				got := slices.Collect(s.All())
				slices.Sort(got)
				if !slices.Equal(got, []int{20, 30}) {
					t.Errorf("expected [20 30] after non-Set RetainAll, got %v", got)
				}
			},
		},
		{
			name: "iterator_full_and_early_break",
			check: func(t *testing.T, s *HashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				got := slices.Collect(s.All())
				slices.Sort(got)
				if !slices.Equal(got, []int{1, 2, 3}) {
					t.Errorf("expected [1 2 3], got %v", got)
				}

				count := 0
				for range s.All() {
					count++
					break
				}
				if count != 1 {
					t.Errorf("expected early break count 1, got %d", count)
				}
			},
		},
		{
			name: "string_empty_single_and_multiple",
			check: func(t *testing.T, s *HashSet[int]) {
				if got := s.String(); got != "[]" {
					t.Errorf("expected [], got %s", got)
				}
				s.Add(1)
				if got := s.String(); got != "[1]" {
					t.Errorf("expected [1], got %s", got)
				}
				s.Add(2)
				got := s.String()
				if got != "[1 2]" && got != "[2 1]" {
					t.Errorf("expected [1 2] or [2 1], got %s", got)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := New[int]()
			tc.check(t, s)
		})
	}
}
