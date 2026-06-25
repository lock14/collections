package comparator_test

import (
	"cmp"
	"fmt"
	"github.com/lock14/collections/comparator"
	"slices"
)

func ExampleReverse() {
	// Let's say we have some numbers we want to sort
	numbers := []int{5, 2, 9, 1, 7}

	// comparator.Reverse takes any existing comparator and negates its output.
	// We can use it to reverse the NaturalOrder to sort in descending order.
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
