// Package heap provides a binary heap priority queue.
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

// Comparator is a function that compares two elements.
type Comparator[T any] func(t1, t2 T) int

// NaturalOrder returns a comparator that orders elements using their natural ordering.
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

// Reversed returns a comparator that reverses the ordering of the given comparator.
func Reversed[T any](comparator Comparator[T]) Comparator[T] {
	return func(t1, t2 T) int {
		return -comparator(t1, t2)
	}
}

// Option is a function that configures a Heap.
type Option[T any] func(config *Config[T])

// Config holds the configuration for a Heap.
type Config[T any] struct {
	capacity   int
	comparator Comparator[T]
}

// WithComparator configures the comparator used by the Heap.
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

// Heap is a binary heap priority queue.
type Heap[T any] struct {
	elements   []T
	comparator Comparator[T]
}

// New creates a new Heap with the given options.
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

// Min creates a new Min-Heap using natural ordering.
func Min[T cmp.Ordered]() *Heap[T] {
	return New[T](WithComparator(NaturalOrder[T]()))
}

// Max creates a new Max-Heap using reversed natural ordering.
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

		p := (index - 1) >> 1
		if index > 0 && h.comparator(h.elements[index], h.elements[p]) <= 0 {
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
	elements := h.elements
	item := elements[cur]
	for cur > 0 {
		p := (cur - 1) >> 1
		if h.comparator(item, elements[p]) <= 0 {
			elements[cur] = elements[p]
			cur = p
		} else {
			break
		}
	}
	elements[cur] = item
}

func (h *Heap[T]) siftDown(cur int) {
	elements := h.elements
	n := len(elements)
	item := elements[cur]
	half := n >> 1
	for cur < half {
		child := 2*cur + 1
		right := child + 1
		if right < n && h.comparator(elements[right], elements[child]) <= 0 {
			child = right
		}
		if h.comparator(item, elements[child]) <= 0 {
			break
		}
		elements[cur] = elements[child]
		cur = child
	}
	elements[cur] = item
}
