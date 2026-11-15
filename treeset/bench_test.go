package treeset_test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/lock14/collections/treeset"
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

func BenchmarkTreeSet_Add(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			keys := shuffledKeys(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s := treeset.NewOrdered[int]()
				for _, k := range keys {
					s.Add(k)
				}
			}
		})
	}
}

func BenchmarkTreeSet_AddRemove(b *testing.B) {
	b.ReportAllocs()
	keys := shuffledKeys(benchSize)
	s := treeset.NewOrdered[int]()
	for _, k := range keys {
		s.Add(k)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := keys[i%benchSize]
		s.RemoveElement(k)
		s.Add(k)
	}
}

func BenchmarkTreeSet_Contains(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			keys := shuffledKeys(size)
			s := treeset.NewOrdered[int]()
			for _, k := range keys {
				s.Add(k)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sinkBool = s.Contains(keys[i%size])
			}
		})
	}
}

func BenchmarkTreeSet_Floor(b *testing.B) {
	b.ReportAllocs()
	s := treeset.NewOrdered[int]()
	for i := 0; i < benchSize; i++ {
		s.Add(i * 2)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt, sinkBool = s.Floor((i % benchSize) * 2)
	}
}

func BenchmarkTreeSet_First(b *testing.B) {
	b.ReportAllocs()
	s := treeset.NewOrdered[int]()
	for i := 0; i < benchSize; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt = s.First()
	}
}

func BenchmarkTreeSet_All(b *testing.B) {
	b.ReportAllocs()
	s := treeset.NewOrdered[int]()
	for i := 0; i < benchSize; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for v := range s.All() {
			sum += v
		}
		sinkInt = sum
	}
}

func BenchmarkTreeSet_Between(b *testing.B) {
	b.ReportAllocs()
	const total = 10000
	s := treeset.NewOrdered[int]()
	for i := 0; i < total; i++ {
		s.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for v := range s.Between(5000, 5020) {
			sum += v
		}
		sinkInt = sum
	}
}
