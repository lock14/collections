// Package trie provides Trie backed map and set implementations.
package trie

import (
	"iter"

	"github.com/lock14/collections"
)

// Map is a Trie that implements collections.MutableMap and provides prefix operations.
type Map[K any, V any] interface {
	collections.MutableMap[K, V]

	// HasPrefix returns true if there is at least one key in the map starting with the given prefix.
	HasPrefix(prefix K) bool

	// KeysWithPrefix returns an iterator over all keys that start with the given prefix.
	KeysWithPrefix(prefix K) iter.Seq[K]

	// ValuesWithPrefix returns an iterator over all values whose keys start with the given prefix.
	ValuesWithPrefix(prefix K) iter.Seq[V]

	// EntriesWithPrefix returns an iterator over all key-value pairs where the key starts with the given prefix.
	EntriesWithPrefix(prefix K) iter.Seq2[K, V]

	// RemovePrefix removes all key-value pairs whose keys start with the given prefix.
	RemovePrefix(prefix K)

	// LongestPrefixOf returns the longest key in the map that is a prefix of the given query.
	LongestPrefixOf(query K) (K, V, bool)

	// ShortestPrefixOf returns the shortest key in the map that is a prefix of the given query.
	ShortestPrefixOf(query K) (K, V, bool)

	// PrefixesOf returns an iterator over all keys in the map that are a prefix of the given query.
	PrefixesOf(query K) iter.Seq2[K, V]
}

// Set is a Trie that implements collections.MutableSet and provides prefix operations.
type Set[K any] interface {
	collections.MutableSet[K]

	// HasPrefix returns true if there is at least one element in the set starting with the given prefix.
	HasPrefix(prefix K) bool

	// ElementsWithPrefix returns an iterator over all elements that start with the given prefix.
	ElementsWithPrefix(prefix K) iter.Seq[K]

	// RemovePrefix removes all elements that start with the given prefix.
	RemovePrefix(prefix K)

	// LongestPrefixOf returns the longest element in the set that is a prefix of the given query.
	LongestPrefixOf(query K) (K, bool)

	// ShortestPrefixOf returns the shortest element in the set that is a prefix of the given query.
	ShortestPrefixOf(query K) (K, bool)

	// PrefixesOf returns an iterator over all elements in the set that are a prefix of the given query.
	PrefixesOf(query K) iter.Seq[K]
}

// NewMap creates an empty Trie map for string keys.
func NewMap[V any]() Map[string, V] {
	return newStringMap[V]()
}

// NewSet creates an empty Trie set for string elements.
func NewSet() Set[string] {
	return newStringSet()
}

// NewSliceMap creates an empty Trie map for slice keys of any comparable element type.
func NewSliceMap[E comparable, V any]() Map[[]E, V] {
	return newSliceMap[E, V]()
}

// NewSliceSet creates an empty Trie set for slice elements of any comparable element type.
func NewSliceSet[E comparable]() Set[[]E] {
	return newSliceSet[E]()
}
