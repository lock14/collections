package labeldgraph

import (
	"github.com/lock14/collections/hashset"
	"iter"
	"maps"
)

// Config holds configuration values for New to use when contracting a LabeledGraph.
type Config struct {
	directed bool
}

// Opt represents a configuration option for constructing a LabeledGraph.
type Opt func(g *Config)

// Directed is an option that configures New to return a directed graph.
func Directed() Opt {
	return func(g *Config) {
		g.directed = true
	}
}

// LabeledGraph is a graph with vertices of type V
// and whose edges are labeled with type L.
type LabeledGraph[V comparable, L any] struct {
	graph     map[V]nodeData[V, L]
	directed  bool
	edgeCount int
}

// New returns a new LabeledGraph constructed according to the given options.
func New[V comparable, L any](opts ...Opt) *LabeledGraph[V, L] {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return &LabeledGraph[V, L]{
		graph:     make(map[V]nodeData[V, L]),
		directed:  config.directed,
		edgeCount: 0,
	}
}

// AddVertex adds the given vertex to the graph.
func (g *LabeledGraph[V, L]) AddVertex(v V) {
	if !g.ContainsVertex(v) {
		g.graph[v] = g.nodeData()
	}
}

// AddEdge adds the given edge with the given label to the graph.
func (g *LabeledGraph[V, L]) AddEdge(u, v V, l L) {
	g.AddVertex(u)
	g.AddVertex(v)
	g.graph[u].AddSuccessor(v, l)
	g.graph[v].AddPredecessor(u, l)
	g.edgeCount++
}

// ContainsEdge returns whether the given edge is contained in the graph.
func (g *LabeledGraph[V, L]) ContainsEdge(u, v V) bool {
	return g.ContainsVertex(u) &&
		g.ContainsVertex(v) &&
		g.graph[u].ContainsSuccessor(v)
}

// ContainsVertex returns whether the given vertex is contained in the graph.
func (g *LabeledGraph[V, L]) ContainsVertex(v V) bool {
	_, ok := g.graph[v]
	return ok
}

// Directed returns whether or not this graph is directed.
func (g *LabeledGraph[V, L]) Directed() bool {
	return g.directed
}

// Order returns the number of vertices in the graph.
func (g *LabeledGraph[V, L]) Order() int {
	return len(g.graph)
}

// Size returns the number of vertices in the graph.
func (g *LabeledGraph[V, L]) Size() int {
	return g.edgeCount
}

// Label returns the label for edge (u, v) if the edge exists in the graph.
// If no such edge exists, the zero value for L is returned.
func (g *LabeledGraph[V, L]) Label(u, v V) (L, bool) {
	var l L
	var ok bool
	if g.ContainsVertex(u) {
		l, ok = g.graph[u].Label(v)
	}
	return l, ok
}

// InDegree return the number of edges coming into the specified vertex in the graph and true.
// If no such vertex exists, the zero value and false are returned.
func (g *LabeledGraph[V, L]) InDegree(u V) (int, bool) {
	var degree int
	var ok bool
	if g.ContainsVertex(u) {
		degree = g.graph[u].InDegree()
		ok = true
	}
	return degree, ok
}

// OutDegree return the number of edges coming out of the specified vertex in the graph and true.
// If no such vertex exists, the zero value and false are returned.
func (g *LabeledGraph[V, L]) OutDegree(u V) (int, bool) {
	var degree int
	var ok bool
	if g.ContainsVertex(u) {
		degree = g.graph[u].OutDegree()
		ok = true
	}
	return degree, ok
}

// Vertices returns an iterator over the vertices of the graph.
func (g *LabeledGraph[V, L]) Vertices() iter.Seq[V] {
	return maps.Keys(g.graph)
}

// Successors returns an iterator over the successors of the given vertex in the graph.
// If the vertex does not exist in the graph an empty iterator is returned.
func (g *LabeledGraph[V, L]) Successors(v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if g.ContainsVertex(v) {
			for u := range g.graph[v].Successors() {
				if !yield(u) {
					return
				}
			}
		}
	}
}

// Predecessors returns an iterator over the predecessors of the given vertex in the graph.
// If the vertex does not exist in the graph an empty iterator is returned.
func (g *LabeledGraph[V, L]) Predecessors(v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if g.ContainsVertex(v) {
			for u := range g.graph[v].Predecessors() {
				if !yield(u) {
					return
				}
			}
		}
	}
}

// Edges returns an iterator over the edges of the graph.
func (g *LabeledGraph[V, L]) Edges() iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for u := range g.Vertices() {
			for v := range g.Successors(u) {
				if !yield(u, v) {
					return
				}
			}
		}
	}
}

// IncidentEdgesIn returns an iterator over the edges coming into the given vertex of the graph.
// if no such vertex exists, then an empty iterator is returned.
func (g *LabeledGraph[V, L]) IncidentEdgesIn(u V) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for v := range g.Successors(u) {
			if !yield(u, v) {
				return
			}
		}
	}
}

// IncidentEdgesOut returns an iterator over the edges coming out of the given vertex of the graph.
// if no such vertex exists, then an empty iterator is returned.
func (g *LabeledGraph[V, L]) IncidentEdgesOut(v V) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for u := range g.Predecessors(v) {
			if !yield(u, v) {
				return
			}
		}
	}
}

func (g *LabeledGraph[V, L]) nodeData() nodeData[V, L] {
	if g.directed {
		return &directedNodeData[V, L]{
			successorLabels: make(map[V]L),
			predecessors:    hashset.New[V](),
		}
	} else {
		return &undirectedNodeData[V, L]{
			adjacentLabels: make(map[V]L),
		}
	}
}

func defaultConfig() *Config {
	return &Config{
		directed: false,
	}
}

type nodeData[V, L any] interface {
	AddPredecessor(V, L)
	AddSuccessor(V, L)
	ContainsPredecessor(V) bool
	ContainsSuccessor(V) bool
	InDegree() int
	Label(V) (L, bool)
	OutDegree() int
	Predecessors() iter.Seq[V]
	RemovePredecessor(V)
	RemoveSuccessor(V)
	Successors() iter.Seq[V]
}

type directedNodeData[V comparable, L any] struct {
	successorLabels map[V]L
	predecessors    *hashset.HashSet[V]
}

func (d *directedNodeData[V, L]) AddPredecessor(v V, l L) {
	// only successors store the label
	d.predecessors.Add(v)
}

func (d *directedNodeData[V, L]) AddSuccessor(v V, l L) {
	d.successorLabels[v] = l
}

func (d *directedNodeData[V, L]) ContainsPredecessor(v V) bool {
	return d.predecessors.Contains(v)
}

func (d *directedNodeData[V, L]) ContainsSuccessor(v V) bool {
	_, ok := d.successorLabels[v]
	return ok
}

func (d *directedNodeData[V, L]) InDegree() int {
	return d.predecessors.Size()
}

func (d *directedNodeData[V, L]) Label(v V) (L, bool) {
	l, ok := d.successorLabels[v]
	return l, ok
}

func (d *directedNodeData[V, L]) OutDegree() int {
	return len(d.successorLabels)
}

func (d *directedNodeData[V, L]) Predecessors() iter.Seq[V] {
	return d.predecessors.All()
}

func (d *directedNodeData[V, L]) RemovePredecessor(v V) {
	d.predecessors.Remove(v)
}

func (d *directedNodeData[V, L]) RemoveSuccessor(v V) {
	delete(d.successorLabels, v)
}

func (d *directedNodeData[V, L]) Successors() iter.Seq[V] {
	return maps.Keys(d.successorLabels)
}

type undirectedNodeData[V comparable, L any] struct {
	adjacentLabels map[V]L
}

func (u *undirectedNodeData[V, L]) AddPredecessor(v V, l L) {
	u.AddSuccessor(v, l)
}

func (u *undirectedNodeData[V, L]) AddSuccessor(v V, l L) {
	u.adjacentLabels[v] = l
}

func (u *undirectedNodeData[V, L]) ContainsPredecessor(v V) bool {
	return u.ContainsSuccessor(v)
}

func (u *undirectedNodeData[V, L]) ContainsSuccessor(v V) bool {
	_, ok := u.adjacentLabels[v]
	return ok
}

func (u *undirectedNodeData[V, L]) InDegree() int {
	return u.OutDegree()
}

func (u *undirectedNodeData[V, L]) Label(v V) (L, bool) {
	l, ok := u.adjacentLabels[v]
	return l, ok
}

func (u *undirectedNodeData[V, L]) OutDegree() int {
	return len(u.adjacentLabels)
}

func (u *undirectedNodeData[V, L]) Predecessors() iter.Seq[V] {
	return u.Successors()
}

func (u *undirectedNodeData[V, L]) RemovePredecessor(v V) {
	u.RemoveSuccessor(v)
}

func (u *undirectedNodeData[V, L]) RemoveSuccessor(v V) {
	delete(u.adjacentLabels, v)
}

func (u *undirectedNodeData[V, L]) Successors() iter.Seq[V] {
	return maps.Keys(u.adjacentLabels)
}
