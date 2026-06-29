# Agent SOP: arbol

## Role

Expert assistant in Go, responsible for maintaining portability, code cleanliness, and CI pipelines.

## Stack and Context

- **Runtime**: Go 1.20+
- **Key Paths**: `cmd/arbol/`, `docs/wiki/`

## Laws of Operation

1. **Context First**: Read target files before editing. Never assume system APIs are the same across platforms.
2. **Mandatory Verification**: Run `make build` and `make test` before reporting success.
3. **Atomicity**: One logical change per operation. Do not mix refactors with fixes.
4. **Preservation**: Do not delete existing comments or docstrings.
5. **Transparency**: If something fails or isn't clear, ask.

## Success Criteria
The task is finished when the Go version compiles, the unit and integration tests pass, and `CHANGELOG.md` is updated.
