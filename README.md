# Go-DomDistiller [![Go Reference](https://pkg.go.dev/badge/github.com/markusmobius/go-domdistiller.svg)](https://pkg.go.dev/github.com/markusmobius/go-domdistiller)

> This main branch is the development version for Go-DomDistiller which incorporates insights from the readability package as well as other improvements. Check the [stable branch][5] for the stable version that is a faithful port of the original DOM Distiller (the stable branch only receives bug fixes).

Go-DomDistiller is a Go package that finds the main readable content and the metadata from a HTML page. It works by removing clutter like buttons, ads, background images, scripts, etc.

This package is based on [DOM Distiller][0] which is part of the Chromium project that is built using Java language. Unlike DOM distiller there are no dependencies on Chromium or GWT which makes it useful to run as a standalone program on a server.

The structure of this package follows the structure of the original Java code. This way, any improvements from Chromium (hopefully) can be implemented easily here.

The port has been [completed][6] and we have used it to process millions of web pages, so it should be stable enough to use.

## Motivations

We are doing computational social science research on news production and consumption as part of [Project Ratio][8]. We collect a lot of news web pages and extract the article inside it using headless Chrome running Readability.js and DOM Distiller. This works fine, but is unbearably slow.

After looking around, we found out that [Readability.js][1] has been [ported to Go][2] by [@RadhiFadlillah] and it has impressive performance. With that said, we decided to ask him to port DOM Distiller to Go language as well. The port was completely done by Radhi.

## Limitations

The algorithm in the original DOM Distiller incorporates some render-level information in both the classification and tree transduction steps. For example, any elements that are hidden from display are not considered as content, images that are too small are not considered as lead images, etc. These render-level checks are a small part of DOM Distiller's strategy.

Unfortunately it's impossible to do that on the server side without running a full headless browser which we don't want to do (we also only want to rely on the HTML without having to download all the style sheets). Therefore, while porting the original code, we exclude parts where we need to compute the stylesheets. These omissions are marked with [`NEED-COMPUTE-CSS`][3].

Fortunately, according to [research][4] by Mohammad Ghasemisharif et al. (2018) they expect that this modification has minimal effects on extraction results, so we feel confident going forward with the port.

## Comparison with the stable branch

The stable branch is the faithful port of original DOM Distiller which only receives bug fixes, while the main branch adds some [insights][7] from Go-Readability.

Both should be stable enough to use, but if you want to replicate the DOM Distiller results as closely as possible you you may prefer to use the stable branch.

## Comparison with other extractors

As far as we know, currently there are three content extractors built for Go:

- Go-DomDistiller
- [Go-Readability][2]
- [Go-Trafilatura][9]

Since every extractors use its own algorithms, their results are a bit different. In general they give satisfactory results, however we found out that there are some cases where DOM Distiller is better and vice versa. Here is the short summary of pros and cons for each extractor:

Dom Distiller:

- Very fast.
- Good at extracting images from article.
- Able to find next page in sites that separated its article to several partial pages.
- Since the original library was embedded in Chromium browser, its tests are pretty thorough.
- CON: has a huge codebase, mostly because it mimics the original Java code.
- CON: the original library is not maintained anymore and has been archived.

Readability:

- Fast, although not as fast as Dom Distiller.
- Better than DOM Distiller at extracting wiki and documentation pages.
- The original library in Readability.js is still actively used and maintained by Firefox.
- The codebase is pretty small.
- CON: the unit tests are not as thorough as the other extractors.

Trafilatura:

- Has the best accuracy compared to other extractors.
- Better at extracting web page's metadata, including its language and publish date.
- Its unit tests are thorough and focused on removing noise while making sure the real contents are still captured.
- Designed to be used in academic domain e.g. natural language processing.
- Actively maintained with new release almost every month.
- CON: slower than the other extractors, mostly because it also looks for language and publish date.
- CON: doesn't really good at extracting images.

The benchmark that compares these extractors is available in [this repository][benchmark]. It uses each extractor to process 983 web pages in single thread. Here is its benchmark result:

|             Extractor             | Time (ms) | Memory (MB) | Mem Allocation (allocs) |
| :-------------------------------: | :-------: | :---------: | :---------------------: |
|            Readability            |   4,212   |    4,412    |       15,261,650        |
|           DomDistiller            |   3,794   |    4,144    |       13,552,246        |
|  DomDistiller+PaginationPrevNext  |   5,263   |    4,598    |       22,744,038        |
| DomDistiller+PaginationPageNumber |   4,156   |    4,222    |       15,669,698        |
|            Trafilatura            |   6,609   |    3,585    |       33,628,972        |
|       Trafilatura+Fallback        |  12,934   |    8,781    |       55,338,023        |

And here is its performance comparison result:

|            Package             | Precision | Recall | Accuracy | F-Score |
| :----------------------------: | :-------: | :----: | :------: | :-----: |
|        `go-readability`        |   0.870   | 0.881  |  0.875   |  0.875  |
|       `go-domdistiller`        |   0.871   | 0.864  |  0.868   |  0.867  |
|        `go-trafilatura`        |   0.909   | 0.886  |  0.899   |  0.897  |
| `go-trafilatura` with fallback |   0.911   | 0.902  |  0.907   |  0.906  |

## Installation

To install the development version of this package, just run `go get` for main branch :

```
go get -u -v github.com/markusmobius/go-domdistiller@main
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

- `LogNothing` will make distiller completely disable the log.
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

- `PrevNext` is the algorithm to find pagination links that works by scoring each anchor in documents using various heuristics on its href, text, class name and ID. It's quite accurate and used as default algorithm. Unfortunately it uses a lot of regular expressions, so it's a bit slow.
- `PageNumber` is algorithm to find pagination links that works by collecting groups of adjacent plain text numbers and outlinks with digital anchor text. It's a lot faster than PrevNext, but also less accurate.

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

	// Node is the *html.Node which contain the distilled content.
	Node *html.Node

	// Text is the string which contains the distilled content in text format.
	Text string

	// ContentImages is list of image URLs that used within the distilled content.
	ContentImages []string
}
```

The `MarkupInfo`, `TimingInfo` and `PaginationInfo` field are defined in `data` package `github.com/markusmobius/go-domdistiller/data` like this :

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

	"github.com/go-shiori/dom"
	distiller "github.com/markusmobius/go-domdistiller"
)

func main() {
	url := "https://arstechnica.com/gadgets/2020/10/iphone-12-and-12-pro-double-review-playing-apples-greatest-hits/"

	// Start distiller
	result, err := distiller.ApplyForURL(url, time.Minute, nil)
	if err != nil {
		panic(err)
	}

	rawHTML := dom.OuterHTML(result.Node)
	fmt.Println(rawHTML)
}
```

### Extracting content from a HTML file

```go
package main

import (
	"fmt"

	"github.com/go-shiori/dom"
	distiller "github.com/markusmobius/go-domdistiller"
)

func main() {
	result, err := distiller.ApplyForFile("example/sample.html", nil)
	if err != nil {
		panic(err)
	}

	rawHTML := dom.OuterHTML(result.Node)
	fmt.Println(rawHTML)
}
```

## Licenses

Go-DomDistiller is distributed under [MIT license](https://choosealicense.com/licenses/mit/) which means you can use and modify it however you want. However, if you make an enhancement for it, if possible please send a pull request.

We are indebted to the Chromium authors for the amazing DOM Distiller. We are equally indebted to Christian Kohlschütter who wrote a content parser called Boilerpipe in 2009 which is based on his PhD thesis and which is also the basis for DOM Distiller (the original Boilerpipe still produces amazing results for most pages). Boilerpipe is licensed under the Apache 2.0 license and DOM Distiller has a BSD-style license. Since our work is derived directly from DOM Distiller and indirectly from Boilerpipe we have included the respective copyright notices at the top of each file as well as the license files for both prior projects.

[0]: https://chromium.googlesource.com/chromium/dom-distiller
[1]: https://github.com/mozilla/readability
[2]: https://github.com/go-shiori/go-readability
[3]: https://github.com/markusmobius/go-domdistiller/search?q=NEED-COMPUTE-CSS
[4]: https://arxiv.org/abs/1811.03661
[5]: https://github.com/markusmobius/go-domdistiller/tree/stable
[6]: https://github.com/markusmobius/go-domdistiller/blob/main/CHANGELOG.md
[7]: https://github.com/markusmobius/go-domdistiller/blob/main/IMPROVEMENTS.md
[8]: https://www.microsoft.com/en-us/research/project/project-ratio/
[9]: https://github.com/markusmobius/go-trafilatura
[@RadhiFadlillah]: https://github.com/RadhiFadlillah
[benchmark]: https://github.com/markusmobius/content-extractor-benchmark
