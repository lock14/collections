package treeset

import (
	"fmt"
	"iter"
	"strings"

	"github.com/lock14/collections"
)

var _ collections.MutableNavigableSet[int] = (*TreeSet[int])(nil)

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

func (s *TreeSet[T]) First() (T, bool) {
	k, _, ok := s.m.First()
	return k, ok
}

func (s *TreeSet[T]) Last() (T, bool) {
	k, _, ok := s.m.Last()
	return k, ok
}

func (s *TreeSet[T]) PollFirst() (T, bool) {
	k, _, ok := s.m.PollFirst()
	return k, ok
}

func (s *TreeSet[T]) PollLast() (T, bool) {
	k, _, ok := s.m.PollLast()
	return k, ok
}

func (s *TreeSet[T]) AddFirst(item T) {
	panic("AddFirst is not supported on SortedSet")
}

func (s *TreeSet[T]) AddLast(item T) {
	panic("AddLast is not supported on SortedSet")
}

func (s *TreeSet[T]) Lower(item T) (T, bool) {
	k, _, ok := s.m.Lower(item)
	return k, ok
}

func (s *TreeSet[T]) Floor(item T) (T, bool) {
	k, _, ok := s.m.Floor(item)
	return k, ok
}

func (s *TreeSet[T]) Ceiling(item T) (T, bool) {
	k, _, ok := s.m.Ceiling(item)
	return k, ok
}

func (s *TreeSet[T]) Higher(item T) (T, bool) {
	k, _, ok := s.m.Higher(item)
	return k, ok
}

func (s *TreeSet[T]) ReversedAll() iter.Seq[T] {
	return s.m.ReversedKeys()
}

func (s *TreeSet[T]) AllFrom(from T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k, _ := range s.m.AllFrom(from) {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *TreeSet[T]) AllTo(to T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k, _ := range s.m.AllTo(to) {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *TreeSet[T]) AllBetween(from T, to T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k, _ := range s.m.AllBetween(from, to) {
			if !yield(k) {
				return
			}
		}
	}
}

// String returns a string representation of the set.
func (s *TreeSet[T]) String() string {
	vals := make([]string, 0, s.Size())
	for item := range s.All() {
		vals = append(vals, fmt.Sprintf("%+v", item))
	}
	return "[" + strings.Join(vals, ", ") + "]"
}
