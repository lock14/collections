package linkedhashset

import (
	"iter"
	"slices"
	"testing"
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

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		opts      []Option
		wantPanic bool
		check     func(*testing.T, *LinkedHashSet[int])
	}{
		{
			name: "default",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				if s.Size() != 0 || !s.Empty() {
					t.Errorf("expected empty set, got size %d", s.Size())
				}
			},
		},
		{
			name: "capacity_zero_and_positive",
			opts: []Option{WithCapacity(0), WithCapacity(100)},
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				if s.Size() != 1 {
					t.Errorf("expected size 1, got %d", s.Size())
				}
			},
		},
		{
			name: "insertion_order",
			opts: []Option{WithInsertionOrder()},
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				if s.Size() != 1 {
					t.Errorf("expected size 1, got %d", s.Size())
				}
			},
		},
		{
			name: "access_order",
			opts: []Option{WithAccessOrder()},
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				if s.Size() != 1 {
					t.Errorf("expected size 1, got %d", s.Size())
				}
			},
		},
		{
			name: "max_elements_eviction",
			opts: []Option{WithMaxElements(2)},
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(3)
				if s.Size() != 2 || s.Contains(1) {
					t.Errorf("expected element 1 evicted and size 2, got %s", s.String())
				}
			},
		},
		{
			name:      "negative_capacity_panics",
			opts:      []Option{WithCapacity(-1)},
			wantPanic: true,
		},
		{
			name:      "zero_max_elements_panics",
			opts:      []Option{WithMaxElements(0)},
			wantPanic: true,
		},
		{
			name:      "negative_max_elements_panics",
			opts:      []Option{WithMaxElements(-1)},
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

func TestLinkedHashSet_Order(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Option
		setup func(*LinkedHashSet[int])
		want  []int
	}{
		{
			name: "insertion_order_default",
			setup: func(s *LinkedHashSet[int]) {
				s.Add(3)
				s.Add(1)
				s.Add(2)
				s.Add(1)
			},
			want: []int{3, 1, 2},
		},
		{
			name: "insertion_order_explicit",
			opts: []Option{WithInsertionOrder()},
			setup: func(s *LinkedHashSet[int]) {
				s.Add(3)
				s.Add(1)
				s.Add(2)
				s.Add(1)
			},
			want: []int{3, 1, 2},
		},
		{
			name: "access_order_moves_to_back_on_reinsert",
			opts: []Option{WithAccessOrder()},
			setup: func(s *LinkedHashSet[int]) {
				s.Add(3)
				s.Add(1)
				s.Add(2)
				s.Add(1)
			},
			want: []int{3, 2, 1},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := New[int](tc.opts...)
			tc.setup(s)
			if got := slices.Collect(s.All()); !slices.Equal(got, tc.want) {
				t.Errorf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestLinkedHashSet_AddRemoveContains(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *LinkedHashSet[int])
	}{
		{
			name: "add_and_contains",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(2)
				if s.Size() != 2 || !s.Contains(1) || !s.Contains(2) || s.Contains(3) {
					t.Errorf("add/contains failed: %s", s.String())
				}
			},
		},
		{
			name: "remove_element_present_and_absent",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.RemoveElement(1)
				s.RemoveElement(99)
				if s.Contains(1) || s.Size() != 1 {
					t.Errorf("remove element failed: %s", s.String())
				}
			},
		},
		{
			name: "remove_returns_first_in_order",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(10)
				s.Add(20)
				if val := s.Remove(); val != 10 {
					t.Errorf("expected Remove() to return first element 10, got %d", val)
				}
				if s.Size() != 1 || !s.Contains(20) {
					t.Errorf("unexpected remaining set: %s", s.String())
				}
			},
		},
		{
			name: "remove_empty_panics",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
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
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
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

func TestLinkedHashSet_BulkAndString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *LinkedHashSet[int])
	}{
		{
			name: "add_all",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				if s.Size() != 3 || !s.Contains(1) {
					t.Errorf("AddAll failed: %s", s.String())
				}
			},
		},
		{
			name: "contains_all_including_self_alias",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				if !s.ContainsAll(s) {
					t.Errorf("expected self ContainsAll to return true")
				}
				other := New[int]()
				other.Add(1)
				other.Add(2)
				if !s.ContainsAll(other) {
					t.Errorf("ContainsAll subset failed")
				}
				other.Add(4)
				if s.ContainsAll(other) {
					t.Errorf("ContainsAll non-subset should return false")
				}
			},
		},
		{
			name: "remove_all_empty_normal_and_self_alias",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				other := New[int]()
				other.Add(1)
				s.RemoveAll(other)
				if !s.Empty() {
					t.Errorf("expected empty set to remain empty")
				}

				s.AddAll(slices.Values([]int{1, 2, 3}))
				other.Add(2)
				s.RemoveAll(other)
				if s.Contains(1) || s.Contains(2) || !s.Contains(3) {
					t.Errorf("RemoveAll failed: %s", s.String())
				}

				s.Add(4)
				s.RemoveAll(s)
				if !s.Empty() {
					t.Errorf("expected empty set after RemoveAll(self)")
				}
			},
		},
		{
			name: "retain_all_empty_set_collection_and_self_alias",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				other := New[int]()
				other.Add(2)
				s.RetainAll(other)
				if !s.Empty() {
					t.Errorf("expected empty set to remain empty")
				}

				s.AddAll(slices.Values([]int{1, 2, 3}))
				s.RetainAll(s)
				if s.Size() != 3 {
					t.Errorf("expected RetainAll(self) to preserve size 3, got %d", s.Size())
				}

				other.Add(3)
				other.Add(4)
				s.RetainAll(other)
				if !slices.Equal(slices.Collect(s.All()), []int{2, 3}) {
					t.Errorf("RetainAll Set failed: %s", s.String())
				}

				s.AddAll(slices.Values([]int{10, 20}))
				s.RetainAll(testSliceCollection[int]{3, 20, 20, 99})
				if !slices.Equal(slices.Collect(s.All()), []int{3, 20}) {
					t.Errorf("RetainAll non-Set Collection failed: %s", s.String())
				}
			},
		},
		{
			name: "string_empty_single_and_multiple",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				if got := s.String(); got != "[]" {
					t.Errorf("expected [], got %s", got)
				}
				s.Add(3)
				if got := s.String(); got != "[3]" {
					t.Errorf("expected [3], got %s", got)
				}
				s.Add(1)
				s.Add(2)
				if got := s.String(); got != "[3 1 2]" {
					t.Errorf("expected [3 1 2], got %s", got)
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

func TestLinkedHashSet_SequencedAndBackward(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		opts      []Option
		wantPanic bool
		check     func(*testing.T, *LinkedHashSet[int])
	}{
		{
			name:      "first_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, s *LinkedHashSet[int]) {
				s.First()
			},
		},
		{
			name:      "last_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, s *LinkedHashSet[int]) {
				s.Last()
			},
		},
		{
			name:      "poll_first_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, s *LinkedHashSet[int]) {
				s.PollFirst()
			},
		},
		{
			name:      "poll_last_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, s *LinkedHashSet[int]) {
				s.PollLast()
			},
		},
		{
			name: "first_last_poll_first_poll_last",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
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
			name: "add_first_and_add_last_insert_and_reposition",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.AddFirst(2)
				s.AddFirst(1)
				s.AddLast(3)
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{1, 2, 3}) {
					t.Errorf("expected [1 2 3], got %v", got)
				}
				s.AddFirst(3)
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{3, 1, 2}) {
					t.Errorf("expected [3 1 2], got %v", got)
				}
				s.AddLast(1)
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{3, 2, 1}) {
					t.Errorf("expected [3 2 1], got %v", got)
				}
			},
		},
		{
			name: "add_first_and_add_last_eviction",
			opts: []Option{WithMaxElements(2)},
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.AddFirst(1)
				s.AddFirst(2)
				s.AddFirst(3) // Evicts tail (1)
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{3, 2}) {
					t.Errorf("expected [3 2], got %v", got)
				}
				s.AddLast(4) // Evicts head (3)
				if got := slices.Collect(s.All()); !slices.Equal(got, []int{2, 4}) {
					t.Errorf("expected [2 4], got %v", got)
				}
			},
		},
		{
			name: "all_and_backward_with_early_break",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				if got := slices.Collect(s.Backward()); !slices.Equal(got, []int{3, 2, 1}) {
					t.Errorf("expected Backward [3 2 1], got %v", got)
				}

				fCount := 0
				for range s.All() {
					fCount++
					break
				}
				bCount := 0
				for range s.Backward() {
					bCount++
					break
				}
				if fCount != 1 || bCount != 1 {
					t.Errorf("expected early break counts 1/1, got %d/%d", fCount, bCount)
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
			s := New[int](tc.opts...)
			tc.check(t, s)
		})
	}
}
