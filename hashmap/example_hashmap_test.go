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

func ExampleHashMap_Get() {
	m := hashmap.New[string, int]()
	m.Put("Alice", 25)

	v, ok := m.Get("Alice")
	fmt.Println("Alice:", v, "found:", ok)

	v, ok = m.Get("Bob")
	fmt.Println("Bob:", v, "found:", ok)

	// Output:
	// Alice: 25 found: true
	// Bob: 0 found: false
}

func ExampleHashMap_Put() {
	m := hashmap.New[string, int]()
	m.Put("Alice", 25)

	v, _ := m.Get("Alice")
	fmt.Println("Alice:", v)

	m.Put("Alice", 30) // Update
	v, _ = m.Get("Alice")
	fmt.Println("Alice updated:", v)

	// Output:
	// Alice: 25
	// Alice updated: 30
}

func ExampleHashMap_Remove() {
	m := hashmap.New[string, int]()
	m.Put("Alice", 25)

	m.Remove("Alice")
	_, ok := m.Get("Alice")
	fmt.Println("Alice found:", ok)

	// Output:
	// Alice found: false
}

func ExampleHashMap_Size() {
	m := hashmap.New[string, int]()
	fmt.Println("Size:", m.Size())

	m.Put("Alice", 25)
	fmt.Println("Size:", m.Size())

	// Output:
	// Size: 0
	// Size: 1
}

func ExampleHashMap_Empty() {
	m := hashmap.New[string, int]()
	fmt.Println("Empty:", m.Empty())

	m.Put("Alice", 25)
	fmt.Println("Empty:", m.Empty())

	// Output:
	// Empty: true
	// Empty: false
}

func ExampleHashMap_Clear() {
	m := hashmap.New[string, int]()
	m.Put("Alice", 25)
	m.Put("Bob", 30)

	m.Clear()
	fmt.Println("Size after clear:", m.Size())
	fmt.Println("Empty after clear:", m.Empty())

	// Output:
	// Size after clear: 0
	// Empty after clear: true
}

func ExampleHashMap_ContainsKey() {
	m := hashmap.New[string, int]()
	m.Put("Alice", 25)

	fmt.Println("Contains Alice:", m.ContainsKey("Alice"))
	fmt.Println("Contains Bob:", m.ContainsKey("Bob"))

	// Output:
	// Contains Alice: true
	// Contains Bob: false
}

func ExampleHashMap_All() {
	m := hashmap.New[string, int]()
	m.Put("Alice", 25)
	m.Put("Bob", 30)

	// Collect and sort to ensure deterministic output
	var entries []string
	for k, v := range m.All() {
		entries = append(entries, fmt.Sprintf("%s: %d", k, v))
	}
	slices.Sort(entries)

	for _, entry := range entries {
		fmt.Println(entry)
	}

	// Output:
	// Alice: 25
	// Bob: 30
}

func ExampleHashMap_Keys() {
	m := hashmap.New[string, int]()
	m.Put("Alice", 25)
	m.Put("Bob", 30)

	keys := slices.Collect(m.Keys())
	slices.Sort(keys)

	for _, k := range keys {
		fmt.Println(k)
	}

	// Output:
	// Alice
	// Bob
}

func ExampleHashMap_Values() {
	m := hashmap.New[string, int]()
	m.Put("Alice", 25)
	m.Put("Bob", 30)

	values := slices.Collect(m.Values())
	slices.Sort(values)

	for _, v := range values {
		fmt.Println(v)
	}

	// Output:
	// 25
	// 30
}

func ExampleHashMap_String() {
	m := hashmap.New[string, int]()
	fmt.Println(m.String())

	m.Put("Alice", 25)
	fmt.Println(m.String())

	// Output:
	// map[]
	// map[Alice:25]
}
