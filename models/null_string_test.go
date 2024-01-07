package models

import (
	"database/sql"
	"testing"
)

func TestNullString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    NullString
		expected string
	}{
		{
			name:     "Valid String",
			input:    NullString{sql.NullString{String: "Hello", Valid: true}},
			expected: `"Hello"`,
		},
		{
			name:     "Valid String",
			input:    NullString{sql.NullString{String: "", Valid: true}},
			expected: `""`,
		},
		{
			name:     "Invalid String",
			input:    NullString{sql.NullString{String: "", Valid: false}},
			expected: `""`,
		},
		// Add more test cases as needed for different scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalJSON()
			if err != nil {
				t.Errorf("Error during MarshalJSON: %s", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestNullString_Scan(t *testing.T) {
	// Test your Scan method here if needed
	// Depending on your use case and how Scan is used in your application
}
