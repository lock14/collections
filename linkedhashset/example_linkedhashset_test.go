package linkedhashset_test

import (
	"fmt"
	"github.com/lock14/collections/linkedhashset"
)

func ExampleLinkedHashSet() {
	// LinkedHashSet provides an ordered collection of unique elements.
	// It guarantees that iteration strictly follows insertion order.
	set := linkedhashset.New[string]()

	set.Add("First")
	set.Add("Second")
	set.Add("Third")

	// Adding an existing element does not change its position in the insertion order
	set.Add("First")

	fmt.Println("Iterating over LinkedHashSet:")
	for e := range set.All() {
		fmt.Println("-", e)
	}

	// Output:
	// Iterating over LinkedHashSet:
	// - First
	// - Second
	// - Third
}
