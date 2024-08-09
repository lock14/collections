package linked_list

import (
	"github.com/lock14/collections"
	"testing"
)

func TestType(t *testing.T) {
	l := New[int]()
	testType[int](l)
}

func testType[T any](deque collections.Deque[T]) {}
