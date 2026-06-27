package graph

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
		check func(*testing.T, *Graph[int])
	}{
		{
			name: "default_undirected",
			check: func(t *testing.T, g *Graph[int]) {
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
			check: func(t *testing.T, g *Graph[int]) {
				if !g.Directed() {
					t.Errorf("expected directed graph")
				}
			},
		},
		{
			name: "capacity",
			opts: []Opt{WithCapacity(100)},
			check: func(t *testing.T, g *Graph[int]) {
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
			g := New[int](tc.opts...)
			tc.check(t, g)
		})
	}
}

func TestGraph_AddRemoveVertex(t *testing.T) {
	t.Parallel()
	g := New[int]()
	g.AddVertex(1)
	if !g.ContainsVertex(1) {
		t.Errorf("expected to contain vertex 1")
	}
	g.RemoveVertex(1)
	if g.ContainsVertex(1) {
		t.Errorf("expected not to contain vertex 1")
	}
}

func TestGraph_AddRemoveEdge(t *testing.T) {
	t.Parallel()
	g := New[int]()
	g.AddEdge(1, 2)
	if !g.ContainsEdge(1, 2) {
		t.Errorf("expected to contain edge")
	}
	if g.Size() != 1 {
		t.Errorf("expected size 1, got %d", g.Size())
	}
	g.RemoveEdge(1, 2)
	if g.ContainsEdge(1, 2) {
		t.Errorf("expected not to contain edge")
	}
}

func TestGraph_Degree(t *testing.T) {
	t.Parallel()
	g := New[int](WithDirected())
	g.AddEdge(1, 2)
	g.AddEdge(2, 1)
	g.AddEdge(1, 3)

	if deg, ok := g.Degree(1); !ok || deg != 3 {
		t.Errorf("expected degree 3, got %d", deg)
	}
	if in, ok := g.InDegree(1); !ok || in != 1 {
		t.Errorf("expected in-degree 1, got %d", in)
	}
	if out, ok := g.OutDegree(1); !ok || out != 2 {
		t.Errorf("expected out-degree 2, got %d", out)
	}
}

type edge struct{ u, v int }

func TestGraph_Iterators(t *testing.T) {
	t.Parallel()
	g := New[int](WithDirected())
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	verts := slices.Collect(g.Vertices())
	slices.Sort(verts)
	if !slices.Equal(verts, []int{1, 2, 3}) {
		t.Errorf("unexpected vertices")
	}

	succs := slices.Collect(g.Successors(2))
	if !slices.Equal(succs, []int{3}) {
		t.Errorf("unexpected successors")
	}

	preds := slices.Collect(g.Predecessors(2))
	if !slices.Equal(preds, []int{1}) {
		t.Errorf("unexpected predecessors")
	}

	neighs := slices.Collect(g.Neighbors(2))
	if !slices.Equal(neighs, []int{3}) {
		t.Errorf("unexpected neighbors")
	}

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
	if !slices.Equal(edges, []edge{{1, 2}, {2, 3}}) {
		t.Errorf("unexpected edges")
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
	if !slices.Equal(incident, []edge{{1, 2}, {2, 3}}) {
		t.Errorf("unexpected incident edges")
	}

	var inIncident []edge
	for u, v := range g.InIncidentEdges(2) {
		inIncident = append(inIncident, edge{u, v})
	}
	if !slices.Equal(inIncident, []edge{{1, 2}}) {
		t.Errorf("unexpected in-incident edges")
	}

	var outIncident []edge
	for u, v := range g.OutIncidentEdges(2) {
		outIncident = append(outIncident, edge{u, v})
	}
	if !slices.Equal(outIncident, []edge{{2, 3}}) {
		t.Errorf("unexpected out-incident edges")
	}
}

func TestGraph_Clear(t *testing.T) {
	t.Parallel()
	g := New[int]()
	g.AddEdge(1, 2)
	g.Clear()
	if g.Order() != 0 || g.Size() != 0 {
		t.Errorf("expected empty graph")
	}
}

func TestGraph_Clone(t *testing.T) {
	t.Parallel()
	g := New[int]()
	g.AddEdge(1, 2)
	clone := g.Clone()
	if clone.Order() != 2 || clone.Size() != 1 {
		t.Errorf("clone has wrong order/size")
	}
	g.AddEdge(2, 3)
	if clone.Size() != 1 {
		t.Errorf("clone affected by original mutation")
	}
}

func TestGraph_Equal(t *testing.T) {
	t.Parallel()
	g1 := New[int]()
	g1.AddEdge(1, 2)

	g2 := New[int]()
	g2.AddEdge(1, 2)

	if !g1.Equal(g2) {
		t.Errorf("expected equal")
	}

	g2.AddEdge(2, 3)
	if g1.Equal(g2) {
		t.Errorf("expected not equal")
	}
}

func TestGraph_String(t *testing.T) {
	t.Parallel()

	g1 := New[int](WithDirected())
	g1.AddEdge(1, 2)
	g1.AddVertex(3)
	str1 := g1.String()
	if !strings.Contains(str1, "1 -> 2") {
		t.Errorf("missing edge: %s", str1)
	}
	if strings.Contains(str1, "{}") {
		t.Errorf("unexpected void struct in string: %s", str1)
	}

	g2 := New[int]()
	g2.AddEdge(1, 2)
	str2 := g2.String()
	if !strings.Contains(str2, "1 - 2") && !strings.Contains(str2, "2 - 1") {
		t.Errorf("missing edge: %s", str2)
	}
}

func TestGraph_Coverage(t *testing.T) {
	g := New[int]()
	g.AddEdge(1, 2)
	_ = g.String()
}
