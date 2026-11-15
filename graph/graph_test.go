package graph

import (
	"slices"
	"strings"
	"testing"
)

type edgePair struct {
	u, v int
}

func collectEdges(seq func(func(int, int) bool)) []edgePair {
	var out []edgePair
	seq(func(u, v int) bool {
		out = append(out, edgePair{u: u, v: v})
		return true
	})
	slices.SortFunc(out, func(a, b edgePair) int {
		if a.u != b.u {
			return a.u - b.u
		}
		return a.v - b.v
	})
	return out
}

func assertPanicsWithValue(t *testing.T, want string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic %q, got none", want)
		}
		if r != want {
			t.Fatalf("expected panic %q, got %v", want, r)
		}
	}()
	fn()
}

func TestNew(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		opts         []Option
		wantDirected bool
		wantPanic    bool
	}{
		{
			name:         "default undirected",
			opts:         nil,
			wantDirected: false,
		},
		{
			name:         "directed",
			opts:         []Option{WithDirected()},
			wantDirected: true,
		},
		{
			name:         "zero capacity",
			opts:         []Option{WithCapacity(0)},
			wantDirected: false,
		},
		{
			name:         "positive capacity and directed",
			opts:         []Option{WithDirected(), WithCapacity(64)},
			wantDirected: true,
		},
		{
			name:      "negative capacity panics",
			opts:      []Option{WithCapacity(-1)},
			wantPanic: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.wantPanic {
				assertPanicsWithValue(t, "graph: capacity cannot be negative", func() {
					New[int](tc.opts...)
				})
				return
			}
			g := New[int](tc.opts...)
			if g == nil {
				t.Fatal("expected non-nil graph")
			}
			if g.Directed() != tc.wantDirected {
				t.Errorf("Directed() = %v, want %v", g.Directed(), tc.wantDirected)
			}
			if g.Order() != 0 || g.Size() != 0 {
				t.Errorf("Order()=%d Size()=%d, want 0, 0", g.Order(), g.Size())
			}
		})
	}
}

func TestGraph_MutationsAndQueries(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *Graph[int])
	}{
		{
			name: "undirected vertex and edge mutations",
			check: func(t *testing.T, g *Graph[int]) {
				g.AddVertex(1)
				g.AddVertex(1)
				if g.Order() != 1 || !g.ContainsVertex(1) || g.ContainsVertex(99) {
					t.Fatalf("unexpected vertex state: Order()=%d", g.Order())
				}

				g.AddEdge(1, 2)
				g.AddEdge(2, 3)
				if g.Order() != 3 || g.Size() != 2 {
					t.Fatalf("Order()=%d Size()=%d, want 3, 2", g.Order(), g.Size())
				}
				if !g.ContainsEdge(1, 2) || !g.ContainsEdge(2, 1) || g.ContainsEdge(1, 3) {
					t.Error("unexpected ContainsEdge results")
				}

				g.RemoveEdge(1, 3) // missing edge no-op
				if g.Size() != 2 {
					t.Errorf("Size() = %d, want 2", g.Size())
				}

				g.RemoveEdge(2, 1)
				if g.Size() != 1 || g.ContainsEdge(1, 2) {
					t.Error("expected edge (1, 2) removed")
				}

				g.RemoveVertex(2)
				if g.Order() != 2 || g.Size() != 0 || g.ContainsVertex(2) {
					t.Errorf("after RemoveVertex(2): Order()=%d Size()=%d", g.Order(), g.Size())
				}

				g.Clear()
				if g.Order() != 0 || g.Size() != 0 {
					t.Errorf("after Clear: Order()=%d Size()=%d, want 0, 0", g.Order(), g.Size())
				}
			},
		},
		{
			name: "directed vertex and edge mutations",
			opts: []Option{WithDirected()},
			check: func(t *testing.T, g *Graph[int]) {
				g.AddEdge(1, 2)
				g.AddEdge(2, 3)
				g.AddEdge(2, 2)
				if g.Order() != 3 || g.Size() != 3 {
					t.Fatalf("Order()=%d Size()=%d, want 3, 3", g.Order(), g.Size())
				}
				if !g.ContainsEdge(1, 2) || g.ContainsEdge(2, 1) {
					t.Error("expected directed edge (1, 2) only")
				}

				g.RemoveVertex(2)
				if g.Order() != 2 || g.Size() != 0 {
					t.Errorf("after RemoveVertex(2): Order()=%d Size()=%d, want 2, 0", g.Order(), g.Size())
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

func TestGraph_Degrees(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		opts       []Option
		setup      func(*Graph[int])
		vertex     int
		wantOK     bool
		wantDeg    int
		wantInDeg  int
		wantOutDeg int
	}{
		{
			name:   "missing vertex",
			setup:  func(_ *Graph[int]) {},
			vertex: 99,
			wantOK: false,
		},
		{
			name: "undirected vertex",
			setup: func(g *Graph[int]) {
				g.AddEdge(1, 2)
				g.AddEdge(1, 3)
			},
			vertex:     1,
			wantOK:     true,
			wantDeg:    2,
			wantInDeg:  2,
			wantOutDeg: 2,
		},
		{
			name: "directed vertex",
			opts: []Option{WithDirected()},
			setup: func(g *Graph[int]) {
				g.AddEdge(1, 2)
				g.AddEdge(3, 2)
				g.AddEdge(2, 4)
			},
			vertex:     2,
			wantOK:     true,
			wantDeg:    3,
			wantInDeg:  2,
			wantOutDeg: 1,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int](tc.opts...)
			tc.setup(g)

			if deg, ok := g.Degree(tc.vertex); ok != tc.wantOK || deg != tc.wantDeg {
				t.Errorf("Degree(%d) = (%d, %v), want (%d, %v)", tc.vertex, deg, ok, tc.wantDeg, tc.wantOK)
			}
			if inDeg, ok := g.InDegree(tc.vertex); ok != tc.wantOK || inDeg != tc.wantInDeg {
				t.Errorf("InDegree(%d) = (%d, %v), want (%d, %v)", tc.vertex, inDeg, ok, tc.wantInDeg, tc.wantOK)
			}
			if outDeg, ok := g.OutDegree(tc.vertex); ok != tc.wantOK || outDeg != tc.wantOutDeg {
				t.Errorf("OutDegree(%d) = (%d, %v), want (%d, %v)", tc.vertex, outDeg, ok, tc.wantOutDeg, tc.wantOK)
			}
		})
	}
}

func TestGraph_CloneAndEqual(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		build func() (*Graph[int], *Graph[int])
		want  bool
	}{
		{
			name: "both nil",
			build: func() (*Graph[int], *Graph[int]) {
				return nil, nil
			},
			want: true,
		},
		{
			name: "nil vs non-nil",
			build: func() (*Graph[int], *Graph[int]) {
				return nil, New[int]()
			},
			want: false,
		},
		{
			name: "non-nil vs nil",
			build: func() (*Graph[int], *Graph[int]) {
				return New[int](), nil
			},
			want: false,
		},
		{
			name: "different edges",
			build: func() (*Graph[int], *Graph[int]) {
				g1 := New[int](WithDirected())
				g1.AddEdge(1, 2)
				g2 := New[int](WithDirected())
				g2.AddEdge(2, 1)
				return g1, g2
			},
			want: false,
		},
		{
			name: "cloned graph with edges is equal and isolated",
			build: func() (*Graph[int], *Graph[int]) {
				g1 := New[int](WithDirected())
				g1.AddEdge(1, 2)
				g1.AddEdge(2, 3)
				g2 := g1.Clone()

				g3 := g2.Clone()
				g3.AddEdge(3, 1)
				return g1, g2
			},
			want: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g1, g2 := tc.build()
			if got := g1.Equal(g2); got != tc.want {
				t.Errorf("Equal() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGraph_Iterators(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *Graph[int])
	}{
		{
			name: "undirected iterators",
			check: func(t *testing.T, g *Graph[int]) {
				g.AddEdge(1, 2)
				g.AddEdge(1, 3)

				verts := slices.Collect(g.Vertices())
				slices.Sort(verts)
				if !slices.Equal(verts, []int{1, 2, 3}) {
					t.Errorf("Vertices() = %v, want [1 2 3]", verts)
				}

				neighbors := slices.Collect(g.Neighbors(1))
				slices.Sort(neighbors)
				if !slices.Equal(neighbors, []int{2, 3}) {
					t.Errorf("Neighbors(1) = %v, want [2 3]", neighbors)
				}

				succs := slices.Collect(g.Successors(1))
				slices.Sort(succs)
				if !slices.Equal(succs, []int{2, 3}) {
					t.Errorf("Successors(1) = %v, want [2 3]", succs)
				}

				preds := slices.Collect(g.Predecessors(1))
				slices.Sort(preds)
				if !slices.Equal(preds, []int{2, 3}) {
					t.Errorf("Predecessors(1) = %v, want [2 3]", preds)
				}

				if got, want := collectEdges(g.Edges()), []edgePair{{1, 2}, {1, 3}, {2, 1}, {3, 1}}; !slices.Equal(got, want) {
					t.Errorf("Edges() = %v, want %v", got, want)
				}
				if got, want := collectEdges(g.IncidentEdges(1)), []edgePair{{1, 2}, {1, 3}}; !slices.Equal(got, want) {
					t.Errorf("IncidentEdges(1) = %v, want %v", got, want)
				}
				if got, want := collectEdges(g.InIncidentEdges(1)), []edgePair{{2, 1}, {3, 1}}; !slices.Equal(got, want) {
					t.Errorf("InIncidentEdges(1) = %v, want %v", got, want)
				}
				if got, want := collectEdges(g.OutIncidentEdges(1)), []edgePair{{1, 2}, {1, 3}}; !slices.Equal(got, want) {
					t.Errorf("OutIncidentEdges(1) = %v, want %v", got, want)
				}
			},
		},
		{
			name: "directed iterators including self-loop deduplication",
			opts: []Option{WithDirected()},
			check: func(t *testing.T, g *Graph[int]) {
				g.AddEdge(1, 2)
				g.AddEdge(2, 2)
				g.AddEdge(2, 3)

				if got, want := collectEdges(g.IncidentEdges(2)), []edgePair{{1, 2}, {2, 2}, {2, 3}}; !slices.Equal(got, want) {
					t.Errorf("IncidentEdges(2) = %v, want %v", got, want)
				}
				if got, want := collectEdges(g.InIncidentEdges(2)), []edgePair{{1, 2}, {2, 2}}; !slices.Equal(got, want) {
					t.Errorf("InIncidentEdges(2) = %v, want %v", got, want)
				}
				if got, want := collectEdges(g.OutIncidentEdges(2)), []edgePair{{2, 2}, {2, 3}}; !slices.Equal(got, want) {
					t.Errorf("OutIncidentEdges(2) = %v, want %v", got, want)
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

func TestGraph_String(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		opts       []Option
		setup      func(*Graph[int])
		wantExact  string
		wantPrefix string
		wantParts  []string
	}{
		{
			name:      "empty undirected",
			setup:     func(_ *Graph[int]) {},
			wantExact: "graph[]",
		},
		{
			name:      "empty directed",
			opts:      []Option{WithDirected()},
			setup:     func(_ *Graph[int]) {},
			wantExact: "digraph[]",
		},
		{
			name: "directed multiple edges and isolated vertices",
			opts: []Option{WithDirected()},
			setup: func(g *Graph[int]) {
				g.AddEdge(1, 2)
				g.AddEdge(2, 3)
				g.AddVertex(10)
				g.AddVertex(20)
			},
			wantPrefix: "digraph[",
			wantParts:  []string{"1 -> 2", "2 -> 3", "10", "20"},
		},
		{
			name: "directed only isolated vertices",
			opts: []Option{WithDirected()},
			setup: func(g *Graph[int]) {
				g.AddVertex(10)
				g.AddVertex(20)
			},
			wantPrefix: "digraph[",
			wantParts:  []string{"10", "20"},
		},
		{
			name: "undirected multiple edges and isolated vertices",
			setup: func(g *Graph[int]) {
				g.AddEdge(1, 2)
				g.AddEdge(2, 3)
				g.AddVertex(10)
				g.AddVertex(20)
			},
			wantPrefix: "graph[",
			wantParts:  []string{"10", "20"},
		},
		{
			name: "undirected only isolated vertices",
			setup: func(g *Graph[int]) {
				g.AddVertex(10)
				g.AddVertex(20)
			},
			wantPrefix: "graph[",
			wantParts:  []string{"10", "20"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int](tc.opts...)
			tc.setup(g)

			got := g.String()
			if tc.wantExact != "" {
				if got != tc.wantExact {
					t.Errorf("String() = %q, want %q", got, tc.wantExact)
				}
				return
			}
			if !strings.HasPrefix(got, tc.wantPrefix) || !strings.HasSuffix(got, "]") {
				t.Fatalf("String() = %q, want prefix %q and suffix \"]\"", got, tc.wantPrefix)
			}
			inner := strings.TrimSuffix(strings.TrimPrefix(got, tc.wantPrefix), "]")
			tokens := strings.Split(inner, ", ")
			for _, part := range tc.wantParts {
				if !slices.Contains(tokens, part) {
					t.Errorf("String() tokens %v missing %q", tokens, part)
				}
			}
		})
	}
}
