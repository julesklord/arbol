#!/usr/bin/env bash
# Minimal fastfetch-like script
# - Minimal dependencies (bash, coreutils, awk, sed)
# - Small, clean ASCII art + aligned key: value output
# - Usage: ./scripts/mini-fetch.sh [--no-ascii]

set -u

NO_ASCII=0
for arg in "$@"; do
  case "$arg" in
    --no-ascii) NO_ASCII=1 ;;
    -h|--help)
      echo "Usage: $0 [--no-ascii]";
      exit 0 ;;
  esac
done

# Basic color helpers (fallback to empty if tput not available)
if command -v tput >/dev/null 2>&1; then
  BOLD=$(tput bold)
  DIM=$(tput dim)
  RESET=$(tput sgr0)
  RED=$(tput setaf 1)
  GREEN=$(tput setaf 2)
  YELLOW=$(tput setaf 3)
  BLUE=$(tput setaf 4)
  MAGENTA=$(tput setaf 5)
  CYAN=$(tput setaf 6)
else
  BOLD=""; DIM=""; RESET=""; RED=""; GREEN=""; YELLOW=""; BLUE=""; MAGENTA=""; CYAN=""
fi

info_val() { printf "%s%s%s" "$BOLD" "$1" "$RESET"; }

# Collect system info with safe fallbacks
HOSTNAME=$(hostname 2>/dev/null || echo "unknown")
OS=$( { lsb_release -ds 2>/dev/null || cat /etc/os-release 2>/dev/null | awk -F= '/^PRETTY_NAME=/{gsub(/\"/,"",$2); print $2; exit}'; } 2>/dev/null || echo "$(uname -s)")
KERNEL=$(uname -r 2>/dev/null || echo "unknown")
UPTIME=$(awk '{print int($1/3600) "h " int(($1%3600)/60) "m"}' /proc/uptime 2>/dev/null || uptime -p 2>/dev/null || echo "n/a")
SHELL_NAME=$(basename "${SHELL:-sh}")

# CPU
CPU_NAME=$(awk -F: '/model name/{print $2; exit}' /proc/cpuinfo 2>/dev/null | sed 's/^[ \t]*//')
if [ -z "$CPU_NAME" ]; then
  CPU_NAME=$(uname -p 2>/dev/null || echo "unknown")
fi

# Memory
read -r MEM_TOTAL MEM_FREE < <(awk '/MemTotal/ {t=$2} /MemAvailable/ {a=$2} END {print t, a}' /proc/meminfo 2>/dev/null || echo "0 0")
if [ "$MEM_TOTAL" -gt 0 ]; then
  MEM_USAGE=$(awk -v t=$MEM_TOTAL -v a=$MEM_FREE 'BEGIN{printf "%.0f%%", (t-a)/t*100}')
  MEM_HUMAN=$(awk -v t=$MEM_TOTAL 'function hr(x){split("KB MB GB TB",u); for(i=1;x>=1024&&i<4;i++) x/=1024; printf "%d%s", x, u[i]} END{hr(t)}')
else
  MEM_USAGE="n/a"
  MEM_HUMAN="n/a"
fi

# Disk (root)
DF_OUT=$(df -h --output=source,size,used,avail,pcent,target / 2>/dev/null | sed -n '2p' || true)
if [ -n "$DF_OUT" ]; then
  DISK=$(echo "$DF_OUT" | awk '{print $1 " " $2 " " $5 " (" $4 " avail)"}')
else
  DISK="n/a"
fi

# Packages (try common package managers)
PKG_COUNT="-"
if command -v dpkg >/dev/null 2>&1; then
  PKG_COUNT=$(dpkg -l 2>/dev/null | awk 'NR>5{c++}END{print c+0}')
elif command -v rpm >/dev/null 2>&1; then
  PKG_COUNT=$(rpm -qa 2>/dev/null | wc -l)
fi

# Resolution (try xdpyinfo, fbset as fallback)
RESOLUTION="-"
if command -v xdpyinfo >/dev/null 2>&1; then
  RESOLUTION=$(xdpyinfo | awk -F: '/dimensions/{gsub(/^[ \t]+/,"",$2); print $2; exit}')
elif [ -n "${DISPLAY-}" ] && command -v wmctrl >/dev/null 2>&1; then
  RESOLUTION=$(wmctrl -d | awk '{print $9; exit}')
fi

# Terminal
TERM_NAME=${TERM:-"-"}

# Git branch for current dir (helpful when running in a project)
GIT_BRANCH="-"
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "-")
fi

# Small ASCII logos to choose from (minimal)
ASCII_SMALL=(
"   .----.  "
"  / .-""-\\ "
" | |     | |")

ASCII_SML=(
"  _  _     "
" | || |___ "
" | __ / -_)"
" |_||_\\___|"
)

print_kv() {
  local key="$1"; shift
  local val="$*"
  printf "%s%-14s %s\n" "$DIM" "$key:" "$RESET$val"
}

render() {
  # Compose right column lines
  local lines=()
  lines+=("Host" "$HOSTNAME")
  lines+=("OS" "$OS")
  lines+=("Kernel" "$KERNEL")
  lines+=("Uptime" "$UPTIME")
  lines+=("Shell" "$SHELL_NAME")
  lines+=("CPU" "$CPU_NAME")
  lines+=("Memory" "$MEM_USAGE ($MEM_HUMAN)")
  lines+=("Disk" "$DISK")
  lines+=("Packages" "$PKG_COUNT")
  lines+=("Resolution" "$RESOLUTION")
  lines+=("Terminal" "$TERM_NAME")
  lines+=("Git branch" "$GIT_BRANCH")

  # Convert to pretty output
  if [ "$NO_ASCII" -eq 1 ]; then
    # simple key-values
    for ((i=0;i<${#lines[@]};i+=2)); do
      print_kv "${lines[i]}" "${lines[i+1]}"
    done
    return
  fi

  # With ASCII art left column
  local left=()
  left=("${ASCII_SML[@]}")
  # ensure left has at least as many lines as right (pad)
  local right_lines=$(( ${#lines[@]} / 2 ))
  local left_lines=${#left[@]}
  if [ $left_lines -lt $right_lines ]; then
    for ((i=left_lines;i<right_lines;i++)); do left+=(" "); done
  fi

  # print side by side
  for ((i=0, j=0; i<${#left[@]}; i++, j+=2)); do
    l=${left[i]}
    if [ $j -lt ${#lines[@]} ]; then
      k=${lines[j]}
      v=${lines[j+1]}
      printf "%s %s %s\n" "$CYAN$l$RESET" "$(printf "%-12s" "$k:")" "$BOLD$v$RESET"
    else
      printf "%s\n" "$CYAN$l$RESET"
    fi
  done
}

# Run
render

exit 0
