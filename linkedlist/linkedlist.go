package linked_list

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"strings"
)

var _ collections.MutableList[int] = (*LinkedList[int])(nil)
var _ collections.MutableDeque[int] = (*LinkedList[int])(nil)

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

func (l *LinkedList[T]) Peek() T {
	return l.PeekFront()
}

func (l *LinkedList[T]) PeekFront() T {
	if l.Empty() {
		panic("cannot peek from an empty list")
	}
	return l.list.next.data
}

func (l *LinkedList[T]) PeekBack() T {
	if l.Empty() {
		panic("cannot peek from an empty list")
	}
	return l.list.prev.data
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

func (l *LinkedList[T]) Get(idx int) T {
	if n := l.get(idx); n != nil {
		return n.data
	}
	panic("cannot get an element from an empty list")
}

func (l *LinkedList[T]) Set(idx int, t T) {
	if n := l.get(idx); n != nil {
		n.data = t
	}
	panic("cannot set an element from an empty list")
}

func (l *LinkedList[T]) Pop() T {
	return l.RemoveFront()
}

func (l *LinkedList[T]) Size() int {
	return l.size
}

func (l *LinkedList[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
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

func (l *LinkedList[T]) get(idx int) *node[T] {
	count := 0
	cur := l.list.next
	for cur != &l.list {
		if count == idx {
			return cur
		}
		count++
	}
	return nil
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
