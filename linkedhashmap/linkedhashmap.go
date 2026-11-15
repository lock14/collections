// Package linkedhashmap provides a hash map that preserves insertion or access order.
package linkedhashmap

import (
	"fmt"
	"iter"
	"math"
	"strings"

	"github.com/lock14/collections"
)

var (
	_ collections.MutableSequencedMap[int, int] = (*LinkedHashMap[int, int])(nil)
	_ fmt.Stringer                              = (*LinkedHashMap[int, int])(nil)
)

// KeyOrder specifies the iteration order maintained by a LinkedHashMap.
type KeyOrder bool

const (
	// InsertionOrder maintains keys in the order they were first inserted.
	InsertionOrder KeyOrder = false
	// AccessOrder moves keys to the back whenever they are accessed or updated.
	AccessOrder KeyOrder = true
)

// config holds configuration values for a LinkedHashMap.
type config struct {
	keyOrder    KeyOrder
	maxElements int
	capacity    int
}

// Option configures a LinkedHashMap.
type Option func(*config)

// WithAccessOrder configures the LinkedHashMap to move entries to the back
// whenever they are accessed via Get or updated via Put.
func WithAccessOrder() Option {
	return func(c *config) {
		c.keyOrder = AccessOrder
	}
}

// WithInsertionOrder configures the LinkedHashMap to maintain entries in insertion order.
func WithInsertionOrder() Option {
	return func(c *config) {
		c.keyOrder = InsertionOrder
	}
}

// WithMaxElements configures the maximum number of entries the LinkedHashMap can hold
// before evicting the eldest entry on insertion.
// Panics in New if maxElements is less than 1.
func WithMaxElements(maxElements int) Option {
	return func(c *config) {
		c.maxElements = maxElements
	}
}

// WithCapacity configures the initial capacity hint of the underlying hash table.
// Panics in New if capacity is negative.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

// LinkedHashMap is a hash map backed by a doubly-linked list that preserves
// insertion or access order and supports optional bounded-size eviction.
type LinkedHashMap[K comparable, V any] struct {
	hashtable   map[K]*node[K, V]
	list        *node[K, V]
	accessOrder KeyOrder
	maxElements int
}

// New creates an empty LinkedHashMap configured with the given options.
// By default, it maintains insertion order and has unbounded size.
// Panics if configured with a negative capacity or maxElements less than 1.
func New[K comparable, V any](opts ...Option) *LinkedHashMap[K, V] {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.capacity < 0 {
		panic("linkedhashmap: capacity cannot be negative")
	}
	if cfg.maxElements <= 0 {
		panic("linkedhashmap: maxElements must be at least 1")
	}
	return &LinkedHashMap[K, V]{
		hashtable:   make(map[K]*node[K, V], cfg.capacity),
		list:        sentinel[K, V](),
		accessOrder: cfg.keyOrder,
		maxElements: cfg.maxElements,
	}
}

// Put associates the specified value with the specified key in the map.
// If the key already exists and AccessOrder is enabled, the entry is moved to the back.
// If the key is new and exceeds maxElements, the front entry is evicted.
func (m *LinkedHashMap[K, V]) Put(key K, value V) {
	n, ok := m.hashtable[key]
	if ok {
		n.value = value
		if m.accessOrder {
			unlink(n)
			insertBefore(m.list, n)
		}
		return
	}
	n = &node[K, V]{
		key:   key,
		value: value,
	}
	m.hashtable[key] = n
	insertBefore(m.list, n)
	if m.removeEldest() {
		eldest := m.list.next
		unlink(eldest)
		delete(m.hashtable, eldest.key)
	}
}

// PutFirst inserts the specified key-value pair at the front of the map.
// If the key already exists, its value is updated and it is moved to the front.
// If the key is new and exceeds maxElements, the back entry is evicted.
func (m *LinkedHashMap[K, V]) PutFirst(key K, value V) {
	n, ok := m.hashtable[key]
	if ok {
		n.value = value
		unlink(n)
		insertBefore(m.list.next, n)
		return
	}
	n = &node[K, V]{
		key:   key,
		value: value,
	}
	m.hashtable[key] = n
	insertBefore(m.list.next, n)
	if m.removeEldest() {
		eldest := m.list.prev
		unlink(eldest)
		delete(m.hashtable, eldest.key)
	}
}

// PutLast inserts the specified key-value pair at the end of the map.
// If the key already exists, its value is updated and it is moved to the back.
// If the key is new and exceeds maxElements, the front entry is evicted.
func (m *LinkedHashMap[K, V]) PutLast(key K, value V) {
	n, ok := m.hashtable[key]
	if ok {
		n.value = value
		unlink(n)
		insertBefore(m.list, n)
		return
	}
	m.Put(key, value)
}

// PollFirst removes and returns the first key-value pair in the map.
// Panics if the map is empty.
func (m *LinkedHashMap[K, V]) PollFirst() (K, V) {
	if m.Size() == 0 {
		panic("linkedhashmap: PollFirst called on empty map")
	}
	e := m.list.next
	unlink(e)
	delete(m.hashtable, e.key)
	return e.key, e.value
}

// PollLast removes and returns the last key-value pair in the map.
// Panics if the map is empty.
func (m *LinkedHashMap[K, V]) PollLast() (K, V) {
	if m.Size() == 0 {
		panic("linkedhashmap: PollLast called on empty map")
	}
	e := m.list.prev
	unlink(e)
	delete(m.hashtable, e.key)
	return e.key, e.value
}

// Remove removes the mapping for the specified key from the map if present.
func (m *LinkedHashMap[K, V]) Remove(key K) {
	n, ok := m.hashtable[key]
	if ok {
		unlink(n)
		delete(m.hashtable, key)
	}
}

// Clear removes all key-value pairs from the map.
func (m *LinkedHashMap[K, V]) Clear() {
	m.hashtable = make(map[K]*node[K, V])
	m.list = sentinel[K, V]()
}

// Get returns the value associated with the specified key, and a boolean indicating if it was found.
// If the key is found and AccessOrder is enabled, the entry is moved to the back.
func (m *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	n, ok := m.hashtable[key]
	if !ok {
		var zero V
		return zero, false
	}
	if m.accessOrder {
		unlink(n)
		insertBefore(m.list, n)
	}
	return n.value, true
}

// First returns the first key-value pair in the map without removing it.
// Panics if the map is empty.
func (m *LinkedHashMap[K, V]) First() (K, V) {
	if m.Size() == 0 {
		panic("linkedhashmap: First called on empty map")
	}
	e := m.list.next
	return e.key, e.value
}

// Last returns the last key-value pair in the map without removing it.
// Panics if the map is empty.
func (m *LinkedHashMap[K, V]) Last() (K, V) {
	if m.Size() == 0 {
		panic("linkedhashmap: Last called on empty map")
	}
	e := m.list.prev
	return e.key, e.value
}

// ContainsKey returns true if the map contains a mapping for the specified key.
func (m *LinkedHashMap[K, V]) ContainsKey(key K) bool {
	_, ok := m.hashtable[key]
	return ok
}

// Size returns the number of key-value pairs in the map.
func (m *LinkedHashMap[K, V]) Size() int {
	return len(m.hashtable)
}

// Empty returns true if the map contains no key-value pairs.
func (m *LinkedHashMap[K, V]) Empty() bool {
	return len(m.hashtable) == 0
}

// All returns an iterator over all key-value pairs in the map in encounter order.
func (m *LinkedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m.list == nil {
			return
		}
		for cur := m.list.next; cur != m.list; {
			next := cur.next
			if !yield(cur.key, cur.value) {
				return
			}
			cur = next
		}
	}
}

// Keys returns an iterator over all keys in the map in encounter order.
func (m *LinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if m.list == nil {
			return
		}
		for cur := m.list.next; cur != m.list; {
			next := cur.next
			if !yield(cur.key) {
				return
			}
			cur = next
		}
	}
}

// Values returns an iterator over all values in the map in encounter order.
func (m *LinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if m.list == nil {
			return
		}
		for cur := m.list.next; cur != m.list; {
			next := cur.next
			if !yield(cur.value) {
				return
			}
			cur = next
		}
	}
}

// Backward returns an iterator over the key-value pairs in reverse encounter order.
func (m *LinkedHashMap[K, V]) Backward() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m.list == nil {
			return
		}
		for cur := m.list.prev; cur != m.list; {
			prev := cur.prev
			if !yield(cur.key, cur.value) {
				return
			}
			cur = prev
		}
	}
}

// BackwardKeys returns an iterator over the keys in reverse encounter order.
func (m *LinkedHashMap[K, V]) BackwardKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if m.list == nil {
			return
		}
		for cur := m.list.prev; cur != m.list; {
			prev := cur.prev
			if !yield(cur.key) {
				return
			}
			cur = prev
		}
	}
}

// BackwardValues returns an iterator over the values in reverse encounter order.
func (m *LinkedHashMap[K, V]) BackwardValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		if m.list == nil {
			return
		}
		for cur := m.list.prev; cur != m.list; {
			prev := cur.prev
			if !yield(cur.value) {
				return
			}
			cur = prev
		}
	}
}

// String returns a string representation of the map in encounter order,
// matching Go's built-in map format.
func (m *LinkedHashMap[K, V]) String() string {
	if m.Empty() {
		return "map[]"
	}
	var sb strings.Builder
	sb.WriteString("map[")
	first := true
	for cur := m.list.next; cur != m.list; cur = cur.next {
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		_, _ = fmt.Fprintf(&sb, "%v:%v", cur.key, cur.value)
	}
	sb.WriteByte(']')
	return sb.String()
}

func (m *LinkedHashMap[K, V]) removeEldest() bool {
	return len(m.hashtable) > m.maxElements
}

func defaultConfig() *config {
	return &config{
		keyOrder:    InsertionOrder,
		maxElements: math.MaxInt,
	}
}

type node[K, V any] struct {
	key   K
	value V
	prev  *node[K, V]
	next  *node[K, V]
}

func sentinel[K, V any]() *node[K, V] {
	n := &node[K, V]{}
	n.prev = n
	n.next = n
	return n
}

func insertBefore[K, V any](n *node[K, V], b *node[K, V]) {
	b.prev = n.prev
	b.next = n
	n.prev.next = b
	n.prev = b
}

func unlink[K, V any](n *node[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
}
