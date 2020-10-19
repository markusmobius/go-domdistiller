// ORIGINAL: javatest/QueryParamPagePatternTest.java

package detector_test

import (
	nurl "net/url"
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/pagination/detector"
	"github.com/stretchr/testify/assert"
)

func Test_Pagination_Detector_QPPP_IsPagingURL(t *testing.T) {
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryA=v1&queryB=4&queryC=v3",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryA=v1&queryB=growl&queryB=5&queryC=v3",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryA=v1&queryC=v3",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryB=2&queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?queryC=v3&queryC=v4",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3&queryC=v4")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b?page=3",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b/",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b.htm",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, true,
		"http://www.foo.com/a/b.html",
		"http://www.foo.com/a/b?page=[*!]")
	qpppAssertPagingURL(t, false,
		"http://www.foo.com/a/b?queryA=v1&queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, false,
		"http://www.foo.com/a/b?queryB=bar&queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3")
	qpppAssertPagingURL(t, false,
		"http://www.foo.com/a/b?queryA=v1",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3")
}

func Test_Pagination_Detector_QPPP_IsPagePatternValid(t *testing.T) {
	qpppAssertPagePatternValid(t, true,
		"http://www.google.com/forum-12",
		"http://www.google.com/forum-12?page=[*!]")
	qpppAssertPagePatternValid(t, true,
		"http://www.google.com/forum-12?sid=12345",
		"http://www.google.com/forum-12?page=[*!]&sort=d")
	qpppAssertPagePatternValid(t, false,
		"http://www.google.com/a/forum-12?sid=12345",
		"http://www.google.com/b/forum-12?page=[*!]&sort=d")
	qpppAssertPagePatternValid(t, false,
		"http://www.google.com/forum-11?sid=12345",
		"http://www.google.com/forum-12?page=[*!]&sort=d")
}

func qpppAssertPagingURL(t *testing.T, expected bool, strURL string, strPattern string) {
	pattern := createQueryParamPagePattern(strPattern)
	assert.NotNil(t, pattern)
	assert.Equal(t, expected, pattern.IsPagingURL(strURL))
}

func qpppAssertPagePatternValid(t *testing.T, expected bool, strURL string, strPattern string) {
	parsedURL, _ := nurl.ParseRequestURI(strURL)
	assert.NotNil(t, parsedURL)

	pattern := createQueryParamPagePattern(strPattern)
	assert.NotNil(t, pattern)

	assert.Equal(t, expected, pattern.IsValidFor(parsedURL))
}

func createQueryParamPagePattern(strPattern string) detector.PagePattern {
	url, err := nurl.ParseRequestURI(strPattern)
	if err != nil {
		return nil
	}

	for key, values := range url.Query() {
		for _, value := range values {
			if value == detector.PageParamPlaceholder {
				pattern, _ := detector.NewQueryParamPagePattern(url, key, "8")
				return pattern
			}
		}
	}

	return nil
}
