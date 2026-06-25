# collections

[![Go Version](https://img.shields.io/github/go-mod/go-version/lock14/collections)](https://go.dev/)
[![Build Status](https://img.shields.io/github/actions/workflow/status/lock14/collections/go.yml?branch=main)](https://github.com/lock14/collections/actions)
[![Benchmarks](https://img.shields.io/github/actions/workflow/status/lock14/collections/benchmark.yml?branch=main&label=Benchmarks)](https://github.com/lock14/collections/actions/workflows/benchmark.yml)
[![Coverage](https://img.shields.io/codecov/c/github/lock14/collections)](https://codecov.io/gh/lock14/collections)
[![Go Report Card](https://goreportcard.com/badge/github.com/lock14/collections)](https://goreportcard.com/report/github.com/lock14/collections)
[![Go Reference](https://pkg.go.dev/badge/github.com/lock14/collections.svg)](https://pkg.go.dev/github.com/lock14/collections)

A comprehensive, generic, and highly optimized data structures and collections library for Go 1.23+.

## Getting Started

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
	// Create a new TreeSet for integers
	set := treeset.New[int](comparator.NaturalOrder[int]())
	
	set.Add(5)
	set.Add(1)
	set.Add(10)
	
	// Iterates in sorted order: 1, 5, 10
	for val := range set.All() {
		fmt.Println(val)
	}
}
```

## Features

This library provides generic implementations of standard data structures, eliminating the need for `interface{}` and type assertions.

*   **Maps**
    *   `hashmap`: Standard hash table based map.
    *   `linkedhashmap`: Hash map that preserves insertion or access order.
    *   `treemap`: B-Tree based sorted map.
*   **Sets**
    *   `hashset`: Standard hash table based set.
    *   `linkedhashset`: Hash set that preserves insertion or access order.
    *   `treeset`: B-Tree based sorted set.
    *   `bitset`: Memory-efficient set of integers.
*   **Lists & Queues**
    *   `arraylist`: Dynamically resizing array.
    *   `linkedlist`: Doubly-linked list.
    *   `arraydeque`: Double-ended queue backed by a ring buffer.
    *   `heap`: Priority queue.
*   **Graphs**
    *   `graph`: Directed and undirected graphs.
    *   `labeledgraph`: Graphs with labeled edges.

## Performance & Testing

This project takes performance seriously. Core data structures are heavily optimized to reduce memory allocations and Garbage Collection (GC) overhead.

*   **100% Table-Driven Tests**: We enforce strict table-driven testing conventions for robustness.
*   **Zero-Allocation Paths**: Read operations (like `Get` and `Contains`) across the library are heavily optimized to be allocation-free where possible.
*   **CI Benchmark Regression**: Every Pull Request is automatically profiled via a GitHub Actions pipeline against the `main` branch to strictly prevent any performance degradations.

## Contributing

Contributions are welcome! We require all tests to pass and benchmarks to show no regression. Ensure any new logic is accompanied by strict table-driven testing.

## License

This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE) file for details.
