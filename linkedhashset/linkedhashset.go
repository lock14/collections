// Package linkedhashset provides a hash-backed set that preserves insertion or access order.
package linkedhashset

import (
	"fmt"
	"github.com/lock14/collections"
	"github.com/lock14/collections/linkedhashmap"
	"iter"
	"strings"
)

var _ collections.MutableSet[int] = (*LinkedHashSet[int])(nil)

// LinkedHashSet represents a set of elements of type T that preserves iteration order.
type LinkedHashSet[T comparable] struct {
	m *linkedhashmap.LinkedHashMap[T, struct{}]
}

// Option configures a LinkedHashSet.
type Option = linkedhashmap.Opt

// WithAccessOrder configures the LinkedHashSet to iterate over elements in the order they were last accessed.
func WithAccessOrder() Option {
	return linkedhashmap.WithAccessOrder()
}

// WithInsertionOrder configures the LinkedHashSet to iterate over elements in the order they were inserted.
func WithInsertionOrder() Option {
	return linkedhashmap.WithInsertionOrder()
}

// WithMaxElements configures the maximum number of elements the LinkedHashSet can hold before evicting.
func WithMaxElements(max int) Option {
	return linkedhashmap.WithMaxElements(max)
}

// WithCapacity configures the initial capacity of the underlying set.
func WithCapacity(capacity int) Option {
	return linkedhashmap.WithCapacity(capacity)
}

// New creates an empty LinkedHashSet.
func New[T comparable](opts ...Option) *LinkedHashSet[T] {
	return &LinkedHashSet[T]{
		m: linkedhashmap.New[T, struct{}](opts...),
	}
}

// Add adds the specified item to the set.
func (s *LinkedHashSet[T]) Add(item T) {
	s.m.Put(item, struct{}{})
}

// Remove removes and returns an arbitrary element from the set.
func (s *LinkedHashSet[T]) Remove() T {
	for k := range s.m.Keys() {
		s.m.Remove(k)
		return k
	}
	panic("cannot remove from an empty set")
}

// RemoveElement removes the specified item from the set.
func (s *LinkedHashSet[T]) RemoveElement(item T) {
	s.m.Remove(item)
}

// Contains returns true if the specified item is in the set.
func (s *LinkedHashSet[T]) Contains(item T) bool {
	return s.m.ContainsKey(item)
}

// ContainsAll returns true if all elements in the given collection are in the set.
func (s *LinkedHashSet[T]) ContainsAll(other collections.Collection[T]) bool {
	for item := range other.All() {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// AddAll adds all elements from the given sequence to the set.
func (s *LinkedHashSet[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		s.Add(t)
	}
}

// RemoveAll removes all elements in the given collection from the set.
func (s *LinkedHashSet[T]) RemoveAll(other collections.Collection[T]) {
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

// RetainAll removes all elements from the set that are not in the given collection.
func (s *LinkedHashSet[T]) RetainAll(other collections.Collection[T]) {
	// Need to check which elements to keep.
	// We iterate over the current set and remove those not in `other`.
	// Since we are iterating and removing, we should be careful or just build a new set.
	// Alternatively, we can collect the keys to remove.
	var toRemove []T
	if set, ok := other.(collections.Set[T]); ok {
		for k := range s.m.Keys() {
			if !set.Contains(k) {
				toRemove = append(toRemove, k)
			}
		}
	} else {
		// Generic collection: convert to a temporary set for fast lookups.
		tempSet := New[T]()
		tempSet.AddAll(other.All())
		for k := range s.m.Keys() {
			if !tempSet.Contains(k) {
				toRemove = append(toRemove, k)
			}
		}
	}
	for _, k := range toRemove {
		s.RemoveElement(k)
	}
}

// Clear removes all elements from the set.
func (s *LinkedHashSet[T]) Clear() {
	s.m.Clear()
}

// Size returns the number of elements in the set.
func (s *LinkedHashSet[T]) Size() int {
	return s.m.Size()
}

// Empty returns true if the set contains no elements.
func (s *LinkedHashSet[T]) Empty() bool {
	return s.m.Empty()
}

// All returns an iterator over the elements in the set.
func (s *LinkedHashSet[T]) All() iter.Seq[T] {
	return s.m.Keys()
}

// String returns a string representation of the set.
func (s *LinkedHashSet[T]) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	i := 0
	for item := range s.All() {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("%v", item))
		i++
	}
	sb.WriteString("]")
	return sb.String()
}
