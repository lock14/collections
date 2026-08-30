package linkedhashmap_test

import (
	"fmt"
	"github.com/lock14/collections/linkedhashmap"
)

func ExampleLinkedHashMap() {
	// LinkedHashMap is a hash map that preserves insertion order.
	m := linkedhashmap.New[string, string]()

	m.Put("First", "A")
	m.Put("Second", "B")
	m.Put("Third", "C")

	// Updating an existing key does not change its position in the insertion order
	m.Put("First", "A+")

	fmt.Println("Iterating over LinkedHashMap:")
	for k, v := range m.All() {
		fmt.Printf("- %s: %s\n", k, v)
	}

	// Output:
	// Iterating over LinkedHashMap:
	// - First: A+
	// - Second: B
	// - Third: C
}

func ExampleLinkedHashMap_lruCache() {
	// A LinkedHashMap can be used as an LRU cache by configuring it with
	// WithAccessOrder and WithMaxElements.
	cache := linkedhashmap.New[int, string](
		linkedhashmap.WithAccessOrder(),
		linkedhashmap.WithMaxElements(3),
	)

	cache.Put(1, "one")
	cache.Put(2, "two")
	cache.Put(3, "three")

	// accessing 1 moves it to the end (most recently used)
	cache.Get(1)

	// inserting 4 will evict the eldest element (which is 2 now)
	cache.Put(4, "four")

	for k, v := range cache.All() {
		fmt.Printf("%d: %s\n", k, v)
	}

	// Output:
	// 3: three
	// 1: one
	// 4: four
}

func ExampleLinkedHashMap_Put() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	fmt.Println(m.Size())
	// Output:
	// 2
}

func ExampleLinkedHashMap_PutFirst() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.PutFirst(2, "B")
	for k, v := range m.All() {
		fmt.Printf("%d: %s\n", k, v)
	}
	// Output:
	// 2: B
	// 1: A
}

func ExampleLinkedHashMap_PutLast() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.PutLast(2, "B")
	for k, v := range m.All() {
		fmt.Printf("%d: %s\n", k, v)
	}
	// Output:
	// 1: A
	// 2: B
}

func ExampleLinkedHashMap_First() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	k, v := m.First()
	fmt.Printf("%d: %s\n", k, v)
	// Output:
	// 1: A
}

func ExampleLinkedHashMap_Last() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	k, v := m.Last()
	fmt.Printf("%d: %s\n", k, v)
	// Output:
	// 2: B
}

func ExampleLinkedHashMap_PollFirst() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	k, v := m.PollFirst()
	fmt.Printf("Polled: %d: %s, Size: %d\n", k, v, m.Size())
	// Output:
	// Polled: 1: A, Size: 1
}

func ExampleLinkedHashMap_PollLast() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	k, v := m.PollLast()
	fmt.Printf("Polled: %d: %s, Size: %d\n", k, v, m.Size())
	// Output:
	// Polled: 2: B, Size: 1
}

func ExampleLinkedHashMap_Get() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	v, ok := m.Get(1)
	fmt.Printf("value: %s, ok: %t\n", v, ok)
	_, ok = m.Get(2)
	fmt.Printf("ok: %t\n", ok)
	// Output:
	// value: A, ok: true
	// ok: false
}

func ExampleLinkedHashMap_Remove() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Remove(1)
	fmt.Println(m.ContainsKey(1))
	// Output:
	// false
}

func ExampleLinkedHashMap_ContainsKey() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	fmt.Println(m.ContainsKey(1))
	fmt.Println(m.ContainsKey(2))
	// Output:
	// true
	// false
}

func ExampleLinkedHashMap_Size() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	fmt.Println(m.Size())
	// Output:
	// 2
}

func ExampleLinkedHashMap_Empty() {
	m := linkedhashmap.New[int, string]()
	fmt.Println(m.Empty())
	m.Put(1, "A")
	fmt.Println(m.Empty())
	// Output:
	// true
	// false
}

func ExampleLinkedHashMap_Clear() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Clear()
	fmt.Println(m.Size())
	// Output:
	// 0
}

func ExampleLinkedHashMap_All() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	for k, v := range m.All() {
		fmt.Printf("%d: %s\n", k, v)
	}
	// Output:
	// 1: A
	// 2: B
}

func ExampleLinkedHashMap_Keys() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	for k := range m.Keys() {
		fmt.Println(k)
	}
	// Output:
	// 1
	// 2
}

func ExampleLinkedHashMap_Values() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	for v := range m.Values() {
		fmt.Println(v)
	}
	// Output:
	// A
	// B
}

func ExampleLinkedHashMap_Backward() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	for k, v := range m.Backward() {
		fmt.Printf("%d: %s\n", k, v)
	}
	// Output:
	// 2: B
	// 1: A
}

func ExampleLinkedHashMap_BackwardKeys() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	for k := range m.BackwardKeys() {
		fmt.Println(k)
	}
	// Output:
	// 2
	// 1
}

func ExampleLinkedHashMap_BackwardValues() {
	m := linkedhashmap.New[int, string]()
	m.Put(1, "A")
	m.Put(2, "B")
	for v := range m.BackwardValues() {
		fmt.Println(v)
	}
	// Output:
	// B
	// A
}

func ExampleLinkedHashMap_String() {
	m := linkedhashmap.New[int, string]()
	fmt.Println(m.String())

	m.Put(1, "A")
	m.Put(2, "B")
	fmt.Println(m.String())

	// Output:
	// map[]
	// map[1:A 2:B]
}
