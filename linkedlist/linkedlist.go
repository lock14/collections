// Package linkedlist provides a doubly-linked list implementation.
package linkedlist

import (
	"fmt"
	"github.com/lock14/collections"
	"iter"
	"strings"
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
	l := &LinkedList[T]{
		size: 0,
	}
	l.list.next = &l.list
	l.list.prev = &l.list
	return l
}

func (l *LinkedList[T]) lazyInit() {
	if l.list.next == nil {
		l.list.next = &l.list
		l.list.prev = &l.list
	}
}

// AddFront inserts the specified element at the front of the list.
func (l *LinkedList[T]) AddFront(t T) {
	l.lazyInit()
	insertBefore(l.list.next, t)
	l.size++
}

// RemoveFront removes and returns the element at the front of the list.
// Panics if the list is empty.
func (l *LinkedList[T]) RemoveFront() T {
	if l.Empty() {
		panic("cannot remove from an empty list")
	}
	n := l.list.next
	unlink(n)
	l.size--
	return n.data
}

// AddBack inserts the specified element at the back of the list.
func (l *LinkedList[T]) AddBack(t T) {
	l.lazyInit()
	insertBefore(&l.list, t)
	l.size++
}

// RemoveBack removes and returns the element at the back of the list.
// Panics if the list is empty.
func (l *LinkedList[T]) RemoveBack() T {
	if l.Empty() {
		panic("cannot remove from an empty list")
	}
	n := l.list.prev
	unlink(n)
	l.size--
	return n.data
}

// Peek returns the element at the front of the list without removing it.
// Panics if the list is empty.
func (l *LinkedList[T]) Peek() T {
	return l.PeekFront()
}

// PeekFront returns the element at the front of the list without removing it.
// Panics if the list is empty.
func (l *LinkedList[T]) PeekFront() T {
	if l.Empty() {
		panic("cannot peek from an empty list")
	}
	return l.list.next.data
}

// PeekBack returns the element at the back of the list without removing it.
// Panics if the list is empty.
func (l *LinkedList[T]) PeekBack() T {
	if l.Empty() {
		panic("cannot peek from an empty list")
	}
	return l.list.prev.data
}

// Add appends the specified element to the end of the list.
func (l *LinkedList[T]) Add(t T) {
	l.AddBack(t)
}

// Remove removes and returns the element at the front of the list.
// Panics if the list is empty.
func (l *LinkedList[T]) Remove() T {
	return l.RemoveFront()
}

// Push adds the specified element to the front of the list.
func (l *LinkedList[T]) Push(t T) {
	l.AddFront(t)
}

// Get returns the element at the specified index.
// Panics if index is out of bounds.
func (l *LinkedList[T]) Get(idx int) T {
	if n := l.get(idx); n != nil {
		return n.data
	}
	panic("index out of bounds")
}

// Set replaces the element at the specified index with the given element.
// Panics if index is out of bounds.
func (l *LinkedList[T]) Set(idx int, t T) {
	if n := l.get(idx); n != nil {
		n.data = t
		return
	}
	panic("index out of bounds")
}

// Pop removes and returns the element at the front of the list.
// Panics if the list is empty.
func (l *LinkedList[T]) Pop() T {
	return l.RemoveFront()
}

// Size returns the number of elements in the list.
func (l *LinkedList[T]) Size() int {
	return l.size
}

// AddAll inserts all elements from the given sequence into the list.
func (l *LinkedList[T]) AddAll(sequence iter.Seq[T]) {
	for t := range sequence {
		l.Add(t)
	}
}

// Empty returns true if the list contains no elements.
func (l *LinkedList[T]) Empty() bool {
	return l.Size() == 0
}

// Clear removes all elements from the list.
func (l *LinkedList[T]) Clear() {
	l.list.next = &l.list
	l.list.prev = &l.list
	l.size = 0
}

// String returns a string representation of the list matching Go slice formatting.
func (l *LinkedList[T]) String() string {
	str := make([]string, 0, l.Size())
	for t := range l.All() {
		str = append(str, fmt.Sprintf("%v", t))
	}
	return "[" + strings.Join(str, " ") + "]"
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

// Backward returns an iterator over all elements of the list in reverse order (from back to front).
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
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
}
