package linked_list

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"strings"
)

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

func (l *LinkedList[T]) AddFront(t T) {
	insertBefore(l.list.next, t)
	l.size++
}

func (l *LinkedList[T]) RemoveFront() T {
	if l.Empty() {
		panic("cannot remove from an empty list")
	}
	n := l.list.next
	unlink(n)
	l.size--
	return n.data
}

func (l *LinkedList[T]) AddBack(t T) {
	insertBefore(&l.list, t)
	l.size++
}

func (l *LinkedList[T]) RemoveBack() T {
	if l.Empty() {
		panic("cannot remove from an empty list")
	}
	n := l.list.prev
	unlink(n)
	l.size--
	return n.data
}

func (l *LinkedList[T]) Add(t T) {
	l.AddBack(t)
}

func (l *LinkedList[T]) Remove() T {
	return l.RemoveFront()
}

func (l *LinkedList[T]) Push(t T) {
	l.AddFront(t)
}

func (l *LinkedList[T]) Pop() T {
	return l.RemoveFront()
}

func (l *LinkedList[T]) Size() int {
	return l.size
}

func (l *LinkedList[T]) AddAll(other collections.Collection[T]) {
	for t := range other.All() {
		l.Add(t)
	}
}

func (l *LinkedList[T]) Empty() bool {
	return l.Size() == 0
}

func (l *LinkedList[T]) String() string {
	str := make([]string, 0, l.Size())
	for t := range l.All() {
		str = append(str, fmt.Sprintf("%+v", t))
	}
	return "[" + strings.Join(str, ", ") + "]"
}

func (l *LinkedList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for cur := l.list.next; cur != &l.list && yield(cur.data); {
			cur = cur.next
		}
	}
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
