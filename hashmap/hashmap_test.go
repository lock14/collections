package hashmap

import (
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		opts      []Option
		wantPanic bool
		check     func(*testing.T, *HashMap[int, string])
	}{
		{
			name: "default",
			check: func(t *testing.T, m *HashMap[int, string]) {
				if m.Size() != 0 || !m.Empty() {
					t.Errorf("expected empty map, got size %d", m.Size())
				}
			},
		},
		{
			name: "zero_capacity",
			opts: []Option{WithCapacity(0)},
			check: func(t *testing.T, m *HashMap[int, string]) {
				if !m.Empty() {
					t.Errorf("expected empty map")
				}
			},
		},
		{
			name: "positive_capacity",
			opts: []Option{WithCapacity(100)},
			check: func(t *testing.T, m *HashMap[int, string]) {
				m.Put(1, "a")
				if m.Size() != 1 {
					t.Errorf("expected size 1, got %d", m.Size())
				}
			},
		},
		{
			name:      "negative_capacity_panics",
			opts:      []Option{WithCapacity(-1)},
			wantPanic: true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				r := recover()
				if tc.wantPanic && r == nil {
					t.Errorf("expected panic, got none")
				}
				if !tc.wantPanic && r != nil {
					t.Errorf("unexpected panic: %v", r)
				}
			}()
			m := New[int, string](tc.opts...)
			if tc.check != nil {
				tc.check(t, m)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T)
	}{
		{
			name: "non_nil_shares_underlying_map",
			check: func(t *testing.T) {
				raw := map[int]string{1: "a", 2: "b"}
				m := Wrap(raw)
				if m.Size() != 2 {
					t.Fatalf("expected size 2, got %d", m.Size())
				}
				if v, ok := m.Get(1); !ok || v != "a" {
					t.Errorf("expected (a, true), got (%q, %v)", v, ok)
				}
				m.Put(3, "c")
				if raw[3] != "c" {
					t.Errorf("expected raw map mutation for key 3, got %q", raw[3])
				}
				raw[4] = "d"
				if v, ok := m.Get(4); !ok || v != "d" {
					t.Errorf("expected wrapped map to observe raw key 4, got (%q, %v)", v, ok)
				}
			},
		},
		{
			name: "nil_map_allocates_writable_map",
			check: func(t *testing.T) {
				var raw map[int]string
				m := Wrap(raw)
				if !m.Empty() || m.Size() != 0 {
					t.Fatalf("expected empty wrapped map, got size %d", m.Size())
				}
				m.Put(1, "a")
				if v, ok := m.Get(1); !ok || v != "a" {
					t.Errorf("expected (a, true) after Put on Wrap(nil), got (%q, %v)", v, ok)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t)
		})
	}
}

func TestHashMap_Operations(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *HashMap[int, string])
	}{
		{
			name: "put_get_contains_overwrite",
			check: func(t *testing.T, m *HashMap[int, string]) {
				if v, ok := m.Get(99); ok || v != "" {
					t.Errorf("expected empty lookup ('', false), got (%q, %v)", v, ok)
				}
				m.Put(1, "a")
				m.Put(2, "b")
				m.Put(1, "c")

				if m.Size() != 2 || m.Empty() {
					t.Errorf("expected non-empty size 2, got size %d", m.Size())
				}
				if v, ok := m.Get(1); !ok || v != "c" {
					t.Errorf("expected ('c', true), got (%q, %v)", v, ok)
				}
				if m.ContainsKey(3) {
					t.Errorf("did not expect key 3")
				}
				if !m.ContainsKey(2) {
					t.Errorf("expected key 2")
				}
			},
		},
		{
			name: "remove_present_and_absent",
			check: func(t *testing.T, m *HashMap[int, string]) {
				m.Put(1, "a")
				m.Remove(1)
				m.Remove(2)

				if !m.Empty() || m.Size() != 0 {
					t.Errorf("expected empty map after remove, got size %d", m.Size())
				}
			},
		},
		{
			name: "clear",
			check: func(t *testing.T, m *HashMap[int, string]) {
				m.Put(1, "a")
				m.Put(2, "b")
				m.Clear()

				if m.Size() != 0 || !m.Empty() {
					t.Errorf("expected empty map after Clear, got size %d", m.Size())
				}
			},
		},
		{
			name: "keys_and_values_iterators_with_early_break",
			check: func(t *testing.T, m *HashMap[int, string]) {
				m.Put(1, "a")
				m.Put(2, "b")
				m.Put(3, "c")

				keys := slices.Collect(m.Keys())
				slices.Sort(keys)
				if !slices.Equal(keys, []int{1, 2, 3}) {
					t.Errorf("unexpected keys: %v", keys)
				}

				vals := slices.Collect(m.Values())
				slices.Sort(vals)
				if !slices.Equal(vals, []string{"a", "b", "c"}) {
					t.Errorf("unexpected values: %v", vals)
				}

				keyCount := 0
				for range m.Keys() {
					keyCount++
					break
				}
				if keyCount != 1 {
					t.Errorf("expected early break keyCount 1, got %d", keyCount)
				}

				valCount := 0
				for range m.Values() {
					valCount++
					break
				}
				if valCount != 1 {
					t.Errorf("expected early break valCount 1, got %d", valCount)
				}
			},
		},
		{
			name: "all_iterator_and_early_break",
			check: func(t *testing.T, m *HashMap[int, string]) {
				m.Put(1, "a")
				m.Put(2, "b")

				var keys []int
				var vals []string
				for k, v := range m.All() {
					keys = append(keys, k)
					vals = append(vals, v)
				}
				slices.Sort(keys)
				slices.Sort(vals)
				if !slices.Equal(keys, []int{1, 2}) {
					t.Errorf("unexpected keys: %v", keys)
				}
				if !slices.Equal(vals, []string{"a", "b"}) {
					t.Errorf("unexpected values: %v", vals)
				}

				count := 0
				for range m.All() {
					count++
					break
				}
				if count != 1 {
					t.Errorf("expected early break count 1, got %d", count)
				}
			},
		},
		{
			name: "string_empty_single_and_multiple",
			check: func(t *testing.T, m *HashMap[int, string]) {
				if got := m.String(); got != "map[]" {
					t.Errorf("expected map[], got %s", got)
				}
				m.Put(1, "a")
				if got := m.String(); got != "map[1:a]" {
					t.Errorf("expected map[1:a], got %s", got)
				}
				m.Put(2, "b")
				got := m.String()
				if got != "map[1:a 2:b]" && got != "map[2:b 1:a]" {
					t.Errorf("expected map[1:a 2:b] or map[2:b 1:a], got %s", got)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := New[int, string]()
			tc.check(t, m)
		})
	}
}
