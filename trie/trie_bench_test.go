package trie_test

import (
	"testing"

	"github.com/lock14/collections/trie"
)

func BenchmarkStringMap_Put(b *testing.B) {
	b.ReportAllocs()
	m := trie.NewMap[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put("some_fairly_long_string_prefix", i)
	}
}

func BenchmarkSliceMap_Put(b *testing.B) {
	b.ReportAllocs()
	m := trie.NewSliceMap[byte, int]()
	key := []byte("some_fairly_long_string_prefix")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(key, i)
	}
}

func BenchmarkStringMap_Get(b *testing.B) {
	b.ReportAllocs()
	m := trie.NewMap[int]()
	m.Put("some_fairly_long_string_prefix", 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get("some_fairly_long_string_prefix")
	}
}

func BenchmarkSliceMap_Get(b *testing.B) {
	b.ReportAllocs()
	m := trie.NewSliceMap[byte, int]()
	key := []byte("some_fairly_long_string_prefix")
	m.Put(key, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(key)
	}
}

func BenchmarkStringMap_KeysWithPrefix(b *testing.B) {
	b.ReportAllocs()
	m := trie.NewMap[int]()
	m.Put("api/v1/users/1", 1)
	m.Put("api/v1/users/2", 2)
	m.Put("api/v1/posts/1", 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range m.KeysWithPrefix("api/v1/users") {
		}
	}
}
