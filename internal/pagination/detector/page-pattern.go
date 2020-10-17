// ORIGINAL: java/PageParameterDetector.java

package detector

import (
	nurl "net/url"
)

// PagePattern is the interface that page pattern handlers must implement to detect
// page parameter from potential pagination URLs.
type PagePattern interface {
	// String returns the string of the URL page pattern.
	String() string

	// PageNumber returns the page number extracted from the URL during creation of
	// object that implements this interface.
	PageNumber() int

	// IsValidFor validates this page pattern according to the current document URL
	// through a pipeline of rules. Returns true if page pattern is valid.
	// docUrl is the current document URL.
	IsValidFor(docURL *nurl.URL) bool

	// IsPagingURL returns true if a URL matches this page pattern based on a pipeline of rules.
	// url is the URL to evaluate.
	IsPagingURL(url string) bool
}
