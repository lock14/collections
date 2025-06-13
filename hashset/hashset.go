package hashset

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"maps"
	"strings"
)

var _ collections.MutableSet[int] = (*HashSet[int])(nil)

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

func (s *HashSet[T]) Remove() T {
	var t T
	for item := range s.m {
		t = item
		break
	}
	return t
}

func (s *HashSet[T]) RemoveElement(item T) {
	delete(s.m, item)
}

func (s *HashSet[T]) Contains(item T) bool {
	_, present := s.m[item]
	return present
}

func (s *HashSet[T]) ContainsAll(other collections.Collection[T]) bool {
	for item := range other.All() {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

func (s *HashSet[T]) AddAll(other collections.Collection[T]) {
	for t := range other.All() {
		s.Add(t)
	}
}

func (s *HashSet[T]) RemoveAll(other collections.Collection[T]) {
	for t := range other.All() {
		s.RemoveElement(t)
	}
}

func (s *HashSet[T]) RetainAll(other collections.Collection[T]) {
	for t := range other.All() {
		if !s.Contains(t) {
			s.RemoveElement(t)
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

func defaultConfig() *Config {
	return &Config{}
}
