// ORIGINAL: java/StringUtil.java

package stringutil

import (
	nurl "net/url"
	"path"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	rxAllDigits = regexp.MustCompile(`^\d+$`)
)

func IsStringAllWhitespace(str string) bool {
	for _, char := range str {
		if !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}

func IsStringAllDigit(str string) bool {
	return rxAllDigits.MatchString(str)
}

// =================================================================================
// Functions below these point are functions that doesn't exist in original code of
// Dom-Distiller, but useful for handling string.
// =================================================================================

// CharCount returns number of char in str.
func CharCount(str string) int {
	return utf8.RuneCountInString(str)
}

// CreateAbsoluteURL convert url to absolute path based on base.
// However, if url is prefixed with hash (#), the url won't be changed.
func CreateAbsoluteURL(url string, base *nurl.URL) string {
	if url == "" || base == nil {
		return url
	}

	// If it is hash tag, return as it is
	if strings.HasPrefix(url, "#") {
		return url
	}

	// If it is data URI, return as it is
	if strings.HasPrefix(url, "data:") {
		return url
	}

	// If it is javascript URI, return as it is
	if strings.HasPrefix(url, "javascript:") {
		return url
	}

	// If it is already an absolute URL, return as it is
	tmp, err := nurl.ParseRequestURI(url)
	if err == nil && tmp.Scheme != "" && tmp.Hostname() != "" {
		return url
	}

	// Otherwise, resolve against base URI.
	// Normalize URL first.
	if !strings.HasPrefix(url, "/") {
		url = path.Join(base.Path, url)
	}

	tmp, err = nurl.Parse(url)
	if err != nil {
		return url
	}

	return base.ResolveReference(tmp).String()
}
