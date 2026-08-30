// Package treeset provides a B-Tree backed set implementation.
package treeset

import (
	"cmp"
	"github.com/lock14/collections/comparator"
	"github.com/lock14/collections/treemap"
)

// TreeSet represents a set of elements backed by a B-Tree.
type TreeSet[T any] struct {
	m *treemap.TreeMap[T, struct{}]
}

// config holds the values for configuring a TreeSet.
type config[T any] struct {
	degree     *int
	comparator comparator.Comparator[T]
}

// Option configures a TreeSet config.
type Option[T any] func(*config[T])

// WithDegree configures the degree of the underlying B-Tree.
func WithDegree[T any](degree int) Option[T] {
	return func(c *config[T]) {
		c.degree = &degree
	}
}

// WithComparator configures the comparator for the TreeSet.
func WithComparator[T any](comp comparator.Comparator[T]) Option[T] {
	return func(c *config[T]) {
		c.comparator = comp
	}
}

// New creates an empty TreeSet.
func New[T any](opts ...Option[T]) *TreeSet[T] {
	config := &config[T]{}
	for _, option := range opts {
		option(config)
	}

	var mapOpts []treemap.Option[T]
	if config.degree != nil {
		mapOpts = append(mapOpts, treemap.WithDegree[T](*config.degree))
	}
	if config.comparator != nil {
		mapOpts = append(mapOpts, treemap.WithComparator[T](config.comparator))
	}

	return &TreeSet[T]{
		m: treemap.New[T, struct{}](mapOpts...),
	}
}

// NewOrdered creates an empty TreeSet for types that implement cmp.Ordered.
func NewOrdered[T cmp.Ordered](opts ...Option[T]) *TreeSet[T] {
	opts = append([]Option[T]{WithComparator(comparator.NaturalOrder[T]())}, opts...)
	return New(opts...)
}
