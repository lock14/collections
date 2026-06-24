package treeset

import (
	"slices"
	"testing"

	"github.com/lock14/collections/arraylist"
)

func TestTreeSet_Operations(t *testing.T) {
	cases := []struct {
		name     string
		ops      func(s *TreeSet[int])
		validate func(t *testing.T, s *TreeSet[int])
	}{
		{
			name: "basic add and contains",
			ops: func(s *TreeSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(3)
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				if s.Size() != 3 {
					t.Errorf("expected size 3, got %d", s.Size())
				}
				if !s.Contains(2) {
					t.Errorf("expected to contain 2")
				}
				if s.Contains(4) {
					t.Errorf("expected not to contain 4")
				}
			},
		},
		{
			name: "add duplicate elements",
			ops: func(s *TreeSet[int]) {
				s.Add(1)
				s.Add(1)
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				if s.Size() != 1 {
					t.Errorf("expected size 1, got %d", s.Size())
				}
			},
		},
		{
			name: "remove element",
			ops: func(s *TreeSet[int]) {
				s.Add(1)
				s.Add(2)
				s.RemoveElement(1)
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				if s.Contains(1) {
					t.Errorf("expected not to contain 1")
				}
				if s.Size() != 1 {
					t.Errorf("expected size 1, got %d", s.Size())
				}
			},
		},
		{
			name: "remove single element",
			ops: func(s *TreeSet[int]) {
				s.Add(1)
				val := s.Remove()
				if val != 1 {
					t.Errorf("expected to remove 1, got %d", val)
				}
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				if !s.Empty() {
					t.Errorf("expected set to be empty")
				}
			},
		},
		{
			name: "clear set",
			ops: func(s *TreeSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Clear()
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				if !s.Empty() || s.Size() != 0 {
					t.Errorf("expected set to be empty")
				}
			},
		},
		{
			name: "contains all",
			ops: func(s *TreeSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(3)
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				other := arraylist.Wrap([]int{})
				other.Add(1)
				other.Add(2)
				if !s.ContainsAll(other) {
					t.Errorf("expected to contain all elements")
				}
				other.Add(4)
				if s.ContainsAll(other) {
					t.Errorf("expected not to contain 4")
				}
			},
		},
		{
			name: "add all",
			ops: func(s *TreeSet[int]) {
				other := arraylist.Wrap([]int{})
				other.Add(1)
				other.Add(2)
				s.AddAll(other.All())
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				if s.Size() != 2 {
					t.Errorf("expected size 2, got %d", s.Size())
				}
				if !s.Contains(1) || !s.Contains(2) {
					t.Errorf("expected to contain 1 and 2")
				}
			},
		},
		{
			name: "remove all",
			ops: func(s *TreeSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(3)
				other := arraylist.Wrap([]int{})
				other.Add(2)
				other.Add(3)
				s.RemoveAll(other)
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				if s.Size() != 1 {
					t.Errorf("expected size 1, got %d", s.Size())
				}
				if !s.Contains(1) {
					t.Errorf("expected to contain 1")
				}
			},
		},
		{
			name: "retain all",
			ops: func(s *TreeSet[int]) {
				s.Add(1)
				s.Add(2)
				s.Add(3)
				other := arraylist.Wrap([]int{})
				other.Add(2)
				other.Add(3)
				other.Add(4)
				s.RetainAll(other)
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				if s.Size() != 2 {
					t.Errorf("expected size 2, got %d", s.Size())
				}
				if !s.Contains(2) || !s.Contains(3) {
					t.Errorf("expected to retain 2 and 3")
				}
			},
		},
		{
			name: "all iterator",
			ops: func(s *TreeSet[int]) {
				s.Add(3)
				s.Add(1)
				s.Add(2)
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				var elements []int
				for e := range s.All() {
					elements = append(elements, e)
				}
				expected := []int{1, 2, 3}
				if !slices.Equal(expected, elements) {
					t.Errorf("expected %v, got %v", expected, elements)
				}
			},
		},
		{
			name: "string representation",
			ops: func(s *TreeSet[int]) {
				s.Add(2)
				s.Add(1)
			},
			validate: func(t *testing.T, s *TreeSet[int]) {
				str := s.String()
				// Should be ordered
				if str != "[1, 2]" {
					t.Errorf("unexpected string representation: %s", str)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := NewOrdered[int]()
			tc.ops(s)
			tc.validate(t, s)
		})
	}
}

func TestTreeSet_Constructors(t *testing.T) {
	cases := []struct {
		name        string
		constructor func()
		expectPanic bool
	}{
		{
			name: "remove from empty panics",
			constructor: func() {
				s := NewOrdered[int]()
				s.Remove()
			},
			expectPanic: true,
		},
		{
			name: "missing comparator panics",
			constructor: func() {
				New[int]()
			},
			expectPanic: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				r := recover()
				if tc.expectPanic && r == nil {
					t.Errorf("expected panic but got none")
				}
				if !tc.expectPanic && r != nil {
					t.Errorf("expected no panic but got: %v", r)
				}
			}()
			tc.constructor()
		})
	}
}
