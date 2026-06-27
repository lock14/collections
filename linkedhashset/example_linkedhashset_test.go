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

func ExampleLinkedHashSet_Add() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")
	set.Add("C")

	fmt.Println(set.Size())
	fmt.Println(set.Contains("B"))

	// Output:
	// 3
	// true
}

func ExampleLinkedHashSet_Remove() {
	set := linkedhashset.New[string]()
	set.Add("A")

	item := set.Remove()
	fmt.Println(item)
	fmt.Println(set.Empty())

	// Output:
	// A
	// true
}

func ExampleLinkedHashSet_First() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")

	fmt.Println(set.First())

	// Output:
	// A
}

func ExampleLinkedHashSet_Last() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")

	fmt.Println(set.Last())

	// Output:
	// B
}

func ExampleLinkedHashSet_PollFirst() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")

	item := set.PollFirst()
	fmt.Println(item)
	fmt.Println(set.Size())

	// Output:
	// A
	// 1
}

func ExampleLinkedHashSet_PollLast() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")

	item := set.PollLast()
	fmt.Println(item)
	fmt.Println(set.Size())

	// Output:
	// B
	// 1
}

func ExampleLinkedHashSet_AddFirst() {
	set := linkedhashset.New[string]()
	set.Add("B")
	set.Add("C")
	set.AddFirst("A")

	for e := range set.All() {
		fmt.Println(e)
	}

	// Output:
	// A
	// B
	// C
}

func ExampleLinkedHashSet_AddLast() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")
	set.AddLast("C")

	for e := range set.All() {
		fmt.Println(e)
	}

	// Output:
	// A
	// B
	// C
}

func ExampleLinkedHashSet_RemoveElement() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")
	set.Add("C")

	set.RemoveElement("B")
	for e := range set.All() {
		fmt.Println(e)
	}

	// Output:
	// A
	// C
}

func ExampleLinkedHashSet_Contains() {
	set := linkedhashset.New[string]()
	set.Add("A")

	fmt.Println(set.Contains("A"))
	fmt.Println(set.Contains("B"))

	// Output:
	// true
	// false
}

func ExampleLinkedHashSet_ContainsAll() {
	set1 := linkedhashset.New[string]()
	set1.Add("A")
	set1.Add("B")
	set1.Add("C")

	set2 := linkedhashset.New[string]()
	set2.Add("A")
	set2.Add("B")

	set3 := linkedhashset.New[string]()
	set3.Add("B")
	set3.Add("D")

	fmt.Println(set1.ContainsAll(set2))
	fmt.Println(set1.ContainsAll(set3))

	// Output:
	// true
	// false
}

func ExampleLinkedHashSet_AddAll() {
	set := linkedhashset.New[string]()
	set.Add("A")

	set2 := linkedhashset.New[string]()
	set2.Add("B")
	set2.Add("C")

	set.AddAll(set2.All())

	for e := range set.All() {
		fmt.Println(e)
	}

	// Output:
	// A
	// B
	// C
}

func ExampleLinkedHashSet_RemoveAll() {
	set1 := linkedhashset.New[string]()
	set1.Add("A")
	set1.Add("B")
	set1.Add("C")

	set2 := linkedhashset.New[string]()
	set2.Add("B")

	set1.RemoveAll(set2)

	for e := range set1.All() {
		fmt.Println(e)
	}

	// Output:
	// A
	// C
}

func ExampleLinkedHashSet_RetainAll() {
	set1 := linkedhashset.New[string]()
	set1.Add("A")
	set1.Add("B")
	set1.Add("C")

	set2 := linkedhashset.New[string]()
	set2.Add("B")
	set2.Add("C")
	set2.Add("D")

	set1.RetainAll(set2)

	for e := range set1.All() {
		fmt.Println(e)
	}

	// Output:
	// B
	// C
}

func ExampleLinkedHashSet_Clear() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Clear()

	fmt.Println(set.Size())
	fmt.Println(set.Empty())

	// Output:
	// 0
	// true
}

func ExampleLinkedHashSet_Size() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")

	fmt.Println(set.Size())

	// Output:
	// 2
}

func ExampleLinkedHashSet_Empty() {
	set := linkedhashset.New[string]()
	fmt.Println(set.Empty())

	set.Add("A")
	fmt.Println(set.Empty())

	// Output:
	// true
	// false
}

func ExampleLinkedHashSet_All() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")
	set.Add("C")

	for e := range set.All() {
		fmt.Println(e)
	}

	// Output:
	// A
	// B
	// C
}

func ExampleLinkedHashSet_Backward() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")
	set.Add("C")

	for e := range set.Backward() {
		fmt.Println(e)
	}

	// Output:
	// C
	// B
	// A
}

func ExampleLinkedHashSet_String() {
	set := linkedhashset.New[string]()
	set.Add("A")
	set.Add("B")
	set.Add("C")

	fmt.Println(set.String())

	// Output:
	// [A B C]
}
