package arraylist

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"slices"
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

func (l *ArrayList[T]) Add(t T) {
	l.slice = append(l.slice, t)
}

func (l *ArrayList[T]) Remove() T {
	if l.Empty() {
		panic("cannot remove from an empty list")
	}
	t := l.slice[l.Size()-1]
	l.slice = l.slice[0 : l.Size()-1]
	return t
}

func (l *ArrayList[T]) Push(t T) {
	l.Add(t)
}

func (l *ArrayList[T]) Pop() T {
	return l.Remove()
}

func (l *ArrayList[T]) AddAll(other collections.Collection[T]) {
	for t := range other.All() {
		l.Add(t)
	}
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

func (l *ArrayList[T]) All() iter.Seq[T] {
	return slices.Values(l.slice[0:l.Size()])
}

// private functions/receivers

func defaultConfig() *Config {
	return &Config{
		capacity: DefaultCapacity,
	}
}
