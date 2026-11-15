package trie_test

import (
	"fmt"
	"slices"

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

func ExampleNewSet() {
	s := trie.NewSet()
	s.Add("go")
	s.Add("golang")
	s.Add("gopher")
	s.Add("rust")

	fmt.Println("Has prefix 'go':", s.HasPrefix("go"))
	for elem := range s.ElementsWithPrefix("go") {
		fmt.Println("-", elem)
	}

	// Output:
	// Has prefix 'go': true
	// - go
	// - golang
	// - gopher
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

func ExampleNewSliceSet() {
	s := trie.NewSliceSet[string]()
	s.Add([]string{"usr", "local"})
	s.Add([]string{"usr", "local", "bin"})

	longest, ok := s.LongestPrefixOf([]string{"usr", "local", "bin", "go"})
	fmt.Println(longest, ok)

	// Output:
	// [usr local bin] true
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

	m.Put("/api/", 1)
	m.Put("/api/v1/", 2)
	m.Put("/api/v1/users", 3)

	query := "/api/v1/users/123/profile"
	k, v, ok := m.LongestPrefixOf(query)
	fmt.Printf("Matched route: %s (value %d, found %v)\n", k, v, ok)

	// Output:
	// Matched route: /api/v1/users (value 3, found true)
}

func ExampleMap_ShortestPrefixOf() {
	m := trie.NewMap[string]()
	m.Put("/api", "root-api")
	m.Put("/api/v1", "v1-api")

	k, v, ok := m.ShortestPrefixOf("/api/v1/items")
	fmt.Println(k, v, ok)

	// Output:
	// /api root-api true
}

func ExampleMap_PrefixesOf() {
	m := trie.NewMap[int]()
	m.Put("a", 1)
	m.Put("ab", 2)
	m.Put("abc", 3)

	var keys []string
	for k := range m.PrefixesOf("abcd") {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	fmt.Println(keys)

	// Output:
	// [a ab abc]
}

func ExampleMap_RemovePrefix() {
	m := trie.NewMap[int]()
	m.Put("tmp/a", 1)
	m.Put("tmp/b", 2)
	m.Put("var/log", 3)

	m.RemovePrefix("tmp/")
	fmt.Println(m)

	// Output:
	// map[var/log:3]
}
