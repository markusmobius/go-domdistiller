// ORIGINAL: java/webdocument/WebTag.java

package webdoc

// Tag represents HTML tags that need to be preserved over.
type Tag struct {
	BaseElement
	Name string
	Type TagType
}

func NewTag(name string, tagType TagType) *Tag {
	return &Tag{Name: name, Type: tagType}
}

func (t *Tag) ElementType() string {
	return "tag"
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
