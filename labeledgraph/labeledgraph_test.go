package labeledgraph

import (
	"testing"
)

func TestLabeledGraph_Directed(t *testing.T) {
	g := New[int, string](Directed())

	g.AddEdge(1, 2, "1->2")
	g.AddEdge(2, 3, "2->3")
	g.AddEdge(3, 1, "3->1")
	g.AddEdge(3, 3, "3->3") // self loop

	if g.Order() != 3 {
		t.Errorf("expected order 3, got %d", g.Order())
	}
	if g.Size() != 4 {
		t.Errorf("expected size 4, got %d", g.Size())
	}

	g.RemoveVertex(3)

	if g.Order() != 2 {
		t.Errorf("expected order 2, got %d", g.Order())
	}
	// Removing 3 should remove 2->3, 3->1, and 3->3. Remaining is 1->2.
	if g.Size() != 1 {
		t.Errorf("expected size 1, got %d", g.Size())
	}

	g.AddEdge(1, 2, "new-label")
	if g.Size() != 1 {
		t.Errorf("expected size 1 after label update, got %d", g.Size())
	}
}

func TestLabeledGraph_Undirected(t *testing.T) {
	g := New[int, string]()

	g.AddEdge(1, 2, "1-2")
	g.AddEdge(2, 3, "2-3")
	g.AddEdge(3, 1, "3-1")
	g.AddEdge(3, 3, "3-3") // self loop

	if g.Order() != 3 {
		t.Errorf("expected order 3, got %d", g.Order())
	}
	if g.Size() != 4 {
		t.Errorf("expected size 4, got %d", g.Size())
	}

	g.RemoveVertex(3)

	if g.Order() != 2 {
		t.Errorf("expected order 2, got %d", g.Order())
	}
	// Removing 3 should remove 2-3, 3-1, 3-3. Remaining: 1-2
	if g.Size() != 1 {
		t.Errorf("expected size 1, got %d", g.Size())
	}
	
	// Removing 1-2 edge
	g.RemoveEdge(1, 2)
	if g.Size() != 0 {
		t.Errorf("expected size 0, got %d", g.Size())
	}
}

func TestLabeledGraph_IncidentEdges(t *testing.T) {
	g := New[int, string](Directed())
	g.AddEdge(1, 2, "1->2")
	count := 0
	for u, v := range g.IncidentEdges(1) {
		if u == 1 && v == 2 {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 incident edge for 1, got %d", count)
	}
}
