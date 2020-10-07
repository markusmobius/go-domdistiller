// ORIGINAL: java/StringUtil.java

package stringutil

import (
	"regexp"
	"unicode"
	"unicode/utf8"
)

var (
	rxNonWhitespace = regexp.MustCompile(`(?i)\S`)
)

func IsStringAllWhitespace(str string) bool {
	for _, char := range str {
		if !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}

// =================================================================================
// Functions below these point are functions that doesn't exist in original code of
// Dom-Distiller, but useful for handling string.
// =================================================================================

// CharCount returns number of char in str.
func CharCount(str string) int {
	return utf8.RuneCountInString(str)
}
