package iterator

// ForwardIterator represents an iterator that moves 'forward' whenever
// Increment() is called.
type ForwardIterator[T any] interface {
	// Empty returns true if this Iterator is empty.
	Empty() bool
	// Increment advances this Iterator one element forwards.
	// If the Iterator is empty, then an error is returned.
	Increment() error
	// MustIncrement advances this Iterator one element forwards.
	// If the Iterator is empty, then a panic occurs.
	MustIncrement()
	// GetFront returns the element at the front of this Iterator.
	// If the Iterator is empty, then an error is returned.
	GetFront() (*T, error)
	// MustGetFront returns the element at the front of this Iterator.
	// If the Iterator is empty, then a panic occurs.
	MustGetFront() *T
}

// BidirectionalIterator represents an iterator that can move both 'forwards'
// and 'backwards'.
type BidirectionalIterator[T any] interface {
	ForwardIterator[T]
	// GetBack returns the element at the back of this Iterator.
	// If the Iterator is empty, then an error is returned.
	GetBack() (*T, error)
	// MustGetBack returns the element at the back of this Iterator.
	// If the Iterator is empty, then a panic occurs.
	MustGetBack() *T
	// Decrement advances this Iterator one element backwards.
	// If the Iterator is empty, then an error is returned.
	Decrement() error
	// MustDecrement advances this Iterator one element backwards.
	// If the Iterator is empty, then a panic occurs.
	MustDecrement()
}

// RandomAccessIterator represents an iterator where the elements of the range
// being iterated over can be randomly accessed.
type RandomAccessIterator[T any] interface {
	BidirectionalIterator[T]
	// Get returns the element at the specified index.
	// If the Iterator is empty, then an error is returned.
	Get(index int) (*T, error)
	// Length returns the length of this Iterator.
	Length() int
}

// ForwardIterable denotes a type that can be iterated over
// by using the ForwardIterator supplied using the Iterator method.
type ForwardIterable[T any] interface {
	// Iterator returns a ForwardIterator over the elements of this Iterable.
	Iterator() ForwardIterator[T]
	// Elements returns a channel containing the elements of this Iterable.
	// The channel provided will be closed automatically after all elements
	// have been read from the channel. If all elements from the channel
	// are not read, then the channel will not be closed. This method is
	// intended to be used primarily with the for...range construct.
	//
	//   for e := range i.Elements {
	//      // do something with e
	//   }
	Elements() chan *T
}

// Range produces a channel that contains the elements of the iterable.
// The channel provided will be closed automatically after all elements
// have been read from the channel. If all elements from the channel
// are not read, then the channel will not be closed. This method is
// intended to be used primarily with the for...range construct.
//
//	for e := range iterator.Range(iterable) {
//	   // do something with e
//	}
func Range[T any](iterable ForwardIterable[T]) chan *T {
	c := make(chan *T)
	go func() {
		for itr := iterable.Iterator(); !itr.Empty(); itr.MustIncrement() {
			c <- itr.MustGetFront()
		}
		close(c)
	}()
	return c

}
