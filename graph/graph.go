// Package graph provides generic graph data structures and algorithms.
package graph

import (
	"fmt"
	"github.com/lock14/collections/labeledgraph"
	"iter"
	"strings"
)

// config holds configuration values for New to use when contracting a Graph.
type config struct {
	delegateOps []labeledgraph.Opt
}

// Opt represents a configuration option for constructing a Graph.
type Opt func(g *config)

// WithDirected configures New to return a directed graph.
func WithDirected() Opt {
	return func(g *config) {
		g.delegateOps = append(g.delegateOps, labeledgraph.WithDirected())
	}
}

// WithCapacity configures New to pre-allocate the graph with the given capacity.
func WithCapacity(n int) Opt {
	return func(g *config) {
		g.delegateOps = append(g.delegateOps, labeledgraph.WithCapacity(n))
	}
}

type void struct{}

// Graph is a graph with vertices of type V.
type Graph[V comparable] struct {
	delegate *labeledgraph.LabeledGraph[V, void]
}

var _ fmt.Stringer = (*Graph[int])(nil)

// New returns a new Graph constructed according to the given options.
func New[V comparable](opts ...Opt) *Graph[V] {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return &Graph[V]{
		delegate: labeledgraph.New[V, void](config.delegateOps...),
	}
}

// ContainsVertex returns whether the given vertex is contained in the graph.
func (g *Graph[V]) ContainsVertex(v V) bool {
	return g.delegate.ContainsVertex(v)
}

// AddVertex adds the given vertex to the graph.
func (g *Graph[V]) AddVertex(v V) {
	g.delegate.AddVertex(v)
}

// RemoveVertex removes the vertex and all incident edges from the graph.
func (g *Graph[V]) RemoveVertex(v V) {
	g.delegate.RemoveVertex(v)
}

// ContainsEdge returns whether the given edge is contained in the graph.
func (g *Graph[V]) ContainsEdge(u, v V) bool {
	return g.delegate.ContainsEdge(u, v)
}

// AddEdge adds the given edge to the graph.
func (g *Graph[V]) AddEdge(u, v V) {
	g.delegate.AddEdge(u, v, void{})
}

// RemoveEdge removes the given edge from the graph.
func (g *Graph[V]) RemoveEdge(u, v V) {
	g.delegate.RemoveEdge(u, v)
}

// Directed returns whether this graph is directed.
func (g *Graph[V]) Directed() bool {
	return g.delegate.Directed()
}

// Order returns the number of vertices in the graph.
func (g *Graph[V]) Order() int {
	return g.delegate.Order()
}

// Size returns the number of edges in the graph.
func (g *Graph[V]) Size() int {
	return g.delegate.Size()
}

// Degree return the number of edges coming into or out of the given vertex.
func (g *Graph[V]) Degree(u V) (int, bool) {
	return g.delegate.Degree(u)
}

// InDegree return the number of edges coming into the given vertex.
func (g *Graph[V]) InDegree(u V) (int, bool) {
	return g.delegate.InDegree(u)
}

// OutDegree return the number of edges coming out of the given vertex.
func (g *Graph[V]) OutDegree(u V) (int, bool) {
	return g.delegate.OutDegree(u)
}

// Vertices returns an iterator over all vertices in the graph.
func (g *Graph[V]) Vertices() iter.Seq[V] {
	return g.delegate.Vertices()
}

// Neighbors is an alias for Successors.
func (g *Graph[V]) Neighbors(v V) iter.Seq[V] {
	return g.delegate.Neighbors(v)
}

// Successors returns an iterator over the successors of the given vertex in the graph.
func (g *Graph[V]) Successors(u V) iter.Seq[V] {
	return g.delegate.Successors(u)
}

// Predecessors returns an iterator over the predecessors of the given vertex in the graph.
func (g *Graph[V]) Predecessors(v V) iter.Seq[V] {
	return g.delegate.Predecessors(v)
}

// Edges returns an iterator over all edges in the graph.
func (g *Graph[V]) Edges() iter.Seq2[V, V] {
	return g.delegate.Edges()
}

// IncidentEdges returns all edges that are incident to the given vertex in the graph.
func (g *Graph[V]) IncidentEdges(v V) iter.Seq2[V, V] {
	return g.delegate.IncidentEdges(v)
}

// InIncidentEdges returns an iterator over the edges coming into the given vertex of the graph.
func (g *Graph[V]) InIncidentEdges(v V) iter.Seq2[V, V] {
	return g.delegate.InIncidentEdges(v)
}

// OutIncidentEdges returns an iterator over the edges coming out of the given vertex of the graph.
func (g *Graph[V]) OutIncidentEdges(u V) iter.Seq2[V, V] {
	return g.delegate.OutIncidentEdges(u)
}

// Clear removes all vertices and edges from the graph.
func (g *Graph[V]) Clear() {
	g.delegate.Clear()
}

// Clone returns a deep copy of the graph.
func (g *Graph[V]) Clone() *Graph[V] {
	return &Graph[V]{
		delegate: g.delegate.Clone(),
	}
}

// Equal returns true if the two graphs are structurally identical.
func (g *Graph[V]) Equal(other *Graph[V]) bool {
	return g.delegate.Equal(other.delegate, func(void, void) bool { return true })
}

// String returns a string representation of the graph.
func (g *Graph[V]) String() string {
	var sb strings.Builder
	if g.Directed() {
		sb.WriteString("digraph[")
	} else {
		sb.WriteString("graph[")
	}
	first := true

	type edgePair struct{ u, v V }
	seen := make(map[edgePair]bool)

	for u, v := range g.Edges() {
		if !g.Directed() {
			if seen[edgePair{u, v}] || seen[edgePair{v, u}] {
				continue
			}
			seen[edgePair{u, v}] = true
		}

		if !first {
			sb.WriteString(", ")
		}
		if g.Directed() {
			sb.WriteString(fmt.Sprintf("%v -> %v", u, v))
		} else {
			sb.WriteString(fmt.Sprintf("%v - %v", u, v))
		}
		first = false
	}

	for v := range g.Vertices() {
		deg, _ := g.Degree(v)
		if deg == 0 {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%v", v))
			first = false
		}
	}

	sb.WriteString("]")
	return sb.String()
}

func defaultConfig() *config {
	return &config{}
}
