package sanitize

import (
	"html"
	"regexp"
	"strings"
)

var (
	htmlTagRegex    = regexp.MustCompile(`(?i)<[^>]*>`)
	scriptTagRegex  = regexp.MustCompile(`(?i)<\s*script[^>]*>.*?<\s*/\s*script\s*>`)
	jsProtocolRegex = regexp.MustCompile(`(?i)javascript:`)
	sqlInjectRegex  = regexp.MustCompile(`(?i)(--|;|/\*|\*/|@@|char\s*\(|nchar\s*\(|varchar\s*\(|alter\s+table|drop\s+table|exec\s*\(|execute\s*\(|union\s+select)`)
)

// CleanText strips HTML/script tags, escapes HTML entities, and trims whitespace.
func CleanText(s string) string {
	s = scriptTagRegex.ReplaceAllString(s, "")
	s = htmlTagRegex.ReplaceAllString(s, "")
	s = jsProtocolRegex.ReplaceAllString(s, "")
	s = html.EscapeString(s)
	return strings.TrimSpace(s)
}

// ContainsInjectionPattern returns true if the string matches dangerous SQL/script injection signatures.
func ContainsInjectionPattern(s string) bool {
	return sqlInjectRegex.MatchString(s) || scriptTagRegex.MatchString(s) || jsProtocolRegex.MatchString(s)
}
