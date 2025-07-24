// calculator/calculator_test.go
package calculator

import (
	"testing"
)

func TestNew(t *testing.T) {
	calc := New()
	if calc == nil {
		t.Error("New() should return a non-nil Calculator")
	}
}

func TestAdd(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 5, 3, 8},
		{"negative numbers", -5, -3, -8},
		{"mixed signs", 5, -3, 2},
		{"zeros", 0, 0, 0},
		{"with zero", 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 5, 3, 2},
		{"negative result", 3, 5, -2},
		{"with zero", 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 5, 3, 15},
		{"with zero", 5, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	calc := New()

	t.Run("normal division", func(t *testing.T) {
		result, err := calc.Divide(10, 2)
		if err != nil {
			t.Errorf("Divide(10, 2) returned unexpected error: %v", err)
		}
		if result != 5 {
			t.Errorf("Divide(10, 2) = %v, want 5", result)
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		_, err := calc.Divide(10, 0)
		if err == nil {
			t.Error("Divide(10, 0) should return an error")
		}
	})
}

func TestSquareRoot(t *testing.T) {
	calc := New()

	t.Run("positive number", func(t *testing.T) {
		result, err := calc.SquareRoot(9)
		if err != nil {
			t.Errorf("SquareRoot(9) returned unexpected error: %v", err)
		}
		if result != 3 {
			t.Errorf("SquareRoot(9) = %v, want 3", result)
		}
	})

	t.Run("negative number", func(t *testing.T) {
		_, err := calc.SquareRoot(-1)
		if err == nil {
			t.Error("SquareRoot(-1) should return an error")
		}
	})
}

func TestIsEven(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		n        int
		expected bool
	}{
		{"even positive", 4, true},
		{"odd positive", 5, false},
		{"zero", 0, true},
		{"even negative", -4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.IsEven(tt.n)
			if result != tt.expected {
				t.Errorf("IsEven(%d) = %v, want %v", tt.n, result, tt.expected)
			}
		})
	}
}

// Note: Factorial and IsPrime tests are intentionally missing to demonstrate coverage gaps

// Add to calculator/calculator_test.go
func TestFactorial(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		n        int
		expected int64
		hasError bool
	}{
		{"factorial of 0", 0, 1, false},
		{"factorial of 1", 1, 1, false},
		{"factorial of 5", 5, 120, false},
		{"negative number", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Factorial(tt.n)
			if tt.hasError {
				if err == nil {
					t.Errorf("Factorial(%d) should return error", tt.n)
				}
			} else {
				if err != nil {
					t.Errorf("Factorial(%d) returned unexpected error: %v", tt.n, err)
				}
				if result != tt.expected {
					t.Errorf("Factorial(%d) = %d, want %d", tt.n, result, tt.expected)
				}
			}
		})
	}
}

func TestIsPrime(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		n        int
		expected bool
	}{
		{"less than 2", 1, false},
		{"prime 2", 2, true},
		{"prime 3", 3, true},
		{"composite 4", 4, false},
		{"prime 17", 17, true},
		{"composite 15", 15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.IsPrime(tt.n)
			if result != tt.expected {
				t.Errorf("IsPrime(%d) = %v, want %v", tt.n, result, tt.expected)
			}
		})
	}
}

func TestPower(t *testing.T) {
	calc := New()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"2^3", 2, 3, 8},
		{"5^0", 5, 0, 1},
		{"10^2", 10, 2, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Power(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Power(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// Add to utils/strings_test.go
func TestIsPalindrome(t *testing.T) {
	utils := NewStringUtils()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple palindrome", "racecar", true},
		{"not palindrome", "hello", false},
		{"palindrome with spaces", "race car", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountVowels(t *testing.T) {
	utils := NewStringUtils()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"simple word", "hello", 2},
		{"all vowels", "aeiou", 5},
		{"no vowels", "xyz", 0},
		{"mixed case", "Hello World", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CountVowels(tt.input)
			if result != tt.expected {
				t.Errorf("CountVowels(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToTitleCase(t *testing.T) {
	utils := NewStringUtils()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "hello world", "Hello World"},
		{"uppercase", "HELLO WORLD", "Hello World"},
		{"mixed", "hELLo WoRLd", "Hello World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ToTitleCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToTitleCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveSpaces(t *testing.T) {
	utils := NewStringUtils()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with spaces", "hello world", "helloworld"},
		{"no spaces", "hello", "hello"},
		{"multiple spaces", "a b c d", "abcd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.RemoveSpaces(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveSpaces(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContainsOnlyLetters(t *testing.T) {
	utils := NewStringUtils()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"only letters", "hello", true},
		{"with numbers", "hello123", false},
		{"with spaces", "hello world", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ContainsOnlyLetters(tt.input)
			if result != tt.expected {
				t.Errorf("ContainsOnlyLetters(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

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
