package util

// MustGet returns the given value of type T if the given error is not nil,
// otherwise it panics.
func MustGet[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// MustDo panics if the given error is not nil.
func MustDo(err error) {
	if err != nil {
		panic(err)
	}
}
