# collections

[![Go Version](https://img.shields.io/github/go-mod/go-version/lock14/collections)](https://go.dev/)
[![Build Status](https://img.shields.io/github/actions/workflow/status/lock14/collections/go.yml?branch=main)](https://github.com/lock14/collections/actions)
[![Benchmarks](https://img.shields.io/github/actions/workflow/status/lock14/collections/benchmark.yml?branch=main&label=Benchmarks)](https://github.com/lock14/collections/actions/workflows/benchmark.yml)
[![Coverage](https://img.shields.io/codecov/c/github/lock14/collections)](https://codecov.io/gh/lock14/collections)
[![Go Report Card](https://goreportcard.com/badge/github.com/lock14/collections)](https://goreportcard.com/report/github.com/lock14/collections)
[![Go Reference](https://pkg.go.dev/badge/github.com/lock14/collections.svg)](https://pkg.go.dev/github.com/lock14/collections)

Generic data structures for Go 1.23+. Focuses on type safety, minimal allocations, and predictable performance.

## Usage

```bash
go get github.com/lock14/collections
```

```go
package main

import (
	"fmt"
	
	"github.com/lock14/collections/comparator"
	"github.com/lock14/collections/treeset"
)

func main() {
	set := treeset.New[int](comparator.NaturalOrder[int]())
	
	set.Add(5)
	set.Add(1)
	set.Add(10)
	
	for val := range set.All() {
		fmt.Println(val)
	}
}
```

## Data Structures

Implementations leverage Go generics to eliminate `interface{}` boxing and runtime type assertions.

*   **Maps**
    *   `hashmap`: Map backed by a hash table.
    *   `linkedhashmap`: Hash map preserving insertion or access order.
    *   `treemap`: Sorted map backed by a B-Tree.
*   **Sets**
    *   `hashset`: Set backed by a hash table.
    *   `linkedhashset`: Hash set preserving insertion or access order.
    *   `treeset`: Sorted set backed by a B-Tree.
    *   `bitset`: Word-aligned dense integer set.
*   **Lists, Queues & Stacks**
    *   `arraylist`: Dynamically resizing array.
    *   `linkedlist`: Doubly-linked list.
    *   `arraydeque`: Double-ended queue backed by a ring buffer.
    *   `heap`: Priority queue.
*   **Graphs**
    *   `graph`: Directed and undirected graphs.
    *   `labeledgraph`: Graphs with labeled edges.

## Performance & Testing

Design prioritizes mechanical sympathy and GC pressure reduction.

*   **Zero-Allocation Reads**: Read paths (`Get`, `Contains`, etc.) bypass heap allocations.
*   **Continuous Benchmarking**: CI gates PRs via `cob`, comparing allocation metrics and execution times against `main`. Regressions fail the build.
*   **Test Coverage**: Table-driven tests are mandatory. Edge cases, bounds checks, and generic fallback paths must be explicitly exercised.

## Contributing

Submit PRs with passing tests and benchmarks. Table-driven testing is required. Performance regressions will not be merged.

## License

Apache 2.0. See [LICENSE](LICENSE).
