package main

import (
	"testing"
)

func BenchmarkGetProcesses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getProcesses()
	}
}

func BenchmarkGetGPU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getGPU()
	}
}

func BenchmarkGetMemory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getMemory()
	}
}

func BenchmarkGetCPUUsage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCPUUsage()
	}
}
