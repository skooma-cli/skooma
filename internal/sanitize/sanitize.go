// Package sanitize contains functions for cleaning and normalizing user input.
package sanitize

import "strings"

// StripHTTPPrefix removes http:// or https:// from the beginning of a string.
func StripHTTPPrefix(str string) string {
	str = strings.TrimPrefix(str, "https://")
	str = strings.TrimPrefix(str, "http://")
	return str
}
