package arraydeque

import (
	"fmt"
	"github.com/lock14/collections/iterator"
	"strings"
)

const (
	DefaultCapacity = 10
)

// ArrayDeque represents a deque of elements of type T backed by an array.
type ArrayDeque[T any] struct {
	slice []T
	front int
	back  int
	size  int
}

// Config holds the values for configuring a ArrayDeque.
type Config struct {
	Capacity int
}

// Option configures a ArrayDeque config
type Option func(*Config)

// New creates an empty ArrayDeque whose initial size is 0.
func New[T any](opts ...Option) *ArrayDeque[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &ArrayDeque[T]{
		slice: make([]T, config.Capacity),
	}
}

// Add is an alias
func (d *ArrayDeque[T]) Add(t T) {
	d.AddBack(t)
}

// Remove is an alias
func (d *ArrayDeque[T]) Remove() T {
	return d.RemoveFront()
}

// Push is an alias
func (d *ArrayDeque[T]) Push(t T) {
	d.AddFront(t)
}

// Pop is an alias
func (d *ArrayDeque[T]) Pop() T {
	return d.RemoveFront()
}

func (d *ArrayDeque[T]) AddFront(t T) {
	if d.size > 0 && d.front == d.back {
		d.resize()
	}
	index := (d.front - 1) % len(d.slice)
	if index < 0 {
		index += len(d.slice)
	}
	d.slice[index] = t
	d.front = index
	d.size++
}

func (d *ArrayDeque[T]) RemoveFront() T {
	if d.isEmpty() {
		panic("cannot remove from an empty ArrayDeque")
	}
	var zero T
	t := d.slice[d.front]
	d.slice[d.front] = zero
	d.front = (d.front + 1) % len(d.slice)
	d.size--
	return t
}

func (d *ArrayDeque[T]) AddBack(t T) {
	if d.size > 0 && d.front == d.back {
		d.resize()
	}
	d.slice[d.back] = t
	d.back = (d.back + 1) % len(d.slice)
	d.size++
}

func (d *ArrayDeque[T]) RemoveBack() T {
	if d.isEmpty() {
		panic("cannot remove from an empty ArrayDeque")
	}
	var zero T
	index := (d.back - 1) % len(d.slice)
	t := d.slice[index]
	d.slice[index] = zero
	d.back = index
	d.size--
	return t
}

func (d *ArrayDeque[T]) Size() int {
	return d.size
}

func (d *ArrayDeque[T]) isEmpty() bool {
	return d.size == 0
}

func (d *ArrayDeque[T]) String() string {
	str := make([]string, 0, d.Size())
	for t := range d.Elements() {
		str = append(str, fmt.Sprintf("%+v", t))
	}
	return "[" + strings.Join(str, ", ") + "]"
}

func (d *ArrayDeque[T]) Iterator() iterator.Iterator[T] {
	return &dequeIterator[T]{
		deque: d,
		index: 0,
	}
}

func (d *ArrayDeque[T]) Elements() chan T {
	return iterator.Elements(d.Iterator())
}

func (d *ArrayDeque[T]) ToSlice() []T {
	slice := make([]T, d.Size())
	for t := range d.Elements() {
		slice = append(slice, t)
	}
	return slice
}

func (d *ArrayDeque[T]) resize() {
	newCap := cap(d.slice) + (cap(d.slice) / 2)
	slice := make([]T, newCap)
	i := 0
	for t := range d.Elements() {
		slice[i] = t
		i++
	}
	d.slice = slice
	d.front = 0
	d.back = i
	d.size = i
}

func defaultConfig() *Config {
	return &Config{
		Capacity: DefaultCapacity,
	}
}

// Iterator

type dequeIterator[T any] struct {
	deque *ArrayDeque[T]
	index int
}

func (itr *dequeIterator[T]) Empty() bool {
	return itr.index >= itr.deque.Size()
}

func (itr *dequeIterator[T]) Next() (T, error) {
	var t T
	if itr.Empty() {
		return t, fmt.Errorf("cannot call Next() on an empty Iterator")
	}
	i := (itr.index + itr.deque.front) % len(itr.deque.slice)
	t = itr.deque.slice[i]
	itr.index++
	return t, nil
}
