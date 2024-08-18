package hashset

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"maps"
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

func (s *HashSet[T]) AddAll(other collections.Collection[T]) {
	for t := range other.All() {
		s.Add(t)
	}
}

func (s *HashSet[T]) RemoveAll(other collections.Collection[T]) {
	for t := range other.All() {
		s.Remove(t)
	}
}

func (s *HashSet[T]) RetainAll(other collections.Collection[T]) {
	for t := range other.All() {
		if !s.Contains(t) {
			s.Remove(t)
		}
	}
}

func (s *HashSet[T]) Size() int {
	return len(s.m)
}

func (s *HashSet[T]) Empty() bool {
	return s.Size() == 0
}

func (s *HashSet[T]) String() string {
	vals := make([]string, s.Size())
	i := 0
	for item := range maps.Keys(s.m) {
		vals[i] = fmt.Sprintf("%+v", item)
		i++
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

func (s *HashSet[T]) All() iter.Seq[T] {
	return maps.Keys(s.m)
}

func (s *HashSet[T]) ToSlice() []T {
	slice := make([]T, 0, s.Size())
	for item := range s.All() {
		slice = append(slice, item)
	}
	return slice
}

func defaultConfig() *Config {
	return &Config{}
}
