// ORIGINAL: java/webdocument/WebTag.java

package webdoc

// Tag represents HTML tags that need to be preserved over.
type Tag struct {
	Name string
	Type TagType
}

func (t *Tag) GenerateOutput(textOnly bool) string {
	if textOnly {
		return ""
	}

	if t.Type == TagStart {
		return "<" + t.Name + ">"
	}
	return "</" + t.Name + ">"
}
