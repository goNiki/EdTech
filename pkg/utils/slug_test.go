package utils_test

import (
	"testing"

	"edtech/pkg/utils"
)

func TestNormalizeSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Golang Basics 101", "golang-basics-101"},
		{"  Web_Dev & Python!! ", "web-dev-python"},
		{"HELLO--WORLD", "hello-world"},
		{"--leading-and-trailing--", "leading-and-trailing"},
		{"Special!@#%^&*()Chars", "special-chars"},
		{"Already-clean-slug", "already-clean-slug"},
		{"", "course"},
		{"   ---   ", "course"},
		{"multiple___underscores", "multiple-underscores"},
	}

	for _, tt := range tests {
		actual := utils.NormalizeSlug(tt.input)
		if actual != tt.expected {
			t.Errorf("NormalizeSlug(%q) = %q; expected %q", tt.input, actual, tt.expected)
		}
	}
}
