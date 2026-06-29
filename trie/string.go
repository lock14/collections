package trie

import (
	"fmt"
	"iter"
	"strings"

	"github.com/lock14/collections"
)

var _ Map[string, int] = (*stringMap[int])(nil)

type stringNode[V any] struct {
	children map[byte]*stringNode[V]
	value    V
	hasValue bool
}

type stringMap[V any] struct {
	root *stringNode[V]
	size int
}

func newStringMap[V any]() *stringMap[V] {
	return &stringMap[V]{
		root: &stringNode[V]{},
	}
}

func (m *stringMap[V]) Get(key string) (V, bool) {
	node := m.getNode(key)
	if node != nil {
		return node.value, node.hasValue
	}
	var zero V
	return zero, false
}

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

func (m *stringMap[V]) Remove(key string) {
	if m.removeNode(m.root, key, 0) {
		m.size--
	}
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
	if removed {
		if !next.hasValue && len(next.children) == 0 {
			delete(node.children, b)
		}
	}
	return removed
}

func (m *stringMap[V]) Size() int {
	return m.size
}

func (m *stringMap[V]) Empty() bool {
	return m.size == 0
}

func (m *stringMap[V]) Clear() {
	m.root = &stringNode[V]{}
	m.size = 0
}

func (m *stringMap[V]) ContainsKey(key string) bool {
	node := m.getNode(key)
	return node != nil && node.hasValue
}

func (m *stringMap[V]) All() iter.Seq2[string, V] {
	return func(yield func(string, V) bool) {
		m.iterate(m.root, nil, yield)
	}
}

func (m *stringMap[V]) Keys() iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func (m *stringMap[V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}

// iterate performs a DFS and iterates over the nodes in lexicographical order.
func (m *stringMap[V]) iterate(node *stringNode[V], prefix []byte, yield func(string, V) bool) bool {
	if node.hasValue {
		if !yield(string(prefix), node.value) {
			return false
		}
	}
	if node.children != nil {
		for i := 0; i < 256; i++ {
			b := byte(i)
			if next, ok := node.children[b]; ok {
				if !m.iterate(next, append(prefix, b), yield) {
					return false
				}
			}
		}
	}
	return true
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

func (m *stringMap[V]) HasPrefix(prefix string) bool {
	node := m.getNode(prefix)
	if node == nil {
		return false
	}
	return m.hasAnyValue(node)
}

func (m *stringMap[V]) hasAnyValue(node *stringNode[V]) bool {
	if node.hasValue {
		return true
	}
	for _, child := range node.children {
		if m.hasAnyValue(child) {
			return true
		}
	}
	return false
}

func (m *stringMap[V]) EntriesWithPrefix(prefix string) iter.Seq2[string, V] {
	return func(yield func(string, V) bool) {
		node := m.getNode(prefix)
		if node == nil {
			return
		}
		m.iterate(node, []byte(prefix), yield)
	}
}

func (m *stringMap[V]) KeysWithPrefix(prefix string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range m.EntriesWithPrefix(prefix) {
			if !yield(k) {
				return
			}
		}
	}
}

func (m *stringMap[V]) ValuesWithPrefix(prefix string) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.EntriesWithPrefix(prefix) {
			if !yield(v) {
				return
			}
		}
	}
}

func (m *stringMap[V]) RemovePrefix(prefix string) {
	if prefix == "" {
		m.Clear()
		return
	}
	m.removePrefixNode(m.root, prefix, 0)
}

func (m *stringMap[V]) countValues(node *stringNode[V]) int {
	if node == nil {
		return 0
	}
	count := 0
	if node.hasValue {
		count++
	}
	for _, child := range node.children {
		count += m.countValues(child)
	}
	return count
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
		removedCount := m.countValues(target)
		m.size -= removedCount
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
	if removed {
		if !next.hasValue && len(next.children) == 0 {
			delete(node.children, b)
		}
	}
	return removed
}

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

// -----------------------------------------------------------------------------
// StringSet Implementation
// -----------------------------------------------------------------------------

var _ Set[string] = (*stringSet)(nil)

type stringSet struct {
	m *stringMap[struct{}]
}

func newStringSet() *stringSet {
	return &stringSet{m: newStringMap[struct{}]()}
}

func (s *stringSet) Size() int {
	return s.m.Size()
}

func (s *stringSet) Empty() bool {
	return s.m.Empty()
}

func (s *stringSet) Clear() {
	s.m.Clear()
}

func (s *stringSet) All() iter.Seq[string] {
	return s.m.Keys()
}

func (s *stringSet) Add(item string) {
	s.m.Put(item, struct{}{})
}

func (s *stringSet) Remove() string {
	for k := range s.m.Keys() {
		s.m.Remove(k)
		return k
	}
	panic("cannot remove from an empty set")
}

func (s *stringSet) AddAll(sequence iter.Seq[string]) {
	for t := range sequence {
		s.Add(t)
	}
}

func (s *stringSet) RemoveElement(item string) {
	s.m.Remove(item)
}

func (s *stringSet) RemoveAll(other collections.Collection[string]) {
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

func (s *stringSet) RetainAll(other collections.Collection[string]) {
	newMap := newStringMap[struct{}]()
	for t := range other.All() {
		if s.Contains(t) {
			newMap.Put(t, struct{}{})
		}
	}
	s.m = newMap
}

func (s *stringSet) Contains(item string) bool {
	return s.m.ContainsKey(item)
}

func (s *stringSet) ContainsAll(other collections.Collection[string]) bool {
	for item := range other.All() {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

func (s *stringSet) HasPrefix(prefix string) bool {
	return s.m.HasPrefix(prefix)
}

func (s *stringSet) ElementsWithPrefix(prefix string) iter.Seq[string] {
	return s.m.KeysWithPrefix(prefix)
}

func (s *stringSet) RemovePrefix(prefix string) {
	s.m.RemovePrefix(prefix)
}

func (s *stringSet) LongestPrefixOf(query string) (string, bool) {
	k, _, ok := s.m.LongestPrefixOf(query)
	return k, ok
}

func (s *stringSet) ShortestPrefixOf(query string) (string, bool) {
	k, _, ok := s.m.ShortestPrefixOf(query)
	return k, ok
}

func (s *stringSet) PrefixesOf(query string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range s.m.PrefixesOf(query) {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *stringSet) String() string {
	vals := make([]string, 0, s.Size())
	for item := range s.All() {
		vals = append(vals, fmt.Sprintf("%+v", item))
	}
	return "[" + strings.Join(vals, ", ") + "]"
}
