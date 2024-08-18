package heap

import (
	"cmp"
	"iter"
)

const (
	DefaultCapacity = 10
)

type Comparator[T any] func(t1, t2 T) int

func NaturalOrder[T cmp.Ordered]() Comparator[T] {
	return func(t1, t2 T) int {
		if t1 < t2 {
			return -1
		} else if t1 > t2 {
			return 1
		} else {
			return 0
		}
	}
}

func Reversed[T any](comparator Comparator[T]) Comparator[T] {
	return func(t1, t2 T) int {
		return -comparator(t1, t2)
	}
}

type Option[T any] func(config *Config[T])

type Config[T any] struct {
	capacity   int
	comparator Comparator[T]
}

func WithComparator[T any](comparator Comparator[T]) Option[T] {
	return func(config *Config[T]) {
		config.comparator = comparator
	}
}

type Heap[T any] struct {
	elements   []T
	size       int
	zero       T
	comparator Comparator[T]
}

func New[T any](opts ...Option[T]) *Heap[T] {
	config := defaultConfig[T]()
	for _, opt := range opts {
		opt(config)
	}
	return &Heap[T]{
		comparator: config.comparator,
	}
}

func Min[T cmp.Ordered]() *Heap[T] {
	return New[T](WithComparator(NaturalOrder[T]()))
}

func Max[T cmp.Ordered]() *Heap[T] {
	return New[T](WithComparator(Reversed(NaturalOrder[T]())))
}

func (h *Heap[T]) Add(t T) {
	if h.size == len(h.elements) {
		h.increaseCapacity()
	}
	h.elements[h.size] = t
	h.size++
	h.siftUp(h.size - 1)
}

func (h *Heap[T]) Remove() T {
	t := h.elements[0]
	h.delete(0)
	return t
}

func (h *Heap[T]) Peek() *T {
	return &h.elements[0]
}

func (h *Heap[T]) Size() int {
	return h.size
}

func (h *Heap[T]) Empty() bool {
	return h.Size() == 0
}

func defaultConfig[T any]() *Config[T] {
	return &Config[T]{
		capacity: DefaultCapacity,
	}
}

func (h *Heap[T]) delete(index int) {
	var zero T
	h.size--
	h.elements[index] = h.elements[h.size]
	h.elements[h.size] = zero
	if h.hasParent(index) && !h.heapCondition(parent(index), index) {
		h.siftUp(index)
	} else {
		h.siftDown(index)
	}

}

func (h *Heap[T]) siftUp(cur int) {
	for h.hasParent(cur) && h.heapCondition(parent(cur), cur) {
		h.swap(parent(cur), cur)
		cur = parent(cur)
	}
}

func (h *Heap[T]) siftDown(cur int) {
	heapConditionViolated := true
	for h.hasLeft(cur) && heapConditionViolated {
		child := left(cur)
		if h.hasRight(cur) && h.heapCondition(right(cur), child) {
			child = right(cur)
		}
		heapConditionViolated = !h.heapCondition(cur, child)
		if heapConditionViolated {
			h.swap(cur, child)
		}
		cur = child
	}
}

func (h *Heap[T]) heapCondition(i, j int) bool {
	return h.comparator(h.elements[i], h.elements[j]) <= 0
}

func (h *Heap[T]) increaseCapacity() {
	newElements := make([]T, len(h.elements)+(len(h.elements)>>1))
	copy(newElements, h.elements)
	h.elements = newElements
}

func (h *Heap[T]) swap(i, j int) {
	h.elements[j], h.elements[i] = h.elements[i], h.elements[j]
}

func (h *Heap[T]) hasParent(index int) bool {
	return index > 0
}

func (h *Heap[T]) hasLeft(index int) bool {
	return left(index) < h.size
}

func (h *Heap[T]) hasRight(index int) bool {
	return right(index) < h.size
}

func parent(index int) int {
	return (index - 1) >> 1
}

func left(index int) int {
	return 2*index + 1
}

func right(index int) int {
	return 2*index + 2
}

func (h *Heap[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		hi := heapIterator[T]{
			heap:  h,
			index: 0,
		}
		for !hi.empty() && yield(hi.next()) {
		}
	}
}

// All

type heapIterator[T any] struct {
	heap  *Heap[T]
	index int
}

func (hi *heapIterator[T]) empty() bool {
	return hi.index >= hi.heap.Size()
}

func (hi *heapIterator[T]) next() T {
	t := hi.heap.elements[hi.index]
	hi.index++
	return t
}
