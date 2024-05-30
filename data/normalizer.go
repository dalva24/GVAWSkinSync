package data

// package data contains everything stored on disk or remotely
// normalizer.go contains normalizing functions for string comparisons

import (
	"regexp"
	"strings"
)

var alphaNumericSpace = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
var alphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// NormTermTrimSpace removes all symbols and spaces, and make strings lowercase
func NormTermTrimSpace(str string) string {
	return strings.ToLower(strings.TrimSpace(alphaNumeric.ReplaceAllString(str, "")))
}

// NormTermLeaveSpace removes all symbols, leaves spaces, and make strings lowercase
func NormTermLeaveSpace(str string) string {
	return strings.ToLower(strings.TrimSpace(alphaNumericSpace.ReplaceAllString(str, "")))
}

// NormString removes all symbols, leaves spaces, and leave strings mixed-case
func NormString(str string) string {
	return strings.TrimSpace(alphaNumericSpace.ReplaceAllString(str, ""))
}
