package arraydeque

import (
	"fmt"
	"iter"
	"strings"
)

const (
	// DefaultCapacity is the capacity assigned if no other is provided.
	DefaultCapacity = 1
	// if an arraydeque's capacity is under this amount its capacity
	// will double when it needs to be resized.
	doublingThreshold = 512
	// if an arraydeque's capacity is under this amount its capacity
	// will increase by 50% when it needs to be resized.
	fiftyPercentThreshold = 2048
)

// ArrayDeque represents a deque of elements of type T backed by an array.
// The zero value for ArrayDeque is an empty deque ready to use.
type ArrayDeque[T any] struct {
	slice []T
	front int
	back  int
	size  int
}

// Config holds the values for configuring a ArrayDeque.
type Config struct {
	Capacity int
}

// Option configures a ArrayDeque config
type Option func(*Config)

// New creates an empty ArrayDeque whose initial size is 0.
func New[T any](opts ...Option) *ArrayDeque[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &ArrayDeque[T]{
		slice: make([]T, config.Capacity),
	}
}

// Peek is an alias for PeekFront
func (d *ArrayDeque[T]) Peek() T {
	return d.PeekFront()
}

// Add is an alias for AddBack
func (d *ArrayDeque[T]) Add(t T) {
	d.AddBack(t)
}

// Remove is an alias for RemoveFront
func (d *ArrayDeque[T]) Remove() T {
	return d.RemoveFront()
}

// Push is an alias for AddFront
func (d *ArrayDeque[T]) Push(t T) {
	d.AddFront(t)
}

// Pop is an alias for RemoveFront
func (d *ArrayDeque[T]) Pop() T {
	return d.RemoveFront()
}

// PeekFront returns the element at the front of this deque.
// If this deque is empty, PeekFront panics.
func (d *ArrayDeque[T]) PeekFront() T {
	if d.Empty() {
		panic("cannot peek from an empty Deque")
	}
	return d.slice[d.front]
}

// AddFront adds the given element to the front of this deque.
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

// RemoveFront removes the given element at the front of this deque
// and returns it. If this deque is empty, RemoveFront panics.
func (d *ArrayDeque[T]) RemoveFront() T {
	if d.Empty() {
		panic("cannot remove from an empty ArrayDeque")
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

// PeekBack returns the element at the back of this deque.
// If this deque is empty, PeekBack panics.
func (d *ArrayDeque[T]) PeekBack() T {
	if d.Empty() {
		panic("cannot peek from an empty Deque")
	}
	i := d.back - 1
	if i < 0 {
		i = len(d.slice) - 1
	}
	return d.slice[i]
}

// AddBack adds the given element to the back of this deque.
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

// RemoveBack removes the given element at the back of this deque
// and returns it. If this deque is empty, RemoveBack panics.
func (d *ArrayDeque[T]) RemoveBack() T {
	if d.Empty() {
		panic("cannot remove from an empty ArrayDeque")
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

// Size returns the number of elements in this deque.
func (d *ArrayDeque[T]) Size() int {
	return d.size
}

// Empty returns true if this deque contains no elements.
// Otherwise, returns false.
func (d *ArrayDeque[T]) Empty() bool {
	return d.size == 0
}

// Clear removes all elements from this deque and
// releases any memory in use by this deque for garbage
// collection.
func (d *ArrayDeque[T]) Clear() {
	d.slice = nil
	d.front = 0
	d.back = 0
	d.size = 0
}

// String returns a string representation of this deque.
func (d *ArrayDeque[T]) String() string {
	str := make([]string, 0, d.Size())
	for t := range d.All() {
		str = append(str, fmt.Sprintf("%+v", t))
	}
	return "[" + strings.Join(str, ", ") + "]"
}

// All returns an iterator over all elements,
// going from front to back in this deque.
func (d *ArrayDeque[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		i := d.front
		for count < d.size {
			if !yield(d.slice[i]) {
				break
			}
			i++
			if i == len(d.slice) {
				i = 0
			}
			count++
		}
	}
}

// Backward returns an iterator over all elements,
// going from back to front in this deque.
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

func (d *ArrayDeque[T]) resize() {
	var newCap int
	if d.slice == nil {
		newCap = DefaultCapacity
	} else if d.size < doublingThreshold {
		newCap = len(d.slice) << 1
	} else if d.size < fiftyPercentThreshold {
		newCap = len(d.slice)
		newCap += len(d.slice) >> 1
	} else { // grow by 25%
		newCap = len(d.slice)
		newCap += len(d.slice) >> 2
	}
	s := make([]T, newCap)
	m := copy(s, d.slice[d.front:])
	n := copy(s[m:], d.slice[0:d.front])
	if m+n != d.size {
		panic("resize algorithm incorrect")
	}
	d.slice = s
	d.front = 0
	d.back = d.size
}

func defaultConfig() *Config {
	return &Config{
		Capacity: DefaultCapacity,
	}
}
