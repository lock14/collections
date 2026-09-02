# Repository Design Principles

*   **Generics & Type Safety:** Leverage Go 1.23+ generics (`any`, `comparable`, `cmp.Ordered`, custom constraints). Avoid `interface{}` / `any` boxing or runtime type assertions where type parameters can be used. Sorted collections must provide both `New` accepting a custom `comparator.Comparator[T]` and `NewOrdered` for `cmp.Ordered` types.
*   **Interface Compliance & Compile-Time Assertions:** Where applicable, concrete data structures should implement the generic collection interfaces defined in `collections.go` (`Collection[T]`, `MutableCollection[T]`, `List[T]`, `Queue[T]`, `Stack[T]`, `Deque[T]`, `Set[T]`, `Map[K, V]`, etc.) and `fmt.Stringer`. Every struct must declare compile-time interface assertions using the standard Go idiom `var _ Interface = (*Type)(nil)` (or `var _ fmt.Stringer = Type{}` for value receivers) to guarantee interface compliance at build time.
*   **Standard Iteration (Go 1.23+ `iter`):**
    *   Collections must expose standard range-over-func iterators using `iter.Seq[T]` (e.g. `All()`) or `iter.Seq2[K, V]` for associative types (e.g. `All()`, `Keys()`, `Values()`).
    *   Iterator implementations must respect early termination immediately when `!yield(...)` returns `false` without leaking resources or performing unnecessary work.
    *   Iterator implementations should minimize or eliminate allocation overhead when consumed in `for ... range` loops.
*   **Non-Concurrent by Design:** Implementations are single-threaded by contract (matching Go's standard slice/map philosophy). Do not add internal mutexes, sync primitives, or goroutines to collection structs. Concurrency is strictly the caller's responsibility.
*   **Zero-Allocation Reads:** Read operations (e.g., `Get`, `Contains`, `Peek`, `PeekFront`, `PeekBack`, `Size`, `Empty`) must avoid heap allocations (guaranteed 0 B/op and 0 allocs/op in benchmarks).
*   **Predictable Memory & Growth:** Growth policies should be amortized O(1) (e.g. geometric resizing) and shrink gracefully where appropriate or support explicit capacity/trim APIs (e.g. `NewWithCapacity`).
*   **String Representations (`fmt.Stringer`):**
    *   **`Collection[T]` Types:** String representations MUST match the format of Go's built-in slice: `[e1 e2 e3]` (space-separated, `[]` when empty).
    *   **`Map[K, V]` Types:** String representations MUST match the format of Go's built-in map: `map[k1:v1 k2:v2]` (prefixed by `map[`, space-separated `k:v` with no space after the colon, `map[]` when empty).
    *   **`Graph[V]` / `LabeledGraph[V, L]` Types:** String representations MUST use `graph[...]` (for undirected graphs) or `digraph[...]` (for directed graphs), containing comma-separated edges (`u - v` / `u -> v: label`) and isolated vertices (`graph[]` / `digraph[]` when empty).
*   **Panic & Error Handling Alignment with Built-in Types:**
    *   **Slice-Like Indexing (`Get`, `Set`):** Bounds check violations (`idx < 0 || idx >= size`) MUST panic with a clear index error (matching built-in Go slice index panic `slice[i]`).
    *   **Consumable Extraction & Positional Access (`Remove()`, `First()`, `Last()`, `Peek()`, `Pop()`, `RemoveFront()`, `RemoveBack()`, `PollFirst()`, `PollLast()`):** Unconditional element extraction or positional access on an empty collection MUST panic with a clear error (matching built-in slice head/tail indexing `s[0]` / `s[len-1]` and `slices.Min`/`slices.Max` semantics). Every `MutableCollection[T]` (including sets) supports `Remove() T` to enable the consumable worklist/drain pattern (popping elements until `Empty()`).
    *   **Targeted Deletions & Map Lookups (`RemoveElement(T)`, `Map.Remove(K)`, `Map.Get(K)`):** Targeted removal of a specific element from a set (`RemoveElement(T)`) or key from a map (`Map.Remove(K)`) MUST be a silent no-op if the item/key is absent (matching built-in `delete(m, k)`). Key lookup MUST use the standard Go comma-ok idiom `(V, bool)` (never panic for missing keys or empty maps).
    *   **No "Unsupported Operation" Runtime Panics:** Interfaces must strictly adhere to the Interface Segregation Principle. Collections must never declare or embed methods that compile but throw runtime panics due to being unsupported (e.g., `AddFirst`/`PutFirst` on sorted collections).

# Package File Structure & Layout Conventions

Every collection package in this repository must maintain a consistent, standardized file structure:

*   **Standard 4-File Package Structure:** Single-collection packages (e.g. `hashset`, `hashmap`, `arraydeque`, `treeset`, `linkedhashmap`, etc.) MUST consist of exactly four files:
    1.  **`<pkg>.go`**: The primary implementation file containing the type declaration, compile-time interface assertions, configuration options, constructors, and public/private methods.
    2.  **`<pkg>_test.go`**: Table-driven unit tests verifying functional correctness, edge cases, and bounds checks (`package <pkg>`).
    3.  **`<pkg>_bench_test.go`**: Performance benchmarks measuring throughput and zero-allocation invariants (`package <pkg>`).
    4.  **`example_<pkg>_test.go`**: Consumer-facing runnable examples with `// Output:` comments (`package <pkg>_test`).
*   **Multi-File / Complex Packages:** When a package implements multiple representations (e.g. `trie` with `string.go` and `slice.go`) or contains a dedicated complex algorithmic engine (e.g. `treemap` with `btree.go`), the entry point, interfaces, and constructors remain in `<pkg>.go`, while variant-specific implementations and tests use clear, consistent naming (e.g. `slice_test.go`). Extraneous or ad-hoc test files (e.g. `coverage_patch_test.go`, `set.go`, `map.go`) are prohibited.
*   **Internal Section Ordering in `<pkg>.go`:**
    To maintain visual and structural consistency across the repository, code within `<pkg>.go` should adhere to the following sequence:
    1.  Package doc comment and `package <pkg>` declaration.
    2.  `import` block (grouped and sorted).
    3.  Compile-time interface compliance assertions (`var _ collections.Interface = (*Type)(nil)` and `var _ fmt.Stringer = (*Type)(nil)`).
    4.  Constants & defaults (e.g. `DefaultCapacity`, `DefaultDegree`, order constants).
    5.  Configuration types and functional options (`config`, `Option`, `With...`).
    6.  Struct definitions (`type Type[T any] struct { ... }`).
    7.  Constructors (`New`, `NewOrdered`, `NewWithCapacity`, etc.).
    8.  Public methods grouped logically (Mutators, Accessors/Queries, Iterators, Stringer).
    9.  Private helper functions and traversal logic.

# Naming & Code Style Conventions

*   **Receiver Naming:** Use short (1-2 letters), mnemonic receiver names consistently across all methods on a struct. Never use `this`, `self`, or vary receiver names within the same type:
    *   `s *HashSet[T]` / `s *TreeSet[T]` / `s *LinkedHashSet[T]` / `s *BitSet`
    *   `m *HashMap[K, V]` / `m *TreeMap[K, V]` / `m *LinkedHashMap[K, V]` / `m *sliceMap[E, V]` / `m *stringMap[V]`
    *   `d *ArrayDeque[T]`
    *   `l *LinkedList[T]` / `l *SliceWrapper[T]`
    *   `h *Heap[T]`
    *   `g *Graph[V]` / `g *LabeledGraph[V, L]`
*   **Avoid Identifier Shadowing:** Never use parameter or local variable names that shadow Go built-in identifiers (`max`, `min`, `len`, `cap`, `new`, `clear`, `copy`, `close`, `delete`). Use descriptive identifiers such as `maxElements`, `capacity`, `cmpFunc`, or `limit`.
*   **Documentation Comments:** All exported packages, types, interfaces, constants, options, and methods must have comprehensive Go doc comments starting with the symbol name and adhering to standard Go / `revive` conventions.
*   **Formatting:** Always run `gofmt -s -w .` after making code changes.

# Testing Conventions

For all tests in this repository, strictly adhere to the following conventions:

*   **Table-Driven Tests:** Always use table-driven tests for unit testing.
*   **Structure:**
    *   Define a `cases` slice of structs.
    *   Each test case struct should have at least a `name string` field.
    *   Iterate through the cases with a `for _, tc := range cases` loop.
    *   Inside the loop, rebind `tc := tc` (for closure safety) and use `t.Run(tc.name, func(t *testing.T) { ... })`.
    *   Use `t.Parallel()` inside sub-tests where appropriate. Ensure each parallel sub-test creates its own independent collection instance.
*   **Exceptions:** Stress tests or randomized large-scale tests (e.g. testing deep merges, randomized mutation sequences) may use procedural/loop-based structures where a table structure would be impractical. Default to table-driven whenever testing specific inputs and expected outputs.

# Benchmarking Conventions

*   **Location & Naming:** Benchmarks belong in `<pkg>_bench_test.go` and must follow `func Benchmark<Type>_<Method>(b *testing.B)`.
*   **Allocation Tracking:** Always call `b.ReportAllocs()` in benchmark functions.
*   **Timer Isolation:** Isolate setup code from measured operations using `b.ResetTimer()` and `b.StopTimer()` / `b.StartTimer()` when appropriate.
*   **Measurement Patterns:**
    *   **Steady-State Mutators (Add/Remove, Push/Pop):** Prepopulate the collection to a fixed size $N$ outside the loop, `b.ResetTimer()`, and execute balanced add/remove operations inside the `for i := 0; i < b.N; i++` loop to measure steady-state per-operation latency rather than unbounded growth or repeated reallocations.
    *   **Zero-Allocation Reads (Get, Contains, Peek, Iteration):** Prepopulate the collection outside the loop, `b.ResetTimer()`, and execute read-only operations inside the loop.
    *   **Bulk / Preallocated Operations:** Test preallocation efficiencies (e.g. `NewWithCapacity`, `Add_Preallocated`) to verify zero amortized allocation overhead.
*   **Regression Guard:** Ensure read-path benchmarks show 0 B/op and 0 allocs/op.

# Examples (`example_*_test.go`)

*   **Runnable Examples:** Exported collection types and significant methods must have runnable examples in an `example_<pkg>_test.go` file.
*   **Verified Output:** Always include an `// Output:` comment at the end of example functions so they are validated by `go test ./...` and rendered cleanly on `pkg.go.dev`.
*   **Package Naming:** Place example tests in the `<pkg>_test` package to demonstrate public API usage from a consumer's perspective.

# Documentation Roles & Boundaries

*   **`README.md` (User-Facing):** Reserved strictly for public/consumer-facing documentation, quickstart examples, feature overviews, high-level capabilities, installation instructions, and package indices. It should **not** delve into internal implementation details, internal algorithmic plumbing, or agent/contributor rules.
*   **`AGENTS.md` (Architecture & Implementation Guidelines):** The canonical place to encode architecture, design invariants, implementation decisions, performance constraints, code guidelines, testing/benchmarking rules, and agent workflows.
*   **Documentation Synchronization:** Whenever any code change modifies, adds, or removes APIs, behaviors, or package structures, all relevant documentation (**Go doc comments**, **`README.md`**, and **`AGENTS.md`**) MUST be updated in the same change to keep documentation strictly in sync with the codebase.

# Pull Request & Workflow Conventions

*   **Accurate & Detailed PR Descriptions:** PR descriptions must comprehensively capture all introduced features, constructors, public APIs, performance benchmarks, and any toolchain or CI workflow modifications made in the branch.
*   **PR Description Synchronization:** Whenever subsequent commits in a PR branch modify, fix, or expand upon the original changeset (e.g. CI adjustments, bug fixes, refactoring), the PR description must be proactively updated via `gh pr edit` to reflect the latest state of the PR.

# Continuous Learning & Self-Correction

*   **Codifying Corrections:** Whenever an agent is corrected by the user regarding an architectural decision, code style, testing pattern, or repository convention, the agent MUST update `AGENTS.md` (or the relevant package documentation) to codify the underlying principle so future agent sessions adhere to it automatically.

# Pre-Completion Verification Checklist

Before concluding any code modification task, agents MUST run and verify the following commands succeed without errors or warnings:
1. `gofmt -s -w .`
2. `go mod tidy && git diff --exit-code go.mod go.sum`
3. `go test -race ./...`
4. `golangci-lint run ./...` (or `revive ./... && go vet ./...`)
5. `govulncheck ./...`


