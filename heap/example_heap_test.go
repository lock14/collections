package heap_test

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/lock14/collections/heap"
)

func ExampleHeap_min() {
	// A basic Min-Heap for ordered types
	h := heap.Min[int]()

	h.Add(10)
	h.Add(5)
	h.Add(20)
	h.Add(1)

	fmt.Println("Elements in ascending order:")
	for !h.Empty() {
		fmt.Println("-", h.Remove())
	}

	// Output:
	// Elements in ascending order:
	// - 1
	// - 5
	// - 10
	// - 20
}

func ExampleHeap_priorityQueue() {
	// A custom Priority Queue using a struct and a custom comparator
	type Task struct {
		Name     string
		Priority int // Higher is more important
	}

	// Custom comparator: We want higher priorities first (Max-Heap by Priority)
	pq := heap.New[Task](heap.WithComparator(func(a, b Task) int {
		// To simulate a Max-Heap, we reverse the comparison: b.Priority - a.Priority
		return cmp.Compare(b.Priority, a.Priority)
	}))

	pq.Add(Task{Name: "Fix typo", Priority: 1})
	pq.Add(Task{Name: "Database outage", Priority: 100})
	pq.Add(Task{Name: "Update dependencies", Priority: 5})

	fmt.Println("Tasks processed by priority:")
	for !pq.Empty() {
		task := pq.Remove()
		fmt.Printf("- %s (Priority: %d)\n", task.Name, task.Priority)
	}

	// Output:
	// Tasks processed by priority:
	// - Database outage (Priority: 100)
	// - Update dependencies (Priority: 5)
	// - Fix typo (Priority: 1)
}

func ExampleHeap_Add() {
	h := heap.Min[int]()
	h.Add(42)
	h.Add(10)

	fmt.Println(h.Remove())
	fmt.Println(h.Remove())

	// Output:
	// 10
	// 42
}

func ExampleHeap_AddAll() {
	h := heap.Min[int]()
	h.AddAll(slices.Values([]int{30, 10, 20}))

	for !h.Empty() {
		fmt.Println(h.Remove())
	}

	// Output:
	// 10
	// 20
	// 30
}

func ExampleHeap_Remove() {
	h := heap.Min[int]()
	h.Add(100)
	h.Add(50)

	fmt.Println(h.Remove())

	// Output:
	// 50
}

func ExampleHeap_Peek() {
	h := heap.Min[int]()
	h.Add(7)
	h.Add(3)

	fmt.Println(h.Peek())
	fmt.Println(h.Size())

	// Output:
	// 3
	// 2
}

func ExampleHeap_Size() {
	h := heap.Min[int]()
	fmt.Println(h.Size())
	h.Add(1)
	fmt.Println(h.Size())

	// Output:
	// 0
	// 1
}

func ExampleHeap_Empty() {
	h := heap.Min[int]()
	fmt.Println(h.Empty())
	h.Add(1)
	fmt.Println(h.Empty())

	// Output:
	// true
	// false
}

func ExampleHeap_Clear() {
	h := heap.Min[int]()
	h.Add(1)
	h.Add(2)
	fmt.Println(h.Size())
	h.Clear()
	fmt.Println(h.Size())
	fmt.Println(h.Empty())

	// Output:
	// 2
	// 0
	// true
}

func ExampleHeap_All() {
	h := heap.Min[int]()
	h.Add(3)
	h.Add(1)
	h.Add(2)

	// The elements are returned in their current heap array order.
	for v := range h.All() {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 3
	// 2
}
