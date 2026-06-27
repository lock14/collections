package treeset_test

import (
	"fmt"
	"slices"

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

func ExampleTreeSet_Add() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	fmt.Println(set.Contains(1))
	fmt.Println(set.Contains(3))
	// Output:
	// true
	// false
}

func ExampleTreeSet_Remove() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	item := set.Remove()
	fmt.Println(item)
	fmt.Println(set.Empty())
	// Output:
	// 1
	// true
}

func ExampleTreeSet_RemoveElement() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.RemoveElement(1)
	fmt.Println(set.Contains(1))
	fmt.Println(set.Contains(2))
	// Output:
	// false
	// true
}

func ExampleTreeSet_Contains() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	fmt.Println(set.Contains(1))
	fmt.Println(set.Contains(2))
	// Output:
	// true
	// false
}

func ExampleTreeSet_ContainsAll() {
	set1 := treeset.NewOrdered[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2 := treeset.NewOrdered[int]()
	set2.Add(1)
	set2.Add(2)

	set3 := treeset.NewOrdered[int]()
	set3.Add(1)
	set3.Add(4)

	fmt.Println(set1.ContainsAll(set2))
	fmt.Println(set1.ContainsAll(set3))
	// Output:
	// true
	// false
}

func ExampleTreeSet_AddAll() {
	set := treeset.NewOrdered[int]()
	set.AddAll(slices.Values([]int{1, 2, 3}))
	fmt.Println(set.Size())
	fmt.Println(set.Contains(2))
	// Output:
	// 3
	// true
}

func ExampleTreeSet_RemoveAll() {
	set1 := treeset.NewOrdered[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2 := treeset.NewOrdered[int]()
	set2.Add(2)
	set2.Add(3)

	set1.RemoveAll(set2)
	fmt.Println(set1.Contains(1))
	fmt.Println(set1.Contains(2))
	// Output:
	// true
	// false
}

func ExampleTreeSet_RetainAll() {
	set1 := treeset.NewOrdered[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2 := treeset.NewOrdered[int]()
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)

	set1.RetainAll(set2)
	fmt.Println(set1.Contains(1))
	fmt.Println(set1.Contains(2))
	fmt.Println(set1.Contains(3))
	fmt.Println(set1.Size())
	// Output:
	// false
	// true
	// true
	// 2
}

func ExampleTreeSet_Clear() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Clear()
	fmt.Println(set.Empty())
	// Output:
	// true
}

func ExampleTreeSet_Size() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	fmt.Println(set.Size())
	// Output:
	// 2
}

func ExampleTreeSet_Empty() {
	set := treeset.NewOrdered[int]()
	fmt.Println(set.Empty())
	set.Add(1)
	fmt.Println(set.Empty())
	// Output:
	// true
	// false
}

func ExampleTreeSet_All() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	for v := range set.All() {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleTreeSet_First() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println(set.First())
	// Output:
	// 1
}

func ExampleTreeSet_Last() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println(set.Last())
	// Output:
	// 3
}

func ExampleTreeSet_PollFirst() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println(set.PollFirst())
	fmt.Println(set.Size())
	// Output:
	// 1
	// 2
}

func ExampleTreeSet_PollLast() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println(set.PollLast())
	fmt.Println(set.Size())
	// Output:
	// 3
	// 2
}

func ExampleTreeSet_AddFirst() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	set := treeset.NewOrdered[int]()
	set.AddFirst(1)
	// Output:
	// AddFirst is not supported on SortedSet
}

func ExampleTreeSet_AddLast() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	set := treeset.NewOrdered[int]()
	set.AddLast(1)
	// Output:
	// AddLast is not supported on SortedSet
}

func ExampleTreeSet_Lower() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(3)
	set.Add(5)

	val, ok := set.Lower(4)
	fmt.Println(val, ok)

	val, ok = set.Lower(1)
	fmt.Println(val, ok)
	// Output:
	// 3 true
	// 0 false
}

func ExampleTreeSet_Floor() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(3)
	set.Add(5)

	val, ok := set.Floor(4)
	fmt.Println(val, ok)

	val, ok = set.Floor(3)
	fmt.Println(val, ok)

	val, ok = set.Floor(0)
	fmt.Println(val, ok)
	// Output:
	// 3 true
	// 3 true
	// 0 false
}

func ExampleTreeSet_Ceiling() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(3)
	set.Add(5)

	val, ok := set.Ceiling(2)
	fmt.Println(val, ok)

	val, ok = set.Ceiling(3)
	fmt.Println(val, ok)

	val, ok = set.Ceiling(6)
	fmt.Println(val, ok)
	// Output:
	// 3 true
	// 3 true
	// 0 false
}

func ExampleTreeSet_Higher() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(3)
	set.Add(5)

	val, ok := set.Higher(2)
	fmt.Println(val, ok)

	val, ok = set.Higher(5)
	fmt.Println(val, ok)
	// Output:
	// 3 true
	// 0 false
}

func ExampleTreeSet_Backward() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	for v := range set.Backward() {
		fmt.Println(v)
	}
	// Output:
	// 3
	// 2
	// 1
}

func ExampleTreeSet_From() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(4)
	set.Add(5)

	for v := range set.From(3) {
		fmt.Println(v)
	}
	// Output:
	// 3
	// 4
	// 5
}

func ExampleTreeSet_To() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(4)
	set.Add(5)

	for v := range set.To(4) {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleTreeSet_Between() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(4)
	set.Add(5)

	for v := range set.Between(2, 5) {
		fmt.Println(v)
	}
	// Output:
	// 2
	// 3
	// 4
}

func ExampleTreeSet_String() {
	set := treeset.NewOrdered[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println(set.String())
	// Output:
	// [1, 2, 3]
}
