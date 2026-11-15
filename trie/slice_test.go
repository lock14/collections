package trie

import (
	"slices"
	"strings"
	"testing"
)

func TestSliceMap_BasicAndRemove(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, Map[[]byte, int])
	}{
		{
			name: "empty_state_and_missing_lookups",
			check: func(t *testing.T, m Map[[]byte, int]) {
				if !m.Empty() || m.Size() != 0 {
					t.Errorf("expected empty map, got size %d", m.Size())
				}
				if v, ok := m.Get([]byte("a")); ok || v != 0 {
					t.Errorf("Get('a') = (%d, %v)", v, ok)
				}
				if m.ContainsKey([]byte("a")) {
					t.Errorf("ContainsKey('a') should be false")
				}
				m.Remove(nil)
				m.Remove([]byte("missing"))
				if !m.Empty() {
					t.Errorf("expected map to remain empty")
				}
			},
		},
		{
			name: "put_overwrite_empty_key_and_remove_pruning",
			check: func(t *testing.T, m Map[[]byte, int]) {
				m.Put([]byte{}, 99)
				m.Put([]byte("a"), 1)
				m.Put([]byte("ab"), 2)
				m.Put([]byte("abc"), 3)
				m.Put([]byte("ab"), 20)

				if m.Size() != 4 || m.Empty() {
					t.Fatalf("expected size 4, got %d", m.Size())
				}
				if v, ok := m.Get(nil); !ok || v != 99 {
					t.Errorf("Get(nil) = (%d, %v)", v, ok)
				}
				if v, ok := m.Get([]byte("ab")); !ok || v != 20 {
					t.Errorf("Get('ab') = (%d, %v)", v, ok)
				}

				m.Remove([]byte("abcd"))
				m.Remove([]byte("ax"))
				if m.Size() != 4 {
					t.Errorf("expected size 4, got %d", m.Size())
				}

				m.Remove([]byte("ab"))
				m.Remove([]byte("ab"))
				if m.ContainsKey([]byte("ab")) || !m.ContainsKey([]byte("abc")) || m.Size() != 3 {
					t.Errorf("unexpected state after removing 'ab'")
				}

				m.Remove([]byte("abc"))
				if m.ContainsKey([]byte("abc")) || !m.ContainsKey([]byte("a")) || m.HasPrefix([]byte("ab")) {
					t.Errorf("expected 'ab' branch pruned after removing 'abc'")
				}

				m.Remove([]byte{})
				if m.ContainsKey([]byte{}) || m.Size() != 1 {
					t.Errorf("expected size 1 after removing empty slice, got %d", m.Size())
				}

				m.Clear()
				if !m.Empty() || m.Size() != 0 {
					t.Errorf("expected empty after Clear")
				}
			},
		},
		{
			name: "string_representation",
			check: func(t *testing.T, m Map[[]byte, int]) {
				sm := m.(*sliceMap[byte, int])
				if got := sm.String(); got != "map[]" {
					t.Errorf("expected map[], got %s", got)
				}
				m.Put([]byte{1}, 10)
				if got := sm.String(); got != "map[[1]:10]" {
					t.Errorf("expected map[[1]:10], got %s", got)
				}
				m.Put([]byte{2}, 20)
				got := sm.String()
				if got != "map[[1]:10 [2]:20]" && got != "map[[2]:20 [1]:10]" {
					t.Errorf("unexpected multi-entry String(): %s", got)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t, NewSliceMap[byte, int]())
		})
	}
}

func TestSliceMap_PrefixOperationsIteratorsAndAliasing(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, Map[[]byte, int])
	}{
		{
			name: "has_prefix_longest_shortest_and_prefixes_of",
			check: func(t *testing.T, m Map[[]byte, int]) {
				if m.HasPrefix(nil) || m.HasPrefix([]byte("a")) {
					t.Errorf("expected HasPrefix false on empty map")
				}
				if _, _, ok := m.LongestPrefixOf([]byte("abc")); ok {
					t.Errorf("expected LongestPrefixOf false on empty map")
				}
				if _, _, ok := m.ShortestPrefixOf([]byte("abc")); ok {
					t.Errorf("expected ShortestPrefixOf false on empty map")
				}

				m.Put([]byte("app"), 1)
				m.Put([]byte("apple"), 2)
				m.Put([]byte("ape"), 3)
				m.Put([]byte("bat"), 4)

				if !m.HasPrefix(nil) || !m.HasPrefix([]byte("ap")) || m.HasPrefix([]byte("applex")) || m.HasPrefix([]byte("z")) {
					t.Errorf("HasPrefix checks failed")
				}

				lk, lv, ok := m.LongestPrefixOf([]byte("apple pie"))
				if !ok || string(lk) != "apple" || lv != 2 {
					t.Errorf("LongestPrefixOf = (%q, %d, %v)", lk, lv, ok)
				}
				// Mutating returned slice must not corrupt trie
				lk[0] = 'X'
				if !m.ContainsKey([]byte("apple")) {
					t.Errorf("LongestPrefixOf returned aliased slice")
				}

				sk, sv, ok := m.ShortestPrefixOf([]byte("apple pie"))
				if !ok || string(sk) != "app" || sv != 1 {
					t.Errorf("ShortestPrefixOf = (%q, %d, %v)", sk, sv, ok)
				}
				sk[0] = 'X'
				if !m.ContainsKey([]byte("app")) {
					t.Errorf("ShortestPrefixOf returned aliased slice")
				}

				if _, _, ok := m.ShortestPrefixOf([]byte("ap")); ok {
					t.Errorf("expected false on non-valued prefix")
				}
				m2 := NewSliceMap[byte, int]()
				m2.Put([]byte("x"), 1)
				if _, _, ok := m2.ShortestPrefixOf([]byte("y")); ok {
					t.Errorf("expected false")
				}

				m.Put([]byte{}, 0)
				if rk, rv, rok := m.ShortestPrefixOf([]byte("anything")); !rok || len(rk) != 0 || rv != 0 {
					t.Errorf("ShortestPrefixOf root = (%v, %d, %v)", rk, rv, rok)
				}
				if rk, rv, rok := m.LongestPrefixOf([]byte("unrelated")); !rok || len(rk) != 0 || rv != 0 {
					t.Errorf("LongestPrefixOf root = (%v, %d, %v)", rk, rv, rok)
				}
			},
		},
		{
			name: "remove_prefix_all_branches",
			check: func(t *testing.T, m Map[[]byte, int]) {
				m.RemovePrefix([]byte("x"))
				m.RemovePrefix([]byte("xyz"))

				m.Put([]byte("a"), 1)
				m.RemovePrefix([]byte("ab"))
				m.RemovePrefix([]byte("abc"))
				if m.Size() != 1 {
					t.Errorf("expected size 1, got %d", m.Size())
				}

				m.Put([]byte("ac"), 2)
				m.RemovePrefix([]byte("ab"))
				m.RemovePrefix([]byte("bx"))
				if m.Size() != 2 {
					t.Errorf("expected size 2, got %d", m.Size())
				}

				m.Put([]byte("p/q/r1"), 10)
				m.Put([]byte("p/q/r2"), 20)
				m.RemovePrefix([]byte("p/q"))
				if m.HasPrefix([]byte("p")) || m.Size() != 2 {
					t.Errorf("expected 'p' branch pruned")
				}

				m.RemovePrefix(nil)
				if !m.Empty() {
					t.Errorf("expected empty after RemovePrefix(nil)")
				}
			},
		},
		{
			name: "iterators_early_exits_and_caller_slice_mutation_safety",
			check: func(t *testing.T, m Map[[]byte, int]) {
				if slices.Collect(m.Keys()) != nil || slices.Collect(m.Values()) != nil {
					t.Errorf("expected empty iterators on empty map")
				}
				if slices.Collect(m.KeysWithPrefix([]byte("a"))) != nil || slices.Collect(m.ValuesWithPrefix([]byte("a"))) != nil {
					t.Errorf("expected empty prefix iterators on empty map")
				}

				keyBuf := []byte("app")
				m.Put(keyBuf, 1)
				keyBuf[0] = 'Z' // Mutate caller slice after Put
				if !m.ContainsKey([]byte("app")) {
					t.Fatalf("Put aliased caller slice")
				}
				m.Put([]byte{}, 0)
				m.Put([]byte("apple"), 2)
				m.Put([]byte("bat"), 3)

				// Mutate prefix slice BEFORE and DURING prefix iteration!
				prefix := []byte("app")
				entriesSeq := m.EntriesWithPrefix(prefix)
				valsSeq := m.ValuesWithPrefix(prefix)
				prefix[0] = 'X'

				var gotEntries []string
				for k := range entriesSeq {
					gotEntries = append(gotEntries, string(k))
					k[0] = 'Y' // Mutate yielded key slice during iteration!
				}
				slices.Sort(gotEntries)
				if !slices.Equal(gotEntries, []string{"app", "apple"}) {
					t.Errorf("EntriesWithPrefix after caller mutation = %v", gotEntries)
				}
				gotVals := slices.Collect(valsSeq)
				slices.Sort(gotVals)
				if !slices.Equal(gotVals, []int{1, 2}) {
					t.Errorf("ValuesWithPrefix after caller mutation = %v", gotVals)
				}

				// Mutate query slice BEFORE and DURING PrefixesOf!
				query := []byte("apple pie")
				prefSeq := m.PrefixesOf(query)
				query[0] = 'Q'
				var gotPref []string
				for k := range prefSeq {
					gotPref = append(gotPref, string(k))
					if len(k) > 0 {
						k[0] = 'W'
					}
				}
				if !slices.Equal(gotPref, []string{"", "app", "apple"}) {
					t.Errorf("PrefixesOf after caller mutation = %v", gotPref)
				}

				// Early exits at stop=1 and stop=2
				for _, stop := range []int{1, 2} {
					c := 0
					for range m.All() {
						c++
						if c == stop {
							break
						}
					}
					if c != stop {
						t.Errorf("All stop=%d got %d", stop, c)
					}

					c = 0
					for range m.Values() {
						c++
						if c == stop {
							break
						}
					}
					if c != stop {
						t.Errorf("Values stop=%d got %d", stop, c)
					}

					c = 0
					for range m.PrefixesOf([]byte("apple_xyz")) {
						c++
						if c == stop {
							break
						}
					}
					if c != stop {
						t.Errorf("PrefixesOf stop=%d got %d", stop, c)
					}
				}

				// Exhaust PrefixesOf past leaf to cover node.children == nil break
				var allPref []string
				for k := range m.PrefixesOf([]byte("apple_xyz")) {
					allPref = append(allPref, string(k))
				}
				if len(allPref) != 3 {
					t.Errorf("expected 3 prefixes, got %v", allPref)
				}
				for range m.PrefixesOf([]byte("zzz")) {
				}

				c := 0
				for range m.Keys() {
					c++
					break
				}
				if c != 1 {
					t.Errorf("Keys early exit failed")
				}
				c = 0
				for range m.KeysWithPrefix([]byte("app")) {
					c++
					break
				}
				if c != 1 {
					t.Errorf("KeysWithPrefix early exit failed")
				}
				c = 0
				for range m.ValuesWithPrefix([]byte("app")) {
					c++
					break
				}
				if c != 1 {
					t.Errorf("ValuesWithPrefix early exit failed")
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t, NewSliceMap[byte, int]())
		})
	}
}

func TestSliceSet_Operations(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, Set[[]byte])
	}{
		{
			name: "add_remove_contains_and_empty_panic",
			check: func(t *testing.T, s Set[[]byte]) {
				if !s.Empty() || s.Size() != 0 {
					t.Errorf("expected empty set")
				}
				s.Add([]byte("hello"))
				s.Add([]byte("world"))
				if s.Size() != 2 || !s.Contains([]byte("hello")) || s.Contains([]byte("missing")) {
					t.Errorf("unexpected set state")
				}
				s.RemoveElement([]byte("world"))
				if s.Contains([]byte("world")) || s.Size() != 1 {
					t.Errorf("RemoveElement failed")
				}
				if got := s.Remove(); string(got) != "hello" || !s.Empty() {
					t.Errorf("Remove() = %q", got)
				}
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic on empty Remove()")
					}
				}()
				s.Remove()
			},
		},
		{
			name: "bulk_operations_and_self_aliasing",
			check: func(t *testing.T, s Set[[]byte]) {
				other := NewSliceSet[byte]()
				other.Add([]byte("a"))
				s.RemoveAll(other)
				s.RetainAll(other)
				if !s.Empty() {
					t.Errorf("expected empty set")
				}

				s.AddAll(slices.Values([][]byte{[]byte("a"), []byte("b"), []byte("c")}))
				if !s.ContainsAll(s) || !s.ContainsAll(other) {
					t.Errorf("ContainsAll failed")
				}
				other.Add([]byte("z"))
				if s.ContainsAll(other) {
					t.Errorf("ContainsAll non-subset should be false")
				}

				s.RetainAll(s)
				if s.Size() != 3 {
					t.Errorf("RetainAll(self) changed size to %d", s.Size())
				}

				keep := NewSliceSet[byte]()
				keep.AddAll(slices.Values([][]byte{[]byte("b"), []byte("c"), []byte("x")}))
				s.RetainAll(keep)
				if s.Size() != 2 || !s.Contains([]byte("b")) || !s.Contains([]byte("c")) {
					t.Errorf("RetainAll failed")
				}

				rem := NewSliceSet[byte]()
				rem.Add([]byte("b"))
				s.RemoveAll(rem)
				if s.Size() != 1 || !s.Contains([]byte("c")) {
					t.Errorf("RemoveAll failed")
				}

				s.RemoveAll(s)
				if !s.Empty() {
					t.Errorf("RemoveAll(self) should empty set")
				}
			},
		},
		{
			name: "prefix_operations_iterators_and_string",
			check: func(t *testing.T, s Set[[]byte]) {
				ss := s.(*sliceSet[byte])
				if got := ss.String(); got != "[]" {
					t.Errorf("expected [], got %s", got)
				}
				s.Add([]byte{120})
				if got := ss.String(); got != "[[120]]" {
					t.Errorf("expected [[120]], got %s", got)
				}
				s.Add([]byte{121})
				gotStr := ss.String()
				if !strings.HasPrefix(gotStr, "[") || !strings.Contains(gotStr, "[120]") || !strings.Contains(gotStr, "[121]") {
					t.Errorf("unexpected multi-element String(): %s", gotStr)
				}
				s.Clear()

				s.AddAll(slices.Values([][]byte{[]byte("app"), []byte("apple"), []byte("bat")}))
				if !s.HasPrefix([]byte("ap")) || s.HasPrefix([]byte("z")) {
					t.Errorf("HasPrefix failed")
				}
				if lk, ok := s.LongestPrefixOf([]byte("apple pie")); !ok || string(lk) != "apple" {
					t.Errorf("LongestPrefixOf = (%q, %v)", lk, ok)
				}
				if sk, ok := s.ShortestPrefixOf([]byte("apple pie")); !ok || string(sk) != "app" {
					t.Errorf("ShortestPrefixOf = (%q, %v)", sk, ok)
				}
				if len(slices.Collect(s.ElementsWithPrefix([]byte("ap")))) != 2 {
					t.Errorf("ElementsWithPrefix failed")
				}
				if len(slices.Collect(s.PrefixesOf([]byte("apple pie")))) != 2 {
					t.Errorf("PrefixesOf failed")
				}

				c1 := 0
				for range s.All() {
					c1++
					break
				}
				c2 := 0
				for range s.PrefixesOf([]byte("apple pie")) {
					c2++
					break
				}
				if c1 != 1 || c2 != 1 {
					t.Errorf("early break failed: %d, %d", c1, c2)
				}

				s.RemovePrefix([]byte("ap"))
				if s.Size() != 1 || !s.Contains([]byte("bat")) {
					t.Errorf("RemovePrefix failed")
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t, NewSliceSet[byte]())
		})
	}
}
