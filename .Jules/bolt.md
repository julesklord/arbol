## 2024-05-18 - Async System Calls Optimization
**Learning:** In Go-based CLI status fetchers, sequential synchronous shell commands (`os/exec`) often become the primary execution bottleneck, bound by the sum of individual command latencies rather than CPU processing time. The time taken to execute `gatherInfo` dropped from ~2.5s down to ~0.7s when running metrics concurrently in a single goroutine batch (even with a sleep delay inside the CPU metric collector).
**Action:** Always verify if independent I/O or shell tasks can be grouped and executed inside goroutines combined with `sync.WaitGroup` to mask blocking latency in performance-sensitive terminal utilities.
## 2024-07-20 - Optimize strings.Split for disk percent collection
**Learning:** Using `strings.Split` on strings with many newlines when only the first few lines are needed allocates a slice for all newlines. Using `strings.SplitN` with the required limit reduces execution time (e.g. from 357.8 ns/op to 279.3 ns/op) and memory allocations (from 192 B/op to 144 B/op).
**Action:** When only a fixed number of lines/elements are needed from a split operation, always prefer `strings.SplitN` over `strings.Split` to optimize memory and CPU usage.
## 2024-10-24 - Avoiding `bash -c` Pipelines in Go
**Learning:** Shelling out to `bash -c "cmd | grep | awk"` creates immense overhead due to spawning multiple subprocesses. `getGPU` execution dropped from ~5.8ms to ~0.05ms simply by invoking the base command natively and handling string manipulation in Go.
**Action:** When retrieving system info, always invoke native commands (e.g., `runCommand("ps", "-A")`) and use pure Go `strings` and `strconv` packages to parse the output instead of relying on shell pipes.
