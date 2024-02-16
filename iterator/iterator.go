package iterator

import "github.com/lock14/collections/util"

// Iterator generates a sequence of elements, one at a time.
type Iterator[T any] interface {
	// Empty returns true if this Iterator is empty.
	Empty() bool
	// Next returns the next element from this Iterator.
	// If the Iterator is empty, then an error is returned.
	Next() (T, error)
}

// Elements produces a channel that contains the elements of the iterator.
// The channel provided will be closed automatically after all elements
// have been read from the channel. If all elements from the channel
// are not read, then the channel will not be closed. This method is
// intended to be used primarily with the for...range construct.
//
//	for e := range iterator.Elements(iterator) {
//	   // do something with e
//	}
func Elements[T any](itr Iterator[T]) chan T {
	c := make(chan T)
	go func() {
		for !itr.Empty() {
			c <- util.MustGet(itr.Next())
		}
		close(c)
	}()
	return c
}
