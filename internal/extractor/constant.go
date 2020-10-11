// ORIGINAL: java/extractors/embeds/*.java

package extractor

var (
	lazyImageAttrs = []string{
		"data-src",
		"data-original",
		"datasrc",
		"data-url",
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
