package arraylist

import (
	"slices"
	"testing"
)

func TestNewAndWrap(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		build     func() *SliceWrapper[int]
		wantPanic bool
		check     func(*testing.T, *SliceWrapper[int])
	}{
		{
			name:  "default_new",
			build: func() *SliceWrapper[int] { return New[int]() },
			check: func(t *testing.T, l *SliceWrapper[int]) {
				if !l.Empty() || l.Size() != 0 || l.String() != "[]" {
					t.Errorf("expected empty list")
				}
				if cap(l.slice) != 0 {
					t.Errorf("expected cap 0, got %d", cap(l.slice))
				}
			},
		},
		{
			name:  "with_capacity",
			build: func() *SliceWrapper[int] { return New[int](WithCapacity(10)) },
			check: func(t *testing.T, l *SliceWrapper[int]) {
				if !l.Empty() || l.Size() != 0 {
					t.Errorf("expected empty list")
				}
				if cap(l.slice) != 10 {
					t.Errorf("expected cap 10, got %d", cap(l.slice))
				}
			},
		},
		{
			name:      "negative_capacity_panics",
			build:     func() *SliceWrapper[int] { return New[int](WithCapacity(-1)) },
			wantPanic: true,
		},
		{
			name:  "wrap_nil_slice",
			build: func() *SliceWrapper[int] { return Wrap[int](nil) },
			check: func(t *testing.T, l *SliceWrapper[int]) {
				if !l.Empty() || l.Size() != 0 {
					t.Errorf("expected empty list")
				}
				l.Add(42)
				if l.Size() != 1 || l.Get(0) != 42 {
					t.Errorf("Wrap(nil) failed after Add")
				}
			},
		},
		{
			name:  "wrap_existing_slice",
			build: func() *SliceWrapper[int] { return Wrap([]int{1, 2, 3}) },
			check: func(t *testing.T, l *SliceWrapper[int]) {
				if l.Size() != 3 || l.Get(0) != 1 || l.Get(1) != 2 || l.Get(2) != 3 {
					t.Errorf("wrong items")
				}
			},
		},
		{
			name: "zero_value_struct",
			build: func() *SliceWrapper[int] {
				var z SliceWrapper[int]
				return &z
			},
			check: func(t *testing.T, l *SliceWrapper[int]) {
				if !l.Empty() || l.Size() != 0 || l.String() != "[]" {
					t.Errorf("expected empty zero-value list")
				}
				l.Add(7)
				if l.Size() != 1 || l.Peek() != 7 {
					t.Errorf("zero-value failed after Add")
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				_ = tc.build()
				return
			}
			l := tc.build()
			tc.check(t, l)
		})
	}
}

func TestSliceWrapper_MutatorsAndPanics(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		wantPanic bool
		check     func(*testing.T, *SliceWrapper[int])
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
			name: "get_and_set",
			check: func(t *testing.T, l *SliceWrapper[int]) {
				l.AddAll(slices.Values([]int{1, 2, 3}))
				if got := l.Get(1); got != 2 {
					t.Errorf("Get(1) = %d, want 2", got)
				}
				l.Set(1, 99)
				if got := l.Get(1); got != 99 {
					t.Errorf("Get(1) = %d, want 99", got)
				}
			},
		},
		{
			name:      "remove_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *SliceWrapper[int]) { l.Remove() },
		},
		{
			name:      "pop_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *SliceWrapper[int]) { l.Pop() },
		},
		{
			name:      "peek_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *SliceWrapper[int]) { l.Peek() },
		},
		{
			name:      "get_negative_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *SliceWrapper[int]) { _ = l.Get(-1) },
		},
		{
			name:      "get_out_of_bounds_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *SliceWrapper[int]) { _ = l.Get(0) },
		},
		{
			name:      "set_out_of_bounds_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *SliceWrapper[int]) { l.Set(0, 1) },
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := New[int]()
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
			}
			tc.check(t, l)
		})
	}
}

func TestSliceWrapper_BulkIteratorsAndClearGC(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T)
	}{
		{
			name: "add_all_and_string",
			check: func(t *testing.T) {
				l := New[int]()
				if l.String() != "[]" {
					t.Errorf("expected [], got %s", l.String())
				}
				l.AddAll(slices.Values([]int{1, 2, 3}))
				if str := l.String(); str != "[1 2 3]" {
					t.Errorf("expected [1 2 3], got %s", str)
				}
			},
		},
		{
			name: "all_full_and_early_break",
			check: func(t *testing.T) {
				l := Wrap([]int{10, 20, 30, 40})
				if got := slices.Collect(l.All()); !slices.Equal(got, []int{10, 20, 30, 40}) {
					t.Errorf("got %v", got)
				}
				var partial []int
				for v := range l.All() {
					partial = append(partial, v)
					if v == 20 {
						break
					}
				}
				if !slices.Equal(partial, []int{10, 20}) {
					t.Errorf("early break got %v", partial)
				}
			},
		},
		{
			name: "clear_and_remove_zero_backing_pointers",
			check: func(t *testing.T) {
				a, b, c := 1, 2, 3
				raw := []*int{&a, &b, &c}
				l := Wrap(raw)
				_ = l.Remove()
				if raw[2] != nil {
					t.Errorf("Remove did not zero vacated pointer slot")
				}
				l.Clear()
				if !l.Empty() || raw[0] != nil || raw[1] != nil {
					t.Errorf("Clear did not zero backing pointer slots")
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
