// ORIGINAL: javatest/PageParamContentInfo.java

package info_test

type PageParamContentType uint

const (
	UnrelatedTerms PageParamContentType = iota
	NumberInPlainText
	NumericOutlink
)

type PageParamContentInfo struct {
	Type      PageParamContentType
	TargetURL string
	Number    int
}

func ppciUnrelatedTerms() *PageParamContentInfo {
	return &PageParamContentInfo{Type: UnrelatedTerms}
}

func ppciNumberInPlainText(number int) *PageParamContentInfo {
	return &PageParamContentInfo{
		Type:   NumberInPlainText,
		Number: number,
	}
}

func ppciNumericOutlink(targetURL string, number int) *PageParamContentInfo {
	return &PageParamContentInfo{
		Type:      NumericOutlink,
		TargetURL: targetURL,
		Number:    number,
	}
}
