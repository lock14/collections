package treemap

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"slices"
	"strings"
)

var (
	_ collections.MutableNavigableMap[int, int] = (*TreeMap[int, int])(nil)
	_ fmt.Stringer                              = (*TreeMap[int, int])(nil)
)

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
	tm.root = tm.newNode(true)
	tm.size = 0
}

func (tm *TreeMap[K, V]) ContainsKey(key K) bool {
	_, ok := tm.Get(key)
	return ok
}

// String returns a string representation of the map matching Go built-in map formatting.
func (tm *TreeMap[K, V]) String() string {
	vals := make([]string, 0, tm.Size())
	for k, v := range tm.All() {
		vals = append(vals, fmt.Sprintf("%v:%v", k, v))
	}
	return "map[" + strings.Join(vals, " ") + "]"
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
		for k := range tm.All() {
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

// First returns the first key-value pair in the map.
func (tm *TreeMap[K, V]) First() (K, V) {
	if tm.Empty() {
		panic("First called on empty map")
	}
	n := tm.root
	for !n.leaf {
		n = n.children[0]
	}
	return n.keys[0], n.values[0]
}

// Last returns the last key-value pair in the map.
func (tm *TreeMap[K, V]) Last() (K, V) {
	if tm.Empty() {
		panic("Last called on empty map")
	}
	n := tm.root
	for !n.leaf {
		n = n.children[len(n.children)-1]
	}
	return n.keys[len(n.keys)-1], n.values[len(n.values)-1]
}

// PollFirst removes and returns the first key-value pair in the map.
func (tm *TreeMap[K, V]) PollFirst() (K, V) {
	if tm.Empty() {
		panic("PollFirst called on empty map")
	}
	k, v := tm.First()
	tm.Remove(k)
	return k, v
}

// PollLast removes and returns the last key-value pair in the map.
func (tm *TreeMap[K, V]) PollLast() (K, V) {
	if tm.Empty() {
		panic("PollLast called on empty map")
	}
	k, v := tm.Last()
	tm.Remove(k)
	return k, v
}

func (tm *TreeMap[K, V]) Lower(key K) (K, V, bool) {
	var bestK K
	var bestV V
	var found bool

	n := tm.root
	for n != nil && len(n.keys) > 0 {
		i, _ := slices.BinarySearchFunc(n.keys, key, tm.comparator)
		if i > 0 {
			bestK = n.keys[i-1]
			bestV = n.values[i-1]
			found = true
		}
		if n.leaf {
			break
		}
		n = n.children[i]
	}
	return bestK, bestV, found
}

func (tm *TreeMap[K, V]) Floor(key K) (K, V, bool) {
	var bestK K
	var bestV V
	var found bool

	n := tm.root
	for n != nil && len(n.keys) > 0 {
		i, match := slices.BinarySearchFunc(n.keys, key, tm.comparator)
		if match {
			return n.keys[i], n.values[i], true
		}
		if i > 0 {
			bestK = n.keys[i-1]
			bestV = n.values[i-1]
			found = true
		}
		if n.leaf {
			break
		}
		n = n.children[i]
	}
	return bestK, bestV, found
}

func (tm *TreeMap[K, V]) Ceiling(key K) (K, V, bool) {
	var bestK K
	var bestV V
	var found bool

	n := tm.root
	for n != nil && len(n.keys) > 0 {
		i, match := slices.BinarySearchFunc(n.keys, key, tm.comparator)
		if match {
			return n.keys[i], n.values[i], true
		}
		if i < len(n.keys) {
			bestK = n.keys[i]
			bestV = n.values[i]
			found = true
		}
		if n.leaf {
			break
		}
		n = n.children[i]
	}
	return bestK, bestV, found
}

func (tm *TreeMap[K, V]) Higher(key K) (K, V, bool) {
	var bestK K
	var bestV V
	var found bool

	n := tm.root
	for n != nil && len(n.keys) > 0 {
		i, match := slices.BinarySearchFunc(n.keys, key, tm.comparator)
		if match {
			i++
		}
		if i < len(n.keys) {
			bestK = n.keys[i]
			bestV = n.values[i]
			found = true
		}
		if n.leaf {
			break
		}
		n = n.children[i]
	}
	return bestK, bestV, found
}

func (tm *TreeMap[K, V]) Backward() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		tm.reverseInOrder(tm.root, yield)
	}
}

func (tm *TreeMap[K, V]) reverseInOrder(n *node[K, V], yield func(K, V) bool) bool {
	if n == nil || len(n.keys) == 0 {
		return true
	}
	if !n.leaf {
		if !tm.reverseInOrder(n.children[len(n.keys)], yield) {
			return false
		}
	}
	for i := len(n.keys) - 1; i >= 0; i-- {
		if !yield(n.keys[i], n.values[i]) {
			return false
		}
		if !n.leaf {
			if !tm.reverseInOrder(n.children[i], yield) {
				return false
			}
		}
	}
	return true
}

func (tm *TreeMap[K, V]) BackwardKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range tm.Backward() {
			if !yield(k) {
				return
			}
		}
	}
}

func (tm *TreeMap[K, V]) BackwardValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range tm.Backward() {
			if !yield(v) {
				return
			}
		}
	}
}

func (tm *TreeMap[K, V]) From(from K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var zero K
		tm.rangeInOrder(tm.root, from, zero, true, false, yield)
	}
}

func (tm *TreeMap[K, V]) To(to K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var zero K
		tm.rangeInOrder(tm.root, zero, to, false, true, yield)
	}
}

func (tm *TreeMap[K, V]) Between(from K, to K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		tm.rangeInOrder(tm.root, from, to, true, true, yield)
	}
}

func (tm *TreeMap[K, V]) rangeInOrder(n *node[K, V], from K, to K, checkFrom, checkTo bool, yield func(K, V) bool) bool {
	if n == nil || len(n.keys) == 0 {
		return true
	}
	for i := 0; i < len(n.keys); i++ {
		cmpFrom := 1
		if checkFrom {
			cmpFrom = tm.comparator(n.keys[i], from)
		}
		cmpTo := -1
		if checkTo {
			cmpTo = tm.comparator(n.keys[i], to)
		}

		if !n.leaf && cmpFrom >= 0 {
			if !tm.rangeInOrder(n.children[i], from, to, checkFrom, checkTo, yield) {
				return false
			}
		}
		if cmpTo >= 0 {
			return false
		}
		if cmpFrom >= 0 {
			if !yield(n.keys[i], n.values[i]) {
				return false
			}
		}
	}
	if !n.leaf {
		if !tm.rangeInOrder(n.children[len(n.keys)], from, to, checkFrom, checkTo, yield) {
			return false
		}
	}
	return true
}
