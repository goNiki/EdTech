package slug

import (
	"regexp"
	"strings"
)

var (
	nonSlugCharsRegex = regexp.MustCompile(`[^a-z0-9-]+`)
	multipleDashRegex = regexp.MustCompile(`-+`)
)

// NormalizeSlug converts any string to a safe URL-friendly slug containing only [a-z0-9-].
func NormalizeSlug(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	s = nonSlugCharsRegex.ReplaceAllString(s, "-")
	s = multipleDashRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "course"
	}
	return s
}
