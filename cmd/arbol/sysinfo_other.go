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
		lines := strings.Split(out, "\n")
		// ps output includes a header line, and we split by newline so the last might be empty
		count := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				count++
			}
		}
		if count > 1 {
			// Subtract 1 for the header
			return strconv.Itoa(count - 1)
		}
	}
	return "n/a"
}

func getKernel() string {
	out := runCommand("uname", "-r")
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
