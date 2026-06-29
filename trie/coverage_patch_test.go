package trie

import (
	"testing"
)

func TestCoveragePatch(t *testing.T) {
	sm := NewSliceMap[byte, int]()
	m := NewMap[int]()

	// 1. LongestPrefixOf on empty slice
	slk, slv, slok := sm.LongestPrefixOf([]byte("xyz"))
	if slok || len(slk) != 0 || slv != 0 {
		t.Errorf("LongestPrefixOf empty slice failed")
	}

	// 2. Remove root values
	m.Put("", 1)
	m.Remove("")
	if m.Size() != 0 {
		t.Errorf("Remove empty string failed")
	}

	sm.Put([]byte{}, 1)
	sm.Remove([]byte{})
	if sm.Size() != 0 {
		t.Errorf("Remove empty slice failed")
	}

	// 3. Remove non-existent long keys to hit removeNode nil children / missing key
	m.Remove("ab")
	m.Remove("xyz")
	sm.Remove([]byte("ab"))
	sm.Remove([]byte("xyz"))

	// 4. Empty trie HasPrefix
	emptyStrMap := NewMap[int]()
	if emptyStrMap.HasPrefix("") {
		t.Errorf("empty map should not have empty prefix")
	}
	emptySliceMap := NewSliceMap[byte, int]()
	if emptySliceMap.HasPrefix([]byte{}) {
		t.Errorf("empty slice map should not have empty prefix")
	}

	// 5. EntriesWithPrefix early exit
	m.Put("abc", 1)
	count := 0
	for range m.EntriesWithPrefix("a") {
		count++
		break
	}
	if count != 1 {
		t.Errorf("EntriesWithPrefix early exit failed")
	}

	sm.Put([]byte("abc"), 1)
	count = 0
	for range sm.EntriesWithPrefix([]byte("a")) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("EntriesWithPrefix early exit failed")
	}

	// 6. PrefixesOf early exit
	count = 0
	for range m.PrefixesOf("abc") {
		count++
		break
	}
	if count != 1 {
		t.Errorf("PrefixesOf early exit failed")
	}

	count = 0
	for range sm.PrefixesOf([]byte("abc")) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("PrefixesOf early exit failed")
	}
}
