package arraylist

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"slices"
	"strings"
)

var _ collections.MutableList[int] = (*SliceWrapper[int])(nil)

type SliceWrapper[T any] struct {
	slice []T
}

func Wrap[T any](slice []T) *SliceWrapper[T] {
	return &SliceWrapper[T]{
		slice: slice,
	}
}

func (l *SliceWrapper[T]) Add(t T) {
	l.slice = append(l.slice, t)
}

func (l *SliceWrapper[T]) Remove() T {
	if l.Empty() {
		panic("cannot remove from an empty list")
	}
	t := l.slice[l.Size()-1]
	l.slice = l.slice[0 : l.Size()-1]
	return t
}

func (l *SliceWrapper[T]) Push(t T) {
	l.Add(t)
}

func (l *SliceWrapper[T]) Pop() T {
	return l.Remove()
}

func (l *SliceWrapper[T]) AddAll(other collections.Collection[T]) {
	for t := range other.All() {
		l.Add(t)
	}
}

func (l *SliceWrapper[T]) Size() int {
	return len(l.slice)
}

func (l *SliceWrapper[T]) Empty() bool {
	return l.Size() == 0
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
