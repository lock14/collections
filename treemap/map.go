package treemap

import (
	"github.com/lock14/collections"
	"iter"
)

var _ collections.MutableMap[int, int] = (*TreeMap[int, int])(nil)

func (tm *TreeMap[K, V]) Get(key K) (V, bool) {
	if tm.root == nil {
		var zero V
		return zero, false
	}
	return tm.get(tm.root, key)
}

func (tm *TreeMap[K, V]) Put(key K, value V) {
	tm.put(key, value)
}

func (tm *TreeMap[K, V]) Remove(key K) {
	tm.remove(key)
}

func (tm *TreeMap[K, V]) Size() int {
	return tm.size
}

func (tm *TreeMap[K, V]) Empty() bool {
	return tm.size == 0
}

func (tm *TreeMap[K, V]) Clear() {
	tm.root = newNode[K, V](true)
	tm.size = 0
}

func (tm *TreeMap[K, V]) ContainsKey(key K) bool {
	_, ok := tm.Get(key)
	return ok
}

func (tm *TreeMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		tm.inOrder(tm.root, yield)
	}
}

func (tm *TreeMap[K, V]) inOrder(n *node[K, V], yield func(K, V) bool) bool {
	if n == nil {
		return true
	}
	for i := 0; i < len(n.keys); i++ {
		if !n.leaf {
			if !tm.inOrder(n.children[i], yield) {
				return false
			}
		}
		if !yield(n.keys[i], n.values[i]) {
			return false
		}
	}
	if !n.leaf {
		if !tm.inOrder(n.children[len(n.keys)], yield) {
			return false
		}
	}
	return true
}

func (tm *TreeMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range tm.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func (tm *TreeMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range tm.All() {
			if !yield(v) {
				return
			}
		}
	}
}
