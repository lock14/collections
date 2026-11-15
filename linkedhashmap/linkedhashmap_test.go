package linkedhashmap

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
		check     func(*testing.T, *LinkedHashMap[int, int])
	}{
		{
			name: "default",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				if m.Size() != 0 || !m.Empty() {
					t.Errorf("expected empty map, got size %d", m.Size())
				}
				if m.accessOrder != InsertionOrder {
					t.Errorf("expected insertion order by default")
				}
				m.Put(1, 10)
				m.Put(2, 20)
				if m.Size() != 2 {
					t.Errorf("expected size 2, got %d", m.Size())
				}
			},
		},
		{
			name: "access_order",
			opts: []Option{WithAccessOrder()},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				if m.accessOrder != AccessOrder {
					t.Errorf("expected access order")
				}
			},
		},
		{
			name: "insertion_order_override",
			opts: []Option{WithAccessOrder(), WithInsertionOrder()},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				if m.accessOrder != InsertionOrder {
					t.Errorf("expected insertion order")
				}
			},
		},
		{
			name: "max_elements",
			opts: []Option{WithMaxElements(2)},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				if m.maxElements != 2 {
					t.Errorf("expected maxElements 2, got %d", m.maxElements)
				}
			},
		},
		{
			name: "capacity_zero_and_positive",
			opts: []Option{WithCapacity(0), WithCapacity(100)},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
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
		{
			name:      "zero_max_elements_panics",
			opts:      []Option{WithMaxElements(0)},
			wantPanic: true,
		},
		{
			name:      "negative_max_elements_panics",
			opts:      []Option{WithMaxElements(-5)},
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
			m := New[int, int](tc.opts...)
			if tc.check != nil {
				tc.check(t, m)
			}
		})
	}
}

func TestLinkedHashMap_GetAndPut(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *LinkedHashMap[int, int])
	}{
		{
			name: "get_missing",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				if val, ok := m.Get(1); ok || val != 0 {
					t.Errorf("expected (0, false), got (%d, %v)", val, ok)
				}
			},
		},
		{
			name: "get_insertion_order_preserves_order",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				if val, ok := m.Get(1); !ok || val != 10 {
					t.Errorf("expected (10, true), got (%d, %v)", val, ok)
				}
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{1, 2}) {
					t.Errorf("expected [1 2], got %v", got)
				}
			},
		},
		{
			name: "get_access_order_moves_to_back",
			opts: []Option{WithAccessOrder()},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				if val, ok := m.Get(1); !ok || val != 10 {
					t.Errorf("expected (10, true), got (%d, %v)", val, ok)
				}
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{2, 1}) {
					t.Errorf("expected [2 1], got %v", got)
				}
			},
		},
		{
			name: "put_update_insertion_order",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Put(1, 100)
				if val, _ := m.Get(1); val != 100 {
					t.Errorf("expected 100, got %d", val)
				}
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{1, 2}) {
					t.Errorf("expected [1 2], got %v", got)
				}
			},
		},
		{
			name: "put_update_access_order",
			opts: []Option{WithAccessOrder()},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Put(1, 100)
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{2, 1}) {
					t.Errorf("expected [2 1], got %v", got)
				}
			},
		},
		{
			name: "put_evicts_eldest_max_2",
			opts: []Option{WithMaxElements(2)},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Put(3, 30)
				if m.Size() != 2 || m.ContainsKey(1) {
					t.Errorf("expected key 1 evicted and size 2, got %s", m.String())
				}
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{2, 3}) {
					t.Errorf("expected [2 3], got %v", got)
				}
			},
		},
		{
			name: "put_evicts_eldest_max_1",
			opts: []Option{WithMaxElements(1)},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				if m.Size() != 1 || m.ContainsKey(1) || !m.ContainsKey(2) {
					t.Errorf("expected only key 2 to remain, got %s", m.String())
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := New[int, int](tc.opts...)
			tc.check(t, m)
		})
	}
}

func TestLinkedHashMap_RemoveAndClear(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *LinkedHashMap[int, int])
	}{
		{
			name: "remove_existing_and_missing",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Remove(1)
				m.Remove(99)
				if m.Size() != 1 || m.ContainsKey(1) || !m.ContainsKey(2) {
					t.Errorf("unexpected state after Remove: %s", m.String())
				}
			},
		},
		{
			name: "clear",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Clear()
				if !m.Empty() || m.Size() != 0 {
					t.Errorf("expected empty map after Clear")
				}
				if slices.Collect(m.Keys()) != nil {
					t.Errorf("expected nil collected keys on cleared map")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := New[int, int]()
			tc.check(t, m)
		})
	}
}

func TestLinkedHashMap_Sequenced(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		opts      []Option
		wantPanic bool
		check     func(*testing.T, *LinkedHashMap[int, string])
	}{
		{
			name:      "first_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, m *LinkedHashMap[int, string]) {
				m.First()
			},
		},
		{
			name:      "last_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, m *LinkedHashMap[int, string]) {
				m.Last()
			},
		},
		{
			name:      "poll_first_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, m *LinkedHashMap[int, string]) {
				m.PollFirst()
			},
		},
		{
			name:      "poll_last_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T, m *LinkedHashMap[int, string]) {
				m.PollLast()
			},
		},
		{
			name: "first_and_last_non_empty",
			check: func(t *testing.T, m *LinkedHashMap[int, string]) {
				m.Put(10, "a")
				m.Put(20, "b")
				m.Put(30, "c")
				if k, v := m.First(); k != 10 || v != "a" {
					t.Errorf("First() = (%d, %q), want (10, 'a')", k, v)
				}
				if k, v := m.Last(); k != 30 || v != "c" {
					t.Errorf("Last() = (%d, %q), want (30, 'c')", k, v)
				}
			},
		},
		{
			name: "poll_first_and_poll_last",
			check: func(t *testing.T, m *LinkedHashMap[int, string]) {
				m.Put(10, "a")
				m.Put(20, "b")
				m.Put(30, "c")
				if k, v := m.PollFirst(); k != 10 || v != "a" || m.ContainsKey(10) {
					t.Errorf("PollFirst failed: (%d, %q)", k, v)
				}
				if k, v := m.PollLast(); k != 30 || v != "c" || m.ContainsKey(30) {
					t.Errorf("PollLast failed: (%d, %q)", k, v)
				}
				if m.Size() != 1 {
					t.Errorf("expected remaining size 1, got %d", m.Size())
				}
			},
		},
		{
			name: "put_first_and_put_last_insert_and_reposition",
			check: func(t *testing.T, m *LinkedHashMap[int, string]) {
				m.PutFirst(2, "B")
				m.PutFirst(1, "A")
				m.PutLast(3, "C")
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{1, 2, 3}) {
					t.Errorf("expected [1 2 3], got %v", got)
				}

				m.PutFirst(3, "C-updated")
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{3, 1, 2}) {
					t.Errorf("expected [3 1 2], got %v", got)
				}
				if v, _ := m.Get(3); v != "C-updated" {
					t.Errorf("expected C-updated, got %q", v)
				}

				m.PutLast(1, "A-updated")
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{3, 2, 1}) {
					t.Errorf("expected [3 2 1], got %v", got)
				}
				if v, _ := m.Get(1); v != "A-updated" {
					t.Errorf("expected A-updated, got %q", v)
				}
			},
		},
		{
			name: "put_first_and_put_last_eviction",
			opts: []Option{WithMaxElements(2)},
			check: func(t *testing.T, m *LinkedHashMap[int, string]) {
				m.PutFirst(1, "A")
				m.PutFirst(2, "B")
				m.PutFirst(3, "C") // Evicts tail (1)
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{3, 2}) {
					t.Errorf("expected [3 2], got %v", got)
				}
				m.PutLast(4, "D") // Evicts head (3)
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{2, 4}) {
					t.Errorf("expected [2 4], got %v", got)
				}
			},
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
			tc.check(t, m)
		})
	}
}

func TestLinkedHashMap_Iterators(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T)
	}{
		{
			name: "forward_iterators_full_and_early_break",
			check: func(t *testing.T) {
				m := New[int, string]()
				m.Put(1, "A")
				m.Put(2, "B")
				m.Put(3, "C")

				var keys []int
				var vals []string
				for k, v := range m.All() {
					keys = append(keys, k)
					vals = append(vals, v)
				}
				if !slices.Equal(keys, []int{1, 2, 3}) || !slices.Equal(vals, []string{"A", "B", "C"}) {
					t.Errorf("unexpected All(): keys=%v vals=%v", keys, vals)
				}
				if !slices.Equal(slices.Collect(m.Keys()), []int{1, 2, 3}) {
					t.Errorf("unexpected Keys()")
				}
				if !slices.Equal(slices.Collect(m.Values()), []string{"A", "B", "C"}) {
					t.Errorf("unexpected Values()")
				}

				c1 := 0
				for range m.All() {
					c1++
					break
				}
				c2 := 0
				for range m.Keys() {
					c2++
					break
				}
				c3 := 0
				for range m.Values() {
					c3++
					break
				}
				if c1 != 1 || c2 != 1 || c3 != 1 {
					t.Errorf("expected early break counts of 1, got %d/%d/%d", c1, c2, c3)
				}
			},
		},
		{
			name: "backward_iterators_full_and_early_break",
			check: func(t *testing.T) {
				m := New[int, string]()
				m.Put(1, "A")
				m.Put(2, "B")
				m.Put(3, "C")

				var keys []int
				var vals []string
				for k, v := range m.Backward() {
					keys = append(keys, k)
					vals = append(vals, v)
				}
				if !slices.Equal(keys, []int{3, 2, 1}) || !slices.Equal(vals, []string{"C", "B", "A"}) {
					t.Errorf("unexpected Backward(): keys=%v vals=%v", keys, vals)
				}
				if !slices.Equal(slices.Collect(m.BackwardKeys()), []int{3, 2, 1}) {
					t.Errorf("unexpected BackwardKeys()")
				}
				if !slices.Equal(slices.Collect(m.BackwardValues()), []string{"C", "B", "A"}) {
					t.Errorf("unexpected BackwardValues()")
				}

				c1 := 0
				for range m.Backward() {
					c1++
					break
				}
				c2 := 0
				for range m.BackwardKeys() {
					c2++
					break
				}
				c3 := 0
				for range m.BackwardValues() {
					c3++
					break
				}
				if c1 != 1 || c2 != 1 || c3 != 1 {
					t.Errorf("expected backward early break counts of 1, got %d/%d/%d", c1, c2, c3)
				}
			},
		},
		{
			name: "zero_value_iterators",
			check: func(t *testing.T) {
				var m LinkedHashMap[int, int]
				for range m.All() {
					t.Fatal("expected empty All() on zero value")
				}
				for range m.Keys() {
					t.Fatal("expected empty Keys() on zero value")
				}
				for range m.Values() {
					t.Fatal("expected empty Values() on zero value")
				}
				for range m.Backward() {
					t.Fatal("expected empty Backward() on zero value")
				}
				for range m.BackwardKeys() {
					t.Fatal("expected empty BackwardKeys() on zero value")
				}
				for range m.BackwardValues() {
					t.Fatal("expected empty BackwardValues() on zero value")
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

func TestLinkedHashMap_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		entries [][2]int
		want    string
	}{
		{
			name:    "empty",
			entries: nil,
			want:    "map[]",
		},
		{
			name:    "single",
			entries: [][2]int{{1, 10}},
			want:    "map[1:10]",
		},
		{
			name:    "multiple_insertion_order",
			entries: [][2]int{{1, 10}, {2, 20}, {3, 30}},
			want:    "map[1:10 2:20 3:30]",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := New[int, int]()
			for _, entry := range tc.entries {
				m.Put(entry[0], entry[1])
			}
			if got := m.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}
