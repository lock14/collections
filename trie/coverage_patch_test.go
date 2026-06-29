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

	count = 0
	for range sm.PrefixesOf([]byte("abc")) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("PrefixesOf early exit failed")
	}

	// 7. Set interface methods coverage
	ss1 := NewSliceSet[byte]()
	ss2 := NewSliceSet[byte]()
	ss2.Add([]byte("a"))
	ss2.Add([]byte("b"))

	ss1.AddAll(ss2.All())
	if !ss1.ContainsAll(ss2) {
		t.Errorf("ContainsAll failed")
	}

	ss1.RemoveAll(ss2)
	if ss1.Size() != 0 {
		t.Errorf("RemoveAll failed")
	}

	ss1.Add([]byte("a"))
	ss1.Add([]byte("b"))
	ss1.Add([]byte("c"))
	ss3 := NewSliceSet[byte]()
	ss3.Add([]byte("a"))
	ss1.RetainAll(ss3)
	if ss1.Size() != 1 || !ss1.Contains([]byte("a")) {
		t.Errorf("RetainAll failed")
	}

	ss1.RemoveElement([]byte("a"))
	if ss1.Size() != 0 {
		t.Errorf("Remove failed")
	}

	// String set methods
	strSet1 := NewSet()
	strSet2 := NewSet()
	strSet2.Add("a")
	strSet1.AddAll(strSet2.All())
	if !strSet1.ContainsAll(strSet2) {
		t.Errorf("ContainsAll string failed")
	}

	// 8. PrefixesOf missing cases
	sm2 := NewSliceMap[byte, int]()
	sm2.Put([]byte("a"), 1)
	for range sm2.PrefixesOf([]byte("ab")) {
	} // no break

	sm3 := NewMap[int]()
	sm3.Put("a", 1)
	for range sm3.PrefixesOf("ab") {
	} // no break

	// 9. ShortestPrefixOf missing cases
	sm2.ShortestPrefixOf([]byte("x"))
	sm3.ShortestPrefixOf("x")
}
