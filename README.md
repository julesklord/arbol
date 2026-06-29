# arbol

> Minimal fastfetch-style status tool — a tiny, focused utility to show system info in a compact, beautiful layout.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Built With](https://img.shields.io/badge/Built%20With-Go-brightgreen.svg)](go.mod)

## Overview

**arbol** is a tiny, dependency-free CLI status utility designed to quickly fetch and display essential system information in your terminal. It is implemented as a high-performance compiled Go binary.

It presents a side-by-side representation of the host OS logo as a TrueColor banner and core system resource metrics structured as a dynamic tree.

## Installation

### From Source

Ensure Go 1.20+ is installed:

```bash
git clone https://github.com/julesklord/arbol.git
cd arbol
make build
```

### System-Wide Installation

Install the compiled binary into `/usr/local/bin`:

```bash
sudo make install
```

## Usage

Run the utility directly from your shell:

```bash
arbol
```

### Command Reference

| Option | Short Alias | Description |
| :--- | :--- | :--- |
| `--help` | `-h` | Display version and usage instructions. |
| `--no-ascii` | | Omit the system ASCII logo. |
| `--minimal` | | Skip extended plugins and display a single info card. |
| `--noframe` | | Omit the box borders and print layout side-by-side using spaces. |
| `--output=FORMAT` | | Serialize system stats and simple plugins into structured output: `json`, `xml`, or `txt`. |
| `--logo=simple,banner` | | Toggle between simple logo styling and banner styling (default is banner). |

## Extensibility & Plugins

**arbol** is fully extensible via custom plugins. It scans the `./plugins` directory for executable scripts or binaries (written in Bash, Python, Go, Node, etc.) and appends their output dynamically to the dashboard.

It supports two types of plugins:
1. **Simple Plugins** (located in `./plugins/`): Status elements that append nested rows/details to the plugins branch.
2. **Extended Plugins** (located in `./plugins/extended/`): Multi-line, complex dashboards that render side-by-side in a separate diagnostics section.

### Creating a Simple Plugin

1. Create an executable file inside the `./plugins/` directory:
   ```bash
   touch plugins/my-plugin.sh
   chmod +x plugins/my-plugin.sh
   ```
2. Your script should output a summary line, followed by detailed stats on subsequent lines (e.g. `Network: Connected` on line 1, then detailed stats like local IP).
3. If a plugin needs to exit early or is not applicable, it should print nothing and exit with code `0`. Any plugin that produces empty output will be cleanly omitted from the dashboard.

### Creating an Extended Plugin

1. Create an executable file inside the `./plugins/extended/` directory:
   ```bash
   touch plugins/extended/my-dashboard.sh
   chmod +x plugins/extended/my-dashboard.sh
   ```
2. Your script can output multiple lines. **arbol** will dynamically calculate widths and align the borders of the third pane symmetrically using rune-count calculations.

### Included Plugins

The repository contains several useful out-of-the-box plugins under [plugins/](file:///home/julesklord/Proyectos/repos/arbol/plugins):
- **Battery** (`battery.sh`): Displays current battery percentage and charge status (supports Linux & macOS).
- **Docker** (`docker.sh`): Shows active containers and daemon status.
- **Git** (`git.sh`): Reports current branch, dirty status, and counters for staged/modified/untracked files with Nerd Fonts.
- **Network** (`ip.sh`): Shows local and external IP addresses (using `icanhazip.com` with a 1s connection timeout).
- **Kubernetes** (`k8s.sh`): Reports active kubectl context and namespace.
- **Packages** (`packages.sh`): Lists installed packages (supports `pacman`, `dpkg`, `rpm`, `flatpak`, `brew`, `snap`).
- **Weather** (`weather.sh`): Fetches temperature and sky status.
- **Media Player** (`media.sh`): Shows currently playing song or media status.

#### Extended Plugins ([plugins/extended/](file:///home/julesklord/Proyectos/repos/arbol/plugins/extended)):
- **Git Commit Graph** (`git_graph.sh`): Displays a beautiful local branch history tree visualization.
- **System Dashboard** (`sys_dashboard.sh`): Displays load averages and top memory-consuming processes.
- **Weather Forecast** (`weather_forecast.sh`): Displays a multi-line weather forecast from `wttr.in`.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
