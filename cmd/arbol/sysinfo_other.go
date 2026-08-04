//go:build !linux

package main

import (
	"errors"
	"strconv"
	"strings"
)

func getProcesses() string {
	// OPTIMIZATION: Avoid shelling out to bash and wc
	out := runCommand("ps", "-ax")
	if out != "" {
		// ps output includes a header line, so count newlines
		count := strings.Count(out, "\n")
		if len(out) > 0 && out[len(out)-1] != '\n' {
			count++
		}
		if count > 1 {
			// Subtract 1 for the header
			return strconv.Itoa(count - 1)
		}
	}
	return "n/a"
}

func getKernel() string {
	out := runCommand("/usr/bin/uname", "-r")
	if out != "" {
		return strings.TrimSpace(out)
	}
	return "n/a"
}

func getSysinfoUptime() (int64, error) {
	return 0, errors.New("not supported")
}

func getSysinfoSwap() (uint64, uint64, error) {
	return 0, 0, errors.New("not supported")
}
