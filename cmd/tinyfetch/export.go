package main

import (
	"fmt"
	"strings"
)

func printJSON(host, osName, kernel, uptime, shell, cpu, mem, disk string, keys, vals []string) {
	fmt.Printf("{\n")
	fmt.Printf("  \"host\": %q,\n", host)
	fmt.Printf("  \"os\": %q,\n", osName)
	fmt.Printf("  \"kernel\": %q,\n", kernel)
	fmt.Printf("  \"uptime\": %q,\n", uptime)
	fmt.Printf("  \"shell\": %q,\n", shell)
	fmt.Printf("  \"cpu\": %q,\n", cpu)
	fmt.Printf("  \"memory\": %q,\n", mem)
	fmt.Printf("  \"disk\": %q", disk)

	if len(keys) > 0 {
		fmt.Printf(",\n  \"plugins\": {\n")
		for i := 0; i < len(keys); i++ {
			cleanVal := stripANSI(vals[i])
			fmt.Printf("    %q: %q", keys[i], cleanVal)
			if i < len(keys)-1 {
				fmt.Printf(",\n")
			} else {
				fmt.Printf("\n")
			}
		}
		fmt.Printf("  }\n")
	} else {
		fmt.Printf("\n")
	}
	fmt.Printf("}\n")
}

func escapeXML(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			sb.WriteString("&amp;")
		case '<':
			sb.WriteString("&lt;")
		case '>':
			sb.WriteString("&gt;")
		case '"':
			sb.WriteString("&quot;")
		case '\'':
			sb.WriteString("&apos;")
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func printXML(host, osName, kernel, uptime, shell, cpu, mem, disk string, keys, vals []string) {
	fmt.Printf("<tinyfetch>\n")
	fmt.Printf("  <host>%s</host>\n", escapeXML(host))
	fmt.Printf("  <os>%s</os>\n", escapeXML(osName))
	fmt.Printf("  <kernel>%s</kernel>\n", escapeXML(kernel))
	fmt.Printf("  <uptime>%s</uptime>\n", escapeXML(uptime))
	fmt.Printf("  <shell>%s</shell>\n", escapeXML(shell))
	fmt.Printf("  <cpu>%s</cpu>\n", escapeXML(cpu))
	fmt.Printf("  <memory>%s</memory>\n", escapeXML(mem))
	fmt.Printf("  <disk>%s</disk>\n", escapeXML(disk))
	if len(keys) > 0 {
		fmt.Printf("  <plugins>\n")
		for i := 0; i < len(keys); i++ {
			tag := strings.ToLower(keys[i])
			var sb strings.Builder
			for _, r := range tag {
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
					sb.WriteRune(r)
				} else {
					sb.WriteRune('_')
				}
			}
			tagStr := sb.String()
			cleanVal := stripANSI(vals[i])
			fmt.Printf("    <%s>%s</%s>\n", tagStr, escapeXML(cleanVal), tagStr)
		}
		fmt.Printf("  </plugins>\n")
	}
	fmt.Printf("</tinyfetch>\n")
}

func printTXT(host, osName, kernel, uptime, shell, cpu, mem, disk string, keys, vals []string) {
	fmt.Printf("Host: %s\n", host)
	fmt.Printf("OS: %s\n", osName)
	fmt.Printf("Kernel: %s\n", kernel)
	fmt.Printf("Uptime: %s\n", uptime)
	fmt.Printf("Shell: %s\n", shell)
	fmt.Printf("CPU: %s\n", cpu)
	fmt.Printf("Memory: %s\n", mem)
	fmt.Printf("Disk: %s\n", disk)
	for i := 0; i < len(keys); i++ {
		fmt.Printf("%s: %s\n", keys[i], stripANSI(vals[i]))
	}
}
