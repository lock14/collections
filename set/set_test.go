package set

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefaultConstruction(t *testing.T) {
	t.Parallel()
	s := New[int]()
	if size := s.Size(); size != 0 {
		t.Errorf("new set has non-zero size: %d", size)
	}
	if !s.isEmpty() {
		t.Error("new set reports not empty")
	}
	if str := s.String(); str != "[]" {
		t.Errorf("new set has wrong String(): %s", str)
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
			name:  "add_duplicates",
			items: []int{1, 2, 2, 1, 2, 1, 1},
			want:  "[1, 2]",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := New[int]()
			for _, item := range tc.items {
				s.Add(item)
			}
			got := s.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}
