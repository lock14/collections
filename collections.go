// Package collections provides generic data structures and interfaces for Go.
package collections

import (
	"iter"
)

// Iterable denotes a type that can be iterated over
// by using the Iterator supplied using the All method.
type Iterable[T any] interface {
	// All returns an Iterator over all the elements of this Iterable.
	All() iter.Seq[T]
}

// Collection represents a group of elements.
type Collection[T any] interface {
	Iterable[T]
	// Size returns the number of elements in the collection.
	Size() int
	// Empty returns true if the collection contains no elements.
	Empty() bool
}

// MutableCollection represents a collection that can be modified.
type MutableCollection[T any] interface {
	Collection[T]
	// Add inserts the specified element into the collection.
	Add(t T)
	// Remove removes and returns a single element from the collection.
	Remove() T
	// AddAll inserts all elements from the given sequence into the collection.
	AddAll(iter.Seq[T])
	// Clear removes all elements from the collection.
	Clear()
}

// List represents an ordered collection (also known as a sequence).
type List[T any] interface {
	Collection[T]
	// Get returns the element at the specified index.
	Get(idx int) T
}

// MutableList represents an ordered collection that can be modified.
type MutableList[T any] interface {
	List[T]
	MutableCollection[T]
	// Set replaces the element at the specified index with the given element.
	Set(idx int, t T)
}

// Queue represents a collection designed for holding elements prior to processing.
type Queue[T any] interface {
	Collection[T]
	// Peek returns the element at the front of the queue without removing it.
	Peek() T
}

// MutableQueue represents a queue that can be modified.
type MutableQueue[T any] interface {
	Queue[T]
	MutableCollection[T]
}

// Stack represents a last-in-first-out (LIFO) stack of objects.
type Stack[T any] interface {
	Collection[T]
	// Peek returns the element at the top of the stack without removing it.
	Peek() T
}

// MutableStack represents a stack that can be modified.
type MutableStack[T any] interface {
	Stack[T]
	MutableCollection[T]
	// Push adds the specified element to the top of the stack.
	Push(t T)
	// Pop removes and returns the element at the top of the stack.
	Pop() T
}

// Deque represents a linear collection that supports element insertion and removal at both ends.
type Deque[T any] interface {
	Stack[T]
	Queue[T]
	// PeekFront returns the element at the front of the deque without removing it.
	PeekFront() T
	// PeekBack returns the element at the back of the deque without removing it.
	PeekBack() T
}

// MutableDeque represents a deque that can be modified.
type MutableDeque[T any] interface {
	Deque[T]
	MutableCollection[T]
	MutableStack[T]
	MutableQueue[T]
	// AddFront inserts the specified element at the front of the deque.
	AddFront(t T)
	// RemoveFront removes and returns the element at the front of the deque.
	RemoveFront() T
	// AddBack inserts the specified element at the back of the deque.
	AddBack(t T)
	// RemoveBack removes and returns the element at the back of the deque.
	RemoveBack() T
}

// Set represents a collection that contains no duplicate elements.
type Set[T any] interface {
	Collection[T]
	// Contains returns true if this set contains the specified element.
	Contains(T) bool
	// ContainsAll returns true if this set contains all elements of the specified collection.
	ContainsAll(Collection[T]) bool
}

// MutableSet represents a set that can be modified.
type MutableSet[T any] interface {
	MutableCollection[T]
	Set[T]
	// RemoveElement removes the specified element from the set.
	RemoveElement(T)
	// RemoveAll removes all elements of the specified collection from this set.
	RemoveAll(Collection[T])
	// RetainAll retains only the elements in this set that are contained in the specified collection.
	RetainAll(Collection[T])
}

// Map represents a collection of key-value pairs where each key is unique.
type Map[K any, V any] interface {
	// Get returns the value associated with the specified key, and a boolean indicating if it was found.
	Get(K) (V, bool)
	// Size returns the number of key-value pairs in the map.
	Size() int
	// Empty returns true if the map contains no key-value pairs.
	Empty() bool
	// ContainsKey returns true if the map contains a mapping for the specified key.
	ContainsKey(K) bool
	// All returns an iterator over all key-value pairs in the map.
	All() iter.Seq2[K, V]
	// Keys returns an iterator over all keys in the map.
	Keys() iter.Seq[K]
	// Values returns an iterator over all values in the map.
	Values() iter.Seq[V]
}

// MutableMap represents a map that can be modified.
type MutableMap[K any, V any] interface {
	Map[K, V]
	// Put associates the specified value with the specified key in the map.
	Put(K, V)
	// Remove removes the mapping for the specified key from the map if present.
	Remove(K)
	// Clear removes all key-value pairs from the map.
	Clear()
}

// SequencedSet represents a set with a defined encounter order.
type SequencedSet[T any] interface {
	Set[T]
	// First returns the first element in the set. Panics if empty.
	First() T
	// Last returns the last element in the set. Panics if empty.
	Last() T
	// Backward returns an iterator over the elements in reverse order.
	Backward() iter.Seq[T]
}

// MutableSequencedSet represents a sequenced set that can be modified.
type MutableSequencedSet[T any] interface {
	SequencedSet[T]
	MutableSet[T]
	// PollFirst removes and returns the first element in the set. Panics if empty.
	PollFirst() T
	// PollLast removes and returns the last element in the set. Panics if empty.
	PollLast() T
	// AddFirst inserts the specified element at the front of the set.
	AddFirst(T)
	// AddLast inserts the specified element at the end of the set.
	AddLast(T)
}

// SequencedMap represents a map with a defined encounter order.
type SequencedMap[K any, V any] interface {
	Map[K, V]
	// First returns the first key-value pair in the map. Panics if empty.
	First() (K, V)
	// Last returns the last key-value pair in the map. Panics if empty.
	Last() (K, V)
	// Backward returns an iterator over the key-value pairs in reverse order.
	Backward() iter.Seq2[K, V]
	// BackwardKeys returns an iterator over the keys in reverse order.
	BackwardKeys() iter.Seq[K]
	// BackwardValues returns an iterator over the values in reverse order.
	BackwardValues() iter.Seq[V]
}

// MutableSequencedMap represents a sequenced map that can be modified.
type MutableSequencedMap[K any, V any] interface {
	SequencedMap[K, V]
	MutableMap[K, V]
	// PollFirst removes and returns the first key-value pair in the map. Panics if empty.
	PollFirst() (K, V)
	// PollLast removes and returns the last key-value pair in the map. Panics if empty.
	PollLast() (K, V)
	// PutFirst inserts the specified key-value pair at the front of the map.
	PutFirst(K, V)
	// PutLast inserts the specified key-value pair at the end of the map.
	PutLast(K, V)
}

// SortedSet represents a set that maintains its elements in ascending order.
type SortedSet[T any] interface {
	SequencedSet[T]
	// From returns an iterator over the elements greater than or equal to the given element.
	From(from T) iter.Seq[T]
	// To returns an iterator over the elements less than the given element.
	To(to T) iter.Seq[T]
	// Between returns an iterator over the elements greater than or equal to 'from' and less than 'to'.
	Between(from, to T) iter.Seq[T]
}

// MutableSortedSet represents a sorted set that can be modified.
type MutableSortedSet[T any] interface {
	SortedSet[T]
	MutableSet[T]
	// PollFirst removes and returns the first element in the set. Panics if empty.
	PollFirst() T
	// PollLast removes and returns the last element in the set. Panics if empty.
	PollLast() T
}

// SortedMap represents a map that maintains its entries in ascending order of keys.
type SortedMap[K any, V any] interface {
	SequencedMap[K, V]
	// From returns an iterator over the entries whose keys are greater than or equal to the given key.
	From(from K) iter.Seq2[K, V]
	// To returns an iterator over the entries whose keys are less than the given key.
	To(to K) iter.Seq2[K, V]
	// Between returns an iterator over the entries whose keys are greater than or equal to 'from' and less than 'to'.
	Between(from, to K) iter.Seq2[K, V]
}

// MutableSortedMap represents a sorted map that can be modified.
type MutableSortedMap[K any, V any] interface {
	SortedMap[K, V]
	MutableMap[K, V]
	// PollFirst removes and returns the first key-value pair in the map. Panics if empty.
	PollFirst() (K, V)
	// PollLast removes and returns the last key-value pair in the map. Panics if empty.
	PollLast() (K, V)
}

// NavigableSet represents a sorted set with routing methods for finding closest matches.
type NavigableSet[T any] interface {
	SortedSet[T]
	// Lower returns the greatest element strictly less than the given element.
	Lower(T) (T, bool)
	// Floor returns the greatest element less than or equal to the given element.
	Floor(T) (T, bool)
	// Ceiling returns the least element greater than or equal to the given element.
	Ceiling(T) (T, bool)
	// Higher returns the least element strictly greater than the given element.
	Higher(T) (T, bool)
}

// MutableNavigableSet represents a navigable set that can be modified.
type MutableNavigableSet[T any] interface {
	NavigableSet[T]
	MutableSortedSet[T]
}

// NavigableMap represents a sorted map with routing methods for finding closest key matches.
type NavigableMap[K any, V any] interface {
	SortedMap[K, V]
	// Lower returns the key-value pair for the greatest key strictly less than the given key.
	Lower(K) (K, V, bool)
	// Floor returns the key-value pair for the greatest key less than or equal to the given key.
	Floor(K) (K, V, bool)
	// Ceiling returns the key-value pair for the least key greater than or equal to the given key.
	Ceiling(K) (K, V, bool)
	// Higher returns the key-value pair for the least key strictly greater than the given key.
	Higher(K) (K, V, bool)
}

// MutableNavigableMap represents a navigable map that can be modified.
type MutableNavigableMap[K any, V any] interface {
	NavigableMap[K, V]
	MutableSortedMap[K, V]
}
