package heap

import (
	"github.com/lock14/collections"
	"github.com/lock14/collections/comparator"
	"slices"
	"testing"
)

func TestHeapImplementsQueue(t *testing.T) {
	queue[int](New[int]())
}

func queue[T any](q collections.Queue[T]) collections.Queue[T] {
	return q
}

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Option[int]
		check func(*testing.T, *Heap[int])
	}{
		{
			name: "default_capacity",
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
			opts: []Option[int]{WithCapacity[int](100)},
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
			h := New[int](tc.opts...)
			tc.check(t, h)
		})
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
	h := New[int]()
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
