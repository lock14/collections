package linked_list

import (
	"github.com/lock14/collections"
	"testing"
)

func TestType(t *testing.T) {
	t.Parallel()

	listType(New[int]())
	dequeType[int](New[int]())
}

func listType[T any](_ collections.List[T]) {}

func dequeType[T any](_ collections.Deque[T]) {}
