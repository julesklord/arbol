package main

import (
	"os"
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	out := runCommand("echo", "hello")
	if out != "hello" {
		t.Errorf("Expected 'hello', got '%s'", out)
	}

	outFail := runCommand("this-command-does-not-exist-12345")
	if outFail != "" {
		t.Errorf("Expected empty string on failure, got '%s'", outFail)
	}
}

func TestRunCommand_ExitCode(t *testing.T) {
	out := runCommand("sh", "-c", "exit 1")
	if out != "" {
		t.Errorf("Expected empty string for non-zero exit code, got %q", out)
	}
}

func TestRunCommandWithTimeout(t *testing.T) {
	out := runCommandWithTimeout(1*time.Second, "echo", "hello")
	if out != "hello" {
		t.Errorf("Expected 'hello', got '%s'", out)
	}

	outFail := runCommandWithTimeout(1*time.Second, "this-command-does-not-exist-12345")
	if outFail != "" {
		t.Errorf("Expected empty string on failure, got '%s'", outFail)
	}
}

func TestRunCommandWithTimeout_Timeout(t *testing.T) {
	out := runCommandWithTimeout(100*time.Millisecond, "sleep", "5")
	if out != "" {
		t.Errorf("Expected empty string on timeout, got %q", out)
	}
}

func TestGetOSName(t *testing.T) {
	name := getOSName()
	if name == "" {
		t.Error("getOSName() returned empty string")
	}
}

func TestGetUptime(t *testing.T) {
	uptime := getUptime()
	if uptime == "" {
		t.Error("getUptime() returned empty string")
	}
}

func TestGetTerminalWidth(t *testing.T) {
	// In a test environment, tput may return empty, so we just check it doesn't panic
	// and returns a positive value
	// Set TERM so tput can run
	os.Setenv("TERM", "xterm-256color")
	w := getTerminalWidth()
	if w <= 0 {
		t.Errorf("getTerminalWidth() returned non-positive value: %d", w)
	}
}
