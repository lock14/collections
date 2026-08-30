// Package optional provides a generic Option[T] type representing optional values.
package optional

import (
	"fmt"
)

// Option represents an optional value: every Option is either empty or contains a value.
type Option[T any] struct {
	value T
	ok    bool
}

var _ fmt.Stringer = Option[int]{}

// Of creates an Option containing the specified non-empty value.
func Of[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		ok:    true,
	}
}

// Empty returns an empty Option with no value.
func Empty[T any]() Option[T] {
	return Option[T]{}
}

// OfPtr creates an Option from a pointer. If ptr is nil, an empty Option is returned;
// otherwise, an Option containing the dereferenced value is returned.
func OfPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return Option[T]{}
	}
	return Option[T]{
		value: *ptr,
		ok:    true,
	}
}

// OfOk creates an Option based on a value and a boolean flag. If ok is false, an empty Option
// is returned; otherwise, an Option containing the value is returned.
func OfOk[T any](value T, ok bool) Option[T] {
	if !ok {
		return Option[T]{}
	}
	return Option[T]{
		value: value,
		ok:    true,
	}
}

// IsPresent returns true if the Option contains a value, or false if it is empty.
func (o Option[T]) IsPresent() bool {
	return o.ok
}

// IsEmpty returns true if the Option is empty, or false if it contains a value.
func (o Option[T]) IsEmpty() bool {
	return !o.ok
}

// Get returns the value and a boolean indicating whether the value is present.
func (o Option[T]) Get() (T, bool) {
	return o.value, o.ok
}

// MustGet returns the value if present, or panics if the Option is empty.
func (o Option[T]) MustGet() T {
	if !o.ok {
		panic("optional: MustGet called on empty Option")
	}
	return o.value
}

// OrElse returns the value if present, or value otherwise.
func (o Option[T]) OrElse(value T) T {
	if o.ok {
		return o.value
	}
	return value
}

// OrElseGet returns the value if present, or the result of invoking supplier otherwise.
func (o Option[T]) OrElseGet(supplier func() T) T {
	if o.ok {
		return o.value
	}
	return supplier()
}

// Map applies the provided mapper function to the contained value if present,
// and returns an Option containing the result. If empty, it returns an empty Option.
func (o Option[T]) Map[U any](mapper func(value T) U) Option[U] {
	if !o.ok {
		return Option[U]{}
	}
	return Option[U]{
		value: mapper(o.value),
		ok:    true,
	}
}

// FlatMap applies the provided Option-bearing mapper function to the contained value if present,
// and returns the resulting Option. If empty, it returns an empty Option.
func (o Option[T]) FlatMap[U any](mapper func(value T) Option[U]) Option[U] {
	if !o.ok {
		return Option[U]{}
	}
	return mapper(o.value)
}

// Filter returns this Option if it contains a value that matches the given predicate;
// otherwise, it returns an empty Option.
func (o Option[T]) Filter(predicate func(value T) bool) Option[T] {
	if o.ok && predicate(o.value) {
		return o
	}
	return Option[T]{}
}

// IfPresent executes the provided action with the contained value if a value is present.
func (o Option[T]) IfPresent(action func(value T)) {
	if o.ok {
		action(o.value)
	}
}

// IfPresentOrElse executes the provided action with the contained value if a value is present;
// otherwise, it executes the provided emptyAction.
func (o Option[T]) IfPresentOrElse(action func(value T), emptyAction func()) {
	if o.ok {
		action(o.value)
	} else {
		emptyAction()
	}
}

// String returns a string representation of the Option.
// If present, it returns "Some(value)"; otherwise, it returns "None".
func (o Option[T]) String() string {
	if o.ok {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}
