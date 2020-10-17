// ORIGINAL: javatest/QueryParamPagePatternTest.java

package detector_test

import (
	nurl "net/url"
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/pagination/detector"
	"github.com/stretchr/testify/assert"
)

func Test_Pagination_Detector_QPPP_IsPagingURL(t *testing.T) {
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?queryA=v1&queryB=4&queryC=v3",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?queryA=v1&queryB=growl&queryB=5&queryC=v3",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?queryA=v1&queryC=v3",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?queryB=2&queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?queryC=v3&queryC=v4",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3&queryC=v4"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b",
		"http://www.foo.com/a/b?page=[*!]"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?page=3",
		"http://www.foo.com/a/b?page=[*!]"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b/",
		"http://www.foo.com/a/b?page=[*!]"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b.htm",
		"http://www.foo.com/a/b?page=[*!]"))
	assert.True(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b.html",
		"http://www.foo.com/a/b?page=[*!]"))
	assert.False(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?queryA=v1&queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3"))
	assert.False(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?queryB=bar&queryC=v3",
		"http://www.foo.com/a/b?queryB=[*!]&queryC=v3"))
	assert.False(t, qpppIsPagingURL(t,
		"http://www.foo.com/a/b?queryA=v1",
		"http://www.foo.com/a/b?queryA=v1&queryB=[*!]&queryC=v3"))
}

func Test_Pagination_Detector_QPPP_IsPagePatternValid(t *testing.T) {
	assert.True(t, qpppIsPagePatternValid(t,
		"http://www.google.com/forum-12",
		"http://www.google.com/forum-12?page=[*!]"))
	assert.True(t, qpppIsPagePatternValid(t,
		"http://www.google.com/forum-12?sid=12345",
		"http://www.google.com/forum-12?page=[*!]&sort=d"))
	assert.False(t, qpppIsPagePatternValid(t,
		"http://www.google.com/a/forum-12?sid=12345",
		"http://www.google.com/b/forum-12?page=[*!]&sort=d"))
	assert.False(t, qpppIsPagePatternValid(t,
		"http://www.google.com/forum-11?sid=12345",
		"http://www.google.com/forum-12?page=[*!]&sort=d"))
}

func qpppIsPagingURL(t *testing.T, strURL string, strPattern string) bool {
	pattern := createQueryParamPagePattern(strPattern)
	assert.NotNil(t, pattern)
	return pattern.IsPagingURL(strURL)
}

func qpppIsPagePatternValid(t *testing.T, strURL string, strPattern string) bool {
	parsedURL, _ := nurl.ParseRequestURI(strURL)
	assert.NotNil(t, parsedURL)

	pattern := createQueryParamPagePattern(strPattern)
	assert.NotNil(t, pattern)

	return pattern.IsValidFor(parsedURL)
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
