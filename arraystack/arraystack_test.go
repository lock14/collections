package arraystack

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefaultConstruction(t *testing.T) {
	t.Parallel()
	s := New[int]()
	if size := s.Size(); size != 0 {
		t.Errorf("new stack has non-zero size: %d", size)
	}
	if !s.isEmpty() {
		t.Error("new stack reports not empty")
	}
	if str := s.String(); str != "[]" {
		t.Errorf("new stack has wrong String(): %s", str)
	}
}

func TestPush(t *testing.T) {
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
			q := New[int]()
			for _, item := range tc.items {
				q.Push(item)
			}
			got := q.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}
