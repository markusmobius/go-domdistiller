// ORIGINAL: java/SchemaOrgParser.java

package schemaorg

const (
	NameProp            = "name"
	URLProp             = "url"
	DescriptionProp     = "description"
	ImageProp           = "image"
	HeadlineProp        = "headline"
	PublisherProp       = "publisher"
	CopyrightHolderProp = "copyrightHolder"
	CopyrightYearProp   = "copyrightYear"
	ContentURLProp      = "contentUrl"
	EncodingFormatProp  = "encodingFormat"
	CaptionProp         = "caption"
	RepresentativeProp  = "representativeOfPage"
	WidthProp           = "width"
	HeightProp          = "height"
	DatePublishedProp   = "datePublished"
	DateModifiedProp    = "dateModified"
	AuthorProp          = "author"
	CreatorProp         = "creator"
	SectionProp         = "articleSection"
	AssociatedMediaProp = "associatedMedia"
	EncodingProp        = "encoding"
	FamilyNameProp      = "familyName"
	GivenNameProp       = "givenName"
	LegalNameProp       = "legalName"
	AuthorRel           = "author"
)

type SchemaType uint

const (
	Unsupported SchemaType = iota
	Image
	Article
	Person
	Organization
)

var schemaTypeURLs = map[string]SchemaType{
	"http://schema.org/ImageObject":             Image,
	"http://schema.org/Article":                 Article,
	"http://schema.org/BlogPosting":             Article,
	"http://schema.org/NewsArticle":             Article,
	"http://schema.org/ScholarlyArticle":        Article,
	"http://schema.org/TechArticle":             Article,
	"http://schema.org/Person":                  Person,
	"http://schema.org/Organization":            Organization,
	"http://schema.org/Corporation":             Organization,
	"http://schema.org/EducationalOrganization": Organization,
	"http://schema.org/GovernmentOrganization":  Organization,
	"http://schema.org/NGO":                     Organization,
}

// The key for `tagAttributeMap` is the tag name, while the entry value is an
// array of attributes in the specified tag from which to extract information:
// - 0th attribute: contains the value for the property specified in itemprop
// - 1st attribute: if available, contains the value for the author property.
var tagAttributeMap = map[string]string{
	"img":    "src",
	"audio":  "src",
	"embed":  "src",
	"iframe": "src",
	"source": "src",
	"track":  "src",
	"video":  "src",
	"a":      "href",
	"link":   "href",
	"area":   "href",
	"meta":   "content",
	"time":   "datetime",
	"object": "data",
	"data":   "value",
	"meter":  "value",
}
