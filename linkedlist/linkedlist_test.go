package linked_list

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
				if l.Size() != 2 {
					t.Errorf("expected size 2")
				}
				if l.PeekFront() != 2 {
					t.Errorf("expected front to be 2")
				}
				if l.PeekBack() != 1 {
					t.Errorf("expected back to be 1")
				}
			},
		},
		{
			name: "add_back",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddBack(1)
				l.AddBack(2)
				if l.PeekFront() != 1 {
					t.Errorf("expected front to be 1")
				}
				if l.PeekBack() != 2 {
					t.Errorf("expected back to be 2")
				}
			},
		},
		{
			name: "remove_front",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddBack(1)
				l.AddBack(2)
				val := l.RemoveFront()
				if val != 1 {
					t.Errorf("expected 1")
				}
				if l.Size() != 1 {
					t.Errorf("expected size 1")
				}
			},
		},
		{
			name: "remove_back",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddBack(1)
				l.AddBack(2)
				val := l.RemoveBack()
				if val != 2 {
					t.Errorf("expected 2")
				}
				if l.Size() != 1 {
					t.Errorf("expected size 1")
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

func TestLinkedList_GetSet(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *LinkedList[int])
	}{
		{
			name: "get",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddBack(10)
				l.AddBack(20)
				l.AddBack(30)
				
				if val := l.Get(0); val != 10 {
					t.Errorf("expected 10, got %d", val)
				}
				if val := l.Get(1); val != 20 {
					t.Errorf("expected 20, got %d", val)
				}
				if val := l.Get(2); val != 30 {
					t.Errorf("expected 30, got %d", val)
				}
			},
		},
		{
			name: "set",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddBack(10)
				l.AddBack(20)
				
				l.Set(1, 100)
				if val := l.Get(1); val != 100 {
					t.Errorf("expected 100, got %d", val)
				}
			},
		},
		{
			name: "get_out_of_bounds",
			check: func(t *testing.T, l *LinkedList[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				l.Get(0)
			},
		},
		{
			name: "set_out_of_bounds",
			check: func(t *testing.T, l *LinkedList[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				l.Set(0, 10)
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

func TestLinkedList_Interfaces(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *LinkedList[int])
	}{
		{
			name: "queue_interface",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.Add(1)
				l.Add(2)
				if l.Peek() != 1 {
					t.Errorf("expected 1")
				}
				if l.Remove() != 1 {
					t.Errorf("expected 1")
				}
			},
		},
		{
			name: "stack_interface",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.Push(1)
				l.Push(2)
				if l.Pop() != 2 {
					t.Errorf("expected 2")
				}
			},
		},
		{
			name: "clear",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.Add(1)
				l.Add(2)
				l.Clear()
				if !l.Empty() || l.Size() != 0 {
					t.Errorf("expected empty")
				}
			},
		},
		{
			name: "iterators",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.AddAll(slices.Values([]int{1, 2, 3}))
				items := slices.Collect(l.All())
				if !slices.Equal(items, []int{1, 2, 3}) {
					t.Errorf("expected [1, 2, 3], got %v", items)
				}
			},
		},
		{
			name: "string",
			check: func(t *testing.T, l *LinkedList[int]) {
				l.Add(1)
				l.Add(2)
				if l.String() != "[1, 2]" {
					t.Errorf("expected [1, 2], got %s", l.String())
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
