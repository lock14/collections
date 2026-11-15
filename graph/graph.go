// Package graph provides generic directed and undirected graph data structures.
package graph

import (
	"fmt"
	"iter"
	"strings"

	"github.com/lock14/collections/labeledgraph"
)

var _ fmt.Stringer = (*Graph[int])(nil)

// config holds configuration values for New to use when constructing a Graph.
type config struct {
	delegateOps []labeledgraph.Option
}

// Option represents a configuration option for constructing a Graph.
type Option func(c *config)

// WithDirected configures New to return a directed graph.
func WithDirected() Option {
	return func(c *config) {
		c.delegateOps = append(c.delegateOps, labeledgraph.WithDirected())
	}
}

// WithCapacity configures New to pre-allocate the graph with the given vertex capacity.
// Panics in New if capacity is negative.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		if capacity < 0 {
			panic("graph: capacity cannot be negative")
		}
		c.delegateOps = append(c.delegateOps, labeledgraph.WithCapacity(capacity))
	}
}

type void struct{}

// Graph is a graph with vertices of type V.
type Graph[V comparable] struct {
	delegate *labeledgraph.LabeledGraph[V, void]
}

// New returns a new Graph constructed according to the given options.
func New[V comparable](opts ...Option) *Graph[V] {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return &Graph[V]{
		delegate: labeledgraph.New[V, void](cfg.delegateOps...),
	}
}

// AddVertex adds the given vertex to the graph if not already present.
func (g *Graph[V]) AddVertex(v V) {
	g.delegate.AddVertex(v)
}

// RemoveVertex removes the given vertex and all incident edges from the graph.
func (g *Graph[V]) RemoveVertex(v V) {
	g.delegate.RemoveVertex(v)
}

// AddEdge adds the edge (u, v) to the graph, inserting missing vertices if needed.
func (g *Graph[V]) AddEdge(u, v V) {
	g.delegate.AddEdge(u, v, void{})
}

// RemoveEdge removes the edge (u, v) from the graph if present.
func (g *Graph[V]) RemoveEdge(u, v V) {
	g.delegate.RemoveEdge(u, v)
}

// Clear removes all vertices and edges from the graph.
func (g *Graph[V]) Clear() {
	g.delegate.Clear()
}

// ContainsVertex returns whether the given vertex is contained in the graph.
func (g *Graph[V]) ContainsVertex(v V) bool {
	return g.delegate.ContainsVertex(v)
}

// ContainsEdge returns whether the edge (u, v) is contained in the graph.
func (g *Graph[V]) ContainsEdge(u, v V) bool {
	return g.delegate.ContainsEdge(u, v)
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

// Degree returns the number of edges incident to the given vertex and true if the vertex exists in the graph.
// If the vertex is not in the graph, 0 and false are returned.
//   - For directed graphs, Degree is the sum of InDegree and OutDegree.
//   - For undirected graphs, Degree, InDegree, and OutDegree are all synonymous.
func (g *Graph[V]) Degree(u V) (int, bool) {
	return g.delegate.Degree(u)
}

// InDegree returns the number of edges coming into the given vertex and true if the vertex exists in the graph.
// If the vertex is not in the graph, 0 and false are returned.
//   - For directed graphs, InDegree is the number of incoming edges terminating at the vertex.
//   - For undirected graphs, Degree, InDegree, and OutDegree are all synonymous.
func (g *Graph[V]) InDegree(u V) (int, bool) {
	return g.delegate.InDegree(u)
}

// OutDegree returns the number of edges going out of the given vertex and true if the vertex exists in the graph.
// If the vertex is not in the graph, 0 and false are returned.
//   - For directed graphs, OutDegree is the number of outgoing edges originating at the vertex.
//   - For undirected graphs, Degree, InDegree, and OutDegree are all synonymous.
func (g *Graph[V]) OutDegree(u V) (int, bool) {
	return g.delegate.OutDegree(u)
}

// Clone returns a deep copy of the graph.
func (g *Graph[V]) Clone() *Graph[V] {
	return &Graph[V]{
		delegate: g.delegate.Clone(),
	}
}

// Equal returns true if the two graphs are structurally identical.
func (g *Graph[V]) Equal(other *Graph[V]) bool {
	if g == nil || other == nil {
		return g == other
	}
	return g.delegate.Equal(other.delegate, func(void, void) bool { return true })
}

// Vertices returns an iterator over all vertices in the graph.
func (g *Graph[V]) Vertices() iter.Seq[V] {
	return g.delegate.Vertices()
}

// Neighbors is an alias for Successors.
func (g *Graph[V]) Neighbors(v V) iter.Seq[V] {
	return g.delegate.Neighbors(v)
}

// Successors returns an iterator over the successors of the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, Successors is the set of vertices v such that (u, v) is an edge in the graph.
//   - For undirected graphs, Predecessors, Successors, and Neighbors are all synonymous.
func (g *Graph[V]) Successors(u V) iter.Seq[V] {
	return g.delegate.Successors(u)
}

// Predecessors returns an iterator over the predecessors of the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, Predecessors is the set of vertices u such that (u, v) is an edge in the graph.
//   - For undirected graphs, Predecessors, Successors, and Neighbors are all synonymous.
func (g *Graph[V]) Predecessors(v V) iter.Seq[V] {
	return g.delegate.Predecessors(v)
}

// Edges returns an iterator over all edges in the graph.
func (g *Graph[V]) Edges() iter.Seq2[V, V] {
	return g.delegate.Edges()
}

// IncidentEdges returns an iterator over all edges incident to the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, IncidentEdges is the union of InIncidentEdges and OutIncidentEdges.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g *Graph[V]) IncidentEdges(v V) iter.Seq2[V, V] {
	return g.delegate.IncidentEdges(v)
}

// InIncidentEdges returns an iterator over the edges coming into the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, this is the set of edges terminating at the given vertex.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g *Graph[V]) InIncidentEdges(v V) iter.Seq2[V, V] {
	return g.delegate.InIncidentEdges(v)
}

// OutIncidentEdges returns an iterator over the edges going out of the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, this is the set of edges originating at the given vertex.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g *Graph[V]) OutIncidentEdges(u V) iter.Seq2[V, V] {
	return g.delegate.OutIncidentEdges(u)
}

// String returns a string representation of the graph.
func (g *Graph[V]) String() string {
	if g.Order() == 0 {
		if g.Directed() {
			return "digraph[]"
		}
		return "graph[]"
	}
	var sb strings.Builder
	if g.Directed() {
		sb.WriteString("digraph[")
		first := true
		for u, v := range g.Edges() {
			if !first {
				sb.WriteString(", ")
			}
			first = false
			fmt.Fprintf(&sb, "%v -> %v", u, v)
		}
		for v := range g.Vertices() {
			if deg, _ := g.Degree(v); deg == 0 {
				if !first {
					sb.WriteString(", ")
				}
				first = false
				fmt.Fprint(&sb, v)
			}
		}
		sb.WriteByte(']')
		return sb.String()
	}

	sb.WriteString("graph[")
	first := true
	type edgePair struct{ u, v V }
	seen := make(map[edgePair]struct{}, g.Size())
	for u, v := range g.Edges() {
		if _, ok := seen[edgePair{v, u}]; ok {
			continue
		}
		seen[edgePair{u, v}] = struct{}{}
		if !first {
			sb.WriteString(", ")
		}
		first = false
		fmt.Fprintf(&sb, "%v - %v", u, v)
	}
	for v := range g.Vertices() {
		if deg, _ := g.Degree(v); deg == 0 {
			if !first {
				sb.WriteString(", ")
			}
			first = false
			fmt.Fprint(&sb, v)
		}
	}
	sb.WriteByte(']')
	return sb.String()
}
