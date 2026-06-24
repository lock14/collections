package treeset_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/lock14/collections/treeset"
)

func BenchmarkTreeSet_Add(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.StopTimer()
			keys := make([]int, size)
			for i := 0; i < size; i++ {
				keys[i] = i
			}
			rand.Shuffle(size, func(i, j int) {
				keys[i], keys[j] = keys[j], keys[i]
			})

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				s := treeset.NewOrdered[int]()
				for _, k := range keys {
					s.Add(k)
				}
			}
		})
	}
}

func BenchmarkTreeSet_Contains(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.StopTimer()
			s := treeset.NewOrdered[int]()
			keys := make([]int, size)
			for i := 0; i < size; i++ {
				keys[i] = i
				s.Add(i)
			}
			rand.Shuffle(size, func(i, j int) {
				keys[i], keys[j] = keys[j], keys[i]
			})

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				for _, k := range keys {
					_ = s.Contains(k)
				}
			}
		})
	}
}
