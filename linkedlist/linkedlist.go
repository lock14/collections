package linked_list

type LinkedList[T any] struct {
	list node[T]
	size int
}

type node[T any] struct {
	data T
	prev *node[T]
	next *node[T]
}

func New[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		list: sentinel[T](),
		size: 0,
	}
}

func (q *LinkedList[T]) Add(t T) {
	insertBefore(&q.list, t)
	q.size++
}

func (q *LinkedList[T]) Remove() T {
	if q.IsEmpty() {
		panic("cannot remove from an empty list")
	}
	n := q.list.next
	unlink(n)
	q.size--
	return n.data
}

func (q *LinkedList[T]) IsEmpty() bool {
	return q.Size() == 0
}

func (q *LinkedList[T]) Size() int {
	return q.size
}

func insertBefore[T any](n *node[T], t T) {
	newNode := node[T]{
		data: t,
		prev: n.prev,
		next: n,
	}
	n.prev.next = &newNode
	n.prev = &newNode
}

func unlink[T any](n *node[T]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
}

func sentinel[T any]() node[T] {
	var n node[T]
	n.next = &n
	n.prev = &n
	return n
}
