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
## 2024-05-24 - Actionable CLI Flag Errors
**Learning:** When manually parsing flags with strings.HasPrefix, missing values often fall through to generic "unknown flag" errors.
**Action:** Always validate exact string matches for keys that expect values to provide actionable error messages.

## 2026-07-07 - Clear Errors for Explicit CLI Flags
**Learning:** Silent fallbacks for invalid explicit CLI flags confuse users.
**Action:** Always validate explicit flag values and provide actionable error messages (e.g., suggesting --help).

## 2024-07-20 - Distinguish Omitted vs Empty Flags
**Learning:** When validating CLI string flags using strict explicit empty checks (like `flag == ""`), you must distinguish between a flag being omitted entirely by the user (which should fall back to a default value safely) and a flag being explicitly provided by the user with an empty value (e.g., `--theme=`). Applying empty-string rejection logic to variables without checking if the prefix was passed causes omissions to falsely trigger error conditions.
**Action:** When enforcing explicit string flag validation, explicitly scope the validation to the block of code executed *when the flag is matched* (e.g., inside `strings.HasPrefix(arg, "--flag=")`) to prevent regressions that break omitted flags.

## 2024-10-24 - Provide actionable CLI error messages instead of generic help pointers
**Learning:** For a CLI tool, simply printing "Run --help for usage" when a user provides an invalid flag value creates friction. It is much more user-friendly to print the explicit list of valid options directly in the error message for that specific flag, saving the user an extra command.
**Action:** When adding or updating CLI flags that require specific values from a set (like themes, styles, or modes), include the valid options in the error output directly rather than redirecting to the general help screen.

## 2024-12-19 - Improved CLI flag error actionability
**Learning:** Invalid enum-style flag inputs simply listed available values without directing users how to get help.
**Action:** Added explicit suggestions to run `--help` in all CLI flag validation error blocks, matching the UX of unknown flag errors.

## 2024-07-25 - Dimming unavailable data reduces visual noise
**Learning:** When displaying a dense layout of system metrics (like in a fastfetch-style tool), 'n/a' or unavailable data points command too much visual attention when styled identically to available data (e.g. using the same bright italic colors). This creates visual clutter and makes it harder for the user to scan for the metrics that actually matter.
**Action:** Always intercept null/unavailable states ('n/a') in CLI or dashboard text outputs and explicitly format them using a muted/dimmed theme color to push them to the background visually, allowing valid data to stand out.

## 2024-07-26 - Prevent CLI Flag Fallthrough Due to Prefix Overlap
**Learning:** When manually implementing CLI flag parsing with `strings.HasPrefix` (or similar methods) in Go, ensure that longer, more specific prefixes (e.g., `--sparkline-style=`) are evaluated before shorter, related prefixes (e.g., `--sparkline`). This prevents the shorter prefix from incorrectly intercepting and mishandling the argument.
**Action:** When manually parsing flags, always place the most specific/longest prefix checks first.

## 2024-07-28 - Explicitly Handling Omitted Values for Explicit CLI Flags
**Learning:** For explicit string CLI flags (e.g. `--theme` or `--bar-style`), relying only on a general catch-all error handling strategy or falling through to default unknown flag error handlers produces generic errors that require users to run `--help` manually. Surface actionable errors with explicit available options for flags that require a value.
**Action:** When evaluating explicit string flags in CLI, explicitly check for exact matches (omitted value scenarios) alongside prefix matches. For exact matches, surface a targeted, actionable error that lists valid inputs for that specific flag, reducing friction for the end user.
## 2024-07-29 - Explicit Flag Matching

**Learning:** When implementing custom CLI flag parsing in Go, using simple prefix matching (like `strings.HasPrefix(arg, "--flag")`) can swallow incorrectly named flags that happen to start with the same prefix (e.g., `--flags=xyz`), preventing an actionable "unknown flag" error.
**Action:** When manually parsing flags, avoid `strings.HasPrefix` for checking boolean flags. Instead, use an explicit string equality check (`arg == "--flag"`) or an explicitly delimited prefix check (`strings.HasPrefix(arg, "--flag=")`).
