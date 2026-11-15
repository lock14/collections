// Package labeledgraph provides directed and undirected graph data structures with labeled edges.
package labeledgraph

import (
	"fmt"
	"iter"
	"maps"
	"strings"
)

var _ fmt.Stringer = (*LabeledGraph[int, string])(nil)

// config holds configuration values for New to use when constructing a LabeledGraph.
type config struct {
	directed bool
	capacity int
}

// Option represents a configuration option for constructing a LabeledGraph.
type Option func(c *config)

// WithDirected configures New to return a directed graph.
func WithDirected() Option {
	return func(c *config) {
		c.directed = true
	}
}

// WithCapacity configures New to pre-allocate the graph with the given vertex capacity.
// Panics in New if capacity is negative.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

type edges[V comparable, L any] struct {
	out map[V]L
	in  map[V]struct{} // nil when !g.isDirected
}

// LabeledGraph is a graph with vertices of type V and edges labeled with type L.
type LabeledGraph[V comparable, L any] struct {
	nodes      map[V]edges[V, L]
	isDirected bool
	edgeCount  int
}

// New returns a new LabeledGraph constructed according to the given options.
// Panics if configured with a negative capacity.
func New[V comparable, L any](opts ...Option) *LabeledGraph[V, L] {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.capacity < 0 {
		panic("labeledgraph: capacity cannot be negative")
	}
	return &LabeledGraph[V, L]{
		nodes:      make(map[V]edges[V, L], cfg.capacity),
		isDirected: cfg.directed,
	}
}

// AddVertex adds the given vertex to the graph if not already present.
func (g *LabeledGraph[V, L]) AddVertex(v V) {
	if _, ok := g.nodes[v]; !ok {
		g.nodes[v] = g.newEdges()
	}
}

// RemoveVertex removes the given vertex and all incident edges from the graph.
func (g *LabeledGraph[V, L]) RemoveVertex(v V) {
	ev, ok := g.nodes[v]
	if !ok {
		return
	}
	if g.isDirected {
		g.edgeCount -= len(ev.out) + len(ev.in)
		if _, selfLoop := ev.out[v]; selfLoop {
			g.edgeCount++
		}
		for u := range ev.in {
			delete(g.nodes[u].out, v)
		}
		for w := range ev.out {
			delete(g.nodes[w].in, v)
		}
	} else {
		g.edgeCount -= len(ev.out)
		for u := range ev.out {
			delete(g.nodes[u].out, v)
		}
	}
	delete(g.nodes, v)
}

// AddEdge adds the edge (u, v) with the given label to the graph, inserting missing vertices if needed.
func (g *LabeledGraph[V, L]) AddEdge(u, v V, label L) {
	eu, ok := g.nodes[u]
	if !ok {
		eu = g.newEdges()
		g.nodes[u] = eu
	}
	ev, ok := g.nodes[v]
	if !ok {
		ev = g.newEdges()
		g.nodes[v] = ev
	}
	if _, exists := eu.out[v]; !exists {
		g.edgeCount++
	}
	eu.out[v] = label
	if g.isDirected {
		ev.in[u] = struct{}{}
	} else {
		ev.out[u] = label
	}
}

// RemoveEdge removes the edge (u, v) from the graph if present.
func (g *LabeledGraph[V, L]) RemoveEdge(u, v V) {
	eu, ok := g.nodes[u]
	if !ok {
		return
	}
	if _, exists := eu.out[v]; !exists {
		return
	}
	delete(eu.out, v)
	if g.isDirected {
		delete(g.nodes[v].in, u)
	} else {
		delete(g.nodes[v].out, u)
	}
	g.edgeCount--
}

// SetLabel assigns the given label to the existing edge (u, v) in the graph.
// If the edge does not exist, SetLabel is a no-op.
func (g *LabeledGraph[V, L]) SetLabel(u, v V, label L) {
	eu, ok := g.nodes[u]
	if !ok {
		return
	}
	if _, exists := eu.out[v]; !exists {
		return
	}
	eu.out[v] = label
	if !g.isDirected {
		g.nodes[v].out[u] = label
	}
}

// Clear removes all vertices and edges from the graph.
func (g *LabeledGraph[V, L]) Clear() {
	clear(g.nodes)
	g.edgeCount = 0
}

// ContainsVertex returns whether the given vertex is contained in the graph.
func (g *LabeledGraph[V, L]) ContainsVertex(v V) bool {
	_, ok := g.nodes[v]
	return ok
}

// ContainsEdge returns whether the edge (u, v) is contained in the graph.
func (g *LabeledGraph[V, L]) ContainsEdge(u, v V) bool {
	eu, ok := g.nodes[u]
	if !ok {
		return false
	}
	_, exists := eu.out[v]
	return exists
}

// Label returns the label for edge (u, v) and true if the edge exists in the graph.
// If no such edge exists, the zero value of L and false are returned.
func (g *LabeledGraph[V, L]) Label(u, v V) (L, bool) {
	eu, ok := g.nodes[u]
	if !ok {
		var zero L
		return zero, false
	}
	l, exists := eu.out[v]
	return l, exists
}

// Directed returns whether this graph is directed.
func (g *LabeledGraph[V, L]) Directed() bool {
	return g.isDirected
}

// Order returns the number of vertices in the graph.
func (g *LabeledGraph[V, L]) Order() int {
	return len(g.nodes)
}

// Size returns the number of edges in the graph.
func (g *LabeledGraph[V, L]) Size() int {
	return g.edgeCount
}

// Degree returns the number of edges incident to the given vertex and true if the vertex exists in the graph.
// If the vertex is not in the graph, 0 and false are returned.
//   - For directed graphs, Degree is the sum of InDegree and OutDegree.
//   - For undirected graphs, Degree, InDegree, and OutDegree are all synonymous.
func (g *LabeledGraph[V, L]) Degree(u V) (int, bool) {
	eu, ok := g.nodes[u]
	if !ok {
		return 0, false
	}
	if g.isDirected {
		return len(eu.in) + len(eu.out), true
	}
	return len(eu.out), true
}

// InDegree returns the number of edges coming into the given vertex and true if the vertex exists in the graph.
// If the vertex is not in the graph, 0 and false are returned.
//   - For directed graphs, InDegree is the number of incoming edges terminating at the vertex.
//   - For undirected graphs, Degree, InDegree, and OutDegree are all synonymous.
func (g *LabeledGraph[V, L]) InDegree(u V) (int, bool) {
	eu, ok := g.nodes[u]
	if !ok {
		return 0, false
	}
	if g.isDirected {
		return len(eu.in), true
	}
	return len(eu.out), true
}

// OutDegree returns the number of edges going out of the given vertex and true if the vertex exists in the graph.
// If the vertex is not in the graph, 0 and false are returned.
//   - For directed graphs, OutDegree is the number of outgoing edges originating at the vertex.
//   - For undirected graphs, Degree, InDegree, and OutDegree are all synonymous.
func (g *LabeledGraph[V, L]) OutDegree(u V) (int, bool) {
	eu, ok := g.nodes[u]
	if !ok {
		return 0, false
	}
	return len(eu.out), true
}

// Clone returns a deep copy of the graph.
func (g *LabeledGraph[V, L]) Clone() *LabeledGraph[V, L] {
	cloned := &LabeledGraph[V, L]{
		nodes:      make(map[V]edges[V, L], len(g.nodes)),
		isDirected: g.isDirected,
		edgeCount:  g.edgeCount,
	}
	for v, e := range g.nodes {
		cloned.nodes[v] = edges[V, L]{
			out: maps.Clone(e.out),
			in:  maps.Clone(e.in),
		}
	}
	return cloned
}

// Equal returns true if the two graphs are structurally identical and all edge labels are equal
// according to the provided labelEq function.
func (g *LabeledGraph[V, L]) Equal(other *LabeledGraph[V, L], labelEq func(L, L) bool) bool {
	if g == nil || other == nil {
		return g == other
	}
	if g.isDirected != other.isDirected || len(g.nodes) != len(other.nodes) || g.edgeCount != other.edgeCount {
		return false
	}
	for u, eu := range g.nodes {
		otherEu, ok := other.nodes[u]
		if !ok || len(eu.out) != len(otherEu.out) {
			return false
		}
		for v, l1 := range eu.out {
			l2, exists := otherEu.out[v]
			if !exists || !labelEq(l1, l2) {
				return false
			}
		}
	}
	return true
}

// Vertices returns an iterator over all vertices in the graph.
func (g *LabeledGraph[V, L]) Vertices() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range g.nodes {
			if !yield(v) {
				return
			}
		}
	}
}

// Neighbors is an alias for Successors.
func (g *LabeledGraph[V, L]) Neighbors(v V) iter.Seq[V] {
	return g.Successors(v)
}

// Successors returns an iterator over the successors of the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, Successors is the set of vertices v such that (u, v) is an edge in the graph.
//   - For undirected graphs, Predecessors, Successors, and Neighbors are all synonymous.
func (g *LabeledGraph[V, L]) Successors(u V) iter.Seq[V] {
	return func(yield func(V) bool) {
		eu, ok := g.nodes[u]
		if !ok {
			return
		}
		for v := range eu.out {
			if !yield(v) {
				return
			}
		}
	}
}

// Predecessors returns an iterator over the predecessors of the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, Predecessors is the set of vertices u such that (u, v) is an edge in the graph.
//   - For undirected graphs, Predecessors, Successors, and Neighbors are all synonymous.
func (g *LabeledGraph[V, L]) Predecessors(v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		ev, ok := g.nodes[v]
		if !ok {
			return
		}
		if g.isDirected {
			for u := range ev.in {
				if !yield(u) {
					return
				}
			}
			return
		}
		for u := range ev.out {
			if !yield(u) {
				return
			}
		}
	}
}

// Edges returns an iterator over all edges in the graph.
func (g *LabeledGraph[V, L]) Edges() iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for u, e := range g.nodes {
			for v := range e.out {
				if !yield(u, v) {
					return
				}
			}
		}
	}
}

// IncidentEdges returns an iterator over all edges incident to the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, IncidentEdges is the union of InIncidentEdges and OutIncidentEdges.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g *LabeledGraph[V, L]) IncidentEdges(v V) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		ev, ok := g.nodes[v]
		if !ok {
			return
		}
		if g.isDirected {
			for u := range ev.in {
				if !yield(u, v) {
					return
				}
			}
			for w := range ev.out {
				if w == v {
					continue
				}
				if !yield(v, w) {
					return
				}
			}
			return
		}
		for w := range ev.out {
			if !yield(v, w) {
				return
			}
		}
	}
}

// InIncidentEdges returns an iterator over the edges coming into the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, this is the set of edges terminating at the given vertex.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g *LabeledGraph[V, L]) InIncidentEdges(v V) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		ev, ok := g.nodes[v]
		if !ok {
			return
		}
		if g.isDirected {
			for u := range ev.in {
				if !yield(u, v) {
					return
				}
			}
			return
		}
		for u := range ev.out {
			if !yield(u, v) {
				return
			}
		}
	}
}

// OutIncidentEdges returns an iterator over the edges going out of the given vertex in the graph,
// or an empty sequence if the vertex is not in the graph.
//   - For directed graphs, this is the set of edges originating at the given vertex.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g *LabeledGraph[V, L]) OutIncidentEdges(u V) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		eu, ok := g.nodes[u]
		if !ok {
			return
		}
		for v := range eu.out {
			if !yield(u, v) {
				return
			}
		}
	}
}

// String returns a string representation of the graph.
func (g *LabeledGraph[V, L]) String() string {
	if len(g.nodes) == 0 {
		if g.isDirected {
			return "digraph[]"
		}
		return "graph[]"
	}
	var sb strings.Builder
	if g.isDirected {
		sb.WriteString("digraph[")
		first := true
		for u, e := range g.nodes {
			for v, l := range e.out {
				if !first {
					sb.WriteString(", ")
				}
				first = false
				fmt.Fprintf(&sb, "%v -> %v: %v", u, v, l)
			}
		}
		for v, e := range g.nodes {
			if len(e.out) == 0 && len(e.in) == 0 {
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
	seen := make(map[edgePair]struct{}, g.edgeCount)
	for u, e := range g.nodes {
		for v, l := range e.out {
			if _, ok := seen[edgePair{v, u}]; ok {
				continue
			}
			seen[edgePair{u, v}] = struct{}{}
			if !first {
				sb.WriteString(", ")
			}
			first = false
			fmt.Fprintf(&sb, "%v - %v: %v", u, v, l)
		}
	}
	for v, e := range g.nodes {
		if len(e.out) == 0 {
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

func (g *LabeledGraph[V, L]) newEdges() edges[V, L] {
	if g.isDirected {
		return edges[V, L]{
			out: make(map[V]L),
			in:  make(map[V]struct{}),
		}
	}
	return edges[V, L]{
		out: make(map[V]L),
	}
}
