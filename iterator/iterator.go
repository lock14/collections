package iterator

import (
	"iter"
)

// Stream produces a channel that contains the elements of the iterator.
// The channel provided will be closed automatically after all elements
// have been read from the channel. If all elements from the channel
// are not read, then the channel will not be closed. This method is
// intended to be used primarily with the for...range construct.
//
//	for e := range iterator.Stream(iterator) {
//	   // do something with e
//	}
func Stream[T any](seq iter.Seq[T]) chan T {
	c := make(chan T)
	go func() {
		for t := range seq {
			c <- t
		}
		close(c)
	}()
	return c
}
