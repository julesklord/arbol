package main

import (
	"syscall"

	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SparklineBuffer holds a rolling window of values for sparkline rendering
type SparklineBuffer struct {
	mu       sync.Mutex
	values   []int
	maxLen   int
	interval time.Duration
	stopCh   chan struct{}
}

// Sparkline characters (braille patterns for 8 levels)
// Each braille char = 2x4 dots, giving us 8 vertical levels
// Using bottom-to-top: ⠁⠂⠄⡀⢀⠠⠐⠈ (but we want filled bottom)
var sparklineChars = []string{
	" ", // 0 - empty
	"▁", // 1 - bottom 1/8
	"▂", // 2 - bottom 2/8
	"▃", // 3 - bottom 3/8
	"▄", // 4 - bottom 4/8
	"▅", // 5 - bottom 5/8
	"▆", // 6 - bottom 6/8
	"▇", // 7 - bottom 7/8
	"█", // 8 - full
}

// Braille sparkline (higher resolution, 2x4 dots per char)
var brailleSparkline = []string{
	"⠀", // 0/8
	"⠁", // 1/8 (dot 1)
	"⠃", // 2/8 (dots 1,2)
	"⠇", // 3/8 (dots 1,2,3)
	"⠏", // 4/8 (dots 1,2,3,4)
	"⠟", // 5/8 (dots 1,2,3,4,5)
	"⠿", // 6/8 (dots 1-6)
	"⠿", // 7/8
	"⠿", // 8/8
}

var (
	cpuSparkline   *SparklineBuffer
	memSparkline   *SparklineBuffer
	swapSparkline  *SparklineBuffer
	diskSparkline  *SparklineBuffer
	sparklineWidth = 20 // default width in characters
)

func initSparklines(width int, interval time.Duration) {
	sparklineWidth = width
	cpuSparkline = NewSparklineBuffer(width, interval)
	memSparkline = NewSparklineBuffer(width, interval)
	swapSparkline = NewSparklineBuffer(width, interval)
	diskSparkline = NewSparklineBuffer(width, interval)

	// Start background collectors
	cpuSparkline.Start(collectCPUPercent)
	memSparkline.Start(collectMemPercent)
	swapSparkline.Start(collectSwapPercent)
	diskSparkline.Start(collectDiskPercent)
}

func NewSparklineBuffer(maxLen int, interval time.Duration) *SparklineBuffer {
	return &SparklineBuffer{
		values:   make([]int, 0, maxLen),
		maxLen:   maxLen,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (s *SparklineBuffer) Start(collector func() int) {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				val := collector()
				s.Add(val)
			case <-s.stopCh:
				return
			}
		}
	}()
}

func (s *SparklineBuffer) Stop() {
	close(s.stopCh)
}

func (s *SparklineBuffer) Add(val int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if val < 0 {
		val = 0
	}
	if val > 100 {
		val = 100
	}
	s.values = append(s.values, val)
	if len(s.values) > s.maxLen {
		s.values = s.values[len(s.values)-s.maxLen:]
	}
}

func (s *SparklineBuffer) Get() []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]int, len(s.values))
	copy(result, s.values)
	return result
}

func (s *SparklineBuffer) Render(style SparklineStyle, color, gray, reset string) string {
	values := s.Get()
	if len(values) == 0 {
		if ColorDisabled {
			return strings.Repeat(" ", sparklineWidth)
		}
		return gray + strings.Repeat(" ", sparklineWidth) + reset
	}

	var sb strings.Builder
	if !ColorDisabled {
		sb.WriteString(color)
	}

	switch style {
	case SparklineBlock:
		for _, v := range values {
			idx := v * 8 / 100
			if idx > 8 {
				idx = 8
			}
			sb.WriteString(sparklineChars[idx])
		}
	case SparklineBraille:
		// Braille gives 2x4 = 8 levels per char, but we render horizontally
		// Each braille char column = 2 dots wide, 4 dots tall
		// We'll use a simpler approach: map 0-100 to 0-8 braille patterns
		for _, v := range values {
			idx := v * 8 / 100
			if idx > 8 {
				idx = 8
			}
			sb.WriteString(brailleSparkline[idx])
		}
	case SparklineDots:
		for _, v := range values {
			if v > 50 {
				sb.WriteString("●")
			} else {
				sb.WriteString("○")
			}
		}
	}

	if !ColorDisabled {
		sb.WriteString(reset)
	}
	return sb.String()
}

type SparklineStyle int

const (
	SparklineBlock SparklineStyle = iota
	SparklineBraille
	SparklineDots
)

var currentSparklineStyle = SparklineBlock

func SetSparklineStyle(style SparklineStyle) {
	currentSparklineStyle = style
}

func GetSparklineStyle() SparklineStyle {
	return currentSparklineStyle
}

func collectCPUPercent() int {
	if runtime.GOOS == "linux" {
		u1, n1, s1, id1, io1, ir1, so1, err1 := getCPUTicks()
		if err1 != nil {
			return 0
		}
		time.Sleep(50 * time.Millisecond)
		u2, n2, s2, id2, io2, ir2, so2, err2 := getCPUTicks()
		if err2 != nil {
			return 0
		}

		idle1 := id1 + io1
		idle2 := id2 + io2

		nonIdle1 := u1 + n1 + s1 + ir1 + so1
		nonIdle2 := u2 + n2 + s2 + ir2 + so2

		total1 := idle1 + nonIdle1
		total2 := idle2 + nonIdle2

		totalDiff := total2 - total1
		idleDiff := idle2 - idle1

		if totalDiff > 0 {
			pct := (totalDiff - idleDiff) * 100 / totalDiff
			return int(pct)
		}
	} else if runtime.GOOS == "darwin" {
		out := runCommand("bash", "-c", "ps -A -o %cpu | awk '{s+=$1} END {print s}'")
		if out != "" {
			if val, err := strconv.ParseFloat(out, 64); err == nil {
				cores := runtime.NumCPU()
				if cores > 0 {
					pct := int(val / float64(cores))
					if pct > 100 {
						pct = 100
					}
					return pct
				}
			}
		}
	}
	return 0
}

func collectMemPercent() int {
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
				return int(usedPct)
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
			return int(pct)
		}
	}
	return 0
}

func collectSwapPercent() int {
	if runtime.GOOS == "linux" {
		file, err := os.Open("/proc/meminfo")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			var total, free int64
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "SwapTotal:") {
					fmt.Sscanf(line, "SwapTotal: %d kB", &total)
				} else if strings.HasPrefix(line, "SwapFree:") {
					fmt.Sscanf(line, "SwapFree: %d kB", &free)
				}
			}
			if total > 0 {
				used := total - free
				pct := used * 100 / total
				return int(pct)
			}
		}
	}
	return 0
}

func collectDiskPercent() int {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err == nil {
		total := stat.Blocks
		free := stat.Bavail
		if total > 0 {
			used := total - free
			return int((used * 100) / total)
		}
	}

	out := runCommand("df", "-Ph", "/")
	if out != "" {
		// Optimize memory by limiting split to 3 items since we only need the 2nd line
		lines := strings.SplitN(out, "\n", 3)
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 5 {
				pctStr := strings.TrimSuffix(fields[4], "%")
				if pct, err := strconv.Atoi(pctStr); err == nil {
					return pct
				}
			}
		}
	}
	return 0
}

func getSparklineCPU(theme Theme) string {
	if cpuSparkline == nil {
		return ""
	}
	return cpuSparkline.Render(currentSparklineStyle, theme.Success, theme.Muted, "\033[0m")
}

func getSparklineMem(theme Theme) string {
	if memSparkline == nil {
		return ""
	}
	return memSparkline.Render(currentSparklineStyle, theme.Primary, theme.Muted, "\033[0m")
}

func getSparklineSwap(theme Theme) string {
	if swapSparkline == nil {
		return ""
	}
	return swapSparkline.Render(currentSparklineStyle, theme.Secondary, theme.Muted, "\033[0m")
}

func getSparklineDisk(theme Theme) string {
	if diskSparkline == nil {
		return ""
	}
	return diskSparkline.Render(currentSparklineStyle, theme.Warning, theme.Muted, "\033[0m")
}
