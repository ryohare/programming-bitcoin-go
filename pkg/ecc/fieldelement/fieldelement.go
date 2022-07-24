package fieldelemen

import (
	"fmt"
	"math"
)

type FieldElement struct {
	Num   int64
	Prime int64
}

func MakeFieldElement(num, prime int64) *FieldElement {
	return &FieldElement{
		Num:   num,
		Prime: prime,
	}
}

func Equal(fe1, fe2 FieldElement) bool {
	return fe1.Num == fe2.Num && fe1.Prime == fe2.Prime
}

func Add(fe1, fe2 *FieldElement) (*FieldElement, error) {
	if fe1.Prime != fe2.Prime {
		return nil, fmt.Errorf("cannot add two numbers in different fields")
	}

	num := (fe1.Num + fe2.Num) % fe1.Prime

	return &FieldElement{
			Num:   num,
			Prime: fe1.Prime,
		},
		nil
}

func Subtract(fe1, fe2 *FieldElement) (*FieldElement, error) {
	if fe1.Prime != fe2.Prime {
		return nil, fmt.Errorf("cannot sub two numbers in different fields")
	}

	num := (fe1.Num - fe2.Num) % fe1.Prime

	return &FieldElement{
			Num:   num,
			Prime: fe1.Prime,
		},
		nil

}

func Multiply(fe1, fe2 *FieldElement) (*FieldElement, error) {
	if fe1.Prime != fe2.Prime {
		return nil, fmt.Errorf("cannot multiply two numbers in different fields")
	}

	num := (fe1.Num * fe2.Num) % fe1.Prime

	return &FieldElement{
			Num:   num,
			Prime: fe1.Prime,
		},
		nil
}

func Exponentiate(fe1 *FieldElement, exponent int64) (*FieldElement, error) {

	n := exponent % (fe1.Prime - 1)
	fmt.Println(n)
	num := int64(math.Pow(float64(fe1.Num), float64(n))) % fe1.Prime
	fmt.Println(num)
	return &FieldElement{
			Num:   num,
			Prime: fe1.Prime,
		},
		nil
}

func Divide(fe1, fe2 *FieldElement) (*FieldElement, error) {
	if fe1.Prime != fe2.Prime {
		return nil, fmt.Errorf("cannot exponentiate two numbers in different fields")
	}

	//num = self.num * pow(other.num, self.prime - 2, self.prime) % self.prime
	num := fe1.Num * int64(math.Pow(float64(fe2.Num), float64(fe1.Prime-2))) & fe1.Prime

	return &FieldElement{
			Num:   num,
			Prime: fe1.Prime,
		},
		nil
}

func (f *FieldElement) ToString() string {
	return fmt.Sprintf(
		"FieldElement_%d(%d)",
		f.Prime,
		f.Num,
	)
}
