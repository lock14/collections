package heap

import (
	"github.com/lock14/collections"
	"testing"
)

func TestHeapImplementsQueue(t *testing.T) {
	queue[int](New[int]())
}

func queue[T any](q collections.Queue[T]) collections.Queue[T] {
	return q
}
