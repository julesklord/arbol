# Agent SOP (Standard Operating Procedure)

This file is the entry point for any AI agent (Gemini, Claude, GPT, etc.) working on this repository. Read it, or don't say we didn't warn you.

## General Instructions

1. **Familiarization**: Before making any changes, read `docs/wiki/index.md` and `FMG-REPO-BIBLE.md` (in the `jules_dev_standard` repository). Know your place.
2. **Compliance**: Strictly follow the laws defined in `docs/wiki/agent-sop.md`. No exceptions.
3. **Identity**: Consult `docs/SOUL.md` to understand the tone and principles of this project.

## Agent Initialization Commands

1. The shell script is located under `scripts/mini-fetch.sh` and has no external dependencies.
2. If working on the Go version:
   - Ensure you are working under the Go modules system (go 1.20+).
   - Use `make build` to verify compiling.
3. Run `shellcheck scripts/mini-fetch.sh` before committing any changes to the shell script.
