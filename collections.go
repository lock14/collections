package collections

import (
	"iter"
)

// Iterable denotes a type that can be iterated over
// by using the Iterator supplied using the Iterator method.
type Iterable[T any] interface {
	// All returns an Iterator over all the elements of this Iterable.
	All() iter.Seq[T]
}

type Collection[T any] interface {
	Iterable[T]
	Size() int
	Empty() bool
}

type MutableCollection[T any] interface {
	Collection[T]
	Add(t T)
	Remove() T
	AddAll(iter.Seq[T])
}

type List[T any] interface {
	Collection[T]
	Get(idx int) T
}

type MutableList[T any] interface {
	List[T]
	MutableCollection[T]
	Set(idx int, t T)
}

type Queue[T any] interface {
	Collection[T]
	Peek() T
}

type MutableQueue[T any] interface {
	Queue[T]
	MutableCollection[T]
}

type Stack[T any] interface {
	Collection[T]
	Peek() T
}

type MutableStack[T any] interface {
	Stack[T]
	MutableCollection[T]
	Push(t T)
	Pop() T
}

type Deque[T any] interface {
	Stack[T]
	Queue[T]
	PeekFront() T
	PeekBack() T
}

type MutableDeque[T any] interface {
	Deque[T]
	MutableCollection[T]
	MutableStack[T]
	MutableQueue[T]
	AddFront(t T)
	RemoveFront() T
	AddBack(t T)
	RemoveBack() T
}

type Set[T any] interface {
	Collection[T]
	Contains(T) bool
	ContainsAll(Collection[T]) bool
}

type MutableSet[T any] interface {
	MutableCollection[T]
	Set[T]
	RemoveElement(T)
	RemoveAll(Collection[T])
	RetainAll(Collection[T])
}

type Map[K any, V any] interface {
	Get(K) (V, bool)
	Size() int
	Empty() bool
	ContainsKey(K) bool
	All() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
}

type MutableMap[K any, V any] interface {
	Map[K, V]
	Put(K, V)
	Remove(K)
}
