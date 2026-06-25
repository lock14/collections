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
