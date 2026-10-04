package utils

import (
	"edtech/pkg/slug"
)

// NormalizeSlug delegates to pkg/slug.NormalizeSlug for backward compatibility.
func NormalizeSlug(input string) string {
	return slug.NormalizeSlug(input)
}
