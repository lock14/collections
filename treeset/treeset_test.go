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
				if str != "[1 2]" {
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
		{
			name: "custom degree and comparator",
			constructor: func() {
				s := New[int](
					WithDegree[int](4),
					WithComparator[int](func(a, b int) int {
						return a - b
					}),
				)
				s.Add(1)
				s.Add(2)
				if s.Size() != 2 || !s.Contains(1) {
					panic("unexpected set state")
				}
			},
			expectPanic: false,
		},
		{
			name: "NewOrdered with degree",
			constructor: func() {
				s := NewOrdered[int](WithDegree[int](3))
				s.Add(10)
				if !s.Contains(10) {
					panic("unexpected set state")
				}
			},
			expectPanic: false,
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

func assertKV(t *testing.T, k int, ok bool, expectedK int, expectedOk bool, name string) {
	t.Helper()
	if ok != expectedOk {
		t.Errorf("%s: expected ok=%v, got %v", name, expectedOk, ok)
	} else if ok && k != expectedK {
		t.Errorf("%s: expected %v, got %v", name, expectedK, k)
	}
}

func TestTreeSet_Navigable(t *testing.T) {
	t.Parallel()
	s := NewOrdered[int]()

	_, ok := s.Lower(10)
	if ok {
		t.Errorf("expected empty set to return false for Lower")
	}

	s.Add(10)
	s.Add(20)
	s.Add(30)
	s.Add(40)
	s.Add(50)

	k, ok := s.Lower(30)
	assertKV(t, k, ok, 20, true, "Lower(30)")
	k, ok = s.Lower(10)
	assertKV(t, k, ok, 0, false, "Lower(10)")

	k, ok = s.Floor(30)
	assertKV(t, k, ok, 30, true, "Floor(30)")
	k, ok = s.Floor(25)
	assertKV(t, k, ok, 20, true, "Floor(25)")
	k, ok = s.Floor(5)
	assertKV(t, k, ok, 0, false, "Floor(5)")

	k, ok = s.Ceiling(30)
	assertKV(t, k, ok, 30, true, "Ceiling(30)")
	k, ok = s.Ceiling(35)
	assertKV(t, k, ok, 40, true, "Ceiling(35)")
	k, ok = s.Ceiling(60)
	assertKV(t, k, ok, 0, false, "Ceiling(60)")

	k, ok = s.Higher(30)
	assertKV(t, k, ok, 40, true, "Higher(30)")
	k, ok = s.Higher(50)
	assertKV(t, k, ok, 0, false, "Higher(50)")
}

func TestTreeSet_Sequenced(t *testing.T) {
	t.Parallel()

	// Test panics
	assertPanic := func(name string, f func()) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s did not panic", name)
			}
		}()
		f()
	}

	s := NewOrdered[int]()
	assertPanic("First", func() { s.First() })
	assertPanic("Last", func() { s.Last() })
	assertPanic("PollFirst", func() { s.PollFirst() })
	assertPanic("PollLast", func() { s.PollLast() })

	s.Add(10)
	s.Add(20)
	s.Add(30)

	v := s.First()
	if v != 10 {
		t.Errorf("First: expected 10, got %v", v)
	}

	v = s.Last()
	if v != 30 {
		t.Errorf("Last: expected 30, got %v", v)
	}

	v = s.PollFirst()
	if v != 10 {
		t.Errorf("PollFirst: expected 10, got %v", v)
	}
	if s.Contains(10) {
		t.Errorf("PollFirst didn't remove")
	}

	v = s.PollLast()
	if v != 30 {
		t.Errorf("PollLast: expected 30, got %v", v)
	}
	if s.Contains(30) {
		t.Errorf("PollLast didn't remove")
	}

	assertPanic("AddFirst", func() { s.AddFirst(1) })
	assertPanic("AddLast", func() { s.AddLast(1) })
}

func TestTreeSet_Backward(t *testing.T) {
	t.Parallel()
	s := NewOrdered[int]()
	for i := 1; i <= 5; i++ {
		s.Add(i * 10)
	}

	var keys []int
	for k := range s.Backward() {
		keys = append(keys, k)
	}
	if !slices.Equal(keys, []int{50, 40, 30, 20, 10}) {
		t.Errorf("Backward: %v", keys)
	}

	// coverage for Backward early exit
	for k := range s.Backward() {
		_ = k
		break
	}
}

func TestTreeSet_Iterators_Bounds(t *testing.T) {
	t.Parallel()
	s := NewOrdered[int]()
	for i := 1; i <= 10; i++ {
		s.Add(i)
	}

	if got := slices.Collect(s.From(5)); !slices.Equal(got, []int{5, 6, 7, 8, 9, 10}) {
		t.Errorf("From(5): %v", got)
	}
	if got := slices.Collect(s.To(5)); !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Errorf("To(5): %v", got)
	}
	if got := slices.Collect(s.Between(4, 8)); !slices.Equal(got, []int{4, 5, 6, 7}) {
		t.Errorf("Between(4, 8): %v", got)
	}

	// Early exit coverage
	for k := range s.Between(1, 10) {
		_ = k
		break
	}
}
