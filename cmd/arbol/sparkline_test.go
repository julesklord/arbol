package main

import (
	"reflect"
	"testing"
	"time"
)

func TestNewSparklineBuffer(t *testing.T) {
	interval := 1 * time.Second
	maxLen := 10

	sb := NewSparklineBuffer(maxLen, interval)

	if sb == nil {
		t.Fatal("NewSparklineBuffer returned nil")
	}
	if sb.maxLen != maxLen {
		t.Errorf("expected maxLen %d, got %d", maxLen, sb.maxLen)
	}
	if sb.interval != interval {
		t.Errorf("expected interval %v, got %v", interval, sb.interval)
	}
	if sb.stopCh == nil {
		t.Error("expected stopCh to be initialized")
	}
	if cap(sb.values) != maxLen {
		t.Errorf("expected values capacity %d, got %d", maxLen, cap(sb.values))
	}
	if len(sb.values) != 0 {
		t.Errorf("expected initial values length 0, got %d", len(sb.values))
	}
}

func TestSparklineBuffer_Add(t *testing.T) {
	tests := []struct {
		name     string
		maxLen   int
		adds     []int
		expected []int
	}{
		{
			name:     "Add within limits",
			maxLen:   5,
			adds:     []int{10, 20, 30},
			expected: []int{10, 20, 30},
		},
		{
			name:     "Clamp below zero",
			maxLen:   5,
			adds:     []int{-5, -10},
			expected: []int{0, 0},
		},
		{
			name:     "Clamp above 100",
			maxLen:   5,
			adds:     []int{105, 200},
			expected: []int{100, 100},
		},
		{
			name:     "Buffer rolling limit",
			maxLen:   3,
			adds:     []int{10, 20, 30, 40, 50},
			expected: []int{30, 40, 50}, // Should keep the last 3 elements
		},
		{
			name:     "Buffer rolling limit with clamp",
			maxLen:   2,
			adds:     []int{50, -5, 105},
			expected: []int{0, 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSparklineBuffer(tt.maxLen, 1*time.Second)
			for _, val := range tt.adds {
				sb.Add(val)
			}
			if !reflect.DeepEqual(sb.values, tt.expected) {
				t.Errorf("expected values %v, got %v", tt.expected, sb.values)
			}
		})
	}
}

func TestSparklineBuffer_Get(t *testing.T) {
	sb := NewSparklineBuffer(5, 1*time.Second)
	sb.Add(10)
	sb.Add(20)

	got := sb.Get()
	expected := []int{10, 20}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}

	// Verify that a copy was returned, not the same slice
	got[0] = 99
	if sb.values[0] == 99 {
		t.Errorf("Get() returned a reference to the internal slice, expected a copy")
	}
}
