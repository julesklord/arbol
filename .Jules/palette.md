## 2024-05-18 - Graceful Exit Cursor Restoration
**Learning:** In Go CLI applications that manually manipulate the terminal cursor (e.g., hiding it with `\033[?25l`), deferred functions (like `defer fmt.Print("\033[?25h")`) will NOT execute if the program is terminated via an unhandled OS signal like `SIGINT` (Ctrl+C). This leaves the user with a permanently hidden cursor, which is a major UX annoyance.
**Action:** Always capture `os.Signal` (`syscall.SIGINT`, `syscall.SIGTERM`) explicitly in long-running terminal loops (like live mode) and return cleanly to guarantee that `defer` statements for terminal state restoration are executed.
## 2024-06-30 - Invalid CLI Flag Feedback
**Learning:** The CLI tool was continuing execution and outputting mangled data when provided with unknown flags, causing user confusion. Providing immediate, clear error feedback with usage instructions improves developer experience.
**Action:** Always validate all input flags and exit early with helpful error messages and usage instructions for unknown arguments.

## 2024-07-01 - Changed default bar style
**Learning:** Braille characters in terminal can be hard to read for some users or terminal emulators without proper font support.
**Action:** Change default bar style to `BarStyleBlock`.

## 2026-07-04 - Support NO_COLOR environment variable
**Learning:** Adding NO_COLOR support is an important accessibility improvement for users with visual sensitivities, colorblindness, or low-contrast requirements, providing them with a way to easily disable ANSI colors.
**Action:** Always implement a `ColorDisabled` mechanism tied to the `NO_COLOR` standard environment variable, modifying rendering functions to strip or omit ANSI color codes based on this flag.
