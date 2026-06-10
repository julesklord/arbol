# mini-fetch

Minimal fastfetch-style status tool — a tiny, focused utility to show system info in a compact, beautiful layout.

Overview
--------

mini-fetch is a minimal, dependency-light system information tool intended to be developed and maintained in this repository. It contains a tiny shell script starter and room for a Go implementation if desired.

Quickstart
----------

1. Clone the repo and develop the script in `scripts/`.
2. Build (if implementing in Go): `go build -o mini-fetch .` from the repo root.
3. Run: `./mini-fetch --no-ascii` or `scripts/mini-fetch.sh`.

Repository Structure
--------------------

- scripts/: place to develop the POSIX/bash implementation.
- docs/: project-level documentation and wiki.
- LICENSE, VERSION, CHANGELOG.md: repo metadata.

Contributing
------------

Follow the repository standard located in `docs/wiki/`.
