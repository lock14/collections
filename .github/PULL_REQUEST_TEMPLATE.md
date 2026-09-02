## Description

Please provide a summary of the changes introduced in this pull request, including motivation and context.

## Type of Change

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] ⚡ Performance optimization (latency, CPU, memory, or allocation improvements)
- [ ] 📝 Documentation update
- [ ] 🔧 Tooling / CI / Build configuration

## Repository Design & Architecture Checklist

- [ ] **Generics & Type Safety:** Uses Go generics cleanly; avoids unnecessary `any` boxing or runtime assertions.
- [ ] **Interface Compliance:** Concrete structures declare compile-time interface assertions (`var _ collections.Interface = (*Type)(nil)`).
- [ ] **Standard Iteration:** Exposes standard Go 1.23+ `iter.Seq` / `iter.Seq2` iterators respecting early yield termination.
- [ ] **Zero-Allocation Reads:** Read operations (`Get`, `Contains`, `Peek`, `Size`, `Empty`) produce 0 B/op and 0 allocs/op in benchmarks.
- [ ] **Stringer:** Implements `fmt.Stringer` matching Go slice `[e1 e2]` or map `map[k:v]` format.
- [ ] **Panic Semantics:** Slice-like indexing panics on out-of-bounds; positional extractions panic on empty; targeted map/set removals are silent no-ops.
- [ ] **Standard Package Layout:** Follows the 4-file package structure (`<pkg>.go`, `<pkg>_test.go`, `<pkg>_bench_test.go`, `example_<pkg>_test.go`).

## Verification & Testing

- [ ] Ran `gofmt -s -w .` (no unformatted files).
- [ ] Ran `go test -race ./...` (all unit tests pass without data races).
- [ ] Ran `golangci-lint run ./...` (or `revive ./... && go vet ./...`) with zero warnings.
- [ ] Added table-driven unit tests in `<pkg>_test.go`.
- [ ] Added/updated benchmarks in `<pkg>_bench_test.go` and verified no regressions with `benchdiff` / `benchstat`.
- [ ] Added/updated runnable examples with `// Output:` in `example_<pkg>_test.go`.
- [ ] Documentation (Go doc comments, `README.md`, `AGENTS.md`) is updated and synchronized.
