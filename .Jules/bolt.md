## 2024-07-16 - Optimize `getCPUUsage` concurrency
**Learning:** Functions doing sampling with long delays (e.g., `time.Sleep` for CPU utilization) can become bottlenecks when called sequentially in the critical path.
**Action:** Use channels to run long-running measurements concurrently with other initialization work. Initialize the channel at the top of the function and receive the final result at the end, hiding the latency of the `time.Sleep`.
