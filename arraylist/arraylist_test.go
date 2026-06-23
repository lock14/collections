package arraylist

import (
	"slices"
	"testing"
)

func TestWrap(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		slice []int
		check func(*testing.T, *SliceWrapper[int])
	}{
		{
			name:  "nil_slice",
			slice: nil,
			check: func(t *testing.T, l *SliceWrapper[int]) {
				if !l.Empty() || l.Size() != 0 {
					t.Errorf("expected empty list")
				}
			},
		},
		{
			name:  "existing_slice",
			slice: []int{1, 2, 3},
			check: func(t *testing.T, l *SliceWrapper[int]) {
				if l.Size() != 3 {
					t.Errorf("expected size 3")
				}
				if l.Get(0) != 1 || l.Get(1) != 2 || l.Get(2) != 3 {
					t.Errorf("wrong items")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := Wrap(tc.slice)
			tc.check(t, l)
		})
	}
}

func TestSliceWrapper_AddRemove(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *SliceWrapper[int])
	}{
		{
			name: "add_and_remove",
			check: func(t *testing.T, l *SliceWrapper[int]) {
				l.Add(1)
				l.Add(2)
				if l.Size() != 2 {
					t.Errorf("expected size 2")
				}
				if val := l.Remove(); val != 2 {
					t.Errorf("expected 2, got %v", val)
				}
				if val := l.Remove(); val != 1 {
					t.Errorf("expected 1, got %v", val)
				}
				if !l.Empty() {
					t.Errorf("expected empty list")
				}
			},
		},
		{
			name: "remove_empty_panic",
			check: func(t *testing.T, l *SliceWrapper[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				l.Remove()
			},
		},
		{
			name: "stack_methods",
			check: func(t *testing.T, l *SliceWrapper[int]) {
				l.Push(1)
				l.Push(2)
				if l.Peek() != 2 {
					t.Errorf("expected 2")
				}
				if val := l.Pop(); val != 2 {
					t.Errorf("expected 2")
				}
				if l.Peek() != 1 {
					t.Errorf("expected 1")
				}
			},
		},
		{
			name: "peek_empty_panic",
			check: func(t *testing.T, l *SliceWrapper[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				l.Peek()
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := Wrap([]int{})
			tc.check(t, l)
		})
	}
}

func TestSliceWrapper_GetSet(t *testing.T) {
	t.Parallel()
	l := Wrap([]int{1, 2, 3})
	if l.Get(1) != 2 {
		t.Errorf("expected 2")
	}
	l.Set(1, 4)
	if l.Get(1) != 4 {
		t.Errorf("expected 4")
	}
}

func TestSliceWrapper_BulkOperations(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *SliceWrapper[int])
	}{
		{
			name: "add_all",
			check: func(t *testing.T, l *SliceWrapper[int]) {
				l.AddAll(slices.Values([]int{1, 2, 3}))
				if l.Size() != 3 {
					t.Errorf("expected 3")
				}
				if str := l.String(); str != "[1, 2, 3]" {
					t.Errorf("expected [1, 2, 3], got %s", str)
				}
			},
		},
		{
			name: "iterators",
			check: func(t *testing.T, l *SliceWrapper[int]) {
				l.AddAll(slices.Values([]int{1, 2, 3}))
				got := slices.Collect(l.All())
				if !slices.Equal(got, []int{1, 2, 3}) {
					t.Errorf("expected [1, 2, 3], got %v", got)
				}
			},
		},
		{
			name: "clear",
			check: func(t *testing.T, l *SliceWrapper[int]) {
				l.AddAll(slices.Values([]int{1, 2, 3}))
				l.Clear()
				if !l.Empty() {
					t.Errorf("expected empty")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := Wrap([]int{})
			tc.check(t, l)
		})
	}
}
