# Changelog

All notable changes to this project will be documented in this file.

## 0.2.3 - Packages manager plugin
- Created `plugins/packages.sh` to count installed package manager details, supporting native `pacman` packages and foreign (AUR) packages via helper detection (`paru`/`yay`).

## 0.2.2 - Documentation improvements
- Updated `docs/wiki/architecture.md` with system overview, component diagram, and Architecture Decision Records (ADRs) for visual alignment and modular extensibility.
- Updated `docs/wiki/development.md` with a comprehensive guide to writing custom plugins, including constraints, stdout format guidelines, and examples.


## 0.2.1 - Visual alignment fixes & Developer plugins
- Fixed box-drawing layout misalignment in Go by using `utf8.RuneCountInString` instead of byte counts for progress bars and unicode characters.
- Upgraded the Git status plugin with Nerd Fonts, staged/modified/untracked counters, and upstream sync indicators.
- Created 3 new developer plugins: Weather (`plugins/weather.sh`), Docker (`plugins/docker.sh`), and Media Player (`plugins/media.sh`).


## 0.2.0 - Multiplatform stability & Go version
- Refactored `scripts/tinyfetch.sh` for Linux & macOS portability under `set -e`.
- Fixed ShellCheck `SC2034` warning.
- Added support for a modular plugins folder (`./plugins/`).
- Added dynamic distro ASCII logo loading from `ascii/` text files with automatic fallbacks.
- Added visual progress bars for memory and disk usage metrics.
- Replaced standard printing with an innovative double-pane terminal card layout with box-drawing borders.
- Implemented compiled Go version in `cmd/tinyfetch/main.go`.
- Added test suite in `tests/test.sh`.
- Created standard FMG files (`docs/AGENT.md`, `docs/GEMINI.md`, etc.).





## 0.1.0 - Initial skeleton
- Repository scaffolded with standard structure and README.

