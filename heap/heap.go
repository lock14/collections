// Package heap provides a binary heap priority queue.
package heap

import (
	"cmp"
	"fmt"
	"github.com/lock14/collections"
	"github.com/lock14/collections/comparator"
	"iter"
	"slices"
	"strings"
)

const (
	// DefaultCapacity is the default initial capacity for a Heap.
	DefaultCapacity = 10
)

var (
	_ collections.MutableQueue[int] = (*Heap[int])(nil)
	_ fmt.Stringer                  = (*Heap[int])(nil)
)

// Option is a function that configures a Heap.
type Option[T any] func(config *config[T])

// config holds the configuration for a Heap.
type config[T any] struct {
	capacity   int
	comparator comparator.Comparator[T]
}

// WithComparator configures the comparator used by the Heap.
func WithComparator[T any](cmpFunc comparator.Comparator[T]) Option[T] {
	return func(config *config[T]) {
		config.comparator = cmpFunc
	}
}

// WithCapacity configures the initial pre-allocated capacity of the heap.
func WithCapacity[T any](capacity int) Option[T] {
	return func(config *config[T]) {
		config.capacity = capacity
	}
}

// Heap is a binary heap priority queue.
type Heap[T any] struct {
	elements   []T
	comparator comparator.Comparator[T]
}

// New creates a new Heap with the given options.
// Panics if no comparator is configured.
func New[T any](opts ...Option[T]) *Heap[T] {
	config := defaultConfig[T]()
	for _, opt := range opts {
		opt(config)
	}
	if config.comparator == nil {
		panic("comparator must be provided or use NewOrdered, Min, or Max")
	}
	return &Heap[T]{
		elements:   make([]T, 0, config.capacity),
		comparator: config.comparator,
	}
}

// NewOrdered creates a new min-heap for ordered types using natural ordering.
func NewOrdered[T cmp.Ordered](opts ...Option[T]) *Heap[T] {
	return New[T](append([]Option[T]{WithComparator(comparator.NaturalOrder[T]())}, opts...)...)
}

// Min creates a new Min-Heap using natural ordering.
func Min[T cmp.Ordered]() *Heap[T] {
	return NewOrdered[T]()
}

// Max creates a new Max-Heap using reversed natural ordering.
func Max[T cmp.Ordered]() *Heap[T] {
	return New[T](WithComparator(comparator.Reverse(comparator.NaturalOrder[T]())))
}

// Add inserts the specified element into the heap.
func (h *Heap[T]) Add(t T) {
	h.elements = append(h.elements, t)
	h.siftUp(len(h.elements) - 1)
}

// AddAll inserts all elements from the given sequence into the heap.
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

// Size returns the number of elements in the heap.
func (h *Heap[T]) Size() int {
	return len(h.elements)
}

// Empty returns true if the heap contains no elements.
func (h *Heap[T]) Empty() bool {
	return len(h.elements) == 0
}

// Clear removes all elements from the heap.
func (h *Heap[T]) Clear() {
	// Zero out elements to allow GC
	var zero T
	for i := range h.elements {
		h.elements[i] = zero
	}
	h.elements = h.elements[:0]
}

// All returns an iterator over all elements in the heap in internal slice order.
func (h *Heap[T]) All() iter.Seq[T] {
	return slices.Values(h.elements)
}

// String returns a string representation of the heap elements matching Go slice formatting.
func (h *Heap[T]) String() string {
	vals := make([]string, 0, h.Size())
	for item := range h.All() {
		vals = append(vals, fmt.Sprintf("%v", item))
	}
	return "[" + strings.Join(vals, " ") + "]"
}

// Private Functions

func defaultConfig[T any]() *config[T] {
	return &config[T]{
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
