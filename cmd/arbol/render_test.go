package main

import (
	"testing"
)

func TestSetBarStyle(t *testing.T) {
	// Save original style to restore after test
	originalStyle := currentBarStyle
	defer func() {
		currentBarStyle = originalStyle
	}()

	tests := []struct {
		name     string
		style    BarStyle
		expected BarStyle
	}{
		{
			name:     "Set to BarStyleBraille",
			style:    BarStyleBraille,
			expected: BarStyleBraille,
		},
		{
			name:     "Set to BarStyleGradient",
			style:    BarStyleGradient,
			expected: BarStyleGradient,
		},
		{
			name:     "Set back to BarStyleBlock",
			style:    BarStyleBlock,
			expected: BarStyleBlock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetBarStyle(tt.style)
			if currentBarStyle != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, currentBarStyle)
			}
		})
	}
}
