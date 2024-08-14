package arraylist

import (
	"fmt"
	"github.com/lock14/collections/iterator"
	"strings"
)

// public constants
const (
	DefaultCapacity = 10
)

// public types/structs

type Config struct {
	capacity int
}

type Option func(*Config)

type ArrayList[T any] struct {
	slice []T
}

type listIterator[T any] struct {
	slice []T
	index int
}

// public functions/recievers

func Capacity(n int) Option {
	return func(config *Config) {
		config.capacity = n
	}
}

func New[T any](opts ...Option) *ArrayList[T] {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return &ArrayList[T]{
		slice: make([]T, 0, config.capacity),
	}
}

func From[T any](slice []T) *ArrayList[T] {
	return &ArrayList[T]{
		slice: slice,
	}
}

func (l *ArrayList[T]) Add(items ...T) {
	l.slice = append(l.slice, items...)
}

func (l *ArrayList[T]) Size() int {
	return len(l.slice)
}

func (l *ArrayList[T]) Empty() bool {
	return l.Size() == 0
}

func (l *ArrayList[T]) Get(index int) T {
	return l.slice[index]
}

func (l *ArrayList[T]) Set(index int, item T) {
	l.slice[index] = item
}

func (l *ArrayList[T]) String() string {
	vals := make([]string, l.Size())
	for i := 0; i < len(l.slice); i++ {
		vals[i] = fmt.Sprintf("%+v", l.slice[i])
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

func (l *ArrayList[T]) Iterator() iterator.Iterator[T] {
	return &listIterator[T]{
		slice: l.slice,
		index: 0,
	}
}

func (l *ArrayList[T]) Elements() chan T {
	return iterator.Elements(l.Iterator())
}

func (itr *listIterator[T]) Empty() bool {
	return itr.index >= len(itr.slice)
}

func (itr *listIterator[T]) Next() (T, error) {
	var t T
	if itr.Empty() {
		return t, fmt.Errorf("listIterator is empty")
	}
	t = itr.slice[itr.index]
	itr.index++
	return t, nil
}

// private functions/receivers

func defaultConfig() *Config {
	return &Config{
		capacity: DefaultCapacity,
	}
}
