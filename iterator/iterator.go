package iterator

// ForwardIterator represents an iterator that moves 'forward' whenever
// PopFront() is called.
type ForwardIterator[T any] interface {
	// Empty returns true if this Iterator is empty.
	Empty() bool
	// PopFront advances this Iterator one element forwards.
	// If the Iterator is empty, then an error is returned.
	PopFront() error
	// MustPopFront advances this Iterator one element forwards.
	// If the Iterator is empty, then a panic occurs.
	MustPopFront()
	// Front returns the element at the front of this Iterator.
	// If the Iterator is empty, then an error is returned.
	Front() (*T, error)
	// MustGetFront returns the element at the front of this Iterator.
	// If the Iterator is empty, then a panic occurs.
	MustGetFront() *T
}

// BidirectionalIterator represents an iterator that can move both 'forwards'
// and 'backwards'.
type BidirectionalIterator[T any] interface {
	ForwardIterator[T]
	// Back returns the element at the back of this Iterator.
	// If the Iterator is empty, then an error is returned.
	Back(*T, error)
	// MustGetBack returns the element at the back of this Iterator.
	// If the Iterator is empty, then a panic occurs.
	MustGetBack() *T
	// PopBack advances this Iterator one element backwards.
	// If the Iterator is empty, then an error is returned.
	PopBack() error
	// MustPopBack advances this Iterator one element backwards.
	// If the Iterator is empty, then a panic occurs.
	MustPopBack()
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
