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

type List[T any] interface {
	Collection[T]
	Add(t T)
	Remove() T
	AddAll(collection Collection[T])
}

type Queue[T any] interface {
	Collection[T]
	Add(t T)
	Remove() T
}

type Stack[T any] interface {
	Collection[T]
	Push(t T)
	Pop() T
}

type Deque[T any] interface {
	Stack[T]
	Queue[T]
	AddFront(t T)
	RemoveFront() T
	AddBack(t T)
	RemoveBack() T
}

type Set[T any] interface {
	Collection[T]
	Add(t T)
	Remove(t T)
	AddAll(collection Collection[T])
	RemoveAll(collection Collection[T])
	RetainAll(collection Collection[T])
}

type Map[K any, V any] interface {
	Put(K, V)
	Get(K) (V, bool)
	Remove(K)
	Size() int
	Empty() bool
	All() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
}
