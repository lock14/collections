package arraydeque_test

import (
	"fmt"
	"slices"

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

func ExampleArrayDeque_Peek() {
	d := arraydeque.New[int]()
	d.Add(1)
	d.Add(2)
	fmt.Println(d.Peek())
	// Output:
	// 1
}

func ExampleArrayDeque_Add() {
	d := arraydeque.New[int]()
	d.Add(1)
	d.Add(2)
	fmt.Println(d.Size())
	// Output:
	// 2
}

func ExampleArrayDeque_Remove() {
	d := arraydeque.New[int]()
	d.Add(1)
	d.Add(2)
	fmt.Println(d.Remove())
	fmt.Println(d.Remove())
	// Output:
	// 1
	// 2
}

func ExampleArrayDeque_Push() {
	d := arraydeque.New[int]()
	d.Push(1)
	d.Push(2)
	fmt.Println(d.Size())
	// Output:
	// 2
}

func ExampleArrayDeque_Pop() {
	d := arraydeque.New[int]()
	d.Push(1)
	d.Push(2)
	fmt.Println(d.Pop())
	fmt.Println(d.Pop())
	// Output:
	// 2
	// 1
}

func ExampleArrayDeque_PeekFront() {
	d := arraydeque.New[int]()
	d.AddFront(1)
	d.AddFront(2)
	fmt.Println(d.PeekFront())
	// Output:
	// 2
}

func ExampleArrayDeque_AddFront() {
	d := arraydeque.New[int]()
	d.AddFront(1)
	d.AddFront(2)
	fmt.Println(d.Size())
	// Output:
	// 2
}

func ExampleArrayDeque_RemoveFront() {
	d := arraydeque.New[int]()
	d.AddFront(1)
	d.AddFront(2)
	fmt.Println(d.RemoveFront())
	fmt.Println(d.RemoveFront())
	// Output:
	// 2
	// 1
}

func ExampleArrayDeque_PeekBack() {
	d := arraydeque.New[int]()
	d.AddBack(1)
	d.AddBack(2)
	fmt.Println(d.PeekBack())
	// Output:
	// 2
}

func ExampleArrayDeque_AddBack() {
	d := arraydeque.New[int]()
	d.AddBack(1)
	d.AddBack(2)
	fmt.Println(d.Size())
	// Output:
	// 2
}

func ExampleArrayDeque_RemoveBack() {
	d := arraydeque.New[int]()
	d.AddBack(1)
	d.AddBack(2)
	fmt.Println(d.RemoveBack())
	fmt.Println(d.RemoveBack())
	// Output:
	// 2
	// 1
}

func ExampleArrayDeque_AddAll() {
	d := arraydeque.New[int]()
	d.AddAll(slices.Values([]int{1, 2, 3}))
	fmt.Println(d.Size())
	// Output:
	// 3
}

func ExampleArrayDeque_Size() {
	d := arraydeque.New[int]()
	fmt.Println(d.Size())
	d.Add(1)
	fmt.Println(d.Size())
	// Output:
	// 0
	// 1
}

func ExampleArrayDeque_Empty() {
	d := arraydeque.New[int]()
	fmt.Println(d.Empty())
	d.Add(1)
	fmt.Println(d.Empty())
	// Output:
	// true
	// false
}

func ExampleArrayDeque_Clear() {
	d := arraydeque.New[int]()
	d.Add(1)
	d.Add(2)
	d.Clear()
	fmt.Println(d.Size())
	fmt.Println(d.Empty())
	// Output:
	// 0
	// true
}

func ExampleArrayDeque_String() {
	d := arraydeque.New[int]()
	d.Add(1)
	d.Add(2)
	d.Add(3)
	fmt.Println(d.String())
	// Output:
	// [1 2 3]
}

func ExampleArrayDeque_All() {
	d := arraydeque.New[int]()
	d.Add(1)
	d.Add(2)
	d.Add(3)
	for v := range d.All() {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleArrayDeque_Backward() {
	d := arraydeque.New[int]()
	d.Add(1)
	d.Add(2)
	d.Add(3)
	for v := range d.Backward() {
		fmt.Println(v)
	}
	// Output:
	// 3
	// 2
	// 1
}
