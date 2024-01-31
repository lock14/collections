package stack

import (
	"fmt"
	"strings"
)

// Stack[T] represents a stack of elements of type T
type Stack[T any] struct {
	slice []T
}

// Config holds the values for configuring a Stack.
type Config struct{}

// Option configures a Stack config
type Option func(*Config)

// New creates a empty Stack whose initial size is 0.
func New[T any](opts ...Option) *Stack[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &Stack[T]{
		slice: make([]T, 0),
	}
}

func (s *Stack[T]) Push(t T) {
	s.slice = append(s.slice, t)
}

func (s *Stack[T]) Pop() T {
	if len(s.slice) == 0 {
		panic("Cannot Pop an empty Stack")
	}
	t := s.slice[len(s.slice)-1]
	s.slice = s.slice[:len(s.slice)-1]
	return t
}

func (s *Stack[T]) Size() int {
	return len(s.slice)
}

func (s *Stack[T]) isEmpty() bool {
	return s.Size() == 0
}

func (s *Stack[T]) String() string {
	str := make([]string, len(s.slice))
	for i := 0; i < len(str); i++ {
		str[i] = fmt.Sprintf("%+v", s.slice[i])
	}
	return "[" + strings.Join(str, ", ") + "]"
}

func (s *Stack[T]) ToSlice() []T {
	slice := make([]T, s.Size())
	copy(slice, s.slice)
	return slice
}

func defaultConfig() *Config {
	return &Config{}
}
