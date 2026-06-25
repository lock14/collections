package heap_test

import (
	"cmp"
	"fmt"
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
