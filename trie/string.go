package trie

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/lock14/collections"
)

var (
	_ Map[string, int]                    = (*stringMap[int])(nil)
	_ collections.MutableMap[string, int] = (*stringMap[int])(nil)
	_ fmt.Stringer                        = (*stringMap[int])(nil)

	_ Set[string]                    = (*stringSet)(nil)
	_ collections.MutableSet[string] = (*stringSet)(nil)
	_ fmt.Stringer                   = (*stringSet)(nil)
)

type stringNode[V any] struct {
	children map[byte]*stringNode[V]
	value    V
	hasValue bool
}

type stringMap[V any] struct {
	root *stringNode[V]
	size int
}

type stringSet struct {
	m *stringMap[struct{}]
}

func newStringMap[V any]() *stringMap[V] {
	return &stringMap[V]{
		root: &stringNode[V]{},
	}
}

func newStringSet() *stringSet {
	return &stringSet{
		m: newStringMap[struct{}](),
	}
}

// Put associates the specified value with the specified key in the map.
func (m *stringMap[V]) Put(key string, value V) {
	node := m.root
	for i := 0; i < len(key); i++ {
		if node.children == nil {
			node.children = make(map[byte]*stringNode[V])
		}
		b := key[i]
		next, ok := node.children[b]
		if !ok {
			next = &stringNode[V]{}
			node.children[b] = next
		}
		node = next
	}
	if !node.hasValue {
		m.size++
		node.hasValue = true
	}
	node.value = value
}

// Remove removes the mapping for the specified key from the map if present.
func (m *stringMap[V]) Remove(key string) {
	if m.removeNode(m.root, key, 0) {
		m.size--
	}
}

// RemovePrefix removes all key-value pairs whose keys start with the given prefix.
func (m *stringMap[V]) RemovePrefix(prefix string) {
	if prefix == "" {
		m.Clear()
		return
	}
	m.removePrefixNode(m.root, prefix, 0)
}

// Clear removes all key-value pairs from the map.
func (m *stringMap[V]) Clear() {
	m.root = &stringNode[V]{}
	m.size = 0
}

// Get returns the value associated with the specified key, and a boolean indicating if it was found.
func (m *stringMap[V]) Get(key string) (V, bool) {
	node := m.getNode(key)
	if node != nil && node.hasValue {
		return node.value, true
	}
	var zero V
	return zero, false
}

// ContainsKey returns true if the map contains a mapping for the specified key.
func (m *stringMap[V]) ContainsKey(key string) bool {
	node := m.getNode(key)
	return node != nil && node.hasValue
}

// Size returns the number of key-value pairs in the map.
func (m *stringMap[V]) Size() int {
	return m.size
}

// Empty returns true if the map contains no key-value pairs.
func (m *stringMap[V]) Empty() bool {
	return m.size == 0
}

// HasPrefix returns true if there is at least one key in the map starting with the given prefix.
func (m *stringMap[V]) HasPrefix(prefix string) bool {
	if prefix == "" {
		return m.size > 0
	}
	return m.getNode(prefix) != nil
}

// LongestPrefixOf returns the longest key in the map that is a prefix of the given query.
func (m *stringMap[V]) LongestPrefixOf(query string) (string, V, bool) {
	node := m.root
	var longestKey string
	var longestVal V
	var found bool

	if node.hasValue {
		longestKey = ""
		longestVal = node.value
		found = true
	}

	for i := 0; i < len(query); i++ {
		if node.children == nil {
			break
		}
		next, ok := node.children[query[i]]
		if !ok {
			break
		}
		node = next
		if node.hasValue {
			longestKey = query[:i+1]
			longestVal = node.value
			found = true
		}
	}
	return longestKey, longestVal, found
}

// ShortestPrefixOf returns the shortest key in the map that is a prefix of the given query.
func (m *stringMap[V]) ShortestPrefixOf(query string) (string, V, bool) {
	node := m.root
	if node.hasValue {
		return "", node.value, true
	}
	for i := 0; i < len(query); i++ {
		if node.children == nil {
			break
		}
		next, ok := node.children[query[i]]
		if !ok {
			break
		}
		node = next
		if node.hasValue {
			return query[:i+1], node.value, true
		}
	}
	var zero V
	return "", zero, false
}

// All returns an iterator over all key-value pairs in the map in lexicographical key order.
func (m *stringMap[V]) All() iter.Seq2[string, V] {
	return func(yield func(string, V) bool) {
		if m.size == 0 {
			return
		}
		buf := make([]byte, 0, 32)
		m.iterate(m.root, &buf, yield)
	}
}

// Keys returns an iterator over all keys in the map in lexicographical order.
func (m *stringMap[V]) Keys() iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator over all values in the map in lexicographical key order.
func (m *stringMap[V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if m.size == 0 {
			return
		}
		m.iterateValues(m.root, yield)
	}
}

// EntriesWithPrefix returns an iterator over all key-value pairs where the key starts with the given prefix.
func (m *stringMap[V]) EntriesWithPrefix(prefix string) iter.Seq2[string, V] {
	return func(yield func(string, V) bool) {
		node := m.getNode(prefix)
		if node == nil {
			return
		}
		buf := make([]byte, len(prefix), len(prefix)+16)
		copy(buf, prefix)
		m.iterate(node, &buf, yield)
	}
}

// KeysWithPrefix returns an iterator over all keys that start with the given prefix.
func (m *stringMap[V]) KeysWithPrefix(prefix string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range m.EntriesWithPrefix(prefix) {
			if !yield(k) {
				return
			}
		}
	}
}

// ValuesWithPrefix returns an iterator over all values whose keys start with the given prefix.
func (m *stringMap[V]) ValuesWithPrefix(prefix string) iter.Seq[V] {
	return func(yield func(V) bool) {
		node := m.getNode(prefix)
		if node == nil {
			return
		}
		m.iterateValues(node, yield)
	}
}

// PrefixesOf returns an iterator over all key-value pairs in the map whose keys are a prefix of the given query.
func (m *stringMap[V]) PrefixesOf(query string) iter.Seq2[string, V] {
	return func(yield func(string, V) bool) {
		node := m.root
		if node.hasValue {
			if !yield("", node.value) {
				return
			}
		}
		for i := 0; i < len(query); i++ {
			if node.children == nil {
				break
			}
			next, ok := node.children[query[i]]
			if !ok {
				break
			}
			node = next
			if node.hasValue {
				if !yield(query[:i+1], node.value) {
					return
				}
			}
		}
	}
}

// String returns the string representation of the map.
func (m *stringMap[V]) String() string {
	if m.size == 0 {
		return "map[]"
	}
	var sb strings.Builder
	sb.WriteString("map[")
	first := true
	for k, v := range m.All() {
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		sb.WriteString(k)
		sb.WriteByte(':')
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte(']')
	return sb.String()
}

// Add inserts the specified element into the set.
func (s *stringSet) Add(item string) {
	s.m.Put(item, struct{}{})
}

// AddAll inserts all elements from the given sequence into the set.
func (s *stringSet) AddAll(sequence iter.Seq[string]) {
	for t := range sequence {
		s.Add(t)
	}
}

// Remove removes and returns an element from the set. It panics if the set is empty.
func (s *stringSet) Remove() string {
	for k := range s.m.Keys() {
		s.m.Remove(k)
		return k
	}
	panic("trie: cannot remove from an empty set")
}

// RemoveElement removes the specified element from the set if present.
func (s *stringSet) RemoveElement(item string) {
	s.m.Remove(item)
}

// RemoveAll removes all elements of the specified collection from this set.
func (s *stringSet) RemoveAll(other collections.Collection[string]) {
	if s.Empty() {
		return
	}
	if otherSet, ok := other.(*stringSet); ok && otherSet == s {
		s.Clear()
		return
	}
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

// RetainAll retains only the elements in this set that are contained in the specified collection.
func (s *stringSet) RetainAll(other collections.Collection[string]) {
	if s.Empty() {
		return
	}
	if otherSet, ok := other.(*stringSet); ok && otherSet == s {
		return
	}
	newMap := newStringMap[struct{}]()
	for t := range other.All() {
		if s.Contains(t) {
			newMap.Put(t, struct{}{})
		}
	}
	s.m = newMap
}

// RemovePrefix removes all elements that start with the given prefix.
func (s *stringSet) RemovePrefix(prefix string) {
	s.m.RemovePrefix(prefix)
}

// Clear removes all elements from the set.
func (s *stringSet) Clear() {
	s.m.Clear()
}

// Contains returns true if this set contains the specified element.
func (s *stringSet) Contains(item string) bool {
	return s.m.ContainsKey(item)
}

// ContainsAll returns true if this set contains all elements of the specified collection.
func (s *stringSet) ContainsAll(other collections.Collection[string]) bool {
	if otherSet, ok := other.(*stringSet); ok && otherSet == s {
		return true
	}
	for item := range other.All() {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// Size returns the number of elements in the set.
func (s *stringSet) Size() int {
	return s.m.Size()
}

// Empty returns true if the set contains no elements.
func (s *stringSet) Empty() bool {
	return s.m.Empty()
}

// HasPrefix returns true if there is at least one element in the set starting with the given prefix.
func (s *stringSet) HasPrefix(prefix string) bool {
	return s.m.HasPrefix(prefix)
}

// LongestPrefixOf returns the longest element in the set that is a prefix of the given query.
func (s *stringSet) LongestPrefixOf(query string) (string, bool) {
	k, _, ok := s.m.LongestPrefixOf(query)
	return k, ok
}

// ShortestPrefixOf returns the shortest element in the set that is a prefix of the given query.
func (s *stringSet) ShortestPrefixOf(query string) (string, bool) {
	k, _, ok := s.m.ShortestPrefixOf(query)
	return k, ok
}

// All returns an iterator over all elements in the set in lexicographical order.
func (s *stringSet) All() iter.Seq[string] {
	return s.m.Keys()
}

// ElementsWithPrefix returns an iterator over all elements that start with the given prefix.
func (s *stringSet) ElementsWithPrefix(prefix string) iter.Seq[string] {
	return s.m.KeysWithPrefix(prefix)
}

// PrefixesOf returns an iterator over all elements in the set that are a prefix of the given query.
func (s *stringSet) PrefixesOf(query string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range s.m.PrefixesOf(query) {
			if !yield(k) {
				return
			}
		}
	}
}

// String returns the string representation of the set.
func (s *stringSet) String() string {
	if s.Empty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for item := range s.All() {
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		sb.WriteString(item)
	}
	sb.WriteByte(']')
	return sb.String()
}

func (m *stringMap[V]) getNode(prefix string) *stringNode[V] {
	node := m.root
	for i := 0; i < len(prefix); i++ {
		if node.children == nil {
			return nil
		}
		next, ok := node.children[prefix[i]]
		if !ok {
			return nil
		}
		node = next
	}
	return node
}

func (m *stringMap[V]) removeNode(node *stringNode[V], key string, depth int) bool {
	if depth == len(key) {
		if !node.hasValue {
			return false
		}
		node.hasValue = false
		var zero V
		node.value = zero
		return true
	}
	if node.children == nil {
		return false
	}
	b := key[depth]
	next, ok := node.children[b]
	if !ok {
		return false
	}
	removed := m.removeNode(next, key, depth+1)
	if removed && !next.hasValue && len(next.children) == 0 {
		delete(node.children, b)
	}
	return removed
}

func (m *stringMap[V]) removePrefixNode(node *stringNode[V], prefix string, depth int) bool {
	if depth == len(prefix)-1 {
		if node.children == nil {
			return false
		}
		b := prefix[depth]
		target, ok := node.children[b]
		if !ok {
			return false
		}
		m.size -= m.countValues(target)
		delete(node.children, b)
		return true
	}
	if node.children == nil {
		return false
	}
	b := prefix[depth]
	next, ok := node.children[b]
	if !ok {
		return false
	}
	removed := m.removePrefixNode(next, prefix, depth+1)
	if removed && !next.hasValue && len(next.children) == 0 {
		delete(node.children, b)
	}
	return removed
}

func (m *stringMap[V]) countValues(node *stringNode[V]) int {
	count := 0
	if node.hasValue {
		count++
	}
	for _, child := range node.children {
		count += m.countValues(child)
	}
	return count
}

func (m *stringMap[V]) iterate(node *stringNode[V], buf *[]byte, yield func(string, V) bool) bool {
	if node.hasValue {
		if !yield(string(*buf), node.value) {
			return false
		}
	}
	if len(node.children) == 1 {
		for b, next := range node.children {
			*buf = append(*buf, b)
			ok := m.iterate(next, buf, yield)
			*buf = (*buf)[:len(*buf)-1]
			return ok
		}
	}
	if len(node.children) > 1 {
		var keys [256]byte
		n := 0
		for b := range node.children {
			keys[n] = b
			n++
		}
		slices.Sort(keys[:n])
		for _, b := range keys[:n] {
			*buf = append(*buf, b)
			if !m.iterate(node.children[b], buf, yield) {
				return false
			}
			*buf = (*buf)[:len(*buf)-1]
		}
	}
	return true
}

func (m *stringMap[V]) iterateValues(node *stringNode[V], yield func(V) bool) bool {
	if node.hasValue {
		if !yield(node.value) {
			return false
		}
	}
	if len(node.children) == 1 {
		for _, next := range node.children {
			return m.iterateValues(next, yield)
		}
	}
	if len(node.children) > 1 {
		var keys [256]byte
		n := 0
		for b := range node.children {
			keys[n] = b
			n++
		}
		slices.Sort(keys[:n])
		for _, b := range keys[:n] {
			if !m.iterateValues(node.children[b], yield) {
				return false
			}
		}
	}
	return true
}
