#!/usr/bin/env bash
# Starter POSIX-compatible mini-fetch script
set -euo pipefail

if [ "${1-}" = "--help" ] || [ "${1-}" = "-h" ]; then
  echo "Usage: $0 [--no-ascii]"
  exit 0
fi

NO_ASCII=0
for a in "$@"; do
  case "$a" in
    --no-ascii) NO_ASCII=1 ;;
  esac
done

printf "Host: %s\n" "$(hostname)"
printf "OS: %s\n" "$(grep '^PRETTY_NAME' /etc/os-release 2>/dev/null | cut -d= -f2 | tr -d '"' || uname -s)"
printf "Kernel: %s\n" "$(uname -r)"
printf "Uptime: %s\n" "$(awk '{print int($1/3600) "h " int(($1%3600)/60) "m"}' /proc/uptime 2>/dev/null || uptime -p)"
printf "Shell: %s\n" "${SHELL-}" 
printf "CPU: %s\n" "$(awk -F: '/model name/{print $2; exit}' /proc/cpuinfo | sed 's/^\s*//')"
printf "Memory: %s\n" "$(awk '/MemTotal/ {t=$2} /MemAvailable/ {a=$2} END {if (t>0) printf "%d%% (%dMB)", (t-a)/t*100, t/1024; else print "n/a"}' /proc/meminfo)"
printf "Disk: %s\n" "$(df -h / --output=source,pcent | sed -n '2p')"

exit 0
