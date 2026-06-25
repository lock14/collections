// Package linkedhashmap provides a hash map that preserves insertion order.
package linkedhashmap

import (
	"github.com/lock14/collections"
	"iter"
	"math"
)

const (
	InsertionOrder = false
	AccessOrder    = true
)

var _ collections.MutableSequencedMap[int, int] = (*LinkedHashMap[int, int])(nil)

// KeyOrder represents the iteration order of the linked hash map.
type KeyOrder bool

// Config holds the configuration for a LinkedHashMap.
type Config struct {
	keyOrder    KeyOrder
	maxElements int
	capacity    int
}

// Opt is a function that configures a LinkedHashMap.
type Opt func(*Config)

// WithAccessOrder configures the LinkedHashMap to iterate over elements in the order they were last accessed.
func WithAccessOrder() Opt {
	return func(config *Config) {
		config.keyOrder = AccessOrder
	}
}

// WithInsertionOrder configures the LinkedHashMap to iterate over elements in the order they were inserted.
func WithInsertionOrder() Opt {
	return func(config *Config) {
		config.keyOrder = InsertionOrder
	}
}

// WithMaxElements configures the maximum number of elements the LinkedHashMap can hold before evicting.
func WithMaxElements(max int) Opt {
	return func(config *Config) {
		config.maxElements = max
	}
}

// WithCapacity configures the initial capacity of the underlying map.
func WithCapacity(capacity int) Opt {
	return func(config *Config) {
		config.capacity = capacity
	}
}

// public functions/receivers

// LinkedHashMap is a hash map that preserves insertion or access order.
type LinkedHashMap[K comparable, V any] struct {
	hashtable   map[K]*node[K, V]
	list        *node[K, V]
	accessOrder KeyOrder
	maxElements int
}

// New creates a new LinkedHashMap with the given options.
func New[K comparable, V any](opts ...Opt) *LinkedHashMap[K, V] {
	c := defaultConfig()
	for _, opt := range opts {
		opt(c)
	}
	return &LinkedHashMap[K, V]{
		hashtable:   make(map[K]*node[K, V], c.capacity),
		list:        sentinel[K, V](),
		accessOrder: c.keyOrder,
		maxElements: c.maxElements,
	}
}

func (hm *LinkedHashMap[K, V]) Put(key K, value V) {
	n, ok := hm.hashtable[key]
	if ok {
		n.value = value
		if hm.accessOrder {
			unlink(n)
			// make n the tail of the list
			insertBefore(hm.list, n)
		}
	} else {
		n = &node[K, V]{
			key:   key,
			value: value,
		}
		hm.hashtable[key] = n
		// make n the tail of the list
		insertBefore(hm.list, n)
		if hm.removeEldest() {
			eldest := hm.list.next
			unlink(eldest)
			delete(hm.hashtable, eldest.key)
		}
	}
}

func (hm *LinkedHashMap[K, V]) PutFirst(key K, value V) {
	n, ok := hm.hashtable[key]
	if ok {
		n.value = value
		unlink(n)
		insertBefore(hm.list.next, n)
	} else {
		n = &node[K, V]{
			key:   key,
			value: value,
		}
		hm.hashtable[key] = n
		insertBefore(hm.list.next, n)
		if hm.removeEldest() {
			eldest := hm.list.prev
			unlink(eldest)
			delete(hm.hashtable, eldest.key)
		}
	}
}

func (hm *LinkedHashMap[K, V]) PutLast(key K, value V) {
	n, ok := hm.hashtable[key]
	if ok {
		n.value = value
		unlink(n)
		insertBefore(hm.list, n)
	} else {
		hm.Put(key, value)
	}
}

func (hm *LinkedHashMap[K, V]) First() (K, V, bool) {
	if hm.Empty() {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return hm.list.next.key, hm.list.next.value, true
}

func (hm *LinkedHashMap[K, V]) Last() (K, V, bool) {
	if hm.Empty() {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return hm.list.prev.key, hm.list.prev.value, true
}

func (hm *LinkedHashMap[K, V]) PollFirst() (K, V, bool) {
	if hm.Empty() {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	n := hm.list.next
	unlink(n)
	delete(hm.hashtable, n.key)
	return n.key, n.value, true
}

func (hm *LinkedHashMap[K, V]) PollLast() (K, V, bool) {
	if hm.Empty() {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	n := hm.list.prev
	unlink(n)
	delete(hm.hashtable, n.key)
	return n.key, n.value, true
}

func (hm *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	n, ok := hm.hashtable[key]
	if !ok {
		var zero V
		return zero, false
	}
	if hm.accessOrder {
		unlink(n)
		// make n the tail of the list
		insertBefore(hm.list, n)
	}
	return n.value, true
}

func (hm *LinkedHashMap[K, V]) Remove(key K) {
	n, ok := hm.hashtable[key]
	if ok {
		unlink(n)
		delete(hm.hashtable, key)
	}
}

func (hm *LinkedHashMap[K, V]) ContainsKey(key K) bool {
	_, ok := hm.hashtable[key]
	return ok
}

func (hm *LinkedHashMap[K, V]) Size() int {
	return len(hm.hashtable)
}

func (hm *LinkedHashMap[K, V]) Empty() bool {
	return hm.Size() == 0
}

func (hm *LinkedHashMap[K, V]) Clear() {
	hm.hashtable = make(map[K]*node[K, V])
	hm.list = sentinel[K, V]()
}

func (hm *LinkedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for cur := hm.list.next; cur != hm.list && yield(cur.key, cur.value); {
			cur = cur.next
		}
	}
}

func (hm *LinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for cur := hm.list.next; cur != hm.list; cur = cur.next {
			if !yield(cur.key) {
				return
			}
		}
	}
}

func (hm *LinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for cur := hm.list.next; cur != hm.list; cur = cur.next {
			if !yield(cur.value) {
				return
			}
		}
	}
}

func (hm *LinkedHashMap[K, V]) ReversedAll() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for cur := hm.list.prev; cur != hm.list && yield(cur.key, cur.value); {
			cur = cur.prev
		}
	}
}

func (hm *LinkedHashMap[K, V]) ReversedKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for cur := hm.list.prev; cur != hm.list; cur = cur.prev {
			if !yield(cur.key) {
				return
			}
		}
	}
}

func (hm *LinkedHashMap[K, V]) ReversedValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		for cur := hm.list.prev; cur != hm.list; cur = cur.prev {
			if !yield(cur.value) {
				return
			}
		}
	}
}

func (hm *LinkedHashMap[K, V]) removeEldest() bool {
	return hm.Size() > hm.maxElements
}

func defaultConfig() *Config {
	return &Config{
		keyOrder:    false,
		maxElements: math.MaxInt,
	}
}

// linked list stuff

type node[K, V any] struct {
	key   K
	value V
	prev  *node[K, V]
	next  *node[K, V]
}

func sentinel[K, V any]() *node[K, V] {
	node := &node[K, V]{}
	node.prev = node
	node.next = node
	return node
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
