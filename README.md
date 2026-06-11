# tinyfetch

> Minimal fastfetch-style status tool — a tiny, focused utility to show system info in a compact, beautiful layout.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Built With](https://img.shields.io/badge/Built%20With-Shell%20%2F%20Go-brightgreen.svg)](go.mod)

## Overview

**tinyfetch** is a tiny, dependency-free CLI status utility designed to quickly fetch and display essential system information in your terminal. It offers two implementations:
1. A robust, portable POSIX-compliant Shell script (`scripts/tinyfetch.sh`).
2. A high-performance compiled Go binary.

Both versions present a side-by-side colorized representation of the host OS logo and core resource metrics (Host, OS, Kernel, Uptime, Shell, CPU, Memory, and Disk usage).

## Installation

### Via Script (Quickstart)

```bash
git clone https://github.com/julesklord/tinyfetch.git
cd tinyfetch
scripts/tinyfetch.sh
```

### From Source (Go Version)

Ensure Go 1.20+ is installed:

```bash
make build
```

### System-Wide Installation

Install the script or compiled binary into `/usr/local/bin`:

```bash
sudo make install
```

## Usage

Run the utility directly from your shell:

```bash
tinyfetch
```

### Command Reference

| Option | Short Alias | Description |
| :--- | :--- | :--- |
| `--help` | `-h` | Display version and usage instructions. |
| `--no-ascii` | | Omit the side-by-side system ASCII logo. |

## Architecture

The utility checks the runtime operating system and dynamically resolves resource usage metrics through the most efficient native queries available.

```mermaid
graph TD
    User([UserActor]) -->|RunCommand| Cli[CliTool]
    Cli -->|QuerySystem| OsCheck{OsCheck}
    OsCheck -->|QueryLinux| LinuxParser[LinuxParser]
    OsCheck -->|QueryDarwin| DarwinParser[DarwinParser]
    LinuxParser -->|ReadFiles| ProcUptime[/proc/uptime]
    LinuxParser -->|ReadFiles| ProcMem[/proc/meminfo]
    LinuxParser -->|ReadFiles| ProcCpu[/proc/cpuinfo]
    DarwinParser -->|CallSysctl| Sysctl[Sysctl]
    DarwinParser -->|CallVMStat| VMStat[VMStat]
    LinuxParser -->|Render| SideBySide[SideBySideRenderer]
    DarwinParser -->|Render| SideBySide
    SideBySide -->|WriteOutput| Stdout[TerminalStdout]
```

## Changelog

Detailed release history is documented in [CHANGELOG.md](CHANGELOG.md).

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

