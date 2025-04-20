// file: helper/slug.go
package helper

import (
	"regexp"
	"strings"
)

var nonAlnumRegex = regexp.MustCompile(`[^a-z0-9\-]+`)
var multipleDashRegex = regexp.MustCompile(`\-{2,}`)

func GenerateSlug(input string) string {
	// 1. Lowercase
	slug := strings.ToLower(input)

	// 2. Trim whitespace
	slug = strings.TrimSpace(slug)

	// 3. Replace spaces with dash
	slug = strings.ReplaceAll(slug, " ", "-")

	// 4. Remove non-alphanumeric characters except dash
	slug = nonAlnumRegex.ReplaceAllString(slug, "")

	// 5. Replace multiple dashes with single dash
	slug = multipleDashRegex.ReplaceAllString(slug, "-")

	return slug
}
