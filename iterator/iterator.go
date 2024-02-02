package iterator

// ForwardIterator represents an iterator that moves 'forward' whenever
// PopFront() is called.
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
