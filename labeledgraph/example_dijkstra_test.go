package labeledgraph_test

import (
	"cmp"
	"fmt"
	"github.com/lock14/collections/heap"
	"github.com/lock14/collections/labeledgraph"
)

type PathNode struct {
	vertex string
	dist   int
}

func ExampleLabeledGraph_dijkstra() {
	// Create a directed graph with string vertices and int edge weights (distances)
	g := labeledgraph.New[string, int](labeledgraph.Directed())

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
