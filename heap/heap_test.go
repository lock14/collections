package heap

import (
	"github.com/lock14/collections"
	"github.com/lock14/collections/comparator"
	"slices"
	"testing"
)

func TestHeapImplementsQueue(t *testing.T) {
	queue[int](NewOrdered[int]())
}

func queue[T any](q collections.Queue[T]) collections.Queue[T] {
	return q
}

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		opts        []Option[int]
		shouldPanic bool
		check       func(*testing.T, *Heap[int])
	}{
		{
			name:        "nil_comparator_panics",
			shouldPanic: true,
		},
		{
			name: "default_capacity",
			opts: []Option[int]{WithComparator(comparator.NaturalOrder[int]())},
			check: func(t *testing.T, h *Heap[int]) {
				if h.Size() != 0 {
					t.Errorf("expected size 0")
				}
				if cap(h.elements) < DefaultCapacity {
					t.Errorf("expected capacity >= %d, got %d", DefaultCapacity, cap(h.elements))
				}
			},
		},
		{
			name: "custom_capacity",
			opts: []Option[int]{WithComparator(comparator.NaturalOrder[int]()), WithCapacity[int](100)},
			check: func(t *testing.T, h *Heap[int]) {
				if cap(h.elements) < 100 {
					t.Errorf("expected capacity >= 100, got %d", cap(h.elements))
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic, but did not panic")
					}
				}()
				_ = New[int](tc.opts...)
				return
			}
			h := New[int](tc.opts...)
			tc.check(t, h)
		})
	}
}

func TestNewOrdered(t *testing.T) {
	t.Parallel()
	h := NewOrdered[int](WithCapacity[int](50))
	if cap(h.elements) < 50 {
		t.Errorf("expected capacity >= 50, got %d", cap(h.elements))
	}
	h.Add(10)
	h.Add(5)
	if h.Peek() != 5 {
		t.Errorf("expected min 5, got %d", h.Peek())
	}
}

func TestHeap_AddRemove(t *testing.T) {
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
				if h.Size() != 5 {
					t.Errorf("expected size 5")
				}
				if h.Peek() != 1 {
					t.Errorf("expected peek 1")
				}
				expected := []int{1, 2, 3, 4, 5}
				for i, exp := range expected {
					val := h.Remove()
					if val != exp {
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
				expected := []int{5, 4, 3, 2, 1}
				for i, exp := range expected {
					val := h.Remove()
					if val != exp {
						t.Errorf("at index %d: expected %d, got %d", i, exp, val)
					}
				}
			},
		},
		{
			name: "sort_large",
			heap: func() *Heap[int] { return Min[int]() },
			check: func(t *testing.T, h *Heap[int]) {
				for i := 1000; i >= 1; i-- {
					h.Add(i)
				}
				for i := 1; i <= 1000; i++ {
					val := h.Remove()
					if val != i {
						t.Fatalf("expected %d, got %d", i, val)
					}
				}
			},
		},
		{
			name: "remove_middle",
			heap: func() *Heap[int] { return Min[int]() },
			check: func(t *testing.T, h *Heap[int]) {
				h.AddAll(slices.Values([]int{10, 20, 30, 40, 50, 60}))
				h.Remove()
				if h.Peek() != 20 {
					t.Errorf("expected 20, got %d", h.Peek())
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := tc.heap()
			tc.check(t, h)
		})
	}
}

func TestHeap_Clear(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		heap  func() *Heap[int]
		check func(*testing.T, *Heap[int])
	}{
		{
			name: "clear_elements",
			heap: func() *Heap[int] { return Min[int]() },
			check: func(t *testing.T, h *Heap[int]) {
				h.Add(1)
				h.Add(2)
				h.Clear()
				if !h.Empty() || h.Size() != 0 {
					t.Errorf("expected empty after clear")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := tc.heap()
			tc.check(t, h)
		})
	}
}

func TestHeap_Iterators(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		heap  func() *Heap[int]
		check func(*testing.T, *Heap[int])
	}{
		{
			name: "all",
			heap: func() *Heap[int] { return Min[int]() },
			check: func(t *testing.T, h *Heap[int]) {
				h.AddAll(slices.Values([]int{1, 2, 3}))
				count := 0
				for range h.All() {
					count++
				}
				if count != 3 {
					t.Errorf("expected 3 items, got %d", count)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := tc.heap()
			tc.check(t, h)
		})
	}
}

func TestHeap_Panics(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		heap  func() *Heap[int]
		check func(*testing.T, *Heap[int])
	}{
		{
			name: "peek_empty",
			heap: func() *Heap[int] { return Min[int]() },
			check: func(t *testing.T, h *Heap[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				h.Peek()
			},
		},
		{
			name: "remove_empty",
			heap: func() *Heap[int] { return Min[int]() },
			check: func(t *testing.T, h *Heap[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				h.Remove()
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := tc.heap()
			tc.check(t, h)
		})
	}
}

func TestHeap_Coverage(t *testing.T) {
	h := NewOrdered[int]()
	assertPanics := func(f func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		f()
	}
	assertPanics(func() { h.Remove() })
	_ = comparator.NaturalOrder[int]()(1, 2)
}

func TestHeap_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		elements []int
		want     string
	}{
		{
			name:     "empty",
			elements: nil,
			want:     "[]",
		},
		{
			name:     "single",
			elements: []int{42},
			want:     "[42]",
		},
		{
			name:     "multiple",
			elements: []int{3, 1, 2},
			want:     "[1 3 2]",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := Min[int]()
			for _, v := range tc.elements {
				h.Add(v)
			}
			if got := h.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}
