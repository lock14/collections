package trie

import (
	"fmt"
	"slices"
	"testing"
)

func TestSliceMap_Comprehensive(t *testing.T) {
	m := NewSliceMap[byte, int]()

	// Test Empty, Size, Clear
	if !m.Empty() {
		t.Errorf("expected empty")
	}
	m.Put([]byte("hello"), 1)
	m.Put([]byte("world"), 2)

	if m.Size() != 2 {
		t.Errorf("expected size 2")
	}

	m.Clear()
	if !m.Empty() {
		t.Errorf("expected empty after clear")
	}

	m.Put([]byte("app"), 1)
	m.Put([]byte("apple"), 2)
	m.Put([]byte("ape"), 3)

	// Test ContainsKey, Get
	if !m.ContainsKey([]byte("app")) {
		t.Errorf("expected to contain app")
	}
	v, ok := m.Get([]byte("apple"))
	if !ok || v != 2 {
		t.Errorf("Get failed")
	}

	// Test Remove
	m.Remove([]byte("ape"))
	if m.ContainsKey([]byte("ape")) {
		t.Errorf("expected ape to be removed")
	}

	// Test Prefix Operations
	if !m.HasPrefix([]byte("ap")) {
		t.Errorf("expected HasPrefix true")
	}

	keys := slices.Collect(m.KeysWithPrefix([]byte("app")))
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}

	vals := slices.Collect(m.ValuesWithPrefix([]byte("app")))
	if len(vals) != 2 {
		t.Errorf("expected 2 vals, got %d", len(vals))
	}

	entries := make(map[string]int)
	for k, v := range m.EntriesWithPrefix([]byte("app")) {
		entries[string(k)] = v
	}
	if entries["app"] != 1 || entries["apple"] != 2 {
		t.Errorf("EntriesWithPrefix failed")
	}

	lk, lv, lok := m.LongestPrefixOf([]byte("apple pie"))
	if !lok || string(lk) != "apple" || lv != 2 {
		t.Errorf("LongestPrefixOf failed")
	}

	sk, sv, sok := m.ShortestPrefixOf([]byte("apple pie"))
	if !sok || string(sk) != "app" || sv != 1 {
		t.Errorf("ShortestPrefixOf failed")
	}

	var prefixes []string
	for k := range m.PrefixesOf([]byte("apple pie")) {
		prefixes = append(prefixes, string(k))
	}
	if len(prefixes) != 2 {
		t.Errorf("PrefixesOf failed")
	}

	// Test RemovePrefix
	m.RemovePrefix([]byte("ap"))
	if !m.Empty() {
		t.Errorf("expected empty after RemovePrefix")
	}
}

func TestSliceSet_Comprehensive(t *testing.T) {
	s := NewSliceSet[byte]()

	s.Add([]byte("hello"))
	s.Add([]byte("world"))

	if s.Size() != 2 {
		t.Errorf("expected size 2")
	}
	if !s.Contains([]byte("hello")) {
		t.Errorf("expected Contains true")
	}

	s.RemoveElement([]byte("world"))
	if s.Contains([]byte("world")) {
		t.Errorf("expected world to be removed")
	}

	s.Clear()
	if !s.Empty() {
		t.Errorf("expected empty")
	}

	s.Add([]byte("app"))
	s.Add([]byte("apple"))

	if !s.HasPrefix([]byte("ap")) {
		t.Errorf("expected HasPrefix true")
	}

	keys := slices.Collect(s.ElementsWithPrefix([]byte("app")))
	if len(keys) != 2 {
		t.Errorf("expected 2 keys")
	}

	lk, lok := s.LongestPrefixOf([]byte("apple pie"))
	if !lok || string(lk) != "apple" {
		t.Errorf("LongestPrefixOf failed")
	}

	sk, sok := s.ShortestPrefixOf([]byte("apple pie"))
	if !sok || string(sk) != "app" {
		t.Errorf("ShortestPrefixOf failed")
	}

	var prefixes []string
	for k := range s.PrefixesOf([]byte("apple pie")) {
		prefixes = append(prefixes, string(k))
	}
	if len(prefixes) != 2 {
		t.Errorf("PrefixesOf failed")
	}

	s.RemovePrefix([]byte("ap"))
	if !s.Empty() {
		t.Errorf("expected empty after RemovePrefix")
	}

	// Test String()
	s.Add([]byte("x"))
	str := fmt.Sprintf("%s", s)
	if str != "[[120]]" {
		t.Errorf("String failed, got %s", str)
	}
}

func TestStringSet_Comprehensive(t *testing.T) {
	s := NewSet()

	s.Add("hello")
	s.Add("world")

	if s.Size() != 2 {
		t.Errorf("expected size 2")
	}
	if !s.Contains("hello") {
		t.Errorf("expected Contains true")
	}

	s.RemoveElement("world")
	if s.Contains("world") {
		t.Errorf("expected world to be removed")
	}

	s.Clear()
	if !s.Empty() {
		t.Errorf("expected empty")
	}

	s.Add("app")
	s.Add("apple")

	if !s.HasPrefix("ap") {
		t.Errorf("expected HasPrefix true")
	}

	keys := slices.Collect(s.ElementsWithPrefix("app"))
	if len(keys) != 2 {
		t.Errorf("expected 2 keys")
	}

	lk, lok := s.LongestPrefixOf("apple pie")
	if !lok || lk != "apple" {
		t.Errorf("LongestPrefixOf failed")
	}

	sk, sok := s.ShortestPrefixOf("apple pie")
	if !sok || sk != "app" {
		t.Errorf("ShortestPrefixOf failed")
	}

	var prefixes []string
	for k := range s.PrefixesOf("apple pie") {
		prefixes = append(prefixes, string(k))
	}
	if len(prefixes) != 2 {
		t.Errorf("PrefixesOf failed")
	}

	s.RemovePrefix("ap")
	if !s.Empty() {
		t.Errorf("expected empty after RemovePrefix")
	}

	// Test String()
	s.Add("x")
	str := fmt.Sprintf("%s", s)
	if str != "[x]" {
		t.Errorf("String failed, got %s", str)
	}
}
