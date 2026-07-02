## 2024-06-25 - ANSI parsing bottlenecks
**Learning:** `visualLength` function allocates a new string via `stripANSI` before measuring. Go's string manipulation can cause a lot of GC pressure.
**Action:** Always check if we can process strings in a single pass without intermediate allocation when calculating lengths or parsing.
