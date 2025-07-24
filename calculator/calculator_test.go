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
