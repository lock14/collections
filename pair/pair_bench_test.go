package pair_test

import (
	"github.com/lock14/collections/pair"
	"testing"
)

func BenchmarkPair_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pair.New(i, i)
	}
}

func BenchmarkPair_Unwrap(b *testing.B) {
	p := pair.New(10, 20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Unwrap()
	}
}
