package treemap

import "slices"

type node[K any, V any] struct {
	leaf     bool
	keys     []K
	values   []V
	children []*node[K, V]
}

func newNode[K any, V any](leaf bool) *node[K, V] {
	return &node[K, V]{
		leaf: leaf,
	}
}

// Get searches for a key in the B-Tree.
func (tm *TreeMap[K, V]) get(n *node[K, V], key K) (V, bool) {
	i, found := slices.BinarySearchFunc(n.keys, key, tm.comparator)
	if found {
		return n.values[i], true
	}
	if n.leaf {
		var zero V
		return zero, false
	}
	return tm.get(n.children[i], key)
}

// Put inserts a new key-value pair or updates an existing one.
func (tm *TreeMap[K, V]) put(key K, value V) {
	// Fast path: try to update an existing key first
	if tm.updateExisting(tm.root, key, value) {
		return
	}

	root := tm.root
	if len(root.keys) == 2*tm.degree-1 {
		s := newNode[K, V](false)
		s.children = append(s.children, root)
		tm.splitChild(s, 0, root)
		tm.root = s
		tm.insertNonFull(s, key, value)
	} else {
		tm.insertNonFull(root, key, value)
	}
	tm.size++
}

func (tm *TreeMap[K, V]) updateExisting(n *node[K, V], key K, value V) bool {
	i, found := slices.BinarySearchFunc(n.keys, key, tm.comparator)
	if found {
		n.values[i] = value
		return true
	}
	if n.leaf {
		return false
	}
	return tm.updateExisting(n.children[i], key, value)
}

func (tm *TreeMap[K, V]) splitChild(x *node[K, V], i int, y *node[K, V]) {
	t := tm.degree
	z := newNode[K, V](y.leaf)

	// Copy the upper t-1 keys and values from y to z
	z.keys = slices.Clone(y.keys[t:])
	z.values = slices.Clone(y.values[t:])
	
	// Truncate y's keys and values
	midKey := y.keys[t-1]
	midVal := y.values[t-1]
	
	var zeroK K
	var zeroV V
	for j := t - 1; j < len(y.keys); j++ {
		y.keys[j] = zeroK
		y.values[j] = zeroV
	}
	y.keys = y.keys[:t-1]
	y.values = y.values[:t-1]

	if !y.leaf {
		z.children = slices.Clone(y.children[t:])
		y.children = y.children[:t]
	}

	x.children = slices.Insert(x.children, i+1, z)
	x.keys = slices.Insert(x.keys, i, midKey)
	x.values = slices.Insert(x.values, i, midVal)
}

func (tm *TreeMap[K, V]) insertNonFull(x *node[K, V], key K, value V) {
	if x.leaf {
		i, _ := slices.BinarySearchFunc(x.keys, key, tm.comparator)
		x.keys = slices.Insert(x.keys, i, key)
		x.values = slices.Insert(x.values, i, value)
	} else {
		i, _ := slices.BinarySearchFunc(x.keys, key, tm.comparator)
		if len(x.children[i].keys) == 2*tm.degree-1 {
			tm.splitChild(x, i, x.children[i])
			if tm.comparator(key, x.keys[i]) > 0 {
				i++
			}
		}
		tm.insertNonFull(x.children[i], key, value)
	}
}

// Remove deletes a key from the B-Tree.
func (tm *TreeMap[K, V]) remove(key K) {
	if tm.root == nil || len(tm.root.keys) == 0 {
		return
	}
	if !tm.ContainsKey(key) {
		return
	}
	tm.deleteNode(tm.root, key)
	if len(tm.root.keys) == 0 {
		if !tm.root.leaf {
			tm.root = tm.root.children[0]
		}
	}
	tm.size--
}

func (tm *TreeMap[K, V]) deleteNode(x *node[K, V], key K) {
	t := tm.degree
	i, found := slices.BinarySearchFunc(x.keys, key, tm.comparator)

	// Key is in this node
	if found {
		if x.leaf {
			// Case 1: The key is in a leaf node
			x.keys = slices.Delete(x.keys, i, i+1)
			x.values = slices.Delete(x.values, i, i+1)
			// Optional: clearing the trailing element to avoid memory leaks
			// But slices.Delete takes care of returning the right slice.
			// Wait, we should zero out the popped elements if we reuse capacity, but slices.Delete doesn't zero them out.
			// Let's rely on standard GC for now, or just leave it.
		} else {
			// Case 2: The key is in an internal node
			y := x.children[i]
			z := x.children[i+1]
			if len(y.keys) >= t {
				// 2a: Predecessor has enough keys
				predKey, predVal := tm.getPredecessor(y)
				x.keys[i] = predKey
				x.values[i] = predVal
				tm.deleteNode(y, predKey)
			} else if len(z.keys) >= t {
				// 2b: Successor has enough keys
				succKey, succVal := tm.getSuccessor(z)
				x.keys[i] = succKey
				x.values[i] = succVal
				tm.deleteNode(z, succKey)
			} else {
				// 2c: Both have t-1 keys, merge them
				tm.merge(x, i)
				tm.deleteNode(y, key)
			}
		}
	} else {
		// Case 3: Key is not in this node
		if x.leaf {
			return
		}
		
		// Ensure the child we descend into has at least t keys
		if len(x.children[i].keys) == t-1 {
			tm.fill(x, i)
			// i might have changed after fill, re-find it
			i, _ = slices.BinarySearchFunc(x.keys, key, tm.comparator)
		}
		tm.deleteNode(x.children[i], key)
	}
}

func (tm *TreeMap[K, V]) getPredecessor(x *node[K, V]) (K, V) {
	for !x.leaf {
		x = x.children[len(x.children)-1]
	}
	return x.keys[len(x.keys)-1], x.values[len(x.values)-1]
}

func (tm *TreeMap[K, V]) getSuccessor(x *node[K, V]) (K, V) {
	for !x.leaf {
		x = x.children[0]
	}
	return x.keys[0], x.values[0]
}

func (tm *TreeMap[K, V]) fill(x *node[K, V], i int) {
	t := tm.degree
	if i != 0 && len(x.children[i-1].keys) >= t {
		tm.borrowFromPrev(x, i)
	} else if i != len(x.children)-1 && len(x.children[i+1].keys) >= t {
		tm.borrowFromNext(x, i)
	} else {
		if i != len(x.children)-1 {
			tm.merge(x, i)
		} else {
			tm.merge(x, i-1)
		}
	}
}

func (tm *TreeMap[K, V]) borrowFromPrev(x *node[K, V], i int) {
	child := x.children[i]
	sibling := x.children[i-1]

	child.keys = slices.Insert(child.keys, 0, x.keys[i-1])
	child.values = slices.Insert(child.values, 0, x.values[i-1])

	if !child.leaf {
		child.children = slices.Insert(child.children, 0, sibling.children[len(sibling.children)-1])
		sibling.children = sibling.children[:len(sibling.children)-1]
	}

	x.keys[i-1] = sibling.keys[len(sibling.keys)-1]
	x.values[i-1] = sibling.values[len(sibling.values)-1]

	sibling.keys = sibling.keys[:len(sibling.keys)-1]
	sibling.values = sibling.values[:len(sibling.values)-1]
}

func (tm *TreeMap[K, V]) borrowFromNext(x *node[K, V], i int) {
	child := x.children[i]
	sibling := x.children[i+1]

	child.keys = append(child.keys, x.keys[i])
	child.values = append(child.values, x.values[i])

	if !child.leaf {
		child.children = append(child.children, sibling.children[0])
		sibling.children = slices.Delete(sibling.children, 0, 1)
	}

	x.keys[i] = sibling.keys[0]
	x.values[i] = sibling.values[0]

	sibling.keys = slices.Delete(sibling.keys, 0, 1)
	sibling.values = slices.Delete(sibling.values, 0, 1)
}

func (tm *TreeMap[K, V]) merge(x *node[K, V], i int) {
	child := x.children[i]
	sibling := x.children[i+1]

	child.keys = append(child.keys, x.keys[i])
	child.values = append(child.values, x.values[i])

	child.keys = append(child.keys, sibling.keys...)
	child.values = append(child.values, sibling.values...)

	if !child.leaf {
		child.children = append(child.children, sibling.children...)
	}

	x.keys = slices.Delete(x.keys, i, i+1)
	x.values = slices.Delete(x.values, i, i+1)
	x.children = slices.Delete(x.children, i+1, i+2)
}
