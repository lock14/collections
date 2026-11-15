package trie

import (
	"slices"
	"testing"
)

type kv struct {
	k string
	v int
}

func TestStringMap_BasicAndRemove(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, Map[string, int])
	}{
		{
			name: "empty_state_and_missing_lookups",
			check: func(t *testing.T, m Map[string, int]) {
				if !m.Empty() || m.Size() != 0 {
					t.Errorf("expected empty map, got size %d", m.Size())
				}
				if v, ok := m.Get("a"); ok || v != 0 {
					t.Errorf("Get('a') = (%d, %v), want (0, false)", v, ok)
				}
				if m.ContainsKey("a") {
					t.Errorf("ContainsKey('a') should be false")
				}
				// Removing missing keys (empty, nil children, missing branch) is a no-op
				m.Remove("")
				m.Remove("missing")
				if !m.Empty() {
					t.Errorf("expected map to stay empty")
				}
			},
		},
		{
			name: "put_overwrite_empty_key_and_remove_pruning",
			check: func(t *testing.T, m Map[string, int]) {
				m.Put("", 99)
				m.Put("a", 1)
				m.Put("ab", 2)
				m.Put("abc", 3)
				m.Put("ab", 20) // overwrite existing key

				if m.Size() != 4 || m.Empty() {
					t.Fatalf("expected size 4, got %d", m.Size())
				}
				if v, ok := m.Get(""); !ok || v != 99 {
					t.Errorf("Get('') = (%d, %v)", v, ok)
				}
				if v, ok := m.Get("ab"); !ok || v != 20 {
					t.Errorf("Get('ab') = (%d, %v)", v, ok)
				}

				// Remove absent deeper key ("abcd" hits nil children, "ax" hits !ok, "a_absent" after removing "a")
				m.Remove("abcd")
				m.Remove("ax")
				if m.Size() != 4 {
					t.Errorf("expected size 4, got %d", m.Size())
				}

				// Remove intermediate key "ab" (keeps child "abc")
				m.Remove("ab")
				m.Remove("ab") // second remove hits !node.hasValue at depth==len(key)
				if m.ContainsKey("ab") || !m.ContainsKey("abc") || m.Size() != 3 {
					t.Errorf("unexpected state after removing 'ab': %s", m)
				}

				// Remove leaf "abc" (prunes 'c' and 'b' up to valued 'a')
				m.Remove("abc")
				if m.ContainsKey("abc") || !m.ContainsKey("a") || m.HasPrefix("ab") {
					t.Errorf("expected 'ab' branch pruned after removing 'abc'")
				}

				// Remove empty key ""
				m.Remove("")
				if m.ContainsKey("") || m.Size() != 1 {
					t.Errorf("expected size 1 after removing '', got %d", m.Size())
				}

				m.Clear()
				if !m.Empty() || m.Size() != 0 {
					t.Errorf("expected empty after Clear")
				}
			},
		},
		{
			name: "string_representation",
			check: func(t *testing.T, m Map[string, int]) {
				sm := m.(*stringMap[int])
				if got := sm.String(); got != "map[]" {
					t.Errorf("expected map[], got %s", got)
				}
				m.Put("b", 2)
				if got := sm.String(); got != "map[b:2]" {
					t.Errorf("expected map[b:2], got %s", got)
				}
				m.Put("a", 1)
				if got := sm.String(); got != "map[a:1 b:2]" {
					t.Errorf("expected map[a:1 b:2], got %s", got)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t, NewMap[int]())
		})
	}
}

func TestStringMap_PrefixOperationsAndIterators(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, Map[string, int])
	}{
		{
			name: "has_prefix_longest_shortest_and_prefixes_of",
			check: func(t *testing.T, m Map[string, int]) {
				if m.HasPrefix("") || m.HasPrefix("a") {
					t.Errorf("expected HasPrefix false on empty map")
				}
				if _, _, ok := m.LongestPrefixOf("abc"); ok {
					t.Errorf("expected LongestPrefixOf false on empty map")
				}
				if _, _, ok := m.ShortestPrefixOf("abc"); ok {
					t.Errorf("expected ShortestPrefixOf false on empty map")
				}

				for _, p := range []kv{{"app", 1}, {"apple", 2}, {"ape", 3}, {"bat", 4}} {
					m.Put(p.k, p.v)
				}
				if !m.HasPrefix("") || !m.HasPrefix("ap") || !m.HasPrefix("apple") || m.HasPrefix("applex") || m.HasPrefix("z") {
					t.Errorf("HasPrefix checks failed")
				}

				if lk, lv, ok := m.LongestPrefixOf("apple pie"); !ok || lk != "apple" || lv != 2 {
					t.Errorf("LongestPrefixOf = (%q, %d, %v)", lk, lv, ok)
				}
				if _, _, ok := m.LongestPrefixOf("zebra"); ok {
					t.Errorf("expected LongestPrefixOf('zebra') false")
				}
				if sk, sv, ok := m.ShortestPrefixOf("apple pie"); !ok || sk != "app" || sv != 1 {
					t.Errorf("ShortestPrefixOf = (%q, %d, %v)", sk, sv, ok)
				}
				if _, _, ok := m.ShortestPrefixOf("ap"); ok {
					t.Errorf("expected ShortestPrefixOf('ap') false")
				}
				if _, _, ok := m.ShortestPrefixOf("apple_extra"); !ok {
					t.Errorf("expected ShortestPrefixOf('apple_extra') true")
				}
				// Leaf nil-children break in ShortestPrefixOf when single char has nil children
				m2 := NewMap[int]()
				m2.Put("x", 1)
				if _, _, ok := m2.ShortestPrefixOf("y"); ok {
					t.Errorf("expected false")
				}
				m3 := NewMap[int]()
				m3.Put("xy", 1)
				if _, _, ok := m3.ShortestPrefixOf("x"); ok {
					t.Errorf("expected false on intermediate non-valued prefix")
				}
				// Put "" at root -> ShortestPrefixOf & LongestPrefixOf match ""
				m.Put("", 0)
				if sk, sv, ok := m.ShortestPrefixOf("anything"); !ok || sk != "" || sv != 0 {
					t.Errorf("ShortestPrefixOf with root key = (%q, %d, %v)", sk, sv, ok)
				}
				if lk, lv, ok := m.LongestPrefixOf("unrelated"); !ok || lk != "" || lv != 0 {
					t.Errorf("LongestPrefixOf with root key = (%q, %d, %v)", lk, lv, ok)
				}
			},
		},
		{
			name: "remove_prefix_all_branches",
			check: func(t *testing.T, m Map[string, int]) {
				// Empty map branches (depth==0 nil children & depth<len-1 nil children)
				m.RemovePrefix("x")
				m.RemovePrefix("xyz")

				m.Put("a", 1)
				// depth==len-1 nil children ("ab" when "a" is leaf)
				m.RemovePrefix("ab")
				// depth<len-1 nil children ("abc" when "a" is leaf)
				m.RemovePrefix("abc")
				if m.Size() != 1 {
					t.Errorf("expected size 1, got %d", m.Size())
				}

				m.Put("ac", 2)
				// depth==len-1 !ok ("ab" when "a" has child 'c')
				m.RemovePrefix("ab")
				// depth<len-1 !ok ("bx" when root only has 'a')
				m.RemovePrefix("bx")
				if m.Size() != 2 {
					t.Errorf("expected size 2, got %d", m.Size())
				}

				// Deep RemovePrefix that prunes unvalued intermediate ancestors
				m.Put("p/q/r1", 10)
				m.Put("p/q/r2", 20)
				m.RemovePrefix("p/q")
				if m.HasPrefix("p") || m.Size() != 2 {
					t.Errorf("expected 'p' ancestor pruned and size 2, got size=%d hasPrefix(p)=%v", m.Size(), m.HasPrefix("p"))
				}

				// Empty prefix clears map
				m.RemovePrefix("")
				if !m.Empty() {
					t.Errorf("expected empty map after RemovePrefix('')")
				}
			},
		},
		{
			name: "iterators_empty_missing_full_and_early_exits",
			check: func(t *testing.T, m Map[string, int]) {
				// Empty map iterators
				if slices.Collect(m.Keys()) != nil || slices.Collect(m.Values()) != nil {
					t.Errorf("expected empty iterators on empty map")
				}
				if slices.Collect(m.KeysWithPrefix("a")) != nil || slices.Collect(m.ValuesWithPrefix("a")) != nil {
					t.Errorf("expected empty prefix iterators on missing prefix")
				}

				m.Put("", 0)
				m.Put("a", 1)
				m.Put("ab", 2)
				m.Put("b", 3)

				if got := slices.Collect(m.Keys()); !slices.Equal(got, []string{"", "a", "ab", "b"}) {
					t.Errorf("Keys() = %v", got)
				}
				if got := slices.Collect(m.Values()); !slices.Equal(got, []int{0, 1, 2, 3}) {
					t.Errorf("Values() = %v", got)
				}
				if got := slices.Collect(m.KeysWithPrefix("a")); !slices.Equal(got, []string{"a", "ab"}) {
					t.Errorf("KeysWithPrefix('a') = %v", got)
				}
				if got := slices.Collect(m.ValuesWithPrefix("a")); !slices.Equal(got, []int{1, 2}) {
					t.Errorf("ValuesWithPrefix('a') = %v", got)
				}
				if slices.Collect(m.KeysWithPrefix("missing")) != nil || slices.Collect(m.ValuesWithPrefix("missing")) != nil {
					t.Errorf("expected nil on missing prefix")
				}

				// Early exits at stop=1 (root) and stop=2 (child) across All, Keys, Values, EntriesWithPrefix, KeysWithPrefix, ValuesWithPrefix, PrefixesOf
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
					for range m.PrefixesOf("abx") {
						c++
						if c == stop {
							break
						}
					}
					if c != stop {
						t.Errorf("PrefixesOf stop=%d got %d", stop, c)
					}
				}

				// PrefixesOf hitting nil children ("abxyz" after leaf "ab") and missing child ("ac" on non-leaf "a")
				var pKeys []string
				for k := range m.PrefixesOf("abxyz") {
					pKeys = append(pKeys, k)
				}
				if !slices.Equal(pKeys, []string{"", "a", "ab"}) {
					t.Errorf("PrefixesOf('abxyz') = %v", pKeys)
				}
				pKeys = nil
				for k := range m.PrefixesOf("ac") {
					pKeys = append(pKeys, k)
				}
				if !slices.Equal(pKeys, []string{"", "a"}) {
					t.Errorf("PrefixesOf('ac') = %v", pKeys)
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
				for range m.KeysWithPrefix("a") {
					c++
					break
				}
				if c != 1 {
					t.Errorf("KeysWithPrefix early exit failed")
				}
				c = 0
				for range m.ValuesWithPrefix("a") {
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
			tc.check(t, NewMap[int]())
		})
	}
}

func TestStringSet_Operations(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, Set[string])
	}{
		{
			name: "add_remove_contains_and_empty_panic",
			check: func(t *testing.T, s Set[string]) {
				if !s.Empty() || s.Size() != 0 {
					t.Errorf("expected empty set")
				}
				s.Add("hello")
				s.Add("world")
				if s.Size() != 2 || !s.Contains("hello") || s.Contains("missing") {
					t.Errorf("unexpected set state")
				}
				s.RemoveElement("world")
				if s.Contains("world") || s.Size() != 1 {
					t.Errorf("RemoveElement failed")
				}
				if got := s.Remove(); got != "hello" || !s.Empty() {
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
			check: func(t *testing.T, s Set[string]) {
				other := NewSet()
				other.Add("a")
				// Empty receiver fast paths
				s.RemoveAll(other)
				s.RetainAll(other)
				if !s.Empty() {
					t.Errorf("expected empty set")
				}

				s.AddAll(slices.Values([]string{"a", "b", "c"}))
				if !s.ContainsAll(s) || !s.ContainsAll(other) {
					t.Errorf("ContainsAll failed")
				}
				other.Add("z")
				if s.ContainsAll(other) {
					t.Errorf("ContainsAll non-subset should be false")
				}

				s.RetainAll(s)
				if s.Size() != 3 {
					t.Errorf("RetainAll(self) changed size to %d", s.Size())
				}

				keep := NewSet()
				keep.AddAll(slices.Values([]string{"b", "c", "x"}))
				s.RetainAll(keep)
				if !slices.Equal(slices.Collect(s.All()), []string{"b", "c"}) {
					t.Errorf("RetainAll failed: %v", slices.Collect(s.All()))
				}

				rem := NewSet()
				rem.Add("b")
				s.RemoveAll(rem)
				if !slices.Equal(slices.Collect(s.All()), []string{"c"}) {
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
			check: func(t *testing.T, s Set[string]) {
				ss := s.(*stringSet)
				if got := ss.String(); got != "[]" {
					t.Errorf("expected [], got %s", got)
				}
				s.AddAll(slices.Values([]string{"app", "apple", "bat"}))
				if got := ss.String(); got != "[app apple bat]" {
					t.Errorf("expected [app apple bat], got %s", got)
				}
				if !s.HasPrefix("ap") || s.HasPrefix("z") {
					t.Errorf("HasPrefix failed")
				}
				if lk, ok := s.LongestPrefixOf("apple pie"); !ok || lk != "apple" {
					t.Errorf("LongestPrefixOf = (%q, %v)", lk, ok)
				}
				if sk, ok := s.ShortestPrefixOf("apple pie"); !ok || sk != "app" {
					t.Errorf("ShortestPrefixOf = (%q, %v)", sk, ok)
				}
				if got := slices.Collect(s.ElementsWithPrefix("ap")); !slices.Equal(got, []string{"app", "apple"}) {
					t.Errorf("ElementsWithPrefix = %v", got)
				}
				if got := slices.Collect(s.PrefixesOf("apple pie")); !slices.Equal(got, []string{"app", "apple"}) {
					t.Errorf("PrefixesOf = %v", got)
				}

				// Early exits
				c1 := 0
				for range s.All() {
					c1++
					break
				}
				c2 := 0
				for range s.PrefixesOf("apple pie") {
					c2++
					break
				}
				if c1 != 1 || c2 != 1 {
					t.Errorf("early break failed: %d, %d", c1, c2)
				}

				s.RemovePrefix("ap")
				if s.Size() != 1 || !s.Contains("bat") {
					t.Errorf("RemovePrefix failed")
				}
				s.Clear()
				if !s.Empty() {
					t.Errorf("Clear failed")
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t, NewSet())
		})
	}
}
