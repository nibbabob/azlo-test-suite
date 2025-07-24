// calculator/calculator.go
package calculator

import (
	"errors"
	"math"
)

// Calculator provides basic mathematical operations
type Calculator struct{}

// New creates a new Calculator instance
func New() *Calculator {
	return &Calculator{}
}

// Add returns the sum of two numbers
func (c *Calculator) Add(a, b float64) float64 {
	return a + b
}

// Subtract returns the difference between two numbers
func (c *Calculator) Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns the product of two numbers
func (c *Calculator) Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns the quotient of two numbers
func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Power returns a raised to the power of b
func (c *Calculator) Power(a, b float64) float64 {
	return math.Pow(a, b)
}

// SquareRoot returns the square root of a number
func (c *Calculator) SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("square root of negative number")
	}
	return math.Sqrt(a), nil
}

// Factorial returns the factorial of a non-negative integer
func (c *Calculator) Factorial(n int) (int64, error) {
	if n < 0 {
		return 0, errors.New("factorial of negative number")
	}
	if n == 0 || n == 1 {
		return 1, nil
	}

	result := int64(1)
	for i := 2; i <= n; i++ {
		result *= int64(i)
	}
	return result, nil
}

// IsEven checks if a number is even
func (c *Calculator) IsEven(n int) bool {
	return n%2 == 0
}

// IsPrime checks if a number is prime
func (c *Calculator) IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}
