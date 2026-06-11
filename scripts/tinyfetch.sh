#!/usr/bin/env bash
# Robust and portable tinyfetch script
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

# Detection of OS
OS_TYPE=$(uname -s)

# Helper functions for portable resource gathering
get_os_name() {
  if [ "$OS_TYPE" = "Darwin" ]; then
    if command -v sw_vers >/dev/null 2>&1; then
      printf "%s %s" "$(sw_vers -productName)" "$(sw_vers -productVersion)"
    else
      echo "macOS"
    fi
  else
    if [ -f /etc/os-release ]; then
      grep '^PRETTY_NAME' /etc/os-release | cut -d= -f2 | tr -d '"'
    else
      echo "$OS_TYPE"
    fi
  fi
}

get_uptime() {
  if [ -f /proc/uptime ]; then
    local uptime_sec
    uptime_sec=$(cut -d. -f1 /proc/uptime)
    local h=$((uptime_sec / 3600))
    local m=$(((uptime_sec % 3600) / 60))
    echo "${h}h ${m}m"
  elif [ "$OS_TYPE" = "Darwin" ] && command -v sysctl >/dev/null 2>&1; then
    local boot_time
    boot_time=$(sysctl -n kern.boottime 2>/dev/null | awk -F'[=,]' '{print $2}' | tr -d ' ')
    if [ -n "$boot_time" ]; then
      local now
      now=$(date +%s)
      local diff=$((now - boot_time))
      local h=$((diff / 3600))
      local m=$(((diff % 3600) / 60))
      echo "${h}h ${m}m"
    else
      uptime | sed -E 's/^.*up[[:space:]]+([^,]+),.*$/\1/'
    fi
  else
    uptime | sed -E 's/^.*up[[:space:]]+([^,]+),.*$/\1/'
  fi
}

get_cpu() {
  if [ -f /proc/cpuinfo ]; then
    awk -F: '/model name/{print $2; exit}' /proc/cpuinfo | sed 's/^\s*//'
  elif [ "$OS_TYPE" = "Darwin" ] && command -v sysctl >/dev/null 2>&1; then
    sysctl -n machdep.cpu.brand_string 2>/dev/null || sysctl -n hw.model 2>/dev/null || echo "Unknown CPU"
  else
    echo "Unknown CPU"
  fi
}

get_memory() {
  if [ -f /proc/meminfo ]; then
    awk '/MemTotal/ {t=$2} /MemAvailable/ {a=$2} END {if (t>0) printf "%d%% (%dMB)", (t-a)/t*100, t/1024; else print "n/a"}' /proc/meminfo
  elif [ "$OS_TYPE" = "Darwin" ] && command -v sysctl >/dev/null 2>&1 && command -v vm_stat >/dev/null 2>&1; then
    local total_bytes
    total_bytes=$(sysctl -n hw.memsize 2>/dev/null)
    local total_mb=$((total_bytes / 1024 / 1024))
    
    local page_size
    page_size=$(vm_stat | awk '/page size of/ {print $8}' | tr -d '.')
    [ -z "$page_size" ] && page_size=4096
    
    local free_pages
    free_pages=$(vm_stat | awk '/Pages free:/ {print $3}' | tr -d '.')
    local inactive_pages
    inactive_pages=$(vm_stat | awk '/Pages inactive:/ {print $3}' | tr -d '.')
    
    if [ -n "$free_pages" ] && [ -n "$inactive_pages" ]; then
      local free_mb=$(((free_pages + inactive_pages) * page_size / 1024 / 1024))
      local used_mb=$((total_mb - free_mb))
      local pct=$((used_mb * 100 / total_mb))
      echo "${pct}% (${total_mb}MB)"
    else
      echo "n/a (${total_mb}MB)"
    fi
  else
    echo "n/a"
  fi
}

# Colors
ESC=$(printf '\033')
RESTORE="${ESC}[0m"
LBLUE="${ESC}[01;34m"
LYELLOW="${ESC}[01;33m"
LCYAN="${ESC}[01;36m"
WHITE="${ESC}[01;37m"
LRED="${ESC}[01;31m"
LGREEN="${ESC}[01;32m"
LIGHTGRAY="${ESC}[00;37m"

# Progress Bar Helper
get_bar() {
  local pct=$1
  local filled=$((pct / 10))
  [ $filled -gt 10 ] && filled=10
  local empty=$((10 - filled))
  local bar=""
  
  local color="$LGREEN"
  if [ "$pct" -gt 80 ]; then
    color="$LRED"
  elif [ "$pct" -gt 50 ]; then
    color="$LYELLOW"
  fi
  
  bar="${color}"
  for ((i=0; i<filled; i++)); do bar="${bar}█"; done
  bar="${bar}${RESTORE}${LIGHTGRAY}"
  for ((i=0; i<empty; i++)); do bar="${bar}░"; done
  bar="${bar}${RESTORE}"
  echo "$bar"
}

# Resolve values safely
HOST=$(hostname)
OS_NAME=$(get_os_name)
KERNEL=$(uname -r)
UPTIME=$(get_uptime)
SHELL_VAL="${SHELL-sh}"
CPU=$(get_cpu)

# Memory with visual bar
MEM_RAW=$(get_memory)
if [[ "$MEM_RAW" == *"%"* ]]; then
  MEM_PCT=$(echo "$MEM_RAW" | cut -d% -f1)
  MEM_BAR=$(get_bar "$MEM_PCT")
  MEMORY="${MEM_BAR} ${MEM_RAW}"
else
  MEMORY="$MEM_RAW"
fi

# Disk with visual bar
DISK_RAW=$(df -h / | awk 'NR==2 {print $1 " (" $5 ")"}')
DISK_PCT=$(echo "$DISK_RAW" | grep -o '[0-9]\+%' | tr -d '%' || echo "0")
DISK_BAR=$(get_bar "$DISK_PCT")
DISK="${DISK_BAR} ${DISK_RAW}"

# Get Distro ID
get_distro_id() {
  if [ "$OS_TYPE" = "Darwin" ]; then
    echo "darwin"
  elif [ -f /etc/os-release ]; then
    local id
    id=$(grep '^ID=' /etc/os-release | cut -d= -f2 | tr -d '"')
    echo "${id:-linux}"
  else
    echo "linux"
  fi
}

DISTRO_ID=$(get_distro_id)

# Find ASCII file path
ASCII_FILE=""
for path in "./ascii/${DISTRO_ID}.txt" "/usr/local/share/tinyfetch/ascii/${DISTRO_ID}.txt" "/usr/share/tinyfetch/ascii/${DISTRO_ID}.txt"; do
  if [ -f "$path" ]; then
    ASCII_FILE="$path"
    break
  fi
done

# If not found, try fallback linux.txt or darwin.txt
if [ -z "$ASCII_FILE" ]; then
  fallback_name="linux"
  if [ "$OS_TYPE" = "Darwin" ]; then
    fallback_name="darwin"
  fi
  for path in "./ascii/${fallback_name}.txt" "/usr/local/share/tinyfetch/ascii/${fallback_name}.txt" "/usr/share/tinyfetch/ascii/${fallback_name}.txt"; do
    if [ -f "$path" ]; then
      ASCII_FILE="$path"
      break
    fi
  done
fi

logo=()
if [ "$NO_ASCII" -eq 0 ]; then
  if [ -n "$ASCII_FILE" ]; then
    while IFS= read -r line || [ -n "$line" ]; do
      logo+=("$line")
    done < "$ASCII_FILE"
  else
    if [ "$OS_TYPE" = "Darwin" ]; then
      logo[0]="${LCYAN}      .---.${RESTORE}"
      logo[1]="${LCYAN}     /     \\${RESTORE}"
      logo[2]="${LCYAN}     \\__   /${RESTORE}"
      logo[3]="${LCYAN}    /   \`-' \\${RESTORE}"
      logo[4]="${LCYAN}   |         |${RESTORE}"
      logo[5]="${LCYAN}    \\       /${RESTORE}"
      logo[6]="${LCYAN}     \`-...-'${RESTORE}"
      logo[7]=""
    else
      logo[0]="${LYELLOW}     .---.${RESTORE}"
      logo[1]="${LYELLOW}    /     \\${RESTORE}"
      logo[2]="${LBLUE}    \\ ${RESTORE}${WHITE}o o${RESTORE}${LBLUE} /${RESTORE}"
      logo[3]="${LYELLOW}    /  \\-/ \\${RESTORE}"
      logo[4]="${LYELLOW}   / /     \\ \\${RESTORE}"
      logo[5]="${LYELLOW}  ( (_     _ ) )${RESTORE}"
      logo[6]="${LYELLOW}   \`(_\`---'_)''${RESTORE}"
      logo[7]=""
    fi
  fi
fi

info=()
info[0]="${LBLUE}Host:${RESTORE}   $HOST"
info[1]="${LBLUE}OS:${RESTORE}     $OS_NAME"
info[2]="${LBLUE}Kernel:${RESTORE} $KERNEL"
info[3]="${LBLUE}Uptime:${RESTORE} $UPTIME"
info[4]="${LBLUE}Shell:${RESTORE}  $SHELL_VAL"
info[5]="${LBLUE}CPU:${RESTORE}    $CPU"
info[6]="${LBLUE}Memory:${RESTORE} $MEMORY"
info[7]="${LBLUE}Disk:${RESTORE}   $DISK"

# Scan ./plugins directory
if [ -d "./plugins" ]; then
  # Enable nullglob to avoid executing literally `./plugins/*` if folder is empty
  shopt -s nullglob
  for p in ./plugins/*; do
    if [ -x "$p" ] && [ -f "$p" ]; then
      plugin_out=$("$p" 2>/dev/null | head -n 1)
      if [ -n "$plugin_out" ]; then
        if [[ "$plugin_out" == *":"* ]]; then
          p_key=$(echo "$plugin_out" | cut -d: -f1)
          p_val=$(echo "$plugin_out" | cut -d: -f2- | sed 's/^\s*//')
          info+=("${LBLUE}${p_key}:${RESTORE} $p_val")
        else
          label=$(basename "$p" | cut -d. -f1 | sed 's/^[0-9]\+-//')
          label="$(tr '[:lower:]' '[:upper:]' <<< "${label:0:1}")${label:1}"
          info+=("${LBLUE}${label}:${RESTORE} $plugin_out")
        fi
      fi
    fi
  done
  shopt -u nullglob
fi

# Scan ./plugins/extended directory
ext_info=()
HAS_EXT=0
if [ -d "./plugins/extended" ]; then
  shopt -s nullglob
  for p in ./plugins/extended/*; do
    if [ -x "$p" ] && [ -f "$p" ]; then
      has_content=0
      # Use temporary array to hold this plugin's output
      tmp_out=()
      while IFS= read -r line || [ -n "$line" ]; do
        tmp_out+=("$line")
        has_content=1
      done < <("$p" 2>/dev/null)
      
      if [ "$has_content" -eq 1 ]; then
        for line in "${tmp_out[@]}"; do
          ext_info+=("$line")
        done
        ext_info+=("") # separator
        HAS_EXT=1
      fi
    fi
  done
  shopt -u nullglob
  # Remove trailing empty line separator if present
  if [ ${#ext_info[@]} -gt 0 ]; then
    last_idx=$((${#ext_info[@]} - 1))
    if [ -z "${ext_info[$last_idx]}" ]; then
      unset "ext_info[$last_idx]"
      # Re-index array after unset to avoid sparse index issues
      ext_info=("${ext_info[@]}")
    fi
  fi
fi

# Print with novel card layout
max_lines=${#info[@]}
if [ "$NO_ASCII" -eq 0 ] && [ ${#logo[@]} -gt "$max_lines" ]; then
  max_lines=${#logo[@]}
fi
if [ "$HAS_EXT" -eq 1 ] && [ ${#ext_info[@]} -gt "$max_lines" ]; then
  max_lines=${#ext_info[@]}
fi

# Calculate maximum logo raw length
left_w=0
if [ "$NO_ASCII" -eq 0 ]; then
  for line in "${logo[@]}"; do
    raw=$(printf "%s" "$line" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
    if [ ${#raw} -gt $left_w ]; then
      left_w=${#raw}
    fi
  done
  [ $left_w -lt 16 ] && left_w=16
fi

# Calculate maximum info raw length
right_w=0
for line in "${info[@]}"; do
  raw=$(printf "%s" "$line" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
  if [ ${#raw} -gt $right_w ]; then
    right_w=${#raw}
  fi
done

# Calculate maximum extended info raw length
ext_w=0
if [ "$HAS_EXT" -eq 1 ]; then
  for line in "${ext_info[@]}"; do
    raw=$(printf "%s" "$line" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g')
    if [ ${#raw} -gt $ext_w ]; then
      ext_w=${#raw}
    fi
  done
  [ $ext_w -lt 24 ] && ext_w=24
fi

strip_ansi() {
  printf "%s" "$1" | sed 's/\x1b\[[0-9;]*[a-zA-Z]//g'
}

repeat_char() {
  local char="$1"
  local count="$2"
  local out=""
  for ((k=0; k<count; k++)); do
    out="${out}${char}"
  done
  echo -n "$out"
}

BORDER_COLOR="$LBLUE"

if [ "$HAS_EXT" -eq 0 ]; then
  if [ "$NO_ASCII" -eq 1 ]; then
    # Case 1: Single pane (Info)
    top_line="${BORDER_COLOR}┌$(repeat_char "─" $((right_w + 2)))┐${RESTORE}"
    bot_line="${BORDER_COLOR}└$(repeat_char "─" $((right_w + 2)))┘${RESTORE}"
    echo -e "$top_line"
    for ((i=0; i<max_lines; i++)); do
      r_line="${info[i]:-}"
      r_raw=$(strip_ansi "$r_line")
      r_pad=$((right_w - ${#r_raw}))
      r_padding=""
      [ $r_pad -gt 0 ] && r_padding=$(printf "%${r_pad}s" "")
      echo -e "${BORDER_COLOR}│${RESTORE} ${r_line}${r_padding} ${BORDER_COLOR}│"
    done
    echo -e "$bot_line"
  else
    # Case 2: Double pane (Logo + Info)
    top_line="${BORDER_COLOR}┌$(repeat_char "─" $((left_w + 2)))┬$(repeat_char "─" $((right_w + 2)))┐${RESTORE}"
    bot_line="${BORDER_COLOR}└$(repeat_char "─" $((left_w + 2)))┴$(repeat_char "─" $((right_w + 2)))┘${RESTORE}"
    echo -e "$top_line"
    for ((i=0; i<max_lines; i++)); do
      l_line="${logo[i]:-}"
      r_line="${info[i]:-}"
      l_raw=$(strip_ansi "$l_line")
      l_pad=$((left_w - ${#l_raw}))
      l_padding=""
      [ $l_pad -gt 0 ] && l_padding=$(printf "%${l_pad}s" "")
      r_raw=$(strip_ansi "$r_line")
      r_pad=$((right_w - ${#r_raw}))
      r_padding=""
      [ $r_pad -gt 0 ] && r_padding=$(printf "%${r_pad}s" "")
      echo -e "${BORDER_COLOR}│${RESTORE} ${l_line}${l_padding} ${BORDER_COLOR}│${RESTORE} ${r_line}${r_padding} ${BORDER_COLOR}│"
    done
    echo -e "$bot_line"
  fi
else
  if [ "$NO_ASCII" -eq 1 ]; then
    # Case 3: Double pane (Info + Extended)
    top_line="${BORDER_COLOR}┌$(repeat_char "─" $((right_w + 2)))┬$(repeat_char "─" $((ext_w + 2)))┐${RESTORE}"
    bot_line="${BORDER_COLOR}└$(repeat_char "─" $((right_w + 2)))┴$(repeat_char "─" $((ext_w + 2)))┘${RESTORE}"
    echo -e "$top_line"
    for ((i=0; i<max_lines; i++)); do
      r_line="${info[i]:-}"
      e_line="${ext_info[i]:-}"
      r_raw=$(strip_ansi "$r_line")
      r_pad=$((right_w - ${#r_raw}))
      r_padding=""
      [ $r_pad -gt 0 ] && r_padding=$(printf "%${r_pad}s" "")
      e_raw=$(strip_ansi "$e_line")
      e_pad=$((ext_w - ${#e_raw}))
      e_padding=""
      [ $e_pad -gt 0 ] && e_padding=$(printf "%${e_pad}s" "")
      echo -e "${BORDER_COLOR}│${RESTORE} ${r_line}${r_padding} ${BORDER_COLOR}│${RESTORE} ${e_line}${e_padding} ${BORDER_COLOR}│"
    done
    echo -e "$bot_line"
  else
    # Case 4: Triple pane (Logo + Info + Extended)
    top_line="${BORDER_COLOR}┌$(repeat_char "─" $((left_w + 2)))┬$(repeat_char "─" $((right_w + 2)))┬$(repeat_char "─" $((ext_w + 2)))┐${RESTORE}"
    bot_line="${BORDER_COLOR}└$(repeat_char "─" $((left_w + 2)))┴$(repeat_char "─" $((right_w + 2)))┴$(repeat_char "─" $((ext_w + 2)))┘${RESTORE}"
    echo -e "$top_line"
    for ((i=0; i<max_lines; i++)); do
      l_line="${logo[i]:-}"
      r_line="${info[i]:-}"
      e_line="${ext_info[i]:-}"
      l_raw=$(strip_ansi "$l_line")
      l_pad=$((left_w - ${#l_raw}))
      l_padding=""
      [ $l_pad -gt 0 ] && l_padding=$(printf "%${l_pad}s" "")
      r_raw=$(strip_ansi "$r_line")
      r_pad=$((right_w - ${#r_raw}))
      r_padding=""
      [ $r_pad -gt 0 ] && r_padding=$(printf "%${r_pad}s" "")
      e_raw=$(strip_ansi "$e_line")
      e_pad=$((ext_w - ${#e_raw}))
      e_padding=""
      [ $e_pad -gt 0 ] && e_padding=$(printf "%${e_pad}s" "")
      echo -e "${BORDER_COLOR}│${RESTORE} ${l_line}${l_padding} ${BORDER_COLOR}│${RESTORE} ${r_line}${r_padding} ${BORDER_COLOR}│${RESTORE} ${e_line}${e_padding} ${BORDER_COLOR}│"
    done
    echo -e "$bot_line"
  fi
fi

exit 0

