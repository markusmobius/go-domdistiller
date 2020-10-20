// ORIGINAL: java/PageParameterParser.java

package pagination

import "regexp"

const (
	// If the numeric value of a link's anchor text is greater than this number,
	// we don't think it represents the page number of the link.
	MaxNumForPageParam = 100
)

var (
	rxNumber               = regexp.MustCompile(`\d`)
	rxLinkNumberCleaner    = regexp.MustCompile(`[()\[\]{}]`)
	rxInvalidParentWrapper = regexp.MustCompile(`(?i)(body)|(html)`)
	rxTerms                = regexp.MustCompile(`(?i)(\S*[\w\x{00C0}-\x{1FFF}\x{2C00}-\x{D7FF}]\S*)`)
	rxSurroundingDigits    = regexp.MustCompile(`(?i)^[\W_]*(\d+)[\W_]*$`)
)
