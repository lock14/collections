package arraylist_test

import (
	"fmt"
	"github.com/lock14/collections/arraylist"
)

func ExampleSliceWrapper() {
	// ArrayList is a dynamically resizing array, providing O(1) random access.
	list := arraylist.Wrap(make([]string, 0))

	list.Add("Apple")
	list.Add("Banana")
	list.Add("Cherry")

	// Set element at a specific index
	list.Set(1, "Blueberry")

	// Random access by index
	fmt.Println("Element at index 2:", list.Get(2))

	// Iteration
	fmt.Println("All fruits:")
	for val := range list.All() {
		fmt.Println("-", val)
	}

	// Output:
	// Element at index 2: Cherry
	// All fruits:
	// - Apple
	// - Blueberry
	// - Cherry
}
