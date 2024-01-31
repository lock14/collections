package util

func MustGet[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func MustDo(err error) {
	if err != nil {
		panic(err)
	}
}
