package trie

import (
	"fmt"
	"iter"
	"strings"

	"github.com/lock14/collections"
)

var (
	_ Map[[]int, int]                    = (*sliceMap[int, int])(nil)
	_ collections.MutableMap[[]int, int] = (*sliceMap[int, int])(nil)
	_ fmt.Stringer                       = (*sliceMap[int, int])(nil)

	_ Set[[]int]                    = (*sliceSet[int])(nil)
	_ collections.MutableSet[[]int] = (*sliceSet[int])(nil)
	_ fmt.Stringer                  = (*sliceSet[int])(nil)
)

type sliceNode[E comparable, V any] struct {
	children map[E]*sliceNode[E, V]
	value    V
	hasValue bool
}

type sliceMap[E comparable, V any] struct {
	root *sliceNode[E, V]
	size int
}

type sliceSet[E comparable] struct {
	m *sliceMap[E, struct{}]
}

func newSliceMap[E comparable, V any]() *sliceMap[E, V] {
	return &sliceMap[E, V]{
		root: &sliceNode[E, V]{},
	}
}

func newSliceSet[E comparable]() *sliceSet[E] {
	return &sliceSet[E]{
		m: newSliceMap[E, struct{}](),
	}
}

// Put associates the specified value with the specified key in the map.
func (m *sliceMap[E, V]) Put(key []E, value V) {
	node := m.root
	for i := 0; i < len(key); i++ {
		if node.children == nil {
			node.children = make(map[E]*sliceNode[E, V])
		}
		elem := key[i]
		next, ok := node.children[elem]
		if !ok {
			next = &sliceNode[E, V]{}
			node.children[elem] = next
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
func (m *sliceMap[E, V]) Remove(key []E) {
	if m.removeNode(m.root, key, 0) {
		m.size--
	}
}

// RemovePrefix removes all key-value pairs whose keys start with the given prefix.
func (m *sliceMap[E, V]) RemovePrefix(prefix []E) {
	if len(prefix) == 0 {
		m.Clear()
		return
	}
	m.removePrefixNode(m.root, prefix, 0)
}

// Clear removes all key-value pairs from the map.
func (m *sliceMap[E, V]) Clear() {
	m.root = &sliceNode[E, V]{}
	m.size = 0
}

// Get returns the value associated with the specified key, and a boolean indicating if it was found.
func (m *sliceMap[E, V]) Get(key []E) (V, bool) {
	node := m.getNode(key)
	if node != nil && node.hasValue {
		return node.value, true
	}
	var zero V
	return zero, false
}

// ContainsKey returns true if the map contains a mapping for the specified key.
func (m *sliceMap[E, V]) ContainsKey(key []E) bool {
	node := m.getNode(key)
	return node != nil && node.hasValue
}

// Size returns the number of key-value pairs in the map.
func (m *sliceMap[E, V]) Size() int {
	return m.size
}

// Empty returns true if the map contains no key-value pairs.
func (m *sliceMap[E, V]) Empty() bool {
	return m.size == 0
}

// HasPrefix returns true if there is at least one key in the map starting with the given prefix.
func (m *sliceMap[E, V]) HasPrefix(prefix []E) bool {
	if len(prefix) == 0 {
		return m.size > 0
	}
	return m.getNode(prefix) != nil
}

// LongestPrefixOf returns the longest key in the map that is a prefix of the given query.
func (m *sliceMap[E, V]) LongestPrefixOf(query []E) ([]E, V, bool) {
	node := m.root
	longestLen := -1
	var longestVal V

	if node.hasValue {
		longestLen = 0
		longestVal = node.value
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
			longestLen = i + 1
			longestVal = node.value
		}
	}

	if longestLen < 0 {
		var zero V
		return nil, zero, false
	}
	longestKey := make([]E, longestLen)
	copy(longestKey, query[:longestLen])
	return longestKey, longestVal, true
}

// ShortestPrefixOf returns the shortest key in the map that is a prefix of the given query.
func (m *sliceMap[E, V]) ShortestPrefixOf(query []E) ([]E, V, bool) {
	node := m.root
	if node.hasValue {
		return []E{}, node.value, true
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
			keyCopy := make([]E, i+1)
			copy(keyCopy, query[:i+1])
			return keyCopy, node.value, true
		}
	}
	var zero V
	return nil, zero, false
}

// All returns an iterator over all key-value pairs in the map.
func (m *sliceMap[E, V]) All() iter.Seq2[[]E, V] {
	return func(yield func([]E, V) bool) {
		if m.size == 0 {
			return
		}
		buf := make([]E, 0, 8)
		m.iterate(m.root, &buf, yield)
	}
}

// Keys returns an iterator over all keys in the map.
func (m *sliceMap[E, V]) Keys() iter.Seq[[]E] {
	return func(yield func([]E) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator over all values in the map.
func (m *sliceMap[E, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if m.size == 0 {
			return
		}
		m.iterateValues(m.root, yield)
	}
}

// EntriesWithPrefix returns an iterator over all key-value pairs where the key starts with the given prefix.
func (m *sliceMap[E, V]) EntriesWithPrefix(prefix []E) iter.Seq2[[]E, V] {
	prefixSnap := make([]E, len(prefix))
	copy(prefixSnap, prefix)
	return func(yield func([]E, V) bool) {
		node := m.getNode(prefixSnap)
		if node == nil {
			return
		}
		buf := make([]E, len(prefixSnap), len(prefixSnap)+8)
		copy(buf, prefixSnap)
		m.iterate(node, &buf, yield)
	}
}

// KeysWithPrefix returns an iterator over all keys that start with the given prefix.
func (m *sliceMap[E, V]) KeysWithPrefix(prefix []E) iter.Seq[[]E] {
	entries := m.EntriesWithPrefix(prefix)
	return func(yield func([]E) bool) {
		for k := range entries {
			if !yield(k) {
				return
			}
		}
	}
}

// ValuesWithPrefix returns an iterator over all values whose keys start with the given prefix.
func (m *sliceMap[E, V]) ValuesWithPrefix(prefix []E) iter.Seq[V] {
	prefixSnap := make([]E, len(prefix))
	copy(prefixSnap, prefix)
	return func(yield func(V) bool) {
		node := m.getNode(prefixSnap)
		if node == nil {
			return
		}
		m.iterateValues(node, yield)
	}
}

// PrefixesOf returns an iterator over all key-value pairs in the map whose keys are a prefix of the given query.
func (m *sliceMap[E, V]) PrefixesOf(query []E) iter.Seq2[[]E, V] {
	querySnap := make([]E, len(query))
	copy(querySnap, query)
	return func(yield func([]E, V) bool) {
		node := m.root
		if node.hasValue {
			if !yield([]E{}, node.value) {
				return
			}
		}
		for i := 0; i < len(querySnap); i++ {
			if node.children == nil {
				break
			}
			next, ok := node.children[querySnap[i]]
			if !ok {
				break
			}
			node = next
			if node.hasValue {
				keyCopy := make([]E, i+1)
				copy(keyCopy, querySnap[:i+1])
				if !yield(keyCopy, node.value) {
					return
				}
			}
		}
	}
}

// String returns the string representation of the map.
func (m *sliceMap[E, V]) String() string {
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
		fmt.Fprintf(&sb, "%v:%v", k, v)
	}
	sb.WriteByte(']')
	return sb.String()
}

// Add inserts the specified element into the set.
func (s *sliceSet[E]) Add(item []E) {
	s.m.Put(item, struct{}{})
}

// AddAll inserts all elements from the given sequence into the set.
func (s *sliceSet[E]) AddAll(sequence iter.Seq[[]E]) {
	for t := range sequence {
		s.Add(t)
	}
}

// Remove removes and returns an element from the set. It panics if the set is empty.
func (s *sliceSet[E]) Remove() []E {
	for k := range s.m.Keys() {
		s.m.Remove(k)
		return k
	}
	panic("trie: cannot remove from an empty set")
}

// RemoveElement removes the specified element from the set if present.
func (s *sliceSet[E]) RemoveElement(item []E) {
	s.m.Remove(item)
}

// RemoveAll removes all elements of the specified collection from this set.
func (s *sliceSet[E]) RemoveAll(other collections.Collection[[]E]) {
	if s.Empty() {
		return
	}
	if otherSet, ok := other.(*sliceSet[E]); ok && otherSet == s {
		s.Clear()
		return
	}
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

// RetainAll retains only the elements in this set that are contained in the specified collection.
func (s *sliceSet[E]) RetainAll(other collections.Collection[[]E]) {
	if s.Empty() {
		return
	}
	if otherSet, ok := other.(*sliceSet[E]); ok && otherSet == s {
		return
	}
	newMap := newSliceMap[E, struct{}]()
	for t := range other.All() {
		if s.Contains(t) {
			newMap.Put(t, struct{}{})
		}
	}
	s.m = newMap
}

// RemovePrefix removes all elements that start with the given prefix.
func (s *sliceSet[E]) RemovePrefix(prefix []E) {
	s.m.RemovePrefix(prefix)
}

// Clear removes all elements from the set.
func (s *sliceSet[E]) Clear() {
	s.m.Clear()
}

// Contains returns true if this set contains the specified element.
func (s *sliceSet[E]) Contains(item []E) bool {
	return s.m.ContainsKey(item)
}

// ContainsAll returns true if this set contains all elements of the specified collection.
func (s *sliceSet[E]) ContainsAll(other collections.Collection[[]E]) bool {
	if otherSet, ok := other.(*sliceSet[E]); ok && otherSet == s {
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
func (s *sliceSet[E]) Size() int {
	return s.m.Size()
}

// Empty returns true if the set contains no elements.
func (s *sliceSet[E]) Empty() bool {
	return s.m.Empty()
}

// HasPrefix returns true if there is at least one element in the set starting with the given prefix.
func (s *sliceSet[E]) HasPrefix(prefix []E) bool {
	return s.m.HasPrefix(prefix)
}

// LongestPrefixOf returns the longest element in the set that is a prefix of the given query.
func (s *sliceSet[E]) LongestPrefixOf(query []E) ([]E, bool) {
	k, _, ok := s.m.LongestPrefixOf(query)
	return k, ok
}

// ShortestPrefixOf returns the shortest element in the set that is a prefix of the given query.
func (s *sliceSet[E]) ShortestPrefixOf(query []E) ([]E, bool) {
	k, _, ok := s.m.ShortestPrefixOf(query)
	return k, ok
}

// All returns an iterator over all elements in the set.
func (s *sliceSet[E]) All() iter.Seq[[]E] {
	return s.m.Keys()
}

// ElementsWithPrefix returns an iterator over all elements that start with the given prefix.
func (s *sliceSet[E]) ElementsWithPrefix(prefix []E) iter.Seq[[]E] {
	return s.m.KeysWithPrefix(prefix)
}

// PrefixesOf returns an iterator over all elements in the set that are a prefix of the given query.
func (s *sliceSet[E]) PrefixesOf(query []E) iter.Seq[[]E] {
	entries := s.m.PrefixesOf(query)
	return func(yield func([]E) bool) {
		for k := range entries {
			if !yield(k) {
				return
			}
		}
	}
}

// String returns the string representation of the set.
func (s *sliceSet[E]) String() string {
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
		fmt.Fprint(&sb, item)
	}
	sb.WriteByte(']')
	return sb.String()
}

func (m *sliceMap[E, V]) getNode(prefix []E) *sliceNode[E, V] {
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

func (m *sliceMap[E, V]) removeNode(node *sliceNode[E, V], key []E, depth int) bool {
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
	elem := key[depth]
	next, ok := node.children[elem]
	if !ok {
		return false
	}
	removed := m.removeNode(next, key, depth+1)
	if removed && !next.hasValue && len(next.children) == 0 {
		delete(node.children, elem)
	}
	return removed
}

func (m *sliceMap[E, V]) removePrefixNode(node *sliceNode[E, V], prefix []E, depth int) bool {
	if depth == len(prefix)-1 {
		if node.children == nil {
			return false
		}
		elem := prefix[depth]
		target, ok := node.children[elem]
		if !ok {
			return false
		}
		m.size -= m.countValues(target)
		delete(node.children, elem)
		return true
	}
	if node.children == nil {
		return false
	}
	elem := prefix[depth]
	next, ok := node.children[elem]
	if !ok {
		return false
	}
	removed := m.removePrefixNode(next, prefix, depth+1)
	if removed && !next.hasValue && len(next.children) == 0 {
		delete(node.children, elem)
	}
	return removed
}

func (m *sliceMap[E, V]) countValues(node *sliceNode[E, V]) int {
	count := 0
	if node.hasValue {
		count++
	}
	for _, child := range node.children {
		count += m.countValues(child)
	}
	return count
}

func (m *sliceMap[E, V]) iterate(node *sliceNode[E, V], buf *[]E, yield func([]E, V) bool) bool {
	if node.hasValue {
		keyCopy := make([]E, len(*buf))
		copy(keyCopy, *buf)
		if !yield(keyCopy, node.value) {
			return false
		}
	}
	for k, next := range node.children {
		*buf = append(*buf, k)
		if !m.iterate(next, buf, yield) {
			return false
		}
		*buf = (*buf)[:len(*buf)-1]
	}
	return true
}

func (m *sliceMap[E, V]) iterateValues(node *sliceNode[E, V], yield func(V) bool) bool {
	if node.hasValue {
		if !yield(node.value) {
			return false
		}
	}
	for _, next := range node.children {
		if !m.iterateValues(next, yield) {
			return false
		}
	}
	return true
}
