package set

import (
	"fmt"
	"strings"
)

// Set[T] represents a set of elements of type T
type Set[T comparable] struct {
	m map[T]struct{}
}

// Config holds the values for configuring a Set.
type Config struct{}

// Option configures a Set config
type Option func(*Config)

// New creates a empty Set.
func New[T comparable](opts ...Option) *Set[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &Set[T]{
		m: make(map[T]struct{}, 0),
	}
}

func (s *Set[T]) Add(item T) {
	s.m[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
	delete(s.m, item)
}

func (s *Set[T]) Contains(item T) bool {
	_, present := s.m[item]
	return present
}

func (s *Set[T]) Size() int {
	return len(s.m)
}

func (s *Set[T]) isEmpty() bool {
	return s.Size() == 0
}

func (s *Set[T]) String() string {
	vals := make([]string, s.Size())
	i := 0
	for item := range s.m {
		vals[i] = fmt.Sprintf("%+v", item)
		i++
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

func (s *Set[T]) ToSlice() []T {
	slice := make([]T, s.Size())
	i := 0
	for item := range s.m {
		slice[i] = item
		i++
	}
	return slice
}

func defaultConfig() *Config {
	return &Config{}
}
