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

func newSliceMap[E comparable, V any]() *sliceMap[E, V] {
	return &sliceMap[E, V]{
		root: &sliceNode[E, V]{},
	}
}

func (m *sliceMap[E, V]) Get(key []E) (V, bool) {
	node := m.getNode(key)
	if node != nil {
		return node.value, node.hasValue
	}
	var zero V
	return zero, false
}

func (m *sliceMap[E, V]) Put(key []E, value V) {
	node := m.root
	for i := 0; i < len(key); i++ {
		if node.children == nil {
			node.children = make(map[E]*sliceNode[E, V])
		}
		b := key[i]
		next, ok := node.children[b]
		if !ok {
			next = &sliceNode[E, V]{}
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

func (m *sliceMap[E, V]) Remove(key []E) {
	if m.removeNode(m.root, key, 0) {
		m.size--
	}
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

func (m *sliceMap[E, V]) Size() int {
	return m.size
}

func (m *sliceMap[E, V]) Empty() bool {
	return m.size == 0
}

func (m *sliceMap[E, V]) Clear() {
	m.root = &sliceNode[E, V]{}
	m.size = 0
}

func (m *sliceMap[E, V]) ContainsKey(key []E) bool {
	node := m.getNode(key)
	return node != nil && node.hasValue
}

func (m *sliceMap[E, V]) All() iter.Seq2[[]E, V] {
	return func(yield func([]E, V) bool) {
		m.iterate(m.root, nil, yield)
	}
}

func (m *sliceMap[E, V]) Keys() iter.Seq[[]E] {
	return func(yield func([]E) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func (m *sliceMap[E, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}

func cloneSlice[E any](s []E) []E {
	if s == nil {
		return nil
	}
	c := make([]E, len(s))
	copy(c, s)
	return c
}

// iterate performs a DFS. Iteration order is not guaranteed because map iteration is random.
func (m *sliceMap[E, V]) iterate(node *sliceNode[E, V], prefix []E, yield func([]E, V) bool) bool {
	if node.hasValue {
		// Yield a cloned slice so the caller cannot modify our internal state or see changes from subsequent iterations.
		if !yield(cloneSlice(prefix), node.value) {
			return false
		}
	}
	if node.children != nil {
		for k, next := range node.children {
			// prefix[:len(prefix):len(prefix)] forces append to allocate a new underlying array
			if !m.iterate(next, append(prefix[:len(prefix):len(prefix)], k), yield) {
				return false
			}
		}
	}
	return true
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

func (m *sliceMap[E, V]) HasPrefix(prefix []E) bool {
	node := m.getNode(prefix)
	if node == nil {
		return false
	}
	return m.hasAnyValue(node)
}

func (m *sliceMap[E, V]) hasAnyValue(node *sliceNode[E, V]) bool {
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

func (m *sliceMap[E, V]) EntriesWithPrefix(prefix []E) iter.Seq2[[]E, V] {
	return func(yield func([]E, V) bool) {
		node := m.getNode(prefix)
		if node == nil {
			return
		}
		m.iterate(node, prefix, yield)
	}
}

func (m *sliceMap[E, V]) KeysWithPrefix(prefix []E) iter.Seq[[]E] {
	return func(yield func([]E) bool) {
		for k := range m.EntriesWithPrefix(prefix) {
			if !yield(k) {
				return
			}
		}
	}
}

func (m *sliceMap[E, V]) ValuesWithPrefix(prefix []E) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.EntriesWithPrefix(prefix) {
			if !yield(v) {
				return
			}
		}
	}
}

func (m *sliceMap[E, V]) RemovePrefix(prefix []E) {
	if len(prefix) == 0 {
		m.Clear()
		return
	}
	m.removePrefixNode(m.root, prefix, 0)
}

func (m *sliceMap[E, V]) countValues(node *sliceNode[E, V]) int {
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

func (m *sliceMap[E, V]) removePrefixNode(node *sliceNode[E, V], prefix []E, depth int) bool {
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

func (m *sliceMap[E, V]) LongestPrefixOf(query []E) ([]E, V, bool) {
	node := m.root
	var longestKey []E
	var longestVal V
	var found bool

	if node.hasValue {
		longestKey = nil
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
			longestKey = cloneSlice(query[:i+1])
			longestVal = node.value
			found = true
		}
	}
	return longestKey, longestVal, found
}

func (m *sliceMap[E, V]) ShortestPrefixOf(query []E) ([]E, V, bool) {
	node := m.root
	if node.hasValue {
		return nil, node.value, true
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
			return cloneSlice(query[:i+1]), node.value, true
		}
	}
	var zero V
	return nil, zero, false
}

func (m *sliceMap[E, V]) PrefixesOf(query []E) iter.Seq2[[]E, V] {
	return func(yield func([]E, V) bool) {
		node := m.root
		if node.hasValue {
			if !yield(nil, node.value) {
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
				if !yield(cloneSlice(query[:i+1]), node.value) {
					return
				}
			}
		}
	}
}

func (m *sliceMap[E, V]) String() string {
	vals := make([]string, 0, m.Size())
	for k, v := range m.All() {
		vals = append(vals, fmt.Sprintf("%v:%v", k, v))
	}
	return "map[" + strings.Join(vals, " ") + "]"
}

// -----------------------------------------------------------------------------
// SliceSet Implementation
// -----------------------------------------------------------------------------

var (
	_ Set[[]int]                    = (*sliceSet[int])(nil)
	_ collections.MutableSet[[]int] = (*sliceSet[int])(nil)
	_ fmt.Stringer                  = (*sliceSet[int])(nil)
)

type sliceSet[E comparable] struct {
	m *sliceMap[E, struct{}]
}

func newSliceSet[E comparable]() *sliceSet[E] {
	return &sliceSet[E]{m: newSliceMap[E, struct{}]()}
}

func (s *sliceSet[E]) Size() int {
	return s.m.Size()
}

func (s *sliceSet[E]) Empty() bool {
	return s.m.Empty()
}

func (s *sliceSet[E]) Clear() {
	s.m.Clear()
}

func (s *sliceSet[E]) All() iter.Seq[[]E] {
	return s.m.Keys()
}

func (s *sliceSet[E]) Add(item []E) {
	s.m.Put(item, struct{}{})
}

func (s *sliceSet[E]) Remove() []E {
	for k := range s.m.Keys() {
		s.m.Remove(k)
		return k
	}
	panic("cannot remove from an empty set")
}

func (s *sliceSet[E]) AddAll(sequence iter.Seq[[]E]) {
	for t := range sequence {
		s.Add(t)
	}
}

func (s *sliceSet[E]) RemoveElement(item []E) {
	s.m.Remove(item)
}

func (s *sliceSet[E]) RemoveAll(other collections.Collection[[]E]) {
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

func (s *sliceSet[E]) RetainAll(other collections.Collection[[]E]) {
	newMap := newSliceMap[E, struct{}]()
	for t := range other.All() {
		if s.Contains(t) {
			newMap.Put(t, struct{}{})
		}
	}
	s.m = newMap
}

func (s *sliceSet[E]) Contains(item []E) bool {
	return s.m.ContainsKey(item)
}

func (s *sliceSet[E]) ContainsAll(other collections.Collection[[]E]) bool {
	for item := range other.All() {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

func (s *sliceSet[E]) HasPrefix(prefix []E) bool {
	return s.m.HasPrefix(prefix)
}

func (s *sliceSet[E]) ElementsWithPrefix(prefix []E) iter.Seq[[]E] {
	return s.m.KeysWithPrefix(prefix)
}

func (s *sliceSet[E]) RemovePrefix(prefix []E) {
	s.m.RemovePrefix(prefix)
}

func (s *sliceSet[E]) LongestPrefixOf(query []E) ([]E, bool) {
	k, _, ok := s.m.LongestPrefixOf(query)
	return k, ok
}

func (s *sliceSet[E]) ShortestPrefixOf(query []E) ([]E, bool) {
	k, _, ok := s.m.ShortestPrefixOf(query)
	return k, ok
}

func (s *sliceSet[E]) PrefixesOf(query []E) iter.Seq[[]E] {
	return func(yield func([]E) bool) {
		for k := range s.m.PrefixesOf(query) {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *sliceSet[E]) String() string {
	vals := make([]string, 0, s.Size())
	for item := range s.All() {
		vals = append(vals, fmt.Sprintf("%v", item))
	}
	return "[" + strings.Join(vals, " ") + "]"
}
