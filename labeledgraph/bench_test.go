package labeledgraph_test

import (
	"testing"

	"github.com/lock14/collections/labeledgraph"
)

const benchSize = 1000

var (
	sinkInt  int
	sinkBool bool
	sinkStr  string
)

func buildBenchDirectedGraph(n int) *labeledgraph.LabeledGraph[int, string] {
	g := labeledgraph.New[int, string](labeledgraph.WithDirected(), labeledgraph.WithCapacity(n))
	for i := range n {
		g.AddEdge(i, (i+1)%n, "w1")
		g.AddEdge(i, (i+7)%n, "w2")
		g.AddEdge(i, (i+31)%n, "w3")
		g.AddEdge(i, (i+127)%n, "w4")
	}
	return g
}

func BenchmarkLabeledGraph_AddRemoveVertex(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := benchSize + (i & 1023)
		g.AddVertex(v)
		g.RemoveVertex(v)
	}
}

func BenchmarkLabeledGraph_AddRemoveEdge(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := i % benchSize
		v := (u + 503) % benchSize
		g.AddEdge(u, v, "dyn")
		g.RemoveEdge(u, v)
	}
}

func BenchmarkLabeledGraph_AddVertex_Preallocated(b *testing.B) {
	b.ReportAllocs()
	g := labeledgraph.New[int, string](labeledgraph.WithCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddVertex(i)
	}
	sinkInt = g.Order()
}

func BenchmarkLabeledGraph_ContainsEdge(b *testing.B) {
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

func BenchmarkLabeledGraph_Label(b *testing.B) {
	g := buildBenchDirectedGraph(benchSize)
	b.ReportAllocs()
	b.ResetTimer()
	var (
		lbl string
		ok  bool
	)
	for i := 0; i < b.N; i++ {
		u := i % benchSize
		lbl, ok = g.Label(u, (u+1)%benchSize)
	}
	sinkStr = lbl
	sinkBool = ok
}

func BenchmarkLabeledGraph_Successors(b *testing.B) {
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

func BenchmarkLabeledGraph_Vertices(b *testing.B) {
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

func BenchmarkLabeledGraph_Edges(b *testing.B) {
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

func BenchmarkLabeledGraph_IncidentEdges_Directed(b *testing.B) {
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

func BenchmarkLabeledGraph_IncidentEdges_Undirected(b *testing.B) {
	g := labeledgraph.New[int, string](labeledgraph.WithCapacity(benchSize))
	for i := range benchSize {
		g.AddEdge(i, (i+1)%benchSize, "w1")
		g.AddEdge(i, (i+7)%benchSize, "w2")
	}
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
