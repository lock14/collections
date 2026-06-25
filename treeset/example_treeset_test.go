package treeset_test

import (
	"fmt"
	"github.com/lock14/collections/treeset"
)

func ExampleTreeSet_navigable() {
	// TreeSet is backed by a B-Tree and implements the NavigableSet interface.
	// Elements are kept sorted and can be queried based on proximity.
	set := treeset.NewOrdered[int]()

	set.Add(10)
	set.Add(20)
	set.Add(30)
	set.Add(40)
	set.Add(50)

	fmt.Println("First element:", set.First())
	fmt.Println("Last element:", set.Last())

	// Floor finds the greatest element less than or equal to the given element
	floor, _ := set.Floor(25)
	fmt.Println("Floor(25):", floor)

	// Ceiling finds the least element greater than or equal to the given element
	ceiling, _ := set.Ceiling(25)
	fmt.Println("Ceiling(25):", ceiling)

	// Iterate over elements between bounds (inclusive min, exclusive max)
	fmt.Println("Elements between 20 and 50:")
	for v := range set.Between(20, 50) {
		fmt.Println("-", v)
	}

	// Output:
	// First element: 10
	// Last element: 50
	// Floor(25): 20
	// Ceiling(25): 30
	// Elements between 20 and 50:
	// - 20
	// - 30
	// - 40
}
