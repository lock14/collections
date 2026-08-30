// Package hashset provides a hash-backed set implementing MutableSet.
package hashset

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"maps"
	"strings"
)

var (
	_ collections.MutableSet[int] = (*HashSet[int])(nil)
	_ fmt.Stringer                = (*HashSet[int])(nil)
)

// HashSet represents a set of elements of type T.
type HashSet[T comparable] struct {
	m map[T]struct{}
}

// config holds the values for configuring a HashSet.
type config struct {
	capacity int
}

// Option configures a HashSet config
type Option func(*config)

// WithCapacity configures the initial capacity of the HashSet.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

// New creates an empty HashSet.
func New[T comparable](opts ...Option) *HashSet[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &HashSet[T]{
		m: make(map[T]struct{}, config.capacity),
	}
}

// Add inserts the specified item into the set.
func (s *HashSet[T]) Add(item T) {
	s.m[item] = struct{}{}
}

// Remove removes and returns an arbitrary element from the set.
// Panics if the set is empty.
func (s *HashSet[T]) Remove() T {
	for item := range s.m {
		delete(s.m, item)
		return item
	}
	panic("cannot remove from an empty set")
}

// RemoveElement removes the specified item from the set.
func (s *HashSet[T]) RemoveElement(item T) {
	delete(s.m, item)
}

// Contains returns true if the set contains the specified item.
func (s *HashSet[T]) Contains(item T) bool {
	_, present := s.m[item]
	return present
}

// ContainsAll returns true if the set contains all elements of the specified collection.
func (s *HashSet[T]) ContainsAll(other collections.Collection[T]) bool {
	for item := range other.All() {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// AddAll inserts all elements from the given sequence into the set.
func (s *HashSet[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		s.Add(t)
	}
}

// RemoveAll removes all elements of the specified collection from this set.
func (s *HashSet[T]) RemoveAll(other collections.Collection[T]) {
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

// RetainAll retains only the elements in this set that are contained in the specified collection.
func (s *HashSet[T]) RetainAll(other collections.Collection[T]) {
	newMap := make(map[T]struct{})
	for t := range other.All() {
		if _, ok := s.m[t]; ok {
			newMap[t] = struct{}{}
		}
	}
	s.m = newMap
}

// Clear removes all elements from the set.
func (s *HashSet[T]) Clear() {
	s.m = make(map[T]struct{})
}

// Size returns the number of elements in the set.
func (s *HashSet[T]) Size() int {
	return len(s.m)
}

// Empty returns true if the set contains no elements.
func (s *HashSet[T]) Empty() bool {
	return len(s.m) == 0
}

// String returns a string representation of the set matching Go slice formatting.
func (s *HashSet[T]) String() string {
	vals := make([]string, 0, len(s.m))
	for item := range s.m {
		vals = append(vals, fmt.Sprintf("%v", item))
	}
	return "[" + strings.Join(vals, " ") + "]"
}

// All returns an iterator over all elements in the set.
func (s *HashSet[T]) All() iter.Seq[T] {
	return maps.Keys(s.m)
}

func defaultConfig() *config {
	return &config{}
}
