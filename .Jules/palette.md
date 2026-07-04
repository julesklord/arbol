## 2024-06-30 - Invalid CLI Flag Feedback
**Learning:** The CLI tool was continuing execution and outputting mangled data when provided with unknown flags, causing user confusion. Providing immediate, clear error feedback with usage instructions improves developer experience.
**Action:** Always validate all input flags and exit early with helpful error messages and usage instructions for unknown arguments.

## 2024-07-01 - Changed default bar style
**Learning:** Braille characters in terminal can be hard to read for some users or terminal emulators without proper font support.
**Action:** Change default bar style to `BarStyleBlock`.

## 2026-07-04 - Support NO_COLOR environment variable
**Learning:** Adding NO_COLOR support is an important accessibility improvement for users with visual sensitivities, colorblindness, or low-contrast requirements, providing them with a way to easily disable ANSI colors.
**Action:** Always implement a `ColorDisabled` mechanism tied to the `NO_COLOR` standard environment variable, modifying rendering functions to strip or omit ANSI color codes based on this flag.
