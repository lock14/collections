package arraydeque

import (
	"github.com/lock14/collections"
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "default",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				if size := d.Size(); size != 0 {
					t.Errorf("new deque has non-zero size: %d", size)
				}
				if !d.Empty() {
					t.Error("new deque reports not empty")
				}
				if str := d.String(); str != "[]" {
					t.Errorf("new deque has wrong String(): %s", str)
				}
			},
		},
		{
			name: "zero_capacity",
			opts: []Option{WithCapacity(0)},
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.Add(1) // Should not panic
				if size := d.Size(); size != 1 {
					t.Errorf("expected size 1, got: %d", size)
				}
			},
		},
		{
			name: "custom_capacity",
			opts: []Option{WithCapacity(10)},
			check: func(t *testing.T, d *ArrayDeque[int]) {
				if cap(d.slice) != 10 {
					t.Errorf("expected capacity 10, got: %d", cap(d.slice))
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int](tc.opts...)
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_Clear(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "clear_and_add",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.Add(1)
				d.Clear()
				d.Add(2) // Should not panic
				if size := d.Size(); size != 1 {
					t.Errorf("expected size 1 after clear and add, got %d", size)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_Add(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "add_none",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				if got := slices.Collect(d.All()); !slices.Equal(got, nil) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, nil)
				}
			},
		},
		{
			name: "add_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.Add(1)
				want := []int{1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "add_up_to_default_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				for _, item := range want {
					d.Add(item)
				}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "add_double_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
				for _, item := range want {
					d.Add(item)
				}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "add_fifty_percent_threshold",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				want := make([]int, 600)
				for i := 0; i < 600; i++ {
					want[i] = i
					d.Add(i)
				}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value")
				}
			},
		},
		{
			name: "add_twenty_five_percent_threshold",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				want := make([]int, 2500)
				for i := 0; i < 2500; i++ {
					want[i] = i
					d.Add(i)
				}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_Push(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "push_none",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				if got := slices.Collect(d.All()); !slices.Equal(got, nil) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, nil)
				}
			},
		},
		{
			name: "push_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.Push(1)
				want := []int{1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "push_up_to_default_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				for _, item := range items {
					d.Push(item)
				}
				want := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "push_double_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
				for _, item := range items {
					d.Push(item)
				}
				want := []int{20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_AddFront(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "add_front_none",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				if got := slices.Collect(d.All()); !slices.Equal(got, nil) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, nil)
				}
			},
		},
		{
			name: "add_front_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddFront(1)
				want := []int{1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "add_front_up_to_default_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				for _, item := range items {
					d.AddFront(item)
				}
				want := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "add_front_double_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
				for _, item := range items {
					d.AddFront(item)
				}
				want := []int{20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_AddBack(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "add_back_none",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				if got := slices.Collect(d.All()); !slices.Equal(got, nil) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, nil)
				}
			},
		},
		{
			name: "add_back_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.Add(1)
				want := []int{1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "add_back_up_to_default_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				for _, item := range want {
					d.Add(item)
				}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "add_back_double_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
				for _, item := range want {
					d.Add(item)
				}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_Rotate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "rotate_none",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				if got := slices.Collect(d.All()); !slices.Equal(got, nil) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, nil)
				}
			},
		},
		{
			name: "rotate_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.Add(1)
				for i := 0; i < d.Size()/2; i++ {
					d.Add(d.Remove())
				}
				want := []int{1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "rotate_even",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				for _, item := range items {
					d.Add(item)
				}
				for i := 0; i < d.Size()/2; i++ {
					d.Add(d.Remove())
				}
				want := []int{6, 7, 8, 9, 10, 1, 2, 3, 4, 5}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "rotate_odd",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
				for _, item := range items {
					d.Add(item)
				}
				for i := 0; i < d.Size()/2; i++ {
					d.Add(d.Remove())
				}
				want := []int{5, 6, 7, 8, 9, 1, 2, 3, 4}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestType(_ *testing.T) {
	l := New[int]()
	testType[int](l)
}

func TestArrayDeque_AddAll(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "add_none",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddAll(slices.Values([]int{}))
				if got := slices.Collect(d.All()); !slices.Equal(got, nil) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, nil)
				}
			},
		},
		{
			name: "add_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddAll(slices.Values([]int{1}))
				want := []int{1}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "add_up_to_default_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				d.AddAll(slices.Values(want))
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "add_double_capacity",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
				d.AddAll(slices.Values(want))
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func testType[T any](_ collections.Deque[T]) {}

func TestArrayDeque_PeekFront(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "peek_empty_panics",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				d.PeekFront()
			},
		},
		{
			name: "peek_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddBack(1)
				if got := d.PeekFront(); got != 1 {
					t.Errorf("expected 1, got %d", got)
				}
				if size := d.Size(); size != 1 {
					t.Errorf("PeekFront should not remove element")
				}
			},
		},
		{
			name: "peek_after_adds",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddBack(1)
				d.AddBack(2)
				d.AddFront(0)
				if got := d.PeekFront(); got != 0 {
					t.Errorf("expected 0, got %d", got)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_PeekBack(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "peek_empty_panics",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				d.PeekBack()
			},
		},
		{
			name: "peek_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddFront(1)
				if got := d.PeekBack(); got != 1 {
					t.Errorf("expected 1, got %d", got)
				}
				if size := d.Size(); size != 1 {
					t.Errorf("PeekBack should not remove element")
				}
			},
		},
		{
			name: "peek_after_adds",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddFront(1)
				d.AddFront(0)
				d.AddBack(2)
				if got := d.PeekBack(); got != 2 {
					t.Errorf("expected 2, got %d", got)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_RemoveFront(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "remove_empty_panics",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				d.RemoveFront()
			},
		},
		{
			name: "remove_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddBack(1)
				if got := d.RemoveFront(); got != 1 {
					t.Errorf("expected 1, got %d", got)
				}
				if size := d.Size(); size != 0 {
					t.Errorf("expected size 0, got %d", size)
				}
			},
		},
		{
			name: "remove_multiple",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				items := []int{1, 2, 3}
				d.AddAll(slices.Values(items))
				for _, want := range items {
					if got := d.RemoveFront(); got != want {
						t.Errorf("expected %d, got %d", want, got)
					}
				}
				if !d.Empty() {
					t.Errorf("expected empty")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_RemoveBack(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "remove_empty_panics",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
				d.RemoveBack()
			},
		},
		{
			name: "remove_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddFront(1)
				if got := d.RemoveBack(); got != 1 {
					t.Errorf("expected 1, got %d", got)
				}
				if size := d.Size(); size != 0 {
					t.Errorf("expected size 0, got %d", size)
				}
			},
		},
		{
			name: "remove_multiple",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				items := []int{1, 2, 3} // front to back
				d.AddAll(slices.Values(items))
				// remove back should be 3, 2, 1
				wantReverse := []int{3, 2, 1}
				for _, want := range wantReverse {
					if got := d.RemoveBack(); got != want {
						t.Errorf("expected %d, got %d", want, got)
					}
				}
				if !d.Empty() {
					t.Errorf("expected empty")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_Aliases(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "Peek_is_PeekFront",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddBack(1)
				d.AddBack(2)
				if d.Peek() != d.PeekFront() {
					t.Errorf("Peek != PeekFront")
				}
			},
		},
		{
			name: "Remove_is_RemoveFront",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddBack(1)
				if d.Remove() != 1 {
					t.Errorf("Remove did not act as RemoveFront")
				}
			},
		},
		{
			name: "Pop_is_RemoveFront",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddBack(1)
				if d.Pop() != 1 {
					t.Errorf("Pop did not act as RemoveFront")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_Backward(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "backward_none",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				if got := slices.Collect(d.Backward()); !slices.Equal(got, nil) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, nil)
				}
			},
		},
		{
			name: "backward_one",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddBack(1)
				want := []int{1}
				if got := slices.Collect(d.Backward()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "backward_multiple",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				items := []int{1, 2, 3, 4, 5}
				d.AddAll(slices.Values(items))
				want := []int{5, 4, 3, 2, 1}
				if got := slices.Collect(d.Backward()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
		{
			name: "backward_wrapped",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				// Capacity is small, wrap around
				d.AddAll(slices.Values([]int{1, 2, 3, 4})) // Cap becomes 4
				d.RemoveFront()                            // size 3, front 1
				d.RemoveFront()                            // size 2, front 2
				d.AddBack(5)                               // size 3, back 1 (wrapped)
				d.AddBack(6)                               // size 4, back 2 (wrapped)
				// d is [3, 4, 5, 6]
				want := []int{6, 5, 4, 3}
				if got := slices.Collect(d.Backward()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_All_Wrapped(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "all_wrapped",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddAll(slices.Values([]int{1, 2, 3, 4})) // Cap becomes 4
				d.RemoveFront()                            // size 3, front 1
				d.RemoveFront()                            // size 2, front 2
				d.AddBack(5)                               // size 3, back 1 (wrapped)
				d.AddBack(6)                               // size 4, back 2 (wrapped)
				want := []int{3, 4, 5, 6}
				if got := slices.Collect(d.All()); !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_All_Break(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "all_break",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddAll(slices.Values([]int{1, 2, 3, 4}))
				got := []int{}
				for v := range d.All() {
					got = append(got, v)
					if v == 2 {
						break
					}
				}
				want := []int{1, 2}
				if !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_Backward_Break(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *ArrayDeque[int])
	}{
		{
			name: "backward_break",
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.AddAll(slices.Values([]int{1, 2, 3, 4}))
				got := []int{}
				for v := range d.Backward() {
					got = append(got, v)
					if v == 3 {
						break
					}
				}
				want := []int{4, 3}
				if !slices.Equal(got, want) {
					t.Errorf("wrong slice value, got: %v, want: %v", got, want)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			tc.check(t, d)
		})
	}
}

func TestArrayDeque_CoveragePanics(t *testing.T) {
	ad := New[int]()
	assertPanics := func(f func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		f()
	}
	assertPanics(func() { ad.RemoveFront() })
	assertPanics(func() { ad.RemoveBack() })
	assertPanics(func() { ad.PeekFront() })
	assertPanics(func() { ad.PeekBack() })
	_ = ad.String()
}
