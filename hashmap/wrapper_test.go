package hashmap

import "testing"

func TestWrapperType(t *testing.T) {
	t.Parallel()

	mapType(Wrap(make(map[int]int)))
}
