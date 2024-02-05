package iterator

import "github.com/lock14/collections/util"

// Iterator represents an iterator that moves 'forward' whenever
// Increment() is called.
type Iterator[T any] interface {
	// Empty returns true if this Iterator is empty.
	Empty() bool
	// Next returns the next element from this Iterator.
	// If the Iterator is empty, then an error is returned.
	Next() (*T, error)
}

// Iterable denotes a type that can be iterated over
// by using the Iterator supplied using the Iterator method.
type Iterable[T any] interface {
	// Iterator returns a Iterator over the elements of this Iterable.
	Iterator() Iterator[T]
	// Elements returns a channel containing the elements of this Iterable.
	// The channel provided will be closed automatically after all elements
	// have been read from the channel. If all elements from the channel
	// are not read, then the channel will not be closed. This method is
	// intended to be used primarily with the for...range construct.
	//
	//   for e := range i.Elements() {
	//      // do something with e
	//   }
	Elements() chan *T
}

// Elements produces a channel that contains the elements of the iterable.
// The channel provided will be closed automatically after all elements
// have been read from the channel. If all elements from the channel
// are not read, then the channel will not be closed. This method is
// intended to be used primarily with the for...range construct.
//
//	for e := range iterator.Elements(iterable) {
//	   // do something with e
//	}
func Elements[T any](iterable Iterable[T]) chan *T {
	c := make(chan *T)
	go func() {
		for itr := iterable.Iterator(); !itr.Empty(); {
			c <- util.MustGet(itr.Next())
		}
		close(c)
	}()
	return c

}
