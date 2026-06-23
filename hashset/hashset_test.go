package hashset

import (
	"slices"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *HashSet[int])
	}{
		{
			name: "default",
			check: func(t *testing.T, s *HashSet[int]) {
				if s.Size() != 0 || !s.Empty() {
					t.Errorf("expected empty set")
				}
			},
		},
		{
			name: "capacity",
			opts: []Option{WithCapacity(100)},
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(1)
				if s.Size() != 1 {
					t.Errorf("expected size 1")
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
			check: func(t *testing.T, s *HashSet[int]) {
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
			check: func(t *testing.T, s *HashSet[int]) {
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
			check: func(t *testing.T, s *HashSet[int]) {
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
			check: func(t *testing.T, s *HashSet[int]) {
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

func TestHashSet_BulkOperations(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *HashSet[int])
	}{
		{
			name: "add_all",
			check: func(t *testing.T, s *HashSet[int]) {
				s.AddAll(slices.Values([]int{1, 2, 3}))
				if s.Size() != 3 {
					t.Errorf("expected 3")
				}
			},
		},
		{
			name: "contains_all",
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(1)
				s.Add(2)
				
				other := New[int]()
				other.Add(1)
				
				if !s.ContainsAll(other) {
					t.Errorf("expected true")
				}
				
				other.Add(3)
				if s.ContainsAll(other) {
					t.Errorf("expected false")
				}
			},
		},
		{
			name: "remove_all",
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(3)
				
				other := New[int]()
				other.Add(2)
				other.Add(3)
				
				s.RemoveAll(other)
				if s.Size() != 1 || !s.Contains(1) {
					t.Errorf("expected only 1 to remain")
				}
			},
		},
		{
			name: "retain_all",
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(3)
				
				other := New[int]()
				other.Add(2)
				other.Add(4)
				
				s.RetainAll(other)
				if s.Size() != 1 || !s.Contains(2) {
					t.Errorf("expected only 2 to remain")
				}
			},
		},
		{
			name: "iterators",
			check: func(t *testing.T, s *HashSet[int]) {
				s.Add(1)
				s.Add(2)
				got := slices.Collect(s.All())
				sort.Ints(got)
				if !slices.Equal(got, []int{1, 2}) {
					t.Errorf("wrong iter results")
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
