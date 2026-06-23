package hashmap

import (
	"slices"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *HashMap[int, string])
	}{
		{
			name: "default",
			check: func(t *testing.T, hm *HashMap[int, string]) {
				if hm.Size() != 0 || !hm.Empty() {
					t.Errorf("expected empty map")
				}
			},
		},
		{
			name: "capacity",
			opts: []Option{WithCapacity(100)},
			check: func(t *testing.T, hm *HashMap[int, string]) {
				hm.Put(1, "a")
				if hm.Size() != 1 {
					t.Errorf("expected size 1")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hm := New[int, string](tc.opts...)
			tc.check(t, hm)
		})
	}
}

func TestWrap(t *testing.T) {
	t.Parallel()
	m := map[int]string{1: "a", 2: "b"}
	hm := Wrap(m)
	if hm.Size() != 2 {
		t.Errorf("expected size 2")
	}
	if v, _ := hm.Get(1); v != "a" {
		t.Errorf("expected 'a'")
	}
}

func TestHashMap_Operations(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *HashMap[int, string])
	}{
		{
			name: "put_get_contains",
			check: func(t *testing.T, hm *HashMap[int, string]) {
				hm.Put(1, "a")
				hm.Put(2, "b")
				hm.Put(1, "c") // overwrite
				
				if hm.Size() != 2 {
					t.Errorf("expected size 2")
				}
				if v, ok := hm.Get(1); !ok || v != "c" {
					t.Errorf("expected 'c', got %v", v)
				}
				if hm.ContainsKey(3) {
					t.Errorf("did not expect 3")
				}
				if !hm.ContainsKey(2) {
					t.Errorf("expected 2")
				}
			},
		},
		{
			name: "remove",
			check: func(t *testing.T, hm *HashMap[int, string]) {
				hm.Put(1, "a")
				hm.Remove(1)
				hm.Remove(2) // non-existent
				
				if !hm.Empty() {
					t.Errorf("expected empty")
				}
			},
		},
		{
			name: "clear",
			check: func(t *testing.T, hm *HashMap[int, string]) {
				hm.Put(1, "a")
				hm.Clear()
				
				if hm.Size() != 0 {
					t.Errorf("expected empty")
				}
			},
		},
		{
			name: "iterators",
			check: func(t *testing.T, hm *HashMap[int, string]) {
				hm.Put(1, "a")
				hm.Put(2, "b")
				
				keys := slices.Collect(hm.Keys())
				sort.Ints(keys)
				if !slices.Equal(keys, []int{1, 2}) {
					t.Errorf("keys wrong")
				}
				
				vals := slices.Collect(hm.Values())
				sort.Strings(vals)
				if !slices.Equal(vals, []string{"a", "b"}) {
					t.Errorf("values wrong")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hm := New[int, string]()
			tc.check(t, hm)
		})
	}
}
