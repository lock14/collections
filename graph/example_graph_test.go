package graph_test

import (
	"fmt"
	"github.com/lock14/collections/graph"
	"slices"
)

func ExampleGraph_directed() {
	// Create a directed graph
	g := graph.New[string](graph.WithDirected())

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

func ExampleGraph_ContainsVertex() {
	g := graph.New[int]()
	g.AddVertex(1)

	fmt.Println(g.ContainsVertex(1))
	fmt.Println(g.ContainsVertex(2))

	// Output:
	// true
	// false
}

func ExampleGraph_AddVertex() {
	g := graph.New[int]()
	g.AddVertex(1)

	fmt.Println(g.ContainsVertex(1))

	// Output:
	// true
}

func ExampleGraph_RemoveVertex() {
	g := graph.New[int]()
	g.AddEdge(1, 2)
	fmt.Println(g.ContainsVertex(1))

	g.RemoveVertex(1)
	fmt.Println(g.ContainsVertex(1))

	// Output:
	// true
	// false
}

func ExampleGraph_ContainsEdge() {
	g := graph.New[int]()
	g.AddEdge(1, 2)

	fmt.Println(g.ContainsEdge(1, 2))
	fmt.Println(g.ContainsEdge(2, 1))
	fmt.Println(g.ContainsEdge(1, 3))

	// Output:
	// true
	// true
	// false
}

func ExampleGraph_AddEdge() {
	g := graph.New[int]()
	g.AddEdge(1, 2)

	fmt.Println(g.ContainsEdge(1, 2))

	// Output:
	// true
}

func ExampleGraph_RemoveEdge() {
	g := graph.New[int]()
	g.AddEdge(1, 2)
	fmt.Println(g.ContainsEdge(1, 2))

	g.RemoveEdge(1, 2)
	fmt.Println(g.ContainsEdge(1, 2))

	// Output:
	// true
	// false
}

func ExampleGraph_Directed() {
	g1 := graph.New[int]()
	g2 := graph.New[int](graph.WithDirected())

	fmt.Println(g1.Directed())
	fmt.Println(g2.Directed())

	// Output:
	// false
	// true
}

func ExampleGraph_Order() {
	g := graph.New[int]()
	g.AddVertex(1)
	g.AddVertex(2)
	g.AddVertex(3)

	fmt.Println(g.Order())

	// Output:
	// 3
}

func ExampleGraph_Size() {
	g := graph.New[int]()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	fmt.Println(g.Size())

	// Output:
	// 2
}

func ExampleGraph_Degree() {
	g := graph.New[int]()
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)

	deg, ok := g.Degree(1)
	fmt.Println(deg, ok)

	deg2, ok2 := g.Degree(4)
	fmt.Println(deg2, ok2)

	// Output:
	// 2 true
	// 0 false
}

func ExampleGraph_InDegree() {
	g := graph.New[int](graph.WithDirected())
	g.AddEdge(1, 2)
	g.AddEdge(3, 2)

	deg, ok := g.InDegree(2)
	fmt.Println(deg, ok)

	// Output:
	// 2 true
}

func ExampleGraph_OutDegree() {
	g := graph.New[int](graph.WithDirected())
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)

	deg, ok := g.OutDegree(1)
	fmt.Println(deg, ok)

	// Output:
	// 2 true
}

func ExampleGraph_Vertices() {
	g := graph.New[int]()
	g.AddVertex(1)
	g.AddVertex(2)
	g.AddVertex(3)

	vertices := slices.Collect(g.Vertices())
	slices.Sort(vertices)
	fmt.Println(vertices)

	// Output:
	// [1 2 3]
}

func ExampleGraph_Neighbors() {
	g := graph.New[int]()
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)

	neighbors := slices.Collect(g.Neighbors(1))
	slices.Sort(neighbors)
	fmt.Println(neighbors)

	// Output:
	// [2 3]
}

func ExampleGraph_Successors() {
	g := graph.New[int](graph.WithDirected())
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)

	succs := slices.Collect(g.Successors(1))
	slices.Sort(succs)
	fmt.Println(succs)

	// Output:
	// [2 3]
}

func ExampleGraph_Predecessors() {
	g := graph.New[int](graph.WithDirected())
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)

	preds := slices.Collect(g.Predecessors(3))
	slices.Sort(preds)
	fmt.Println(preds)

	// Output:
	// [1 2]
}

func ExampleGraph_Edges() {
	g := graph.New[int](graph.WithDirected())
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	var edges [][2]int
	for u, v := range g.Edges() {
		edges = append(edges, [2]int{u, v})
	}
	slices.SortFunc(edges, func(a, b [2]int) int {
		if a[0] != b[0] {
			return a[0] - b[0]
		}
		return a[1] - b[1]
	})
	fmt.Println(edges)

	// Output:
	// [[1 2] [2 3]]
}

func ExampleGraph_IncidentEdges() {
	g := graph.New[int]()
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)

	var edges [][2]int
	for u, v := range g.IncidentEdges(1) {
		// To make the output deterministic, sort the edge's endpoints
		if u > v {
			u, v = v, u
		}
		edges = append(edges, [2]int{u, v})
	}
	slices.SortFunc(edges, func(a, b [2]int) int {
		if a[0] != b[0] {
			return a[0] - b[0]
		}
		return a[1] - b[1]
	})
	fmt.Println(edges)

	// Output:
	// [[1 2] [1 3]]
}

func ExampleGraph_InIncidentEdges() {
	g := graph.New[int](graph.WithDirected())
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)

	var edges [][2]int
	for u, v := range g.InIncidentEdges(3) {
		edges = append(edges, [2]int{u, v})
	}
	slices.SortFunc(edges, func(a, b [2]int) int {
		if a[0] != b[0] {
			return a[0] - b[0]
		}
		return a[1] - b[1]
	})
	fmt.Println(edges)

	// Output:
	// [[1 3] [2 3]]
}

func ExampleGraph_OutIncidentEdges() {
	g := graph.New[int](graph.WithDirected())
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)

	var edges [][2]int
	for u, v := range g.OutIncidentEdges(1) {
		edges = append(edges, [2]int{u, v})
	}
	slices.SortFunc(edges, func(a, b [2]int) int {
		if a[0] != b[0] {
			return a[0] - b[0]
		}
		return a[1] - b[1]
	})
	fmt.Println(edges)

	// Output:
	// [[1 2] [1 3]]
}

func ExampleGraph_Clear() {
	g := graph.New[int]()
	g.AddEdge(1, 2)
	fmt.Println(g.Order(), g.Size())

	g.Clear()
	fmt.Println(g.Order(), g.Size())

	// Output:
	// 2 1
	// 0 0
}

func ExampleGraph_Clone() {
	g := graph.New[int]()
	g.AddEdge(1, 2)

	g2 := g.Clone()
	fmt.Println(g2.ContainsEdge(1, 2))

	g2.AddEdge(2, 3)
	fmt.Println(g.ContainsEdge(2, 3))
	fmt.Println(g2.ContainsEdge(2, 3))

	// Output:
	// true
	// false
	// true
}

func ExampleGraph_Equal() {
	g1 := graph.New[int]()
	g1.AddEdge(1, 2)

	g2 := graph.New[int]()
	g2.AddEdge(1, 2)

	g3 := graph.New[int]()
	g3.AddEdge(1, 3)

	fmt.Println(g1.Equal(g2))
	fmt.Println(g1.Equal(g3))

	// Output:
	// true
	// false
}

func ExampleGraph_String() {
	g := graph.New[int](graph.WithDirected())
	g.AddEdge(1, 2)
	g.AddVertex(3)

	fmt.Println(g.String())

	// Output:
	// digraph[1 -> 2, 3]
}
