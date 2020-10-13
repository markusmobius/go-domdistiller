// ORIGINAL: javatest/TerminatingBlocksFinderTest.java

package english

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Filter_English_TerminatingBlocks_Positives(t *testing.T) {
	texts := []string{
		// Startswith cases.
		"comments foo", "© reuters", "© reuters foo bar", "please rate this",
		"please rate this foo", "post a comment", "post a comment foo", "123 comments",
		"9 comments foo", "1346213423 users responded in", "1346213423 users responded in foo",

		// Contains cases.
		"foo what you think... bar", "what you think...", "foo what you think...",
		"add your comment", "foo add your comment", "add comment bar", "reader views bar",
		"have your say bar", "foo reader comments", "foo rätta artikeln",

		// Equals cases.
		"thanks for your comments - this feedback is now closed",

		// Check some case insensitivity.
		"Thanks for your comments - this feedback is now closed", "Add Comment Bar",
		"READER VIEWS BAR", "Comments FOO",
	}

	terminatingBlocksFinder := NewTerminatingBlocksFinder()
	builder := testutil.NewTextBlockBuilder(stringutil.FastWordCounter{})

	for _, text := range texts {
		tb := builder.CreateForText(text)
		assert.True(t, terminatingBlocksFinder.isTerminating(tb))
	}
}

func Test_Filter_English_TerminatingBlocks_Negatives(t *testing.T) {
	texts := []string{
		// Startswith cases.
		"lcomments foo", "xd© reuters", "not please rate this", "xx post a comment",
		"users responded in", "123users responded in foo",

		// Contains cases.
		"what you think..", "addyour comment", "ad comment", "readerviews",

		// Equals cases.
		"thanks for your comments - this feedback is now closed foo",
		"foo thanks for your comments - this feedback is now closed",

		// Long case.
		"1 2 3 4 5 6 7 8 9 10 11 12 13 14 15",
	}

	terminatingBlocksFinder := NewTerminatingBlocksFinder()
	builder := testutil.NewTextBlockBuilder(stringutil.FastWordCounter{})

	for _, text := range texts {
		tb := builder.CreateForText(text)
		assert.False(t, terminatingBlocksFinder.isTerminating(tb))
	}
}

func Test_Filter_English_TerminatingBlocks_CommentsLink(t *testing.T) {
}
