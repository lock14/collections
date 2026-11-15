// Package treemap provides a sorted, navigable B-Tree map implementation.
package treemap

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/lock14/collections"
	"github.com/lock14/collections/comparator"
)

var (
	_ collections.MutableNavigableMap[int, int] = (*TreeMap[int, int])(nil)
	_ fmt.Stringer                              = (*TreeMap[int, int])(nil)
)

const (
	// DefaultDegree is the default minimum degree (t) for the B-Tree.
	DefaultDegree = 32
)

// config holds the configuration values for constructing a TreeMap.
type config[K any] struct {
	degree     int
	comparator comparator.Comparator[K]
}

// Option configures a TreeMap.
type Option[K any] func(*config[K])

// WithDegree configures the minimum degree (t) of the B-Tree.
// Panics in New if degree is less than 2.
func WithDegree[K any](degree int) Option[K] {
	return func(c *config[K]) {
		c.degree = degree
	}
}

// WithComparator configures the key comparator used by the TreeMap.
// Panics in New if cmpFunc is nil.
func WithComparator[K any](cmpFunc comparator.Comparator[K]) Option[K] {
	return func(c *config[K]) {
		c.comparator = cmpFunc
	}
}

// TreeMap is a B-Tree backed map that implements collections.MutableNavigableMap.
type TreeMap[K any, V any] struct {
	root       *node[K, V]
	size       int
	degree     int
	comparator comparator.Comparator[K]
}

// New creates an empty TreeMap configured with the given options.
// Panics if degree is less than 2 or if no non-nil comparator is provided.
func New[K any, V any](opts ...Option[K]) *TreeMap[K, V] {
	cfg := &config[K]{
		degree: DefaultDegree,
	}
	for _, option := range opts {
		option(cfg)
	}
	if cfg.degree < 2 {
		panic("treemap: degree must be at least 2")
	}
	if cfg.comparator == nil {
		panic("treemap: comparator must be provided or use NewOrdered")
	}
	m := &TreeMap[K, V]{
		degree:     cfg.degree,
		comparator: cfg.comparator,
	}
	m.root = m.newNode(true)
	return m
}

// NewOrdered creates an empty TreeMap for cmp.Ordered keys using natural ascending order.
// Any options provided in opts are applied after the default natural-order comparator.
func NewOrdered[K cmp.Ordered, V any](opts ...Option[K]) *TreeMap[K, V] {
	allOpts := make([]Option[K], 0, len(opts)+1)
	allOpts = append(allOpts, WithComparator(comparator.NaturalOrder[K]()))
	allOpts = append(allOpts, opts...)
	return New[K, V](allOpts...)
}

// Put associates the specified value with the specified key in the map.
func (m *TreeMap[K, V]) Put(key K, value V) {
	m.put(key, value)
}

// Remove removes the mapping for the specified key from the map if present.
func (m *TreeMap[K, V]) Remove(key K) {
	m.remove(key)
}

// PollFirst removes and returns the smallest key-value pair in the map.
// Panics if the map is empty.
func (m *TreeMap[K, V]) PollFirst() (K, V) {
	if m.Empty() {
		panic("treemap: PollFirst called on empty map")
	}
	k, v := m.First()
	m.Remove(k)
	return k, v
}

// PollLast removes and returns the largest key-value pair in the map.
// Panics if the map is empty.
func (m *TreeMap[K, V]) PollLast() (K, V) {
	if m.Empty() {
		panic("treemap: PollLast called on empty map")
	}
	k, v := m.Last()
	m.Remove(k)
	return k, v
}

// Clear removes all key-value pairs from the map.
func (m *TreeMap[K, V]) Clear() {
	m.root = m.newNode(true)
	m.size = 0
}

// Get returns the value associated with the specified key, and a boolean indicating if it was found.
func (m *TreeMap[K, V]) Get(key K) (V, bool) {
	if m.root == nil {
		var zero V
		return zero, false
	}
	return m.get(m.root, key)
}

// ContainsKey returns true if the map contains a mapping for the specified key.
func (m *TreeMap[K, V]) ContainsKey(key K) bool {
	_, ok := m.Get(key)
	return ok
}

// Size returns the number of key-value pairs in the map.
func (m *TreeMap[K, V]) Size() int {
	return m.size
}

// Empty returns true if the map contains no key-value pairs.
func (m *TreeMap[K, V]) Empty() bool {
	return m.size == 0
}

// First returns the smallest key-value pair in the map without removing it.
// Panics if the map is empty.
func (m *TreeMap[K, V]) First() (K, V) {
	if m.Empty() {
		panic("treemap: First called on empty map")
	}
	n := m.root
	for !n.leaf {
		n = n.children[0]
	}
	return n.keys[0], n.values[0]
}

// Last returns the largest key-value pair in the map without removing it.
// Panics if the map is empty.
func (m *TreeMap[K, V]) Last() (K, V) {
	if m.Empty() {
		panic("treemap: Last called on empty map")
	}
	n := m.root
	for !n.leaf {
		n = n.children[len(n.children)-1]
	}
	return n.keys[len(n.keys)-1], n.values[len(n.values)-1]
}

// Lower returns the key-value pair for the greatest key strictly less than key, or false if not found.
func (m *TreeMap[K, V]) Lower(key K) (K, V, bool) {
	var bestK K
	var bestV V
	var found bool

	n := m.root
	for n != nil && len(n.keys) > 0 {
		i, _ := slices.BinarySearchFunc(n.keys, key, m.comparator)
		if i > 0 {
			bestK = n.keys[i-1]
			bestV = n.values[i-1]
			found = true
		}
		if n.leaf {
			break
		}
		n = n.children[i]
	}
	return bestK, bestV, found
}

// Floor returns the key-value pair for the greatest key less than or equal to key, or false if not found.
func (m *TreeMap[K, V]) Floor(key K) (K, V, bool) {
	var bestK K
	var bestV V
	var found bool

	n := m.root
	for n != nil && len(n.keys) > 0 {
		i, match := slices.BinarySearchFunc(n.keys, key, m.comparator)
		if match {
			return n.keys[i], n.values[i], true
		}
		if i > 0 {
			bestK = n.keys[i-1]
			bestV = n.values[i-1]
			found = true
		}
		if n.leaf {
			break
		}
		n = n.children[i]
	}
	return bestK, bestV, found
}

// Ceiling returns the key-value pair for the least key greater than or equal to key, or false if not found.
func (m *TreeMap[K, V]) Ceiling(key K) (K, V, bool) {
	var bestK K
	var bestV V
	var found bool

	n := m.root
	for n != nil && len(n.keys) > 0 {
		i, match := slices.BinarySearchFunc(n.keys, key, m.comparator)
		if match {
			return n.keys[i], n.values[i], true
		}
		if i < len(n.keys) {
			bestK = n.keys[i]
			bestV = n.values[i]
			found = true
		}
		if n.leaf {
			break
		}
		n = n.children[i]
	}
	return bestK, bestV, found
}

// Higher returns the key-value pair for the least key strictly greater than key, or false if not found.
func (m *TreeMap[K, V]) Higher(key K) (K, V, bool) {
	var bestK K
	var bestV V
	var found bool

	n := m.root
	for n != nil && len(n.keys) > 0 {
		i, match := slices.BinarySearchFunc(n.keys, key, m.comparator)
		if match {
			i++
		}
		if i < len(n.keys) {
			bestK = n.keys[i]
			bestV = n.values[i]
			found = true
		}
		if n.leaf {
			break
		}
		n = n.children[i]
	}
	return bestK, bestV, found
}

// All returns an iterator over all key-value pairs in ascending key order.
func (m *TreeMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.inOrder(m.root, yield)
	}
}

// Keys returns an iterator over all keys in ascending order.
func (m *TreeMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator over all values in ascending key order.
func (m *TreeMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}

// Backward returns an iterator over the key-value pairs in descending key order.
func (m *TreeMap[K, V]) Backward() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.reverseInOrder(m.root, yield)
	}
}

// BackwardKeys returns an iterator over the keys in descending order.
func (m *TreeMap[K, V]) BackwardKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m.Backward() {
			if !yield(k) {
				return
			}
		}
	}
}

// BackwardValues returns an iterator over the values in descending key order.
func (m *TreeMap[K, V]) BackwardValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.Backward() {
			if !yield(v) {
				return
			}
		}
	}
}

// From returns an iterator over the key-value pairs whose keys are greater than or equal to from.
func (m *TreeMap[K, V]) From(from K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var zero K
		m.rangeInOrder(m.root, from, zero, true, false, yield)
	}
}

// To returns an iterator over the key-value pairs whose keys are strictly less than to.
func (m *TreeMap[K, V]) To(to K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var zero K
		m.rangeInOrder(m.root, zero, to, false, true, yield)
	}
}

// Between returns an iterator over the key-value pairs in the half-open interval [from, to).
func (m *TreeMap[K, V]) Between(from, to K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m.root == nil || m.comparator(from, to) >= 0 {
			return
		}
		m.rangeInOrder(m.root, from, to, true, true, yield)
	}
}

// String returns a string representation of the map matching Go built-in map formatting.
func (m *TreeMap[K, V]) String() string {
	if m.Empty() {
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
		_, _ = fmt.Fprintf(&sb, "%v:%v", k, v)
	}
	sb.WriteByte(']')
	return sb.String()
}

func (m *TreeMap[K, V]) inOrder(n *node[K, V], yield func(K, V) bool) bool {
	if n == nil {
		return true
	}
	for i := 0; i < len(n.keys); i++ {
		if !n.leaf {
			if !m.inOrder(n.children[i], yield) {
				return false
			}
		}
		if !yield(n.keys[i], n.values[i]) {
			return false
		}
	}
	if !n.leaf {
		if !m.inOrder(n.children[len(n.keys)], yield) {
			return false
		}
	}
	return true
}

func (m *TreeMap[K, V]) reverseInOrder(n *node[K, V], yield func(K, V) bool) bool {
	if n == nil || len(n.keys) == 0 {
		return true
	}
	if !n.leaf {
		if !m.reverseInOrder(n.children[len(n.keys)], yield) {
			return false
		}
	}
	for i := len(n.keys) - 1; i >= 0; i-- {
		if !yield(n.keys[i], n.values[i]) {
			return false
		}
		if !n.leaf {
			if !m.reverseInOrder(n.children[i], yield) {
				return false
			}
		}
	}
	return true
}

func (m *TreeMap[K, V]) rangeInOrder(n *node[K, V], from, to K, checkFrom, checkTo bool, yield func(K, V) bool) bool {
	if n == nil || len(n.keys) == 0 {
		return true
	}
	startIdx := 0
	foundFrom := false
	if checkFrom {
		startIdx, foundFrom = slices.BinarySearchFunc(n.keys, from, m.comparator)
	}

	for i := startIdx; i < len(n.keys); i++ {
		if !n.leaf && !(i == startIdx && foundFrom) {
			cmpTo := -1
			if checkTo {
				cmpTo = m.comparator(n.keys[i], to)
			}
			if !m.rangeInOrder(n.children[i], from, to, checkFrom && i == startIdx, cmpTo >= 0, yield) {
				return false
			}
			if cmpTo >= 0 {
				return false
			}
		} else if checkTo && m.comparator(n.keys[i], to) >= 0 {
			return false
		}

		if !yield(n.keys[i], n.values[i]) {
			return false
		}
	}

	if !n.leaf {
		if !m.rangeInOrder(n.children[len(n.keys)], from, to, checkFrom && startIdx == len(n.keys), checkTo, yield) {
			return false
		}
	}
	return true
}
