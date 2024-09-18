package graph

import (
	"iter"
	"maps"
)

type LabeledGraph[V comparable, L any] struct {
	graph     map[V]nodeData[V, L]
	directed  bool
	edgeCount int
}

func (g *LabeledGraph[V, L]) AddVertex(v V) {
	if !g.ContainsVertex(v) {
		g.graph[v] = nil // TODO: add nodeData impl
	}
}

func (g *LabeledGraph[V, L]) AddEdge(u, v V, l L) {
	g.AddVertex(u)
	g.AddVertex(v)
	g.graph[u].AddSuccessor(v, l)
	g.graph[v].AddPredecessor(u, l)
}

func (g *LabeledGraph[V, L]) ContainsEdge(u, v V) bool {
	return g.ContainsVertex(u) &&
		g.ContainsVertex(v) &&
		g.graph[u].ContainsSuccessor(v)
}

func (g *LabeledGraph[V, L]) ContainsVertex(v V) bool {
	_, ok := g.graph[v]
	return ok
}

func (g *LabeledGraph[V, L]) Directed() bool {
	return g.directed
}

func (g *LabeledGraph[V, L]) EdgeCount() int {
	return g.edgeCount
}

func (g *LabeledGraph[V, L]) Edges() iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for u := range g.Vertices() {
			for v := range g.graph[u].Successors() {
				if !yield(u, v) {
					return
				}
			}
		}
	}
}

func (g *LabeledGraph[V, L]) Label(u, v V) L {
	return g.graph[u].Label(v)
}

func (g *LabeledGraph[V, L]) RemoveVertex(v V) {

}

func (g *LabeledGraph[V, L]) Vertices() iter.Seq[V] {
	return maps.Keys(g.graph)
}

type nodeData[V, L any] interface {
	AddPredecessor(V, L)
	AddSuccessor(V, L)
	ContainsPredecessor(V) bool
	ContainsSuccessor(V) bool
	Label(V) L
	Predecessors() iter.Seq[V]
	RemovePredecessor(V)
	RemoveSuccessor(V)
	Successors() iter.Seq[V]
}
