// ORIGINAL: java/OpenGraphProtocolParser.java

package opengraph

const (
	TitleProp                 = "title"
	TypeProp                  = "type"
	ImageProp                 = "image"
	URLProp                   = "url"
	DescriptionProp           = "description"
	SiteNameProp              = "site_name"
	ImageStructPropPfx        = "image:"
	ImageURLProp              = "image:url"
	ImageSecureURLProp        = "image:secure_url"
	ImageTypeProp             = "image:type"
	ImageWidthProp            = "image:width"
	ImageHeightProp           = "image:height"
	ProfileFirstnameProp      = "first_name"
	ProfileLastnameProp       = "last_name"
	ArticleSectionProp        = "section"
	ArticlePublishedTimeProp  = "published_time"
	ArticleModifiedTimeProp   = "modified_time"
	ArticleExpirationTimeProp = "expiration_time"
	ArticleAuthorProp         = "author"
	ProfileObjtype            = "profile"
	ArticleObjtype            = "article"

	doPrefixFiltering = true
)

var importantProperties = []struct {
	Name   string
	Prefix Prefix
	Type   string
}{
	{TitleProp, OG, ""},
	{TypeProp, OG, ""},
	{URLProp, OG, ""},
	{DescriptionProp, OG, ""},
	{SiteNameProp, OG, ""},
	{ImageProp, OG, "image"},
	{ImageStructPropPfx, OG, "image"},
	{ProfileFirstnameProp, Profile, "profile"},
	{ProfileLastnameProp, Profile, "profile"},
	{ArticleSectionProp, Article, "article"},
	{ArticlePublishedTimeProp, Article, "article"},
	{ArticleModifiedTimeProp, Article, "article"},
	{ArticleExpirationTimeProp, Article, "article"},
	{ArticleAuthorProp, Article, "article"},
}
