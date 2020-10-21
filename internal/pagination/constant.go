// ORIGINAL: java/PageParameterParser.java

package pagination

import "regexp"

const (
	// If the numeric value of a link's anchor text is greater than this number,
	// we don't think it represents the page number of the link.
	MaxNumForPageParam = 100
)

var (
	rxNumber        = regexp.MustCompile(`\d`)
	rxNumberAtStart = regexp.MustCompile(`^\d+`)

	// Regex for page number finder
	rxLinkNumberCleaner    = regexp.MustCompile(`[()\[\]{}]`)
	rxInvalidParentWrapper = regexp.MustCompile(`(?i)(body)|(html)`)
	rxTerms                = regexp.MustCompile(`(?i)(\S*[\w\x{00C0}-\x{1FFF}\x{2C00}-\x{D7FF}]\S*)`)
	rxSurroundingDigits    = regexp.MustCompile(`(?i)^[\W_]*(\d+)[\W_]*$`)

	// Regex for prev next finder
	rxNextLink       = regexp.MustCompile(`(?i)(next|weiter|continue|>([^\|]|$)|»([^\|]|$))`)
	rxPrevLink       = regexp.MustCompile(`(?i)(prev|early|old|new|<|«)`)
	rxPositive       = regexp.MustCompile(`(?i)article|body|content|entry|hentry|main|page|pagination|post|text|blog|story`)
	rxNegative       = regexp.MustCompile(`(?i)combx|comment|com-|contact|foot|footer|footnote|masthead|media|meta|outbrain|promo|related|shoutbox|sidebar|sponsor|shopping|tags|tool|widget`)
	rxExtraneous     = regexp.MustCompile(`(?i)print|archive|comment|discuss|e[\-]?mail|share|reply|all|login|sign|single|as one|article|post|篇`)
	rxPagination     = regexp.MustCompile(`(?i)pag(e|ing|inat)`)
	rxLinkPagination = regexp.MustCompile(`(?i)p(a|g|ag)?(e|ing|ination)?(=|\/)[0-9]{1,2}$`)
	rxFirstLast      = regexp.MustCompile(`(?i)(first|last)`)
)
