// ORIGINAL: Protobuf model in proto/dom_distiller.proto

package data

type PaginationInfo struct {
	NextPage string
	PrevPage string
}

// MarkupArticle is object to contains the properties of an article document.
type MarkupArticle struct {
	PublishedTime  string
	ModifiedTime   string
	ExpirationTime string
	Section        string
	Authors        []string
}

// MarkupImage is used to contains the properties of an image in the document.
type MarkupImage struct {
	Root      string
	URL       string
	SecureURL string
	Type      string
	Caption   string
	Width     int
	Height    int
}

type MarkupInfo struct {
	Title       string
	Type        string
	URL         string
	Description string
	Publisher   string
	Copyright   string
	Author      string
	Article     MarkupArticle
	Images      []MarkupImage
}
