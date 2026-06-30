package main

import (
	"strings"
	"unicode/utf8"
)


const (
	ansiStateNormal = iota
	ansiStateEscape
	ansiStateCSI
	ansiStateOSC
)

func runeWidth(r rune) int {
	// Zero-width space, joiners, control chars, variation selectors
	if r == '\u200d' || r == '\u200c' || (r >= '\ufe00' && r <= '\ufe0f') {
		return 0
	}
	// Combining diacritical marks
	if r >= 0x0300 && r <= 0x036F {
		return 0
	}
	// Wide ranges (2 columns)
	// Emojis / Pictographs in SMP (Plane 1): U+1F000 to U+1FAFF
	if r >= 0x1F000 && r <= 0x1FAFF {
		return 2
	}
	// Miscellaneous Symbols and Pictographs, Emoticons, Ornamental Dingbats, etc. in BMP
	if r >= 0x2600 && r <= 0x27BF {
		return 2
	}
	// CJK ranges
	if (r >= 0x2E80 && r <= 0x2FDF) || // CJK Radicals
		(r >= 0x3000 && r <= 0x9FFF) || // Hiragana, Katakana, CJK Unified Ideographs
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility
		(r >= 0xFF01 && r <= 0xFF60) || // Fullwidth Forms
		(r >= 0xFFE0 && r <= 0xFFE6) {
		return 2
	}
	// Default
	return 1
}

func visualLength(s string) int {
	raw := stripANSI(s)
	length := 0
	for _, r := range raw {
		length += runeWidth(r)
	}
	return length
}

func truncateANSI(s string, limit int) string {
	if visualLength(s) <= limit {
		return s
	}

	var builder strings.Builder
	visualLen := 0
	restoreCode := "[0m"
	targetLen := limit - 1
	if targetLen < 0 {
		targetLen = 0
	}

	state := ansiStateNormal

	for i := 0; i < len(s); i++ {
		ch := s[i]

		switch state {
		case ansiStateNormal:
			if ch == '' {
				state = ansiStateEscape
				builder.WriteByte(ch)
			} else {
				if visualLen < targetLen {
					r, size := utf8.DecodeRuneInString(s[i:])
					w := runeWidth(r)
					if visualLen+w <= targetLen {
						builder.WriteRune(r)
						visualLen += w
					} else {
						visualLen = targetLen
					}
					i += size - 1
				}
			}
		case ansiStateEscape:
			builder.WriteByte(ch)
			if ch == '[' {
				state = ansiStateCSI
			} else if ch == ']' {
				state = ansiStateOSC
			} else {
				state = ansiStateNormal
			}
		case ansiStateCSI:
			builder.WriteByte(ch)
			if ch >= 0x40 && ch <= 0x7E {
				state = ansiStateNormal
			}
		case ansiStateOSC:
			builder.WriteByte(ch)
			if ch == '' {
				state = ansiStateNormal
			} else if ch == '' {
				state = ansiStateEscape
			}
		}
	}
	builder.WriteString("…")
	builder.WriteString(restoreCode)
	return builder.String()
}

func stripANSI(s string) string {
	var builder strings.Builder
	state := ansiStateNormal

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch state {
		case ansiStateNormal:
			if ch == '' {
				state = ansiStateEscape
			} else {
				builder.WriteByte(ch)
			}
		case ansiStateEscape:
			if ch == '[' {
				state = ansiStateCSI
			} else if ch == ']' {
				state = ansiStateOSC
			} else {
				state = ansiStateNormal
			}
		case ansiStateCSI:
			if ch >= 0x40 && ch <= 0x7E {
				state = ansiStateNormal
			}
		case ansiStateOSC:
			if ch == '' {
				state = ansiStateNormal
			} else if ch == '' {
				state = ansiStateEscape
			}
		}
	}
	return builder.String()
}

func getBar(pct int) string {
	if pct < 0 {
		pct = 0
	}
	filled := pct / 10
	if filled > 10 {
		filled = 10
	}
	empty := 10 - filled
	color := "\033[01;32m" // Green
	if pct > 80 {
		color = "\033[01;31m" // Red
	} else if pct > 50 {
		color = "\033[01;33m" // Yellow
	}
	restore := "\033[0m"
	gray := "\033[00;37m"

	var sb strings.Builder
	sb.WriteString(color)
	for i := 0; i < filled; i++ {
		sb.WriteString("█")
	}
	sb.WriteString(restore + gray)
	for i := 0; i < empty; i++ {
		sb.WriteString("░")
	}
	sb.WriteString(restore)
	return sb.String()
}

func padString(s string, width int) string {
	rawLen := visualLength(s)
	if rawLen >= width {
		return ""
	}
	return strings.Repeat(" ", width-rawLen)
}
