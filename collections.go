package collections

import (
	"iter"
)

// Iterable denotes a type that can be iterated over
// by using the Iterator supplied using the Iterator method.
type Iterable[T any] interface {
	// All returns an Iterator over all the elements of this Iterable.
	All() iter.Seq[T]
	// Stream returns a channel containing the elements of this Iterable.
	// The channel provided will be closed automatically after all elements
	// have been read from the channel. If all elements from the channel
	// are not read, then the channel will not be closed. This method is
	// intended to be used primarily with the for...range construct.
	//
	//   for e := range i.Stream() {
	//      // do something with e
	//   }
	Stream() chan T
}

type Collection[T any] interface {
	Iterable[T]
	Size() int
	Empty() bool
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
}

type Map[K any, V any] interface {
	Put(K, V)
	Get(K) (V, bool)
	Remove(K)
	Size() int
	Empty() bool
	Entries() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
}
