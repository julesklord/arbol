# Development Guide

This document guides developers on local setup, running tests, and creating custom plugins.

## Prerequisites

- Go 1.20+

## Local Setup

1. Clone the repository.
2. Build the binary using the Makefile:
   ```bash
   make build
   ```
3. Run the test suite to verify everything compiles and behaves correctly:
   ```bash
   make test
   ```

## Creating Custom Plugins

`arbol` scans the `./plugins` directory and its `extended/` subdirectory for executable scripts or binaries. You can write plugins in Bash, Python, Go, Node, or any other scripting language.

### Simple Plugins

1. **Location**: Place your script under `./plugins/` (e.g., `plugins/battery.sh`).
2. **Executability**: The file must be executable. Run `chmod +x plugins/my-plugin` to enable it.
3. **Stdout Format**: The plugin can output a single line or multiple lines.
   - The first line can follow the `Label: Value` pattern. If a colon is detected, the key (`Git`) will be printed in blue, and the value in default colors.
   - If multiple lines are printed, subsequent lines will automatically be parsed and rendered as nested sub-branches under the parent node in the tree.
4. **Error Handling**: If the plugin fails, it must exit silently (`exit 0`) and print nothing. If a plugin prints nothing, the row/node is omitted from the tree.

### Extended Plugins

1. **Location**: Place your script under `./plugins/extended/` (e.g., `plugins/extended/sys_dashboard.sh`).
2. **Executability**: Must be executable (`chmod +x`).
3. **Stdout Format**: Can output multiple lines. **arbol** will dynamically calculate widths and align the borders of the third pane symmetrically.
4. **Error Handling**: If the plugin fails or is not applicable, it must exit silently (`exit 0`) and print nothing. If all extended plugins print nothing, the third column will be cleanly omitted from the output.

### Example Simple Plugin (Shell)

`plugins/battery.sh`:
```bash
#!/usr/bin/env bash
set -euo pipefail

# Check if battery path exists (Linux)
if [ -d /sys/class/power_supply/BAT0 ]; then
  capacity=$(cat /sys/class/power_supply/BAT0/capacity)
  status=$(cat /sys/class/power_supply/BAT0/status)
  echo "Battery: 🔋 ${capacity}%"
  echo "Status: ${status}"
else
  exit 0
fi
```
Make it executable:
```bash
chmod +x plugins/battery.sh
```

## Git Workflow & Conventions

Follow the rules in [hygiene.md](file:///home/julesklord/Proyectos/repos/arbol/docs/wiki/hygiene.md) for conventional commit messages. Keep changes atomic.
