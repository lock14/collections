// Package comparator provides types and functions for comparing ordered types.
package comparator

import "cmp"

// Comparator is a function that compares two elements.
// It returns a negative number if t1 < t2, a positive number if t1 > t2, and zero if they are equal.
type Comparator[T any] func(t1, t2 T) int

// NaturalOrder returns a comparator that orders elements using their natural ordering.
func NaturalOrder[T cmp.Ordered]() Comparator[T] {
	return cmp.Compare[T]
}

// Reverse returns a comparator that reverses the ordering of the given comparator.
func Reverse[T any](comparator Comparator[T]) Comparator[T] {
	return func(t1, t2 T) int {
		return -comparator(t1, t2)
	}
}

// benchmark matrix test trigger
