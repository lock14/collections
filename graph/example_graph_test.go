package graph_test

import (
	"fmt"
	"github.com/lock14/collections/graph"
	"slices"
)

func ExampleGraph_directed() {
	// Create a directed graph
	g := graph.New[string](graph.Directed())

	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")

	fmt.Println("Successors of A:")
	succs := slices.Collect(g.Successors("A"))
	slices.Sort(succs)
	for _, v := range succs {
		fmt.Println("-", v)
	}

	fmt.Println("Predecessors of D:")
	preds := slices.Collect(g.Predecessors("D"))
	slices.Sort(preds)
	for _, v := range preds {
		fmt.Println("-", v)
	}

	// Output:
	// Successors of A:
	// - B
	// - C
	// Predecessors of D:
	// - B
	// - C
}

func ExampleGraph_undirected() {
	// Create an undirected graph (default)
	g := graph.New[string]()

	g.AddEdge("A", "B")
	g.AddEdge("A", "C")

	// In an undirected graph, Successors, Predecessors, and Neighbors are identical
	fmt.Println("Neighbors of A:")
	neighbors := slices.Collect(g.Neighbors("A"))
	slices.Sort(neighbors)
	for _, v := range neighbors {
		fmt.Println("-", v)
	}

	// Output:
	// Neighbors of A:
	// - B
	// - C
}
