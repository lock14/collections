package hashmap

import (
	"github.com/lock14/collections"
	"iter"
	"maps"
)

var _ collections.MutableMap[int, int] = (*HashMap[int, int])(nil)

// HashMap is a wrapper around the built-in map that implements collections.MutableMap.
type HashMap[K comparable, V any] struct {
	m map[K]V
}

// Config holds the values for configuring a HashMap.
type Config struct {
	capacity int
}

// Option configures a HashMap config
type Option func(*Config)

// WithCapacity configures the initial capacity of the HashMap.
func WithCapacity(capacity int) Option {
	return func(c *Config) {
		c.capacity = capacity
	}
}

// New creates an empty HashMap.
func New[K comparable, V any](opts ...Option) *HashMap[K, V] {
	config := &Config{}
	for _, option := range opts {
		option(config)
	}
	return &HashMap[K, V]{
		m: make(map[K]V, config.capacity),
	}
}

// Wrap wraps an existing built-in map.
func Wrap[K comparable, V any](m map[K]V) *HashMap[K, V] {
	return &HashMap[K, V]{
		m: m,
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
	return len(hm.m) == 0
}

func (hm *HashMap[K, V]) Clear() {
	hm.m = make(map[K]V)
}

func (hm *HashMap[K, V]) ContainsKey(key K) bool {
	_, ok := hm.m[key]
	return ok
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
