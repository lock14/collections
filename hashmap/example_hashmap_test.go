package hashmap_test

import (
	"fmt"
	"github.com/lock14/collections/hashmap"
	"slices"
)

func ExampleHashMap() {
	// HashMap provides O(1) average time complexity for key-value lookups.
	m := hashmap.New[string, int]()

	m.Put("Alice", 25)
	m.Put("Bob", 30)
	m.Put("Charlie", 35)

	age, ok := m.Get("Bob")
	fmt.Printf("Bob is %d years old (found: %v)\n", age, ok)

	// Update an existing key
	m.Put("Alice", 26)

	// Since HashMap is unordered, we collect and sort to ensure deterministic output
	keys := slices.Collect(m.Keys())
	slices.Sort(keys)

	fmt.Println("All Entries:")
	for _, k := range keys {
		v, _ := m.Get(k)
		fmt.Printf("- %s: %d\n", k, v)
	}

	// Output:
	// Bob is 30 years old (found: true)
	// All Entries:
	// - Alice: 26
	// - Bob: 30
	// - Charlie: 35
}
