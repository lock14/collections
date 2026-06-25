package linkedhashmap_test

import (
	"fmt"
	"github.com/lock14/collections/linkedhashmap"
)

func ExampleLinkedHashMap() {
	// LinkedHashMap is a hash map that preserves insertion order.
	m := linkedhashmap.New[string, string]()

	m.Put("First", "A")
	m.Put("Second", "B")
	m.Put("Third", "C")

	// Updating an existing key does not change its position in the insertion order
	m.Put("First", "A+")

	fmt.Println("Iterating over LinkedHashMap:")
	for k, v := range m.All() {
		fmt.Printf("- %s: %s\n", k, v)
	}

	// Output:
	// Iterating over LinkedHashMap:
	// - First: A+
	// - Second: B
	// - Third: C
}
