// ORIGINAL: javatest/SchemaOrgParserAccessorTest.java

package schemaorg_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/go-shiori/dom"
	"github.com/markusmobius/go-domdistiller/internal/markup/schemaorg"
	"github.com/markusmobius/go-domdistiller/internal/testutil"
	"golang.org/x/net/html"
)

func Test_SchemaOrg_ImageWithEmbeddedPublisher(t *testing.T) {
	expectedURL := "http://dummy/Test_SchemaOrg_image_with_embedded_item.html"
	expectedFormat := "jpeg"
	expectedCaption := "A test for IMAGE with embedded publisher"
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/ImageObject">` +
		`	<h1 itemprop="headline">Testcase for IMAGE` +
		`	</h1>` +
		`	<h2 itemprop="description">Testing IMAGE with embedded publisher` +
		`	</h2>` +
		`	<a itemprop="contentUrl" href="` + expectedURL + `">test results` +
		`	</a>` +
		`	<div id="2" itemscope itemtype="http://schema.org/Organization"` +
		`		 itemprop="publisher">Publisher: ` +
		`		<span itemprop="name">Whatever Image Incorporated` +
		`		</span>` +
		`	</div>` +
		`	<div id="3">` +
		`		<span itemprop="copyrightYear">1999-2022` +
		`		</span>` +
		`		<span itemprop="copyrightHolder">Whoever Image Copyrighted` +
		`		</span>` +
		`	</div>` +
		`	<span itemprop="encodingFormat">` + expectedFormat +
		`	</span>` +
		`	<span itemprop="caption">` + expectedCaption +
		`	</span>` +
		`	<meta itemprop="representativeOfPage" content="true">` +
		`	<meta itemprop="width" content="600">` +
		`	<meta itemprop="height" content="400">` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "", parser.Type())
	assert.Equal(t, "", parser.Title())
	assert.Equal(t, "", parser.Description())
	assert.Equal(t, "", parser.URL())
	assert.Equal(t, "", parser.Publisher())
	assert.Nil(t, parser.Article())
	assert.Equal(t, "", parser.Author())
	assert.Equal(t, "", parser.Copyright())

	images := parser.Images()
	assert.Equal(t, 1, len(images))

	image := images[0]
	assert.Equal(t, expectedURL, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, expectedFormat, image.Type)
	assert.Equal(t, expectedCaption, image.Caption)
	assert.Equal(t, 600, image.Width)
	assert.Equal(t, 400, image.Height)
}

func Test_SchemaOrg_2Images(t *testing.T) {
	expectedURL1 := "http://dummy/Test_SchemaOrg_1st image.html"
	expectedPublisher1 := "Whatever 1st Image Incorporated"
	expectedFormat1 := "jpeg"
	expectedCaption1 := "A test for 1st IMAGE"
	expectedURL2 := "http://dummy/Test_SchemaOrg_2nd image.html"
	expectedFormat2 := "gif"
	expectedCaption2 := "A test for 2nd IMAGE"
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/ImageObject">` +
		`	<h1 itemprop="headline">Testcase for 1st IMAGE` +
		`	</h1>` +
		`	<h2 itemprop="description">Testing 1st IMAGE` +
		`	</h2>` +
		`	<a itemprop="contentUrl" href="` + expectedURL1 + `">1st test results` +
		`	</a>` +
		`	<div id="2" itemprop="publisher">` + expectedPublisher1 +
		`	</div>` +
		`	<div id="3">` +
		`		<span itemprop="copyrightYear">1000-1999` +
		`		</span>` +
		`		<span itemprop="copyrightHolder">Whoever 1st Image Copyrighted` +
		`		</span>` +
		`	</div>` +
		`	<span itemprop="encodingFormat">` + expectedFormat1 +
		`	</span>` +
		`	<span itemprop="caption">` + expectedCaption1 +
		`	</span>` +
		`	<meta itemprop="representativeOfPage" content="false">` +
		`	<meta itemprop="width" content="400">` +
		`	<meta itemprop="height" content="300">` +
		`</div>` +
		`<div id="4" itemscope itemtype="http://schema.org/ImageObject">` +
		`	<h3 itemprop="headline">Testcase for 2nd IMAGE` +
		`	</h3>` +
		`	<h4 itemprop="description">Testing 2nd IMAGE` +
		`	</h4>` +
		`	<a itemprop="contentUrl" href="` + expectedURL2 + `">2nd test results` +
		`	</a>` +
		`	<div id="5" itemprop="publisher">Whatever 2nd Image Incorporated` +
		`	</div>` +
		`	<div id="6">` +
		`		<span itemprop="copyrightYear">2000-2999` +
		`		</span>` +
		`		<span itemprop="copyrightHolder">Whoever 2nd Image Copyrighted` +
		`		</span>` +
		`	</div>` +
		`	<span itemprop="encodingFormat">` + expectedFormat2 +
		`	</span>` +
		`	<span itemprop="caption">` + expectedCaption2 +
		`	</span>` +
		`	<meta itemprop="representativeOfPage" content="true">` +
		`	<meta itemprop="width" content="1000">` +
		`	<meta itemprop="height" content="600">` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)

	// The basic properties of Thing should be from the first
	// image that was inserted.
	assert.Equal(t, "", parser.Type())
	assert.Equal(t, "", parser.Title())
	assert.Equal(t, "", parser.Description())
	assert.Equal(t, "", parser.URL())
	assert.Equal(t, "", parser.Publisher())
	assert.Nil(t, parser.Article())
	assert.Equal(t, "", parser.Author())
	assert.Equal(t, "", parser.Copyright())

	images := parser.Images()
	assert.Equal(t, 2, len(images))

	// The 2nd image that was inserted is representative of page, so the
	// images should be swapped in `images`.
	image := images[0]
	assert.Equal(t, expectedURL2, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, expectedFormat2, image.Type)
	assert.Equal(t, expectedCaption2, image.Caption)
	assert.Equal(t, 1000, image.Width)
	assert.Equal(t, 600, image.Height)

	image = images[1]
	assert.Equal(t, expectedURL1, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, expectedFormat1, image.Type)
	assert.Equal(t, expectedCaption1, image.Caption)
	assert.Equal(t, 400, image.Width)
	assert.Equal(t, 300, image.Height)
}

func Test_SchemaOrg_ArticleWithEmbeddedAuthorAndPublisher(t *testing.T) {
	expectedTitle := "Testcase for ARTICLE"
	expectedDescription := "Testing ARTICLE with embedded author and publisher"
	expectedUrl := "http://dummy/Test_SchemaOrg_article_with_embedded_items.html"
	expectedImage := "http://dummy/Test_SchemaOrg_article_with_embedded_items.jpeg"
	expectedAuthor := "Whoever authored"
	expectedPublisher := "Whatever Article Incorporated"
	expectedDatePublished := "April 15, 2014"
	expectedTimeModified := "2014-04-16T23:59"
	expectedCopyrightYear := "2000-2014"
	expectedCopyrightHolder := "Whoever Article Copyrighted"
	expectedSection := "Romance thriller"
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/Article">` +
		`	<h1 itemprop="headline">` + expectedTitle +
		`	</h1>` +
		`	<h2 itemprop="description">` + expectedDescription +
		`	</h2>` +
		`	<a itemprop="url" href="` + expectedUrl + `">test results` +
		`	</a>` +
		`	<img itemprop="image" src="` + expectedImage + `">` +
		`	<div id="2" itemscope itemtype="http://schema.org/Person"` +
		`		 itemprop="author">Author: ` +
		`		<span itemprop="name">` + expectedAuthor +
		`		</span>` +
		`	</div>` +
		`	<div id="3" itemscope itemtype="http://schema.org/Organization"` +
		`		 itemprop="publisher">Publisher: ` +
		`		<span itemprop="name">` + expectedPublisher +
		`		</span>` +
		`	</div>` +
		`	<span itemprop="datePublished">` + expectedDatePublished +
		`	</span>` +
		`	<time itemprop="dateModified" datetime="` + expectedTimeModified +
		`		">April 16, 2014 11:59pm` +
		`	</time>` +
		`	<span itemprop="copyrightYear">` + expectedCopyrightYear +
		`	</span>` +
		`	<span itemprop="copyrightHolder">` + expectedCopyrightHolder +
		`	</span>` +
		`	<span itemprop="articleSection">` + expectedSection +
		`	</span>` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "Article", parser.Type())
	assert.Equal(t, expectedTitle, parser.Title())
	assert.Equal(t, expectedDescription, parser.Description())
	assert.Equal(t, expectedUrl, parser.URL())
	assert.Equal(t, expectedAuthor, parser.Author())
	assert.Equal(t, expectedPublisher, parser.Publisher())
	assert.Equal(t, "Copyright "+expectedCopyrightYear+" "+expectedCopyrightHolder, parser.Copyright())

	images := parser.Images()
	assert.Equal(t, 1, len(images))
	assert.Equal(t, expectedImage, images[0].URL)

	article := parser.Article()
	assert.NotNil(t, article)
	assert.Equal(t, expectedDatePublished, article.PublishedTime)
	assert.Equal(t, expectedTimeModified, article.ModifiedTime)
	assert.Equal(t, "", article.ExpirationTime)
	assert.Equal(t, expectedSection, article.Section)
	assert.Equal(t, 1, len(article.Authors))
	assert.Equal(t, expectedAuthor, article.Authors[0])
}

func Test_SchemaOrg_ArticleWithEmbeddedAndTopLevelImages(t *testing.T) {
	expectedTitle := "Testcase for ARTICLE with Embedded and Top-Level IMAGEs"
	expectedDescription := "Testing ARTICLE with embedded and top-level images"
	expectedUrl := "http://dummy/Test_SchemaOrg_article_with_embedded_and_toplevel_images.html"
	expectedImage1 := "http://dummy/Test_SchemaOrg_toplevel image.html"
	expectedFormat1 := "gif"
	expectedCaption1 := "A test for top-level IMAGE"
	expectedImage2 := "http://dummy/Test_SchemaOrg_article_with_embedded_and_toplevel_images.html"
	expectedFormat2 := "jpeg"
	expectedCaption2 := "A test for embedded IMAGE in ARTICLE"
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/ImageObject">` +
		`	<span itemprop="headline">Title should be ignored` +
		`	</span>` +
		`	<span itemprop="description">Testing top-level IMAGE` +
		`	</span>` +
		`	<a itemprop="url" href="http://dummy/to_be_ignored_url.html">test results` +
		`	</a>` +
		`	<a itemprop="contentUrl" href="` + expectedImage1 + `">top-level image` +
		`	</a>` +
		`	<span itemprop="encodingFormat">` + expectedFormat1 +
		`	</span>` +
		`	<span itemprop="caption">` + expectedCaption1 +
		`	</span>` +
		`	<meta itemprop="representativeOfPage" content="true">` +
		`	<meta itemprop="width" content="1000">` +
		`	<meta itemprop="height" content="600">` +
		`</div>` +
		`<div id="2" itemscope itemtype="http://schema.org/Article">` +
		`	<span itemprop="headline">` + expectedTitle +
		`	</span>` +
		`	<span itemprop="description">` + expectedDescription +
		`	</span>` +
		`	<a itemprop="url" href="` + expectedUrl + `">test results` +
		`	</a>` +
		`	<img itemprop="image" src="http://dummy/should_be_ignored_image.jpeg">` +
		`	<div id="3" itemscope itemtype="http://schema.org/ImageObject"` +
		`		 itemprop="associatedMedia">` +
		`		<a itemprop="url" href="` + expectedImage2 + `">associated image` +
		`		</a>` +
		`		<span itemprop="encodingFormat">` + expectedFormat2 +
		`		</span>` +
		`		<span itemprop="caption">` + expectedCaption2 +
		`		</span>` +
		`		<meta itemprop="representativeOfPage" content="false">` +
		`		<meta itemprop="width" content="600">` +
		`		<meta itemprop="height" content="400">` +
		`	</div>` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "Article", parser.Type())
	assert.Equal(t, expectedTitle, parser.Title())
	assert.Equal(t, expectedDescription, parser.Description())
	assert.Equal(t, expectedUrl, parser.URL())

	images := parser.Images()
	assert.Equal(t, 2, len(images))

	image := images[0]
	assert.Equal(t, expectedImage2, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, expectedFormat2, image.Type)
	assert.Equal(t, expectedCaption2, image.Caption)
	assert.Equal(t, 600, image.Width)
	assert.Equal(t, 400, image.Height)

	image = images[1]
	assert.Equal(t, expectedImage1, image.URL)
	assert.Equal(t, "", image.SecureURL)
	assert.Equal(t, expectedFormat1, image.Type)
	assert.Equal(t, expectedCaption1, image.Caption)
	assert.Equal(t, 1000, image.Width)
	assert.Equal(t, 600, image.Height)
}

func Test_SchemaOrg_ItemscopeInHTMLTag(t *testing.T) {
	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	setItemScopeAndType(doc, "Article")

	expectedTitle := "Testcase for ItemScope in HTML tag"
	h := testutil.CreateHeading(1, expectedTitle)
	setItemProp(h, "headline")
	dom.AppendChild(body, h)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "Article", parser.Type())
	assert.Equal(t, expectedTitle, parser.Title())
	assert.NotNil(t, parser.Article())
}

func Test_SchemaOrg_SupportedWithUnsupportedItemprop(t *testing.T) {
	expectedTitle := "Testcase for Supported With Unsupported Itemprop"
	expectedSection := "Testing"
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/Article">` +
		`	<span itemprop="headline">` + expectedTitle +
		`	</span>` +
		`	<span itemprop="articleSection">` + expectedSection +
		`	</span>` +
		// Add unsupported AggregateRating to supported Article as itemprop.
		`	<div id="2" itemscope itemtype="http://schema.org/AggregateRating"` +
		`		 itemprop="aggregateRating">Ratings: ` +
		`		<span itemprop="ratingValue">9.9` +
		`		</span>` +
		// Add supported Person to unsupported AggregateRating as itemprop.
		`		<div id="3" itemscope itemtype="http://schema.org/Person"` +
		`			 itemprop="author">Author: ` +
		`			<span itemprop="name">Whoever authored` +
		`			</span>` +
		`		</div>` +
		`	</div>` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "Article", parser.Type())
	assert.Equal(t, expectedTitle, parser.Title())
	assert.Equal(t, "", parser.Description())
	assert.Equal(t, "", parser.URL())
	assert.Equal(t, "", parser.Author())
	assert.Equal(t, "", parser.Publisher())
	assert.Equal(t, "", parser.Copyright())

	article := parser.Article()
	assert.NotNil(t, article)
	assert.Equal(t, "", article.PublishedTime)
	assert.Equal(t, "", article.ModifiedTime)
	assert.Equal(t, "", article.ExpirationTime)
	assert.Equal(t, expectedSection, article.Section)
	assert.Equal(t, 0, len(article.Authors))
}

func Test_SchemaOrg_UnsupportedWithSupportedItemprop(t *testing.T) {
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/Movie">` +
		`	<span itemprop="headline">Testcase for Unsupported With Supported Itemprop` +
		`	</span>` +
		// Add supported Person to unsupported Movie as itemprop.
		`	<div id="3" itemscope itemtype="http://schema.org/Person"` +
		`		 itemprop="publisher">Publisher: ` +
		`		<span itemprop="name">Whoever published` +
		`		</span>` +
		`	</div>` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "", parser.Type())
	assert.Equal(t, "", parser.Title())
	assert.Equal(t, "", parser.Description())
	assert.Equal(t, "", parser.URL())
	assert.Equal(t, "", parser.Author())
	assert.Equal(t, "", parser.Publisher())
	assert.Equal(t, "", parser.Copyright())
	assert.Nil(t, parser.Article())
	assert.Equal(t, 0, len(parser.Images()))
}

func Test_SchemaOrg_UnsupportedWithNestedSupported(t *testing.T) {
	expectedTitle := "Testcase for ARTICLE nested in Unsupported Type"
	expectedDescription := "Testing ARTICLE that is nested within unsupported type"
	expectedUrl := "http://dummy/Test_SchemaOrg_article_with_embedded_items.html"
	expectedImage := "http://dummy/Test_SchemaOrg_article_with_embedded_items.jpeg"
	expectedAuthor := "Whoever authored"
	expectedPublisher := "Whatever Article Incorporated"
	expectedDatePublished := "April 15, 2014"
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/Movie">` +
		`	<span itemprop="headline">Testcase for Unsupported With Supported Itemprop` +
		`	</span>` +
		// Add supported Article to unsupported Movie as a non-itemprop.
		`	<div id="2" itemscope itemtype="http://schema.org/Article">` +
		`		<span itemprop="headline">` + expectedTitle +
		`		</span>` +
		`		<span itemprop="description">` + expectedDescription +
		`		</span>` +
		`		<a itemprop="url" href="` + expectedUrl + `">test results` +
		`		</a>` +
		`		<img itemprop="image" src="` + expectedImage + `">` +
		`		<div id="3" itemscope itemtype="http://schema.org/Person"` +
		`			 itemprop="author">Author: ` +
		`			<span itemprop="name">` + expectedAuthor +
		`			</span>` +
		`		</div>` +
		`		<div id="4" itemscope itemtype="http://schema.org/Organization"` +
		`			 itemprop="publisher">Publisher: ` +
		`			<span itemprop="name">` + expectedPublisher +
		`			</span>` +
		`		</div>` +
		`		<span itemprop="datePublished">` + expectedDatePublished +
		`		</span>` +
		`	</div>` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "Article", parser.Type())
	assert.Equal(t, expectedTitle, parser.Title())
	assert.Equal(t, expectedDescription, parser.Description())
	assert.Equal(t, expectedUrl, parser.URL())
	assert.Equal(t, expectedAuthor, parser.Author())
	assert.Equal(t, expectedPublisher, parser.Publisher())
	assert.Equal(t, "", parser.Copyright())

	images := parser.Images()
	assert.Equal(t, 1, len(images))
	assert.Equal(t, expectedImage, images[0].URL)

	article := parser.Article()
	assert.NotNil(t, article)
	assert.Equal(t, expectedDatePublished, article.PublishedTime)
	assert.Equal(t, "", article.ExpirationTime)
	assert.Equal(t, 1, len(article.Authors))
	assert.Equal(t, expectedAuthor, article.Authors[0])
}

func Test_SchemaOrg_SameItempropDifferentValues(t *testing.T) {
	expectedAuthor := "Author 1"
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/Article">` +
		`	<div id="2" itemscope itemtype="http://schema.org/Person"` +
		`		 itemprop="author">Authors: ` +
		`		<span itemprop="name">` + expectedAuthor +
		`		</span>` +
		`		<span itemprop="name">Author 2` +
		`		</span>` +
		`	</div>` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "Article", parser.Type())
	assert.Equal(t, expectedAuthor, parser.Author())

	article := parser.Article()
	assert.NotNil(t, article)
	assert.Equal(t, 1, len(article.Authors))
	assert.Equal(t, expectedAuthor, article.Authors[0])
}

func Test_SchemaOrg_ItempropWithMultiProperties(t *testing.T) {
	expectedPerson := "Person foo"
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/Article">` +
		`	<div id="2" itemscope itemtype="http://schema.org/Person"` +
		`		 itemprop="author publisher">` +
		`		<span itemprop="name">` + expectedPerson +
		`		</span>` +
		`	</div>` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "Article", parser.Type())
	assert.Equal(t, expectedPerson, parser.Author())
	assert.Equal(t, expectedPerson, parser.Publisher())

	article := parser.Article()
	assert.NotNil(t, article)
	assert.Equal(t, 1, len(article.Authors))
	assert.Equal(t, expectedPerson, article.Authors[0])
}

func Test_SchemaOrg_AuthorPropertyFromDifferentSources(t *testing.T) {
	// Test that "creator" property is used when "author" property doesn't exist.
	expectedCreator := "Whoever created"
	htmlStr := `<div id="1" itemscope itemtype="http://schema.org/Article">` +
		`	<div id="2" itemscope itemtype="http://schema.org/Person"` +
		`		 itemprop="author">Creator: ` +
		`		<span itemprop="name">` + expectedCreator +
		`		</span>` +
		`	</div>` +
		`</div>`

	rootDiv := testutil.CreateDiv(0)
	dom.SetInnerHTML(rootDiv, htmlStr)

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.AppendChild(body, rootDiv)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, "Article", parser.Type())
	assert.Equal(t, expectedCreator, parser.Author())

	article := parser.Article()
	assert.NotNil(t, article)
	assert.Equal(t, 1, len(article.Authors))
	assert.Equal(t, expectedCreator, article.Authors[0])

	// Remove article item from parent, to clear the state for the
	// next rel="author" test.
	body.RemoveChild(rootDiv)

	// Test that rel="author" attribute in an anchor element is used
	// in the absence of "author" or "creator" properties.
	expectedAuthor := "Chromium Authors"
	link := testutil.CreateAnchor("http://dummy/rel_author.html", expectedAuthor)
	dom.SetAttribute(link, "rel", "author")
	dom.AppendChild(body, link)

	parser = schemaorg.NewParser(doc, nil)
	assert.Equal(t, "", parser.Type())
	assert.Equal(t, expectedAuthor, parser.Author())
	assert.Nil(t, parser.Article())
}

func Test_SchemaOrg_GetTitleWhenTheMainArticleDoesntHaveHeadline(t *testing.T) {
	// In original dom-distiller there are title several elements with varying size.
	// Their test is passed if element with `itemprop=name` with biggest area is
	// selected. Unfortunately it's not possible with Go, so here we just test to
	// make sure the parser can get title from element with `itemprop=name`.
	// NEED-COMPUTE-CSS
	expectedTitle := "This is a headline"
	elements := `<div itemscope itemtype="http://schema.org/Article" ` +
		`	style="width: 200; height: 400">` +
		`	<span itemprop="name">` + expectedTitle + `</span>` +
		`</div>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, elements)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, expectedTitle, parser.Title())
}

func Test_SchemaOrg_GetTitleWhenTheMainArticleHasHeadline(t *testing.T) {
	// In original dom-distiller there are several title elements with varying size.
	// Their test is passed if element with `itemprop=headline` with biggest area
	// is selected. Unfortunately it's not possible with Go, so here we just test
	// to make sure the parser can get title from element with `itemprop=headline`.
	// NEED-COMPUTE-CSS
	expectedTitle := "This is a headline"
	elements := `<div itemscope itemtype="http://schema.org/Article" ` +
		`	style="width: 200; height: 400">` +
		`	<span itemprop="name headline">` + expectedTitle + `</span>` +
		`</div>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, elements)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, expectedTitle, parser.Title())
}

func Test_SchemaOrg_GetTitleWithNestedArticles(t *testing.T) {
	// In original dom-distiller within the article there are several title
	// elements with varying size. Their test is passed if element with
	// `itemprop=headline` with biggest area is selected. Unfortunately it's
	// not possible with Go, so here we just test to make sure the parser
	// can get title from nested article. NEED-COMPUTE-CSS
	expectedTitle := "This is a headline"
	elements := `<div itemscope itemtype="http://schema.org/Article" ` +
		`	style="width: 200; height: 100">` +
		`	<div itemscope itemtype="http://schema.org/Article" ` +
		`	style="width: 400; height: 300">` +
		`		<span itemprop="name headline">` + expectedTitle + `</span>` +
		`	</div>` +
		`</div>` +
		`<div itemscope itemtype="http://schema.org/Article" ` +
		`	style="width: 200; height: 200">` +
		`	<span itemprop="name headline">A headline</span>` +
		`</div>`

	doc := testutil.CreateHTML()
	body := dom.QuerySelector(doc, "body")
	dom.SetInnerHTML(body, elements)

	parser := schemaorg.NewParser(doc, nil)
	assert.Equal(t, expectedTitle, parser.Title())
}

func setItemScopeAndType(node *html.Node, schemaType string) {
	dom.SetAttribute(node, "itemscope", "")
	dom.SetAttribute(node, "itemtype", "http://schema.org/"+schemaType)
}

func setItemProp(node *html.Node, name string) {
	dom.SetAttribute(node, "itemprop", name)
}
