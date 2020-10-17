// ORIGINAL: java/PageParamInfo.java

package pagination

import (
	"fmt"
)

// LinearFormula stores the coefficient and delta values of the linear formula:
// pageParamValue = coefficient * pageNum + delta.
type LinearFormula struct {
	coefficient int
	delta       int
}

func NewLinearFormula(coefficient, delta int) *LinearFormula {
	return &LinearFormula{
		coefficient: coefficient,
		delta:       delta,
	}
}

func (lf *LinearFormula) String() string {
	return fmt.Sprintf("coefficient=%d, delta=%d", lf.coefficient, lf.delta)
}
