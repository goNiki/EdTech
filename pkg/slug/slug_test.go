package slug_test

import (
	"testing"

	"edtech/pkg/slug"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "regular english title",
			input:    "Introduction to Go Programming",
			expected: "introduction-to-go-programming",
		},
		{
			name:     "special characters and spaces",
			input:    "Golang: Best Practices & Tips! (2024)",
			expected: "golang-best-practices-tips-2024",
		},
		{
			name:     "underscores and multiple dashes",
			input:    "my__course___name---here",
			expected: "my-course-name-here",
		},
		{
			name:     "empty string defaults to course",
			input:    "   ",
			expected: "course",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, slug.NormalizeSlug(tt.input))
		})
	}
}
