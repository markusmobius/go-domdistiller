# Go-DomDistiller

Go-DomDistiller is a Go package that finds the main readable content and the metadata from a HTML page. It works by removing clutter like buttons, ads, background images, script, etc.

This package is based on [DOM Distiller][0] which is part of the Chromium project that is built using Java language. The structure of this package is arranged following the structure of original Java code. This way, any improvements from Chromium can be implemented easily here. Another advantage, hopefully all web page that can be parsed by the original Dom Distiller can be parsed by this package as well with identical result.

## Status

This package is still in development and the port process is still not finished. There are 29 files with 4,627 lines of code that haven’t been ported, so there is still long way to go.

## Changelog

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

[0]: https://chromium.googlesource.com/chromium/dom-distiller
[1]: https://github.com/stretchr/testify