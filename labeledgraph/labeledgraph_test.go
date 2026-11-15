package labeledgraph

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
				assertPanicsWithValue(t, "labeledgraph: capacity cannot be negative", func() {
					New[int, string](tc.opts...)
				})
				return
			}
			g := New[int, string](tc.opts...)
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

func TestLabeledGraph_VertexOperations(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		opts         []Option
		setup        func(*LabeledGraph[int, string])
		remove       int
		wantOrder    int
		wantSize     int
		wantPresent  []int
		wantAbsent   []int
		wantEdgeGone [][2]int
		wantEdgeKept [][2]int
	}{
		{
			name: "add new and duplicate vertices then remove missing",
			setup: func(g *LabeledGraph[int, string]) {
				g.AddVertex(1)
				g.AddVertex(1)
				g.AddVertex(2)
			},
			remove:      99,
			wantOrder:   2,
			wantSize:    0,
			wantPresent: []int{1, 2},
			wantAbsent:  []int{99},
		},
		{
			name: "undirected remove vertex with incident edges and self loop",
			setup: func(g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "a")
				g.AddEdge(2, 3, "b")
				g.AddEdge(2, 2, "self")
				g.AddEdge(3, 4, "c")
			},
			remove:       2,
			wantOrder:    3,
			wantSize:     1,
			wantPresent:  []int{1, 3, 4},
			wantAbsent:   []int{2},
			wantEdgeGone: [][2]int{{1, 2}, {2, 1}, {2, 3}, {3, 2}, {2, 2}},
			wantEdgeKept: [][2]int{{3, 4}, {4, 3}},
		},
		{
			name: "directed remove vertex with incoming outgoing and self loop",
			opts: []Option{WithDirected()},
			setup: func(g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "in")
				g.AddEdge(2, 3, "out")
				g.AddEdge(2, 2, "self")
				g.AddEdge(3, 4, "keep")
			},
			remove:       2,
			wantOrder:    3,
			wantSize:     1,
			wantPresent:  []int{1, 3, 4},
			wantAbsent:   []int{2},
			wantEdgeGone: [][2]int{{1, 2}, {2, 3}, {2, 2}},
			wantEdgeKept: [][2]int{{3, 4}},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := New[int, string](tc.opts...)
			tc.setup(g)
			g.RemoveVertex(tc.remove)

			if g.Order() != tc.wantOrder {
				t.Errorf("Order() = %d, want %d", g.Order(), tc.wantOrder)
			}
			if g.Size() != tc.wantSize {
				t.Errorf("Size() = %d, want %d", g.Size(), tc.wantSize)
			}
			for _, v := range tc.wantPresent {
				if !g.ContainsVertex(v) {
					t.Errorf("expected vertex %d to be present", v)
				}
			}
			for _, v := range tc.wantAbsent {
				if g.ContainsVertex(v) {
					t.Errorf("expected vertex %d to be absent", v)
				}
			}
			for _, e := range tc.wantEdgeGone {
				if g.ContainsEdge(e[0], e[1]) {
					t.Errorf("expected edge (%d, %d) to be absent", e[0], e[1])
				}
			}
			for _, e := range tc.wantEdgeKept {
				if !g.ContainsEdge(e[0], e[1]) {
					t.Errorf("expected edge (%d, %d) to be present", e[0], e[1])
				}
			}
		})
	}
}

func TestLabeledGraph_EdgeAndLabelOperations(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "undirected add overwrite setLabel and removeEdge",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddVertex(1)
				g.AddEdge(1, 2, "initial")
				if g.Order() != 2 || g.Size() != 1 {
					t.Fatalf("Order()=%d Size()=%d, want 2, 1", g.Order(), g.Size())
				}
				if !g.ContainsEdge(1, 2) || !g.ContainsEdge(2, 1) {
					t.Error("expected undirected edge (1, 2) and (2, 1)")
				}

				// Overwrite existing edge label without incrementing Size.
				g.AddEdge(2, 1, "overwritten")
				if g.Size() != 1 {
					t.Errorf("Size() after overwrite = %d, want 1", g.Size())
				}
				if lbl, ok := g.Label(1, 2); !ok || lbl != "overwritten" {
					t.Errorf("Label(1, 2) = (%q, %v), want (\"overwritten\", true)", lbl, ok)
				}

				// SetLabel updates both directions in undirected graph.
				g.SetLabel(1, 2, "updated")
				if lblRev, ok := g.Label(2, 1); !ok || lblRev != "updated" {
					t.Errorf("Label(2, 1) = (%q, %v), want (\"updated\", true)", lblRev, ok)
				}

				// SetLabel and Label on missing vertex / missing edge.
				g.SetLabel(99, 1, "noop")
				g.SetLabel(1, 99, "noop")
				if _, ok := g.Label(99, 1); ok {
					t.Error("expected Label(99, 1) ok=false")
				}
				if _, ok := g.Label(1, 99); ok {
					t.Error("expected Label(1, 99) ok=false")
				}
				if g.ContainsEdge(99, 1) || g.ContainsEdge(1, 99) {
					t.Error("expected missing edges to return false")
				}

				// RemoveEdge on missing vertex / missing edge is a no-op.
				g.RemoveEdge(99, 1)
				g.RemoveEdge(1, 99)
				if g.Size() != 1 {
					t.Errorf("Size() = %d, want 1", g.Size())
				}

				g.RemoveEdge(1, 2)
				if g.Size() != 0 || g.ContainsEdge(1, 2) || g.ContainsEdge(2, 1) {
					t.Error("expected edge (1, 2) removed")
				}
			},
		},
		{
			name: "directed add overwrite setLabel removeEdge and clear",
			opts: []Option{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "forward")
				if g.Order() != 2 || g.Size() != 1 {
					t.Fatalf("Order()=%d Size()=%d, want 2, 1", g.Order(), g.Size())
				}
				if !g.ContainsEdge(1, 2) || g.ContainsEdge(2, 1) {
					t.Error("expected directed edge (1, 2) only")
				}

				// Overwrite existing directed edge.
				g.AddEdge(1, 2, "forward2")
				if g.Size() != 1 {
					t.Errorf("Size() after overwrite = %d, want 1", g.Size())
				}

				// SetLabel only updates existing directed direction.
				g.SetLabel(1, 2, "forward3")
				g.SetLabel(2, 1, "ignored")
				if lbl, ok := g.Label(1, 2); !ok || lbl != "forward3" {
					t.Errorf("Label(1, 2) = (%q, %v), want (\"forward3\", true)", lbl, ok)
				}
				if _, ok := g.Label(2, 1); ok {
					t.Error("expected Label(2, 1) ok=false")
				}

				// Reverse edge is distinct in directed graph.
				g.AddEdge(2, 1, "backward")
				if g.Size() != 2 {
					t.Errorf("Size() = %d, want 2", g.Size())
				}
				g.RemoveEdge(1, 2)
				if g.Size() != 1 || g.ContainsEdge(1, 2) || !g.ContainsEdge(2, 1) {
					t.Error("expected only (2, 1) remaining")
				}

				g.Clear()
				if g.Order() != 0 || g.Size() != 0 {
					t.Errorf("after Clear: Order()=%d Size()=%d, want 0, 0", g.Order(), g.Size())
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

func TestLabeledGraph_Degrees(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		opts       []Option
		setup      func(*LabeledGraph[int, string])
		vertex     int
		wantOK     bool
		wantDeg    int
		wantInDeg  int
		wantOutDeg int
	}{
		{
			name:   "missing vertex undirected",
			setup:  func(_ *LabeledGraph[int, string]) {},
			vertex: 42,
			wantOK: false,
		},
		{
			name:   "missing vertex directed",
			opts:   []Option{WithDirected()},
			setup:  func(_ *LabeledGraph[int, string]) {},
			vertex: 42,
			wantOK: false,
		},
		{
			name: "undirected vertex degrees are synonymous",
			setup: func(g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "a")
				g.AddEdge(1, 3, "b")
			},
			vertex:     1,
			wantOK:     true,
			wantDeg:    2,
			wantInDeg:  2,
			wantOutDeg: 2,
		},
		{
			name: "directed vertex distinct in and out degrees",
			opts: []Option{WithDirected()},
			setup: func(g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "a")
				g.AddEdge(3, 2, "b")
				g.AddEdge(2, 4, "c")
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
			g := New[int, string](tc.opts...)
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

func TestLabeledGraph_CloneAndEqual(t *testing.T) {
	t.Parallel()

	strEq := func(a, b string) bool { return a == b }

	cases := []struct {
		name  string
		build func() (g1, g2 *LabeledGraph[int, string])
		want  bool
	}{
		{
			name: "both nil",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				return nil, nil
			},
			want: true,
		},
		{
			name: "nil receiver vs non-nil",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				return nil, New[int, string]()
			},
			want: false,
		},
		{
			name: "non-nil receiver vs nil",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				return New[int, string](), nil
			},
			want: false,
		},
		{
			name: "directed vs undirected",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				return New[int, string](WithDirected()), New[int, string]()
			},
			want: false,
		},
		{
			name: "different vertex order",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				g1 := New[int, string]()
				g1.AddVertex(1)
				return g1, New[int, string]()
			},
			want: false,
		},
		{
			name: "different edge count",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				g1 := New[int, string]()
				g1.AddEdge(1, 2, "a")
				g2 := New[int, string]()
				g2.AddVertex(1)
				g2.AddVertex(2)
				return g1, g2
			},
			want: false,
		},
		{
			name: "different vertex sets with same order and size",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				g1 := New[int, string]()
				g1.AddVertex(1)
				g2 := New[int, string]()
				g2.AddVertex(2)
				return g1, g2
			},
			want: false,
		},
		{
			name: "different out-degree distribution with same vertices and size",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				g1 := New[int, string](WithDirected())
				g1.AddEdge(1, 2, "a")
				g2 := New[int, string](WithDirected())
				g2.AddEdge(2, 1, "a")
				return g1, g2
			},
			want: false,
		},
		{
			name: "different target edges with identical per-vertex out-degree",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				g1 := New[int, string](WithDirected())
				g1.AddEdge(1, 2, "a")
				g1.AddEdge(2, 3, "a")
				g1.AddEdge(3, 1, "a")
				g2 := New[int, string](WithDirected())
				g2.AddEdge(1, 3, "a")
				g2.AddEdge(3, 2, "a")
				g2.AddEdge(2, 1, "a")
				return g1, g2
			},
			want: false,
		},
		{
			name: "different edge labels",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				g1 := New[int, string](WithDirected())
				g1.AddEdge(1, 2, "a")
				g2 := New[int, string](WithDirected())
				g2.AddEdge(1, 2, "b")
				return g1, g2
			},
			want: false,
		},
		{
			name: "cloned directed graph is equal and deep-isolated",
			build: func() (*LabeledGraph[int, string], *LabeledGraph[int, string]) {
				g1 := New[int, string](WithDirected())
				g1.AddEdge(1, 2, "a")
				g1.AddEdge(2, 3, "b")
				g2 := g1.Clone()

				// Mutating a second clone must not affect g1 or g2.
				g3 := g2.Clone()
				g3.SetLabel(1, 2, "mutated")
				g3.AddEdge(3, 1, "extra")
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
			if got := g1.Equal(g2, strEq); got != tc.want {
				t.Errorf("Equal() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLabeledGraph_Iterators(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		opts  []Option
		check func(*testing.T, *LabeledGraph[int, string])
	}{
		{
			name: "undirected iterators full and early termination",
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "a")
				g.AddEdge(1, 3, "b")

				// Missing vertex iterators yield nothing.
				if got := slices.Collect(g.Neighbors(99)); len(got) != 0 {
					t.Errorf("Neighbors(99) = %v, want empty", got)
				}
				if got := slices.Collect(g.Successors(99)); len(got) != 0 {
					t.Errorf("Successors(99) = %v, want empty", got)
				}
				if got := slices.Collect(g.Predecessors(99)); len(got) != 0 {
					t.Errorf("Predecessors(99) = %v, want empty", got)
				}
				if got := collectEdges(g.IncidentEdges(99)); len(got) != 0 {
					t.Errorf("IncidentEdges(99) = %v, want empty", got)
				}
				if got := collectEdges(g.InIncidentEdges(99)); len(got) != 0 {
					t.Errorf("InIncidentEdges(99) = %v, want empty", got)
				}
				if got := collectEdges(g.OutIncidentEdges(99)); len(got) != 0 {
					t.Errorf("OutIncidentEdges(99) = %v, want empty", got)
				}

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

				// Early termination across all undirected iterators.
				count := 0
				for range g.Vertices() {
					count++
					break
				}
				if count != 1 {
					t.Errorf("Vertices early break count = %d, want 1", count)
				}

				count = 0
				for range g.Successors(1) {
					count++
					break
				}
				if count != 1 {
					t.Errorf("Successors early break count = %d, want 1", count)
				}

				count = 0
				for range g.Predecessors(1) {
					count++
					break
				}
				if count != 1 {
					t.Errorf("Predecessors early break count = %d, want 1", count)
				}

				count = 0
				for range g.Edges() {
					count++
					break
				}
				if count != 1 {
					t.Errorf("Edges early break count = %d, want 1", count)
				}

				count = 0
				for range g.IncidentEdges(1) {
					count++
					break
				}
				if count != 1 {
					t.Errorf("IncidentEdges early break count = %d, want 1", count)
				}

				count = 0
				for range g.InIncidentEdges(1) {
					count++
					break
				}
				if count != 1 {
					t.Errorf("InIncidentEdges early break count = %d, want 1", count)
				}

				count = 0
				for range g.OutIncidentEdges(1) {
					count++
					break
				}
				if count != 1 {
					t.Errorf("OutIncidentEdges early break count = %d, want 1", count)
				}
			},
		},
		{
			name: "directed iterators self-loop deduplication and early termination",
			opts: []Option{WithDirected()},
			check: func(t *testing.T, g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "in1")
				g.AddEdge(3, 2, "in2")
				g.AddEdge(2, 2, "self")
				g.AddEdge(2, 4, "out1")
				g.AddEdge(2, 5, "out2")

				preds := slices.Collect(g.Predecessors(2))
				slices.Sort(preds)
				if !slices.Equal(preds, []int{1, 2, 3}) {
					t.Errorf("Predecessors(2) = %v, want [1 2 3]", preds)
				}

				succs := slices.Collect(g.Successors(2))
				slices.Sort(succs)
				if !slices.Equal(succs, []int{2, 4, 5}) {
					t.Errorf("Successors(2) = %v, want [2 4 5]", succs)
				}

				// Self-loop (2, 2) must be yielded exactly once by IncidentEdges(2).
				if got, want := collectEdges(g.IncidentEdges(2)), []edgePair{{1, 2}, {2, 2}, {2, 4}, {2, 5}, {3, 2}}; !slices.Equal(got, want) {
					t.Errorf("IncidentEdges(2) = %v, want %v", got, want)
				}
				if got, want := collectEdges(g.InIncidentEdges(2)), []edgePair{{1, 2}, {2, 2}, {3, 2}}; !slices.Equal(got, want) {
					t.Errorf("InIncidentEdges(2) = %v, want %v", got, want)
				}
				if got, want := collectEdges(g.OutIncidentEdges(2)), []edgePair{{2, 2}, {2, 4}, {2, 5}}; !slices.Equal(got, want) {
					t.Errorf("OutIncidentEdges(2) = %v, want %v", got, want)
				}

				// Early termination in directed Predecessors, InIncidentEdges, IncidentEdges (in loop), and IncidentEdges (out loop).
				count := 0
				for range g.Predecessors(2) {
					count++
					break
				}
				if count != 1 {
					t.Errorf("Predecessors early break count = %d, want 1", count)
				}

				count = 0
				for range g.InIncidentEdges(2) {
					count++
					break
				}
				if count != 1 {
					t.Errorf("InIncidentEdges early break count = %d, want 1", count)
				}

				count = 0
				for range g.IncidentEdges(2) {
					count++
					break
				}
				if count != 1 {
					t.Errorf("IncidentEdges(2) early break count = %d, want 1", count)
				}

				// Vertex 1 has only outgoing edges; add a second outgoing edge so early break inside ev.out loop is exercised.
				g.AddEdge(1, 6, "outFrom1")
				count = 0
				for range g.IncidentEdges(1) {
					count++
					break
				}
				if count != 1 {
					t.Errorf("IncidentEdges(1) early break count = %d, want 1", count)
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
		name       string
		opts       []Option
		setup      func(*LabeledGraph[int, string])
		wantExact  string
		wantPrefix string
		wantParts  []string
	}{
		{
			name:      "empty undirected",
			setup:     func(_ *LabeledGraph[int, string]) {},
			wantExact: "graph[]",
		},
		{
			name:      "empty directed",
			opts:      []Option{WithDirected()},
			setup:     func(_ *LabeledGraph[int, string]) {},
			wantExact: "digraph[]",
		},
		{
			name: "directed multiple edges and isolated vertices",
			opts: []Option{WithDirected()},
			setup: func(g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "a")
				g.AddEdge(2, 3, "b")
				g.AddVertex(10)
				g.AddVertex(20)
			},
			wantPrefix: "digraph[",
			wantParts:  []string{"1 -> 2: a", "2 -> 3: b", "10", "20"},
		},
		{
			name: "directed only isolated vertices",
			opts: []Option{WithDirected()},
			setup: func(g *LabeledGraph[int, string]) {
				g.AddVertex(10)
				g.AddVertex(20)
			},
			wantPrefix: "digraph[",
			wantParts:  []string{"10", "20"},
		},
		{
			name: "undirected multiple edges and isolated vertices",
			setup: func(g *LabeledGraph[int, string]) {
				g.AddEdge(1, 2, "a")
				g.AddEdge(2, 3, "b")
				g.AddVertex(10)
				g.AddVertex(20)
			},
			wantPrefix: "graph[",
			wantParts:  []string{"10", "20"},
		},
		{
			name: "undirected only isolated vertices",
			setup: func(g *LabeledGraph[int, string]) {
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
			g := New[int, string](tc.opts...)
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
