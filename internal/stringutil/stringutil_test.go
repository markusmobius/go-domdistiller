// ORIGINAL: Part of javatest/StringUtilTest.java

package stringutil_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/markusmobius/go-domdistiller/internal/stringutil"
)

func Test_IsStringAllWhitespace(t *testing.T) {
	assert.True(t, stringutil.IsStringAllWhitespace(""))
	assert.True(t, stringutil.IsStringAllWhitespace(" \t\r\n"))
	assert.True(t, stringutil.IsStringAllWhitespace(" \u00a0     \t\t\t"))
	assert.False(t, stringutil.IsStringAllWhitespace("a"))
	assert.False(t, stringutil.IsStringAllWhitespace("     a  "))
	assert.False(t, stringutil.IsStringAllWhitespace("\u00a0\u0460"))
	assert.False(t, stringutil.IsStringAllWhitespace("\n\t_ "))
}
