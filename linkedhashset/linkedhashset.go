// Package linkedhashset provides a hash-backed set that preserves insertion or access order.
package linkedhashset

import (
	"fmt"
	"iter"
	"strings"

	"github.com/lock14/collections"
	"github.com/lock14/collections/linkedhashmap"
)

var (
	_ collections.MutableSequencedSet[int] = (*LinkedHashSet[int])(nil)
	_ fmt.Stringer                         = (*LinkedHashSet[int])(nil)
)

// Option configures a LinkedHashSet.
type Option = linkedhashmap.Option

// WithAccessOrder configures the LinkedHashSet to move elements to the back
// whenever they are re-inserted via Add.
func WithAccessOrder() Option {
	return linkedhashmap.WithAccessOrder()
}

// WithInsertionOrder configures the LinkedHashSet to maintain elements in insertion order.
func WithInsertionOrder() Option {
	return linkedhashmap.WithInsertionOrder()
}

// WithMaxElements configures the maximum number of elements the LinkedHashSet can hold
// before evicting the eldest element on insertion.
// Panics in New if maxElements is less than 1.
func WithMaxElements(maxElements int) Option {
	return linkedhashmap.WithMaxElements(maxElements)
}

// WithCapacity configures the initial capacity hint of the underlying hash table.
// Panics in New if capacity is negative.
func WithCapacity(capacity int) Option {
	return linkedhashmap.WithCapacity(capacity)
}

// LinkedHashSet represents a set of unique elements that preserves insertion or access order.
type LinkedHashSet[T comparable] struct {
	m *linkedhashmap.LinkedHashMap[T, struct{}]
}

// New creates an empty LinkedHashSet configured with the given options.
// By default, it maintains insertion order and has unbounded size.
// Panics if configured with a negative capacity or maxElements less than 1.
func New[T comparable](opts ...Option) *LinkedHashSet[T] {
	return &LinkedHashSet[T]{
		m: linkedhashmap.New[T, struct{}](opts...),
	}
}

// Add inserts the specified item into the set.
// If the item already exists and AccessOrder is enabled, it is moved to the back.
// If the item is new and exceeds maxElements, the front element is evicted.
func (s *LinkedHashSet[T]) Add(item T) {
	s.m.Put(item, struct{}{})
}

// AddFirst inserts the specified item at the front of the set.
// If the item already exists, it is moved to the front.
// If the item is new and exceeds maxElements, the back element is evicted.
func (s *LinkedHashSet[T]) AddFirst(item T) {
	s.m.PutFirst(item, struct{}{})
}

// AddLast inserts the specified item at the end of the set.
// If the item already exists, it is moved to the back.
// If the item is new and exceeds maxElements, the front element is evicted.
func (s *LinkedHashSet[T]) AddLast(item T) {
	s.m.PutLast(item, struct{}{})
}

// AddAll inserts all elements from the given sequence into the set.
func (s *LinkedHashSet[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		s.Add(t)
	}
}

// Remove removes and returns the first element from the set.
// Panics if the set is empty.
func (s *LinkedHashSet[T]) Remove() T {
	return s.PollFirst()
}

// PollFirst removes and returns the first element in the set.
// Panics if the set is empty.
func (s *LinkedHashSet[T]) PollFirst() T {
	k, _ := s.m.PollFirst()
	return k
}

// PollLast removes and returns the last element in the set.
// Panics if the set is empty.
func (s *LinkedHashSet[T]) PollLast() T {
	k, _ := s.m.PollLast()
	return k
}

// RemoveElement removes the specified item from the set if present.
func (s *LinkedHashSet[T]) RemoveElement(item T) {
	s.m.Remove(item)
}

// RemoveAll removes all elements of the specified collection from this set.
func (s *LinkedHashSet[T]) RemoveAll(other collections.Collection[T]) {
	if s.Empty() {
		return
	}
	if otherSet, ok := other.(*LinkedHashSet[T]); ok && otherSet == s {
		s.Clear()
		return
	}
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

// RetainAll retains only the elements in this set that are contained in the specified collection.
func (s *LinkedHashSet[T]) RetainAll(other collections.Collection[T]) {
	if s.Empty() {
		return
	}
	if otherSet, ok := other.(*LinkedHashSet[T]); ok && otherSet == s {
		return
	}
	if set, ok := other.(collections.Set[T]); ok {
		for k := range s.m.Keys() {
			if !set.Contains(k) {
				s.RemoveElement(k)
			}
		}
		return
	}
	keep := make(map[T]struct{})
	for t := range other.All() {
		if s.Contains(t) {
			keep[t] = struct{}{}
		}
	}
	for k := range s.m.Keys() {
		if _, ok := keep[k]; !ok {
			s.RemoveElement(k)
		}
	}
}

// Clear removes all elements from the set.
func (s *LinkedHashSet[T]) Clear() {
	s.m.Clear()
}

// First returns the first element in the set without removing it.
// Panics if the set is empty.
func (s *LinkedHashSet[T]) First() T {
	k, _ := s.m.First()
	return k
}

// Last returns the last element in the set without removing it.
// Panics if the set is empty.
func (s *LinkedHashSet[T]) Last() T {
	k, _ := s.m.Last()
	return k
}

// Contains returns true if the specified item is in the set.
func (s *LinkedHashSet[T]) Contains(item T) bool {
	return s.m.ContainsKey(item)
}

// ContainsAll returns true if all elements in the given collection are in the set.
func (s *LinkedHashSet[T]) ContainsAll(other collections.Collection[T]) bool {
	if otherSet, ok := other.(*LinkedHashSet[T]); ok && otherSet == s {
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
func (s *LinkedHashSet[T]) Size() int {
	return s.m.Size()
}

// Empty returns true if the set contains no elements.
func (s *LinkedHashSet[T]) Empty() bool {
	return s.m.Empty()
}

// All returns an iterator over the elements in the set in encounter order.
func (s *LinkedHashSet[T]) All() iter.Seq[T] {
	return s.m.Keys()
}

// Backward returns an iterator over the elements in reverse encounter order.
func (s *LinkedHashSet[T]) Backward() iter.Seq[T] {
	return s.m.BackwardKeys()
}

// String returns a string representation of the set in encounter order,
// matching Go's built-in slice format.
func (s *LinkedHashSet[T]) String() string {
	if s.Empty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for item := range s.All() {
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		_, _ = fmt.Fprintf(&sb, "%v", item)
	}
	sb.WriteByte(']')
	return sb.String()
}
