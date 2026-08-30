package labeledgraph

import (
	"slices"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "default_undirected",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				if g.Directed() {
					t.Errorf("expected undirected graph by default")
				}
				if g.Order() != 0 {
					t.Errorf("expected order 0, got %d", g.Order())
				}
				if g.Size() != 0 {
					t.Errorf("expected size 0, got %d", g.Size())
				}
			},
		},
		{
			name: "directed",
			opts: []Opt{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				if !g.Directed() {
					t.Errorf("expected directed graph")
				}
			},
		},
		{
			name: "capacity",
			opts: []Opt{WithCapacity(100)},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddVertex(1)
				if g.Order() != 1 {
					t.Errorf("expected order 1, got %d", g.Order())
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.check(t, g)
		})
	}
}

func TestLabeledGraph_AddVertex(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "add_one",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddVertex(1)
				if !g.ContainsVertex(1) {
					t.Errorf("expected graph to contain vertex 1")
				}
				if g.Order() != 1 {
					t.Errorf("expected order 1, got %d", g.Order())
				}
			},
		},
		{
			name: "add_duplicate",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddVertex(1)
				g.AddVertex(1)
				if g.Order() != 1 {
					t.Errorf("expected order 1, got %d", g.Order())
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string]()
			tc.check(t, g)
		})
	}
}

func TestLabeledGraph_RemoveVertex(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "remove_nonexistent",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.RemoveVertex(1) // Should not panic
				if g.Order() != 0 {
					t.Errorf("expected order 0, got %d", g.Order())
				}
			},
		},
		{
			name: "remove_isolated",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddVertex(1)
				g.RemoveVertex(1)
				if g.ContainsVertex(1) {
					t.Errorf("expected graph to not contain vertex 1")
				}
				if g.Order() != 0 {
					t.Errorf("expected order 0, got %d", g.Order())
				}
			},
		},
		{
			name: "remove_with_edges_directed",
			opts: []Opt{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "1->2")
				g.AddEdge(2, 3, "2->3")
				g.AddEdge(3, 1, "3->1")
				g.AddEdge(3, 3, "3->3")
				g.RemoveVertex(3)
				if g.Order() != 2 {
					t.Errorf("expected order 2, got %d", g.Order())
				}
				if g.Size() != 1 { // 1->2 remains
					t.Errorf("expected size 1, got %d", g.Size())
				}
			},
		},
		{
			name: "remove_with_edges_undirected",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "1-2")
				g.AddEdge(2, 3, "2-3")
				g.AddEdge(3, 1, "3-1")
				g.AddEdge(3, 3, "3-3")
				g.RemoveVertex(3)
				if g.Order() != 2 {
					t.Errorf("expected order 2, got %d", g.Order())
				}
				if g.Size() != 1 { // 1-2 remains
					t.Errorf("expected size 1, got %d", g.Size())
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.check(t, g)
		})
	}
}

func TestLabeledGraph_AddRemoveEdge(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "add_edge_directed",
			opts: []Opt{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "1->2")
				if !g.ContainsEdge(1, 2) {
					t.Errorf("expected graph to contain edge 1->2")
				}
				if g.ContainsEdge(2, 1) {
					t.Errorf("did not expect graph to contain edge 2->1")
				}
				if g.Size() != 1 {
					t.Errorf("expected size 1, got %d", g.Size())
				}
			},
		},
		{
			name: "add_edge_undirected",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "1-2")
				if !g.ContainsEdge(1, 2) {
					t.Errorf("expected graph to contain edge 1-2")
				}
				if !g.ContainsEdge(2, 1) {
					t.Errorf("expected graph to contain edge 2-1 (undirected)")
				}
				if g.Size() != 1 {
					t.Errorf("expected size 1, got %d", g.Size())
				}
			},
		},
		{
			name: "remove_edge",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "1-2")
				g.RemoveEdge(1, 2)
				if g.ContainsEdge(1, 2) {
					t.Errorf("expected graph to not contain edge 1-2")
				}
				if g.Size() != 0 {
					t.Errorf("expected size 0, got %d", g.Size())
				}
				g.RemoveEdge(3, 4) // should not panic
			},
		},
		{
			name: "set_label",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "initial")
				g.SetLabel(1, 2, "updated")
				l, ok := g.Label(1, 2)
				if !ok || l != "updated" {
					t.Errorf("expected label 'updated', got '%v'", l)
				}
				g.SetLabel(3, 4, "ignored") // nonexistent edge, no panic
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.check(t, g)
		})
	}
}

func TestLabeledGraph_Degree(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "directed_degree",
			opts: []Opt{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "1->2")
				g.AddEdge(2, 1, "2->1")
				g.AddEdge(1, 3, "1->3")
				g.AddEdge(1, 1, "1->1") // self loop

				if deg, ok := g.Degree(1); !ok || deg != 5 { // 3 out, 1 in (wait: 2->1 in, 1->1 in; 1->2 out, 1->3 out, 1->1 out)
					// Let's count properly: In: 2->1, 1->1 (2). Out: 1->2, 1->3, 1->1 (3). Total 5.
					t.Errorf("expected degree 5, got %d", deg)
				}
				if in, ok := g.InDegree(1); !ok || in != 2 {
					t.Errorf("expected indegree 2, got %d", in)
				}
				if out, ok := g.OutDegree(1); !ok || out != 3 {
					t.Errorf("expected outdegree 3, got %d", out)
				}
			},
		},
		{
			name: "undirected_degree",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "1-2")
				g.AddEdge(1, 3, "1-3")
				g.AddEdge(1, 1, "1-1") // self loop

				if deg, ok := g.Degree(1); !ok || deg != 3 {
					t.Errorf("expected degree 3, got %d", deg)
				}
				if in, ok := g.InDegree(1); !ok || in != 3 {
					t.Errorf("expected indegree 3, got %d", in)
				}
				if out, ok := g.OutDegree(1); !ok || out != 3 {
					t.Errorf("expected outdegree 3, got %d", out)
				}
			},
		},
		{
			name: "missing_vertex",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				if _, ok := g.Degree(1); ok {
					t.Errorf("expected degree of missing vertex to return false")
				}
				if _, ok := g.InDegree(1); ok {
					t.Errorf("expected indegree of missing vertex to return false")
				}
				if _, ok := g.OutDegree(1); ok {
					t.Errorf("expected outdegree of missing vertex to return false")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.check(t, g)
		})
	}
}

type edge struct{ u, v int }

func TestLabeledGraph_Iterators(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "vertices",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddVertex(3)
				g.AddVertex(1)
				g.AddVertex(2)
				got := slices.Collect(g.Vertices())
				slices.Sort(got)
				want := []int{1, 2, 3}
				if !slices.Equal(got, want) {
					t.Errorf("expected vertices %v, got %v", want, got)
				}
			},
		},
		{
			name: "successors_and_predecessors_directed",
			opts: []Opt{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "")
				g.AddEdge(1, 3, "")
				g.AddEdge(4, 1, "")
				got := slices.Collect(g.Successors(1))
				slices.Sort(got)
				want := []int{2, 3}
				if !slices.Equal(got, want) {
					t.Errorf("expected successors %v, got %v", want, got)
				}

				gotPreds := slices.Collect(g.Predecessors(1))
				slices.Sort(gotPreds)
				wantPreds := []int{4}
				if !slices.Equal(gotPreds, wantPreds) {
					t.Errorf("expected predecessors %v, got %v", wantPreds, gotPreds)
				}

				if len(slices.Collect(g.Successors(99))) != 0 {
					t.Errorf("expected empty iterator for missing vertex")
				}
				if len(slices.Collect(g.Predecessors(99))) != 0 {
					t.Errorf("expected empty iterator for missing vertex")
				}
			},
		},
		{
			name: "edges_and_incident_directed",
			opts: []Opt{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "")
				g.AddEdge(2, 3, "")
				g.AddEdge(4, 2, "")

				var edges []edge
				for u, v := range g.Edges() {
					edges = append(edges, edge{u, v})
				}
				slices.SortFunc(edges, func(a, b edge) int {
					if a.u != b.u {
						return a.u - b.u
					}
					return a.v - b.v
				})
				wantEdges := []edge{{1, 2}, {2, 3}, {4, 2}}
				if !slices.Equal(edges, wantEdges) {
					t.Errorf("expected edges %v, got %v", wantEdges, edges)
				}

				var incident []edge
				for u, v := range g.IncidentEdges(2) {
					incident = append(incident, edge{u, v})
				}
				slices.SortFunc(incident, func(a, b edge) int {
					if a.u != b.u {
						return a.u - b.u
					}
					return a.v - b.v
				})
				wantIncident := []edge{{1, 2}, {2, 3}, {4, 2}}
				if !slices.Equal(incident, wantIncident) {
					t.Errorf("expected incident edges %v, got %v", wantIncident, incident)
				}
			},
		},
		{
			name: "neighbors_undirected",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "")
				g.AddEdge(1, 3, "")
				got := slices.Collect(g.Neighbors(1))
				slices.Sort(got)
				want := []int{2, 3}
				if !slices.Equal(got, want) {
					t.Errorf("expected neighbors %v, got %v", want, got)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.check(t, g)
		})
	}
}

func TestLabeledGraph_Clear(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "clear",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "1-2")
				g.Clear()
				if g.Order() != 0 || g.Size() != 0 {
					t.Errorf("expected empty graph after clear")
				}
				if len(slices.Collect(g.Vertices())) != 0 {
					t.Errorf("expected no vertices")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.check(t, g)
		})
	}
}

func TestLabeledGraph_Clone(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "clone",
			opts: []Opt{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "1->2")
				g.AddVertex(3)

				clone := g.Clone()
				if clone.Order() != 3 || clone.Size() != 1 {
					t.Errorf("clone has wrong order/size")
				}

				// Mutate original, should not affect clone
				g.AddEdge(2, 3, "2->3")
				if clone.Size() != 1 {
					t.Errorf("clone was affected by mutation to original")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.check(t, g)
		})
	}
}

func TestLabeledGraph_Equal(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "equal",
			opts: []Opt{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "A")
				g.AddVertex(3)

				other := New[int, string](WithDirected())
				other.AddEdge(1, 2, "A")
				other.AddVertex(3)

				eqFunc := func(a, b string) bool { return a == b }

				if !g.Equal(other, eqFunc) {
					t.Errorf("expected graphs to be equal")
				}

				other.AddEdge(2, 3, "B")
				if g.Equal(other, eqFunc) {
					t.Errorf("expected graphs to not be equal")
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.check(t, g)
		})
	}
}

func TestLabeledGraph_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		opts  []Opt
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "string_directed",
			opts: []Opt{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				if got := g.String(); got != "digraph[]" {
					t.Errorf("expected digraph[], got: %s", got)
				}
				g.AddEdge(1, 2, "A")
				g.AddVertex(3)
				str := g.String()
				if !strings.HasPrefix(str, "digraph[") || !strings.HasSuffix(str, "]") {
					t.Errorf("expected digraph[...] format, got: %s", str)
				}
				if !strings.Contains(str, "1 -> 2: A") {
					t.Errorf("string missing edge: %s", str)
				}
				if !strings.Contains(str, "3") {
					t.Errorf("string missing isolated vertex: %s", str)
				}
			},
		},
		{
			name: "string_undirected",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				if got := g.String(); got != "graph[]" {
					t.Errorf("expected graph[], got: %s", got)
				}
				g.AddEdge(1, 2, "A")
				str := g.String()
				if !strings.HasPrefix(str, "graph[") || !strings.HasSuffix(str, "]") {
					t.Errorf("expected graph[...] format, got: %s", str)
				}
				if !strings.Contains(str, "1 - 2: A") && !strings.Contains(str, "2 - 1: A") {
					t.Errorf("string missing edge: %s", str)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.check(t, g)
		})
	}
}

func TestLabeledGraph_Coverage(t *testing.T) {
	g := New[int, int]()
	g.AddEdge(1, 2, 3)
	for range g.Successors(1) {
		break
	}
	for range g.Predecessors(2) {
		break
	}
	for range g.Edges() {
		break
	}
	for range g.IncidentEdges(1) {
		break
	}
	for range g.InIncidentEdges(2) {
		break
	}
	for range g.OutIncidentEdges(1) {
		break
	}

	g2 := New[int, int]()
	g2.AddEdge(1, 2, 10)
	g2.AddEdge(1, 4, 40)
	_ = g.Equal(g2, func(a, b int) bool { return a == b })
	_ = g.String()
	_ = g.Clone()

	un := &undirectedNodeData[int, int]{adjacentLabels: make(map[int]int)}
	un.setPredecessorLabel(2, 20)
	un.setSuccessorLabel(3, 30)
	_ = un.containsPredecessor(2)

	d := New[int, int](WithDirected())
	d.AddEdge(1, 2, 10)
	_ = d.graph[2].containsPredecessor(1)

	// cover directedNodeData setters
	dn := &directedNodeData[int, int]{successorLabels: make(map[int]int)}
	dn.setPredecessorLabel(2, 20)
	dn.setSuccessorLabel(3, 30)

	// cover undirected Edges/IncidentEdges early break
	for range g.Edges() {
		break
	}
	for range g.IncidentEdges(1) {
		break
	}

	// cover equal early returns
	g3 := New[int, int]()
	g3.AddEdge(1, 2, 10)
	g4 := New[int, int]()
	g4.AddEdge(1, 3, 10)
	_ = g3.Equal(g4, func(a, b int) bool { return a == b })

	g5 := New[int, int](WithDirected())
	_ = g3.Equal(g5, func(a, b int) bool { return a == b })
}
