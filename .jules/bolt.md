## 2024-07-05 - Kernel Version Retrieval Performance
**Learning:** Shelling out to `uname -r` via `exec.Command` in `runCommand` takes ~2.5ms per call. While this only happens once per run, using the `syscall.Uname` natively available in Go takes less than ~0.001ms (a 2500x speedup). Calling external commands in a compiled fastfetch alternative is a known bottleneck that should be replaced by native syscalls where possible.
**Action:** Replace `runCommand("uname", "-r")` with native `syscall.Uname` on Linux to improve start-up performance and reduce unnecessary sub-processes.

## 2024-08-11 - Disk Usage Retrieval Performance
**Learning:** Shelling out to `df -Ph /` via `exec.Command` inside `getDisk` and `collectDiskPercent` takes ~3.5ms per call. In contrast, using the native `syscall.Statfs` takes ~2.5µs. This >1000x speedup drastically reduces startup latency and CPU overhead, especially during high-frequency live mode updates where `collectDiskPercent` is called continuously.
**Action:** Replace shell outs to `df` with `syscall.Statfs` to calculate disk percentage usage natively, and parse `/proc/mounts` on Linux to get the filesystem name instead of spawning a new process.

## 2026-07-20 - Uptime and Swap Retrieval Performance
**Learning:** Parsing `/proc/uptime` and `/proc/meminfo` involves file I/O and string allocations, which takes ~15-30µs per call. Using native `syscall.Sysinfo` accesses system metrics directly and executes in ~1µs (15x-30x speedup), significantly reducing overhead especially during frequent live updates.
**Action:** Replaced `/proc/uptime` and `/proc/meminfo` file parsing with `syscall.Sysinfo` in `getUptime`, `getSwap` and `collectSwapPercent` for Linux to improve performance.

## 2024-07-26 - Sparkline and getCPU/getMem Performance on macOS (Darwin)
**Learning:** Shelling out to `bash -c` with `awk` and `tr` pipelines (like `ps -A -o %cpu | awk ...` or `vm_stat | awk ...`) incurs immense subprocess overhead, taking roughly ~50ms+ per invocation. Given that `collectCPUPercent` and `collectMemPercent` are called on a tight loop for sparklines (e.g. 1-second interval), this eats a lot of unnecessary CPU and memory.
**Action:** Replace shell pipelines involving `bash -c` with direct `exec.Command` calls (e.g. `ps -A -o %cpu` and `vm_stat`) and process their output natively in Go using `strings.Split` and `strconv.ParseFloat`. This improves cross-platform speed significantly and prevents blocking the main thread during live updates.

## 2024-05-24 - sysinfo.go getCPUTicks Optimization
**Learning:** Found an opportunity to optimize `/proc/stat` parsing in Go by eliminating the use of `strings.Fields()` which allocates a slice and strings for each token, and instead parsing integers in a single pass over the string buffer using a custom loop.
**Action:** When extracting data from Linux procfs files in high-frequency/hot paths, skip generic string-splitting functions in favor of hand-rolled indexing loops for scanning string values.
## 2024-05-18 - Cache static system metrics to avoid live mode overhead
**Learning:** Fetching static system metrics (like CPU, OS, Distro) from /proc files or subprocesses repeatedly during live mode iterations results in significant I/O and process allocation overhead.
**Action:** Use `sync.Once` and static caching to ensure properties that do not change over the lifecycle of the application are fetched only once.
