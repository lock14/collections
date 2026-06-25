package linkedhashset

import (
	"github.com/lock14/collections/arraylist"
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *LinkedHashSet[int])
	}{
		{
			name: "default",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				if s.Size() != 0 || !s.Empty() {
					t.Errorf("expected empty set")
				}
			},
		},
		{
			name: "capacity",
			opts: []Option{WithCapacity(100)},
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				if s.Size() != 1 {
					t.Errorf("expected size 1")
				}
			},
		},
		{
			name: "insertion_order",
			opts: []Option{WithInsertionOrder()},
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				if s.Size() != 1 {
					t.Errorf("expected size 1")
				}
			},
		},
		{
			name: "access_order",
			opts: []Option{WithAccessOrder()},
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				if s.Size() != 1 {
					t.Errorf("expected size 1")
				}
			},
		},
		{
			name: "max_elements",
			opts: []Option{WithMaxElements(2)},
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(3)
				if s.Size() != 2 {
					t.Errorf("expected size 2, got %v", s.Size())
				}
				if s.Contains(1) {
					t.Errorf("expected element 1 to be evicted")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := New[int](tc.opts...)
			tc.check(t, s)
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
			name: "insertion order (default)",
			opts: nil,
			setup: func(s *LinkedHashSet[int]) {
				s.Add(3)
				s.Add(1)
				s.Add(2)
				s.Add(1) // re-insert shouldn't change order if insertion order
			},
			want: []int{3, 1, 2},
		},
		{
			name: "insertion order (explicit)",
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
			name: "access order",
			opts: []Option{WithAccessOrder()},
			setup: func(s *LinkedHashSet[int]) {
				s.Add(3)
				s.Add(1)
				s.Add(2)
				s.Add(1) // re-insert should move 1 to the end
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
			got := slices.Collect(s.All())
			if !slices.Equal(got, tc.want) {
				t.Errorf("expected order %v, got %v", tc.want, got)
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
				s.Add(2) // duplicate
				if s.Size() != 2 {
					t.Errorf("expected size 2")
				}
				if !s.Contains(1) || !s.Contains(2) || s.Contains(3) {
					t.Errorf("contains failed")
				}
			},
		},
		{
			name: "remove_element",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.RemoveElement(1)
				if s.Contains(1) || s.Size() != 1 {
					t.Errorf("remove failed")
				}
			},
		},
		{
			name: "remove_arbitrary",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(1)
				val := s.Remove()
				if val != 1 {
					t.Errorf("expected to remove 1")
				}
				if !s.Empty() {
					t.Errorf("expected empty")
				}
			},
		},
		{
			name: "remove_empty_panic",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
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
				if !s.Empty() {
					t.Errorf("expected empty")
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

func TestLinkedHashSet_Bulk(t *testing.T) {
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
					t.Errorf("AddAll failed")
				}
			},
		},
		{
			name: "contains_all",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				other := New[int]()
				other.Add(1)
				other.Add(2)
				if !s.ContainsAll(other) {
					t.Errorf("ContainsAll failed")
				}
				other.Add(4)
				if s.ContainsAll(other) {
					t.Errorf("ContainsAll should fail")
				}
			},
		},
		{
			name: "remove_all_set",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				other := New[int]()
				other.Add(1)
				other.Add(2)
				s.RemoveAll(other)
				if s.Contains(1) || s.Contains(2) || !s.Contains(3) {
					t.Errorf("RemoveAll failed")
				}
			},
		},
		{
			name: "retain_all_set",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				other := New[int]()
				other.Add(2)
				other.Add(3)
				other.Add(4)
				s.RetainAll(other)
				if s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
					t.Errorf("RetainAll failed")
				}
			},
		},
		{
			name: "retain_all_generic",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				list := arraylist.Wrap([]int{2, 3, 4})
				s.RetainAll(list)
				if s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
					t.Errorf("RetainAll generic failed")
				}
			},
		},
		{
			name: "string",
			check: func(t *testing.T, s *LinkedHashSet[int]) {
				s.Add(3)
				s.Add(1)
				s.Add(2)
				if s.String() != "[3 1 2]" {
					t.Errorf("expected string [3 1 2], got %s", s.String())
				}
				s.Clear()
				if s.String() != "[]" {
					t.Errorf("expected string [], got %s", s.String())
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
