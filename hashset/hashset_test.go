package hashset

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefaultConstruction(t *testing.T) {
	t.Parallel()
	s := New[int]()
	if size := s.Size(); size != 0 {
		t.Errorf("new hashset has non-zero size: %d", size)
	}
	if !s.isEmpty() {
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
			parts := strings.Split(got[1:len(got)-1], ",")
			slices.SortFunc(parts, func(a, b string) int {
				return strings.Compare(b, a)
			})
			got = strings.Join([]string{"[", "]"}, strings.Join(parts, ","))
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("wrong string value, -got,+want: %s", diff)
			}
		})
	}
}
