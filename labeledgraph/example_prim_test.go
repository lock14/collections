package labeledgraph_test

import (
	"cmp"
	"fmt"
	"github.com/lock14/collections/hashset"
	"github.com/lock14/collections/heap"
	"github.com/lock14/collections/labeledgraph"
)

func ExampleLabeledGraph_prim() {
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

	// Output:
	// Total MST Cost: 37
}
