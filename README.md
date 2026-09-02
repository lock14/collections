# collections

[![Go Version](https://img.shields.io/github/go-mod/go-version/lock14/collections)](https://go.dev/)
[![CI](https://img.shields.io/github/actions/workflow/status/lock14/collections/ci.yml?branch=main&label=CI)](https://github.com/lock14/collections/actions/workflows/ci.yml)
[![Benchmarks](https://img.shields.io/github/actions/workflow/status/lock14/collections/benchmark.yml?branch=main&label=Benchmarks)](https://github.com/lock14/collections/actions/workflows/benchmark.yml)
[![Security](https://img.shields.io/github/actions/workflow/status/lock14/collections/security.yml?branch=main&label=Security)](https://github.com/lock14/collections/actions/workflows/security.yml)
[![Coverage](https://img.shields.io/codecov/c/github/lock14/collections)](https://codecov.io/gh/lock14/collections)
[![Go Reference](https://pkg.go.dev/badge/github.com/lock14/collections.svg)](https://pkg.go.dev/github.com/lock14/collections)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Generic data structures for Go 1.27+. Focuses on type safety, minimal allocations, and predictable performance.

## Usage

```bash
go get github.com/lock14/collections
```

```go
package main

import (
	"fmt"

	"github.com/lock14/collections/treeset"
)

func main() {
	set := treeset.NewOrdered[int]()

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
*   **Lists, Queues, & Stacks**
    *   `arraylist`: Dynamically resizing array.
    *   `linkedlist`: Doubly-linked list.
    *   `arraydeque`: Double-ended queue backed by a ring buffer.
    *   `heap`: Priority queue.
*   **Strings & Prefixes**
    *   `trie`: String and generic slice (`[]E`) prefix trees with prefix queries (`KeysWithPrefix`, `LongestPrefixOf`, etc.).
*   **Graphs**
    *   `graph`: Directed and undirected graphs.
    *   `labeledgraph`: Graphs with labeled edges.
*   **Comparators**
    *   `comparator`: Type-safe element comparison functions (`NaturalOrder`, `Reverse`).

## Performance & Testing

Design prioritizes mechanical sympathy and GC pressure reduction.

*   **Zero-Allocation Reads**: Read paths (`Get`, `Contains`, etc.) bypass heap allocations.
*   **Continuous Benchmarking**: CI evaluates PRs via `benchdiff` (powered by `benchstat`), performing statistical comparisons across allocation metrics and execution times against `main`.
*   **Test Coverage**: Table-driven tests are mandatory. Edge cases, bounds checks, and generic fallback paths must be explicitly exercised.

## Concurrency

Implementations in this library are **not thread-safe** by design, matching Go standard library types like slices and maps. If a collection is accessed concurrently by multiple goroutines and at least one modifies it, access must be synchronized externally (e.g. using `sync.RWMutex` or `sync.Mutex`).

## Contributing

We welcome contributions! Please review our [Contributing Guidelines](CONTRIBUTING.md) for architectural invariants, testing standards, and pull request verification instructions.

## Security

Please report security issues responsibly according to our [Security Policy](SECURITY.md).

## License

Apache 2.0. See [LICENSE](LICENSE).
