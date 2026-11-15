package treemap_test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/lock14/collections/treemap"
)

const benchSize = 1000

var (
	sinkInt  int
	sinkBool bool
)

func shuffledKeys(size int) []int {
	rng := rand.New(rand.NewPCG(42, uint64(size)))
	keys := make([]int, size)
	for i := range keys {
		keys[i] = i
	}
	rng.Shuffle(size, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	return keys
}

func BenchmarkTreeMap_Put(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			keys := shuffledKeys(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m := treemap.NewOrdered[int, int]()
				for _, k := range keys {
					m.Put(k, k)
				}
			}
		})
	}
}

func BenchmarkTreeMap_PutRemove(b *testing.B) {
	b.ReportAllocs()
	keys := shuffledKeys(benchSize)
	m := treemap.NewOrdered[int, int]()
	for _, k := range keys {
		m.Put(k, k)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := keys[i%benchSize]
		m.Remove(k)
		m.Put(k, i)
	}
}

func BenchmarkTreeMap_Get(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			keys := shuffledKeys(size)
			m := treemap.NewOrdered[int, int]()
			for _, k := range keys {
				m.Put(k, k)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sinkInt, sinkBool = m.Get(keys[i%size])
			}
		})
	}
}

func BenchmarkTreeMap_Floor(b *testing.B) {
	b.ReportAllocs()
	m := treemap.NewOrdered[int, int]()
	for i := 0; i < benchSize; i++ {
		m.Put(i*2, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v int
		sinkInt, v, sinkBool = m.Floor((i % benchSize) * 2)
		sinkInt += v
	}
}

func BenchmarkTreeMap_Higher(b *testing.B) {
	b.ReportAllocs()
	m := treemap.NewOrdered[int, int]()
	for i := 0; i < benchSize; i++ {
		m.Put(i*2, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v int
		sinkInt, v, sinkBool = m.Higher((i % benchSize) * 2)
		sinkInt += v
	}
}

func BenchmarkTreeMap_First(b *testing.B) {
	b.ReportAllocs()
	m := treemap.NewOrdered[int, int]()
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, v := m.First()
		sinkInt = k + v
	}
}

func BenchmarkTreeMap_All(b *testing.B) {
	b.ReportAllocs()
	m := treemap.NewOrdered[int, int]()
	for i := 0; i < benchSize; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for k, v := range m.All() {
			sum += k + v
		}
		sinkInt = sum
	}
}

func BenchmarkTreeMap_Between(b *testing.B) {
	b.ReportAllocs()
	const total = 10000
	m := treemap.NewOrdered[int, int]()
	for i := 0; i < total; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for k, v := range m.Between(5000, 5020) {
			sum += k + v
		}
		sinkInt = sum
	}
}
