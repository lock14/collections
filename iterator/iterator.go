package iterator

type ForwardIterator[T any] interface {
	Empty() bool
	PopFront() error
	Front() (*T, error)
}

type BidirectionalIterator[T any] interface {
	ForwardIterator[T]
	Back(*T, error)
	PopBack() error
}

type RandomAccessIterator[T any] interface {
	BidirectionalIterator[T]
	Get(index int) (*T, error)
	Length() int
}
