package arraydeque_test

import (
	"fmt"
	"github.com/lock14/collections/arraydeque"
)

func ExampleArrayDeque_queue() {
	// ArrayDeque implements a fast, ring-buffer backed Queue
	// Elements are added to the back and removed from the front (FIFO).
	queue := arraydeque.New[string]()

	queue.AddBack("Alice")
	queue.AddBack("Bob")
	queue.AddBack("Charlie")

	for !queue.Empty() {
		person := queue.RemoveFront()
		fmt.Println("Served:", person)
	}

	// Output:
	// Served: Alice
	// Served: Bob
	// Served: Charlie
}

func ExampleArrayDeque_stack() {
	// ArrayDeque can also be used efficiently as a Stack
	// Elements are added to the back and removed from the back (LIFO).
	stack := arraydeque.New[string]()

	stack.Push("Page 1")
	stack.Push("Page 2")
	stack.Push("Page 3")

	for !stack.Empty() {
		page := stack.Pop()
		fmt.Println("Navigating back from:", page)
	}

	// Output:
	// Navigating back from: Page 3
	// Navigating back from: Page 2
	// Navigating back from: Page 1
}
