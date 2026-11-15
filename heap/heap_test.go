package heap

import (
	"slices"
	"testing"

	"github.com/lock14/collections"
	"github.com/lock14/collections/comparator"
)

func TestNewConstructors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		build     func() *Heap[int]
		wantPanic bool
		check     func(*testing.T, *Heap[int])
	}{
		{
			name:      "missing_comparator_panics",
			build:     func() *Heap[int] { return New[int]() },
			wantPanic: true,
		},
		{
			name:      "explicit_nil_comparator_panics",
			build:     func() *Heap[int] { return NewOrdered[int](WithComparator[int](nil)) },
			wantPanic: true,
		},
		{
			name:      "negative_capacity_panics",
			build:     func() *Heap[int] { return Min[int](WithCapacity[int](-1)) },
			wantPanic: true,
		},
		{
			name:  "default_capacity",
			build: func() *Heap[int] { return New[int](WithComparator(comparator.NaturalOrder[int]())) },
			check: func(t *testing.T, h *Heap[int]) {
				var _ collections.MutableQueue[int] = h
				if h.Size() != 0 || !h.Empty() {
					t.Errorf("expected empty heap")
				}
				if cap(h.elements) < DefaultCapacity {
					t.Errorf("expected capacity >= %d, got %d", DefaultCapacity, cap(h.elements))
				}
			},
		},
		{
			name:  "min_with_capacity",
			build: func() *Heap[int] { return Min[int](WithCapacity[int](64)) },
			check: func(t *testing.T, h *Heap[int]) {
				if cap(h.elements) < 64 {
					t.Errorf("expected capacity >= 64, got %d", cap(h.elements))
				}
				h.Add(10)
				h.Add(5)
				if h.Peek() != 5 {
					t.Errorf("expected min 5, got %d", h.Peek())
				}
			},
		},
		{
			name:  "max_with_capacity",
			build: func() *Heap[int] { return Max[int](WithCapacity[int](64)) },
			check: func(t *testing.T, h *Heap[int]) {
				if cap(h.elements) < 64 {
					t.Errorf("expected capacity >= 64, got %d", cap(h.elements))
				}
				h.Add(5)
				h.Add(10)
				if h.Peek() != 10 {
					t.Errorf("expected max 10, got %d", h.Peek())
				}
			},
		},
		{
			name: "new_ordered_option_precedence_and_no_slice_aliasing",
			build: func() *Heap[int] {
				opts := make([]Option[int], 0, 4)
				opts = append(opts, WithComparator(comparator.Reverse(comparator.NaturalOrder[int]())))
				sentinelOpt := WithCapacity[int](123)
				opts = append(opts, sentinelOpt)[:1]
				h := NewOrdered[int](opts...)
				// Verify caller's spare slice slot was not overwritten by NewOrdered
				full := opts[:2]
				cfg := &config[int]{}
				full[1](cfg)
				if cfg.capacity != 123 {
					panic("caller slice backing array was mutated")
				}
				return h
			},
			check: func(t *testing.T, h *Heap[int]) {
				h.AddAll(slices.Values([]int{1, 5, 3}))
				if got := h.Remove(); got != 5 {
					t.Errorf("expected custom reverse comparator to win, got %d", got)
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
			h := tc.build()
			tc.check(t, h)
		})
	}
}

func TestHeap_AddRemoveAndSort(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		heap  func() *Heap[int]
		check func(*testing.T, *Heap[int])
	}{
		{
			name: "min_heap",
			heap: func() *Heap[int] { return Min[int]() },
			check: func(t *testing.T, h *Heap[int]) {
				h.AddAll(slices.Values([]int{5, 1, 3, 2, 4}))
				if h.Size() != 5 || h.Peek() != 1 {
					t.Errorf("unexpected size/peek")
				}
				for i, exp := range []int{1, 2, 3, 4, 5} {
					if val := h.Remove(); val != exp {
						t.Errorf("at index %d: expected %d, got %d", i, exp, val)
					}
				}
			},
		},
		{
			name: "max_heap",
			heap: func() *Heap[int] { return Max[int]() },
			check: func(t *testing.T, h *Heap[int]) {
				h.AddAll(slices.Values([]int{5, 1, 3, 2, 4}))
				for i, exp := range []int{5, 4, 3, 2, 1} {
					if val := h.Remove(); val != exp {
						t.Errorf("at index %d: expected %d, got %d", i, exp, val)
					}
				}
			},
		},
		{
			name: "sort_1000_reverse_inputs",
			heap: func() *Heap[int] { return Min[int]() },
			check: func(t *testing.T, h *Heap[int]) {
				for i := 1000; i >= 1; i-- {
					h.Add(i)
				}
				for i := 1; i <= 1000; i++ {
					if val := h.Remove(); val != i {
						t.Fatalf("expected %d, got %d", i, val)
					}
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t, tc.heap())
		})
	}
}

func TestHeap_ClearIteratorsPanicsAndString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		wantPanic bool
		run       func(*testing.T)
	}{
		{
			name:      "peek_empty_panics",
			wantPanic: true,
			run:       func(_ *testing.T) { _ = Min[int]().Peek() },
		},
		{
			name:      "remove_empty_panics",
			wantPanic: true,
			run:       func(_ *testing.T) { _ = Min[int]().Remove() },
		},
		{
			name: "clear_and_remove_zero_pointer_slots",
			run: func(t *testing.T) {
				a, b, c := 10, 20, 30
				h := New[*int](WithComparator(func(x, y *int) int { return *x - *y }))
				h.Add(&c)
				h.Add(&a)
				h.Add(&b)
				backing := h.elements[:3]
				_ = h.Remove()
				if backing[2] != nil {
					t.Errorf("Remove did not zero vacated pointer slot")
				}
				h.Clear()
				if !h.Empty() || backing[0] != nil || backing[1] != nil {
					t.Errorf("Clear did not zero pointer slots")
				}
			},
		},
		{
			name: "all_full_and_early_break",
			run: func(t *testing.T) {
				h := Min[int]()
				h.AddAll(slices.Values([]int{4, 1, 3, 2}))
				if len(slices.Collect(h.All())) != 4 {
					t.Errorf("expected 4 elements")
				}
				count := 0
				for range h.All() {
					count++
					if count == 2 {
						break
					}
				}
				if count != 2 {
					t.Errorf("expected early break at 2, got %d", count)
				}
			},
		},
		{
			name: "string_formatting",
			run: func(t *testing.T) {
				h := Min[int]()
				if h.String() != "[]" {
					t.Errorf("got %q, want []", h.String())
				}
				h.Add(42)
				if h.String() != "[42]" {
					t.Errorf("got %q, want [42]", h.String())
				}
				h.Clear()
				h.AddAll(slices.Values([]int{3, 1, 2}))
				if h.String() != "[1 3 2]" {
					t.Errorf("got %q, want [1 3 2]", h.String())
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
			}
			tc.run(t)
		})
	}
}
