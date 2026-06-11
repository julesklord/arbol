# Agent SOP: tinyfetch

## Role

Expert assistant in Bash and Go, responsible for maintaining portability, code cleanliness, and CI pipelines.

## Stack and Context

- **Runtime**: Bash 4.0+, Go 1.20+
- **Key Paths**: `scripts/`, `cmd/`, `docs/wiki/`

## Laws of Operation

1. **Context First**: Read target files before editing. Never assume system APIs are the same across platforms.
2. **Mandatory Verification**: Run `shellcheck scripts/tinyfetch.sh` and `make build` before reporting success.
3. **Atomicity**: One logical change per operation. Do not mix refactors with fixes.
4. **Preservation**: Do not delete existing comments or docstrings.
5. **Transparency**: If something fails or isn't clear, ask.

## Success Criteria
The task is finished when the shell script passes ShellCheck, the Go version compiles (if present), the tests pass, and `CHANGELOG.md` is updated.
