package treeset

import (
	"iter"
	"slices"
	"testing"

	"github.com/lock14/collections/comparator"
)

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

func TestTreeSet_Constructors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		wantPanic bool
		check     func(*testing.T)
	}{
		{
			name:      "missing_comparator_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				New[int]()
			},
		},
		{
			name:      "invalid_degree_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				NewOrdered[int](WithDegree[int](1))
			},
		},
		{
			name: "custom_degree_and_comparator",
			check: func(t *testing.T) {
				s := New[int](
					WithDegree[int](4),
					WithComparator(comparator.Reverse(comparator.NaturalOrder[int]())),
				)
				s.Add(1)
				s.Add(3)
				s.Add(2)
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{3, 2, 1}) {
					t.Errorf("expected [3 2 1], got %v", got)
				}
			},
		},
		{
			name: "new_ordered_option_precedence_and_slice_non_aliasing",
			check: func(t *testing.T) {
				opts := make([]Option[int], 1, 4)
				opts[0] = WithComparator(comparator.Reverse(comparator.NaturalOrder[int]()))
				sentinelCalled := false
				fullSlice := opts[:2]
				fullSlice[1] = func(*config[int]) { sentinelCalled = true }

				s := NewOrdered[int](opts...)
				if sentinelCalled {
					t.Fatalf("NewOrdered overwrote or read spare capacity of caller slice")
				}
				s.Add(10)
				s.Add(20)
				if s.First() != 20 {
					t.Errorf("expected custom ReverseOrder comparator to win in NewOrdered, got First()=%d", s.First())
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				r := recover()
				if tc.wantPanic && r == nil {
					t.Errorf("expected panic but got none")
				}
				if !tc.wantPanic && r != nil {
					t.Errorf("expected no panic but got: %v", r)
				}
			}()
			tc.check(t)
		})
	}
}

func TestTreeSet_OperationsAndBulk(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *TreeSet[int])
	}{
		{
			name: "add_duplicate_contains_remove_element_and_clear",
			check: func(t *testing.T, s *TreeSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(2)
				s.Add(3)
				if s.Size() != 3 || !s.Contains(2) || s.Contains(4) {
					t.Errorf("unexpected set state: %s", s.String())
				}
				s.RemoveElement(2)
				s.RemoveElement(99)
				if s.Contains(2) || s.Size() != 2 {
					t.Errorf("RemoveElement failed: %s", s.String())
				}
				s.Clear()
				if !s.Empty() || s.Size() != 0 {
					t.Errorf("Clear failed")
				}
			},
		},
		{
			name: "remove_returns_smallest_and_panics_when_empty",
			check: func(t *testing.T, s *TreeSet[int]) {
				s.AddAll(slices.Values([]int{30, 10, 20}))
				if got := s.Remove(); got != 10 || s.Size() != 2 {
					t.Errorf("Remove() = %d, want 10", got)
				}
				s.Clear()
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic on empty Remove()")
					}
				}()
				s.Remove()
			},
		},
		{
			name: "contains_all_including_self_alias",
			check: func(t *testing.T, s *TreeSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				if !s.ContainsAll(s) {
					t.Errorf("expected ContainsAll(self) to return true")
				}
				if !s.ContainsAll(testSliceCollection[int]{1, 2}) {
					t.Errorf("expected subset ContainsAll to return true")
				}
				if s.ContainsAll(testSliceCollection[int]{1, 4}) {
					t.Errorf("expected non-subset ContainsAll to return false")
				}
			},
		},
		{
			name: "remove_all_empty_normal_and_self_alias",
			check: func(t *testing.T, s *TreeSet[int]) {
				s.RemoveAll(testSliceCollection[int]{1, 2})
				if !s.Empty() {
					t.Errorf("expected empty set to remain empty")
				}

				s.AddAll(slices.Values([]int{1, 2, 3}))
				s.RemoveAll(testSliceCollection[int]{2, 3})
				if s.Size() != 1 || !s.Contains(1) {
					t.Errorf("expected [1], got %s", s.String())
				}

				s.Add(2)
				s.RemoveAll(s)
				if !s.Empty() {
					t.Errorf("expected empty set after RemoveAll(self)")
				}
			},
		},
		{
			name: "retain_all_empty_normal_and_self_alias",
			check: func(t *testing.T, s *TreeSet[int]) {
				s.RetainAll(testSliceCollection[int]{1, 2})
				if !s.Empty() {
					t.Errorf("expected empty set to remain empty")
				}

				s.AddAll(slices.Values([]int{1, 2, 3}))
				s.RetainAll(s)
				if s.Size() != 3 {
					t.Errorf("expected RetainAll(self) to preserve size 3, got %d", s.Size())
				}

				s.RetainAll(testSliceCollection[int]{2, 2, 3, 4})
				if !slices.Equal(slices.Collect(s.All()), []int{2, 3}) {
					t.Errorf("expected [2 3], got %s", s.String())
				}
			},
		},
		{
			name: "string_empty_single_and_multiple",
			check: func(t *testing.T, s *TreeSet[int]) {
				if got := s.String(); got != "[]" {
					t.Errorf("expected [], got %s", got)
				}
				s.Add(2)
				if got := s.String(); got != "[2]" {
					t.Errorf("expected [2], got %s", got)
				}
				s.Add(1)
				s.Add(3)
				if got := s.String(); got != "[1 2 3]" {
					t.Errorf("expected [1 2 3], got %s", got)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := NewOrdered[int]()
			tc.check(t, s)
		})
	}
}

func TestTreeSet_SequencedNavigableAndIterators(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		wantPanic bool
		check     func(*testing.T, *TreeSet[int])
	}{
		{
			name:      "first_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, s *TreeSet[int]) {
				s.First()
			},
		},
		{
			name:      "last_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, s *TreeSet[int]) {
				s.Last()
			},
		},
		{
			name:      "poll_first_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, s *TreeSet[int]) {
				s.PollFirst()
			},
		},
		{
			name:      "poll_last_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, s *TreeSet[int]) {
				s.PollLast()
			},
		},
		{
			name: "first_last_poll_first_poll_last",
			check: func(t *testing.T, s *TreeSet[int]) {
				s.AddAll(slices.Values([]int{10, 20, 30}))
				if v := s.First(); v != 10 {
					t.Errorf("First() = %d, want 10", v)
				}
				if v := s.Last(); v != 30 {
					t.Errorf("Last() = %d, want 30", v)
				}
				if v := s.PollFirst(); v != 10 || s.Contains(10) {
					t.Errorf("PollFirst() = %d", v)
				}
				if v := s.PollLast(); v != 30 || s.Contains(30) {
					t.Errorf("PollLast() = %d", v)
				}
			},
		},
		{
			name: "navigable_lower_floor_ceiling_higher",
			check: func(t *testing.T, s *TreeSet[int]) {
				if _, ok := s.Lower(10); ok {
					t.Errorf("expected false on empty Lower")
				}
				s.AddAll(slices.Values([]int{10, 20, 30, 40, 50}))

				if k, ok := s.Lower(30); !ok || k != 20 {
					t.Errorf("Lower(30) = (%d, %v)", k, ok)
				}
				if _, ok := s.Lower(10); ok {
					t.Errorf("Lower(10) expected false")
				}
				if k, ok := s.Floor(30); !ok || k != 30 {
					t.Errorf("Floor(30) = (%d, %v)", k, ok)
				}
				if k, ok := s.Floor(25); !ok || k != 20 {
					t.Errorf("Floor(25) = (%d, %v)", k, ok)
				}
				if _, ok := s.Floor(5); ok {
					t.Errorf("Floor(5) expected false")
				}
				if k, ok := s.Ceiling(30); !ok || k != 30 {
					t.Errorf("Ceiling(30) = (%d, %v)", k, ok)
				}
				if k, ok := s.Ceiling(35); !ok || k != 40 {
					t.Errorf("Ceiling(35) = (%d, %v)", k, ok)
				}
				if _, ok := s.Ceiling(60); ok {
					t.Errorf("Ceiling(60) expected false")
				}
				if k, ok := s.Higher(30); !ok || k != 40 {
					t.Errorf("Higher(30) = (%d, %v)", k, ok)
				}
				if _, ok := s.Higher(50); ok {
					t.Errorf("Higher(50) expected false")
				}
			},
		},
		{
			name: "all_backward_from_to_between_and_early_breaks",
			check: func(t *testing.T, s *TreeSet[int]) {
				for i := 1; i <= 10; i++ {
					s.Add(i)
				}
				if got := slices.Collect(s.Backward()); !slices.Equal(got, []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}) {
					t.Errorf("Backward() = %v", got)
				}
				if got := slices.Collect(s.From(5)); !slices.Equal(got, []int{5, 6, 7, 8, 9, 10}) {
					t.Errorf("From(5) = %v", got)
				}
				if got := slices.Collect(s.To(5)); !slices.Equal(got, []int{1, 2, 3, 4}) {
					t.Errorf("To(5) = %v", got)
				}
				if got := slices.Collect(s.Between(4, 8)); !slices.Equal(got, []int{4, 5, 6, 7}) {
					t.Errorf("Between(4, 8) = %v", got)
				}

				for _, seq := range []iter.Seq[int]{
					s.All(),
					s.Backward(),
					s.From(2),
					s.To(9),
					s.Between(2, 9),
				} {
					count := 0
					for range seq {
						count++
						break
					}
					if count != 1 {
						t.Errorf("expected early break count 1, got %d", count)
					}
				}
			},
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
			s := NewOrdered[int](WithDegree[int](2))
			tc.check(t, s)
		})
	}
}
