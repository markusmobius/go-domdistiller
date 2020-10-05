// ORIGINAL: java/StringUtil.java

package stringutil

import "unicode/utf8"

// =================================================================================
// Functions below these point are functions that doesn't exist in original code of
// Dom-Distiller, but useful for handling string.
// =================================================================================

// CharCount returns number of char in str.
func CharCount(str string) int {
	return utf8.RuneCountInString(str)
}
