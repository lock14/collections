package hashset

import (
	"fmt"
	"github.com/lock14/collections/hashmap"
	"github.com/lock14/collections/iterator"
	"strings"
)

// HashSet represents a set of elements of type T.
type HashSet[T comparable] struct {
	m *hashmap.HashMap[T, struct{}]
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
		m: hashmap.New[T, struct{}](),
	}
}

func (s *HashSet[T]) Add(item T) {
	s.m.Put(item, struct{}{})
}

func (s *HashSet[T]) Remove(item T) {
	s.m.Remove(item)
}

func (s *HashSet[T]) Contains(item T) bool {
	_, present := s.m.Get(item)
	return present
}

func (s *HashSet[T]) Size() int {
	return s.m.Size()
}

func (s *HashSet[T]) Empty() bool {
	return s.m.Empty()
}

func (s *HashSet[T]) String() string {
	vals := make([]string, s.Size())
	i := 0
	for item := range iterator.Elements(s.m.Keys()) {
		vals[i] = fmt.Sprintf("%+v", item)
		i++
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

func (s *HashSet[T]) Iterator() iterator.Iterator[T] {
	return s.m.Keys()
}

func (s *HashSet[T]) Elements() chan T {
	return iterator.Elements(s.Iterator())
}

func (s *HashSet[T]) ToSlice() []T {
	slice := make([]T, 0, s.Size())
	for item := range iterator.Elements(s.m.Keys()) {
		slice = append(slice, item)
	}
	return slice
}

func defaultConfig() *Config {
	return &Config{}
}
