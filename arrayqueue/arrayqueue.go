package arrayqueue

import (
	"fmt"
	"strings"
)

const (
	DefaultCapacity = 10
)

// ArrayQueue represents a queue of elements of type T backed by an array.
type ArrayQueue[T any] struct {
	slice []T
	front int
	back  int
	size  int
}

// Config holds the values for configuring a ArrayQueue.
type Config struct {
	Capacity int
}

// Option configures a ArrayQueue config
type Option func(*Config)

// New creates a empty ArrayQueue whose initial size is 0.
func New[T any](opts ...Option) *ArrayQueue[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &ArrayQueue[T]{
		slice: make([]T, config.Capacity),
	}
}

func (q *ArrayQueue[T]) Add(t T) {
	if q.size > 0 && q.back == q.front {
		q.resize()
	}
	q.slice[q.back] = t
	q.back = (q.back + 1) % len(q.slice)
	q.size++
}

func (q *ArrayQueue[T]) Remove() T {
	if q.isEmpty() {
		panic("cannot remove from an empty arrayqueue")
	}
	t := q.slice[q.front]
	q.front = (q.front + 1) % len(q.slice)
	q.size--
	return t
}

func (q *ArrayQueue[T]) Size() int {
	return q.size
}

func (q *ArrayQueue[T]) isEmpty() bool {
	return q.size == 0
}

func (q *ArrayQueue[T]) String() string {
	str := make([]string, q.Size())
	cur := q.front
	for i := 0; i < len(str); i++ {
		str[i] = fmt.Sprintf("%+v", q.slice[cur])
		cur = (cur + 1) % len(q.slice)
	}
	return "[" + strings.Join(str, ", ") + "]"
}

func (q *ArrayQueue[T]) ToSlice() []T {
	slice := make([]T, q.Size())
	for i := 0; i < cap(slice); i++ {
		t := q.Remove()
		slice = append(slice, t)
		q.Add(t)
	}
	return slice
}

func (q *ArrayQueue[T]) resize() {
	newCap := cap(q.slice) + (cap(q.slice) / 2)
	slice := make([]T, newCap)
	i := 0
	for !q.isEmpty() {
		slice[i] = q.Remove()
		i++
	}
	q.slice = slice
	q.front = 0
	q.back = i
	q.size = i
}

func defaultConfig() *Config {
	return &Config{
		Capacity: DefaultCapacity,
	}
}
