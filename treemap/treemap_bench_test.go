package treemap_test

import (
	"fmt"
	"github.com/lock14/collections/treemap"
	"math/rand"
	"testing"
)

func BenchmarkTreeMap_Put(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
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
				tm := treemap.NewOrdered[int, int]()
				for _, k := range keys {
					tm.Put(k, k)
				}
			}
		})
	}
}

func BenchmarkTreeMap_Get(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.StopTimer()
			tm := treemap.NewOrdered[int, int]()
			keys := make([]int, size)
			for i := 0; i < size; i++ {
				keys[i] = i
				tm.Put(i, i)
			}
			rand.Shuffle(size, func(i, j int) {
				keys[i], keys[j] = keys[j], keys[i]
			})

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				for _, k := range keys {
					_, _ = tm.Get(k)
				}
			}
		})
	}
}

func BenchmarkBuiltinMap_Get(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.StopTimer()
			m := make(map[int]int, size)
			keys := make([]int, size)
			for i := 0; i < size; i++ {
				keys[i] = i
				m[i] = i
			}
			rand.Shuffle(size, func(i, j int) {
				keys[i], keys[j] = keys[j], keys[i]
			})

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				for _, k := range keys {
					_, _ = m[k]
				}
			}
		})
	}
}
