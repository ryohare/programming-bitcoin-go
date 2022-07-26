package fieldelement

import (
	"fmt"
	"math/big"
)

type FieldElement struct {
	Num   *big.Int
	Prime *big.Int
}

func (f FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%s(%s)", f.Prime.String(), f.Num.String())
}

func MakeFieldElement(num, prime big.Int) *FieldElement {
	return &FieldElement{
		Num:   &num,
		Prime: &prime,
	}
}

func Equal(fe1, fe2 *FieldElement) bool {
	if fe1.Num.Cmp(fe2.Num) == 0 && fe1.Prime.Cmp(fe2.Prime) == 0 {
		return true
	}
	return false
}

func Add(fe1, fe2 *FieldElement) (*FieldElement, error) {
	if fe1.Prime.Cmp(fe2.Prime) != 0 {
		return nil, fmt.Errorf("cannot add two numbers in different fields")
	}

	num := big.NewInt(0)
	num.Add(fe1.Num, fe2.Num)
	num.Mod(num, fe1.Prime)

	return &FieldElement{
			Num:   num,
			Prime: fe1.Prime,
		},
		nil
}

func Subtract(fe1, fe2 *FieldElement) (*FieldElement, error) {
	if fe1.Prime.Cmp(fe2.Prime) != 0 {
		return nil, fmt.Errorf("cannot sub two numbers in different fields")
	}

	// big int will store the result in the lhs of the object, so if we need to preserve
	// the object, we want to create a new empty one for holding the result of the sub
	sub := big.NewInt(0)

	// sub = fe1.Num - fe2.Num
	sub.Sub(fe1.Num, fe2.Num)

	// sub % fe1.Prime
	num := sub.Mod(sub, fe1.Prime)

	return &FieldElement{
			Num:   num,
			Prime: fe1.Prime,
		},
		nil
}

func Multiply(fe1, fe2 *FieldElement) (*FieldElement, error) {
	if fe1.Prime.Cmp(fe2.Prime) != 0 {
		return nil, fmt.Errorf("cannot multiply two numbers in different fields")
	}

	mul := big.NewInt(0)
	mod := big.NewInt(0)
	mul.Mul(fe1.Num, fe2.Num)
	mod.Mod(mul, fe1.Prime)
	num := mod

	return &FieldElement{
			Num:   num,
			Prime: fe1.Prime,
		},
		nil
}

func Exponentiate(fe1 *FieldElement, exponent big.Int) (*FieldElement, error) {
	fmt.Println(fe1.Num.String())
	// n := big.NewInt(0)
	// primeMinus1 := big.NewInt(0)

	// primeMinus1 = primeMinus1.Sub(fe1.Prime, big.NewInt(1))
	// n = n.Mod(exponent, primeMinus1)

	// //n := exponent % (fe1.Prime - 1)

	// // num := int64(math.Pow(float64(fe1.Num), float64(n))) % fe1.Prime

	// // num = fe1.num^n % fe1.Prime
	// fe1ToTheN := big.NewInt(fe1.Num.Int64())
	// fe1ToTheN = fe1ToTheN.Exp(fe1ToTheN, n, fe1.Prime)

	// n = exponent % (self.prime - 1)
	// num = pow(self.num, n, self.prime)
	tmp := new(big.Int)
	prime := tmp.Sub(fe1.Prime, big.NewInt(1))
	exp := exponent
	n := exp.Mod(&exponent, prime)
	tmp = new(big.Int)
	num := tmp.Exp(fe1.Num, n, fe1.Prime)

	return &FieldElement{
			Num:   num,
			Prime: fe1.Prime,
		},
		nil
}

func Divide(fe1, fe2 *FieldElement) (*FieldElement, error) {
	if fe1.Prime.Cmp(fe2.Prime) != 0 {
		return nil, fmt.Errorf("cannot multiply two numbers in different fields")
	}
	// num = (self.num * pow(other.num, self.prime - 2, self.prime)) % self.prime

	// solve the pow section first
	// pow(other.num, self.prime - 2, self.prime)
	pow := big.NewInt(0)
	primeMinus2 := big.NewInt(0)
	primeMinus2 = primeMinus2.Sub(fe1.Prime, big.NewInt(2))
	pow = pow.Exp(fe2.Num, primeMinus2, fe1.Prime)

	// distribute out the constant self.num
	pow = pow.Mul(pow, fe1.Num)

	// final modulo
	pow = pow.Mod(pow, fe1.Prime)

	return &FieldElement{
			Num:   pow,
			Prime: fe1.Prime,
		},
		nil
}

// def __rmul__(self, coefficient):
// num = (self.num * coefficient) % self.prime
// return self.__class__(num=num, prime=self.prime)
func RMultiply(fe1 *FieldElement, coefficient *big.Int) (*FieldElement, error) {
	num := big.NewInt(0)
	num = num.Mul(fe1.Num, coefficient)
	num = num.Mod(num, fe1.Prime)
	return &FieldElement{Num: num, Prime: fe1.Prime}, nil
}
