package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func runCommand(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func runCommandWithTimeout(timeout time.Duration, name string, arg ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, arg...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getTerminalWidth() int {
	out := runCommand("tput", "cols")
	if out != "" {
		if w, err := strconv.Atoi(out); err == nil {
			return w
		}
	}
	return 80
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
	out := runCommand("df", "-Ph", "/")
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
