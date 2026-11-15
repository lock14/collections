// Package treeset provides a sorted, navigable B-Tree set implementation.
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

// config holds configuration values for constructing a TreeSet.
type config[T any] struct {
	degree     *int
	comparator comparator.Comparator[T]
}

// Option configures a TreeSet.
type Option[T any] func(*config[T])

// WithDegree configures the minimum degree (t) of the underlying B-Tree.
// Panics in New if degree is less than 2.
func WithDegree[T any](degree int) Option[T] {
	return func(c *config[T]) {
		c.degree = &degree
	}
}

// WithComparator configures the comparator used by the TreeSet.
// Panics in New if cmpFunc is nil.
func WithComparator[T any](cmpFunc comparator.Comparator[T]) Option[T] {
	return func(c *config[T]) {
		c.comparator = cmpFunc
	}
}

// TreeSet represents a sorted set of unique elements backed by a B-Tree.
type TreeSet[T any] struct {
	m *treemap.TreeMap[T, struct{}]
}

// New creates an empty TreeSet configured with the given options.
// Panics if degree is less than 2 or if no non-nil comparator is provided.
func New[T any](opts ...Option[T]) *TreeSet[T] {
	cfg := &config[T]{}
	for _, option := range opts {
		option(cfg)
	}

	var mapOpts []treemap.Option[T]
	if cfg.degree != nil {
		mapOpts = append(mapOpts, treemap.WithDegree[T](*cfg.degree))
	}
	if cfg.comparator != nil {
		mapOpts = append(mapOpts, treemap.WithComparator[T](cfg.comparator))
	}

	return &TreeSet[T]{
		m: treemap.New[T, struct{}](mapOpts...),
	}
}

// NewOrdered creates an empty TreeSet for cmp.Ordered elements using natural ascending order.
// Any options provided in opts are applied after the default natural-order comparator.
func NewOrdered[T cmp.Ordered](opts ...Option[T]) *TreeSet[T] {
	allOpts := make([]Option[T], 0, len(opts)+1)
	allOpts = append(allOpts, WithComparator(comparator.NaturalOrder[T]()))
	allOpts = append(allOpts, opts...)
	return New[T](allOpts...)
}

// Add inserts the specified element into the set.
func (s *TreeSet[T]) Add(item T) {
	s.m.Put(item, struct{}{})
}

// AddAll inserts all elements from the given sequence into the set.
func (s *TreeSet[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		s.Add(t)
	}
}

// Remove removes and returns the smallest element in the set.
// Panics if the set is empty.
func (s *TreeSet[T]) Remove() T {
	return s.PollFirst()
}

// PollFirst removes and returns the smallest element in the set.
// Panics if the set is empty.
func (s *TreeSet[T]) PollFirst() T {
	k, _ := s.m.PollFirst()
	return k
}

// PollLast removes and returns the largest element in the set.
// Panics if the set is empty.
func (s *TreeSet[T]) PollLast() T {
	k, _ := s.m.PollLast()
	return k
}

// RemoveElement removes the specified element from the set if present.
func (s *TreeSet[T]) RemoveElement(item T) {
	s.m.Remove(item)
}

// RemoveAll removes all elements of the specified collection from this set.
func (s *TreeSet[T]) RemoveAll(other collections.Collection[T]) {
	if s.Empty() {
		return
	}
	if otherSet, ok := other.(*TreeSet[T]); ok && otherSet == s {
		s.Clear()
		return
	}
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

// RetainAll retains only the elements in this set that are contained in the specified collection.
func (s *TreeSet[T]) RetainAll(other collections.Collection[T]) {
	if s.Empty() {
		return
	}
	if otherSet, ok := other.(*TreeSet[T]); ok && otherSet == s {
		return
	}
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

// Contains returns true if this set contains the specified element.
func (s *TreeSet[T]) Contains(item T) bool {
	return s.m.ContainsKey(item)
}

// ContainsAll returns true if this set contains all elements of the specified collection.
func (s *TreeSet[T]) ContainsAll(other collections.Collection[T]) bool {
	if otherSet, ok := other.(*TreeSet[T]); ok && otherSet == s {
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
func (s *TreeSet[T]) Size() int {
	return s.m.Size()
}

// Empty returns true if the set contains no elements.
func (s *TreeSet[T]) Empty() bool {
	return s.m.Empty()
}

// First returns the smallest element in the set without removing it.
// Panics if the set is empty.
func (s *TreeSet[T]) First() T {
	k, _ := s.m.First()
	return k
}

// Last returns the largest element in the set without removing it.
// Panics if the set is empty.
func (s *TreeSet[T]) Last() T {
	k, _ := s.m.Last()
	return k
}

// Lower returns the greatest element strictly less than item, or false if not found.
func (s *TreeSet[T]) Lower(item T) (T, bool) {
	k, _, ok := s.m.Lower(item)
	return k, ok
}

// Floor returns the greatest element less than or equal to item, or false if not found.
func (s *TreeSet[T]) Floor(item T) (T, bool) {
	k, _, ok := s.m.Floor(item)
	return k, ok
}

// Ceiling returns the least element greater than or equal to item, or false if not found.
func (s *TreeSet[T]) Ceiling(item T) (T, bool) {
	k, _, ok := s.m.Ceiling(item)
	return k, ok
}

// Higher returns the least element strictly greater than item, or false if not found.
func (s *TreeSet[T]) Higher(item T) (T, bool) {
	k, _, ok := s.m.Higher(item)
	return k, ok
}

// All returns an iterator over all elements of this set in ascending order.
func (s *TreeSet[T]) All() iter.Seq[T] {
	return s.m.Keys()
}

// Backward returns an iterator over the elements in descending order.
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
