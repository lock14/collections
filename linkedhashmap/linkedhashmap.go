// Package linkedhashmap provides a hash map that preserves insertion order.
package linkedhashmap

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"math"
	"strings"
)

const (
	InsertionOrder = false
	AccessOrder    = true
)

var (
	_ collections.MutableSequencedMap[int, int] = (*LinkedHashMap[int, int])(nil)
	_ fmt.Stringer                              = (*LinkedHashMap[int, int])(nil)
)

// KeyOrder represents the iteration order of the linked hash map.
type KeyOrder bool

// config holds the configuration for a LinkedHashMap.
type config struct {
	keyOrder    KeyOrder
	maxElements int
	capacity    int
}

// Opt is a function that configures a LinkedHashMap.
type Opt func(*config)

// WithAccessOrder configures the LinkedHashMap to iterate over elements in the order they were last accessed.
func WithAccessOrder() Opt {
	return func(config *config) {
		config.keyOrder = AccessOrder
	}
}

// WithInsertionOrder configures the LinkedHashMap to iterate over elements in the order they were inserted.
func WithInsertionOrder() Opt {
	return func(config *config) {
		config.keyOrder = InsertionOrder
	}
}

// WithMaxElements configures the maximum number of elements the LinkedHashMap can hold before evicting.
func WithMaxElements(max int) Opt {
	return func(config *config) {
		config.maxElements = max
	}
}

// WithCapacity configures the initial capacity of the underlying map.
func WithCapacity(capacity int) Opt {
	return func(config *config) {
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

// First returns the first key-value pair in the map.
func (hm *LinkedHashMap[K, V]) First() (K, V) {
	if hm.Size() == 0 {
		panic("First called on empty map")
	}
	e := hm.list.next
	return e.key, e.value
}

// Last returns the last key-value pair in the map.
func (hm *LinkedHashMap[K, V]) Last() (K, V) {
	if hm.Size() == 0 {
		panic("Last called on empty map")
	}
	e := hm.list.prev
	return e.key, e.value
}

// PollFirst removes and returns the first key-value pair in the map.
func (hm *LinkedHashMap[K, V]) PollFirst() (K, V) {
	if hm.Size() == 0 {
		panic("PollFirst called on empty map")
	}
	e := hm.list.next
	unlink(e)
	delete(hm.hashtable, e.key)
	return e.key, e.value
}

// PollLast removes and returns the last key-value pair in the map.
func (hm *LinkedHashMap[K, V]) PollLast() (K, V) {
	if hm.Size() == 0 {
		panic("PollLast called on empty map")
	}
	e := hm.list.prev
	unlink(e)
	delete(hm.hashtable, e.key)
	return e.key, e.value
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
		if hm.list == nil {
			return
		}
		for cur := hm.list.next; cur != hm.list; cur = cur.next {
			if !yield(cur.key, cur.value) {
				return
			}
		}
	}
}

func (hm *LinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if hm.list == nil {
			return
		}
		for cur := hm.list.next; cur != hm.list; cur = cur.next {
			if !yield(cur.key) {
				return
			}
		}
	}
}

func (hm *LinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if hm.list == nil {
			return
		}
		for cur := hm.list.next; cur != hm.list; cur = cur.next {
			if !yield(cur.value) {
				return
			}
		}
	}
}

func (hm *LinkedHashMap[K, V]) Backward() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if hm.list == nil {
			return
		}
		for cur := hm.list.prev; cur != hm.list; cur = cur.prev {
			if !yield(cur.key, cur.value) {
				return
			}
		}
	}
}

func (hm *LinkedHashMap[K, V]) BackwardKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if hm.list == nil {
			return
		}
		for cur := hm.list.prev; cur != hm.list; cur = cur.prev {
			if !yield(cur.key) {
				return
			}
		}
	}
}

func (hm *LinkedHashMap[K, V]) BackwardValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		if hm.list == nil {
			return
		}
		for cur := hm.list.prev; cur != hm.list; cur = cur.prev {
			if !yield(cur.value) {
				return
			}
		}
	}
}

// String returns a string representation of the map matching Go built-in map formatting.
func (hm *LinkedHashMap[K, V]) String() string {
	vals := make([]string, 0, hm.Size())
	for k, v := range hm.All() {
		vals = append(vals, fmt.Sprintf("%v:%v", k, v))
	}
	return "map[" + strings.Join(vals, " ") + "]"
}

func (hm *LinkedHashMap[K, V]) removeEldest() bool {
	return hm.Size() > hm.maxElements
}

func defaultConfig() *config {
	return &config{
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

// benchmark matrix test trigger
