## 2024-07-05 - Kernel Version Retrieval Performance
**Learning:** Shelling out to `uname -r` via `exec.Command` in `runCommand` takes ~2.5ms per call. While this only happens once per run, using the `syscall.Uname` natively available in Go takes less than ~0.001ms (a 2500x speedup). Calling external commands in a compiled fastfetch alternative is a known bottleneck that should be replaced by native syscalls where possible.
**Action:** Replace `runCommand("uname", "-r")` with native `syscall.Uname` on Linux to improve start-up performance and reduce unnecessary sub-processes.

## 2024-08-11 - Disk Usage Retrieval Performance
**Learning:** Shelling out to `df -Ph /` via `exec.Command` inside `getDisk` and `collectDiskPercent` takes ~3.5ms per call. In contrast, using the native `syscall.Statfs` takes ~2.5µs. This >1000x speedup drastically reduces startup latency and CPU overhead, especially during high-frequency live mode updates where `collectDiskPercent` is called continuously.
**Action:** Replace shell outs to `df` with `syscall.Statfs` to calculate disk percentage usage natively, and parse `/proc/mounts` on Linux to get the filesystem name instead of spawning a new process.

## 2026-07-20 - Uptime and Swap Retrieval Performance
**Learning:** Parsing `/proc/uptime` and `/proc/meminfo` involves file I/O and string allocations, which takes ~15-30µs per call. Using native `syscall.Sysinfo` accesses system metrics directly and executes in ~1µs (15x-30x speedup), significantly reducing overhead especially during frequent live updates.
**Action:** Replaced `/proc/uptime` and `/proc/meminfo` file parsing with `syscall.Sysinfo` in `getUptime`, `getSwap` and `collectSwapPercent` for Linux to improve performance.
