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
type SliceWrapper[T any] struct {
	slice []T
}

// Wrap creates a new SliceWrapper around the given slice.
func Wrap[T any](slice []T) *SliceWrapper[T] {
	return &SliceWrapper[T]{
		slice: slice,
	}
}

func (l *SliceWrapper[T]) Add(t T) {
	l.slice = append(l.slice, t)
}

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

func (l *SliceWrapper[T]) Push(t T) {
	l.Add(t)
}

func (l *SliceWrapper[T]) Pop() T {
	return l.Remove()
}

func (l *SliceWrapper[T]) Peek() T {
	if len(l.slice) == 0 {
		panic("cannot peek from an empty list")
	}
	return l.slice[len(l.slice)-1]
}

func (l *SliceWrapper[T]) Clear() {
	// Zero out to avoid memory leaks
	var zero T
	for i := range l.slice {
		l.slice[i] = zero
	}
	l.slice = l.slice[:0]
}

func (l *SliceWrapper[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		l.Add(t)
	}
}

func (l *SliceWrapper[T]) Size() int {
	return len(l.slice)
}

func (l *SliceWrapper[T]) Empty() bool {
	return len(l.slice) == 0
}

func (l *SliceWrapper[T]) Get(index int) T {
	return l.slice[index]
}

func (l *SliceWrapper[T]) Set(index int, item T) {
	l.slice[index] = item
}

func (l *SliceWrapper[T]) String() string {
	vals := make([]string, 0, l.Size())
	for v := range l.All() {
		vals = append(vals, fmt.Sprintf("%+v", v))
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

func (l *SliceWrapper[T]) All() iter.Seq[T] {
	return slices.Values(l.slice[0:l.Size()])
}
