package main

import (
	"sort"
	"testing"
)

func TestGetTheme(t *testing.T) {
	// Ensure we start with default theme
	originalTheme := GetTheme()
	defer func() { currentTheme = originalTheme }()

	currentTheme = themes["default"]
	theme := GetTheme()
	if theme.Name != "default" {
		t.Errorf("GetTheme() returned theme Name %q, expected %q", theme.Name, "default")
	}
}

func TestSetTheme(t *testing.T) {
	// Restore original theme after test
	originalTheme := GetTheme()
	defer func() { currentTheme = originalTheme }()

	tests := []struct {
		name          string
		themeName     string
		expectedBool  bool
		expectedTheme string
	}{
		{
			name:          "Valid theme existing",
			themeName:     "catppuccin",
			expectedBool:  true,
			expectedTheme: "catppuccin",
		},
		{
			name:          "Valid theme case insensitive",
			themeName:     "CaTpPuCcIn",
			expectedBool:  true,
			expectedTheme: "catppuccin",
		},
		{
			name:          "Invalid theme",
			themeName:     "nonexistent",
			expectedBool:  false,
			expectedTheme: "catppuccin", // should not change from previous test
		},
		{
			name:          "Set back to default",
			themeName:     "default",
			expectedBool:  true,
			expectedTheme: "default",
		},
	}

	// Set initial state
	currentTheme = themes["default"]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SetTheme(tt.themeName)
			if result != tt.expectedBool {
				t.Errorf("SetTheme(%q) returned %v, expected %v", tt.themeName, result, tt.expectedBool)
			}

			current := GetTheme()
			if current.Name != tt.expectedTheme {
				t.Errorf("After SetTheme(%q), current theme is %q, expected %q", tt.themeName, current.Name, tt.expectedTheme)
			}
		})
	}
}

func TestListThemes(t *testing.T) {
	themeList := ListThemes()

	// Check length matches the map
	if len(themeList) != len(themes) {
		t.Errorf("ListThemes() returned %d themes, expected %d", len(themeList), len(themes))
	}

	// Check that specific known themes are in the list
	expectedThemes := []string{"default", "catppuccin"}

	for _, expected := range expectedThemes {
		found := false
		for _, name := range themeList {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ListThemes() missing expected theme: %s", expected)
		}
	}

	// Verify there are no duplicates in the returned list
	sort.Strings(themeList)
	for i := 1; i < len(themeList); i++ {
		if themeList[i] == themeList[i-1] {
			t.Errorf("ListThemes() returned duplicate theme: %s", themeList[i])
		}
	}
}
