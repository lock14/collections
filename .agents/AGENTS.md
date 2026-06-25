# Testing Conventions

For all tests in this repository, strictly adhere to the following convention:

*   **Table-Driven Tests:** Always use table-driven tests for unit testing. This is the idiomatic Go way and makes tests extremely easy to read, extend, and debug.
*   **Structure:**
    *   Define a `cases` slice of structs.
# Testing Conventions

For all tests in this repository, strictly adhere to the following convention:

*   **Table-Driven Tests:** Always use table-driven tests for unit testing. This is the idiomatic Go way and makes tests extremely easy to read, extend, and debug.
*   **Structure:**
    *   Define a `cases` slice of structs.
    *   Each test case struct should have at least a `name string` field.
    *   Iterate through the cases with a `for _, tc := range cases` loop.
    *   Inside the loop, rebind `tc := tc` (for pre-Go 1.22 closure safety, though fine to do generally) and use `t.Run(tc.name, func(t *testing.T) { ... })`.
    *   Use `t.Parallel()` inside the sub-test where appropriate.
*   **Exception:** Stress tests, benchmarks, or highly randomized large-scale tests (e.g. testing B-Tree deep merges with 1000s of randomized elements) may use procedural/loop-based structures if a table structure would be overly verbose or impractical. But default to table-driven whenever testing specific inputs and expected outputs.

# Code Style and Formatting

*   **Formatting:** Always run `gofmt -s -w .` after making any code changes, before committing or concluding a task. This ensures all Go files are properly formatted according to standard Go conventions and simplified.
