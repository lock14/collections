// Package labeledgraph provides graph data structures with edge labels.
package labeledgraph

import (
	"fmt"
	"github.com/lock14/collections/hashset"
	"iter"
	"maps"
	"strings"
)

// config holds configuration values for New to use when contracting a LabeledGraph.
type config struct {
	directed bool
	capacity int
}

// Opt represents a configuration option for constructing a LabeledGraph.
type Opt func(g *config)

// Directed is an option that configures New to return a directed graph.
func Directed() Opt {
	return func(g *config) {
		g.directed = true
	}
}

// Capacity is an option that configures New to pre-allocate the graph with the given capacity.
func Capacity(n int) Opt {
	return func(g *config) {
		g.capacity = n
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
		graph:     make(map[V]nodeData[V, L], config.capacity),
		directed:  config.directed,
		edgeCount: 0,
	}
}

// ContainsVertex returns whether the given vertex is contained in the graph.
func (g *LabeledGraph[V, L]) ContainsVertex(v V) bool {
	_, ok := g.graph[v]
	return ok
}

// AddVertex adds the given vertex to the graph.
func (g *LabeledGraph[V, L]) AddVertex(v V) {
	if !g.ContainsVertex(v) {
		g.graph[v] = g.nodeData()
	}
}

// RemoveVertex removes the vertex and all incident edges from the graph.
func (g *LabeledGraph[V, L]) RemoveVertex(v V) {
	if g.ContainsVertex(v) {
		if g.directed {
			g.edgeCount -= g.graph[v].outDegree()
			g.edgeCount -= g.graph[v].inDegree()
			if g.ContainsEdge(v, v) {
				g.edgeCount++
			}
		} else {
			g.edgeCount -= g.graph[v].outDegree()
		}

		for u := range g.graph[v].predecessors() {
			g.graph[u].removeSuccessor(v)
		}
		for u := range g.graph[v].successors() {
			g.graph[u].removePredecessor(v)
		}
		delete(g.graph, v)
	}
}

// ContainsEdge returns whether the given edge is contained in the graph.
func (g *LabeledGraph[V, L]) ContainsEdge(u, v V) bool {
	return g.ContainsVertex(u) &&
		g.ContainsVertex(v) &&
		g.graph[u].containsSuccessor(v)
}

// AddEdge adds the given edge with the given label to the graph.
func (g *LabeledGraph[V, L]) AddEdge(u, v V, l L) {
	g.AddVertex(u)
	g.AddVertex(v)
	if !g.ContainsEdge(u, v) {
		g.edgeCount++
	}
	g.graph[u].addSuccessor(v, l)
	g.graph[v].addPredecessor(u, l)
}

// RemoveEdge removes the given edge from the graph.
func (g *LabeledGraph[V, L]) RemoveEdge(u, v V) {
	if g.ContainsEdge(u, v) {
		g.graph[u].removeSuccessor(v)
		g.graph[v].removePredecessor(u)
		g.edgeCount--
	}
}

// SetLabel assigns the given label to the given edge in the graph.
func (g *LabeledGraph[V, L]) SetLabel(u, v V, label L) {
	if g.ContainsEdge(u, v) {
		g.graph[u].setSuccessorLabel(v, label)
		g.graph[v].setPredecessorLabel(u, label)
	}
}

// Directed returns whether this graph is directed.
func (g *LabeledGraph[V, L]) Directed() bool {
	return g.directed
}

// Order returns the number of vertices in the graph.
func (g *LabeledGraph[V, L]) Order() int {
	return len(g.graph)
}

// Size returns the number of edges in the graph.
func (g *LabeledGraph[V, L]) Size() int {
	return g.edgeCount
}

// Label returns the label for edge (u, v) and true if the edge exists in the graph.
// If no such edge exists, the zero value and false are returned.
func (g *LabeledGraph[V, L]) Label(u, v V) (L, bool) {
	var l L
	var ok bool
	if g.ContainsVertex(u) {
		l, ok = g.graph[u].label(v)
	}
	return l, ok
}

// Degree return the number of edges coming into or out of the given vertex
// and true if the vertex exists in the graph. If the given vertex is not in
// the graph, an empty iterator is returned.
//   - For directed graphs, Degree is the sum of InDegree and OutDegree.
//   - For undirected graphs, Degree, InDegree, and OutDegree are all synonymous.
func (g *LabeledGraph[V, L]) Degree(u V) (int, bool) {
	if g.directed {
		var degree int
		var ok bool
		if g.ContainsVertex(u) {
			degree = g.graph[u].inDegree() + g.graph[u].outDegree()
			ok = true
		}
		return degree, ok
	} else {
		return g.OutDegree(u)
	}
}

// InDegree return the number of edges coming into the given vertex
// and true if the vertex exists in the graph.
func (g *LabeledGraph[V, L]) InDegree(u V) (int, bool) {
	var degree int
	var ok bool
	if g.ContainsVertex(u) {
		degree = g.graph[u].inDegree()
		ok = true
	}
	return degree, ok
}

// OutDegree return the number of edges coming out of the given vertex
// and true if the vertex exists in the graph. If no such vertex exists,
// the zero value and false are returned.
func (g *LabeledGraph[V, L]) OutDegree(u V) (int, bool) {
	var degree int
	var ok bool
	if g.ContainsVertex(u) {
		degree = g.graph[u].outDegree()
		ok = true
	}
	return degree, ok
}

// Vertices returns an iterator over all vertices in the graph.
func (g *LabeledGraph[V, L]) Vertices() iter.Seq[V] {
	return maps.Keys(g.graph)
}

// Neighbors is an alias for Successors.
func (g *LabeledGraph[V, L]) Neighbors(v V) iter.Seq[V] {
	return g.Successors(v)
}

// Successors returns an iterator over the successors of the given vertex in the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, Successors is the set of vertices u, such that (u, v) is
//     an edge in the graph for the given vertex u.
//   - For undirected graphs, Predecessors, Successors, and Neighbors are all synonymous.
func (g *LabeledGraph[V, L]) Successors(u V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if g.ContainsVertex(u) {
			for v := range g.graph[u].successors() {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Predecessors returns an iterator over the predecessors of the given vertex in the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, Predecessors is the set of vertices u, such that (u, v) is
//     an edge in the graph for the given vertex v.
//   - For undirected graphs, Predecessors, Successors, and Neighbors are all synonymous.
func (g *LabeledGraph[V, L]) Predecessors(v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if g.ContainsVertex(v) {
			for u := range g.graph[v].predecessors() {
				if !yield(u) {
					return
				}
			}
		}
	}
}

// Edges returns an iterator over all edges in the graph.
func (g *LabeledGraph[V, L]) Edges() iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for u, node := range g.graph {
			if g.directed {
				dNode := node.(*directedNodeData[V, L])
				for v := range dNode.successorLabels {
					if !yield(u, v) {
						return
					}
				}
			} else {
				uNode := node.(*undirectedNodeData[V, L])
				for v := range uNode.adjacentLabels {
					if !yield(u, v) {
						return
					}
				}
			}
		}
	}
}

// IncidentEdges returns all edges that are incident to the given vertex in the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, IncidentEdges is the union of InIncidentEdges and OutIncidentEdges.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g *LabeledGraph[V, L]) IncidentEdges(v V) iter.Seq2[V, V] {
	if g.directed {
		return func(yield func(V, V) bool) {
			for x, y := range g.InIncidentEdges(v) {
				if !yield(x, y) {
					return
				}
			}
			for x, y := range g.OutIncidentEdges(v) {
				if !yield(x, y) {
					return
				}
			}
		}
	} else {
		return g.OutIncidentEdges(v)
	}
}

// InIncidentEdges returns an iterator over the edges coming into the given vertex of the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, this is the set of edges that terminate at the given vertex.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g *LabeledGraph[V, L]) InIncidentEdges(v V) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for u := range g.Predecessors(v) {
			if !yield(u, v) {
				return
			}
		}
	}
}

// OutIncidentEdges returns an iterator over the edges coming out of the given vertex of the graph.
// If the given vertex is not in the graph, an empty iterator is returned.
//   - For directed graphs, this is the set of edges that originate at the given vertex.
//   - For undirected graphs, IncidentEdges, InIncidentEdges, and OutIncidentEdges are all synonymous.
func (g *LabeledGraph[V, L]) OutIncidentEdges(u V) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for v := range g.Successors(u) {
			if !yield(u, v) {
				return
			}
		}
	}
}

// Clear removes all vertices and edges from the graph.
func (g *LabeledGraph[V, L]) Clear() {
	clear(g.graph)
	g.edgeCount = 0
}

// Clone returns a deep copy of the graph.
func (g *LabeledGraph[V, L]) Clone() *LabeledGraph[V, L] {
	clone := &LabeledGraph[V, L]{
		graph:     make(map[V]nodeData[V, L], len(g.graph)),
		directed:  g.directed,
		edgeCount: g.edgeCount,
	}

	for v, data := range g.graph {
		if g.directed {
			dNode := data.(*directedNodeData[V, L])
			clonedSuccessors := make(map[V]L, len(dNode.successorLabels))
			for s, l := range dNode.successorLabels {
				clonedSuccessors[s] = l
			}
			clonedPreds := hashset.New[V]()
			for p := range dNode.preds.All() {
				clonedPreds.Add(p)
			}
			clone.graph[v] = &directedNodeData[V, L]{
				successorLabels: clonedSuccessors,
				preds:           clonedPreds,
			}
		} else {
			uNode := data.(*undirectedNodeData[V, L])
			clonedAdjacent := make(map[V]L, len(uNode.adjacentLabels))
			for s, l := range uNode.adjacentLabels {
				clonedAdjacent[s] = l
			}
			clone.graph[v] = &undirectedNodeData[V, L]{
				adjacentLabels: clonedAdjacent,
			}
		}
	}
	return clone
}

// Equal returns true if the two graphs are structurally identical and all edge labels are equal according to the provided eq function.
func (g *LabeledGraph[V, L]) Equal(other *LabeledGraph[V, L], eq func(L, L) bool) bool {
	if g.directed != other.directed || g.Order() != other.Order() || g.Size() != other.Size() {
		return false
	}

	for v := range g.Vertices() {
		if !other.ContainsVertex(v) {
			return false
		}
	}

	for u, v := range g.Edges() {
		if !other.ContainsEdge(u, v) {
			return false
		}
		l1, _ := g.Label(u, v)
		l2, _ := other.Label(u, v)
		if !eq(l1, l2) {
			return false
		}
	}

	return true
}

// String returns a string representation of the graph.
func (g *LabeledGraph[V, L]) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	first := true

	type edgePair struct{ u, v V }
	seen := make(map[edgePair]bool)

	for u, v := range g.Edges() {
		if !g.directed {
			if seen[edgePair{u, v}] || seen[edgePair{v, u}] {
				continue
			}
			seen[edgePair{u, v}] = true
		}

		if !first {
			sb.WriteString(", ")
		}
		l, _ := g.Label(u, v)
		if g.directed {
			sb.WriteString(fmt.Sprintf("%v -> %v: %v", u, v, l))
		} else {
			sb.WriteString(fmt.Sprintf("%v - %v: %v", u, v, l))
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

func (g *LabeledGraph[V, L]) nodeData() nodeData[V, L] {
	if g.directed {
		return &directedNodeData[V, L]{
			successorLabels: make(map[V]L),
			preds:           hashset.New[V](),
		}
	} else {
		return &undirectedNodeData[V, L]{
			adjacentLabels: make(map[V]L),
		}
	}
}

func defaultConfig() *config {
	return &config{
		directed: false,
	}
}

type nodeData[V, L any] interface {
	addPredecessor(V, L)
	addSuccessor(V, L)
	containsPredecessor(V) bool
	containsSuccessor(V) bool
	inDegree() int
	label(V) (L, bool)
	outDegree() int
	predecessors() iter.Seq[V]
	removePredecessor(V)
	removeSuccessor(V)
	setPredecessorLabel(V, L)
	setSuccessorLabel(V, L)
	successors() iter.Seq[V]
}

type directedNodeData[V comparable, L any] struct {
	successorLabels map[V]L
	preds           *hashset.HashSet[V]
}

func (d *directedNodeData[V, L]) addPredecessor(v V, l L) {
	// only successors store the label
	d.preds.Add(v)
}

func (d *directedNodeData[V, L]) addSuccessor(v V, l L) {
	d.successorLabels[v] = l
}

func (d *directedNodeData[V, L]) containsPredecessor(v V) bool {
	return d.preds.Contains(v)
}

func (d *directedNodeData[V, L]) containsSuccessor(v V) bool {
	_, ok := d.successorLabels[v]
	return ok
}

func (d *directedNodeData[V, L]) inDegree() int {
	return d.preds.Size()
}

func (d *directedNodeData[V, L]) label(v V) (L, bool) {
	l, ok := d.successorLabels[v]
	return l, ok
}

func (d *directedNodeData[V, L]) outDegree() int {
	return len(d.successorLabels)
}

func (d *directedNodeData[V, L]) predecessors() iter.Seq[V] {
	return d.preds.All()
}

func (d *directedNodeData[V, L]) removePredecessor(v V) {
	d.preds.RemoveElement(v)
}

func (d *directedNodeData[V, L]) removeSuccessor(v V) {
	delete(d.successorLabels, v)
}

func (d *directedNodeData[V, L]) setPredecessorLabel(v V, label L) {
	// no op
}

func (d *directedNodeData[V, L]) setSuccessorLabel(v V, label L) {
	d.successorLabels[v] = label
}

func (d *directedNodeData[V, L]) successors() iter.Seq[V] {
	return maps.Keys(d.successorLabels)
}

type undirectedNodeData[V comparable, L any] struct {
	adjacentLabels map[V]L
}

func (u *undirectedNodeData[V, L]) addPredecessor(v V, l L) {
	u.addSuccessor(v, l)
}

func (u *undirectedNodeData[V, L]) addSuccessor(v V, l L) {
	u.adjacentLabels[v] = l
}

func (u *undirectedNodeData[V, L]) containsPredecessor(v V) bool {
	return u.containsSuccessor(v)
}

func (u *undirectedNodeData[V, L]) containsSuccessor(v V) bool {
	_, ok := u.adjacentLabels[v]
	return ok
}

func (u *undirectedNodeData[V, L]) inDegree() int {
	return u.outDegree()
}

func (u *undirectedNodeData[V, L]) label(v V) (L, bool) {
	l, ok := u.adjacentLabels[v]
	return l, ok
}

func (u *undirectedNodeData[V, L]) outDegree() int {
	return len(u.adjacentLabels)
}

func (u *undirectedNodeData[V, L]) predecessors() iter.Seq[V] {
	return u.successors()
}

func (u *undirectedNodeData[V, L]) removePredecessor(v V) {
	u.removeSuccessor(v)
}

func (u *undirectedNodeData[V, L]) removeSuccessor(v V) {
	delete(u.adjacentLabels, v)
}

func (u *undirectedNodeData[V, L]) setPredecessorLabel(v V, label L) {
	u.setSuccessorLabel(v, label)
}

func (u *undirectedNodeData[V, L]) setSuccessorLabel(v V, label L) {
	u.adjacentLabels[v] = label
}

func (u *undirectedNodeData[V, L]) successors() iter.Seq[V] {
	return maps.Keys(u.adjacentLabels)
}
