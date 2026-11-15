package graph_test

import (
	"testing"

	"github.com/lock14/collections/graph"
)

const benchSize = 1000

var (
	sinkInt  int
	sinkBool bool
)

func buildBenchDirectedGraph(n int) *graph.Graph[int] {
	g := graph.New[int](graph.WithDirected(), graph.WithCapacity(n))
	for i := range n {
		g.AddEdge(i, (i+1)%n)
		g.AddEdge(i, (i+7)%n)
		g.AddEdge(i, (i+31)%n)
		g.AddEdge(i, (i+127)%n)
	}
	return g
}

func BenchmarkGraph_AddRemoveVertex(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := benchSize + (i & 1023)
		g.AddVertex(v)
		g.RemoveVertex(v)
	}
}

func BenchmarkGraph_AddRemoveEdge(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := i % benchSize
		v := (u + 503) % benchSize
		g.AddEdge(u, v)
		g.RemoveEdge(u, v)
	}
}

func BenchmarkGraph_AddVertex_Preallocated(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int](graph.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
	}
	sinkInt = g.Order()
}

func BenchmarkGraph_ContainsEdge(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	var ok bool
	for i := 0; i < b.N; i++ {
		u := i % benchSize
		ok = g.ContainsEdge(u, (u+1)%benchSize)
	}
	sinkBool = ok
}

func BenchmarkGraph_Successors(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	var sum int
	for i := 0; i < b.N; i++ {
		for v := range g.Successors(i % benchSize) {
			sum += v
		}
	}
	sinkInt = sum
}

func BenchmarkGraph_Vertices(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	var sum int
	for i := 0; i < b.N; i++ {
		for v := range g.Vertices() {
			sum += v
		}
	}
	sinkInt = sum
}

func BenchmarkGraph_Edges(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	var sum int
	for i := 0; i < b.N; i++ {
		for u, v := range g.Edges() {
			sum += u + v
		}
	}
	sinkInt = sum
}

func BenchmarkGraph_IncidentEdges_Directed(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	var sum int
	for i := 0; i < b.N; i++ {
		for u, v := range g.IncidentEdges(i % benchSize) {
			sum += u + v
		}
	}
	sinkInt = sum
}
