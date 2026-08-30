// Package hashmap provides a wrapper for the built-in map implementing MutableMap.
package hashmap

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"maps"
	"strings"
)

var (
	_ collections.MutableMap[int, int] = (*HashMap[int, int])(nil)
	_ fmt.Stringer                     = (*HashMap[int, int])(nil)
)

// HashMap is a wrapper around the built-in map that implements collections.MutableMap.
type HashMap[K comparable, V any] struct {
	m map[K]V
}

// config holds the values for configuring a HashMap.
type config struct {
	capacity int
}

// Option configures a HashMap config
type Option func(*config)

// WithCapacity configures the initial capacity of the HashMap.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

// New creates an empty HashMap.
func New[K comparable, V any](opts ...Option) *HashMap[K, V] {
	config := &config{}
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

// Get returns the value associated with the specified key, and a boolean indicating if it was found.
func (hm *HashMap[K, V]) Get(key K) (V, bool) {
	v, ok := hm.m[key]
	return v, ok
}

// Put associates the specified value with the specified key in the map.
func (hm *HashMap[K, V]) Put(key K, value V) {
	hm.m[key] = value
}

// Remove removes the mapping for the specified key from the map if present.
func (hm *HashMap[K, V]) Remove(key K) {
	delete(hm.m, key)
}

// Size returns the number of key-value pairs in the map.
func (hm *HashMap[K, V]) Size() int {
	return len(hm.m)
}

// Empty returns true if the map contains no key-value pairs.
func (hm *HashMap[K, V]) Empty() bool {
	return len(hm.m) == 0
}

// Clear removes all key-value pairs from the map.
func (hm *HashMap[K, V]) Clear() {
	hm.m = make(map[K]V)
}

// ContainsKey returns true if the map contains a mapping for the specified key.
func (hm *HashMap[K, V]) ContainsKey(key K) bool {
	_, ok := hm.m[key]
	return ok
}

// All returns an iterator over all key-value pairs in the map.
func (hm *HashMap[K, V]) All() iter.Seq2[K, V] {
	return maps.All(hm.m)
}

// Keys returns an iterator over all keys in the map.
func (hm *HashMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(hm.m)
}

// Values returns an iterator over all values in the map.
func (hm *HashMap[K, V]) Values() iter.Seq[V] {
	return maps.Values(hm.m)
}

// String returns a string representation of the map matching Go built-in map formatting.
func (hm *HashMap[K, V]) String() string {
	vals := make([]string, 0, hm.Size())
	for k, v := range hm.All() {
		vals = append(vals, fmt.Sprintf("%v:%v", k, v))
	}
	return "map[" + strings.Join(vals, " ") + "]"
}
