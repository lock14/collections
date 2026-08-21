package result_test

import (
	"testing"

	"github.com/lock14/collections/result"
)

func BenchmarkResult_Ok(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = result.Ok[int, string](i)
	}
}

func BenchmarkResult_Err(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = result.Err[int, string]("err")
	}
}

func BenchmarkResult_IsOk(b *testing.B) {
	res := result.Ok[int, string](42)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = res.IsOk()
	}
}

func BenchmarkResult_IsErr(b *testing.B) {
	res := result.Err[int, string]("err")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = res.IsErr()
	}
}

func BenchmarkResult_Ok_Extract(b *testing.B) {
	res := result.Ok[int, string](42)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = res.Ok()
	}
}

func BenchmarkResult_Err_Extract(b *testing.B) {
	res := result.Err[int, string]("err")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = res.Err()
	}
}

func BenchmarkResult_Unwrap(b *testing.B) {
	res := result.Ok[int, string](42)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = res.Unwrap()
	}
}

func BenchmarkResult_OrElse(b *testing.B) {
	res := result.Err[int, string]("err")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = res.OrElse(100)
	}
}

func BenchmarkResult_OrElseGet(b *testing.B) {
	res := result.Ok[int, string](42)
	supplier := func() int { return 100 }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = res.OrElseGet(supplier)
	}
}

func BenchmarkResult_Map(b *testing.B) {
	res := result.Ok[int, string](42)
	mapper := func(value int) int { return value * 2 }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = res.Map(mapper)
	}
}

func BenchmarkResult_MapErr(b *testing.B) {
	res := result.Err[int, int](404)
	mapper := func(err int) int { return err + 1 }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = res.MapErr(mapper)
	}
}

func BenchmarkResult_FlatMap(b *testing.B) {
	res := result.Ok[int, string](42)
	mapper := func(value int) result.Result[int, string] {
		return result.Ok[int, string](value * 2)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = res.FlatMap(mapper)
	}
}
