package main

import (
	"testing"
)

func TestTreeStyleFunctions(t *testing.T) {
	// Save the original style to restore it later
	originalStyle := GetTreeStyle()
	defer SetTreeStyle(originalStyle)

	// Test 1: Check the default style
	// We assume that the original style is the default one because these tests run independently.
	// But to be sure, we can just test if GetTreeStyle() returns something valid.
	// It is better to test explicitly after setting.
	SetTreeStyle(TreeStyleDefault)
	if got := GetTreeStyle(); got != TreeStyleDefault {
		t.Errorf("GetTreeStyle() = %v, want %v", got, TreeStyleDefault)
	}

	// Test 2: Set and Get a different style
	SetTreeStyle(TreeStyleRounded)
	if got := GetTreeStyle(); got != TreeStyleRounded {
		t.Errorf("GetTreeStyle() = %v, want %v", got, TreeStyleRounded)
	}

	// Test 3: Set and Get another style
	SetTreeStyle(TreeStyleASCII)
	if got := GetTreeStyle(); got != TreeStyleASCII {
		t.Errorf("GetTreeStyle() = %v, want %v", got, TreeStyleASCII)
	}
}

func TestGetConnectors(t *testing.T) {
	// Save the original style to restore it later
	originalStyle := GetTreeStyle()
	defer SetTreeStyle(originalStyle)

	// Test 1: Connectors for TreeStyleDefault
	SetTreeStyle(TreeStyleDefault)
	connectors := GetConnectors()
	if connectors.TopLeft != "┌" {
		t.Errorf("GetConnectors() for TreeStyleDefault returned TopLeft = %v, want %v", connectors.TopLeft, "┌")
	}
	if connectors.Branch != "├── " {
		t.Errorf("GetConnectors() for TreeStyleDefault returned Branch = %v, want %v", connectors.Branch, "├── ")
	}

	// Test 2: Connectors for TreeStyleRounded
	SetTreeStyle(TreeStyleRounded)
	connectors = GetConnectors()
	if connectors.TopLeft != "╭" {
		t.Errorf("GetConnectors() for TreeStyleRounded returned TopLeft = %v, want %v", connectors.TopLeft, "╭")
	}

	// Test 3: Connectors for TreeStyleHeavy
	SetTreeStyle(TreeStyleHeavy)
	connectors = GetConnectors()
	if connectors.TopLeft != "┏" {
		t.Errorf("GetConnectors() for TreeStyleHeavy returned TopLeft = %v, want %v", connectors.TopLeft, "┏")
	}

	// Test 4: Connectors for TreeStyleDouble
	SetTreeStyle(TreeStyleDouble)
	connectors = GetConnectors()
	if connectors.TopLeft != "╔" {
		t.Errorf("GetConnectors() for TreeStyleDouble returned TopLeft = %v, want %v", connectors.TopLeft, "╔")
	}

	// Test 5: Connectors for TreeStyleASCII
	SetTreeStyle(TreeStyleASCII)
	connectors = GetConnectors()
	if connectors.TopLeft != "+" {
		t.Errorf("GetConnectors() for TreeStyleASCII returned TopLeft = %v, want %v", connectors.TopLeft, "+")
	}

	// Test 6: Connectors for TreeStyleDotted
	SetTreeStyle(TreeStyleDotted)
	connectors = GetConnectors()
	if connectors.TopLeft != "┆" {
		t.Errorf("GetConnectors() for TreeStyleDotted returned TopLeft = %v, want %v", connectors.TopLeft, "┆")
	}
}
