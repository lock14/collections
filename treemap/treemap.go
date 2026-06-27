// Package treemap provides a B-Tree backed map implementation.
package treemap

import (
	"cmp"
	"github.com/lock14/collections/comparator"
)

const (
	// DefaultDegree is the default minimum degree (t) for the B-Tree.
	DefaultDegree = 32
)

// config holds the values for configuring a TreeMap.
type config[K any] struct {
	degree     int
	comparator comparator.Comparator[K]
}

// Option configures a TreeMap config
type Option[K any] func(*config[K])

// WithDegree configures the minimum degree of the B-Tree.
func WithDegree[K any](degree int) Option[K] {
	return func(c *config[K]) {
		c.degree = degree
	}
}

// WithComparator configures the comparator used by the TreeMap.
func WithComparator[K any](cmpFunc comparator.Comparator[K]) Option[K] {
	return func(c *config[K]) {
		c.comparator = cmpFunc
	}
}

// TreeMap is a B-Tree backed map that implements collections.MutableMap.
type TreeMap[K any, V any] struct {
	root       *node[K, V]
	size       int
	degree     int
	comparator comparator.Comparator[K]
}

// New creates an empty TreeMap with the given options.
func New[K any, V any](opts ...Option[K]) *TreeMap[K, V] {
	config := &config[K]{
		degree: DefaultDegree,
	}
	for _, option := range opts {
		option(config)
	}
	if config.degree < 2 {
		panic("degree must be at least 2")
	}
	if config.comparator == nil {
		panic("comparator must be provided or use NewOrdered")
	}
	tm := &TreeMap[K, V]{
		degree:     config.degree,
		comparator: config.comparator,
	}
	tm.root = tm.newNode(true)
	return tm
}

// NewOrdered creates an empty TreeMap for keys that satisfy cmp.Ordered using natural ordering.
func NewOrdered[K cmp.Ordered, V any](opts ...Option[K]) *TreeMap[K, V] {
	return New[K, V](append(opts, WithComparator(comparator.NaturalOrder[K]()))...)
}
