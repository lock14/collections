package arraydeque

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	t.Parallel()
	d := New[int]()
	if size := d.Size(); size != 0 {
		t.Errorf("new deque has non-zero size: %d", size)
	}
	if !d.isEmpty() {
		t.Error("new deque reports not empty")
	}
	if str := d.String(); str != "[]" {
		t.Errorf("new deque has wrong String(): %s", str)
	}
}

func TestArrayDeque_Add(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		items []int
		want  string
	}{
		{
			name:  "add_none",
			items: []int{},
			want:  "[]",
		},
		{
			name:  "add_one",
			items: []int{1},
			want:  "[1]",
		},
		{
			name:  "add_up_to_default_capacity",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:  "[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]",
		},
		{
			name:  "add_double_capacity",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:  "[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20]",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			for _, item := range tc.items {
				d.Add(item)
			}
			got := d.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}

func TestArrayDeque_Push(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		items []int
		want  string
	}{
		{
			name:  "push_none",
			items: []int{},
			want:  "[]",
		},
		{
			name:  "push_one",
			items: []int{1},
			want:  "[1]",
		},
		{
			name:  "push_up_to_default_capacity",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:  "[10, 9, 8, 7, 6, 5, 4, 3, 2, 1]",
		},
		{
			name:  "push_double_capacity",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:  "[20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1]",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			for _, item := range tc.items {
				d.Push(item)
			}
			got := d.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}

func TestArrayDeque_AddFront(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		items []int
		want  string
	}{
		{
			name:  "add_front_none",
			items: []int{},
			want:  "[]",
		},
		{
			name:  "add_front_one",
			items: []int{1},
			want:  "[1]",
		},
		{
			name:  "add_front_up_to_default_capacity",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:  "[10, 9, 8, 7, 6, 5, 4, 3, 2, 1]",
		},
		{
			name:  "add_front_double_capacity",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:  "[20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1]",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			for _, item := range tc.items {
				d.AddFront(item)
			}
			got := d.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}

func TestArrayDeque_AddBack(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		items []int
		want  string
	}{
		{
			name:  "add_back_none",
			items: []int{},
			want:  "[]",
		},
		{
			name:  "add_back_one",
			items: []int{1},
			want:  "[1]",
		},
		{
			name:  "add_back_up_to_default_capacity",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:  "[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]",
		},
		{
			name:  "add_back_double_capacity",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:  "[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20]",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			for _, item := range tc.items {
				d.Add(item)
			}
			got := d.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}

func TestArrayDeque_Rotate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		items []int
		want  string
	}{
		{
			name:  "rotate_none",
			items: []int{},
			want:  "[]",
		},
		{
			name:  "rotate_one",
			items: []int{1},
			want:  "[1]",
		},
		{
			name:  "rotate_even",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:  "[6, 7, 8, 9, 10, 1, 2, 3, 4, 5]",
		},
		{
			name:  "rotate_odd",
			items: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			want:  "[5, 6, 7, 8, 9, 1, 2, 3, 4]",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New[int]()
			for _, item := range tc.items {
				d.Add(item)
			}
			for i := 0; i < d.Size()/2; i++ {
				d.Add(d.Remove())
			}
			got := d.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}
