package labeledgraph_test

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/lock14/collections/hashset"
	"github.com/lock14/collections/heap"
	"github.com/lock14/collections/labeledgraph"
)

func ExampleNew() {
	g := labeledgraph.New[string, int](labeledgraph.WithDirected(), labeledgraph.WithCapacity(16))
	g.AddEdge("A", "B", 5)
	fmt.Println(g.Directed(), g.Order(), g.Size())
	// Output:
	// true 2 1
}

func ExampleLabeledGraph() {
	// Create a directed graph with string vertices and string edge labels (relationships)
	g := labeledgraph.New[string, string](labeledgraph.WithDirected())

	g.AddEdge("Alice", "Bob", "Knows")
	g.AddEdge("Bob", "Charlie", "Likes")
	g.AddEdge("Alice", "Charlie", "Dislikes")

	for target := range g.Successors("Alice") {
		relationship, _ := g.Label("Alice", target)
		fmt.Printf("Alice %s %s\n", relationship, target)
	}

	// Unordered output:
	// Alice Knows Bob
	// Alice Dislikes Charlie
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

	for v := range g.Vertices() {
		fmt.Println(v)
	}
	// Unordered output:
	// 1
	// 2
	// 3
}

func ExampleLabeledGraph_Neighbors() {
	g := labeledgraph.New[int, string]()
	g.AddEdge(1, 2, "A")
	g.AddEdge(1, 3, "B")

	for v := range g.Neighbors(1) {
		fmt.Println(v)
	}
	// Unordered output:
	// 2
	// 3
}

func ExampleLabeledGraph_Successors() {
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	g.AddEdge(1, 2, "A")
	g.AddEdge(1, 3, "B")
	g.AddEdge(3, 1, "C") // In-edge shouldn't be counted in successors

	for v := range g.Successors(1) {
		fmt.Println(v)
	}
	// Unordered output:
	// 2
	// 3
}

func ExampleLabeledGraph_Predecessors() {
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	g.AddEdge(2, 4, "C") // Out-edge shouldn't be counted in predecessors

	for u := range g.Predecessors(2) {
		fmt.Println(u)
	}
	// Unordered output:
	// 1
	// 3
}

func ExampleLabeledGraph_Edges() {
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	g.AddEdge(1, 2, "A")
	g.AddEdge(2, 3, "B")

	for u, v := range g.Edges() {
		l, _ := g.Label(u, v)
		fmt.Printf("%d->%d: %s\n", u, v, l)
	}
	// Unordered output:
	// 1->2: A
	// 2->3: B
}

func ExampleLabeledGraph_IncidentEdges() {
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	g.AddEdge(2, 4, "C")

	for u, v := range g.IncidentEdges(2) {
		l, _ := g.Label(u, v)
		fmt.Printf("%d->%d: %s\n", u, v, l)
	}
	// Unordered output:
	// 1->2: A
	// 2->4: C
	// 3->2: B
}

func ExampleLabeledGraph_InIncidentEdges() {
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	g.AddEdge(2, 4, "C")

	for u, v := range g.InIncidentEdges(2) {
		l, _ := g.Label(u, v)
		fmt.Printf("%d->%d: %s\n", u, v, l)
	}
	// Unordered output:
	// 1->2: A
	// 3->2: B
}

func ExampleLabeledGraph_OutIncidentEdges() {
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	g.AddEdge(1, 2, "A")
	g.AddEdge(3, 2, "B")
	g.AddEdge(2, 4, "C")
	g.AddEdge(2, 5, "D")

	for u, v := range g.OutIncidentEdges(2) {
		l, _ := g.Label(u, v)
		fmt.Printf("%d->%d: %s\n", u, v, l)
	}
	// Unordered output:
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
	g := labeledgraph.New[int, string](labeledgraph.WithDirected())
	g.AddEdge(1, 2, "A")
	fmt.Println(g.String())
	// Output:
	// digraph[1 -> 2: A]
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

	type pathNode struct {
		vertex string
		dist   int
	}

	// We want to find the shortest path from "A" to all other nodes.
	start := "A"

	// Priority queue to select the node with the minimum distance
	pq := heap.New[pathNode](heap.WithComparator(func(a, b pathNode) int {
		return cmp.Compare(a.dist, b.dist)
	}))

	// Map to store the shortest known distance to each node
	distances := make(map[string]int)

	// Map to store the previous node to reconstruct the path
	previous := make(map[string]string)

	distances[start] = 0
	pq.Add(pathNode{vertex: start, dist: 0})

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
				pq.Add(pathNode{vertex: neighbor, dist: newDist})
			}
		}
	}

	// Print shortest distances from A
	for _, v := range []string{"A", "B", "C", "D", "E", "F"} {
		fmt.Printf("Distance from %s to %s: %d\n", start, v, distances[v])
	}

	// Reconstruct shortest path to F
	var path []string
	for curr := "F"; curr != ""; curr = previous[curr] {
		path = append(path, curr)
	}
	slices.Reverse(path)
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

	type edge struct {
		u, v   string
		weight int
	}

	// Priority queue to select the minimum weight edge
	pq := heap.New[edge](heap.WithComparator(func(a, b edge) int {
		return cmp.Compare(a.weight, b.weight)
	}))

	visited := hashset.New[string]()
	mstEdges := make([]edge, 0)
	totalCost := 0

	// Start from arbitrary vertex "A"
	start := "A"
	visited.Add(start)

	// Add all edges from the start vertex to the queue
	for neighbor := range g.Neighbors(start) {
		weight, _ := g.Label(start, neighbor)
		pq.Add(edge{u: start, v: neighbor, weight: weight})
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
				pq.Add(edge{u: e.v, v: neighbor, weight: weight})
			}
		}
	}

	fmt.Printf("Total MST Cost: %d\n", totalCost)
	fmt.Printf("MST Edge Count: %d\n", len(mstEdges))

	// Output:
	// Total MST Cost: 37
	// MST Edge Count: 8
}
