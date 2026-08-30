package labeledgraph_test

import (
	"github.com/lock14/collections/labeledgraph"
	"testing"
)

func BenchmarkLabeledGraph_AddVertex(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
	}
}

func BenchmarkLabeledGraph_AddVertex_Preallocated(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string](labeledgraph.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
	}
}

func BenchmarkLabeledGraph_AddEdge_Directed(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1, "label")
	}
}

func BenchmarkLabeledGraph_AddEdge_Directed_Preallocated(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string](labeledgraph.WithDirected(), labeledgraph.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1, "label")
	}
}

func BenchmarkLabeledGraph_AddEdge_Undirected(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1, "label")
	}
}

func BenchmarkLabeledGraph_AddEdge_Undirected_Preallocated(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string](labeledgraph.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1, "label")
	}
}

func BenchmarkLabeledGraph_AddRemoveVertex(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
		g.RemoveVertex(i)
	}
}

func BenchmarkLabeledGraph_AddRemoveEdge(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string]()
	g.AddVertex(1)
	g.AddVertex(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(1, 2, "l")
		g.RemoveEdge(1, 2)
	}
}

func BenchmarkLabeledGraph_IterateVertices(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string]()
	for i := 0; i < 1000; i++ {
		g.AddVertex(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range g.Vertices() {
		}
	}
}

func BenchmarkLabeledGraph_IterateEdges(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string]()
	for i := 0; i < 1000; i++ {
		g.AddEdge(i, i+1, "l")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range g.Edges() {
		}
	}
}

func BenchmarkLabeledGraph_IncidentEdges(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string]()
	for i := 0; i < 1000; i++ {
		g.AddEdge(0, i, "l")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range g.IncidentEdges(0) {
		}
	}
}
