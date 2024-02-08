package arraystack

import (
	"fmt"
	"github.com/lock14/collections/iterator"
	"strings"
)

const (
	DefaultCapacity = 10
)

// ArrayStack represents a stack of elements of type T backed by an array.
type ArrayStack[T any] struct {
	slice []T
}

// Config holds the values for configuring a ArrayStack.
type Config struct {
	Capacity int
}

// Option configures a ArrayStack config
type Option func(*Config)

// New creates a empty ArrayStack whose initial size is 0.
func New[T any](opts ...Option) *ArrayStack[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &ArrayStack[T]{
		slice: make([]T, 0, config.Capacity),
	}
}

func (s *ArrayStack[T]) Push(t T) {
	s.slice = append(s.slice, t)
}

func (s *ArrayStack[T]) Pop() T {
	if len(s.slice) == 0 {
		panic("Cannot Pop an empty ArrayStack")
	}
	t := s.slice[len(s.slice)-1]
	s.slice = s.slice[:len(s.slice)-1]
	return t
}

func (s *ArrayStack[T]) Size() int {
	return len(s.slice)
}

func (s *ArrayStack[T]) isEmpty() bool {
	return s.Size() == 0
}

func (s *ArrayStack[T]) String() string {
	str := make([]string, 0, len(s.slice))
	for t := range s.Elements() {
		str = append(str, fmt.Sprintf("%+v", *t))
	}
	return "[" + strings.Join(str, ", ") + "]"
}

func (s *ArrayStack[T]) Iterator() iterator.Iterator[T] {
	return &stackIterator[T]{
		stack: s,
		index: len(s.slice) - 1,
	}
}

func (s *ArrayStack[T]) Elements() chan *T {
	return iterator.Elements(s.Iterator())
}

func (s *ArrayStack[T]) ToSlice() []T {
	slice := make([]T, s.Size())
	copy(slice, s.slice)
	return slice
}

func defaultConfig() *Config {
	return &Config{
		Capacity: DefaultCapacity,
	}
}

// Iterator

type stackIterator[T any] struct {
	stack *ArrayStack[T]
	index int
}

func (itr *stackIterator[T]) Empty() bool {
	return itr.index < 0
}

func (itr *stackIterator[T]) Next() (*T, error) {
	if itr.Empty() {
		return nil, fmt.Errorf("cannot call Next() on an empty Iterator")
	}
	t := &itr.stack.slice[itr.index]
	itr.index--
	return t, nil
}
