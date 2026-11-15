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

// Reverse returns a comparator that imposes the reverse ordering of cmpFunc.
func Reverse[T any](cmpFunc Comparator[T]) Comparator[T] {
	return func(t1, t2 T) int {
		return cmpFunc(t2, t1)
	}
}
