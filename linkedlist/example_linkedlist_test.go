package linked_list_test

import (
	"fmt"
	"github.com/lock14/collections/linkedlist"
)

func ExampleLinkedList() {
	// LinkedList is a doubly-linked list, providing fast O(1) insertions
	// and deletions at the front and back.
	list := linked_list.New[string]()

	list.AddBack("Middle")
	list.AddFront("Start")
	list.AddBack("End")

	// Set elements at a specific index (requires O(N) traversal to index)
	list.Set(1, "Second")

	// Iteration
	for val := range list.All() {
		fmt.Println(val)
	}

	// Output:
	// Start
	// Second
	// End
}
