package hashset

import (
	"slices"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefaultConstruction(t *testing.T) {
	t.Parallel()
	s := New[int]()
	if size := s.Size(); size != 0 {
		t.Errorf("new hashset has non-zero size: %d", size)
	}
	if !s.Empty() {
		t.Error("new hashset reports not empty")
	}
	if str := s.String(); str != "[]" {
		t.Errorf("new hashset has wrong String(): %s", str)
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		items []int
		want  []int
	}{
		{
			name:  "add_none",
			items: []int{},
			want:  nil,
		},
		{
			name:  "add_one",
			items: []int{1},
			want:  []int{1},
		},
		{
			name:  "add_duplicates",
			items: []int{1, 2, 2, 1, 2, 1, 1},
			want:  []int{1, 2},
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
			got := slices.Collect(s.All())
			sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}
