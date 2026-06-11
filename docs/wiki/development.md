# Development Guide

This document guides developers on local setup, running tests, and creating custom plugins.

## Prerequisites

- Go 1.20+ (optional, for the Go binary)
- Bash 4.0+
- `shellcheck` (for code quality linting)

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

`mini-fetch` scans the `./plugins` directory for executable scripts or binaries. You can write plugins in Bash, Python, Go, Node, or any other scripting language.

### Plugin Requirements

1. **Location**: Place your script under `./plugins/` (e.g., `plugins/battery.sh`).
2. **Executability**: The file must be executable. Run `chmod +x plugins/my-plugin` to enable it.
3. **Stdout Format**: The plugin must output exactly one line. It can follow one of two patterns:
   - **Label format**: `Label: Value` (e.g., `Git: main`). If a colon is detected, the key (`Git`) will be printed in blue, and the value (`main`) in default colors.
   - **Plain format**: `Value` (e.g., `☀️ +20°C`). The tool will automatically use the capitalized filename as the label (e.g., `weather.sh` becomes `Weather: ☀️ +20°C`).
4. **Error Handling**: If the plugin fails (e.g., no internet, missing tools), it must exit silently (`exit 0`) and print nothing. If a plugin prints nothing, the row is omitted from the dashboard.

### Example Plugin (Shell)

`plugins/battery.sh`:
```bash
#!/usr/bin/env bash
set -euo pipefail

# Check if battery path exists (Linux)
if [ -d /sys/class/power_supply/BAT0 ]; then
  capacity=$(cat /sys/class/power_supply/BAT0/capacity)
  status=$(cat /sys/class/power_supply/BAT0/status)
  echo "Battery: 🔋 ${capacity}% (${status})"
else
  exit 0
fi
```
Make it executable:
```bash
chmod +x plugins/battery.sh
```

## Git Workflow & Conventions

Follow the rules in [hygiene.md](file:///home/julesklord/Proyectos/repos/mini-fetch/docs/wiki/hygiene.md) for conventional commit messages. Keep changes atomic.

