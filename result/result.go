// Package result provides a generic Result[T, E] type representing either success (Ok) or failure (Err).
package result

import (
	"fmt"
)

// Result represents either success (Ok) with value of type T or failure (Err) with error of type E.
type Result[T any, E any] struct {
	value T
	err   E
	isOk  bool
}

var _ fmt.Stringer = Result[int, error]{}

// Ok creates a success Result containing value.
func Ok[T any, E any](value T) Result[T, E] {
	return Result[T, E]{
		value: value,
		isOk:  true,
	}
}

// Err creates a failure Result containing err.
func Err[T any, E any](err E) Result[T, E] {
	return Result[T, E]{
		err:  err,
		isOk: false,
	}
}

// Of creates a Result from a value and standard error. If err is not nil, an Err Result is returned;
// otherwise, an Ok Result containing value is returned.
func Of[T any](value T, err error) Result[T, error] {
	if err != nil {
		return Result[T, error]{
			err:  err,
			isOk: false,
		}
	}
	return Result[T, error]{
		value: value,
		isOk:  true,
	}
}

// IsOk returns true if the Result is Ok.
func (r Result[T, E]) IsOk() bool {
	return r.isOk
}

// IsErr returns true if the Result is Err.
func (r Result[T, E]) IsErr() bool {
	return !r.isOk
}

// Ok returns the value and a boolean indicating if the Result is Ok.
func (r Result[T, E]) Ok() (T, bool) {
	return r.value, r.isOk
}

// Err returns the error and a boolean indicating if the Result is Err.
func (r Result[T, E]) Err() (E, bool) {
	return r.err, !r.isOk
}

// Unwrap returns both the value and error as a tuple.
func (r Result[T, E]) Unwrap() (T, E) {
	return r.value, r.err
}

// MustGet returns the value if the Result is Ok, or panics if the Result is Err.
func (r Result[T, E]) MustGet() T {
	if !r.isOk {
		panic(fmt.Sprintf("result: MustGet called on Err(%v)", r.err))
	}
	return r.value
}

// OrElse returns the value if the Result is Ok, or value otherwise.
func (r Result[T, E]) OrElse(value T) T {
	if r.isOk {
		return r.value
	}
	return value
}

// OrElseGet returns the value if the Result is Ok, or the result of invoking supplier otherwise.
func (r Result[T, E]) OrElseGet(supplier func() T) T {
	if r.isOk {
		return r.value
	}
	return supplier()
}

// Map applies the provided mapper function to the contained value if Ok,
// and returns a Result containing the mapped value. If Err, it returns the original error.
func (r Result[T, E]) Map[U any](mapper func(value T) U) Result[U, E] {
	if !r.isOk {
		return Result[U, E]{
			err:  r.err,
			isOk: false,
		}
	}
	return Result[U, E]{
		value: mapper(r.value),
		isOk:  true,
	}
}

// MapErr applies the provided mapper function to the contained error if Err,
// and returns a Result containing the mapped error. If Ok, it returns the original value.
func (r Result[T, E]) MapErr[F any](mapper func(err E) F) Result[T, F] {
	if r.isOk {
		return Result[T, F]{
			value: r.value,
			isOk:  true,
		}
	}
	return Result[T, F]{
		err:  mapper(r.err),
		isOk: false,
	}
}

// FlatMap applies the provided Result-bearing mapper function to the contained value if Ok,
// and returns the resulting Result. If Err, it returns the original error.
func (r Result[T, E]) FlatMap[U any](mapper func(value T) Result[U, E]) Result[U, E] {
	if !r.isOk {
		return Result[U, E]{
			err:  r.err,
			isOk: false,
		}
	}
	return mapper(r.value)
}

// IfOk executes the provided action with the contained value if the Result is Ok.
func (r Result[T, E]) IfOk(action func(value T)) {
	if r.isOk {
		action(r.value)
	}
}

// IfErr executes the provided action with the contained error if the Result is Err.
func (r Result[T, E]) IfErr(action func(err E)) {
	if !r.isOk {
		action(r.err)
	}
}

// IfOkOrElse executes okAction with the contained value if Ok; otherwise executes errAction with the contained error.
func (r Result[T, E]) IfOkOrElse(okAction func(value T), errAction func(err E)) {
	if r.isOk {
		okAction(r.value)
	} else {
		errAction(r.err)
	}
}

// String returns a string representation of the Result.
// If Ok, it returns "Ok(value)"; if Err, it returns "Err(err)".
func (r Result[T, E]) String() string {
	if r.isOk {
		return fmt.Sprintf("Ok(%v)", r.value)
	}
	return fmt.Sprintf("Err(%v)", r.err)
}

// benchmark matrix test trigger
