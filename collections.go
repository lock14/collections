package collections

import "github.com/lock14/collections/iterator"

// Iterable denotes a type that can be iterated over
// by using the Iterator supplied using the Iterator method.
type Iterable[T any] interface {
	// Iterator returns an Iterator over the elements of this Iterable.
	Iterator() iterator.Iterator[T]
	// Elements returns a channel containing the elements of this Iterable.
	// The channel provided will be closed automatically after all elements
	// have been read from the channel. If all elements from the channel
	// are not read, then the channel will not be closed. This method is
	// intended to be used primarily with the for...range construct.
	//
	//   for e := range i.Elements() {
	//      // do something with e
	//   }
	Elements() chan T
}
