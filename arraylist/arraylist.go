// Package arraylist provides an array-backed list and stack wrapper.
package arraylist

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/lock14/collections"
)

var (
	_ collections.MutableList[int]  = (*SliceWrapper[int])(nil)
	_ collections.MutableStack[int] = (*SliceWrapper[int])(nil)
	_ fmt.Stringer                  = (*SliceWrapper[int])(nil)
)

// config holds the configuration options for a SliceWrapper.
type config struct {
	capacity int
}

// Option configures a SliceWrapper.
type Option func(*config)

// WithCapacity configures the initial capacity of the SliceWrapper.
// Passing a negative capacity causes New to panic.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

// SliceWrapper is a wrapper around the built-in slice that implements collections.MutableList and collections.MutableStack.
// The zero value for SliceWrapper is an empty list ready to use.
type SliceWrapper[T any] struct {
	slice []T
}

// New creates a new empty SliceWrapper with optional configuration.
// Panics if a negative capacity is provided via WithCapacity.
func New[T any](opts ...Option) *SliceWrapper[T] {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.capacity < 0 {
		panic("arraylist: capacity cannot be negative")
	}
	return &SliceWrapper[T]{
		slice: make([]T, 0, cfg.capacity),
	}
}

// Wrap creates a new SliceWrapper around the given slice.
// If slice is nil, the returned SliceWrapper is empty and ready to use.
func Wrap[T any](slice []T) *SliceWrapper[T] {
	return &SliceWrapper[T]{
		slice: slice,
	}
}

// Add appends the given element to the end of the list.
func (l *SliceWrapper[T]) Add(t T) {
	l.slice = append(l.slice, t)
}

// Push adds the given element to the top of the stack (end of the list).
func (l *SliceWrapper[T]) Push(t T) {
	l.Add(t)
}

// AddAll appends all elements from the given sequence to the end of the list.
func (l *SliceWrapper[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		l.Add(t)
	}
}

// Set replaces the element at the specified index with the given item.
// Panics if index is out of bounds (index < 0 || index >= Size()).
func (l *SliceWrapper[T]) Set(index int, item T) {
	l.slice[index] = item
}

// Remove removes and returns the last element from the list.
// Panics if the list is empty.
func (l *SliceWrapper[T]) Remove() T {
	if len(l.slice) == 0 {
		panic("arraylist: cannot remove from an empty list")
	}
	idx := len(l.slice) - 1
	t := l.slice[idx]
	var zero T
	l.slice[idx] = zero
	l.slice = l.slice[:idx]
	return t
}

// Pop removes and returns the element at the top of the stack (end of the list).
// Panics if the stack is empty.
func (l *SliceWrapper[T]) Pop() T {
	return l.Remove()
}

// Clear removes all elements from the list while retaining its capacity.
func (l *SliceWrapper[T]) Clear() {
	clear(l.slice)
	l.slice = l.slice[:0]
}

// Get returns the element at the specified index.
// Panics if index is out of bounds (index < 0 || index >= Size()).
func (l *SliceWrapper[T]) Get(index int) T {
	return l.slice[index]
}

// Peek returns the element at the top of the stack (end of the list) without removing it.
// Panics if the stack is empty.
func (l *SliceWrapper[T]) Peek() T {
	if len(l.slice) == 0 {
		panic("arraylist: cannot peek from an empty list")
	}
	return l.slice[len(l.slice)-1]
}

// Size returns the number of elements in the list.
func (l *SliceWrapper[T]) Size() int {
	return len(l.slice)
}

// Empty returns true if the list contains no elements.
func (l *SliceWrapper[T]) Empty() bool {
	return len(l.slice) == 0
}

// All returns an iterator over all elements in the list from first to last.
func (l *SliceWrapper[T]) All() iter.Seq[T] {
	return slices.Values(l.slice)
}

// String returns a string representation of the list matching Go slice formatting.
func (l *SliceWrapper[T]) String() string {
	if l.Empty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range l.slice {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte(']')
	return sb.String()
}
