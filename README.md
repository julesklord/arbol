# mini-fetch

Minimal fastfetch-style status tool — a tiny, focused utility to show system info in a compact, beautiful layout.

Overview
--------

mini-fetch is a small CLI utility for quickly viewing system information in your terminal. This repo follows the Jules Dev Standard. Use the `template/` directory in the org standard as reference for documentation and structure.

Quickstart
----------

Clone and run the provided shell starter:

```bash
git clone <this-repo>
cd mini-fetch
scripts/mini-fetch.sh
```

If you implement a Go version, build from the repo root:

```bash
make build
```

Install the script/binary system-wide:

```bash
sudo make install
```

Repository Structure
--------------------

- scripts/: POSIX-compatible implementation. (Primary development area.)
- docs/: repository wiki and developer guides.
- VERSION, CHANGELOG.md, LICENSE: project metadata.

Contributing
------------

Follow the standards in `repos/jules_dev_standard/template`. Create a branch, open a PR, and include tests or shellcheck where appropriate.
