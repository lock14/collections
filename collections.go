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
