package arraydeque

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefaultConstruction(t *testing.T) {
	t.Parallel()
	q := New[int]()
	if size := q.Size(); size != 0 {
		t.Errorf("new deque has non-zero size: %d", size)
	}
	if !q.isEmpty() {
		t.Error("new deque reports not empty")
	}
	if str := q.String(); str != "[]" {
		t.Errorf("new deque has wrong String(): %s", str)
	}
}

func TestAdd(t *testing.T) {
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
			q := New[int]()
			for _, item := range tc.items {
				q.Add(item)
			}
			got := q.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}

func TestRotate(t *testing.T) {
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
			q := New[int]()
			for _, item := range tc.items {
				q.Add(item)
			}
			for i := 0; i < q.Size()/2; i++ {
				q.Add(q.Remove())
			}
			got := q.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}
