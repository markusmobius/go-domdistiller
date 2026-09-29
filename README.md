# go-domdistiller

`go-domdistiller` extracts article text, HTML, images, metadata and pagination
links from web pages. It is a native Go port of
[chromium/dom-distiller](https://chromium.googlesource.com/chromium/dom-distiller),
with additional extraction heuristics on the main branch.

## Philosophy

Our extractor packages share three principles:

1. **Bring your own HTML.** Keep page acquisition separate from extraction.
   The primary workflow uses HTML supplied by the caller, who controls fetching,
   caching, rendering, retries and scheduling.
2. **Stay close to upstream.** Preserve the algorithms and behavior of each
   package's declared upstream reference as closely as possible. Document
   deliberate differences and compatibility limits in [UPSTREAM.md](UPSTREAM.md)
   rather than claiming exact equivalence on every page.
3. **Provide very fast Go and Rust packages.** Run extraction natively, without
   a Python or Java runtime. Improve throughput and allocation efficiency while
   preserving intended behavior, and substantiate performance with reproducible
   benchmarks that report quality alongside speed.

## Overview

The current `go-domdistiller` release is **v1.0.0**. It accepts an HTML tree,
an `io.Reader`, a file or a URL. Reader and tree extraction do not fetch pages;
the explicit URL helper downloads the page before extracting it.

Results include the article node and plain text, title, available metadata,
image URLs, word count and previous/next page links. Metadata depends on what
the page provides; extraction does not run JavaScript or compute browser layout.

The corresponding Rust package,
[rust-domdistiller](https://github.com/markusmobius/rust-domdistiller), uses
`go-domdistiller` as its behavioral reference.

## Installation

```sh
go get github.com/markusmobius/go-domdistiller@v1.0.0
```

Import the package as `distiller`. See [go.mod](go.mod) for the module's Go and
dependency requirements and [CHANGELOG.md](CHANGELOG.md) for release changes.

## Usage

Extract text from HTML already held in memory:

```go
package main

import (
	"fmt"
	"net/url"
	"strings"

	distiller "github.com/markusmobius/go-domdistiller"
)

func main() {
	pageURL, err := url.Parse("https://example.org/research")
	if err != nil {
		panic(err)
	}

	source := `<html><head><title>Research results</title></head><body><article>
<h1>Research results</h1>
<p>The research team compared several methods for extracting articles from saved
web pages. Every method received the same original HTML, and the evaluation
kept the reference text separate from the input supplied to each extractor.</p>
<p>The report records the complete experiment, including errors and repeated
measurements. Its results describe this collection of pages and do not promise
the same quality or execution time for every website.</p>
</article></body></html>`

	result, err := distiller.ApplyForReader(strings.NewReader(source), &distiller.Options{
		OriginalURL:    pageURL,
		SkipPagination: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Text)
}
```

| Entry Point | Input |
| --- | --- |
| `Apply` | An existing `*html.Node` |
| `ApplyForReader` | HTML from an `io.Reader` |
| `ApplyForFile` | A local HTML file path |
| `ApplyForURL` | A URL to download, with a caller-supplied timeout |

Each returns `(*Result, error)`. Use `result.Text` for plain text and
`result.Node` for the extracted HTML tree. `MarkupInfo`, `ContentImages` and
`PaginationInfo` expose metadata, article images and pagination links.
Full types and signatures are in the
[go-domdistiller API reference](https://pkg.go.dev/github.com/markusmobius/go-domdistiller).

## Options

Pass `nil` for default options, or supply `distiller.Options`:

| Option | Default | Effect |
| --- | --- | --- |
| `OriginalURL` | Unset | Supplies page context for relative links and pagination. `ApplyForURL` sets it from its URL argument. |
| `SkipPagination` | `false` | Set to `true` to omit pagination detection when only the current article is needed. |
| `PaginationAlgo` | `PrevNext` | Select `PrevNext` for scored previous/next links or `PageNumber` for groups of numbered page links. |
| `LogFlags` | `LogNothing` | Enable extraction, visibility, pagination or timing logs; combine flags with bitwise OR or use `LogEverything`. |

Pagination detection identifies links; it does not assemble a multi-page
article. The comparison below disables pagination for `go-domdistiller` and
`rust-domdistiller`.

## Current Quality and Speed

The [2026-09-29 shared benchmark](https://github.com/markusmobius/content-extractor-benchmark/blob/ec719092d12f4d2a438dd29d9f4405aab6e0a321/README.md#results-2026-09-29)
compares the six packages below on **2,659 saved pages**: 983 LegoNews,
181 ScrapingHub and 1,495 WCXB.

### Extraction Speed

| Go Package (Measured Version) | Rust Package (Measured Version) | Go ms/page | Rust ms/page | Go/Rust |
| --- | --- | ---: | ---: | ---: |
| `go-readabilityV2` 0.6.0 | `rust-readability-v2` 0.6.5 | 4.755 | 3.945 | 1.21x |
| `go-domdistiller` 1.0.0 | `rust-domdistiller` 1.0.1 | 6.159 | 3.400 | 1.81x |
| `go-trafilatura` 2.2.6 (FAST) | `rust-trafilatura` 2.2.6 (FAST) | 11.329 | 6.570 | 1.72x |

Times are means of **all four measured passes after one warmup**. Go/Rust is
the named Go package's time divided by the named Rust package's time, not an
old/new release speedup. Later documentation-only releases do not change the
versions actually measured.

The run used Windows 11, Ryzen AI 7 PRO 350, Go 1.27.1 and Rust 1.98.1 GNU
with ThinLTO/mimalloc. Extraction includes required working copies, metadata
and text rendering. File I/O, startup, IPC, response serialization and scoring
are excluded. Comments and pagination are off; tables are on.
`go-trafilatura` and `rust-trafilatura` use FAST with external fallback disabled.
Power and sleep checks passed.

Parsing is separate: **Go 11.283 / Rust 6.386 ms/page**, charged once per
language/page for the shared suite. It includes decoding, DOM construction and
the separate `go-trafilatura` / `rust-trafilatura` noscript tree when needed.
These are extraction-stage comparisons, not complete request latencies.

### Text Quality

Each named pair has equal text scores. Errors are listed in LegoNews /
ScrapingHub / WCXB order and remain in the scoring denominators.

| Go Package | Rust Package | LegoNews F1 | ScrapingHub F1 | WCXB F1 | Errors |
| --- | --- | ---: | ---: | ---: | --- |
| `go-readabilityV2` | `rust-readability-v2` | 87.82711% | 95.20557% | 78.47603% | 7 / 0 / 28 |
| `go-domdistiller` | `rust-domdistiller` | 86.74080% | 92.74280% | 74.39696% | 0 / 0 / 0 |
| `go-trafilatura` (FAST) | `rust-trafilatura` (FAST) | 90.91534% | 96.15663% | 78.51703% | 4 / 0 / 10 |

The corpora use different scoring rules; their F1 scores must not be averaged.
Equal text scores do not imply identical metadata: `go-trafilatura` and
`rust-trafilatura` differ on one title and one author field. The
[full report](https://github.com/markusmobius/content-extractor-benchmark/blob/49c426d6135df81b7d492bea7e6aec8e6d77d80c/go_rust_shared_performance_2026_09_29.json)
contains metadata scores, differences, every pass and source/build identities.

## Compatibility and Limitations

- **Main and stable branches differ.** The `go-domdistiller` v1.0.0 release
  identifies the main-branch implementation, including its
  [documented improvements](IMPROVEMENTS.md). The separate
  [go-domdistiller stable branch](https://github.com/markusmobius/go-domdistiller/tree/stable)
  stays closer to the original `chromium/dom-distiller` algorithm.
- **No browser rendering.** CSS/layout-dependent parts of `chromium/dom-distiller`
  are omitted. Saved HTML may lack content generated by JavaScript.
  `go-domdistiller` does not promise identical output to `chromium/dom-distiller`.
- **Heuristic extraction.** Boilerplate can remain and article content can be
  missed. Corpus scores do not guarantee results for an arbitrary page.
- **Not a sanitizer.** Treat extracted HTML as untrusted and sanitize it before
  displaying it in an application.

See [UPSTREAM.md](UPSTREAM.md) for source ancestry, compatibility boundaries,
verification evidence and historical comparisons.

## Development

From a checkout with Go installed:

```sh
go test ./...
go vet ./...
```

Keep extraction, pagination and metadata changes covered by the existing tests.
Documentation and release maintenance rules are in [AGENTS.md](AGENTS.md).

## License and Credits

`go-domdistiller` is [MIT-licensed](LICENSE). It incorporates work from
`chromium/dom-distiller` under its
[BSD and Apache terms](https://chromium.googlesource.com/chromium/dom-distiller/+/2a180397710719913340a12804affc65b789275e/LICENSE)
and `kohlschutter/boilerpipe` under [Apache-2.0](LICENSE-boilerpipe.txt), with
the inherited [notice](NOTICE-boilerpipe.txt).

Radhi Fadlillah implemented the original `go-domdistiller` port for Project Ratio.
The Chromium Authors created `chromium/dom-distiller`, building on
Christian Kohlschuetter's `kohlschutter/boilerpipe`. Markus Mobius maintains this
package and holds its MIT copyright. These upstream creators and contributors made
the port possible; their licenses and notices are retained.