package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetFilename(t *testing.T) {
	// Setup
	pathProject = "/test/project"

	tests := []struct {
		input    string
		expected string
	}{
		{"/test/project/file.go", "/file.go"},
		{"/test/project/subdir/file.go", "/subdir/file.go"},
		{"/test/project\\subdir\\file.go", "/subdir/file.go"},
	}

	for _, test := range tests {
		result := getFilename(test.input)
		if result != test.expected {
			t.Errorf("getFilename(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestIsExtFilter(t *testing.T) {
	// Setup
	configProject = config{
		ExtFilter: []string{".go", ".js", ".md"},
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"file.go", true},
		{"file.js", true},
		{"file.md", true},
		{"file.txt", false},
		{"file.py", false},
	}

	for _, test := range tests {
		result := isExtFilter(test.input)
		if result != test.expected {
			t.Errorf("isExtFilter(%s) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestLoadFileData(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

func main() {
	fmt.Println("Hello")
}

func test() {
	return "test"
}`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Setup config
	configProject = config{
		MethodFilter: []string{`func \(.*\) .*\(.*\).*{`, `func .*\(.*\).*{`},
	}

	// Initialize filesData map
	filesData = make(map[string]fileData)

	// Test loading file data
	err = loadFileData(testFile)
	if err != nil {
		t.Errorf("loadFileData() failed: %v", err)
	}

	// Verify data was loaded
	filename := getFilename(testFile)
	data, exists := filesData[filename]
	if !exists {
		t.Errorf("File data not loaded for %s", filename)
	}

	// Check content
	expectedContent := strings.Split(content, "\n")
	if len(data.Content) != len(expectedContent) {
		t.Errorf("Content length mismatch: got %d, expected %d", len(data.Content), len(expectedContent))
	}

	// Check methods detection
	expectedMethods := []int{2, 6} // Lines where func declarations are
	if len(data.Methods) != len(expectedMethods) {
		t.Errorf("Methods count mismatch: got %d, expected %d", len(data.Methods), len(expectedMethods))
	}
}

func TestLoadFileDataTooLarge(t *testing.T) {
	// Create temporary large test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "large.go")

	// Create file larger than 10MB limit
	largeContent := strings.Repeat("x", 11*1024*1024) // 11MB

	err := os.WriteFile(testFile, []byte(largeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	// Test loading large file data
	err = loadFileData(testFile)
	if err == nil {
		t.Error("Expected error for large file, got nil")
	}

	expectedError := "file too large"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

func TestChangedUserField(t *testing.T) {
	// Setup
	pathProject = "/test/project"

	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"file.go<>main<>checked", "1", false},
		{"file.go<>test<>notes", "test notes", false},
		{"invalid-format", "1", true},
		{"file.go<>main", "1", true},   // Missing field
		{"<>main<>checked", "1", true}, // Missing filename
	}

	for _, test := range tests {
		result := changedUserField(test.name, test.value)

		if test.expectErr && result {
			t.Errorf("changedUserField(%s, %s) expected error but got success", test.name, test.value)
		}

		if !test.expectErr && !result {
			t.Errorf("changedUserField(%s, %s) expected success but got error", test.name, test.value)
		}
	}
}

func TestFromWindows1252(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\x80", "€"},      // Euro sign
		{"\x93", "\u201C"}, // Left double quote
		{"\x94", "\u201D"}, // Right double quote
		{"\x85", "\u2026"}, // Horizontal ellipsis
		{"\xA9", "©"},      // Copyright (unchanged)
		{"Hello", "Hello"}, // ASCII unchanged
	}

	for _, test := range tests {
		result := fromWindows1252(test.input)
		if result != test.expected {
			t.Errorf("fromWindows1252(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestScanProject(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"main.go": `package main
func main() {}`,
		"utils/helper.go": `package utils
func helper() {}`,
		"docs/readme.md": `# README`,
		"test.txt":       "text file",
	}

	for filePath, content := range testFiles {
		fullPath := filepath.Join(tmpDir, filePath)
		dir := filepath.Dir(fullPath)

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Setup config
	configProject = config{
		ExtFilter: []string{".go"},
	}

	// Test scanning project
	files, err := scanProject(tmpDir+string(filepath.Separator), nil)
	if err != nil {
		t.Errorf("scanProject() failed: %v", err)
	}

	// Should find only .go files
	expectedCount := 2 // main.go and helper.go
	if len(files) < expectedCount {
		t.Errorf("Expected at least %d files, got %d", expectedCount, len(files))
	}

	// Check that filesData was populated
	if len(filesData) < expectedCount {
		t.Errorf("Expected at least %d entries in filesData, got %d", expectedCount, len(filesData))
	}
}
