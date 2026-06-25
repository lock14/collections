package treemap_test

import (
	"fmt"
	"github.com/lock14/collections/treemap"
)

func ExampleTreeMap() {
	// TreeMap is a map backed by a B-Tree, maintaining entries sorted by key.
	// It implements the NavigableMap interface.
	m := treemap.NewOrdered[int, string]()

	m.Put(10, "A")
	m.Put(20, "B")
	m.Put(30, "C")
	m.Put(40, "D")
	m.Put(50, "E")

	// Iteration is always in sorted key order
	fmt.Println("Entries:")
	for k, v := range m.All() {
		fmt.Printf("- %d: %s\n", k, v)
	}

	// Navigable lookups
	k, v, _ := m.Floor(25)
	fmt.Printf("Floor(25): %d -> %s\n", k, v)

	k, v, _ = m.Ceiling(25)
	fmt.Printf("Ceiling(25): %d -> %s\n", k, v)

	// Output:
	// Entries:
	// - 10: A
	// - 20: B
	// - 30: C
	// - 40: D
	// - 50: E
	// Floor(25): 20 -> B
	// Ceiling(25): 30 -> C
}
