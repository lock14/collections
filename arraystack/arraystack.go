package arraystack

import (
	"fmt"
	"strings"
)

// ArrayStack represents a stack of elements of type T backed by an array.
type ArrayStack[T any] struct {
	slice []T
}

// Config holds the values for configuring a ArrayStack.
type Config struct{}

// Option configures a ArrayStack config
type Option func(*Config)

// New creates a empty ArrayStack whose initial size is 0.
func New[T any](opts ...Option) *ArrayStack[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &ArrayStack[T]{
		slice: make([]T, 0),
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
	str := make([]string, len(s.slice))
	for i := 0; i < len(str); i++ {
		str[i] = fmt.Sprintf("%+v", s.slice[i])
	}
	return "[" + strings.Join(str, ", ") + "]"
}

func (s *ArrayStack[T]) ToSlice() []T {
	slice := make([]T, s.Size())
	copy(slice, s.slice)
	return slice
}

func defaultConfig() *Config {
	return &Config{}
}
