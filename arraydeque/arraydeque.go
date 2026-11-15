// Package arraydeque provides an array-backed double-ended queue.
package arraydeque

import (
	"fmt"
	"iter"
	"strings"

	"github.com/lock14/collections"
)

var (
	_ collections.MutableDeque[int] = (*ArrayDeque[int])(nil)
	_ fmt.Stringer                  = (*ArrayDeque[int])(nil)
)

const (
	// DefaultCapacity is the capacity assigned if no other is provided.
	DefaultCapacity = 16
	// doublingThreshold is the capacity threshold below which an ArrayDeque
	// doubles its capacity when resized.
	doublingThreshold = 512
	// fiftyPercentThreshold is the capacity threshold below which (and at or above
	// doublingThreshold) an ArrayDeque increases its capacity by 50% when resized.
	// At or above this threshold, capacity increases by 25%.
	fiftyPercentThreshold = 2048
)

// config holds the values for configuring an ArrayDeque.
type config struct {
	capacity int
}

// Option configures an ArrayDeque.
type Option func(*config)

// WithCapacity configures the initial capacity of the ArrayDeque.
// Passing a negative capacity causes New to panic.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

// ArrayDeque represents a double-ended queue of elements of type T backed by a circular slice.
// The zero value for ArrayDeque is an empty deque ready to use.
type ArrayDeque[T any] struct {
	slice []T
	front int
	back  int
	size  int
}

// New creates an empty ArrayDeque whose initial size is 0.
// Panics if a negative capacity is provided via WithCapacity.
func New[T any](opts ...Option) *ArrayDeque[T] {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.capacity < 0 {
		panic("arraydeque: capacity cannot be negative")
	}
	return &ArrayDeque[T]{
		slice: make([]T, cfg.capacity),
	}
}

// Add inserts the given element at the back of this deque.
// It is an alias for AddBack.
func (d *ArrayDeque[T]) Add(t T) {
	d.AddBack(t)
}

// AddFront inserts the given element at the front of this deque.
func (d *ArrayDeque[T]) AddFront(t T) {
	if d.size == len(d.slice) {
		d.resize()
	}
	d.front--
	if d.front == -1 {
		d.front = len(d.slice) - 1
	}
	d.slice[d.front] = t
	d.size++
}

// AddBack inserts the given element at the back of this deque.
func (d *ArrayDeque[T]) AddBack(t T) {
	if d.size == len(d.slice) {
		d.resize()
	}
	d.slice[d.back] = t
	d.back++
	if d.back == len(d.slice) {
		d.back = 0
	}
	d.size++
}

// Push inserts the given element at the front of this deque.
// It is an alias for AddFront.
func (d *ArrayDeque[T]) Push(t T) {
	d.AddFront(t)
}

// AddAll appends all elements from the given sequence to the back of this deque.
func (d *ArrayDeque[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		d.Add(t)
	}
}

// Remove removes and returns the element at the front of this deque.
// It is an alias for RemoveFront. Panics if this deque is empty.
func (d *ArrayDeque[T]) Remove() T {
	return d.RemoveFront()
}

// RemoveFront removes and returns the element at the front of this deque.
// Panics if this deque is empty.
func (d *ArrayDeque[T]) RemoveFront() T {
	if d.Empty() {
		panic("arraydeque: cannot remove from an empty deque")
	}
	var zero T
	t := d.slice[d.front]
	d.slice[d.front] = zero
	d.front++
	if d.front == len(d.slice) {
		d.front = 0
	}
	d.size--
	return t
}

// RemoveBack removes and returns the element at the back of this deque.
// Panics if this deque is empty.
func (d *ArrayDeque[T]) RemoveBack() T {
	if d.Empty() {
		panic("arraydeque: cannot remove from an empty deque")
	}
	var zero T
	d.back--
	if d.back == -1 {
		d.back = len(d.slice) - 1
	}
	t := d.slice[d.back]
	d.slice[d.back] = zero
	d.size--
	return t
}

// Pop removes and returns the element at the front of this deque.
// It is an alias for RemoveFront. Panics if this deque is empty.
func (d *ArrayDeque[T]) Pop() T {
	return d.RemoveFront()
}

// Clear removes all elements from this deque and releases backing memory for garbage collection.
func (d *ArrayDeque[T]) Clear() {
	clear(d.slice)
	d.slice = nil
	d.front = 0
	d.back = 0
	d.size = 0
}

// Peek returns the element at the front of this deque without removing it.
// It is an alias for PeekFront. Panics if this deque is empty.
func (d *ArrayDeque[T]) Peek() T {
	return d.PeekFront()
}

// PeekFront returns the element at the front of this deque without removing it.
// Panics if this deque is empty.
func (d *ArrayDeque[T]) PeekFront() T {
	if d.Empty() {
		panic("arraydeque: cannot peek from an empty deque")
	}
	return d.slice[d.front]
}

// PeekBack returns the element at the back of this deque without removing it.
// Panics if this deque is empty.
func (d *ArrayDeque[T]) PeekBack() T {
	if d.Empty() {
		panic("arraydeque: cannot peek from an empty deque")
	}
	i := d.back - 1
	if i < 0 {
		i = len(d.slice) - 1
	}
	return d.slice[i]
}

// Size returns the number of elements in this deque.
func (d *ArrayDeque[T]) Size() int {
	return d.size
}

// Empty returns true if this deque contains no elements.
func (d *ArrayDeque[T]) Empty() bool {
	return d.size == 0
}

// All returns an iterator over all elements in this deque from front to back.
func (d *ArrayDeque[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		i := d.front
		for count < d.size {
			if !yield(d.slice[i]) {
				return
			}
			i++
			if i == len(d.slice) {
				i = 0
			}
			count++
		}
	}
}

// Backward returns an iterator over all elements in this deque from back to front.
func (d *ArrayDeque[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		i := d.back - 1
		if i < 0 {
			i = len(d.slice) - 1
		}
		for count < d.size {
			if !yield(d.slice[i]) {
				return
			}
			i--
			if i < 0 {
				i = len(d.slice) - 1
			}
			count++
		}
	}
}

// String returns a string representation of this deque matching Go slice formatting.
func (d *ArrayDeque[T]) String() string {
	if d.Empty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for v := range d.All() {
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte(']')
	return sb.String()
}

func defaultConfig() *config {
	return &config{
		capacity: DefaultCapacity,
	}
}

func (d *ArrayDeque[T]) resize() {
	var newCap int
	switch {
	case len(d.slice) == 0:
		newCap = DefaultCapacity
	case d.size < doublingThreshold:
		newCap = len(d.slice) << 1
	case d.size < fiftyPercentThreshold:
		newCap = len(d.slice) + (len(d.slice) >> 1)
	default:
		newCap = len(d.slice) + (len(d.slice) >> 2)
	}
	s := make([]T, newCap)
	m := copy(s, d.slice[d.front:])
	copy(s[m:], d.slice[:d.front])
	d.slice = s
	d.front = 0
	d.back = d.size
}
