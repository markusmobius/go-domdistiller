// ORIGINAL: Part of javatest/StringUtilTest.java

package stringutil_test

import (
	nurl "net/url"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
)

func Test_StringUtil_IsStringAllWhitespace(t *testing.T) {
	assert.True(t, stringutil.IsStringAllWhitespace(""))
	assert.True(t, stringutil.IsStringAllWhitespace(" \t\r\n"))
	assert.True(t, stringutil.IsStringAllWhitespace(" \u00a0     \t\t\t"))
	assert.False(t, stringutil.IsStringAllWhitespace("a"))
	assert.False(t, stringutil.IsStringAllWhitespace("     a  "))
	assert.False(t, stringutil.IsStringAllWhitespace("\u00a0\u0460"))
	assert.False(t, stringutil.IsStringAllWhitespace("\n\t_ "))
}

// =================================================================================
// Tests below these point are test for function that doesn't exist in original code
// =================================================================================

func Test_StringUtil_CreateAbsoluteURL(t *testing.T) {
	relURL1 := "#here"
	relURL2 := "/test/123"
	relURL3 := "test/123"
	relURL4 := "//www.google.com"
	relURL5 := "https://www.google.com"
	relURL6 := "ftp://ftp.server.com"
	relURL7 := "www.google.com"
	relURL8 := "http//www.google.com"
	relURL9 := "../hello/relative"

	absURL1 := "#here"
	absURL2 := "http://example.com/test/123"
	absURL3 := "http://example.com/page/test/123"
	absURL4 := "http://www.google.com"
	absURL5 := "https://www.google.com"
	absURL6 := "ftp://ftp.server.com"
	absURL7 := "http://example.com/page/www.google.com"
	absURL8 := "http://example.com/page/http/www.google.com"
	absURL9 := "http://example.com/hello/relative"

	baseURL, _ := nurl.ParseRequestURI("http://example.com/page/")
	assert.Equal(t, absURL1, stringutil.CreateAbsoluteURL(relURL1, baseURL))
	assert.Equal(t, absURL2, stringutil.CreateAbsoluteURL(relURL2, baseURL))
	assert.Equal(t, absURL3, stringutil.CreateAbsoluteURL(relURL3, baseURL))
	assert.Equal(t, absURL4, stringutil.CreateAbsoluteURL(relURL4, baseURL))
	assert.Equal(t, absURL5, stringutil.CreateAbsoluteURL(relURL5, baseURL))
	assert.Equal(t, absURL6, stringutil.CreateAbsoluteURL(relURL6, baseURL))
	assert.Equal(t, absURL7, stringutil.CreateAbsoluteURL(relURL7, baseURL))
	assert.Equal(t, absURL8, stringutil.CreateAbsoluteURL(relURL8, baseURL))
	assert.Equal(t, absURL9, stringutil.CreateAbsoluteURL(relURL9, baseURL))
}
