// ORIGINAL: javatest/PageParamContentInfo.java

package testutil

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

func PPCIUnrelatedTerms() *PageParamContentInfo {
	return &PageParamContentInfo{Type: UnrelatedTerms}
}

func PPCINumberInPlainText(number int) *PageParamContentInfo {
	return &PageParamContentInfo{
		Type:   NumberInPlainText,
		Number: number,
	}
}

func PPCINumericOutlink(targetURL string, number int) *PageParamContentInfo {
	return &PageParamContentInfo{
		Type:      NumericOutlink,
		TargetURL: targetURL,
		Number:    number,
	}
}
