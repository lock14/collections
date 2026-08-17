package labeledgraph_test

import (
	"cmp"
	"fmt"
	"github.com/lock14/collections/hashset"
	"github.com/lock14/collections/heap"
	"github.com/lock14/collections/labeledgraph"
	"slices"
)

func ExampleLabeledGraph() {
	// Create a directed graph with string vertices and string edge labels (relationships)
	g := labeledgraph.New[string, string](labeledgraph.WithDirected())

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
	g2 := labeledgraph.New[int, string](labeledgraph.WithDirected())
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
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	deg, ok := g.InDegree(2)
	fmt.Println(deg, ok)
	// Output:
	// 2 true
}

func ExampleLabeledGraph_OutDegree() {
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
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
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
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
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
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
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
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
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
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
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
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
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
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
	g2 := labeledgraph.New[int, string](labeledgraph.WithDirected())
	g2.AddEdge(1, 2, "A")
	fmt.Println(g2.String())
	// Output:
	// [1 -> 2: A]
}

type PathNode struct {
	vertex string
	dist   int
}

func ExampleLabeledGraph_shortestPath() {
	// Create a directed graph with string vertices and int edge weights (distances)
	g := labeledgraph.New[string, int](labeledgraph.WithDirected())

	// Add edges and their weights
	g.AddEdge("A", "B", 4)
	g.AddEdge("A", "C", 2)
	g.AddEdge("B", "C", 5)
	g.AddEdge("B", "D", 10)
	g.AddEdge("C", "E", 3)
	g.AddEdge("E", "D", 4)
	g.AddEdge("D", "F", 11)

	// We want to find the shortest path from "A" to all other nodes.
	start := "A"

	// Priority queue to select the node with the minimum distance
	pq := heap.New[PathNode](heap.WithComparator(func(a, b PathNode) int {
		return cmp.Compare(a.dist, b.dist)
	}))

	// Map to store the shortest known distance to each node
	distances := make(map[string]int)

	// Map to store the previous node to reconstruct the path
	previous := make(map[string]string)

	distances[start] = 0
	pq.Add(PathNode{vertex: start, dist: 0})

	for !pq.Empty() {
		current := pq.Remove()

		// If we already found a shorter path to this node, skip processing
		if current.dist > distances[current.vertex] {
			continue
		}

		for neighbor := range g.Successors(current.vertex) {
			weight, _ := g.Label(current.vertex, neighbor)
			newDist := current.dist + weight

			// If we haven't visited the neighbor, or found a shorter path
			if d, visited := distances[neighbor]; !visited || newDist < d {
				distances[neighbor] = newDist
				previous[neighbor] = current.vertex
				pq.Add(PathNode{vertex: neighbor, dist: newDist})
			}
		}
	}

	// Print shortest distances from A
	for _, v := range []string{"A", "B", "C", "D", "E", "F"} {
		fmt.Printf("Distance from %s to %s: %d\n", start, v, distances[v])
	}

	// Print shortest path to F
	path := []string{"F"}
	curr := "F"
	for curr != start {
		curr = previous[curr]
		path = append([]string{curr}, path...)
	}
	fmt.Printf("Shortest path to F: %v\n", path)

	// Output:
	// Distance from A to A: 0
	// Distance from A to B: 4
	// Distance from A to C: 2
	// Distance from A to D: 9
	// Distance from A to E: 5
	// Distance from A to F: 20
	// Shortest path to F: [A C E D F]
}

func ExampleLabeledGraph_minimumSpanningTree() {
	// Create an undirected graph with string vertices and int edge weights
	g := labeledgraph.New[string, int]()

	// Add edges for a connected network
	g.AddEdge("A", "B", 4)
	g.AddEdge("A", "H", 8)
	g.AddEdge("B", "C", 8)
	g.AddEdge("B", "H", 11)
	g.AddEdge("C", "I", 2)
	g.AddEdge("C", "F", 4)
	g.AddEdge("C", "D", 7)
	g.AddEdge("D", "E", 9)
	g.AddEdge("D", "F", 14)
	g.AddEdge("E", "F", 10)
	g.AddEdge("F", "G", 2)
	g.AddEdge("G", "I", 6)
	g.AddEdge("G", "H", 1)
	g.AddEdge("H", "I", 7)

	type Edge struct {
		u, v   string
		weight int
	}

	// Priority queue to select the minimum weight edge
	pq := heap.New[Edge](heap.WithComparator(func(a, b Edge) int {
		return cmp.Compare(a.weight, b.weight)
	}))

	visited := hashset.New[string]()
	mstEdges := make([]Edge, 0)
	totalCost := 0

	// Start from arbitrary vertex "A"
	start := "A"
	visited.Add(start)

	// Add all edges from the start vertex to the queue
	for neighbor := range g.Neighbors(start) {
		weight, _ := g.Label(start, neighbor)
		pq.Add(Edge{u: start, v: neighbor, weight: weight})
	}

	for !pq.Empty() {
		e := pq.Remove()

		// If the destination is already visited, this edge forms a cycle in the MST
		if visited.Contains(e.v) {
			continue
		}

		// Include edge in MST
		visited.Add(e.v)
		mstEdges = append(mstEdges, e)
		totalCost += e.weight

		// Add new outgoing edges from the newly visited vertex to the queue
		for neighbor := range g.Neighbors(e.v) {
			if !visited.Contains(neighbor) {
				weight, _ := g.Label(e.v, neighbor)
				pq.Add(Edge{u: e.v, v: neighbor, weight: weight})
			}
		}
	}

	fmt.Printf("Total MST Cost: %d\n", totalCost)
	fmt.Printf("MST Edge Count: %d\n", len(mstEdges))

	// Output:
	// Total MST Cost: 37
	// MST Edge Count: 8
}
