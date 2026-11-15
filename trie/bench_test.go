package trie_test

import (
	"fmt"
	"testing"

	"github.com/lock14/collections/trie"
)

const benchSize = 1000

var (
	sinkInt    int
	sinkBool   bool
	sinkString string
)

func buildStringPool(n int) []string {
	keys := make([]string, n)
	for i := range keys {
		keys[i] = fmt.Sprintf("api/v1/resource/%04d/detail", i)
	}
	return keys
}

func buildBytePool(keys []string) [][]byte {
	out := make([][]byte, len(keys))
	for i, k := range keys {
		out[i] = []byte(k)
	}
	return out
}

func BenchmarkStringMap_PutPool(b *testing.B) {
	b.ReportAllocs()
	keys := buildStringPool(benchSize)
	m := trie.NewMap[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(keys[i%benchSize], i)
	}
}

func BenchmarkStringMap_PutRemove(b *testing.B) {
	b.ReportAllocs()
	keys := buildStringPool(benchSize)
	m := trie.NewMap[int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := keys[i%benchSize]
		m.Remove(k)
		m.Put(k, i)
	}
}

func BenchmarkStringMap_Get(b *testing.B) {
	b.ReportAllocs()
	keys := buildStringPool(benchSize)
	m := trie.NewMap[int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt, sinkBool = m.Get(keys[i%benchSize])
	}
}

func BenchmarkStringMap_HasPrefix(b *testing.B) {
	b.ReportAllocs()
	keys := buildStringPool(benchSize)
	m := trie.NewMap[int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBool = m.HasPrefix("api/v1/resource/")
	}
}

func BenchmarkStringMap_LongestPrefixOf(b *testing.B) {
	b.ReportAllocs()
	keys := buildStringPool(benchSize)
	m := trie.NewMap[int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	query := "api/v1/resource/0500/detail/subpath"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString, sinkInt, sinkBool = m.LongestPrefixOf(query)
	}
}

func BenchmarkStringMap_Values(b *testing.B) {
	b.ReportAllocs()
	keys := buildStringPool(benchSize)
	m := trie.NewMap[int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for v := range m.Values() {
			sum += v
		}
		sinkInt = sum
	}
}

func BenchmarkStringMap_KeysWithPrefixPool(b *testing.B) {
	b.ReportAllocs()
	keys := buildStringPool(benchSize)
	m := trie.NewMap[int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		for range m.KeysWithPrefix("api/v1/resource/01") {
			count++
		}
		sinkInt = count
	}
}

func BenchmarkSliceMap_PutPool(b *testing.B) {
	b.ReportAllocs()
	keys := buildBytePool(buildStringPool(benchSize))
	m := trie.NewSliceMap[byte, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(keys[i%benchSize], i)
	}
}

func BenchmarkSliceMap_PutRemove(b *testing.B) {
	b.ReportAllocs()
	keys := buildBytePool(buildStringPool(benchSize))
	m := trie.NewSliceMap[byte, int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := keys[i%benchSize]
		m.Remove(k)
		m.Put(k, i)
	}
}

func BenchmarkSliceMap_Get(b *testing.B) {
	b.ReportAllocs()
	keys := buildBytePool(buildStringPool(benchSize))
	m := trie.NewSliceMap[byte, int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt, sinkBool = m.Get(keys[i%benchSize])
	}
}

func BenchmarkSliceMap_HasPrefix(b *testing.B) {
	b.ReportAllocs()
	keys := buildBytePool(buildStringPool(benchSize))
	m := trie.NewSliceMap[byte, int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	prefix := []byte("api/v1/resource/")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBool = m.HasPrefix(prefix)
	}
}

func BenchmarkSliceMap_Values(b *testing.B) {
	b.ReportAllocs()
	keys := buildBytePool(buildStringPool(benchSize))
	m := trie.NewSliceMap[byte, int]()
	for i, k := range keys {
		m.Put(k, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for v := range m.Values() {
			sum += v
		}
		sinkInt = sum
	}
}

func BenchmarkStringSet_Contains(b *testing.B) {
	b.ReportAllocs()
	keys := buildStringPool(benchSize)
	s := trie.NewSet()
	for _, k := range keys {
		s.Add(k)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBool = s.Contains(keys[i%benchSize])
	}
}
