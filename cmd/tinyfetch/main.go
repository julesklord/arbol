package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func runCommand(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getOSName() string {
	if runtime.GOOS == "darwin" {
		name := runCommand("sw_vers", "-productName")
		ver := runCommand("sw_vers", "-productVersion")
		if name != "" && ver != "" {
			return name + " " + ver
		}
		return "macOS"
	}
	if runtime.GOOS == "linux" {
		file, err := os.Open("/etc/os-release")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					val := strings.TrimPrefix(line, "PRETTY_NAME=")
					return strings.Trim(val, "\"")
				}
			}
		}
		return "Linux"
	}
	return runtime.GOOS
}

func getDistroID() string {
	if runtime.GOOS == "darwin" {
		return "darwin"
	}
	if runtime.GOOS == "linux" {
		file, err := os.Open("/etc/os-release")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "ID=") {
					val := strings.TrimPrefix(line, "ID=")
					return strings.Trim(val, "\"")
				}
			}
		}
	}
	return "linux"
}

func getUptime() string {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/uptime")
		if err == nil {
			parts := strings.Fields(string(data))
			if len(parts) > 0 {
				if sec, err := strconv.ParseFloat(parts[0], 64); err == nil {
					h := int(sec) / 3600
					m := (int(sec) % 3600) / 60
					return fmt.Sprintf("%dh %dm", h, m)
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		out := runCommand("sysctl", "-n", "kern.boottime")
		if out != "" {
			idx := strings.Index(out, "sec = ")
			if idx != -1 {
				s := out[idx+6:]
				comma := strings.Index(s, ",")
				if comma != -1 {
					secStr := strings.TrimSpace(s[:comma])
					if sec, err := strconv.ParseInt(secStr, 10, 64); err == nil {
						diff := time.Now().Unix() - sec
						h := diff / 3600
						m := (diff % 3600) / 60
						return fmt.Sprintf("%dh %dm", h, m)
					}
				}
			}
		}
	}
	// Generic fallback
	uptimeStr := runCommand("uptime")
	if uptimeStr != "" {
		// Very simple parser for uptime output: look for "up"
		idx := strings.Index(uptimeStr, "up ")
		if idx != -1 {
			s := uptimeStr[idx+3:]
			comma := strings.Index(s, ",")
			if comma != -1 {
				return strings.TrimSpace(s[:comma])
			}
		}
	}
	return "n/a"
}

func getCPU() string {
	if runtime.GOOS == "linux" {
		file, err := os.Open("/proc/cpuinfo")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "model name") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		brand := runCommand("sysctl", "-n", "machdep.cpu.brand_string")
		if brand != "" {
			return brand
		}
		model := runCommand("sysctl", "-n", "hw.model")
		if model != "" {
			return model
		}
	}
	return "Unknown CPU"
}

func getMemory() string {
	if runtime.GOOS == "linux" {
		file, err := os.Open("/proc/meminfo")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			var total, avail int64
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "MemTotal:") {
					fmt.Sscanf(line, "MemTotal: %d kB", &total)
				} else if strings.HasPrefix(line, "MemAvailable:") {
					fmt.Sscanf(line, "MemAvailable: %d kB", &avail)
				}
			}
			if total > 0 {
				usedPct := (total - avail) * 100 / total
				return fmt.Sprintf("%d%% (%dMB)", usedPct, total/1024)
			}
		}
	} else if runtime.GOOS == "darwin" {
		totalBytesStr := runCommand("sysctl", "-n", "hw.memsize")
		totalBytes, _ := strconv.ParseInt(totalBytesStr, 10, 64)
		totalMB := totalBytes / 1024 / 1024

		pageSizeStr := runCommand("bash", "-c", "vm_stat | awk '/page size of/ {print $8}' | tr -d '.'")
		pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 64)
		if pageSize == 0 {
			pageSize = 4096
		}

		freePagesStr := runCommand("bash", "-c", "vm_stat | awk '/Pages free:/ {print $3}' | tr -d '.'")
		freePages, _ := strconv.ParseInt(freePagesStr, 10, 64)

		inactivePagesStr := runCommand("bash", "-c", "vm_stat | awk '/Pages inactive:/ {print $3}' | tr -d '.'")
		inactivePages, _ := strconv.ParseInt(inactivePagesStr, 10, 64)

		if freePages > 0 && totalMB > 0 {
			freeMB := (freePages + inactivePages) * pageSize / 1024 / 1024
			usedMB := totalMB - freeMB
			pct := usedMB * 100 / totalMB
			return fmt.Sprintf("%d%% (%dMB)", pct, totalMB)
		}
		if totalMB > 0 {
			return fmt.Sprintf("n/a (%dMB)", totalMB)
		}
	}
	return "n/a"
}

func getDisk() string {
	out := runCommand("df", "-h", "/")
	if out != "" {
		lines := strings.Split(out, "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 5 {
				return fmt.Sprintf("%s (%s)", fields[0], fields[4])
			}
		}
	}
	return "n/a"
}

func getBar(pct int) string {
	filled := pct / 10
	if filled > 10 {
		filled = 10
	}
	empty := 10 - filled
	color := "\033[01;32m" // Green
	if pct > 80 {
		color = "\033[01;31m" // Red
	} else if pct > 50 {
		color = "\033[01;33m" // Yellow
	}
	restore := "\033[0m"
	gray := "\033[00;37m"

	var sb strings.Builder
	sb.WriteString(color)
	for i := 0; i < filled; i++ {
		sb.WriteString("█")
	}
	sb.WriteString(restore + gray)
	for i := 0; i < empty; i++ {
		sb.WriteString("░")
	}
	sb.WriteString(restore)
	return sb.String()
}

func stripANSI(s string) string {
	var builder strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		builder.WriteByte(s[i])
	}
	return builder.String()
}

func main() {
	noASCII := false
	for _, arg := range os.Args[1:] {
		if arg == "--no-ascii" {
			noASCII = true
		} else if arg == "--help" || arg == "-h" {
			fmt.Printf("Usage: %s [--no-ascii]\n", os.Args[0])
			os.Exit(0)
		}
	}

	hostname, _ := os.Hostname()
	osName := getOSName()
	kernel := runCommand("uname", "-r")
	uptimeVal := getUptime()
	shellVal := os.Getenv("SHELL")
	if shellVal == "" {
		shellVal = "sh"
	}
	cpuVal := getCPU()

	// Memory & Progress Bar
	memRaw := getMemory()
	memVal := memRaw
	if strings.Contains(memRaw, "%") {
		pctPart := strings.Split(memRaw, "%")[0]
		if pct, err := strconv.Atoi(strings.TrimSpace(pctPart)); err == nil {
			memVal = getBar(pct) + " " + memRaw
		}
	}

	// Disk & Progress Bar
	diskRaw := getDisk()
	diskVal := diskRaw
	if strings.Contains(diskRaw, "%") {
		idx := strings.Index(diskRaw, "%")
		start := idx
		for start > 0 && diskRaw[start-1] >= '0' && diskRaw[start-1] <= '9' {
			start--
		}
		if pctStr := diskRaw[start:idx]; pctStr != "" {
			if pct, err := strconv.Atoi(pctStr); err == nil {
				diskVal = getBar(pct) + " " + diskRaw
			}
		}
	}

	// Colors
	restore := "\033[0m"
	lblue := "\033[01;34m"
	lyellow := "\033[01;33m"
	lcyan := "\033[01;36m"
	white := "\033[01;37m"

	// Setup Logo
	var logo []string
	if !noASCII {
		distroID := getDistroID()
		// Paths to search
		searchPaths := []string{
			"./ascii/" + distroID + ".txt",
			"/usr/local/share/tinyfetch/ascii/" + distroID + ".txt",
			"/usr/share/tinyfetch/ascii/" + distroID + ".txt",
		}

		asciiPath := ""
		for _, path := range searchPaths {
			if _, err := os.Stat(path); err == nil {
				asciiPath = path
				break
			}
		}

		// Fallback to generic if not found
		if asciiPath == "" {
			fallback := "linux"
			if runtime.GOOS == "darwin" {
				fallback = "darwin"
			}
			fallbackPaths := []string{
				"./ascii/" + fallback + ".txt",
				"/usr/local/share/tinyfetch/ascii/" + fallback + ".txt",
				"/usr/share/tinyfetch/ascii/" + fallback + ".txt",
			}
			for _, path := range fallbackPaths {
				if _, err := os.Stat(path); err == nil {
					asciiPath = path
					break
				}
			}
		}

		if asciiPath != "" {
			file, err := os.Open(asciiPath)
			if err == nil {
				defer file.Close()
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					logo = append(logo, scanner.Text())
				}
			}
		}

		// Hardcoded fallbacks if no file is available
		if len(logo) == 0 {
			if runtime.GOOS == "darwin" {
				logo = []string{
					lcyan + "      .---." + restore,
					lcyan + "     /     \\" + restore,
					lcyan + "     \\__   /" + restore,
					lcyan + "    /   `-' \\" + restore,
					lcyan + "   |         |" + restore,
					lcyan + "    \\       /" + restore,
					lcyan + "     `-...-'" + restore,
				}
			} else {
				logo = []string{
					lyellow + "     .---." + restore,
					lyellow + "    /     \\" + restore,
					lblue + "    \\ " + restore + white + "o o" + restore + lblue + " /" + restore,
					lyellow + "    /  \\-/ \\" + restore,
					lyellow + "   / /     \\ \\" + restore,
					lyellow + "  ( (_     _ ) )" + restore,
					lyellow + "   `(_`---'_)''" + restore,
				}
			}
		}
	}

	// Setup Info
	info := []string{
		lblue + "Host:" + restore + "   " + hostname,
		lblue + "OS:" + restore + "     " + osName,
		lblue + "Kernel:" + restore + " " + kernel,
		lblue + "Uptime:" + restore + " " + uptimeVal,
		lblue + "Shell:" + restore + "  " + shellVal,
		lblue + "CPU:" + restore + "    " + cpuVal,
		lblue + "Memory:" + restore + " " + memVal,
		lblue + "Disk:" + restore + "   " + diskVal,
	}

	// Scan ./plugins directory
	if entries, err := os.ReadDir("./plugins"); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				infoPath := "./plugins/" + entry.Name()
				fileInfo, err := entry.Info()
				if err == nil && (fileInfo.Mode()&0111 != 0) {
					out := runCommand(infoPath)
					if out != "" {
						lines := strings.Split(out, "\n")
						pluginOut := strings.TrimSpace(lines[0])
						if pluginOut != "" {
							if strings.Contains(pluginOut, ":") {
								parts := strings.SplitN(pluginOut, ":", 2)
								info = append(info, lblue+parts[0]+":"+restore+" "+strings.TrimSpace(parts[1]))
							} else {
								name := entry.Name()
								if idx := strings.Index(name, "."); idx != -1 {
									name = name[:idx]
								}
								if len(name) > 0 {
									name = strings.ToUpper(name[:1]) + name[1:]
								}
								info = append(info, lblue+name+":"+restore+" "+pluginOut)
							}
						}
					}
				}
			}
		}
	}

	maxLines := len(info)
	if !noASCII && len(logo) > maxLines {
		maxLines = len(logo)
	}

	// Calculate maximum logo raw length
	leftW := 0
	if !noASCII {
		for _, line := range logo {
			raw := stripANSI(line)
			rawLen := utf8.RuneCountInString(raw)
			if rawLen > leftW {
				leftW = rawLen
			}
		}
		if leftW < 16 {
			leftW = 16
		}
	}

	// Calculate maximum info raw length
	rightW := 0
	for _, line := range info {
		raw := stripANSI(line)
		rawLen := utf8.RuneCountInString(raw)
		if rawLen > rightW {
			rightW = rawLen
		}
	}

	borderCol := lblue

	if noASCII {
		topLine := borderCol + "┌" + strings.Repeat("─", rightW+2) + "┐" + restore
		botLine := borderCol + "└" + strings.Repeat("─", rightW+2) + "┘" + restore
		fmt.Println(topLine)
		for _, line := range info {
			rawLen := utf8.RuneCountInString(stripANSI(line))
			padCount := rightW - rawLen
			padding := ""
			if padCount > 0 {
				padding = strings.Repeat(" ", padCount)
			}
			fmt.Printf("%s│%s %s%s %s│\n", borderCol, restore, line, padding, borderCol)
		}
		fmt.Println(botLine)
	} else {
		topLine := borderCol + "┌" + strings.Repeat("─", leftW+2) + "┬" + strings.Repeat("─", rightW+2) + "┐" + restore
		botLine := borderCol + "└" + strings.Repeat("─", leftW+2) + "┴" + strings.Repeat("─", rightW+2) + "┘" + restore
		fmt.Println(topLine)
		for i := 0; i < maxLines; i++ {
			logoPrint := ""
			if i < len(logo) {
				logoPrint = logo[i]
			}
			lRaw := utf8.RuneCountInString(stripANSI(logoPrint))
			lPadCount := leftW - lRaw
			lPadding := ""
			if lPadCount > 0 {
				lPadding = strings.Repeat(" ", lPadCount)
			}

			infoPrint := ""
			if i < len(info) {
				infoPrint = info[i]
			}
			rRaw := utf8.RuneCountInString(stripANSI(infoPrint))
			rPadCount := rightW - rRaw
			rPadding := ""
			if rPadCount > 0 {
				rPadding = strings.Repeat(" ", rPadCount)
			}

			fmt.Printf("%s│%s %s%s %s│%s %s%s %s│\n",
				borderCol, restore, logoPrint, lPadding,
				borderCol, restore, infoPrint, rPadding,
				borderCol)
		}
		fmt.Println(botLine)
	}
}
