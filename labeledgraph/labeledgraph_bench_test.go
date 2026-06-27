package labeledgraph_test

import (
	"github.com/lock14/collections/labeledgraph"
	"testing"
)

func BenchmarkLabeledGraph_AddVertex(b *testing.B) {
	g := labeledgraph.New[int, string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
	}
}

func BenchmarkLabeledGraph_AddEdge_Directed(b *testing.B) {
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1, "label")
	}
}

func BenchmarkLabeledGraph_AddEdge_Undirected(b *testing.B) {
	g := labeledgraph.New[int, string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1, "label")
	}
}

func BenchmarkLabeledGraph_RemoveVertex(b *testing.B) {
	g := labeledgraph.New[int, string]()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.RemoveVertex(i)
	}
}

func BenchmarkLabeledGraph_RemoveEdge(b *testing.B) {
	g := labeledgraph.New[int, string]()
	for i := 0; i < b.N; i++ {
		g.AddEdge(i, i+1, "l")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.RemoveEdge(i, i+1)
	}
}

func BenchmarkLabeledGraph_IterateVertices(b *testing.B) {
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
