package comparator_test

import (
	"github.com/lock14/collections/comparator"
	"testing"
)

func BenchmarkNaturalOrder_Compare(b *testing.B) {
	b.ReportAllocs()
	cmp := comparator.NaturalOrder[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cmp(i, i+1)
	}
}

func BenchmarkReverse_Compare(b *testing.B) {
	b.ReportAllocs()
	cmp := comparator.Reverse(comparator.NaturalOrder[int]())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cmp(i, i+1)
	}
}
