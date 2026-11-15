package comparator_test

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/lock14/collections/comparator"
)

func ExampleReverse() {
	// Let's say we have some numbers we want to sort
	numbers := []int{5, 2, 9, 1, 7}

	// comparator.Reverse takes any existing comparator and reverses its ordering.
	// We can use it to reverse NaturalOrder to sort in descending order.
	descendingCmp := comparator.Reverse(comparator.NaturalOrder[int]())

	slices.SortFunc(numbers, descendingCmp)

	fmt.Println("Descending:", numbers)

	// Custom comparators can also be reversed
	type Person struct {
		Age int
	}
	people := []Person{{Age: 30}, {Age: 20}, {Age: 40}}

	// Compare by Age ascending
	ageAsc := func(a, b Person) int {
		return cmp.Compare(a.Age, b.Age)
	}

	// Sort by Age descending
	slices.SortFunc(people, comparator.Reverse(ageAsc))

	fmt.Println("People by Age Descending:")
	for _, p := range people {
		fmt.Println("-", p.Age)
	}

	// Output:
	// Descending: [9 7 5 2 1]
	// People by Age Descending:
	// - 40
	// - 30
	// - 20
}

func ExampleNaturalOrder() {
	numbers := []int{5, 2, 9, 1, 7}

	// comparator.NaturalOrder returns a comparator that sorts in ascending order
	// based on the natural ordering of cmp.Ordered types.
	slices.SortFunc(numbers, comparator.NaturalOrder[int]())

	fmt.Println("Ascending:", numbers)

	words := []string{"banana", "apple", "cherry"}
	slices.SortFunc(words, comparator.NaturalOrder[string]())

	fmt.Println("Strings Ascending:", words)

	// Output:
	// Ascending: [1 2 5 7 9]
	// Strings Ascending: [apple banana cherry]
}
