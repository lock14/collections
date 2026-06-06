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
			opts: []Option{func(c *Config) { c.Capacity = 0 }},
			check: func(t *testing.T, d *ArrayDeque[int]) {
				d.Add(1) // Should not panic
				if size := d.Size(); size != 1 {
					t.Errorf("expected size 1, got: %d", size)
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

func TestType(t *testing.T) {
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

func testType[T any](deque collections.Deque[T]) {}
