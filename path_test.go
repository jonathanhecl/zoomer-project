package main

import (
	"os"
	"testing"
)

func TestIsValidFile(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "testfile_*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create temporary directory
	tmpDir := t.TempDir()

	tests := []struct {
		input    string
		expected bool
	}{
		{tmpFile.Name(), true},     // Existing file
		{tmpDir, false},            // Directory
		{"nonexistent.txt", false}, // Non-existent file
		{"", false},                // Empty path
	}

	for _, test := range tests {
		result := isValidFile(test.input)
		if result != test.expected {
			t.Errorf("isValidFile(%s) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidPath(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "testfile_*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		input    string
		expected bool
	}{
		{tmpDir, true},          // Existing directory
		{tmpFile.Name(), false}, // File (not directory)
		{"nonexistent", false},  // Non-existent directory
		{"", false},             // Empty path
	}

	for _, test := range tests {
		result := isValidPath(test.input)
		if result != test.expected {
			t.Errorf("isValidPath(%s) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidPort(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"80", true},      // Valid port
		{"443", true},     // Valid port
		{"8080", true},    // Valid port
		{"65535", true},   // Max valid port
		{"1", true},       // Min valid port
		{"0", false},      // Below min
		{"65536", false},  // Above max
		{"-1", false},     // Negative
		{"abc", false},    // Non-numeric
		{"", false},       // Empty
		{"80.5", false},   // Float
		{"8080 ", false},  // With space
		{" 8080", false},  // Leading space
		{"8080\n", false}, // With newline
		{"99999", false},  // Too large
	}

	for _, test := range tests {
		result := isValidPort(test.input)
		if result != test.expected {
			t.Errorf("isValidPort(%s) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestEdgeCases(t *testing.T) {
	// Test with special characters in paths
	specialDir := t.TempDir()
	specialFile := specialDir + "/file with spaces.go"

	err := os.WriteFile(specialFile, []byte("package main"), 0644)
	if err != nil {
		t.Fatalf("Failed to create special file: %v", err)
	}

	if !isValidFile(specialFile) {
		t.Error("isValidFile failed for file with spaces")
	}

	// Test with very long port string
	longPort := "12345678901234567890"
	if isValidPort(longPort) {
		t.Error("isValidPort should fail for very long port string")
	}

	// Test with port-like but invalid strings
	invalidPorts := []string{"80a", "a80", "8-0", "8_0", "80.0"}
	for _, port := range invalidPorts {
		if isValidPort(port) {
			t.Errorf("isValidPort should fail for invalid port: %s", port)
		}
	}
}
