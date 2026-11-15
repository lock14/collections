// Package heap provides a binary heap priority queue.
package heap

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/lock14/collections"
	"github.com/lock14/collections/comparator"
)

var (
	_ collections.MutableQueue[int] = (*Heap[int])(nil)
	_ fmt.Stringer                  = (*Heap[int])(nil)
)

const (
	// DefaultCapacity is the default initial capacity for a Heap.
	DefaultCapacity = 10
)

// config holds the configuration for a Heap.
type config[T any] struct {
	capacity   int
	comparator comparator.Comparator[T]
}

// Option is a function that configures a Heap.
type Option[T any] func(*config[T])

// WithComparator configures the comparator used by the Heap.
// Passing a nil comparator causes New to panic.
func WithComparator[T any](cmpFunc comparator.Comparator[T]) Option[T] {
	return func(cfg *config[T]) {
		cfg.comparator = cmpFunc
	}
}

// WithCapacity configures the initial pre-allocated capacity of the Heap.
// Passing a negative capacity causes New to panic.
func WithCapacity[T any](capacity int) Option[T] {
	return func(cfg *config[T]) {
		cfg.capacity = capacity
	}
}

// Heap is a binary heap priority queue.
type Heap[T any] struct {
	elements   []T
	comparator comparator.Comparator[T]
}

// New creates a new Heap with the given options.
// Panics if no comparator (or a nil comparator) is configured, or if capacity is negative.
func New[T any](opts ...Option[T]) *Heap[T] {
	cfg := defaultConfig[T]()
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.capacity < 0 {
		panic("heap: capacity cannot be negative")
	}
	if cfg.comparator == nil {
		panic("heap: comparator must be provided or use NewOrdered, Min, or Max")
	}
	return &Heap[T]{
		elements:   make([]T, 0, cfg.capacity),
		comparator: cfg.comparator,
	}
}

// NewOrdered creates a new min-heap for ordered types using natural ordering and optional configuration.
// Panics if a subsequent option sets a nil comparator or negative capacity.
func NewOrdered[T cmp.Ordered](opts ...Option[T]) *Heap[T] {
	allOpts := make([]Option[T], 0, len(opts)+1)
	allOpts = append(allOpts, WithComparator(comparator.NaturalOrder[T]()))
	allOpts = append(allOpts, opts...)
	return New[T](allOpts...)
}

// Min creates a new Min-Heap using natural ordering and optional configuration.
// Panics if a subsequent option sets a nil comparator or negative capacity.
func Min[T cmp.Ordered](opts ...Option[T]) *Heap[T] {
	return NewOrdered[T](opts...)
}

// Max creates a new Max-Heap using reversed natural ordering and optional configuration.
// Panics if a subsequent option sets a nil comparator or negative capacity.
func Max[T cmp.Ordered](opts ...Option[T]) *Heap[T] {
	allOpts := make([]Option[T], 0, len(opts)+1)
	allOpts = append(allOpts, WithComparator(comparator.Reverse(comparator.NaturalOrder[T]())))
	allOpts = append(allOpts, opts...)
	return New[T](allOpts...)
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
		panic("heap: cannot remove from an empty heap")
	}
	return h.popRoot()
}

// Clear removes all elements from the heap while retaining its capacity.
func (h *Heap[T]) Clear() {
	clear(h.elements)
	h.elements = h.elements[:0]
}

// Peek returns the top element from the heap without removing it.
// Panics if the heap is empty.
func (h *Heap[T]) Peek() T {
	if h.Empty() {
		panic("heap: cannot peek from an empty heap")
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

// All returns an iterator over all elements in the heap in internal slice order.
func (h *Heap[T]) All() iter.Seq[T] {
	return slices.Values(h.elements)
}

// String returns a string representation of the heap elements matching Go slice formatting.
func (h *Heap[T]) String() string {
	if h.Empty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range h.elements {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte(']')
	return sb.String()
}

func defaultConfig[T any]() *config[T] {
	return &config[T]{
		capacity: DefaultCapacity,
	}
}

func (h *Heap[T]) popRoot() T {
	root := h.elements[0]
	last := len(h.elements) - 1
	h.elements[0] = h.elements[last]
	var zero T
	h.elements[last] = zero
	h.elements = h.elements[:last]
	if last > 0 {
		h.siftDown(0)
	}
	return root
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
