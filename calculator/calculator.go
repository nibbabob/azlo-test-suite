// calculator/calculator.go
package calculator

import (
	"errors"
	"math"
	"strings"
	"unicode"
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

// StringUtils provides utility functions for string operations
type StringUtils struct{}

// NewStringUtils creates a new StringUtils instance
func NewStringUtils() *StringUtils {
	return &StringUtils{}
}

// Reverse returns the reversed version of a string
func (s *StringUtils) Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome checks if a string is a palindrome
func (s *StringUtils) IsPalindrome(str string) bool {
	cleaned := strings.ToLower(strings.ReplaceAll(str, " ", ""))
	return cleaned == s.Reverse(cleaned)
}

// CountWords returns the number of words in a string
func (s *StringUtils) CountWords(str string) int {
	fields := strings.Fields(str)
	return len(fields)
}

// CountVowels returns the number of vowels in a string
func (s *StringUtils) CountVowels(str string) int {
	vowels := "aeiouAEIOU"
	count := 0
	for _, char := range str {
		if strings.ContainsRune(vowels, char) {
			count++
		}
	}
	return count
}

// ToTitleCase converts a string to title case
func (s *StringUtils) ToTitleCase(str string) string {
	return strings.Title(strings.ToLower(str))
}

// RemoveSpaces removes all spaces from a string
func (s *StringUtils) RemoveSpaces(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

// IsNumeric checks if a string contains only numeric characters
func (s *StringUtils) IsNumeric(str string) bool {
	if str == "" {
		return false
	}
	for _, char := range str {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

// ContainsOnlyLetters checks if a string contains only letters
func (s *StringUtils) ContainsOnlyLetters(str string) bool {
	if str == "" {
		return false
	}
	for _, char := range str {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}
