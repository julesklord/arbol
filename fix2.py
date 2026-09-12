import sys

with open("cmd/arbol/sysinfo.go", "r") as f:
    text = f.read()

# Only fix the duplicated imports/vars in HEAD, don't remove existing needed vars
text = text.replace('import (\n\t"sync"\n\t"bufio"', 'import (\n\t"bufio"')

bad_var_block = """var (
	osNameCache   string
	distroIDCache string
	cpuCache      string
	gpuCache      string
	osReleaseOnce sync.Once
	cpuOnce       sync.Once
	gpuOnce       sync.Once
)"""

good_var_block = """var (
	osNameCache        string
	distroIDCache      string
	cpuCache           string
	gpuCache           string
	gpuOnce            sync.Once
	darwinTotalMemMB   int64
	darwinTotalMemOnce sync.Once
)"""

text = text.replace(bad_var_block, good_var_block)

sysinfo_func = """func getMemory() string {"""
sysinfo_func_replacement = """func getDarwinTotalMemMB() int64 {
	darwinTotalMemOnce.Do(func() {
		totalBytesStr := runCommand("sysctl", "-n", "hw.memsize")
		totalBytes, _ := strconv.ParseInt(totalBytesStr, 10, 64)
		darwinTotalMemMB = totalBytes / 1024 / 1024
	})
	return darwinTotalMemMB
}

func getMemory() string {"""
text = text.replace(sysinfo_func, sysinfo_func_replacement, 1) # Only replace FIRST instance, sysinfo.go might have it twice by accident earlier

sysinfo_mem_old = """	} else if runtime.GOOS == "darwin" {
		totalBytesStr := runCommand("sysctl", "-n", "hw.memsize")
		totalBytes, _ := strconv.ParseInt(totalBytesStr, 10, 64)
		totalMB := totalBytes / 1024 / 1024

		// OPTIMIZATION: Consolidate multiple subprocesses into a single native call"""
sysinfo_mem_new = """	} else if runtime.GOOS == "darwin" {
		totalMB := getDarwinTotalMemMB()

		// OPTIMIZATION: Consolidate multiple subprocesses into a single native call"""
text = text.replace(sysinfo_mem_old, sysinfo_mem_new)

with open("cmd/arbol/sysinfo.go", "w") as f:
    f.write(text)
