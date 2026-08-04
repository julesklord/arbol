package main

import (
	"os"
	"testing"
)

func BenchmarkPrintJSON(b *testing.B) {
	info := SystemInfo{
		Host:   "myhost",
		OSName: "Linux",
		Plugins: []PluginInfo{
			{
				Key:     "Plug1",
				Val:     "Val1",
				Details: []string{"detail1", "detail2", "detail3", "detail4", "detail5"},
			},
			{
				Key:     "Plug2",
				Val:     "Val2",
				Details: []string{"detail1", "detail2", "detail3", "detail4", "detail5"},
			},
		},
	}

	// redirect stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Note: rather than redirecting stdout with a pipe which can leak memory or goroutines,
	// we just discard stdout for this benchmark.
	// However, os.Stdout is used directly in export.go, so we DO need a pipe,
	// but we should just read and discard it.
	r, w, _ := os.Pipe()
	os.Stdout = w

	go func() {
		// Just read and discard to prevent blocking
		buf := make([]byte, 1024)
		for {
			_, err := r.Read(buf)
			if err != nil {
				break
			}
		}
	}()
	defer r.Close()
	defer w.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		printJSON(info)
	}
}
