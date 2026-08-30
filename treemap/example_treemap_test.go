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

func ExampleTreeMap_Get() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")

	v, ok := m.Get(1)
	fmt.Printf("Get(1): %s, %t\n", v, ok)

	v, ok = m.Get(2)
	fmt.Printf("Get(2): %s, %t\n", v, ok)

	// Output:
	// Get(1): A, true
	// Get(2): , false
}

func ExampleTreeMap_Put() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")

	fmt.Println(m.Size())

	// Output:
	// 2
}

func ExampleTreeMap_Remove() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")

	m.Remove(1)
	fmt.Println(m.ContainsKey(1))
	fmt.Println(m.Size())

	// Output:
	// false
	// 1
}

func ExampleTreeMap_Size() {
	m := treemap.NewOrdered[int, string]()
	fmt.Println(m.Size())
	m.Put(1, "A")
	fmt.Println(m.Size())

	// Output:
	// 0
	// 1
}

func ExampleTreeMap_Empty() {
	m := treemap.NewOrdered[int, string]()
	fmt.Println(m.Empty())
	m.Put(1, "A")
	fmt.Println(m.Empty())

	// Output:
	// true
	// false
}

func ExampleTreeMap_Clear() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	m.Clear()

	fmt.Println(m.Empty())
	fmt.Println(m.Size())

	// Output:
	// true
	// 0
}

func ExampleTreeMap_ContainsKey() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")

	fmt.Println(m.ContainsKey(1))
	fmt.Println(m.ContainsKey(2))

	// Output:
	// true
	// false
}

func ExampleTreeMap_All() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	m.Put(3, "C")

	for k, v := range m.All() {
		fmt.Printf("%d: %s\n", k, v)
	}

	// Output:
	// 1: A
	// 2: B
	// 3: C
}

func ExampleTreeMap_Keys() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	m.Put(3, "C")

	for k := range m.Keys() {
		fmt.Println(k)
	}

	// Output:
	// 1
	// 2
	// 3
}

func ExampleTreeMap_Values() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	m.Put(3, "C")

	for v := range m.Values() {
		fmt.Println(v)
	}

	// Output:
	// A
	// B
	// C
}

func ExampleTreeMap_First() {
	m := treemap.NewOrdered[int, string]()
	m.Put(2, "B")
	m.Put(1, "A")
	m.Put(3, "C")

	k, v := m.First()
	fmt.Printf("%d: %s\n", k, v)

	// Output:
	// 1: A
}

func ExampleTreeMap_Last() {
	m := treemap.NewOrdered[int, string]()
	m.Put(2, "B")
	m.Put(1, "A")
	m.Put(3, "C")

	k, v := m.Last()
	fmt.Printf("%d: %s\n", k, v)

	// Output:
	// 3: C
}

func ExampleTreeMap_PollFirst() {
	m := treemap.NewOrdered[int, string]()
	m.Put(2, "B")
	m.Put(1, "A")
	m.Put(3, "C")

	k, v := m.PollFirst()
	fmt.Printf("Polled %d: %s\n", k, v)
	fmt.Printf("New First: %d\n", func() int { k, _ := m.First(); return k }())

	// Output:
	// Polled 1: A
	// New First: 2
}

func ExampleTreeMap_PollLast() {
	m := treemap.NewOrdered[int, string]()
	m.Put(2, "B")
	m.Put(1, "A")
	m.Put(3, "C")

	k, v := m.PollLast()
	fmt.Printf("Polled %d: %s\n", k, v)
	fmt.Printf("New Last: %d\n", func() int { k, _ := m.Last(); return k }())

	// Output:
	// Polled 3: C
	// New Last: 2
}

func ExampleTreeMap_PutFirst() {
	m := treemap.NewOrdered[int, string]()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	m.PutFirst(1, "A")

	// Output:
	// PutFirst is not supported on SortedMap
}

func ExampleTreeMap_PutLast() {
	m := treemap.NewOrdered[int, string]()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	m.PutLast(1, "A")

	// Output:
	// PutLast is not supported on SortedMap
}

func ExampleTreeMap_Lower() {
	m := treemap.NewOrdered[int, string]()
	m.Put(10, "A")
	m.Put(20, "B")
	m.Put(30, "C")

	k, v, ok := m.Lower(20)
	fmt.Printf("Lower(20): %d -> %s, %t\n", k, v, ok)

	k, v, ok = m.Lower(10)
	fmt.Printf("Lower(10): %d -> %s, %t\n", k, v, ok)

	// Output:
	// Lower(20): 10 -> A, true
	// Lower(10): 0 -> , false
}

func ExampleTreeMap_Floor() {
	m := treemap.NewOrdered[int, string]()
	m.Put(10, "A")
	m.Put(20, "B")
	m.Put(30, "C")

	k, v, ok := m.Floor(25)
	fmt.Printf("Floor(25): %d -> %s, %t\n", k, v, ok)

	k, v, ok = m.Floor(20)
	fmt.Printf("Floor(20): %d -> %s, %t\n", k, v, ok)

	// Output:
	// Floor(25): 20 -> B, true
	// Floor(20): 20 -> B, true
}

func ExampleTreeMap_Ceiling() {
	m := treemap.NewOrdered[int, string]()
	m.Put(10, "A")
	m.Put(20, "B")
	m.Put(30, "C")

	k, v, ok := m.Ceiling(25)
	fmt.Printf("Ceiling(25): %d -> %s, %t\n", k, v, ok)

	k, v, ok = m.Ceiling(30)
	fmt.Printf("Ceiling(30): %d -> %s, %t\n", k, v, ok)

	// Output:
	// Ceiling(25): 30 -> C, true
	// Ceiling(30): 30 -> C, true
}

func ExampleTreeMap_Higher() {
	m := treemap.NewOrdered[int, string]()
	m.Put(10, "A")
	m.Put(20, "B")
	m.Put(30, "C")

	k, v, ok := m.Higher(20)
	fmt.Printf("Higher(20): %d -> %s, %t\n", k, v, ok)

	k, v, ok = m.Higher(30)
	fmt.Printf("Higher(30): %d -> %s, %t\n", k, v, ok)

	// Output:
	// Higher(20): 30 -> C, true
	// Higher(30): 0 -> , false
}

func ExampleTreeMap_Backward() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	m.Put(3, "C")

	for k, v := range m.Backward() {
		fmt.Printf("%d: %s\n", k, v)
	}

	// Output:
	// 3: C
	// 2: B
	// 1: A
}

func ExampleTreeMap_BackwardKeys() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	m.Put(3, "C")

	for k := range m.BackwardKeys() {
		fmt.Println(k)
	}

	// Output:
	// 3
	// 2
	// 1
}

func ExampleTreeMap_BackwardValues() {
	m := treemap.NewOrdered[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	m.Put(3, "C")

	for v := range m.BackwardValues() {
		fmt.Println(v)
	}

	// Output:
	// C
	// B
	// A
}

func ExampleTreeMap_From() {
	m := treemap.NewOrdered[int, string]()
	m.Put(10, "A")
	m.Put(20, "B")
	m.Put(30, "C")
	m.Put(40, "D")

	for k, v := range m.From(20) {
		fmt.Printf("%d: %s\n", k, v)
	}

	// Output:
	// 20: B
	// 30: C
	// 40: D
}

func ExampleTreeMap_To() {
	m := treemap.NewOrdered[int, string]()
	m.Put(10, "A")
	m.Put(20, "B")
	m.Put(30, "C")
	m.Put(40, "D")

	for k, v := range m.To(30) {
		fmt.Printf("%d: %s\n", k, v)
	}

	// Output:
	// 10: A
	// 20: B
}

func ExampleTreeMap_Between() {
	m := treemap.NewOrdered[int, string]()
	m.Put(10, "A")
	m.Put(20, "B")
	m.Put(30, "C")
	m.Put(40, "D")
	m.Put(50, "E")

	for k, v := range m.Between(20, 40) {
		fmt.Printf("%d: %s\n", k, v)
	}

	// Output:
	// 20: B
	// 30: C
}

func ExampleTreeMap_String() {
	m := treemap.NewOrdered[int, string]()
	fmt.Println(m.String())

	m.Put(10, "A")
	m.Put(20, "B")
	m.Put(30, "C")
	fmt.Println(m.String())

	// Output:
	// map[]
	// map[10:A 20:B 30:C]
}
