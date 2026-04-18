package services

import (
	"regexp"
	"strings"
)

var nonAlnumRegex = regexp.MustCompile(`[^a-z0-9\-]+`)
var multipleDashRegex = regexp.MustCompile(`\-{2,}`)

func GenerateSlug(input string) string {
	slug := strings.ToLower(input)
	slug = strings.TrimSpace(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = nonAlnumRegex.ReplaceAllString(slug, "")
	slug = multipleDashRegex.ReplaceAllString(slug, "-")
	return slug
}