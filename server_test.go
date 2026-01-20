package main

import (
	"testing"
)

func TestDisassemblyFieldName(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
		filename    string
		method      string
		field       string
	}{
		{
			input:    "file.go<>main<>checked",
			filename: "file.go",
			method:   "main",
			field:    "checked",
		},
		{
			input:    "utils/helper.go<>testFunc<>notes",
			filename: "utils/helper.go",
			method:   "testFunc",
			field:    "notes",
		},
		{
			input:       "invalid-format",
			expectError: true,
		},
		{
			input:       "file.go<>main",
			expectError: true,
		},
		{
			input:       "file.go<>main<>checked<>extra",
			expectError: true,
		},
	}

	for _, test := range tests {
		filename, method, field := disassemblyFieldName(test.input)

		if test.expectError {
			if filename != "" || method != "" || field != "" {
				t.Errorf("disassemblyFieldName(%s) expected empty results for invalid input, got (%s, %s, %s)",
					test.input, filename, method, field)
			}
		} else {
			if filename != test.filename || method != test.method || field != test.field {
				t.Errorf("disassemblyFieldName(%s) = (%s, %s, %s); expected (%s, %s, %s)",
					test.input, filename, method, field, test.filename, test.method, test.field)
			}
		}
	}
}

func TestCreateFieldName(t *testing.T) {
	tests := []struct {
		filename string
		method   string
		field    string
		expected string
	}{
		{
			filename: "file.go",
			method:   "main",
			field:    "checked",
			expected: "file.go<>main<>checked",
		},
		{
			filename: "utils/helper.go",
			method:   "testFunc",
			field:    "notes",
			expected: "utils/helper.go<>testFunc<>notes",
		},
		{
			filename: "",
			method:   "main",
			field:    "checked",
			expected: "<>main<>checked",
		},
	}

	for _, test := range tests {
		result := createFieldName(test.filename, test.method, test.field)
		if result != test.expected {
			t.Errorf("createFieldName(%s, %s, %s) = %s; expected %s",
				test.filename, test.method, test.field, result, test.expected)
		}
	}
}

func TestRoundTripFieldName(t *testing.T) {
	// Test that createFieldName and disassemblyFieldName are inverse operations
	testCases := []struct {
		filename string
		method   string
		field    string
	}{
		{"file.go", "main", "checked"},
		{"utils/helper.go", "testFunc", "notes"},
		{"src/app/routes.go", "setupRoutes", "description"},
	}

	for _, tc := range testCases {
		// Create field name
		fieldName := createFieldName(tc.filename, tc.method, tc.field)

		// Disassemble it
		filename, method, field := disassemblyFieldName(fieldName)

		// Should match original
		if filename != tc.filename || method != tc.method || field != tc.field {
			t.Errorf("Round trip failed for (%s, %s, %s): got (%s, %s, %s)",
				tc.filename, tc.method, tc.field, filename, method, field)
		}
	}
}

func TestGetFileID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"file.go", "file.go"},
		{"utils/helper.go", "utils.helper.go"},
		{"src/app/routes.go", "src.app.routes.go"},
		{"", ""},
		{"/path/to/file.go", ".path.to.file.go"},
	}

	for _, test := range tests {
		result := getFileID(test.input)
		if result != test.expected {
			t.Errorf("getFileID(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestParseEscapeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// html.EscapeString escapa más caracteres que la implementación anterior (más seguro)
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"<div>content</div>", "&lt;div&gt;content&lt;/div&gt;"},
		{"normal text", "normal text"},
		{"a < b && c > d", "a &lt; b &amp;&amp; c &gt; d"}, // & también se escapa
		{"", ""},
		{"&amp;", "&amp;amp;"}, // & existente se escapa correctamente
	}

	for _, test := range tests {
		result := parseEscapeHTML(test.input)
		if result != test.expected {
			t.Errorf("parseEscapeHTML(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}
