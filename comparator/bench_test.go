package comparator_test

import (
	"testing"

	"github.com/lock14/collections/comparator"
)

var sinkInt int

func BenchmarkComparator_NaturalOrder(b *testing.B) {
	b.ReportAllocs()
	cmpFunc := comparator.NaturalOrder[int]()
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = cmpFunc(i, i+1)
	}
	sinkInt = r
}

func BenchmarkComparator_Reverse(b *testing.B) {
	b.ReportAllocs()
	cmpFunc := comparator.Reverse(comparator.NaturalOrder[int]())
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = cmpFunc(i, i+1)
	}
	sinkInt = r
}
