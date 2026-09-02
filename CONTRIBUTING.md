# Contributing to collections

Thank you for your interest in contributing to `github.com/lock14/collections`! This document outlines our repository standards, architectural guidelines, development workflow, and testing requirements.

---

## 1. Repository Design Principles

*   **Generics & Type Safety:** Leverage Go generics (`any`, `comparable`, `cmp.Ordered`, custom constraints). Avoid `interface{}` / `any` boxing or runtime type assertions. Sorted collections must provide both `New` accepting a custom `comparator.Comparator[T]` and `NewOrdered` for `cmp.Ordered` types.
*   **Interface Compliance & Compile-Time Assertions:** Concrete data structures should implement the generic collection interfaces defined in `collections.go` (`Collection[T]`, `MutableCollection[T]`, `List[T]`, `Queue[T]`, `Stack[T]`, `Deque[T]`, `Set[T]`, `Map[K, V]`, etc.) and `fmt.Stringer`. Every struct must declare compile-time interface assertions:
    ```go
    var (
        _ collections.MutableSet[T] = (*HashSet[T])(nil)
        _ fmt.Stringer              = (*HashSet[T])(nil)
    )
    ```
*   **Standard Iteration (Go 1.23+ `iter`):**
    *   Expose standard range-over-func iterators using `iter.Seq[T]` (e.g. `All()`) or `iter.Seq2[K, V]` (e.g. `All()`, `Keys()`, `Values()`).
    *   Iterators must terminate immediately when `!yield(...)` returns `false` without leaking resources.
    *   Iterator implementations should minimize or eliminate allocation overhead in loops.
*   **Non-Concurrent by Design:** Implementations are single-threaded by contract (matching Go's standard slice/map philosophy). Do not add internal mutexes or sync primitives to collection structs. Concurrency is strictly the caller's responsibility.
*   **Zero-Allocation Reads:** Read operations (`Get`, `Contains`, `Peek`, `PeekFront`, `PeekBack`, `Size`, `Empty`) must avoid heap allocations (0 B/op and 0 allocs/op in benchmarks).
*   **Predictable Memory & Growth:** Growth policies should be amortized O(1) (e.g. geometric resizing) and support explicit capacity APIs (`NewWithCapacity`).
*   **Panic & Error Handling Alignment with Built-in Types:**
    *   **Slice-Like Indexing (`Get`, `Set`):** Bounds check violations (`idx < 0 || idx >= size`) MUST panic with a clear index error (matching built-in Go `slice[i]` panic).
    *   **Consumable Extraction & Positional Access (`Remove()`, `First()`, `Last()`, `Peek()`, `Pop()`, `RemoveFront()`, `RemoveBack()`, `PollFirst()`, `PollLast()`):** Extraction or access on an empty collection MUST panic with a clear error (matching built-in `s[0]` / `s[len-1]` semantics).
    *   **Targeted Deletions & Map Lookups (`RemoveElement(T)`, `Map.Remove(K)`, `Map.Get(K)`):** Targeted removal of a specific element or key is a silent no-op if absent (matching `delete(m, k)`). Key lookup uses the comma-ok idiom `(V, bool)` without panicking on missing keys.
*   **String Representations (`fmt.Stringer`):**
    *   `Collection[T]`: Formatted as `[e1 e2 e3]` (space-separated, `[]` when empty).
    *   `Map[K, V]`: Formatted as `map[k1:v1 k2:v2]` (`map[]` when empty).
    *   `Graph[V]` / `LabeledGraph[V, L]`: Formatted as `graph[...]` or `digraph[...]`.

---

## 2. Package File Structure

Single-collection packages (e.g. `hashset`, `hashmap`, `arraydeque`, `treeset`, etc.) MUST consist of exactly four files:

1.  **`<pkg>.go`**: Primary implementation file containing type declarations, compile-time assertions, configuration options, constructors, and public/private methods.
2.  **`<pkg>_test.go`**: Table-driven unit tests verifying functional correctness, edge cases, and bounds checks (`package <pkg>`).
3.  **`<pkg>_bench_test.go`**: Performance benchmarks measuring throughput and zero-allocation invariants (`package <pkg>`).
4.  **`example_<pkg>_test.go`**: Consumer-facing runnable examples with `// Output:` comments (`package <pkg>_test`).

---

## 3. Development Setup & Workflow

### Prerequisites

*   Go 1.27 or higher
*   Git

### Clone & Test

```bash
# Clone the repository
git clone https://github.com/lock14/collections.git
cd collections

# Verify dependencies
go mod verify

# Run all unit tests with race detection
go test -race ./...

# Run all benchmarks
go test -bench=. -benchmem ./...
```

### Table-Driven Testing

All unit tests must be table-driven:
*   Define a `cases` slice of structs containing `name string`.
*   Iterate with `for _, tc := range cases`.
*   Rebind `tc := tc` and invoke `t.Run(tc.name, func(t *testing.T) { ... })`.
*   Use `t.Parallel()` where appropriate with independent collection instances.

### Benchmarking & Performance Regression Verification

*   Benchmark functions belong in `<pkg>_bench_test.go` named `func Benchmark<Type>_<Method>(b *testing.B)`.
*   Always call `b.ReportAllocs()`.
*   Isolate setup with `b.ResetTimer()`.
*   Read paths must show 0 B/op and 0 allocs/op.
*   To verify no performance regressions against the `main` branch locally:
    ```bash
    go install filippo.io/mostly-harmless/benchdiff@latest
    benchdiff -base-ref origin/main -- -run '^$' -bench . -benchmem -count=6 ./...
    ```

---

## 4. Pre-Completion Verification Checklist

Before opening a pull request or submitting code, ensure all of the following commands execute without errors or warnings:

```bash
# 1. Format code
gofmt -s -w .

# 2. Verify module hygiene
go mod tidy
git diff --exit-code go.mod go.sum

# 3. Static analysis & Vet
go vet ./...
golangci-lint run ./...

# 4. Security vulnerability check
govulncheck ./...

# 5. Run tests with race detection
go test -race ./...
```

---

## 5. Pull Request Guidelines

1.  **Branch Naming:** Use clear branch names such as `feat/new-collection`, `fix/bounds-check`, `perf/hash-distribution`.
2.  **Commit Messages:** Write concise, descriptive commit messages following the Conventional Commits specification (e.g. `feat: add LinkedHashSet`, `fix: correct deque wrap-around index`).
3.  **PR Description:** Fill out the [Pull Request Template](.github/PULL_REQUEST_TEMPLATE.md) completely, detailing changes, architectural adherence, and benchmark comparisons.
4.  **Documentation Sync:** When adding or modifying APIs, update the corresponding Go doc comments, `README.md`, and `AGENTS.md` in the same pull request.
