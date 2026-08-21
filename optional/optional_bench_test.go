package optional_test

import (
	"testing"

	"github.com/lock14/collections/optional"
)

func BenchmarkOption_Creation(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = optional.Of(i)
	}
}

func BenchmarkOption_Get(b *testing.B) {
	opt := optional.Of(42)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = opt.Get()
	}
}

func BenchmarkOption_IsPresent(b *testing.B) {
	opt := optional.Of(42)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.IsPresent()
	}
}

func BenchmarkOption_OrElse(b *testing.B) {
	opt := optional.Empty[int]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.OrElse(100)
	}
}

func BenchmarkOption_OrElseGet(b *testing.B) {
	opt := optional.Of(42)
	supplier := func() int { return 100 }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.OrElseGet(supplier)
	}
}

func BenchmarkOption_Filter(b *testing.B) {
	opt := optional.Of(42)
	predicate := func(value int) bool { return value > 0 }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.Filter(predicate)
	}
}

func BenchmarkOption_Map(b *testing.B) {
	opt := optional.Of(42)
	mapper := func(value int) int { return value * 2 }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.Map(mapper)
	}
}

func BenchmarkOption_FlatMap(b *testing.B) {
	opt := optional.Of(42)
	mapper := func(value int) optional.Option[int] {
		return optional.Of(value * 2)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.FlatMap(mapper)
	}
}
