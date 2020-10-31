# Go-DomDistiller

Go-DomDistiller is a Go package that finds the main readable content and the metadata from a HTML page. It works by removing clutter like buttons, ads, background images, scripts, etc.

This package is based on [DOM Distiller][0] which is part of the Chromium project that is built using Java language. The structure of this package is arranged following the structure of original Java code. This way, any improvements from Chromium (hopefully) can be implemented easily here.

## Motivations

We are doing computational social science research on news consumption, so we collect a lot of web pages and extract the article inside it using headless Chrome running Readability.js and DOM Distiller. This works fine, but unbearably slow.

After looking around, we found out that [Readability.js][1] has been [ported to Go][2] and it has an impressive performance. With that said, we decided to port DOM Distiller to Go language as well.

## Limitations

In original DOM Distiller, they consider some render-level information in both the classification and tree transduction steps. For example, any elements that are hidden from display are not considered as content, images that are too small are not considered as lead images, etc. These render-level checks are a small part of DOM Distiller's strategy.

Unfortunately it's impossible to do that on the server side, and we don't want to use a headless browser anymore. So, while porting the original code, we exclude parts where we need to compute the stylesheets. These omits are marked with [`NEED-COMPUTE-CSS`][3].

Fortunately, according to [research][4] by Mohammad Ghasemisharif et al. (2018) they expect that this modification has minimal effects on extraction results, so we feel confident going forward with the port.

## Comparison with Go-Readability

Since Readability and DOM Distiller work using different algorithms, their results are a bit different. In general they give satisfactory results, however we found out that there are some cases where DOM Distiller is better and vice versa. In practice we use both of them then use some kind of scoring to find out which extraction result is more suitable for our use case.

The pros of Dom Distiller :

- better at extracting images;
- better at extracting article's metadata;
- able to find next page in sites that separated its article to several partial pages;
- suitable for processing news articles.

The pros of Readability :

- faster extraction speed;
- better than DOM Distiller at extracting wiki and documentation pages.

Here is the benchmark result between DOM Distiller and Readability :

```
BenchmarkReadability-8                     	       1	22270423614 ns/op	5134614848 B/op	21071083 allocs/op
BenchmarkDistillerWithoutPagination-8      	       1	24248745284 ns/op	7987711256 B/op	30309028 allocs/op
BenchmarkDistillerPageNumberPagination-8   	       1	33292305569 ns/op	8080458848 B/op	32918938 allocs/op
BenchmarkDistillerPrevNextPagination-8     	       1	47737605918 ns/op	8378848776 B/op	36243299 allocs/op
```

## Installation

To install this package, just run `go get` :

```
go get -u -v github.com/markusmobius/go-domdistiller
```

## API Documentation

Dom Distiller has four functions :

- `Apply(doc *html.Node, opts *Options) (*Result, error)`

	This function will apply distiller to the specified HTML node.

- `ApplyForReader(r io.Reader, opts *Options) (*Result, error)`

	This function parses input that received from the specified reader into a HTML node then pass it into the `Apply` function.

- `ApplyForFile(path string, opts *Options) (*Result, error)`

	This function open the file at specified path then pass it into the `ApplyForReader` function.

- `ApplyForURL(url string, timeout time.Duration, opts *Options) (*Result, error)`

	This function download the web page at specified URL then pass it into the `ApplyForReader` function.

Each function accept custom `Option` which is a struct that defined like this :

```go
type Options struct {
	// Whether to extract only the text (or to include the containing html).
	ExtractTextOnly bool

	// Flags to specify which info to dump to log.
	LogFlags LogFlag

	// Original URL of the page, which is used in the heuristics in detecting
	// next/prev page links. Will be ignored if Option is used in ApplyForURL.
	OriginalURL *url.URL

	// Set to true to skip process for finding pagination.
	SkipPagination bool

	// Algorithm to use for next page detection.
	PaginationAlgo PaginationAlgo
}
```

There are several flags available for `LogFlags` :

- `LogExtraction` will make distiller print info of each process when extracting article.
- `LogVisibility` will make distiller print info on why an element is visible.
- `LogPagination` will make distiller print info of pagination process.
- `LogTiming` will make distiller print info of duration of each process when extracting article.

Since `LogFlag` is bit, you can use several flags using bitwise operator `OR` like this :

```go
opts := &distiller.Options{
	LogFlags: distiller.LogExtraction | distiller.LogVisibility,
}
```

Or if you want to log everything, you can use `LogEverything` flag :

```go
opts := &distiller.Options{	LogFlags: distiller.LogEverything }
```

There are two values available for `PaginationAlgo` :

- `PrevNext` is the algorithm to find pagination links that work by scoring  each anchor in documents using various heuristics on its href, text, class name and ID. It's quite accurate and used as default algorithm. Unfortunately it uses a lot of regular expressions, so it's a bit slow. 
- `PageNumber` is algorithm to find pagination links that work by collecting groups of adjacent plain text numbers and outlinks with digital anchor text. A lot faster than PrevNext, but also less accurate.

The distillation result is defined as struct like this :

```go
type Result struct {
	// URL is the URL of the processed page.
	URL string

	// Title is the title of the processed page.
	Title string

	// MarkupInfo is the metadata of the page. The metadata is extracted following three markup
	// specifications: OpenGraphProtocol, IEReadingView and SchemaOrg. For now, OpenGraph protocol
	// takes precedence because it uses specific meta tags and hence the fastest. The other
	// specifications is used as fallback in case some metadata not found.
	MarkupInfo data.MarkupInfo

	// TimingInfo is the record of the time it takes to do each step in the process of content extraction.
	TimingInfo data.TimingInfo

	// PaginationInfo contains link to previous and next partial page. This is useful for long article or
	// that may be partitioned into several partial pages by its webmaster.
	PaginationInfo data.PaginationInfo

	// WordCount is the count of words within document.
	WordCount int

	// HTML is the string which contains the distilled content in HTML format.
	HTML string

	// ContentImages is list of image URLs that used within the distilled content.
	ContentImages []string
}
```

The `MarkupInfo`, `TimingInfo` and `PaginationInfo` field are defined in package `github.com/markusmobius/go-domdistiller/data` like this :

```go
type PaginationInfo struct {
	NextPage string
	PrevPage string
}

type MarkupArticle struct {
	PublishedTime  string
	ModifiedTime   string
	ExpirationTime string
	Section        string
	Authors        []string
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

type MarkupImage struct {
	Root      string
	URL       string
	SecureURL string
	Type      string
	Caption   string
	Width     int
	Height    int
}
```

## Examples

### Extracting web page from an URL

```go
package main

import (
	"fmt"
	"time"

	distiller "github.com/markusmobius/go-domdistiller"
)

func main() {
	url := "https://arstechnica.com/gadgets/2020/10/iphone-12-and-12-pro-double-review-playing-apples-greatest-hits/"

	// Start distiller
	result, err := distiller.ApplyForURL(url, time.Minute, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.HTML)
}
```

### Extracting content from a HTML file

```go
package main

import (
	"fmt"

	distiller "github.com/markusmobius/go-domdistiller"
)

func main() {
	result, err := distiller.ApplyForFile("example/sample.html", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.HTML)
}
```

## Licenses

Go-DomDistiller is distributed under [MIT license](https://choosealicense.com/licenses/mit/) which means you can use and modify it however you want. However, if you make an enhancement for it, if possible please send a pull request.

[0]: https://chromium.googlesource.com/chromium/dom-distiller
[1]: https://github.com/mozilla/readability
[2]: https://github.com/go-shiori/go-readability
[3]: https://github.com/markusmobius/go-domdistiller/search?q=NEED-COMPUTE-CSS
[4]: https://arxiv.org/abs/1811.03661