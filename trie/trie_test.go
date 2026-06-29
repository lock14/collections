package trie

import (
	"reflect"
	"slices"
	"testing"
)

type kv struct {
	k string
	v int
}

func TestStringMap_BasicOperations(t *testing.T) {
	cases := []struct {
		name      string
		puts      []kv
		removes   []string
		gets      []string
		wantVals  []int
		wantFound []bool
		wantSize  int
		wantEmpty bool
	}{
		{
			name:      "empty",
			gets:      []string{"a"},
			wantVals:  []int{0},
			wantFound: []bool{false},
			wantSize:  0,
			wantEmpty: true,
		},
		{
			name:      "single_element",
			puts:      []kv{{"hello", 1}},
			gets:      []string{"hello", "hell", "helloo"},
			wantVals:  []int{1, 0, 0},
			wantFound: []bool{true, false, false},
			wantSize:  1,
			wantEmpty: false,
		},
		{
			name:      "overwrite",
			puts:      []kv{{"hello", 1}, {"hello", 2}},
			gets:      []string{"hello"},
			wantVals:  []int{2},
			wantFound: []bool{true},
			wantSize:  1,
		},
		{
			name:      "remove",
			puts:      []kv{{"a", 1}, {"ab", 2}, {"abc", 3}},
			removes:   []string{"ab", "x"},
			gets:      []string{"a", "ab", "abc"},
			wantVals:  []int{1, 0, 3},
			wantFound: []bool{true, false, true},
			wantSize:  2,
		},
		{
			name:      "clear",
			puts:      []kv{{"a", 1}, {"b", 2}},
			removes:   []string{"_clear"},
			gets:      []string{"a", "b"},
			wantVals:  []int{0, 0},
			wantFound: []bool{false, false},
			wantSize:  0,
			wantEmpty: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := NewMap[int]()
			for _, p := range tc.puts {
				m.Put(p.k, p.v)
			}
			for _, r := range tc.removes {
				if r == "_clear" {
					m.Clear()
				} else {
					m.Remove(r)
				}
			}

			if m.Size() != tc.wantSize {
				t.Errorf("got size %d, want %d", m.Size(), tc.wantSize)
			}
			if m.Empty() != tc.wantEmpty {
				t.Errorf("got empty %v, want %v", m.Empty(), tc.wantEmpty)
			}

			for i, g := range tc.gets {
				v, ok := m.Get(g)
				if v != tc.wantVals[i] {
					t.Errorf("Get(%q) val = %d, want %d", g, v, tc.wantVals[i])
				}
				if ok != tc.wantFound[i] {
					t.Errorf("Get(%q) ok = %v, want %v", g, ok, tc.wantFound[i])
				}
				if m.ContainsKey(g) != tc.wantFound[i] {
					t.Errorf("ContainsKey(%q) = %v, want %v", g, m.ContainsKey(g), tc.wantFound[i])
				}
			}
		})
	}
}

func TestStringMap_PrefixOperations(t *testing.T) {
	cases := []struct {
		name               string
		puts               []kv
		hasPrefix          string
		wantHasPrefix      bool
		prefixQuery        string
		wantKeysWithPrefix []string
		longestQuery       string
		wantLongestKey     string
		wantLongestFound   bool
		shortestQuery      string
		wantShortestKey    string
		wantShortestFound  bool
		prefixesOfQuery    string
		wantPrefixesOf     []string
	}{
		{
			name:               "empty_trie",
			hasPrefix:          "a",
			wantHasPrefix:      false,
			prefixQuery:        "a",
			wantKeysWithPrefix: []string{},
			longestQuery:       "a",
			shortestQuery:      "a",
			prefixesOfQuery:    "a",
			wantPrefixesOf:     []string{},
		},
		{
			name:               "prefix_matches",
			puts:               []kv{{"app", 1}, {"apple", 2}, {"ape", 3}, {"bat", 4}},
			hasPrefix:          "ap",
			wantHasPrefix:      true,
			prefixQuery:        "ap",
			wantKeysWithPrefix: []string{"ape", "app", "apple"},
			longestQuery:       "apple pie",
			wantLongestKey:     "apple",
			wantLongestFound:   true,
			shortestQuery:      "apple pie",
			wantShortestKey:    "app",
			wantShortestFound:  true,
			prefixesOfQuery:    "apple pie",
			wantPrefixesOf:     []string{"app", "apple"},
		},
		{
			name:               "exact_match_as_prefix",
			puts:               []kv{{"cat", 1}},
			hasPrefix:          "cat",
			wantHasPrefix:      true,
			prefixQuery:        "cat",
			wantKeysWithPrefix: []string{"cat"},
			longestQuery:       "cat",
			wantLongestKey:     "cat",
			wantLongestFound:   true,
		},
		{
			name:               "no_prefix_match",
			puts:               []kv{{"dog", 1}},
			hasPrefix:          "cat",
			wantHasPrefix:      false,
			prefixQuery:        "cat",
			wantKeysWithPrefix: []string{},
			longestQuery:       "cat",
			wantLongestFound:   false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := NewMap[int]()
			for _, p := range tc.puts {
				m.Put(p.k, p.v)
			}

			if m.HasPrefix(tc.hasPrefix) != tc.wantHasPrefix {
				t.Errorf("HasPrefix(%q) = %v, want %v", tc.hasPrefix, m.HasPrefix(tc.hasPrefix), tc.wantHasPrefix)
			}

			gotKeys := slices.Collect(m.KeysWithPrefix(tc.prefixQuery))
			if !reflect.DeepEqual(gotKeys, tc.wantKeysWithPrefix) {
				// if both are empty slices, DeepEqual might fail if one is nil and one is empty, so handle that
				if len(gotKeys) != 0 || len(tc.wantKeysWithPrefix) != 0 {
					t.Errorf("KeysWithPrefix(%q) = %v, want %v", tc.prefixQuery, gotKeys, tc.wantKeysWithPrefix)
				}
			}

			lk, _, lok := m.LongestPrefixOf(tc.longestQuery)
			if lok != tc.wantLongestFound || (lok && lk != tc.wantLongestKey) {
				t.Errorf("LongestPrefixOf(%q) = (%q, %v), want (%q, %v)", tc.longestQuery, lk, lok, tc.wantLongestKey, tc.wantLongestFound)
			}

			sk, _, sok := m.ShortestPrefixOf(tc.shortestQuery)
			if sok != tc.wantShortestFound || (sok && sk != tc.wantShortestKey) {
				t.Errorf("ShortestPrefixOf(%q) = (%q, %v), want (%q, %v)", tc.shortestQuery, sk, sok, tc.wantShortestKey, tc.wantShortestFound)
			}

			var gotPrefixes []string
			for k := range m.PrefixesOf(tc.prefixesOfQuery) {
				gotPrefixes = append(gotPrefixes, k)
			}
			if !reflect.DeepEqual(gotPrefixes, tc.wantPrefixesOf) {
				if len(gotPrefixes) != 0 || len(tc.wantPrefixesOf) != 0 {
					t.Errorf("PrefixesOf(%q) = %v, want %v", tc.prefixesOfQuery, gotPrefixes, tc.wantPrefixesOf)
				}
			}
		})
	}
}

func TestStringMap_RemovePrefix(t *testing.T) {
	m := NewMap[int]()
	m.Put("app", 1)
	m.Put("apple", 2)
	m.Put("ape", 3)
	m.Put("bat", 4)

	m.RemovePrefix("ap")

	if m.Size() != 1 {
		t.Errorf("got size %d, want 1", m.Size())
	}
	if !m.ContainsKey("bat") {
		t.Errorf("expected bat to remain")
	}
	if m.ContainsKey("app") || m.ContainsKey("ape") {
		t.Errorf("expected ap* keys to be removed")
	}
}

func TestSliceMap_Basic(t *testing.T) {
	m := NewSliceMap[int, string]()
	m.Put([]int{1, 2, 3}, "a")
	m.Put([]int{1, 2}, "b")

	if m.Size() != 2 {
		t.Errorf("got size %d, want 2", m.Size())
	}

	v, ok := m.Get([]int{1, 2})
	if !ok || v != "b" {
		t.Errorf("Get([1,2]) = %v, %v, want b, true", v, ok)
	}

	lk, lv, lok := m.LongestPrefixOf([]int{1, 2, 3, 4})
	if !lok || lv != "a" || !reflect.DeepEqual(lk, []int{1, 2, 3}) {
		t.Errorf("LongestPrefixOf = %v, %v, %v", lk, lv, lok)
	}
}

func TestStringSet_Basic(t *testing.T) {
	s := NewSet()
	s.Add("hello")
	s.Add("world")

	if s.Size() != 2 {
		t.Errorf("size = %d", s.Size())
	}
	if !s.Contains("hello") {
		t.Errorf("Contains(hello) = false")
	}

	lk, lok := s.LongestPrefixOf("hello world")
	if !lok || lk != "hello" {
		t.Errorf("LongestPrefixOf = %v, %v", lk, lok)
	}
}

func TestStringMap_Values(t *testing.T) {
	m := NewMap[int]()
	m.Put("a", 1)
	m.Put("b", 2)

	vals := slices.Collect(m.Values())
	if !slices.Contains(vals, 1) || !slices.Contains(vals, 2) {
		t.Errorf("Values failed")
	}

	vp := slices.Collect(m.ValuesWithPrefix("a"))
	if len(vp) != 1 || vp[0] != 1 {
		t.Errorf("ValuesWithPrefix failed")
	}
}

func TestStringSet_Remaining(t *testing.T) {
	s := NewSet()
	s.Add("x")

	s2 := NewSet()
	s2.AddAll(s.All())
	if !s2.ContainsAll(s) {
		t.Errorf("ContainsAll failed")
	}

	s2.RemoveAll(s)
	if !s2.Empty() {
		t.Errorf("RemoveAll failed")
	}

	s2.Add("x")
	s2.Add("y")
	s2.RetainAll(s)
	if s2.Size() != 1 {
		t.Errorf("RetainAll failed")
	}

	v := s2.Remove()
	if v != "x" {
		t.Errorf("Remove failed")
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		s2.Remove()
	}()
}

func TestStringMap_EarlyExit(t *testing.T) {
	m := NewMap[int]()
	m.Put("a", 1)
	m.Put("ab", 2)
	m.Put("b", 3)

	count := 0
	for range m.All() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("All() early exit failed")
	}

	count = 0
	for range m.KeysWithPrefix("a") {
		count++
		break
	}
	if count != 1 {
		t.Errorf("KeysWithPrefix early exit failed")
	}

	count = 0
	for range m.PrefixesOf("abc") {
		count++
		break
	}
	if count != 1 {
		t.Errorf("PrefixesOf early exit failed")
	}
}

func TestSliceMap_EarlyExit(t *testing.T) {
	m := NewSliceMap[byte, int]()
	m.Put([]byte("a"), 1)
	m.Put([]byte("ab"), 2)
	m.Put([]byte("b"), 3)

	count := 0
	for range m.All() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("All() early exit failed")
	}

	count = 0
	for range m.KeysWithPrefix([]byte("a")) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("KeysWithPrefix early exit failed")
	}

	count = 0
	for range m.PrefixesOf([]byte("abc")) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("PrefixesOf early exit failed")
	}
}

func TestStringSet_EarlyExit(t *testing.T) {
	s := NewSet()
	s.Add("a")
	s.Add("b")
	count := 0
	for range s.All() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("All early exit failed")
	}
	count = 0
	for range s.PrefixesOf("abc") {
		count++
		break
	}
	if count != 1 {
		t.Errorf("PrefixesOf early exit failed")
	}
}

func TestSliceSet_EarlyExit(t *testing.T) {
	s := NewSliceSet[byte]()
	s.Add([]byte("a"))
	s.Add([]byte("b"))
	count := 0
	for range s.All() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("All early exit failed")
	}
	count = 0
	for range s.PrefixesOf([]byte("abc")) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("PrefixesOf early exit failed")
	}
}

func TestStringMap_EdgeCases(t *testing.T) {
	m := NewMap[int]()
	// ShortestPrefixOf missing characters / nil children
	_, _, ok := m.ShortestPrefixOf("abc")
	if ok {
		t.Errorf("expected false")
	}
	m.Put("a", 1)
	_, _, ok = m.ShortestPrefixOf("xyz")
	if ok {
		t.Errorf("expected false")
	}

	// RemovePrefix empty string clears map
	m.RemovePrefix("")
	if !m.Empty() {
		t.Errorf("expected empty map")
	}

	// RemovePrefix where prefix doesn't exist
	m.Put("a", 1)
	m.RemovePrefix("ab")
	if m.Size() != 1 {
		t.Errorf("size should be 1")
	}
	m.RemovePrefix("x")
	if m.Size() != 1 {
		t.Errorf("size should be 1")
	}
}

func TestSliceMap_EdgeCases(t *testing.T) {
	m := NewSliceMap[byte, int]()
	_, _, ok := m.ShortestPrefixOf([]byte("abc"))
	if ok {
		t.Errorf("expected false")
	}
	m.Put([]byte("a"), 1)
	_, _, ok = m.ShortestPrefixOf([]byte("xyz"))
	if ok {
		t.Errorf("expected false")
	}

	m.RemovePrefix(nil) // should clear map
	if !m.Empty() {
		t.Errorf("expected empty map")
	}
	m.Put([]byte("a"), 1)
	m.RemovePrefix([]byte{}) // also clears map
	if !m.Empty() {
		t.Errorf("expected empty map")
	}

	m.Put([]byte("a"), 1)
	m.RemovePrefix([]byte("ab"))
	if m.Size() != 1 {
		t.Errorf("size should be 1")
	}
	m.RemovePrefix([]byte("x"))
	if m.Size() != 1 {
		t.Errorf("size should be 1")
	}
}

func TestStringMap_MoreEdgeCases(t *testing.T) {
	m := NewMap[int]()
	m.Put("a", 1)
	// Try to remove "ab", depth reaches 'a', but 'b' is not in 'a's children
	m.RemovePrefix("ab")

	// Values tests with early exit
	count := 0
	for range m.Values() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("Values early exit failed")
	}

	count = 0
	for range m.ValuesWithPrefix("a") {
		count++
		break
	}
	if count != 1 {
		t.Errorf("ValuesWithPrefix early exit failed")
	}
}

func TestSliceMap_MoreEdgeCases(t *testing.T) {
	m := NewSliceMap[byte, int]()
	m.Put([]byte("a"), 1)
	m.RemovePrefix([]byte("ab"))

	count := 0
	for range m.Values() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("Values early exit failed")
	}

	count = 0
	for range m.ValuesWithPrefix([]byte("a")) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("ValuesWithPrefix early exit failed")
	}
}

func TestMoreCoverage(t *testing.T) {
	// StringMap edge cases
	m := NewMap[int]()
	// HasPrefix on nil children
	if m.HasPrefix("a") {
		t.Errorf("expected false")
	}
	m.Put("a", 1)
	if m.HasPrefix("ab") {
		t.Errorf("expected false")
	}

	// EntriesWithPrefix on missing
	count := 0
	for range m.EntriesWithPrefix("ab") {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0")
	}
	for range m.KeysWithPrefix("ab") {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0")
	}

	// SliceMap edge cases
	sm := NewSliceMap[byte, int]()
	if sm.HasPrefix([]byte("a")) {
		t.Errorf("expected false")
	}
	sm.Put([]byte("a"), 1)
	if sm.HasPrefix([]byte("ab")) {
		t.Errorf("expected false")
	}
	for range sm.EntriesWithPrefix([]byte("ab")) {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0")
	}
	for range sm.KeysWithPrefix([]byte("ab")) {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0")
	}

	// SliceMap missing key Get
	if _, ok := sm.Get([]byte("missing")); ok {
		t.Errorf("expected false")
	}

	// Root values
	m.Put("", 999)
	if m.Size() != 2 {
		t.Errorf("size should be 2")
	}
	if v, ok := m.Get(""); !ok || v != 999 {
		t.Errorf("Get empty string failed")
	}
	if !m.HasPrefix("") {
		t.Errorf("HasPrefix empty string failed")
	}
	lk, lv, lok := m.LongestPrefixOf("xyz")
	if !lok || lk != "" || lv != 999 {
		t.Errorf("LongestPrefixOf empty failed")
	}

	sm.Put([]byte{}, 999)
	if sm.Size() != 2 {
		t.Errorf("size should be 2")
	}
	if v, ok := sm.Get([]byte{}); !ok || v != 999 {
		t.Errorf("Get empty slice failed")
	}
	if !sm.HasPrefix([]byte{}) {
		t.Errorf("HasPrefix empty slice failed")
	}
	slk, slv, slok := sm.LongestPrefixOf([]byte("xyz"))
	if !slok || len(slk) != 0 || slv != 999 {
		t.Errorf("LongestPrefixOf empty slice failed")
	}

	// Remove root values
	m.Remove("")
	if m.Size() != 1 {
		t.Errorf("Remove empty string failed")
	}
	sm.Remove([]byte{})
	if sm.Size() != 1 {
		t.Errorf("Remove empty slice failed")
	}
}

func BenchmarkStringMap_Put(b *testing.B) {
	m := NewMap[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put("some_fairly_long_string_prefix", i)
	}
}

func BenchmarkSliceMap_Put(b *testing.B) {
	m := NewSliceMap[byte, int]()
	key := []byte("some_fairly_long_string_prefix")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(key, i)
	}
}

func BenchmarkStringMap_Get(b *testing.B) {
	m := NewMap[int]()
	m.Put("some_fairly_long_string_prefix", 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get("some_fairly_long_string_prefix")
	}
}

func BenchmarkSliceMap_Get(b *testing.B) {
	m := NewSliceMap[byte, int]()
	key := []byte("some_fairly_long_string_prefix")
	m.Put(key, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(key)
	}
}

func BenchmarkStringMap_KeysWithPrefix(b *testing.B) {
	m := NewMap[int]()
	m.Put("api/v1/users/1", 1)
	m.Put("api/v1/users/2", 2)
	m.Put("api/v1/posts/1", 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range m.KeysWithPrefix("api/v1/users") {
		}
	}
}
