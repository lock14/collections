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

func ExampleHashSet_Add() {
	set := hashset.New[string]()
	set.Add("Apple")
	fmt.Println(set.Size())
	// Output: 1
}

func ExampleHashSet_Remove() {
	set := hashset.New[string]()
	set.Add("Apple")
	item := set.Remove()
	fmt.Println(item)
	fmt.Println(set.Size())
	// Output:
	// Apple
	// 0
}

func ExampleHashSet_RemoveElement() {
	set := hashset.New[string]()
	set.Add("Apple")
	set.RemoveElement("Apple")
	fmt.Println(set.Size())
	// Output: 0
}

func ExampleHashSet_Contains() {
	set := hashset.New[string]()
	set.Add("Apple")
	fmt.Println(set.Contains("Apple"))
	fmt.Println(set.Contains("Banana"))
	// Output:
	// true
	// false
}

func ExampleHashSet_ContainsAll() {
	set1 := hashset.New[string]()
	set1.Add("Apple")
	set1.Add("Banana")

	set2 := hashset.New[string]()
	set2.Add("Apple")

	fmt.Println(set1.ContainsAll(set2))

	set2.Add("Cherry")
	fmt.Println(set1.ContainsAll(set2))
	// Output:
	// true
	// false
}

func ExampleHashSet_AddAll() {
	set := hashset.New[string]()
	elements := slices.Values([]string{"Apple", "Banana"})
	set.AddAll(elements)
	fmt.Println(set.Size())
	// Output: 2
}

func ExampleHashSet_RemoveAll() {
	set1 := hashset.New[string]()
	set1.Add("Apple")
	set1.Add("Banana")

	set2 := hashset.New[string]()
	set2.Add("Apple")

	set1.RemoveAll(set2)
	fmt.Println(set1.Contains("Banana"))
	fmt.Println(set1.Size())
	// Output:
	// true
	// 1
}

func ExampleHashSet_RetainAll() {
	set1 := hashset.New[string]()
	set1.Add("Apple")
	set1.Add("Banana")

	set2 := hashset.New[string]()
	set2.Add("Apple")
	set2.Add("Cherry")

	set1.RetainAll(set2)
	fmt.Println(set1.Contains("Apple"))
	fmt.Println(set1.Contains("Banana"))
	fmt.Println(set1.Size())
	// Output:
	// true
	// false
	// 1
}

func ExampleHashSet_Clear() {
	set := hashset.New[string]()
	set.Add("Apple")
	set.Clear()
	fmt.Println(set.Size())
	// Output: 0
}

func ExampleHashSet_Size() {
	set := hashset.New[string]()
	set.Add("Apple")
	fmt.Println(set.Size())
	// Output: 1
}

func ExampleHashSet_Empty() {
	set := hashset.New[string]()
	fmt.Println(set.Empty())
	set.Add("Apple")
	fmt.Println(set.Empty())
	// Output:
	// true
	// false
}

func ExampleHashSet_String() {
	set := hashset.New[string]()
	set.Add("Apple")
	fmt.Println(set.String())
	// Output: [Apple]
}

func ExampleHashSet_All() {
	set := hashset.New[string]()
	set.Add("Apple")
	set.Add("Banana")

	elements := slices.Collect(set.All())
	slices.Sort(elements)

	for _, e := range elements {
		fmt.Println(e)
	}
	// Output:
	// Apple
	// Banana
}
