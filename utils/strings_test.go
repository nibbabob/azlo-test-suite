// utils/strings_test.go
package utils

import (
	"testing"
)

func TestNewStringUtils(t *testing.T) {
	utils := NewStringUtils()
	if utils == nil {
		t.Error("NewStringUtils() should return a non-nil StringUtils")
	}
}

func TestReverse(t *testing.T) {
	utils := NewStringUtils()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello", "olleh"},
		{"empty string", "", ""},
		{"single character", "a", "a"},
		{"with spaces", "hello world", "dlrow olleh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.Reverse(tt.input)
			if result != tt.expected {
				t.Errorf("Reverse(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	utils := NewStringUtils()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"single word", "hello", 1},
		{"multiple words", "hello world", 2},
		{"empty string", "", 0},
		{"extra spaces", "hello  world   test", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CountWords(tt.input)
			if result != tt.expected {
				t.Errorf("CountWords(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	utils := NewStringUtils()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"numeric string", "12345", true},
		{"empty string", "", false},
		{"mixed", "123abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("IsNumeric(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Note: Several methods are not tested to show coverage gaps:
// - IsPalindrome
// - CountVowels
// - ToTitleCase
// - RemoveSpaces
// - ContainsOnlyLetters
