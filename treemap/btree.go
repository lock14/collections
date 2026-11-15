package treemap

import "slices"

type node[K any, V any] struct {
	leaf     bool
	keys     []K
	values   []V
	children []*node[K, V]
}

func (m *TreeMap[K, V]) newNode(leaf bool) *node[K, V] {
	var children []*node[K, V]
	if !leaf {
		children = make([]*node[K, V], 0, 2*m.degree)
	}
	return &node[K, V]{
		leaf:     leaf,
		keys:     make([]K, 0, 2*m.degree-1),
		values:   make([]V, 0, 2*m.degree-1),
		children: children,
	}
}

func (m *TreeMap[K, V]) get(n *node[K, V], key K) (V, bool) {
	for {
		i, found := slices.BinarySearchFunc(n.keys, key, m.comparator)
		if found {
			return n.values[i], true
		}
		if n.leaf {
			var zero V
			return zero, false
		}
		n = n.children[i]
	}
}

func (m *TreeMap[K, V]) put(key K, value V) {
	if m.root == nil {
		m.root = m.newNode(true)
	}
	root := m.root
	if len(root.keys) == 2*m.degree-1 {
		if i, found := slices.BinarySearchFunc(root.keys, key, m.comparator); found {
			root.values[i] = value
			return
		}
		s := m.newNode(false)
		s.children = append(s.children, root)
		m.splitChild(s, 0, root)
		m.root = s
		m.putNonFull(s, key, value)
		return
	}
	m.putNonFull(root, key, value)
}

func (m *TreeMap[K, V]) splitChild(x *node[K, V], i int, y *node[K, V]) {
	t := m.degree
	z := m.newNode(y.leaf)

	z.keys = append(z.keys, y.keys[t:]...)
	z.values = append(z.values, y.values[t:]...)

	midKey := y.keys[t-1]
	midVal := y.values[t-1]

	clear(y.keys[t-1:])
	clear(y.values[t-1:])
	y.keys = y.keys[:t-1]
	y.values = y.values[:t-1]

	if !y.leaf {
		z.children = append(z.children, y.children[t:]...)
		clear(y.children[t:])
		y.children = y.children[:t]
	}

	x.children = slices.Insert(x.children, i+1, z)
	x.keys = slices.Insert(x.keys, i, midKey)
	x.values = slices.Insert(x.values, i, midVal)
}

func (m *TreeMap[K, V]) putNonFull(x *node[K, V], key K, value V) {
	for {
		i, found := slices.BinarySearchFunc(x.keys, key, m.comparator)
		if found {
			x.values[i] = value
			return
		}
		if x.leaf {
			x.keys = slices.Insert(x.keys, i, key)
			x.values = slices.Insert(x.values, i, value)
			m.size++
			return
		}
		child := x.children[i]
		if len(child.keys) == 2*m.degree-1 {
			cmp := m.comparator(key, child.keys[m.degree-1])
			if cmp == 0 {
				child.values[m.degree-1] = value
				return
			}
			m.splitChild(x, i, child)
			if cmp > 0 {
				i++
			}
		}
		x = x.children[i]
	}
}

func (m *TreeMap[K, V]) remove(key K) {
	if m.root == nil || len(m.root.keys) == 0 {
		return
	}
	deleted := m.deleteNode(m.root, key)
	if len(m.root.keys) == 0 && !m.root.leaf {
		m.root = m.root.children[0]
	}
	if deleted {
		m.size--
	}
}

func (m *TreeMap[K, V]) deleteNode(x *node[K, V], key K) bool {
	t := m.degree
	i, found := slices.BinarySearchFunc(x.keys, key, m.comparator)

	if found {
		if x.leaf {
			x.keys = slices.Delete(x.keys, i, i+1)
			x.values = slices.Delete(x.values, i, i+1)
			return true
		}
		y := x.children[i]
		z := x.children[i+1]
		if len(y.keys) >= t {
			predKey, predVal := m.getPredecessor(y)
			x.keys[i] = predKey
			x.values[i] = predVal
			return m.deleteNode(y, predKey)
		}
		if len(z.keys) >= t {
			succKey, succVal := m.getSuccessor(z)
			x.keys[i] = succKey
			x.values[i] = succVal
			return m.deleteNode(z, succKey)
		}
		m.merge(x, i)
		return m.deleteNode(y, key)
	}

	if x.leaf {
		return false
	}
	if len(x.children[i].keys) == t-1 {
		m.fill(x, i)
		if i > len(x.keys) {
			i--
		}
	}
	return m.deleteNode(x.children[i], key)
}

func (m *TreeMap[K, V]) getPredecessor(x *node[K, V]) (K, V) {
	for !x.leaf {
		x = x.children[len(x.children)-1]
	}
	return x.keys[len(x.keys)-1], x.values[len(x.values)-1]
}

func (m *TreeMap[K, V]) getSuccessor(x *node[K, V]) (K, V) {
	for !x.leaf {
		x = x.children[0]
	}
	return x.keys[0], x.values[0]
}

func (m *TreeMap[K, V]) fill(x *node[K, V], i int) {
	t := m.degree
	if i != 0 && len(x.children[i-1].keys) >= t {
		m.borrowFromPrev(x, i)
	} else if i != len(x.children)-1 && len(x.children[i+1].keys) >= t {
		m.borrowFromNext(x, i)
	} else {
		if i != len(x.children)-1 {
			m.merge(x, i)
		} else {
			m.merge(x, i-1)
		}
	}
}

func (m *TreeMap[K, V]) borrowFromPrev(x *node[K, V], i int) {
	child := x.children[i]
	sibling := x.children[i-1]

	child.keys = slices.Insert(child.keys, 0, x.keys[i-1])
	child.values = slices.Insert(child.values, 0, x.values[i-1])

	if !child.leaf {
		lastChild := len(sibling.children) - 1
		child.children = slices.Insert(child.children, 0, sibling.children[lastChild])
		sibling.children[lastChild] = nil
		sibling.children = sibling.children[:lastChild]
	}

	lastIdx := len(sibling.keys) - 1
	x.keys[i-1] = sibling.keys[lastIdx]
	x.values[i-1] = sibling.values[lastIdx]

	clear(sibling.keys[lastIdx:])
	clear(sibling.values[lastIdx:])
	sibling.keys = sibling.keys[:lastIdx]
	sibling.values = sibling.values[:lastIdx]
}

func (m *TreeMap[K, V]) borrowFromNext(x *node[K, V], i int) {
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

func (m *TreeMap[K, V]) merge(x *node[K, V], i int) {
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
