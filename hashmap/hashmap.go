package hashmap

import (
	"fmt"
	"github.com/lock14/collections/iterator"
	"github.com/lock14/collections/pair"
	"hash/maphash"
	"unsafe"
)

const (
	DefaultCapacity   = 10
	DefaultLoadFactor = 0.75
)

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
		hm.insert(key, value)
	}
}

func (hm *HashMap[K, V]) Remove(key K) {
	index := hm.fix(hm.hashBytes(key))
	node := hm.hashtable[index]
	var prev *hashNode[K, V]
	for node != nil {
		if node.key == key {
			if prev == nil {
				hm.hashtable[index] = nil
			} else {
				prev.next = node.next
			}
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

func (hm *HashMap[K, V]) Entries() iterator.Iterator[*pair.Pair[K, V]] {
	ei := &entryIterator[K, V]{
		hashMap: hm,
		index:   0,
	}
	ei.current = ei.getNext(nil)
	return ei
}

func (hm *HashMap[K, V]) Keys() iterator.Iterator[K] {
	return &keyIterator[K, V]{
		ei: hm.Entries(),
	}
}

func (hm *HashMap[K, V]) Values() iterator.Iterator[V] {
	return &valueIterator[K, V]{
		ei: hm.Entries(),
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

func (hm *HashMap[K, V]) insert(key K, value V) {
	if hm.loadFactor() > hm.maxLoadFactor {
		hm.resize()
	}
	index := hm.fix(hm.hashBytes(key))
	current := hm.hashtable[index]
	if current == nil {
		hm.hashtable[index] = &hashNode[K, V]{
			key:   key,
			value: value,
		}
	} else {
		for current.next != nil {
			current = current.next
		}
		current.next = &hashNode[K, V]{
			key:   key,
			value: value,
		}
	}
	hm.size++
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
			hm.Put(node.key, node.value)
			node = node.next
		}
	}
}

func defaultConfig() *Config {
	return &Config{
		capacity:   DefaultCapacity,
		loadFactor: DefaultLoadFactor,
	}
}

// iterator stuff

type entryIterator[K comparable, V any] struct {
	hashMap *HashMap[K, V]
	index   int
	current *hashNode[K, V]
}

func (ei *entryIterator[K, V]) Empty() bool {
	return ei.current == nil
}

func (ei *entryIterator[K, V]) Next() (*pair.Pair[K, V], error) {
	if ei.Empty() {
		return nil, fmt.Errorf("iterator is empty")
	}
	cur := ei.current
	ei.current = ei.getNext(ei.current.next)
	return pair.New(cur.key, cur.value), nil
}

func (ei *entryIterator[K, V]) getNext(next *hashNode[K, V]) *hashNode[K, V] {
	if next == nil {
		for ei.index < cap(ei.hashMap.hashtable) && ei.hashMap.hashtable[ei.index] == nil {
			ei.index++
		}
		if ei.index < cap(ei.hashMap.hashtable) {
			next = ei.hashMap.hashtable[ei.index]
			ei.index++
		}
	}
	return next
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
