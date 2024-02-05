package hashset

import (
	"fmt"
	"strings"
)

// HashSet represents a set of elements of type T.
type HashSet[T comparable] struct {
	m map[T]struct{}
}

// Config holds the values for configuring a HashSet.
type Config struct{}

// Option configures a HashSet config
type Option func(*Config)

// New creates an empty HashSet.
func New[T comparable](opts ...Option) *HashSet[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &HashSet[T]{
		m: make(map[T]struct{}),
	}
}

func (s *HashSet[T]) Add(item T) {
	s.m[item] = struct{}{}
}

func (s *HashSet[T]) Remove(item T) {
	delete(s.m, item)
}

func (s *HashSet[T]) Contains(item T) bool {
	_, present := s.m[item]
	return present
}

func (s *HashSet[T]) Size() int {
	return len(s.m)
}

func (s *HashSet[T]) isEmpty() bool {
	return s.Size() == 0
}

func (s *HashSet[T]) String() string {
	vals := make([]string, s.Size())
	i := 0
	for item := range s.m {
		vals[i] = fmt.Sprintf("%+v", item)
		i++
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

func (s *HashSet[T]) ToSlice() []T {
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
