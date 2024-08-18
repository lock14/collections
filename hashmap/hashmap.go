package hashmap

import (
	"iter"
	"maps"
)

const (
	DefaultCapacity = 10
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

type HashMap[K comparable, V any] struct {
	m map[K]V
}

func New[K comparable, V any](opts ...Opt) *HashMap[K, V] {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return &HashMap[K, V]{
		m: make(map[K]V, config.capacity),
	}
}

func (hm *HashMap[K, V]) Get(key K) (V, bool) {
	v, ok := hm.m[key]
	return v, ok
}

func (hm *HashMap[K, V]) Put(key K, value V) {
	hm.m[key] = value
}

func (hm *HashMap[K, V]) Remove(key K) {
	delete(hm.m, key)
}

func (hm *HashMap[K, V]) Size() int {
	return len(hm.m)
}

func (hm *HashMap[K, V]) Empty() bool {
	return hm.Size() == 0
}

func (hm *HashMap[K, V]) All() iter.Seq2[K, V] {
	return maps.All(hm.m)
}

func (hm *HashMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(hm.m)
}

func (hm *HashMap[K, V]) Values() iter.Seq[V] {
	return maps.Values(hm.m)
}

func defaultConfig() *Config {
	return &Config{
		capacity: DefaultCapacity,
	}
}
