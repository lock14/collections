package graph_test

import (
	"testing"

	"github.com/lock14/collections/graph"
)

func BenchmarkGraph_AddVertex(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
	}
}

func BenchmarkGraph_AddVertex_Preallocated(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int](graph.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
	}
}

func BenchmarkGraph_AddEdge_Directed(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int](graph.WithDirected())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1)
	}
}

func BenchmarkGraph_AddEdge_Directed_Preallocated(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int](graph.WithDirected(), graph.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1)
	}
}

func BenchmarkGraph_AddEdge_Undirected(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1)
	}
}

func BenchmarkGraph_AddEdge_Undirected_Preallocated(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int](graph.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1)
	}
}

func BenchmarkGraph_RemoveVertex(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int]()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.RemoveVertex(i)
	}
}

func BenchmarkGraph_RemoveEdge(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int]()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.RemoveEdge(i, i+1)
	}
}

func BenchmarkGraph_IterateVertices(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int]()
	for i := 0; i < 1000; i++ {
		g.AddVertex(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range g.Vertices() {
		}
	}
}

func BenchmarkGraph_IterateEdges(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int]()
	for i := 0; i < 1000; i++ {
		g.AddEdge(i, i+1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range g.Edges() {
		}
	}
}

func BenchmarkGraph_IncidentEdges(b *testing.B) {
	b.ReportAllocs()
	g := graph.New[int]()
	for i := 0; i < 1000; i++ {
		g.AddEdge(0, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range g.IncidentEdges(0) {
		}
	}
}
