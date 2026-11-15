// Package hashmap provides a wrapper around Go's built-in map implementing MutableMap.
package hashmap

import (
	"fmt"
	"iter"
	"maps"
	"strings"

	"github.com/lock14/collections"
)

var (
	_ collections.MutableMap[int, int] = (*HashMap[int, int])(nil)
	_ fmt.Stringer                     = (*HashMap[int, int])(nil)
)

// config holds configuration values for a HashMap.
type config struct {
	capacity int
}

// Option configures a HashMap.
type Option func(*config)

// WithCapacity configures the initial capacity hint of the HashMap.
// Panics in New if capacity is negative.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

// HashMap is a wrapper around Go's built-in map that implements collections.MutableMap.
type HashMap[K comparable, V any] struct {
	m map[K]V
}

// New creates an empty HashMap configured with the given options.
// Panics if configured with a negative capacity.
func New[K comparable, V any](opts ...Option) *HashMap[K, V] {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.capacity < 0 {
		panic("hashmap: capacity cannot be negative")
	}
	return &HashMap[K, V]{
		m: make(map[K]V, cfg.capacity),
	}
}

// Wrap wraps an existing built-in map without copying its entries.
// Mutations through the returned HashMap directly modify raw, and vice versa.
// If raw is nil, a new empty map is allocated so that subsequent writes do not panic.
func Wrap[K comparable, V any](raw map[K]V) *HashMap[K, V] {
	if raw == nil {
		raw = make(map[K]V)
	}
	return &HashMap[K, V]{
		m: raw,
	}
}

// Put associates the specified value with the specified key in the map.
func (m *HashMap[K, V]) Put(key K, value V) {
	m.m[key] = value
}

// Remove removes the mapping for the specified key from the map if present.
func (m *HashMap[K, V]) Remove(key K) {
	delete(m.m, key)
}

// Clear removes all key-value pairs from the map.
func (m *HashMap[K, V]) Clear() {
	m.m = make(map[K]V)
}

// Get returns the value associated with the specified key, and a boolean indicating if it was found.
func (m *HashMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.m[key]
	return v, ok
}

// ContainsKey returns true if the map contains a mapping for the specified key.
func (m *HashMap[K, V]) ContainsKey(key K) bool {
	_, ok := m.m[key]
	return ok
}

// Size returns the number of key-value pairs in the map.
func (m *HashMap[K, V]) Size() int {
	return len(m.m)
}

// Empty returns true if the map contains no key-value pairs.
func (m *HashMap[K, V]) Empty() bool {
	return len(m.m) == 0
}

// All returns an iterator over all key-value pairs in the map.
func (m *HashMap[K, V]) All() iter.Seq2[K, V] {
	return maps.All(m.m)
}

// Keys returns an iterator over all keys in the map.
func (m *HashMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.m)
}

// Values returns an iterator over all values in the map.
func (m *HashMap[K, V]) Values() iter.Seq[V] {
	return maps.Values(m.m)
}

// String returns a string representation of the map matching Go's built-in map format.
func (m *HashMap[K, V]) String() string {
	if len(m.m) == 0 {
		return "map[]"
	}
	var sb strings.Builder
	sb.WriteString("map[")
	first := true
	for k, v := range m.m {
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		_, _ = fmt.Fprintf(&sb, "%v:%v", k, v)
	}
	sb.WriteByte(']')
	return sb.String()
}
