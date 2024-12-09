package graph

import (
	"github.com/lock14/collections/labeldgraph"
	"iter"
)

// Config holds configuration values for New to use when contracting a LabeledGraph.
type Config struct {
	delegateOps []labeldgraph.Opt
}

// Opt represents a configuration option for constructing a LabeledGraph.
type Opt func(g *Config)

// Directed is an option that configures New to return a directed graph.
func Directed() Opt {
	return func(g *Config) {
		g.delegateOps = append(g.delegateOps, labeldgraph.Directed())
	}
}

type void struct{}

// Graph is a graph with vertices of type V.
type Graph[V comparable] struct {
	delegate *labeldgraph.LabeledGraph[V, void]
}

// New returns a new LabeledGraph constructed according to the given options.
func New[V comparable](opts ...Opt) *Graph[V] {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return &Graph[V]{
		delegate: labeldgraph.New[V, void](config.delegateOps...),
	}
}

// ContainsVertex returns whether the given vertex is contained in the graph.
func (g Graph[V]) ContainsVertex(v V) bool {
	return g.delegate.ContainsVertex(v)
}

// AddVertex adds the given vertex to the graph.
func (g Graph[V]) AddVertex(v V) {
	g.delegate.AddVertex(v)
}

// RemoveVertex removes the vertex and all incident edges from the graph.
func (g Graph[V]) RemoveVertex(v V) {
	g.delegate.RemoveVertex(v)
}

// ContainsEdge returns whether the given edge is contained in the graph.
func (g Graph[V]) ContainsEdge(u, v V) bool {
	return g.delegate.ContainsEdge(u, v)
}

// AddEdge adds the given edge with the given label to the graph.
func (g Graph[V]) AddEdge(u, v V) {
	g.delegate.AddEdge(u, v, void{})
}

// RemoveEdge removes the given edge from the graph.
func (g Graph[V]) RemoveEdge(u, v V) {
	g.delegate.RemoveEdge(u, v)
}

// Directed returns whether this graph is directed.
func (g Graph[V]) Directed() bool {
	return g.delegate.Directed()
}

// Order returns the number of vertices in the graph.
func (g Graph[V]) Order() int {
	return g.delegate.Order()
}

// Size returns the number of edges in the graph.
func (g Graph[V]) Size() int {
	return g.delegate.Size()
}

// Degree return the number of edges coming into or out of the given vertex
// and true if the vertex exists in the graph. If the given vertex is not in
// the graph, an empty iterator is returned.
//   - For directed graphs, Degree is the sum of InDegree and OutDegree.
//   - For undirected graphs, Degree, InDegree, and OutDegree are all synonymous.
func (g Graph[V]) Degree(u V) (int, bool) {
	return g.delegate.Degree(u)
}

// InDegree return the number of edges coming into the given vertex
// and true if the vertex exists in the graph.
func (g Graph[V]) InDegree(u V) (int, bool) {
	return g.delegate.InDegree(u)
}

// OutDegree return the number of edges coming out of the given vertex
// and true if the vertex exists in the graph. If no such vertex exists,
// the zero value and false are returned.
func (g Graph[V]) OutDegree(u V) (int, bool) {
	return g.delegate.OutDegree(u)
}

// Vertices returns an iterator over all vertices in the graph.
func (g Graph[V]) Vertices() iter.Seq[V] {
	return g.delegate.Vertices()
}

// Neighbors is an alias for Successors.
func (g Graph[V]) Neighbors(v V) iter.Seq[V] {
	return g.delegate.Neighbors(v)
}

// Successors returns an iterator over the successors of the given vertex in the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, Successors is the set of vertices u, such that (u, v) is
//     an edge in the graph for the given vertex u.
//   - For undirected graphs, Predecessors, Successors, and Neighbors are all synonymous.
func (g Graph[V]) Successors(u V) iter.Seq[V] {
	return g.delegate.Successors(u)
}

// Predecessors returns an iterator over the predecessors of the given vertex in the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, Predecessors is the set of vertices u, such that (u, v) is
//     an edge in the graph for the given vertex v.
//   - For undirected graphs, Predecessors, Successors, and Neighbors are all synonymous.
func (g Graph[V]) Predecessors(v V) iter.Seq[V] {
	return g.delegate.Predecessors(v)
}

// Edges returns an iterator over all edges in the graph.
func (g Graph[V]) Edges() iter.Seq2[V, V] {
	return g.delegate.Edges()
}

// IncidentEdges returns all edges that are incident to the given vertex in the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, IncidentEdges is the union of InIncidentEdges and OutIncidentEdges.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g Graph[V]) IncidentEdges(v V) iter.Seq2[V, V] {
	return g.delegate.IncidentEdges(v)
}

// InIncidentEdges returns an iterator over the edges coming into the given vertex of the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, this is the set of edges that terminate at the given vertex.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g Graph[V]) InIncidentEdges(v V) iter.Seq2[V, V] {
	return g.delegate.IncidentEdges(v)
}

// OutIncidentEdges returns an iterator over the edges coming out of the given vertex of the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, this is the set of edges that originate at the given vertex.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g Graph[V]) OutIncidentEdges(u V) iter.Seq2[V, V] {
	return g.delegate.OutIncidentEdges(u)
}

func defaultConfig() *Config {
	return &Config{
		delegateOps: []labeldgraph.Opt{labeldgraph.Directed()},
	}
}
