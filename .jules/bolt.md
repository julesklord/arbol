## 2024-07-05 - Kernel Version Retrieval Performance
**Learning:** Shelling out to `uname -r` via `exec.Command` in `runCommand` takes ~2.5ms per call. While this only happens once per run, using the `syscall.Uname` natively available in Go takes less than ~0.001ms (a 2500x speedup). Calling external commands in a compiled fastfetch alternative is a known bottleneck that should be replaced by native syscalls where possible.
**Action:** Replace `runCommand("uname", "-r")` with native `syscall.Uname` on Linux to improve start-up performance and reduce unnecessary sub-processes.
