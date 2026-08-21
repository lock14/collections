# Repository Design Principles

*   **Generics & Type Safety:** Leverage Go 1.23+ generics (`any`, `comparable`, `cmp.Ordered`, custom constraints). Avoid `interface{}` / `any` boxing or runtime type assertions where type parameters can be used. Sorted collections must provide both `New` accepting a custom `comparator.Comparator[T]` and `NewOrdered` for `cmp.Ordered` types.
*   **Interface Compliance:** Where applicable, concrete data structures should implement the generic collection interfaces defined in `collections.go` (`Collection[T]`, `MutableCollection[T]`, `List[T]`, `Queue[T]`, `Stack[T]`, `Deque[T]`, `Set[T]`, `Map[K, V]`, etc.).
*   **Standard Iteration (Go 1.23+ `iter`):**
    *   Collections must expose standard range-over-func iterators using `iter.Seq[T]` (e.g. `All()`) or `iter.Seq2[K, V]` for associative types (e.g. `All()`, `Keys()`, `Values()`).
    *   Iterator implementations should minimize or eliminate allocation overhead when consumed in `for ... range` loops.
*   **Non-Concurrent by Design:** Implementations are single-threaded by contract (matching Go's standard slice/map philosophy). Do not add internal mutexes, sync primitives, or goroutines to collection structs. Concurrency is strictly the caller's responsibility.
*   **Zero-Allocation Reads:** Read operations (e.g., `Get`, `Contains`, `Peek`, `PeekFront`, `PeekBack`, `Size`, `Empty`) must avoid heap allocations (guaranteed 0 B/op and 0 allocs/op in benchmarks).
*   **Predictable Memory & Growth:** Growth policies should be amortized O(1) (e.g. geometric resizing) and shrink gracefully where appropriate or support explicit capacity/trim APIs (e.g. `NewWithCapacity`).

# Go Version & Toolchain Consistency

*   **Single Source of Truth (`go.mod`):** The Go version specified in `go.mod` is the canonical version for the entire repository.
*   **CI Workflow Consistency:** GitHub Actions workflows must use `go-version-file: 'go.mod'` with `actions/setup-go@v5` rather than hardcoded version strings.

# Testing Conventions

For all tests in this repository, strictly adhere to the following conventions:

*   **Table-Driven Tests:** Always use table-driven tests for unit testing. This is the idiomatic Go way and makes tests easy to read, extend, and maintain.
*   **Structure:**
    *   Define a `cases` slice of structs.
    *   Each test case struct should have at least a `name string` field.
    *   Iterate through the cases with a `for _, tc := range cases` loop.
    *   Inside the loop, rebind `tc := tc` (for closure safety) and use `t.Run(tc.name, func(t *testing.T) { ... })`.
    *   Use `t.Parallel()` inside sub-tests where appropriate. Ensure each parallel sub-test creates its own independent collection instance.
*   **Exceptions:** Stress tests or randomized large-scale tests (e.g. testing deep merges, randomized mutation sequences) may use procedural/loop-based structures where a table structure would be impractical. Default to table-driven whenever testing specific inputs and expected outputs.

# Examples (`example_*_test.go`)

*   **Runnable Examples:** Exported collection types and significant methods must have runnable examples in an `example_<pkg>_test.go` file.
*   **Verified Output:** Always include an `// Output:` comment at the end of example functions so they are validated by `go test ./...` and rendered cleanly on `pkg.go.dev`.
*   **Package Naming:** Place example tests in the `<pkg>_test` package to demonstrate public API usage from a consumer's perspective.

# Benchmarks (`*_bench_test.go`)

*   **Location & Naming:** Benchmarks belong in `<pkg>_bench_test.go` and must follow `func Benchmark<Type>_<Method>(b *testing.B)`.
*   **Allocation Tracking:** Always call `b.ReportAllocs()` in benchmark functions.
*   **Timer Control:** Isolate setup code from measured code using `b.ResetTimer()` and `b.StopTimer()` / `b.StartTimer()` when appropriate.
*   **Regression Guard:** Ensure read-path benchmarks show 0 B/op and 0 allocs/op.

# Documentation Roles & Boundaries

*   **`README.md` (User-Facing):** Reserved strictly for public/consumer-facing documentation, quickstart examples, feature overviews, high-level capabilities, installation instructions, and package indices. It should **not** delve into internal implementation details, internal algorithmic plumbing, or agent/contributor rules.
*   **`AGENTS.md` (Architecture & Implementation Guidelines):** The canonical place to encode architecture, design invariants, implementation decisions, performance constraints, code guidelines, testing/benchmarking rules, and agent workflows.
*   **Documentation Synchronization:** Whenever any code change modifies, adds, or removes APIs, behaviors, or package structures, all relevant documentation (**Go doc comments**, **`README.md`**, and **`AGENTS.md`**) MUST be updated in the same change to keep documentation strictly in sync with the codebase.

# Code Style and Quality

*   **Formatting:** Always run `gofmt -s -w .` after making code changes.
*   **Documentation:** All exported packages, types, interfaces, constants, and functions must have comprehensive, up-to-date Go doc comments adhering to standard Go conventions (checked by `revive`).
*   **Linting:** Verify code with `golangci-lint run ./...` before completing tasks.

# Pull Request & Workflow Conventions

*   **Accurate & Detailed PR Descriptions:** PR descriptions must comprehensively capture all introduced features, constructors, public APIs, performance benchmarks, and any toolchain or CI workflow modifications made in the branch.
*   **PR Description Synchronization:** Whenever subsequent commits in a PR branch modify, fix, or expand upon the original changeset (e.g. CI adjustments, bug fixes, refactoring), the PR description must be proactively updated via `gh pr edit` to reflect the latest state of the PR.

# Continuous Learning & Self-Correction

*   **Codifying Corrections:** Whenever an agent is corrected by the user regarding an architectural decision, code style, testing pattern, or repository convention, the agent MUST update `AGENTS.md` (or the relevant package documentation) to codify the underlying principle so future agent sessions adhere to it automatically.

# Pre-Completion Verification Checklist

Before concluding any code modification task, agents MUST run and verify the following commands succeed without errors or warnings:
1. `gofmt -s -w .`
2. `go test -race ./...`
3. `golangci-lint run ./...`
