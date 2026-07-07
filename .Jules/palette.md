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

## 2024-07-05 - Validate CLI Flags and Exit Early on Errors
**Learning:** In CLI applications, silently ignoring unknown flag values and falling back to defaults can be confusing for users. When provided with invalid configurations, providing immediate error feedback is a core UX improvement.
**Action:** When creating CLI apps, validate flag values like themes and styles directly after parsing. If invalid, output a clear error message to `stderr` and exit with a non-zero status code rather than silently failing to a default.

## 2024-12-06 - Error Message Clarity
**Learning:** Command-line tools should return actionable and friendly error messages rather than terse generic errors. Suggesting using `--help` immediately directs the user to the correct next step.
**Action:** When updating or reviewing CLI argument parsers, always ensure error messages not only describe the problem but also provide the user with clear instructions on how to find the supported options or correct their mistake.

## 2024-07-06 - Change bar colors to theme colors
**Learning:** Hardcoding standard ANSI colors for progress bars (green/yellow/red) breaks the visual coherence of customized color themes. A UI component should inherit its semantic colors (success/warning/error) from the active theme palette rather than using fixed ANSI escape codes to ensure a consistent look and feel across different themes.
**Action:** Replace hardcoded ANSI colors (\033[01;32m etc) with dynamic lookups from the active theme configuration (theme.BarColors[0], theme.Muted, etc) for UI components like progress bars.

## 2024-07-06 - Error Message Actionability
**Learning:** Terse CLI error messages like 'Unknown logo mode: foo' are confusing and require the user to guess next steps. Error messages should point to the solution.
**Action:** Update error message text to provide clear steps to resolve the issue (e.g., suggesting '--help').
