package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestSetUserValue(t *testing.T) {
	// Reset userFields before test
	userFields = make([]fieldsData, 0)

	// Test setting new value
	setUserValue("file.go", "main", "checked", "1")

	if len(userFields) != 1 {
		t.Errorf("Expected 1 user field, got %d", len(userFields))
	}

	field := userFields[0]
	if field.Filename != "file.go" || field.Method != "main" || field.Field != "checked" || field.Value != "1" {
		t.Errorf("Field data mismatch: got %+v", field)
	}

	// Test updating existing value
	setUserValue("file.go", "main", "checked", "0")

	if len(userFields) != 1 {
		t.Errorf("Expected still 1 user field, got %d", len(userFields))
	}

	field = userFields[0]
	if field.Value != "0" {
		t.Errorf("Expected updated value '0', got '%s'", field.Value)
	}

	// Test adding another field
	setUserValue("file.go", "main", "notes", "test notes")

	if len(userFields) != 2 {
		t.Errorf("Expected 2 user fields, got %d", len(userFields))
	}
}

func TestGetUserValue(t *testing.T) {
	// Reset and setup test data
	userFields = []fieldsData{
		{"file.go", "main", "checked", "1"},
		{"file.go", "main", "notes", "test notes"},
		{"utils.go", "helper", "checked", "0"},
	}

	// Test getting existing values
	tests := []struct {
		filename string
		method   string
		field    string
		expected string
	}{
		{"file.go", "main", "checked", "1"},
		{"file.go", "main", "notes", "test notes"},
		{"utils.go", "helper", "checked", "0"},
	}

	for _, test := range tests {
		result := getUserValue(test.filename, test.method, test.field)
		if result != test.expected {
			t.Errorf("getUserValue(%s, %s, %s) = %s; expected %s",
				test.filename, test.method, test.field, result, test.expected)
		}
	}

	// Test getting non-existing value
	result := getUserValue("nonexistent.go", "main", "checked")
	if result != "" {
		t.Errorf("Expected empty string for non-existent field, got '%s'", result)
	}
}

func TestConcurrentUserFieldAccess(t *testing.T) {
	// Reset userFields
	userFields = make([]fieldsData, 0)

	// Test concurrent access to userFields
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// Start multiple goroutines setting values
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				filename := "file.go"
				method := "main"
				field := "checked"
				value := "1"

				setUserValue(filename, method, field, value)

				// Also test reading
				getUserValue(filename, method, field)
			}
		}(i)
	}

	wg.Wait()

	// Verify we didn't lose any data due to race conditions
	if len(userFields) == 0 {
		t.Error("Expected userFields to contain data after concurrent operations")
	}

	// Verify data integrity
	for _, field := range userFields {
		if field.Filename == "" || field.Method == "" || field.Field == "" {
			t.Error("Found incomplete field data due to race condition")
		}
	}
}

func TestLoadUserFields(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	pathProject = tmpDir

	// Create test user fields file
	testFile := filepath.Join(tmpDir, "zoomer-userfields.json")
	testData := []fieldsData{
		{"file.go", "main", "checked", "1"},
		{"utils.go", "helper", "notes", "test notes"},
	}

	// Write test data to file
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	// Use JSON encoding (same as in the actual function)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(testData)
	if err != nil {
		t.Fatalf("Failed to encode test data: %v", err)
	}

	// Reset userFields
	userFields = make([]fieldsData, 0)

	// Test loading
	result := loadUserFields()
	if !result {
		t.Error("loadUserFields() returned false")
	}

	if len(userFields) != 2 {
		t.Errorf("Expected 2 user fields, got %d", len(userFields))
	}

	// Verify data was loaded correctly
	if userFields[0].Filename != "file.go" || userFields[0].Value != "1" {
		t.Errorf("First field data mismatch: got %+v", userFields[0])
	}
}

func TestSaveFileUserFields(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	pathProject = tmpDir

	// Setup test data
	userFields = []fieldsData{
		{"file.go", "main", "checked", "1"},
		{"utils.go", "helper", "notes", "test notes"},
	}

	// Reset timestamps
	lastChange = time.Now()
	lastSave = time.Now().Add(-time.Minute) // Make it older than lastChange

	// Test saving
	saveFileUserFields()

	// Verify file was created
	testFile := filepath.Join(tmpDir, "zoomer-userfields.json")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("User fields file was not created")
	}

	// Verify file content
	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open saved file: %v", err)
	}
	defer file.Close()

	var loadedData []fieldsData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&loadedData)
	if err != nil {
		t.Fatalf("Failed to decode saved file: %v", err)
	}

	if len(loadedData) != 2 {
		t.Errorf("Expected 2 fields in saved file, got %d", len(loadedData))
	}

	// Verify timestamps were updated
	if !lastSave.After(lastChange) {
		t.Error("lastSave timestamp was not updated")
	}
}

func TestSaveFileUserFieldsNoChanges(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	pathProject = tmpDir

	// Setup test data
	userFields = []fieldsData{
		{"file.go", "main", "checked", "1"},
	}

	// Set timestamps so no save is needed
	now := time.Now()
	lastChange = now.Add(-time.Minute)
	lastSave = now // lastSave is more recent than lastChange

	// Test saving (should not save)
	saveFileUserFields()

	// Verify file was NOT created
	testFile := filepath.Join(tmpDir, "zoomer-userfields.json")
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("User fields file was created when it shouldn't have been")
	}
}
