package linkedlist

import (
	"slices"
	"testing"
)

func TestLinkedList_AddRemoveFrontBack(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *LinkedList[int])
	}{
		{
			name: "add_front",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddFront(1)
				l.AddFront(2)
				if l.Size() != 2 || l.PeekFront() != 2 || l.PeekBack() != 1 {
					t.Errorf("unexpected state after AddFront")
				}
			},
		},
		{
			name: "add_back",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddBack(1)
				l.AddBack(2)
				if l.PeekFront() != 1 || l.PeekBack() != 2 {
					t.Errorf("unexpected state after AddBack")
				}
			},
		},
		{
			name: "remove_front",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddBack(1)
				l.AddBack(2)
				if val := l.RemoveFront(); val != 1 || l.Size() != 1 {
					t.Errorf("expected 1 and size 1, got %d", val)
				}
			},
		},
		{
			name: "remove_back",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddBack(1)
				l.AddBack(2)
				if val := l.RemoveBack(); val != 2 || l.Size() != 1 {
					t.Errorf("expected 2 and size 1, got %d", val)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := New[int]()
			tc.check(t, l)
		})
	}
}

func TestLinkedList_GetSetAndBounds(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		wantPanic bool
		check     func(*testing.T, *LinkedList[int])
	}{
		{
			name: "forward_and_backward_half_traversal",
			check: func(t *testing.T, l *LinkedList[int]) {
				// 8 elements ensures idx=2 exercises cur=cur.next loop (idx < size/2)
				// and idx=5 exercises cur=cur.prev loop (idx >= size/2).
				l.AddAll(slices.Values([]int{10, 20, 30, 40, 50, 60, 70, 80}))
				for i := 0; i < 8; i++ {
					if got := l.Get(i); got != (i+1)*10 {
						t.Errorf("Get(%d) = %d, want %d", i, got, (i+1)*10)
					}
				}
				l.Set(2, 300)
				l.Set(5, 600)
				if l.Get(2) != 300 || l.Get(5) != 600 {
					t.Errorf("Set failed: got %d, %d", l.Get(2), l.Get(5))
				}
			},
		},
		{
			name:      "get_negative_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { _ = l.Get(-1) },
		},
		{
			name:      "get_out_of_bounds_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { _ = l.Get(0) },
		},
		{
			name:      "set_out_of_bounds_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { l.Set(0, 10) },
		},
		{
			name:      "remove_front_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { _ = l.RemoveFront() },
		},
		{
			name:      "remove_back_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { _ = l.RemoveBack() },
		},
		{
			name:      "remove_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { _ = l.Remove() },
		},
		{
			name:      "pop_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { _ = l.Pop() },
		},
		{
			name:      "peek_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { _ = l.Peek() },
		},
		{
			name:      "peek_front_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { _ = l.PeekFront() },
		},
		{
			name:      "peek_back_empty_panics",
			wantPanic: true,
			check:     func(_ *testing.T, l *LinkedList[int]) { _ = l.PeekBack() },
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

func TestLinkedList_InterfacesIteratorsAndZeroValue(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T)
	}{
		{
			name: "queue_and_stack_interfaces",
			check: func(t *testing.T) {
				l := New[int]()
				l.Add(1)
				l.Add(2)
				if l.Peek() != 1 || l.Remove() != 1 {
					t.Errorf("queue interface mismatch")
				}
				l.Clear()
				l.Push(10)
				l.Push(20)
				if l.Pop() != 20 || !slices.Equal(slices.Collect(l.All()), []int{10}) {
					t.Errorf("stack interface mismatch")
				}
			},
		},
		{
			name: "string_empty_and_non_empty",
			check: func(t *testing.T) {
				l := New[int]()
				if l.String() != "[]" {
					t.Errorf("got %q, want []", l.String())
				}
				l.AddAll(slices.Values([]int{1, 2, 3}))
				if l.String() != "[1 2 3]" {
					t.Errorf("got %q, want [1 2 3]", l.String())
				}
			},
		},
		{
			name: "all_and_backward_with_early_break",
			check: func(t *testing.T) {
				l := New[int]()
				l.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
				if got := slices.Collect(l.Backward()); !slices.Equal(got, []int{5, 4, 3, 2, 1}) {
					t.Errorf("Backward() = %v", got)
				}
				var fwd []int
				for v := range l.All() {
					fwd = append(fwd, v)
					if v == 2 {
						break
					}
				}
				if !slices.Equal(fwd, []int{1, 2}) {
					t.Errorf("All break = %v", fwd)
				}
				var bwd []int
				for v := range l.Backward() {
					bwd = append(bwd, v)
					if v == 4 {
						break
					}
				}
				if !slices.Equal(bwd, []int{5, 4}) {
					t.Errorf("Backward break = %v", bwd)
				}
			},
		},
		{
			name: "zero_value_operations",
			check: func(t *testing.T) {
				var z LinkedList[int]
				if !z.Empty() || z.Size() != 0 || z.String() != "[]" {
					t.Errorf("unexpected zero-value state")
				}
				if len(slices.Collect(z.All())) != 0 || len(slices.Collect(z.Backward())) != 0 {
					t.Errorf("expected empty iterators on zero value")
				}
				z.AddFront(2)
				z.AddFront(1)
				z.AddBack(3)
				if !slices.Equal(slices.Collect(z.All()), []int{1, 2, 3}) {
					t.Errorf("zero value after adds failed")
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
