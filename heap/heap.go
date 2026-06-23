package heap

import (
	"cmp"
	"github.com/lock14/collections"
	"iter"
	"slices"
)

const (
	DefaultCapacity = 10
)

var _ collections.MutableQueue[int] = (*Heap[int])(nil)

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

// Capacity configures the initial pre-allocated capacity of the heap.
func Capacity[T any](capacity int) Option[T] {
	return func(config *Config[T]) {
		config.capacity = capacity
	}
}

type Heap[T any] struct {
	elements   []T
	comparator Comparator[T]
}

func New[T any](opts ...Option[T]) *Heap[T] {
	config := defaultConfig[T]()
	for _, opt := range opts {
		opt(config)
	}
	return &Heap[T]{
		elements:   make([]T, 0, config.capacity),
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
	h.elements = append(h.elements, t)
	h.siftUp(len(h.elements) - 1)
}

func (h *Heap[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		h.Add(t)
	}
}

// Remove removes and returns the top element from the heap.
// Panics if the heap is empty.
func (h *Heap[T]) Remove() T {
	if h.Empty() {
		panic("heap is empty")
	}
	t := h.elements[0]
	h.delete(0)
	return t
}

// Peek returns the top element from the heap without removing it.
// Panics if the heap is empty.
func (h *Heap[T]) Peek() T {
	if h.Empty() {
		panic("heap is empty")
	}
	return h.elements[0]
}

func (h *Heap[T]) Size() int {
	return len(h.elements)
}

func (h *Heap[T]) Empty() bool {
	return len(h.elements) == 0
}

func (h *Heap[T]) Clear() {
	// Zero out elements to allow GC
	var zero T
	for i := range h.elements {
		h.elements[i] = zero
	}
	h.elements = h.elements[:0]
}

func (h *Heap[T]) All() iter.Seq[T] {
	return slices.Values(h.elements)
}

// Private Functions

func defaultConfig[T any]() *Config[T] {
	return &Config[T]{
		capacity: DefaultCapacity,
	}
}

func (h *Heap[T]) delete(index int) {
	last := len(h.elements) - 1
	var zero T
	if index != last {
		h.elements[index] = h.elements[last]
		h.elements[last] = zero
		h.elements = h.elements[:last]
		if h.hasParent(index) && !h.heapCondition(parent(index), index) {
			h.siftUp(index)
		} else {
			h.siftDown(index)
		}
	} else {
		h.elements[last] = zero
		h.elements = h.elements[:last]
	}
}

func (h *Heap[T]) siftUp(cur int) {
	for h.hasParent(cur) && h.heapCondition(cur, parent(cur)) {
		h.swap(parent(cur), cur)
		cur = parent(cur)
	}
}

func (h *Heap[T]) siftDown(cur int) {
	for {
		best := cur
		if h.hasLeft(cur) && h.heapCondition(left(cur), best) {
			best = left(cur)
		}
		if h.hasRight(cur) && h.heapCondition(right(cur), best) {
			best = right(cur)
		}
		if best == cur {
			break
		}
		h.swap(cur, best)
		cur = best
	}
}

// heapCondition returns true if the element at index i should be placed higher in the heap than the element at index j.
func (h *Heap[T]) heapCondition(i, j int) bool {
	return h.comparator(h.elements[i], h.elements[j]) <= 0
}

func (h *Heap[T]) swap(i, j int) {
	h.elements[j], h.elements[i] = h.elements[i], h.elements[j]
}

func (h *Heap[T]) hasParent(index int) bool {
	return index > 0
}

func (h *Heap[T]) hasLeft(index int) bool {
	return left(index) < len(h.elements)
}

func (h *Heap[T]) hasRight(index int) bool {
	return right(index) < len(h.elements)
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
