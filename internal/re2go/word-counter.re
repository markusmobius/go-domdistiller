/*!include:re2c "base.re" */

// Original pattern: [\x{3040}-\x{A4CF}]
func UseFullWordCounter(input string) bool {
	var cursor, marker int
	input += string(rune(0)) // add terminating null
	limit := len(input) - 1  // limit points at the terminating null
	_ = marker

	for { /*!use:re2c:base_template
		re2c:case-insensitive = 1;

		[\u3040-\uA4CF] { return true }
		* { continue }
		$ { return false }
		*/
	}
}

// Original pattern: [\x{AC00}-\x{D7AF}]
func UseLetterWordCounter(input string) bool {
	var cursor, marker int
	input += string(rune(0)) // add terminating null
	limit := len(input) - 1  // limit points at the terminating null
	_ = marker

	for { /*!use:re2c:base_template
		re2c:case-insensitive = 1;

		[\uAC00-\uD7AF] { return true }
		* { continue }
		$ { return false }
		*/
	}
}
