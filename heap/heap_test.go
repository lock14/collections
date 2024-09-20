package heap

import (
	"github.com/lock14/collections"
	"testing"
)

func TestHeapImplementsQueue(t *testing.T) {
	t.Parallel()

	queueType(New[int]())
}

func queueType[T any](q collections.Queue[T]) collections.Queue[T] {
	return q
}
