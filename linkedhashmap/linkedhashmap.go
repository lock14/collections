package linkedhashmap

import (
	"fmt"
	"github.com/lock14/collections/iterator"
	"github.com/lock14/collections/pair"
	"math"
)

type KeyOrder bool

type Config struct {
	keyOrder    KeyOrder
	maxElements int
}

func WithAccessOrder() KeyOrder {
	return true
}

func WithInsertionOrder() KeyOrder {
	return false
}

func WithMaxElements(max int) int {
	return max
}

type Opt func(*Config)

// public functions/receivers

type LinkedHashMap[K comparable, V any] struct {
	hashtable   map[K]*node[K, V]
	list        *node[K, V]
	accessOrder KeyOrder
	maxElements int
}

func New[K comparable, V any](opts ...Opt) *LinkedHashMap[K, V] {
	c := defaultConfig()
	for _, opt := range opts {
		opt(c)
	}
	return &LinkedHashMap[K, V]{
		hashtable:   make(map[K]*node[K, V]),
		list:        sentinel[K, V](),
		accessOrder: c.keyOrder,
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

func (hm *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	n, ok := hm.hashtable[key]
	if ok && bool(hm.accessOrder) {
		unlink(n)
		// make n the tail of the list
		insertBefore(hm.list, n)
	}
	return n.value, ok
}

func (hm *LinkedHashMap[K, V]) Remove(key K) {
	n, ok := hm.hashtable[key]
	if ok {
		unlink(n)
		delete(hm.hashtable, key)
	}
}

func (hm *LinkedHashMap[K, V]) Size() int {
	return len(hm.hashtable)
}

func (hm *LinkedHashMap[K, V]) Empty() bool {
	return hm.Size() == 0
}

func (hm *LinkedHashMap[K, V]) Entries() iterator.Iterator[*pair.Pair[K, V]] {
	return &entryIterator[K, V]{
		hashMap: hm,
		current: hm.list.next,
	}
}

func (hm *LinkedHashMap[K, V]) Keys() iterator.Iterator[K] {
	return &keyIterator[K, V]{
		ei: hm.Entries(),
	}
}

func (hm *LinkedHashMap[K, V]) Values() iterator.Iterator[V] {
	return &valueIterator[K, V]{
		ei: hm.Entries(),
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

// iterator stuff

type entryIterator[K comparable, V any] struct {
	hashMap *LinkedHashMap[K, V]
	current *node[K, V]
}

func (ei *entryIterator[K, V]) Empty() bool {
	return ei.current == ei.hashMap.list
}

func (ei *entryIterator[K, V]) Next() (*pair.Pair[K, V], error) {
	if ei.Empty() {
		return nil, fmt.Errorf("iterator is empty")
	}
	cur := ei.current
	ei.current = ei.current.next
	return pair.New(cur.key, cur.value), nil
}

type keyIterator[K any, V any] struct {
	ei iterator.Iterator[*pair.Pair[K, V]]
}

func (ki *keyIterator[K, V]) Empty() bool {
	return ki.ei.Empty()
}

func (ki *keyIterator[K, V]) Next() (K, error) {
	p, err := ki.ei.Next()
	if err != nil {
		var k K
		return k, err
	}
	return p.Fst(), err
}

type valueIterator[K any, V any] struct {
	ei iterator.Iterator[*pair.Pair[K, V]]
}

func (vi *valueIterator[K, V]) Empty() bool {
	return vi.ei.Empty()
}

func (vi *valueIterator[K, V]) Next() (V, error) {
	p, err := vi.ei.Next()
	if err != nil {
		var v V
		return v, err
	}
	return p.Snd(), err
}
