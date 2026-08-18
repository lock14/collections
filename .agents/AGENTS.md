# Repository Design Principles

*   **Generics & Type Safety:** Leverage Go 1.23+ generics (`any`, `comparable`, custom constraints). Avoid `interface{}` / `any` boxing or runtime type assertions where type parameters can be used.
*   **Non-Concurrent by Design:** Implementations are single-threaded by design (matching Go's standard slice/map philosophy). Do not add internal mutexes or goroutines to collection structs.
*   **Zero-Allocation Reads:** Read operations (e.g., `Get`, `Contains`, `Peek`, `Size`, `Empty`) must avoid heap allocations.

# Testing Conventions

For all tests in this repository, strictly adhere to the following convention:

*   **Table-Driven Tests:** Always use table-driven tests for unit testing. This is the idiomatic Go way and makes tests easy to read, extend, and maintain.
*   **Structure:**
    *   Define a `cases` slice of structs.
    *   Each test case struct should have at least a `name string` field.
    *   Iterate through the cases with a `for _, tc := range cases` loop.
    *   Inside the loop, rebind `tc := tc` (for pre-Go 1.22 closure safety, good practice generally) and use `t.Run(tc.name, func(t *testing.T) { ... })`.
    *   Use `t.Parallel()` inside the sub-test where appropriate. Ensure each parallel sub-test creates its own independent collection instance.
*   **Exception:** Stress tests or randomized large-scale tests (e.g. testing deep merges with random elements) may use procedural/loop-based structures where a table structure would be impractical. Default to table-driven whenever testing specific inputs and expected outputs.

# Examples (`example_*_test.go`)

*   **Runnable Examples:** Exported collection types and significant methods must have runnable examples in an `example_<pkg>_test.go` file.
*   **Verified Output:** Always include an `// Output:` comment at the end of example functions so they are validated by `go test ./...` and rendered cleanly on `pkg.go.dev`.
*   **Package Naming:** Place example tests in the `<pkg>_test` package to demonstrate public API usage from a consumer's perspective.

# Benchmarks (`*_bench_test.go`)

*   **Location & Naming:** Benchmarks belong in `<pkg>_bench_test.go` and must follow `func Benchmark<Type>_<Method>(b *testing.B)`.
*   **Allocation Tracking:** Always call `b.ReportAllocs()` in benchmark functions.
*   **Timer Control:** Isolate setup code from measured code using `b.ResetTimer()` and `b.StopTimer()` / `b.StartTimer()` when appropriate.
*   **Regression Guard:** Ensure read-path benchmarks show 0 B/op and 0 allocs/op.

# Code Style and Quality

*   **Formatting:** Always run `gofmt -s -w .` after making code changes.
*   **Documentation:** All exported packages, types, interfaces, constants, and functions must have proper Go doc comments adhering to standard Go conventions (checked by `revive`).
*   **Linting:** Verify code with `golangci-lint run ./...` before completing tasks.

# Pre-Completion Verification Checklist

Before concluding any code modification task, run:
1. `gofmt -s -w .`
2. `go test ./...`
3. `golangci-lint run ./...`
