# Changelog

### 31 October 2020

- Separate stable version to its own branch.

### 30 October 2020

- From Readability: strip identification and presentational attributes from each nodes.
- Improve lazy load image extractor.
- Mark large blocks around main content as content as well.

### 29 October 2020

- From Readability: exclude nodes which roles indicate that it's not a content.

### 28 October 2020

- From Readability: skip byline, empty divs and unlikely elements.
- From Readability: convert anchor that uses Javascript URL and only contains a single text node into an ordinary text node.
- From Readability: convert `<font>` elements to `<span>`.
- Make sure figure caption doesn't contain `<noscript>` tags.

### 27 October 2020

- Mark `<acronym>` and `<tt>` as inline elements. At this point the port process has finished so I tagged it as v1.0.0.
- From Readability: check if node is probably invisible by using class name and `aria-hidden` attribute.
- From Readability: exclude form and input element.

### 26 October 2020

- Simplify function for getting display style.
- Fix fatal error in doc builder which caused missing contents.

### 25 October 2020

- Add initial test files.
- Improve lazy-loaded image replacer in image extractor.

### 24 October 2020

- Port `LogUtil` froml `LogUtil.java`
- Fix pagination finder `PrevNextFinder` ignores page number in URL queries.

### 22 October 2020

- Merge pagination finder. Now pagination link to previous and next partial page is accessible via `Result.PaginationInfo`.
- Improve page number pagination finder to also find page numbers in web page where its page number are not all consecutive (like in ArsTechnica).
- Move all models from `internal/model` directory (which can't be imported by package's user) to `data` directory.

### 21 October 2020

- Port `PagingLinksFinder` from `PagingLinksFinder.java`
- Restructure models by moving distiller `Result` out of internal directory and remove unused data fields.

### 20 October 2020

- Port `PageParameterParser` from `PageParameterParser.java`
- Fix panic when generating image output.
- Implement test for `testutil.TextDocumentBuilder` following `javatest/TextDocumentConstructionTest`.
- Implement test for `webdoc.TextDocument` following `javatest/TextDocumentStatisticsTest`.

### 19 October 2020

- Port `PathComponentPagePattern` from `PathComponentPagePattern.java`
- Port `PageParameterDetector` from `PageParameterDetector.java`

### 17 October 2020

- Port `PageLinkInfo` from `PageLinkInfo.java`
- Port `PageParamInfo` from `PageParamInfo.java`
- Port `MonotonicPageInfosGroups` from `MonotonicPageInfosGroups.java`
- Port `PagePattern` interface from `PageParameterDetector.java`
- Port `QueryParamPagePattern` from `QueryParamPagePattern.java`

### 16 October 2020

- Port `ContentExtractor` from `ContentExtractor.java`
- Remove `JsTestEntryGenerator` from `javatest/JsTestEntryGenerator.java` because it's only used in Java to prepare the unit tests.

### 14 October 2020

- Port all `ImageScorer` in `webdocuments/filters/images/`
- Port `LeadImageFinder` from `webdocuments/filters/LeadImageFinder.java`
- Port `NestedElementRetainer` from `webdocuments/filters/NestedElementRetainer.java`
- Port `RelevantElements` from `webdocuments/filters/RelevantElements.java`
- Port `TestWebDocumentBuilder` from `javatest/webdocument/TestWebDocumentBuilder.java`

### 13 October 2020

- Port `NumWordsRulesClassifier` from `filters/english/NumWordsRulesClassifier.java`
- Port `TerminatingBlocksFinder` from `filters/english/TerminatingBlocksFinder.java`
- Port `BlockProximityFusion` from `filters/heuristics/BlockProximityFusion.java`
- Port `DocumentTitleMatchClassifier` from `filters/heuristics/DocumentTitleMatchClassifier.java`
- Port `ExpandTitleToContentFilter` from `filters/heuristics/ExpandTitleToContentFilter.java`
- Port `HeadingFusion` from `filters/heuristics/HeadingFusion.java`
- Port `KeepLargestBlockFilter` from `filters/heuristics/KeepLargestBlockFilter.java`
- Port `LargeBlockSameTagLevelToContentFilter` from `filters/heuristics/LargeBlockSameTagLevelToContentFilter.java`
- Port `ListAtEndFilter` from `filters/heuristics/ListAtEndFilter.java`
- Port `SimilarSiblingContentExpansion` from `filters/heuristics/SimilarSiblingContentExpansion.java`
- Port `BoilerplateBlockFilter` from `filters/simple/BoilerplateBlockFilter.java`
- Port `LabelToBoilerplateFilter` from `filters/simple/LabelToBoilerplateFilter.java`
- Port `TestTextBlockBuilder` from `javatest/TestTextBlockBuilder.java`
- Port `TestTextDocumentBuilder` from `javatest/TestTextDocumentBuilder.java`
- Port `TextDocumentTestUtil` from `javatest/document/TextDocumentTestUtil.java`
- Port `TestWebTextBuilder` from `javatest/webdocument/TestWebTextBuilder.java`
- Port `ArticleExtractor` from `extractors/ArticleExtractor.java`
- Remove `filters/simple/MarkEverythingBoilerplateFilter.java` since it's not used anywhere.
- Remove `filters/simple/MarkEverythingContentFilter.java` and `filters/simple/MinWordsFilter.java` since it's only used in `KeepEverythingExtractor.java` and `KeepEverythingWithMinKWordsExtractor.java` that we already removed back in 8 October.

### 12 October 2020

- Port `DomConverter` from `webdocument/DomConverter.java`
- Port `FakeWebDocumentBuilder` from `javatest/webdocument/FakeWebDocumentBuilder.java`
- Replace `alecthomas/assert` with `stretchr/testify/assert`. Nothing wrong with the former but the latter is better since it prints the log as raw text instead of formatted one. Might be useful if in later days we decide to set CI for testing.

### 11 October 2020

- Port `WebDocument` from `webdocument/WebDocument.java`
- Port `WebDocumentBuilder` from `webdocument/WebDocumentBuilder.java`
- Port `EmbedExtractor` from `extractors/embed/EmbedExtractor.java`
- Port `ImageExtractor` from `extractors/embed/ImageExtractor.java`
- Port `TwitterExtractor` from `extractors/embed/TwitterExtractor.java`
- Port `VimeoExtractor` from `extractors/embed/Vimeotractor.java`
- Port `YouTubeExtractor` from `extractors/embed/YouTubeExtractor.java`
- Remove `JavaScript.java` because functions inside it already available in Go standard library.
- Remove `GwtOverlayProtoTest.java` because it's only test model for Protobuf which we don't use.
- Remove `KeepEverythingExtractor.java` and `KeepEverythingWithMinKWordsExtractor.java` because it's not used anywhere.

### 9 October 2020

- Port `WebVideo` from `webdocument/WebVideo.java`
- Port `TextBlock` from `document/TextBlock.java`
- Port `TextDocument` from `document/TextDocument.java` and `document/TextDocumentStatistics.java`
- Add initial MIT license.

### 8 October 2020

- Port `WebTag` from `webdocument/WebTag.java`
- Port `WebText` from `webdocument/WebText.java`
- Port `WebEmbed` from `webdocument/WebEmbed.java`
- Port `WebImage` from `webdocument/WebImage.java`
- Port `WebTable` from `webdocument/WebTable.java`
- Port `WebFigure` from `webdocument/WebFigure.java`
- Port `WebTextBuilder` from `webdocument/WebTextBuilder.java`
- Port `ElementAction` from `webdocument/ElementAction.java`
- Port `DomWalker` from `DomWalker.java`
- Remove `NodeListExpander` since it has identical result as `TreeCloneBuilder` and we already port the latter (even their unit tests are similar).
- Remove `NodeTree` since it's only used in `NodeListExpander`. Besides that, it also requires us to compute stylesheet which is impossible to implement right now.
- Remove `OrderedNodeMatcher` since it's only used in `NodeListExpander` and `TreeCloneBuilder` and our implementation of `TreeCloneBuilder` doesn't require it.

### 7 October 2020

- Port `TableClassifier` from `TableClassifier.java`
- Remove `DomDistillerEntry` since it's useless for our case.
- Remove `Assert` because we already use [`testify`][1] package that provide assertion utilities.
- Remove `JsTestCase`, `JsTestEntry`, `JsTestSuitBase` and `DomDistillerJsTestCase` because it's only used in Java to prepare the unit tests.

### 6 October 2020

- Port `CreateDivTree` from `TestUtil.java`
- Port `BuildTreeClone` from `TreeCloneBuilder.java`

### 5 October 2020

- Port `SchemaOrgParser` and `SchemaOrgParserAccessor` from `SchemaOrg.java`
- Port `MarkupParser` from `MarkupParser.java`
- Port `getDocumentTitle` from `DocumentTitleGetter.java`

### 4 October 2020

- Port `IEReadingViewParser` from `IEReadingViewParser.java`

### 3 October 2020

- Porting process started
- Port `WordCounter` interface from `StringUtil.java`
- Port `OpenGraphParser` and `OpenGraphParserAccessor` from `OpenGraphParser.java`

[1]: https://github.com/stretchr/testify