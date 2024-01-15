package queue

import (
	"fmt"
	"strings"
)

const (
	DefaultCapacity = 10
)

// Queue[T] represents a queue of elements of type T
type Queue[T any] struct {
	slice []T
	front int
	back  int
	size  int
}

// Config holds the values for configuring a Queue.
type Config struct {
	Capacity int
}

// Option configures a Queue config
type Option func(*Config)

// New creates a empty Queue whose initial size is 0.
func New[T any](opts ...Option) *Queue[T] {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &Queue[T]{
		slice: make([]T, config.Capacity),
	}
}

func (q *Queue[T]) Add(t T) {
	if q.size > 0 && q.back == q.front {
		q.resize()
	}
	q.slice[q.back] = t
	q.back = (q.back + 1) % len(q.slice)
	q.size++
}

func (q *Queue[T]) Remove() T {
	if q.isEmpty() {
		panic("cannot remove from an empty queue")
	}
	t := q.slice[q.front]
	q.front = (q.front + 1) % len(q.slice)
	q.size--
	return t
}

func (q *Queue[T]) Size() int {
	return q.size
}

func (q *Queue[T]) isEmpty() bool {
	return q.size == 0
}

func (q *Queue[T]) String() string {
	str := make([]string, q.Size())
	cur := q.front
	for i := 0; i < len(str); i++ {
		str[i] = fmt.Sprintf("%+v", q.slice[cur])
		cur = (cur + 1) % len(q.slice)
	}
	return "[" + strings.Join(str, ", ") + "]"
}

func (q *Queue[T]) resize() {
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
