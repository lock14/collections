// Package treeset provides a B-Tree backed set implementation.
package treeset

import (
	"cmp"
	"fmt"
	"iter"
	"strings"

	"github.com/lock14/collections"
	"github.com/lock14/collections/comparator"
	"github.com/lock14/collections/treemap"
)

var (
	_ collections.MutableNavigableSet[int] = (*TreeSet[int])(nil)
	_ fmt.Stringer                         = (*TreeSet[int])(nil)
)

// config holds the values for configuring a TreeSet.
type config[T any] struct {
	degree     *int
	comparator comparator.Comparator[T]
}

// Option configures a TreeSet config.
type Option[T any] func(*config[T])

// WithDegree configures the degree of the underlying B-Tree.
func WithDegree[T any](degree int) Option[T] {
	return func(c *config[T]) {
		c.degree = &degree
	}
}

// WithComparator configures the comparator for the TreeSet.
func WithComparator[T any](comp comparator.Comparator[T]) Option[T] {
	return func(c *config[T]) {
		c.comparator = comp
	}
}

// TreeSet represents a set of elements backed by a B-Tree.
type TreeSet[T any] struct {
	m *treemap.TreeMap[T, struct{}]
}

// New creates an empty TreeSet.
func New[T any](opts ...Option[T]) *TreeSet[T] {
	config := &config[T]{}
	for _, option := range opts {
		option(config)
	}

	var mapOpts []treemap.Option[T]
	if config.degree != nil {
		mapOpts = append(mapOpts, treemap.WithDegree[T](*config.degree))
	}
	if config.comparator != nil {
		mapOpts = append(mapOpts, treemap.WithComparator[T](config.comparator))
	}

	return &TreeSet[T]{
		m: treemap.New[T, struct{}](mapOpts...),
	}
}

// NewOrdered creates an empty TreeSet for types that implement cmp.Ordered.
func NewOrdered[T cmp.Ordered](opts ...Option[T]) *TreeSet[T] {
	opts = append([]Option[T]{WithComparator(comparator.NaturalOrder[T]())}, opts...)
	return New(opts...)
}

// Add inserts the specified element into the set.
func (s *TreeSet[T]) Add(item T) {
	s.m.Put(item, struct{}{})
}

// Remove removes and returns a single element from the set.
func (s *TreeSet[T]) Remove() T {
	for item := range s.m.Keys() {
		s.m.Remove(item)
		return item
	}
	panic("cannot remove from an empty set")
}

// RemoveElement removes the specified element from the set.
func (s *TreeSet[T]) RemoveElement(item T) {
	s.m.Remove(item)
}

// Contains returns true if this set contains the specified element.
func (s *TreeSet[T]) Contains(item T) bool {
	return s.m.ContainsKey(item)
}

// ContainsAll returns true if this set contains all elements of the specified collection.
func (s *TreeSet[T]) ContainsAll(other collections.Collection[T]) bool {
	for item := range other.All() {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// AddAll inserts all elements from the given sequence into the set.
func (s *TreeSet[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		s.Add(t)
	}
}

// RemoveAll removes all elements of the specified collection from this set.
func (s *TreeSet[T]) RemoveAll(other collections.Collection[T]) {
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

// RetainAll retains only the elements in this set that are contained in the specified collection.
func (s *TreeSet[T]) RetainAll(other collections.Collection[T]) {
	intersection := make([]T, 0, min(s.Size(), other.Size()))
	for t := range other.All() {
		if s.Contains(t) {
			intersection = append(intersection, t)
		}
	}
	s.Clear()
	for _, t := range intersection {
		s.Add(t)
	}
}

// Clear removes all elements from the set.
func (s *TreeSet[T]) Clear() {
	s.m.Clear()
}

// Size returns the number of elements in the set.
func (s *TreeSet[T]) Size() int {
	return s.m.Size()
}

// Empty returns true if the set contains no elements.
func (s *TreeSet[T]) Empty() bool {
	return s.m.Empty()
}

// All returns an Iterator over all the elements of this set.
func (s *TreeSet[T]) All() iter.Seq[T] {
	return s.m.Keys()
}

// First returns the first element in the set.
func (s *TreeSet[T]) First() T {
	k, _ := s.m.First()
	return k
}

// Last returns the last element in the set.
func (s *TreeSet[T]) Last() T {
	k, _ := s.m.Last()
	return k
}

// PollFirst removes and returns the first element in the set.
func (s *TreeSet[T]) PollFirst() T {
	k, _ := s.m.PollFirst()
	return k
}

// PollLast removes and returns the last element in the set.
func (s *TreeSet[T]) PollLast() T {
	k, _ := s.m.PollLast()
	return k
}

// Lower returns the greatest element strictly less than the given element, or false if not found.
func (s *TreeSet[T]) Lower(item T) (T, bool) {
	k, _, ok := s.m.Lower(item)
	return k, ok
}

// Floor returns the greatest element less than or equal to the given element, or false if not found.
func (s *TreeSet[T]) Floor(item T) (T, bool) {
	k, _, ok := s.m.Floor(item)
	return k, ok
}

// Ceiling returns the least element greater than or equal to the given element, or false if not found.
func (s *TreeSet[T]) Ceiling(item T) (T, bool) {
	k, _, ok := s.m.Ceiling(item)
	return k, ok
}

// Higher returns the least element strictly greater than the given element, or false if not found.
func (s *TreeSet[T]) Higher(item T) (T, bool) {
	k, _, ok := s.m.Higher(item)
	return k, ok
}

// Backward returns an iterator over the elements in reverse (descending) order.
func (s *TreeSet[T]) Backward() iter.Seq[T] {
	return s.m.BackwardKeys()
}

// From returns an iterator over the elements greater than or equal to from in ascending order.
func (s *TreeSet[T]) From(from T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s.m.From(from) {
			if !yield(k) {
				return
			}
		}
	}
}

// To returns an iterator over the elements strictly less than to in ascending order.
func (s *TreeSet[T]) To(to T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s.m.To(to) {
			if !yield(k) {
				return
			}
		}
	}
}

// Between returns an iterator over the elements in the half-open interval [from, to) in ascending order.
func (s *TreeSet[T]) Between(from, to T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s.m.Between(from, to) {
			if !yield(k) {
				return
			}
		}
	}
}

// String returns a string representation of the set matching Go slice formatting.
func (s *TreeSet[T]) String() string {
	vals := make([]string, 0, s.Size())
	for item := range s.All() {
		vals = append(vals, fmt.Sprintf("%v", item))
	}
	return "[" + strings.Join(vals, " ") + "]"
}
