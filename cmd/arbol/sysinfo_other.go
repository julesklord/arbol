//go:build !linux

package main

import (
	"errors"
	"strings"
)

func getProcesses() string {
	out := runCommand("bash", "-c", "ps -ax | wc -l")
	if out != "" {
		return strings.TrimSpace(out)
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
