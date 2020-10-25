// ORIGINAL: java/extractors/embeds/*.java

package embed

import "regexp"

var (
	rxB64DataURL      = regexp.MustCompile(`(?i)^data:\s*([^\s;,]+)\s*;\s*base64\s*`)
	rxSrcsetURL       = regexp.MustCompile(`(?i)(\S+)(\s+[\d.]+[xw])?(\s*(?:,|$))`)
	rxImgExtensions   = regexp.MustCompile(`(?i)\.(jpg|jpeg|png|webp)`)
	rxLazyImageSrcset = regexp.MustCompile(`(?i)\.(jpg|jpeg|png|webp)\s+\d`)
	rxLazyImageSrc    = regexp.MustCompile(`(?i)^\s*\S+\.(jpg|jpeg|png|webp)\S*\s*$`)

	figureImageSelectors = []string{
		"noscript picture",
		"noscript img",
		"picture",
		"img",
	}

	lazyImageSrcAttrs = []string{
		"data-src",
		"data-original",
		"datasrc",
		"data-url",
	}

	lazyImageSrcsetAttrs = []string{
		"data-srcset",
		"datasrcset",
	}

	relevantImageTags = map[string]struct{}{
		// TODO: Add "div" to this list for css images and possibly captions.
		"img":     {},
		"picture": {},
		"figure":  {},
		"span":    {},
	}

	relevantTwitterTags = map[string]struct{}{
		"blockquote": {},
		"iframe":     {},
	}

	relevantVimeoTags = map[string]struct{}{
		"iframe": {},
	}

	relevantYouTubeTags = map[string]struct{}{
		"iframe": {},
		"object": {},
	}
)
