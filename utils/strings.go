// utils/strings.go
package utils

import (
	"strings"
	"unicode"
)

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
