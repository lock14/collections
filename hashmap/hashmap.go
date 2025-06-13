package hashmap

import (
	"fmt"
	"github.com/lock14/collections"
	"hash/maphash"
	"iter"
	"unsafe"
)

const (
	DefaultCapacity   = 10
	DefaultLoadFactor = 0.75
)

var _ collections.MutableMap[int, int] = (*HashMap[int, int])(nil)

type Config struct {
	capacity   int
	loadFactor float64
}

type Opt func(*Config)

// public functions/receivers

func Capacity(capacity int) Opt {
	return func(config *Config) {
		config.capacity = capacity
	}
}

func LoadFactor(loadFactor float64) Opt {
	return func(config *Config) {
		config.loadFactor = loadFactor
	}
}

type HashMap[K comparable, V any] struct {
	hashtable     []*hashNode[K, V]
	size          int
	maxLoadFactor float64
	hash          maphash.Hash
}

type hashNode[K any, V any] struct {
	key   K
	value V
	next  *hashNode[K, V]
}

func New[K comparable, V any](opts ...Opt) *HashMap[K, V] {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return &HashMap[K, V]{
		hashtable:     make([]*hashNode[K, V], config.capacity),
		size:          0,
		maxLoadFactor: config.loadFactor,
	}
}

func (hm *HashMap[K, V]) ContainsKey(key K) bool {
	return hm.get(key) != nil
}

func (hm *HashMap[K, V]) Get(key K) (V, bool) {
	var v V
	var ok bool
	if n := hm.get(key); n != nil {
		v, ok = n.value, true
	}
	return v, ok
}

func (hm *HashMap[K, V]) Put(key K, value V) {
	if n := hm.get(key); n != nil {
		n.value = value
	} else {
		hm.insert(&hashNode[K, V]{key: key, value: value})
		hm.size++
	}
}

func (hm *HashMap[K, V]) Remove(key K) {
	var prev *hashNode[K, V]
	index := hm.fix(hm.hashBytes(key))
	node := hm.hashtable[index]
	for node != nil {
		if node.key == key {
			if prev == nil {
				hm.hashtable[index] = node.next
			} else {
				prev.next = node.next
			}
			node.next = nil
			hm.size--
			return
		}
		prev = node
		node = node.next
	}
}

func (hm *HashMap[K, V]) Size() int {
	return hm.size
}

func (hm *HashMap[K, V]) Empty() bool {
	return hm.Size() == 0
}

func (hm *HashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := 0; i < len(hm.hashtable); i++ {
			cur := hm.hashtable[i]
			for cur != nil {
				if !yield(cur.key, cur.value) {
					return
				}
				cur = cur.next
			}
		}
	}
}

func (hm *HashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range hm.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func (hm *HashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range hm.All() {
			if !yield(v) {
				return
			}
		}
	}
}

// private functions/receivers

func (hm *HashMap[K, V]) get(key K) *hashNode[K, V] {
	index := hm.fix(hm.hashBytes(key))
	node := hm.hashtable[index]
	for node != nil {
		if node.key == key {
			return node
		}
		node = node.next
	}
	return node
}

func (hm *HashMap[K, V]) insert(node *hashNode[K, V]) {
	if hm.loadFactor() > hm.maxLoadFactor {
		hm.resize()
	}
	index := hm.fix(hm.hashBytes(node.key))
	current := hm.hashtable[index]
	if current == nil {
		hm.hashtable[index] = node
	} else {
		for current.next != nil {
			current = current.next
		}
		current.next = node
	}
}

func (hm *HashMap[K, V]) fix(hash uint64) uint64 {
	return hash % uint64(cap(hm.hashtable))
}

func (hm *HashMap[K, V]) hashBytes(key K) uint64 {
	data := (*(*[1<<31 - 1]byte)(unsafe.Pointer(&key)))[:unsafe.Sizeof(key)]
	if _, err := hm.hash.Write(data); err != nil {
		panic(fmt.Sprintf("cannot hash key: %v", err))
	}
	hash := hm.hash.Sum64()
	hm.hash.Reset()
	return hash
}

func (hm *HashMap[K, V]) loadFactor() float64 {
	return float64(hm.size) / float64(cap(hm.hashtable))
}

func (hm *HashMap[K, V]) resize() {
	oldHashTable := hm.hashtable
	hm.hashtable = make([]*hashNode[K, V], 2*cap(oldHashTable))
	for _, node := range oldHashTable {
		for node != nil {
			next := node.next
			node.next = nil
			hm.insert(node)
			node = next
		}
	}
}

func defaultConfig() *Config {
	return &Config{
		capacity:   DefaultCapacity,
		loadFactor: DefaultLoadFactor,
	}
}
