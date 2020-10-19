// ORIGINAL: java/PageParameterDetector.java

package pattern

import "regexp"

const (
	PageParamPlaceholder = "[*!]"
)

var (
	rxNumber             = regexp.MustCompile(`(?i)(\d+)`)
	rxEndOrHasSHTML      = regexp.MustCompile(`(?i)(.s?html?)?$`)
	rxLastPathComponent  = regexp.MustCompile(`(?i)([^/]*)/$`)
	rxTrailingSlashHTML  = regexp.MustCompile(`(?i)(?:/|(.html?))$`)
	rxPageParamSeparator = regexp.MustCompile(`[-_;,]`)
)

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
