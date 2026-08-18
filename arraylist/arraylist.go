// Package arraylist provides an array-backed list and stack wrapper.
package arraylist

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"slices"
	"strings"
)

var (
	_ collections.MutableList[int]  = (*SliceWrapper[int])(nil)
	_ collections.MutableStack[int] = (*SliceWrapper[int])(nil)
)

// SliceWrapper is a wrapper around the built-in slice that implements collections.MutableList and collections.MutableStack.
// The zero value for SliceWrapper is an empty list ready to use.
type SliceWrapper[T any] struct {
	slice []T
}

// config holds the configuration options for a SliceWrapper.
type config struct {
	capacity int
}

// Option configures a SliceWrapper.
type Option func(*config)

// WithCapacity configures the initial capacity of the SliceWrapper.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

// New creates a new empty SliceWrapper.
func New[T any](opts ...Option) *SliceWrapper[T] {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return &SliceWrapper[T]{
		slice: make([]T, 0, cfg.capacity),
	}
}

// Wrap creates a new SliceWrapper around the given slice.
func Wrap[T any](slice []T) *SliceWrapper[T] {
	return &SliceWrapper[T]{
		slice: slice,
	}
}

// Add appends the given element to the end of the list.
func (l *SliceWrapper[T]) Add(t T) {
	l.slice = append(l.slice, t)
}

// Remove removes and returns the last element from the list.
// If the list is empty, Remove panics.
func (l *SliceWrapper[T]) Remove() T {
	if len(l.slice) == 0 {
		panic("cannot remove from an empty list")
	}
	idx := len(l.slice) - 1
	t := l.slice[idx]
	var zero T
	l.slice[idx] = zero // avoid memory leak
	l.slice = l.slice[:idx]
	return t
}

// Push adds the given element to the top of the stack (end of the list).
func (l *SliceWrapper[T]) Push(t T) {
	l.Add(t)
}

// Pop removes and returns the element at the top of the stack (end of the list).
// If the stack is empty, Pop panics.
func (l *SliceWrapper[T]) Pop() T {
	return l.Remove()
}

// Peek returns the element at the top of the stack (end of the list) without removing it.
// If the stack is empty, Peek panics.
func (l *SliceWrapper[T]) Peek() T {
	if len(l.slice) == 0 {
		panic("cannot peek from an empty list")
	}
	return l.slice[len(l.slice)-1]
}

// Clear removes all elements from the list while retaining its capacity.
func (l *SliceWrapper[T]) Clear() {
	// Zero out to avoid memory leaks
	var zero T
	for i := range l.slice {
		l.slice[i] = zero
	}
	l.slice = l.slice[:0]
}

// AddAll appends all elements from the given sequence to the end of the list.
func (l *SliceWrapper[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		l.Add(t)
	}
}

// Size returns the number of elements in the list.
func (l *SliceWrapper[T]) Size() int {
	return len(l.slice)
}

// Empty returns true if the list contains no elements.
func (l *SliceWrapper[T]) Empty() bool {
	return len(l.slice) == 0
}

// Get returns the element at the specified index.
func (l *SliceWrapper[T]) Get(index int) T {
	return l.slice[index]
}

// Set replaces the element at the specified index with the given item.
func (l *SliceWrapper[T]) Set(index int, item T) {
	l.slice[index] = item
}

// String returns a string representation of the list.
func (l *SliceWrapper[T]) String() string {
	vals := make([]string, 0, l.Size())
	for v := range l.All() {
		vals = append(vals, fmt.Sprintf("%+v", v))
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

// All returns an iterator over all elements in the list from first to last.
func (l *SliceWrapper[T]) All() iter.Seq[T] {
	return slices.Values(l.slice)
}
