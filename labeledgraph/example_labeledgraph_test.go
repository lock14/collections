package labeledgraph_test

import (
	"fmt"
	"github.com/lock14/collections/labeledgraph"
	"slices"
)

func ExampleLabeledGraph() {
	// Create a directed graph with string vertices and string edge labels (relationships)
	g := labeledgraph.New[string, string](labeledgraph.Directed())

	g.AddEdge("Alice", "Bob", "Knows")
	g.AddEdge("Bob", "Charlie", "Likes")
	g.AddEdge("Alice", "Charlie", "Dislikes")

	fmt.Println("Alice's outgoing relationships:")
	succs := slices.Collect(g.Successors("Alice"))
	slices.Sort(succs)

	for _, target := range succs {
		relationship, _ := g.Label("Alice", target)
		fmt.Printf("- Alice %s %s\n", relationship, target)
	}

	// Output:
	// Alice's outgoing relationships:
	// - Alice Knows Bob
	// - Alice Dislikes Charlie
}

func ExampleLabeledGraph_ContainsVertex() {
	g := labeledgraph.New[int, string]()
	g.AddVertex(1)
	fmt.Println(g.ContainsVertex(1))
	fmt.Println(g.ContainsVertex(2))
	// Output:
	// true
	// false
}

func ExampleLabeledGraph_AddVertex() {
	g := labeledgraph.New[int, string]()
	g.AddVertex(1)
	fmt.Println(g.ContainsVertex(1))
	// Output:
	// true
}

func ExampleLabeledGraph_RemoveVertex() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	g.RemoveVertex(1)
	fmt.Println(g.ContainsVertex(1))
	fmt.Println(g.ContainsVertex(2))
	fmt.Println(g.ContainsEdge(1, 2))
	// Output:
	// false
	// true
	// false
}

func ExampleLabeledGraph_ContainsEdge() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	fmt.Println(g.ContainsEdge(1, 2))
	fmt.Println(g.ContainsEdge(2, 3))
	// Output:
	// true
	// false
}

func ExampleLabeledGraph_AddEdge() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	fmt.Println(g.ContainsEdge(1, 2))
	// Output:
	// true
}

func ExampleLabeledGraph_RemoveEdge() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	g.RemoveEdge(1, 2)
	fmt.Println(g.ContainsEdge(1, 2))
	// Output:
	// false
}

func ExampleLabeledGraph_SetLabel() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	g.SetLabel(1, 2, "B")
	l, _ := g.Label(1, 2)
	fmt.Println(l)
	// Output:
	// B
}

func ExampleLabeledGraph_Directed() {
	g1 := labeledgraph.New[int, string]()
	g2 := labeledgraph.New[int, string](labeledgraph.Directed())
	fmt.Println(g1.Directed())
	fmt.Println(g2.Directed())
	// Output:
	// false
	// true
}

func ExampleLabeledGraph_Order() {
	g := labeledgraph.New[int, string]()
	g.AddVertex(1)
	g.AddVertex(2)
	fmt.Println(g.Order())
	// Output:
	// 2
}

func ExampleLabeledGraph_Size() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	g.AddEdge(2, 3, "B")
	fmt.Println(g.Size())
	// Output:
	// 2
}

func ExampleLabeledGraph_Label() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	l1, ok1 := g.Label(1, 2)
	l2, ok2 := g.Label(2, 3)
	fmt.Println(l1, ok1)
	fmt.Println(l2, ok2)
	// Output:
	// A true
	//  false
}

func ExampleLabeledGraph_Degree() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	g.AddEdge(1, 3, "B")
	deg, ok := g.Degree(1)
	fmt.Println(deg, ok)
	// Output:
	// 2 true
}

func ExampleLabeledGraph_InDegree() {
	g := labeledgraph.New[int, string](labeledgraph.Directed())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	deg, ok := g.InDegree(2)
	fmt.Println(deg, ok)
	// Output:
	// 2 true
}

func ExampleLabeledGraph_OutDegree() {
	g := labeledgraph.New[int, string](labeledgraph.Directed())
	g.AddEdge(1, 2, "A")
	g.AddEdge(1, 3, "B")
	deg, ok := g.OutDegree(1)
	fmt.Println(deg, ok)
	// Output:
	// 2 true
}

func ExampleLabeledGraph_Vertices() {
	g := labeledgraph.New[int, string]()
	g.AddVertex(1)
	g.AddVertex(2)
	g.AddVertex(3)

	verts := slices.Collect(g.Vertices())
	slices.Sort(verts)
	fmt.Println(verts)
	// Output:
	// [1 2 3]
}

func ExampleLabeledGraph_Neighbors() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	g.AddEdge(1, 3, "B")

	neighbors := slices.Collect(g.Neighbors(1))
	slices.Sort(neighbors)
	fmt.Println(neighbors)
	// Output:
	// [2 3]
}

func ExampleLabeledGraph_Successors() {
	g := labeledgraph.New[int, string](labeledgraph.Directed())
	g.AddEdge(1, 2, "A")
	g.AddEdge(1, 3, "B")
	g.AddEdge(3, 1, "C") // In-edge shouldn't be counted in successors

	succs := slices.Collect(g.Successors(1))
	slices.Sort(succs)
	fmt.Println(succs)
	// Output:
	// [2 3]
}

func ExampleLabeledGraph_Predecessors() {
	g := labeledgraph.New[int, string](labeledgraph.Directed())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	g.AddEdge(2, 4, "C") // Out-edge shouldn't be counted in predecessors

	preds := slices.Collect(g.Predecessors(2))
	slices.Sort(preds)
	fmt.Println(preds)
	// Output:
	// [1 3]
}

func ExampleLabeledGraph_Edges() {
	g := labeledgraph.New[int, string](labeledgraph.Directed())
	g.AddEdge(1, 2, "A")
	g.AddEdge(2, 3, "B")

	// Collect into string format for stable sorting/output
	var edges []string
	for u, v := range g.Edges() {
		l, _ := g.Label(u, v)
		edges = append(edges, fmt.Sprintf("%d->%d: %s", u, v, l))
	}
	slices.Sort(edges)
	for _, e := range edges {
		fmt.Println(e)
	}
	// Output:
	// 1->2: A
	// 2->3: B
}

func ExampleLabeledGraph_IncidentEdges() {
	g := labeledgraph.New[int, string](labeledgraph.Directed())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	g.AddEdge(2, 4, "C")

	var edges []string
	for u, v := range g.IncidentEdges(2) {
		l, _ := g.Label(u, v)
		edges = append(edges, fmt.Sprintf("%d->%d: %s", u, v, l))
	}
	slices.Sort(edges)
	for _, e := range edges {
		fmt.Println(e)
	}
	// Output:
	// 1->2: A
	// 2->4: C
	// 3->2: B
}

func ExampleLabeledGraph_InIncidentEdges() {
	g := labeledgraph.New[int, string](labeledgraph.Directed())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	g.AddEdge(2, 4, "C")

	var edges []string
	for u, v := range g.InIncidentEdges(2) {
		l, _ := g.Label(u, v)
		edges = append(edges, fmt.Sprintf("%d->%d: %s", u, v, l))
	}
	slices.Sort(edges)
	for _, e := range edges {
		fmt.Println(e)
	}
	// Output:
	// 1->2: A
	// 3->2: B
}

func ExampleLabeledGraph_OutIncidentEdges() {
	g := labeledgraph.New[int, string](labeledgraph.Directed())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	g.AddEdge(2, 4, "C")
	g.AddEdge(2, 5, "D")

	var edges []string
	for u, v := range g.OutIncidentEdges(2) {
		l, _ := g.Label(u, v)
		edges = append(edges, fmt.Sprintf("%d->%d: %s", u, v, l))
	}
	slices.Sort(edges)
	for _, e := range edges {
		fmt.Println(e)
	}
	// Output:
	// 2->4: C
	// 2->5: D
}

func ExampleLabeledGraph_Clear() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	g.Clear()
	fmt.Println(g.Order())
	fmt.Println(g.Size())
	// Output:
	// 0
	// 0
}

func ExampleLabeledGraph_Clone() {
	g1 := labeledgraph.New[int, string]()
	g1.AddEdge(1, 2, "A")
	g2 := g1.Clone()
	fmt.Println(g1.ContainsEdge(1, 2))
	fmt.Println(g2.ContainsEdge(1, 2))

	g1.RemoveEdge(1, 2)
	fmt.Println(g1.ContainsEdge(1, 2))
	fmt.Println(g2.ContainsEdge(1, 2))
	// Output:
	// true
	// true
	// false
	// true
}

func ExampleLabeledGraph_Equal() {
	g1 := labeledgraph.New[int, string]()
	g1.AddEdge(1, 2, "A")

	g2 := labeledgraph.New[int, string]()
	g2.AddEdge(1, 2, "A")

	g3 := labeledgraph.New[int, string]()
	g3.AddEdge(1, 2, "B")

	eq := func(l1, l2 string) bool { return l1 == l2 }
	fmt.Println(g1.Equal(g2, eq))
	fmt.Println(g1.Equal(g3, eq))
	// Output:
	// true
	// false
}

func ExampleLabeledGraph_String() {
	g2 := labeledgraph.New[int, string](labeledgraph.Directed())
	g2.AddEdge(1, 2, "A")
	fmt.Println(g2.String())
	// Output:
	// [1 -> 2: A]
}
