// Package pair provides a generic 2-element tuple.
package pair

import "fmt"

// Pair represents a generic 2-element tuple.
type Pair[T1 any, T2 any] struct {
	fst T1
	snd T2
}

var _ fmt.Stringer = Pair[int, string]{}

// New creates a new Pair.
func New[T1 any, T2 any](t1 T1, t2 T2) Pair[T1, T2] {
	return Pair[T1, T2]{
		fst: t1,
		snd: t2,
	}
}

// Fst returns the first element of the pair.
func (p Pair[T1, T2]) Fst() T1 {
	return p.fst
}

// Snd returns the second element of the pair.
func (p Pair[T1, T2]) Snd() T2 {
	return p.snd
}

// Unwrap returns both elements of the pair as a tuple (fst, snd).
func (p Pair[T1, T2]) Unwrap() (T1, T2) {
	return p.fst, p.snd
}

// String returns the string representation of the pair formatted as "(fst, snd)".
func (p Pair[T1, T2]) String() string {
	return fmt.Sprintf("(%v, %v)", p.fst, p.snd)
}

// benchmark matrix test trigger
