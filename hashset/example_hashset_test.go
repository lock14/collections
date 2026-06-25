package hashset_test

import (
	"fmt"
	"github.com/lock14/collections/hashset"
	"slices"
)

func ExampleHashSet() {
	// HashSet provides an unordered collection of unique elements
	// with O(1) average time complexity for Add, Remove, and Contains.
	set := hashset.New[string]()

	set.Add("Apple")
	set.Add("Banana")
	set.Add("Cherry")
	set.Add("Apple") // Ignored, already exists

	fmt.Println("Contains Apple:", set.Contains("Apple"))
	fmt.Println("Size of set:", set.Size())

	// Since HashSet is unordered, we collect and sort to ensure deterministic output
	elements := slices.Collect(set.All())
	slices.Sort(elements)

	fmt.Println("Elements:")
	for _, e := range elements {
		fmt.Println("-", e)
	}

	// Output:
	// Contains Apple: true
	// Size of set: 3
	// Elements:
	// - Apple
	// - Banana
	// - Cherry
}
