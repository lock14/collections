package linkedhashmap

import (
	"github.com/lock14/collections"
	"testing"
)

func TestType(t *testing.T) {
	t.Parallel()

	mapType(New[int, int]())
}

func mapType[K, V any](_ collections.Map[K, V]) {}
