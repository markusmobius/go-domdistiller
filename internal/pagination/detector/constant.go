// ORIGINAL: java/PageParameterDetector.java

package detector

import "regexp"

const (
	PageParamPlaceholder = "[*!]"
)

var rxTrailingSlashHTML = regexp.MustCompile(`(?i)(?:/|(.html?))$`)

var badPageParamNames = map[string]struct{}{
	"baixar-gratis":  {},
	"category":       {},
	"content":        {},
	"day":            {},
	"date":           {},
	"definition":     {},
	"etiket":         {},
	"film-seyret":    {},
	"key":            {},
	"keys":           {},
	"keyword":        {},
	"label":          {},
	"news":           {},
	"q":              {},
	"query":          {},
	"rating":         {},
	"s":              {},
	"search":         {},
	"seasons":        {},
	"search_keyword": {},
	"search_query":   {},
	"sortby":         {},
	"subscriptions":  {},
	"tag":            {},
	"tags":           {},
	"video":          {},
	"videos":         {},
	"w":              {},
	"wiki":           {},
}
