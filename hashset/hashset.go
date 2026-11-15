// Package hashset provides a hash-backed set implementing MutableSet.
package hashset

import (
	"fmt"
	"iter"
	"maps"
	"strings"

	"github.com/lock14/collections"
)

var (
	_ collections.MutableSet[int] = (*HashSet[int])(nil)
	_ fmt.Stringer                = (*HashSet[int])(nil)
)

// config holds configuration values for a HashSet.
type config struct {
	capacity int
}

// Option configures a HashSet.
type Option func(*config)

// WithCapacity configures the initial capacity hint of the HashSet.
// Panics in New if capacity is negative.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

// HashSet represents an unordered set of unique elements backed by a hash table.
type HashSet[T comparable] struct {
	m map[T]struct{}
}

// New creates an empty HashSet configured with the given options.
// Panics if configured with a negative capacity.
func New[T comparable](opts ...Option) *HashSet[T] {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.capacity < 0 {
		panic("hashset: capacity cannot be negative")
	}
	return &HashSet[T]{
		m: make(map[T]struct{}, cfg.capacity),
	}
}

// Add inserts the specified item into the set.
func (s *HashSet[T]) Add(item T) {
	s.m[item] = struct{}{}
}

// AddAll inserts all elements from the given sequence into the set.
func (s *HashSet[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		s.Add(t)
	}
}

// Remove removes and returns an arbitrary element from the set.
// Panics if the set is empty.
func (s *HashSet[T]) Remove() T {
	for item := range s.m {
		delete(s.m, item)
		return item
	}
	panic("hashset: cannot remove from an empty set")
}

// RemoveElement removes the specified item from the set if present.
func (s *HashSet[T]) RemoveElement(item T) {
	delete(s.m, item)
}

// RemoveAll removes all elements of the specified collection from this set.
func (s *HashSet[T]) RemoveAll(other collections.Collection[T]) {
	if len(s.m) == 0 {
		return
	}
	if otherSet, ok := other.(*HashSet[T]); ok && otherSet == s {
		s.Clear()
		return
	}
	for t := range other.All() {
		delete(s.m, t)
	}
}

// RetainAll retains only the elements in this set that are contained in the specified collection.
func (s *HashSet[T]) RetainAll(other collections.Collection[T]) {
	if len(s.m) == 0 {
		return
	}
	if otherSet, ok := other.(*HashSet[T]); ok && otherSet == s {
		return
	}
	if set, ok := other.(collections.Set[T]); ok {
		for item := range s.m {
			if !set.Contains(item) {
				delete(s.m, item)
			}
		}
		return
	}
	keep := make(map[T]struct{})
	for t := range other.All() {
		if _, ok := s.m[t]; ok {
			keep[t] = struct{}{}
		}
	}
	s.m = keep
}

// Clear removes all elements from the set.
func (s *HashSet[T]) Clear() {
	s.m = make(map[T]struct{})
}

// Contains returns true if the set contains the specified item.
func (s *HashSet[T]) Contains(item T) bool {
	_, present := s.m[item]
	return present
}

// ContainsAll returns true if the set contains all elements of the specified collection.
func (s *HashSet[T]) ContainsAll(other collections.Collection[T]) bool {
	if otherSet, ok := other.(*HashSet[T]); ok && otherSet == s {
		return true
	}
	for item := range other.All() {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// Size returns the number of elements in the set.
func (s *HashSet[T]) Size() int {
	return len(s.m)
}

// Empty returns true if the set contains no elements.
func (s *HashSet[T]) Empty() bool {
	return len(s.m) == 0
}

// All returns an iterator over all elements in the set.
func (s *HashSet[T]) All() iter.Seq[T] {
	return maps.Keys(s.m)
}

// String returns a string representation of the set matching Go's built-in slice format.
func (s *HashSet[T]) String() string {
	if len(s.m) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for item := range s.m {
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		_, _ = fmt.Fprintf(&sb, "%v", item)
	}
	sb.WriteByte(']')
	return sb.String()
}
