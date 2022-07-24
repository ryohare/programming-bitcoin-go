package finitefield

import "fmt"

type FieldElement struct {
	Num   int
	Prime int
}

func Make(num, prime int) *FieldElement {
	return &FieldElement{
		Num:   num,
		Prime: prime,
	}
}

func Equal(fe1, fe2 FieldElement) bool {
	return fe1.Num == fe2.Num && fe1.Prime == fe2.Prime
}

func (f *FieldElement) ToString() string {
	return fmt.Sprintf(
		"FieldElement_%d(%d)",
		f.Prime,
		f.Num,
	)
}
