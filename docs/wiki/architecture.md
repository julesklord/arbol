# Architecture

This project is intentionally small. The primary artifact will be a single, dependency-free shell script in `scripts/` and optionally a small Go binary at the repo root for performance.

Design goals:
- Minimal dependencies
- Easy install (Makefile)
- Clear, testable outputs for CI
