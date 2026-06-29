package trie_test

import (
	"fmt"

	"github.com/lock14/collections/trie"
)

func ExampleNewMap() {
	m := trie.NewMap[int]()

	m.Put("apple", 1)
	m.Put("app", 2)
	m.Put("ape", 3)
	m.Put("bat", 4)

	v, ok := m.Get("app")
	fmt.Printf("Get('app'): %d, %v\n", v, ok)

	fmt.Printf("Size: %d\n", m.Size())

	// Output:
	// Get('app'): 2, true
	// Size: 4
}

func ExampleMap_KeysWithPrefix() {
	m := trie.NewMap[int]()

	m.Put("apple", 1)
	m.Put("app", 2)
	m.Put("ape", 3)
	m.Put("bat", 4)

	fmt.Println("Keys with prefix 'ap':")
	for k := range m.KeysWithPrefix("ap") {
		fmt.Println(k)
	}

	// Output:
	// Keys with prefix 'ap':
	// ape
	// app
	// apple
}

func ExampleMap_LongestPrefixOf() {
	m := trie.NewMap[int]()

	// E.g., setting up routing rules
	m.Put("/api/", 1)
	m.Put("/api/v1/", 2)
	m.Put("/api/v1/users", 3)

	query := "/api/v1/users/123/profile"
	k, v, ok := m.LongestPrefixOf(query)
	fmt.Printf("Matched route: %s (value %d, found %v)\n", k, v, ok)

	// Output:
	// Matched route: /api/v1/users (value 3, found true)
}

func ExampleNewSliceMap() {
	m := trie.NewSliceMap[byte, string]()

	m.Put([]byte{192, 168, 0, 1}, "router")
	m.Put([]byte{192, 168, 0, 100}, "laptop")

	k, v, ok := m.LongestPrefixOf([]byte{192, 168, 0, 100, 255})
	fmt.Printf("Matched: %v -> %s (found: %v)\n", k, v, ok)

	// Output:
	// Matched: [192 168 0 100] -> laptop (found: true)
}
