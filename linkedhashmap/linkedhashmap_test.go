package linkedhashmap

import (
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LinkedHashMap[int, int])
	}{
		{
			name: "default",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				if m.Size() != 0 {
					t.Errorf("expected size 0")
				}
				if m.accessOrder {
					t.Errorf("expected insertion order by default")
				}
				// Should be able to add multiple elements
				m.Put(1, 10)
				m.Put(2, 20)
				if m.Size() != 2 {
					t.Errorf("expected size 2, got %d", m.Size())
				}
			},
		},
		{
			name: "access_order",
			opts: []Opt{WithAccessOrder()},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				if !m.accessOrder {
					t.Errorf("expected access order")
				}
			},
		},
		{
			name: "insertion_order",
			opts: []Opt{WithInsertionOrder()},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				if m.accessOrder {
					t.Errorf("expected insertion order")
				}
			},
		},
		{
			name: "max_elements",
			opts: []Opt{WithMaxElements(2)},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				if m.maxElements != 2 {
					t.Errorf("expected maxElements 2, got %d", m.maxElements)
				}
			},
		},
		{
			name: "capacity",
			opts: []Opt{WithCapacity(100)},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				// No exposed capacity check, but verifies option doesn't panic
				m.Put(1, 10)
				if m.Size() != 1 {
					t.Errorf("expected size 1")
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

func TestLinkedHashMap_Get(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LinkedHashMap[int, int])
	}{
		{
			name: "get_missing",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				_, ok := m.Get(1)
				if ok {
					t.Errorf("expected missing key")
				}
			},
		},
		{
			name: "get_insertion_order",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				val, ok := m.Get(1)
				if !ok || val != 10 {
					t.Errorf("expected 10")
				}

				keys := slices.Collect(m.Keys())
				if !slices.Equal(keys, []int{1, 2}) {
					t.Errorf("expected [1, 2], got %v", keys)
				}
			},
		},
		{
			name: "get_access_order",
			opts: []Opt{WithAccessOrder()},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				val, ok := m.Get(1)
				if !ok || val != 10 {
					t.Errorf("expected 10")
				}

				// 1 should now be the most recently accessed (at the tail)
				keys := slices.Collect(m.Keys())
				if !slices.Equal(keys, []int{2, 1}) {
					t.Errorf("expected [2, 1], got %v", keys)
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

func TestLinkedHashMap_Put(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LinkedHashMap[int, int])
	}{
		{
			name: "put_update_insertion_order",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Put(1, 100) // update

				if val, _ := m.Get(1); val != 100 {
					t.Errorf("expected update to 100")
				}

				keys := slices.Collect(m.Keys())
				if !slices.Equal(keys, []int{1, 2}) {
					t.Errorf("expected [1, 2]")
				}
			},
		},
		{
			name: "put_update_access_order",
			opts: []Opt{WithAccessOrder()},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Put(1, 100) // update

				// 1 should move to the tail
				keys := slices.Collect(m.Keys())
				if !slices.Equal(keys, []int{2, 1}) {
					t.Errorf("expected [2, 1]")
				}
			},
		},
		{
			name: "put_evicts_eldest",
			opts: []Opt{WithMaxElements(2)},
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Put(3, 30)

				if m.Size() != 2 {
					t.Errorf("expected size 2, got %d", m.Size())
				}

				keys := slices.Collect(m.Keys())
				if !slices.Equal(keys, []int{2, 3}) {
					t.Errorf("expected [2, 3], got %v", keys)
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

func TestLinkedHashMap_Remove(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LinkedHashMap[int, int])
	}{
		{
			name: "remove_existing",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Remove(1)

				if m.Size() != 1 {
					t.Errorf("expected size 1")
				}
				if m.ContainsKey(1) {
					t.Errorf("expected 1 to be removed")
				}

				keys := slices.Collect(m.Keys())
				if !slices.Equal(keys, []int{2}) {
					t.Errorf("expected [2], got %v", keys)
				}
			},
		},
		{
			name: "remove_missing",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Remove(2) // should not panic
				if m.Size() != 1 {
					t.Errorf("expected size 1")
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

func TestLinkedHashMap_Iterators(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LinkedHashMap[int, int])
	}{
		{
			name: "all_keys_values",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Put(3, 30)

				var keys []int
				var vals []int
				for k, v := range m.All() {
					keys = append(keys, k)
					vals = append(vals, v)
				}
				if !slices.Equal(keys, []int{1, 2, 3}) {
					t.Errorf("unexpected keys")
				}
				if !slices.Equal(vals, []int{10, 20, 30}) {
					t.Errorf("unexpected values")
				}

				if !slices.Equal(slices.Collect(m.Keys()), []int{1, 2, 3}) {
					t.Errorf("unexpected Keys()")
				}
				if !slices.Equal(slices.Collect(m.Values()), []int{10, 20, 30}) {
					t.Errorf("unexpected Values()")
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

func TestLinkedHashMap_Clear(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LinkedHashMap[int, int])
	}{
		{
			name: "clear",
			check: func(t *testing.T, m *LinkedHashMap[int, int]) {
				m.Put(1, 10)
				m.Put(2, 20)
				m.Clear()

				if !m.Empty() {
					t.Errorf("expected empty map")
				}
				if slices.Collect(m.Keys()) != nil {
					t.Errorf("expected no keys")
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

func TestLinkedHashMap_CoverageIters(t *testing.T) {
	hm := New[int, int]()
	hm.Put(1, 1)
	hm.Put(2, 2)
	for k := range hm.Keys() {
		_ = k
		break
	}
	for v := range hm.Values() {
		_ = v
		break
	}
}

func TestLinkedHashMap_Sequenced(t *testing.T) {
	t.Parallel()
	
	t.Run("first_last", func(t *testing.T) {
		m := New[int, string]()
		_, _, ok := m.First()
		if ok {
			t.Errorf("expected empty map to have no first")
		}
		_, _, ok = m.Last()
		if ok {
			t.Errorf("expected empty map to have no last")
		}
		
		m.Put(1, "A")
		m.Put(2, "B")
		m.Put(3, "C")
		
		k, v, ok := m.First()
		if !ok || k != 1 || v != "A" {
			t.Errorf("expected First to be (1, A), got (%v, %v, %v)", k, v, ok)
		}
		
		k, v, ok = m.Last()
		if !ok || k != 3 || v != "C" {
			t.Errorf("expected Last to be (3, C), got (%v, %v, %v)", k, v, ok)
		}
	})
	
	t.Run("poll_first_last", func(t *testing.T) {
		m := New[int, string]()
		_, _, ok := m.PollFirst()
		if ok {
			t.Errorf("expected false")
		}
		_, _, ok = m.PollLast()
		if ok {
			t.Errorf("expected false")
		}
		
		m.Put(1, "A")
		m.Put(2, "B")
		m.Put(3, "C")
		
		k, v, ok := m.PollFirst()
		if !ok || k != 1 || v != "A" || m.Size() != 2 {
			t.Errorf("expected PollFirst to be (1, A) and size 2")
		}
		
		k, v, ok = m.PollLast()
		if !ok || k != 3 || v != "C" || m.Size() != 1 {
			t.Errorf("expected PollLast to be (3, C) and size 1")
		}
	})
	
	t.Run("put_first_last", func(t *testing.T) {
		m := New[int, string]()
		m.PutFirst(2, "B")
		m.PutFirst(1, "A")
		m.PutLast(3, "C")
		
		got := slices.Collect(m.Keys())
		expected := []int{1, 2, 3}
		if !slices.Equal(got, expected) {
			t.Errorf("expected keys %v, got %v", expected, got)
		}
		
		// Update existing with PutFirst
		m.PutFirst(3, "C-updated")
		got = slices.Collect(m.Keys())
		expected = []int{3, 1, 2}
		if !slices.Equal(got, expected) {
			t.Errorf("expected keys %v, got %v", expected, got)
		}
		v, _ := m.Get(3)
		if v != "C-updated" {
			t.Errorf("expected value updated")
		}

		// Update existing with PutLast
		m.PutLast(1, "A-updated")
		got = slices.Collect(m.Keys())
		expected = []int{3, 2, 1}
		if !slices.Equal(got, expected) {
			t.Errorf("expected keys %v, got %v", expected, got)
		}
		v, _ = m.Get(1)
		if v != "A-updated" {
			t.Errorf("expected value updated")
		}
	})

	t.Run("put_first_eviction", func(t *testing.T) {
		m := New[int, string](WithMaxElements(2))
		m.PutFirst(1, "A")
		m.PutFirst(2, "B")
		m.PutFirst(3, "C") // Should evict tail (1, "A")

		got := slices.Collect(m.Keys())
		expected := []int{3, 2}
		if !slices.Equal(got, expected) {
			t.Errorf("expected keys %v, got %v", expected, got)
		}
	})
}

func TestLinkedHashMap_Reversed(t *testing.T) {
	t.Parallel()
	m := New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	m.Put(3, "C")

	var keys []int
	var values []string
	for k, v := range m.ReversedAll() {
		keys = append(keys, k)
		values = append(values, v)
	}
	if !slices.Equal(keys, []int{3, 2, 1}) {
		t.Errorf("expected ReversedAll keys [3, 2, 1], got %v", keys)
	}
	if !slices.Equal(values, []string{"C", "B", "A"}) {
		t.Errorf("expected ReversedAll values [C, B, A], got %v", values)
	}

	keys2 := slices.Collect(m.ReversedKeys())
	if !slices.Equal(keys2, []int{3, 2, 1}) {
		t.Errorf("expected ReversedKeys [3, 2, 1], got %v", keys2)
	}

	values2 := slices.Collect(m.ReversedValues())
	if !slices.Equal(values2, []string{"C", "B", "A"}) {
		t.Errorf("expected ReversedValues [C, B, A], got %v", values2)
	}
	
	// Test coverage for early exit iterator
	for k := range m.ReversedKeys() {
		_ = k
		break
	}
	for v := range m.ReversedValues() {
		_ = v
		break
	}
	for k, v := range m.ReversedAll() {
		_ = k
		_ = v
		break
	}
}
