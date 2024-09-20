package arraylist

import (
	"github.com/lock14/collections"
	"testing"
)

func TestType(t *testing.T) {
	t.Parallel()

	listType(Wrap(make([]int, 0)))
}

func listType[T any](_ collections.List[T]) {}
