/*!include:re2c "base.re" */

// Original pattern: (?i)(^(comments|© reuters|please rate this|post a comment|\d+\s+(comments|users responded in))|what you think\.\.\.|add your comment|add comment|reader views|have your say|reader comments|rätta artikeln|^thanks for your comments - this feedback is now closed$)
// For convenience it will be separated into 3 regexes:
// - ^(comments|© reuters|please rate this|post a comment|\d+\s+(comments|users responded in))
// - what you think\.\.\.|add your comment|add comment|reader views|have your say|reader comments|rätta artikeln
// - ^thanks for your comments - this feedback is now closed$

// Handle (^(comments|© reuters|please rate this|post a comment|\d+\s+(comments|users responded in))
func isTerminatingBlocks1(input string) bool {
	var cursor, marker int
	input += string(rune(0)) // add terminating null
	limit := len(input) - 1  // limit points at the terminating null
	_ = marker

	for { /*!use:re2c:base_template
		re2c:case-insensitive = 1;

		tb1a = comments;
		tb1b = ©[ ]reuters;
		tb1c = please[ ]rate[ ]this;
		tb1d = post[ ]a[ ]comment;

		tb1eQuant1 = [0-9]+;
		tb1eQuant2 = [0-9]+[\t\n\f\r ]+;
		tb1e       = [0-9]+[\t\n\f\r ]+(comments|users[ ]responded[ ]in);

		{tb1a} { return true }
		{tb1b} { return true }
		{tb1c} { return true }
		{tb1d} { return true }
		{tb1e} { return true }

		{tb1eQuant1} { return false }
		{tb1eQuant2} { return false }

		* { return false }
		$ { return false }
		*/
	}
}

// Handle what you think\.\.\.|add your comment|add comment|reader views|have your say|reader comments|rätta artikeln
func isTerminatingBlocks2(input string) bool {
	var cursor, marker int
	input += string(rune(0)) // add terminating null
	limit := len(input) - 1  // limit points at the terminating null
	_ = marker

	for { /*!use:re2c:base_template
		re2c:case-insensitive = 1;

		tb2a = what[ ]you[ ]think[.]{3};
		tb2b = add[ ]your[ ]comment;
		tb2c = add[ ]comment;
		tb2d = reader[ ]views;
		tb2e = have[ ]your[ ]say;
		tb2f = reader[ ]comments;
		tb2g = r[äÄ]tta[ ]artikeln;

		{tb2a} { return true }
		{tb2b} { return true }
		{tb2c} { return true }
		{tb2d} { return true }
		{tb2e} { return true }
		{tb2f} { return true }
		{tb2g} { return true }

		* { continue }
		$ { return false }
		*/
	}
}

// Handle ^thanks for your comments - this feedback is now closed$
func isTerminatingBlocks3(input string) bool {
	var cursor, marker int
	input += string(rune(0)) // add terminating null
	limit := len(input) - 1  // limit points at the terminating null
	_ = marker

	var found bool
	for { /*!use:re2c:base_template
		re2c:case-insensitive = 1;

		thanks[ ]for[ ]your[ ]comments[ ]-[ ]this[ ]feedback[ ]is[ ]now[ ]closed {
			found = true
			continue
		}

		* { return false }
		$ { return found }
		*/
	}
}

func IsTerminatingBlocks(input string) bool {
	return isTerminatingBlocks1(input) ||
		isTerminatingBlocks2(input) ||
		isTerminatingBlocks3(input)
}