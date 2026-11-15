// Package linkedlist provides a doubly-linked list implementation.
package linkedlist

import (
	"fmt"
	"iter"
	"strings"

	"github.com/lock14/collections"
)

var (
	_ collections.MutableList[int]  = (*LinkedList[int])(nil)
	_ collections.MutableDeque[int] = (*LinkedList[int])(nil)
	_ fmt.Stringer                  = (*LinkedList[int])(nil)
)

// LinkedList is a doubly-linked list implementation.
// The zero value for LinkedList is an empty list ready to use.
type LinkedList[T any] struct {
	list node[T]
	size int
}

type node[T any] struct {
	data T
	prev *node[T]
	next *node[T]
}

// New creates an empty LinkedList.
func New[T any]() *LinkedList[T] {
	l := &LinkedList[T]{}
	l.list.next = &l.list
	l.list.prev = &l.list
	return l
}

// Add appends the specified element to the end of the list.
// It is an alias for AddBack.
func (l *LinkedList[T]) Add(t T) {
	l.AddBack(t)
}

// AddFront inserts the specified element at the front of the list.
func (l *LinkedList[T]) AddFront(t T) {
	l.lazyInit()
	insertBefore(l.list.next, t)
	l.size++
}

// AddBack inserts the specified element at the back of the list.
func (l *LinkedList[T]) AddBack(t T) {
	l.lazyInit()
	insertBefore(&l.list, t)
	l.size++
}

// Push adds the specified element to the front of the list.
// It is an alias for AddFront.
func (l *LinkedList[T]) Push(t T) {
	l.AddFront(t)
}

// AddAll appends all elements from the given sequence to the back of the list.
func (l *LinkedList[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		l.Add(t)
	}
}

// Set replaces the element at the specified index with the given element.
// Panics if idx is out of bounds (idx < 0 || idx >= Size()).
func (l *LinkedList[T]) Set(idx int, t T) {
	if n := l.get(idx); n != nil {
		n.data = t
		return
	}
	panic("linkedlist: index out of bounds")
}

// Remove removes and returns the element at the front of the list.
// It is an alias for RemoveFront. Panics if the list is empty.
func (l *LinkedList[T]) Remove() T {
	return l.RemoveFront()
}

// RemoveFront removes and returns the element at the front of the list.
// Panics if the list is empty.
func (l *LinkedList[T]) RemoveFront() T {
	if l.Empty() {
		panic("linkedlist: cannot remove from an empty list")
	}
	n := l.list.next
	val := n.data
	unlink(n)
	l.size--
	return val
}

// RemoveBack removes and returns the element at the back of the list.
// Panics if the list is empty.
func (l *LinkedList[T]) RemoveBack() T {
	if l.Empty() {
		panic("linkedlist: cannot remove from an empty list")
	}
	n := l.list.prev
	val := n.data
	unlink(n)
	l.size--
	return val
}

// Pop removes and returns the element at the front of the list.
// It is an alias for RemoveFront. Panics if the list is empty.
func (l *LinkedList[T]) Pop() T {
	return l.RemoveFront()
}

// Clear removes all elements from the list.
func (l *LinkedList[T]) Clear() {
	l.list.next = &l.list
	l.list.prev = &l.list
	l.size = 0
}

// Get returns the element at the specified index.
// Panics if idx is out of bounds (idx < 0 || idx >= Size()).
func (l *LinkedList[T]) Get(idx int) T {
	if n := l.get(idx); n != nil {
		return n.data
	}
	panic("linkedlist: index out of bounds")
}

// Peek returns the element at the front of the list without removing it.
// It is an alias for PeekFront. Panics if the list is empty.
func (l *LinkedList[T]) Peek() T {
	return l.PeekFront()
}

// PeekFront returns the element at the front of the list without removing it.
// Panics if the list is empty.
func (l *LinkedList[T]) PeekFront() T {
	if l.Empty() {
		panic("linkedlist: cannot peek from an empty list")
	}
	return l.list.next.data
}

// PeekBack returns the element at the back of the list without removing it.
// Panics if the list is empty.
func (l *LinkedList[T]) PeekBack() T {
	if l.Empty() {
		panic("linkedlist: cannot peek from an empty list")
	}
	return l.list.prev.data
}

// Size returns the number of elements in the list.
func (l *LinkedList[T]) Size() int {
	return l.size
}

// Empty returns true if the list contains no elements.
func (l *LinkedList[T]) Empty() bool {
	return l.size == 0
}

// All returns an iterator over all elements of the list from front to back.
func (l *LinkedList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if l.list.next == nil {
			return
		}
		for cur := l.list.next; cur != &l.list; cur = cur.next {
			if !yield(cur.data) {
				return
			}
		}
	}
}

// Backward returns an iterator over all elements of the list in reverse order from back to front.
func (l *LinkedList[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		if l.list.prev == nil {
			return
		}
		for cur := l.list.prev; cur != &l.list; cur = cur.prev {
			if !yield(cur.data) {
				return
			}
		}
	}
}

// String returns a string representation of the list matching Go slice formatting.
func (l *LinkedList[T]) String() string {
	if l.Empty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for v := range l.All() {
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte(']')
	return sb.String()
}

func (l *LinkedList[T]) lazyInit() {
	if l.list.next == nil {
		l.list.next = &l.list
		l.list.prev = &l.list
	}
}

func (l *LinkedList[T]) get(idx int) *node[T] {
	if idx < 0 || idx >= l.size {
		return nil
	}
	if idx < l.size/2 {
		cur := l.list.next
		for i := 0; i < idx; i++ {
			cur = cur.next
		}
		return cur
	}
	cur := l.list.prev
	for i := l.size - 1; i > idx; i-- {
		cur = cur.prev
	}
	return cur
}

func insertBefore[T any](n *node[T], t T) {
	newNode := &node[T]{
		data: t,
		prev: n.prev,
		next: n,
	}
	n.prev.next = newNode
	n.prev = newNode
}

func unlink[T any](n *node[T]) {
	var zero T
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
	n.data = zero
}
