package opcodes

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"

	S256 "github.com/ryohare/programming-bitcoin-go/pkg/ecc/curves/secp256k1"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
	"golang.org/x/crypto/ripemd160"
)

type StackElement struct {
	Bytes  []byte
	OpCode bool
}

type Stack struct {
	Elements []StackElement
}

// Abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// copy an array of stack elements into a newly allocated array
func scopy(s []StackElement) []StackElement {
	r := make([]StackElement, len(s))
	copy(r, s)
	return r
}

// encode and int into a byte array
func encode(num int) []byte {

	// If the num is 0, return an empty byte array
	if num == 0 {
		return []byte{}
	}

	// absolute value of the number
	absNum := abs(num)

	// flag indicating if the number is negative
	negative := (num < 0)

	// results array
	result := []byte{}

	// Shift in the bytes
	for absNum > 0 {
		result = append(result, byte(absNum)&0xff)
		absNum >>= 8
	}

	// if the top bit is set,
	// for negative numbers we ensure that the top bit is set
	// for positive numbers we ensure that the top bit is not set
	res := result[len(result)-1] & 0x80
	if res > 0 {
		if negative {
			result = append(result, 0x80)
		} else {
			result = append(result, 0x00)
		}
	} else if negative {
		result[len(result)-1] |= 0x80
	}
	return result
}

// decode and encoded byte array into and int
func decode(element []byte) int {

	// check of the byte array is empty
	// if so the result is 0
	if len(element) == 0 {
		return 0
	}

	// reverse the element to be in big endian
	// (was previously encoded as little endian)
	bigEndian := utils.ReorderBytes(element)

	// negative flag
	var negative bool

	// store the result as a byte and convert to an integer to retun
	var result byte

	// top bit being 1 means its negative
	res := bigEndian[0] & 0x80
	if res > 0 {
		negative = true
	} else {
		result = bigEndian[0] & 0x7f
	}

	for _, b := range bigEndian[1:] {
		result <<= 8
		result += b
	}

	if negative {
		return -int(result)
	} else {
		return int(result)
	}
}

// Push in a raw byte array as a stack element
func (s *Stack) Push(b []byte) {
	s.Elements = append(s.Elements, StackElement{Bytes: b})
}

// Push in a raw byte array and flag as an opcode
func (s *Stack) PushOp(b []byte) {
	s.Elements = append(s.Elements, StackElement{Bytes: b, OpCode: true})
}

// Pops to top stack item
func (s *Stack) Pop() StackElement {
	toPop := s.Elements[len(s.Elements)-1]
	s.Elements = s.Elements[:len(s.Elements)-1]
	return toPop
}

// "Pops" off the head of the stack
func (s *Stack) Head() StackElement {
	toHead := s.Elements[0]
	s.Elements = s.Elements[1:]
	return toHead
}

func (s *Stack) Op0() bool {
	b := encode(0)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op1Negate() bool {
	b := encode(-1)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op1() bool {
	b := encode(2)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op2() bool {
	b := encode(2)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op3() bool {
	b := encode(3)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op4() bool {
	b := encode(4)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op5() bool {
	b := encode(5)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op6() bool {
	b := encode(6)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op7() bool {
	b := encode(7)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op8() bool {
	b := encode(8)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op9() bool {
	b := encode(9)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op10() bool {
	b := encode(10)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op11() bool {
	b := encode(11)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op12() bool {
	b := encode(12)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op13() bool {
	b := encode(13)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op14() bool {
	b := encode(14)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op15() bool {
	b := encode(15)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) Op16() bool {
	b := encode(16)
	s.Elements = append(s.Elements, StackElement{Bytes: b})
	return true
}

func (s *Stack) OpNop() bool {
	return true
}

func (s *Stack) Length() int {
	return len(s.Elements)
}

// If the top stack value is not False, the statements are executed. The top stack value is removed.
func (s *Stack) OpIf(cmds *Stack) bool {
	if len(s.Elements) < 1 {
		return false
	}

	// go through
	var trueItems []StackElement
	var falseItems []StackElement
	currentArray := trueItems
	found := false
	numEndIfsNeeded := 1

	// iterate of the the commands stack
	for i := 0; i < cmds.Length(); i++ {

		// get the top item of the existing stack
		cmd := cmds.Pop()

		// check the opcode of the cmd
		if cmd.OpCode {
			// This is an op code, check if its additional control flow commands
			opCode := binary.BigEndian.Uint32(cmd.Bytes)
			if opCode == 99 || opCode == 100 {
				// if and not if op codes
				numEndIfsNeeded += 1
				currentArray = append(currentArray, StackElement{Bytes: cmd.Bytes, OpCode: cmd.OpCode})
			} else if numEndIfsNeeded == 1 && opCode == 103 {
				// else op code
				currentArray = falseItems
			} else if opCode == 104 {
				// end if op code
				if numEndIfsNeeded == 1 {
					// found the match, time to abort this shit
					found = true
					break
				} else {
					// decrement the counter as we closed out one if statement
					numEndIfsNeeded -= 1
					currentArray = append(currentArray, cmd)
				}
			} else {
				// still looking, save the progress so far
				currentArray = append(currentArray, cmd)
			}
		}
	}
	if !found {
		// indicates we have a mismatched if/else block
		return false
	}
	// we pop off the if statement since we paired it with a closure
	element := s.Pop()

	// check if the stack is 0, indicating false was evaluated of
	// the if statement and assign it to the cmds overwriting the
	// the passed in stacks elements. Otherwise set it to the true branch
	if decode(element.Bytes) == 0 {
		cmds.Elements = falseItems
	} else {
		cmds.Elements = trueItems
	}
	return true
}

func (s *Stack) OpNotIf(cmds *Stack) bool {
	if len(s.Elements) < 1 {
		return false
	}

	// go through
	var trueItems []StackElement
	var falseItems []StackElement
	currentArray := trueItems
	found := false
	numEndIfsNeeded := 1

	// iterate of the the commands stack
	for i := 0; i < cmds.Length(); i++ {

		// get the top item of the existing stack
		cmd := cmds.Pop()

		// check the opcode of the cmd
		if cmd.OpCode {
			// This is an op code, check if its additional control flow commands
			opCode := binary.BigEndian.Uint32(cmd.Bytes)
			if opCode == 99 || opCode == 100 {
				// if and not if op codes
				numEndIfsNeeded += 1
				currentArray = append(currentArray, StackElement{Bytes: cmd.Bytes, OpCode: cmd.OpCode})
			} else if numEndIfsNeeded == 1 && opCode == 103 {
				// else op code
				currentArray = falseItems
			} else if opCode == 104 {
				// end if op code
				if numEndIfsNeeded == 1 {
					// found the match, time to abort this shit
					found = true
					break
				} else {
					// decrement the counter as we closed out one if statement
					numEndIfsNeeded -= 1
					currentArray = append(currentArray, cmd)
				}
			} else {
				// still looking, save the progress so far
				currentArray = append(currentArray, cmd)
			}
		}
	}
	if !found {
		// indicates we have a mismatched if/else block
		return false
	}
	// we pop off the if statement since we paired it with a closure
	element := s.Pop()

	// check if the stack is 0, indicating false was evaluated of
	// the if statement and assign it to the cmds overwriting the
	// the passed in stacks elements. Otherwise set it to the true branch
	if decode(element.Bytes) == 0 {
		cmds.Elements = trueItems
	} else {
		cmds.Elements = falseItems
	}
	return true
}

func (s *Stack) OpVerify() bool {
	if len(s.Elements) < 1 {
		return false
	}

	element := s.Pop()
	if decode(element.Bytes) == 0 {
		return false
	}
	return true
}

func (s *Stack) OpReturn() bool {
	return false
}

func (s *Stack) OpToAltStack(altStack Stack) bool {
	if len(altStack.Elements) < 1 {
		return false
	}
	altStack.Elements = append(altStack.Elements, s.Pop())
	return true
}

func (s *Stack) OpFromAltStack(altStack Stack) bool {
	if len(altStack.Elements) < 1 {
		return false
	}
	s.Elements = append(s.Elements, altStack.Pop())
	return true
}

// Removes the top two stack items.
func (s *Stack) Op2Drop() bool {
	if len(s.Elements) < 2 {
		return false
	}
	s.Pop()
	s.Pop()
	return true
}

// Duplicates the top two stack items.
func (s *Stack) Op2Dup() bool {
	if len(s.Elements) < 2 {
		return false
	}
	pop1 := s.Elements[len(s.Elements)-1]
	pop2 := s.Elements[len(s.Elements)-2]
	s.Elements = append(s.Elements, pop2)
	s.Elements = append(s.Elements, pop1)
	return true
}

// Duplicates the top three stack items.
func (s *Stack) Op3Dup() bool {
	if len(s.Elements) < 3 {
		return false
	}
	pop1 := s.Elements[len(s.Elements)-1]
	pop2 := s.Elements[len(s.Elements)-2]
	pop3 := s.Elements[len(s.Elements)-2]
	s.Elements = append(s.Elements, pop3)
	s.Elements = append(s.Elements, pop2)
	s.Elements = append(s.Elements, pop1)
	return true
}

// Copies the pair of items two spaces back in the stack to the front.
func (s *Stack) Op2Over() bool {
	if len(s.Elements) < 4 {
		return false
	}
	pop3 := s.Elements[len(s.Elements)-3]
	pop4 := s.Elements[len(s.Elements)-4]
	s.Elements = append(s.Elements, pop4)
	s.Elements = append(s.Elements, pop3)
	return true
}

// The fifth and sixth items back are moved to the top of the stack.
func (s *Stack) Op2Rot() bool {
	if len(s.Elements) < 6 {
		return false
	}
	head1 := s.Head()
	head2 := s.Head()
	s.Elements = append(s.Elements, head1)
	s.Elements = append(s.Elements, head2)
	return false
}

// Swaps the top two pairs of items.
func (s *Stack) Op2Swap() bool {
	if len(s.Elements) < 4 {
		return false
	}

	slice1 := s.Elements[:2]
	slice2 := s.Elements[2:4]
	s.Elements = append(slice2, slice1...)

	return true
}

// If the top stack value is not 0, duplicate it.
func (s *Stack) OpIfDup() bool {
	if len(s.Elements) < 1 {
		return false
	}
	if decode(s.Elements[len(s.Elements)-1].Bytes) != 0 {
		s.Elements = append(s.Elements, s.Elements[len(s.Elements)-1])
	}
	return true
}

// Puts the number of stack items onto the stack.
func (s *Stack) OpDepth() bool {
	bn := encode(len(s.Elements))
	s.Elements = append(s.Elements, StackElement{Bytes: bn})
	return true
}

// Removes the top stack item.
func (s *Stack) OpDrop() bool {
	if len(s.Elements) < 1 {
		return false
	}
	s.Pop()
	return true
}

func (s *Stack) OpDup() bool {
	// OpDup will duplidate the top option on the stack

	// check that the stack is not empty
	if len(s.Elements) < 1 {
		return false
	}

	// get the top element of the stack
	toDup := s.Elements[len(s.Elements)-1]
	s.Elements = append(s.Elements, toDup)

	return true
}

// Removes the second-to-top stack item.
func (s *Stack) OpNip() bool {
	if len(s.Elements) < 2 {
		return false
	}
	slice1 := s.Elements[:len(s.Elements)-2]
	s.Elements = append(slice1, s.Elements[len(s.Elements)-1])
	return true
}

// The item n back in the stack is copied to the top.
func (s *Stack) OpPick() bool {
	if len(s.Elements) < 1 {
		return false
	}
	n := decode(s.Pop().Bytes)
	if len(s.Elements) < n+1 {
		return false
	}
	s.Elements = append(s.Elements, s.Elements[-n-1])
	return true
}

// The item n back in the stack is moved to the top.
func (s *Stack) OpRoll() bool {
	if len(s.Elements) < 1 {
		return false
	}
	n := decode(s.Pop().Bytes)
	if len(s.Elements) < n+1 {
		return false
	}
	if n == 0 {
		return true
	}
	idx := -n - 1
	element := s.Elements[idx]
	slice1 := s.Elements[:idx]
	slice2 := s.Elements[idx+1:]
	s.Elements = append(slice1, slice2...)
	s.Elements = append(s.Elements, element)
	return true
}

// The 3rd item down the stack is moved to the top.
func (s *Stack) OpRot() bool {
	if len(s.Elements) < 3 {
		return false
	}

	element := s.Elements[len(s.Elements)-3]
	slice1 := s.Elements[:len(s.Elements)-3]
	slice2 := s.Elements[len(s.Elements)-2 : len(s.Elements)]

	s.Elements = append(slice1, slice2...)
	s.Elements = append(s.Elements, element)
	return true
}

// Copies the second-to-top stack item to the top.
func (s *Stack) OpOver() bool {
	if len(s.Elements) < 2 {
		return false
	}
	s.Elements = append(s.Elements, s.Elements[len(s.Elements)-2])
	return true
}

// The top two items on the stack are swapped.
func (s *Stack) OpSwap() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element := s.Elements[len(s.Elements)-2]
	s.Elements[len(s.Elements)-2] = s.Elements[len(s.Elements)-1]
	s.Elements[len(s.Elements)-1] = element
	return true
}

// The item at the top of the stack is copied and inserted before the second-to-top item.
func (s *Stack) OpTuck() bool {
	if len(s.Elements) < 2 {
		return false
	}

	// get the top item off the stack
	top := s.Elements[len(s.Elements)-1]

	// create two sides of the array
	slice1 := scopy(s.Elements[:len(s.Elements)-2])
	slice2 := scopy(s.Elements[len(s.Elements)-1:])

	// create the first part of the new array
	slice1 = append(slice1, top)

	// append the backend of the slice back in
	slice1 = append(slice1, slice2...)
	s.Elements = slice1
	return true
}

// Pushes the string length of the top element of the stack (without popping it).
func (s *Stack) OpSize() bool {
	if len(s.Elements) < 1 {
		return false
	}
	s.Elements = append(s.Elements, StackElement{Bytes: encode(len(s.Elements[len(s.Elements)-1].Bytes))})
	return true
}

// Returns 1 if the inputs are exactly equal, 0 otherwise.
func (s *Stack) OpEqual() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := s.Pop()
	element2 := s.Pop()

	if bytes.Equal(element1.Bytes, element2.Bytes) {
		// elements are equal
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// Same as OP_EQUAL, but runs OP_VERIFY afterward.
func (s *Stack) OpEqualVerify() bool {
	eq := s.OpEqual()
	vr := s.OpVerify()

	if eq && vr {
		return true
	}
	return false
}

// 1 is added to the input.
func (s *Stack) Op1Add() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := int(decode(s.Pop().Bytes))
	s.Elements = append(s.Elements, StackElement{Bytes: encode(element + 1)})
	return true
}

// 1 is added to the input.
func (s *Stack) Op1Sub() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := int(decode(s.Pop().Bytes))
	s.Elements = append(s.Elements, StackElement{Bytes: encode(element - 1)})
	return true
}

// The sign of the input is flipped.
func (s *Stack) OpNegate() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := int(decode(s.Pop().Bytes))
	s.Elements = append(s.Elements, StackElement{Bytes: encode(-element)})
	return true
}

// The input is made positive.
func (s *Stack) OpAbs() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := int(decode(s.Pop().Bytes))
	if element < 0 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(-element)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(element)})
	}
	return true
}

// If the input is 0 or 1, it is flipped. Otherwise the output will be 0.
func (s *Stack) OpNot() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := int(decode(s.Pop().Bytes))
	if element == 0 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// Returns 0 if the input is 0. 1 otherwise.
func (s *Stack) Op0NotEqual() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := int(decode(s.Pop().Bytes))
	if element == 0 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	}
	return true
}

// Adds the first two inputs together
func (s *Stack) OpAdd() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))
	s.Elements = append(s.Elements, StackElement{Bytes: encode(element1 + element2)})
	return true
}

// Substracts the first two inputs together
func (s *Stack) OpSub() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))
	s.Elements = append(s.Elements, StackElement{Bytes: encode(element2 + element1)})
	return true
}

// If both a and b are not 0, the output is 1. Otherwise 0.
func (s *Stack) OpBoolAnd() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element1 != 0 && element2 != 0 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// If both a and b are not 0, the output is 1. Otherwise 0.
func (s *Stack) OpBoolOr() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element1 != 0 || element2 != 0 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// Returns 1 if the numbers are equal, 0 otherwise.
func (s *Stack) OpNumEqual() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element1 == element2 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// Same as OP_NUMEQUAL, but runs OP_VERIFY afterward.
func (s *Stack) OpNumEqualVerify() bool {
	return s.OpNumEqual() && s.OpVerify()
}

// Returns 1 if the numbers are not equal, 0 otherwise.
func (s *Stack) OpNumNotEqual() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element1 == element2 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	}
	return true
}

// Returns 1 if a is less than b, 0 otherwise.
func (s *Stack) OpLessThan() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element2 < element1 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// Returns 1 if a is greater than b, 0 otherwise.
func (s *Stack) OpGreaterThan() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element2 < element1 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// Returns 1 if a is less than or equal to b, 0 otherwise.
func (s *Stack) OpLessOrEqualThan() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element2 >= element1 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// Returns 1 if a is greater than or equal to b, 0 otherwise.
func (s *Stack) OpGreaterOrEqualThan() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element2 <= element1 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// Returns the smaller of a and b.
func (s *Stack) OpMin() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element1 < element2 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(element1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(element2)})
	}
	return true
}

// Returns the larger of a and b.
func (s *Stack) OpMax() bool {
	if len(s.Elements) < 2 {
		return false
	}
	element1 := int(decode(s.Pop().Bytes))
	element2 := int(decode(s.Pop().Bytes))

	if element1 > element2 {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(element1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(element2)})
	}
	return true
}

// Returns 1 if x is within the specified range (left-inclusive), 0 otherwise.
func (s *Stack) OpWithIn() bool {
	if len(s.Elements) < 3 {
		return false
	}
	max := int(decode(s.Pop().Bytes))
	min := int(decode(s.Pop().Bytes))
	element := int(decode(s.Head().Bytes))

	if element >= min && element < max {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// The input is hashed using RIPEMD-160
func (s *Stack) OpRipeMd160() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := s.Pop()
	sum := ripemd160.New().Sum(element.Bytes)
	s.Elements = append(s.Elements, StackElement{Bytes: sum[:]})
	return true
}

// The input is hashed using SHA-1.
func (s *Stack) OpSha1() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := s.Pop()
	sum := sha1.Sum(element.Bytes)
	s.Elements = append(s.Elements, StackElement{Bytes: sum[:]})
	return true
}

// The input is hashed using SHA-256.
func (s *Stack) OpSha256() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := s.Pop()
	sum := sha256.Sum256(element.Bytes)
	s.Elements = append(s.Elements, StackElement{Bytes: sum[:]})
	return true
}

// The input is hashed using HASH-160.
func (s *Stack) OpHash160() bool {
	if len(s.Elements) < 1 {
		return false
	}
	element := s.Pop()
	sum := utils.Hash160(element.Bytes)
	s.Elements = append(s.Elements, StackElement{Bytes: sum[:]})
	return true
}

func (s *Stack) OpHash256() bool {

	// Check that the stack is not empty
	if len(s.Elements) < 1 {
		return false
	}

	// get the top element off the stack
	se := s.Pop()

	// hash the element
	s.Push(utils.Hash256(se.Bytes))

	return true

}

// The entire transaction's outputs, inputs, and script (from the most recently-executed OP_CODESEPARATOR to the end)
// are hashed. The signature used by OP_CHECKSIG must be a valid signature for this hash and public key.
// If it is, 1 is returned, 0 otherwise.
func (s *Stack) OpCheckSig(z *big.Int) bool {
	if len(s.Elements) > 2 {
		return false
	}

	// get the sec formatted pub key
	secPubKey := s.Pop()

	// get the signature in der formation
	derSignature := s.Pop()

	// shave off the SIGHASH flag
	derSignature.Bytes = derSignature.Bytes[:len(derSignature.Bytes)-1]

	// using ECC lib, verify the sec pubkey and the associated der signature
	point, err := S256.ParseSec(secPubKey.Bytes)
	if err != nil {
		fmt.Printf("Failed to parse the secPubKey because %v\n", err.Error())
		return false
	}
	sig, err := S256.ParseSignature(derSignature.Bytes)
	if err != nil {
		fmt.Printf("Failed to parse the signature because %v\n", err.Error())
		return false
	}

	// With the secPubKey and the Der signature, we can verify that point, z,
	// is on the curve and ready to go
	result, err := point.Verify(*z, *sig)
	if err != nil {
		fmt.Printf("failed to verify point because %s\n", err.Error())
		return false
	}

	if result {
		// signature is validates, so append a true (1) value to the stack
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}

	return true
}

// OP_CHECKSIG coupled with OP_VERIFY
func (s *Stack) OpCheckSigVerify(z *big.Int) bool {
	return s.OpCheckSig(z) && s.OpVerify()
}

// Compares the first signature against each public key until it finds an ECDSA match.
// Starting with the subsequent public key, it compares the second signature against each
// remaining public key until it finds an ECDSA match. The process is repeated until all
// signatures have been checked or not enough public keys remain to produce a successful result.
// All signatures need to match a public key. Because public keys are not checked again if
// they fail any signature comparison, signatures must be placed in the scriptSig using
// the same order as their corresponding public keys were placed in the scriptPubKey or
// redeemScript. If all signatures are valid, 1 is returned, 0 otherwise. Due to a bug,
// one extra unused value is removed from the stack.
func (s *Stack) OpCheckMultisig(z *big.Int) bool {
	if len(s.Elements) < 1 {
		return false
	}

	// get the n number of multi sig addresses required for a validation
	n := decode(s.Pop().Bytes)

	// make sure the stack is large enough to hold n multi sigs in the first place
	if len(s.Elements) < n+1 {
		return false
	}

	// pop all the pub keys into an array
	secPubKeys := []StackElement{}
	for i := 0; i < n; i++ {
		secPubKeys = append(secPubKeys, s.Pop())
	}

	// get the number of multi sigs required (m)
	m := decode(s.Pop().Bytes)

	// make sure the stack has enough m signatures for the op
	if len(s.Elements) < m+1 {
		return false
	}

	// get all the signatures pushed onto the stack for the op
	derSigs := []StackElement{}
	for i := 0; i < m; i++ {
		// one thing to note is the last big the der signatures is the
		// SIGHASH_ALL flag (or variant of)
		derSigs = append(derSigs, s.Pop())
	}

	// off by one error coded into the bitcoin standard.
	// a 1 is always pushed in, but can really be anything
	// we will just pop it off anyway
	s.Pop()

	// verify that there are n of m valid signatures.
	matches := 0
	for _, d := range derSigs {
		// parse the der sig into a der object
		// Strip of the SIGHASH_ALL flag from the signatures
		sig, err := S256.ParseSignature(d.Bytes[:len(d.Bytes)-1])
		if err != nil {
			fmt.Printf("failed to parse der siganature because %s\n", sig)
			return false
		}

		// range over the sec pub keys and see if we can find a match
		// e.g. verify passes.
		for _, s := range secPubKeys {

			// parse the secPubKey into the struct
			p, err := S256.ParseSec(s.Bytes)
			if err != nil {
				fmt.Printf("failed to parse SEC point because %s\n", err.Error())
				return false
			}

			// try and verify the signature
			result, err := p.Verify(*z, *sig)
			if err != nil {
				fmt.Printf("failed verification of der signature because %s\n", err.Error())
				return false
			}

			if result {
				matches += 1
				break
			}
		}
	}

	// check if we have n of m sigantures
	if matches >= n {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(1)})
	} else {
		s.Elements = append(s.Elements, StackElement{Bytes: encode(0)})
	}
	return true
}

// Same as OP_CHECKMULTISIG, but OP_VERIFY is executed afterward.
func (s *Stack) OpCheckMultisigVerify(z *big.Int) bool {
	return s.OpCheckMultisig(z) && s.OpVerify()
}

// Marks transaction as invalid if the top stack item is greater than the transaction's nLockTime field,
// otherwise script evaluation continues as though an OP_NOP was executed. Transaction is also invalid
// if 1. the stack is empty; or 2. the top stack item is negative; or 3. the top stack item is greater
// than or equal to 500000000 while the transaction's nLockTime field is less than 500000000, or vice versa;
//  or 4. the input's nSequence field is equal to 0xffffffff. The precise semantics are described in BIP 0065.
func (s *Stack) OpCheckTimeLockVerify(locktime, sequence uint64) bool {
	if sequence == 0xffffffff {
		return false
	}
	if len(s.Elements) < 1 {
		return false
	}
	element := decode(s.Elements[len(s.Elements)-1].Bytes)
	if element < 0 {
		return false
	}
	if element < 500000000 && locktime > 500000000 {
		return false
	}
	if locktime < uint64(element) {
		return false
	}
	return true
}

// Marks transaction as invalid if the relative lock time of the input
// (enforced by BIP 0068 with nSequence) is not equal to or longer than the value of the top stack item.
// the precise semantics are described in BIP 0112.
func (s *Stack) OpCheckSequenceVerify(version, sequence uint64) bool {
	if sequence&(1<<32) == (1 << 31) {
		return false
	}
	if len(s.Elements) < 1 {
		return false
	}
	element := decode(s.Elements[len(s.Elements)-1].Bytes)

	if element < 0 {
		return false
	}
	if element&(1<<31) == (1 << 31) {
		if version < 2 {
			return false
		} else if sequence&(1<<31) == (1 << 31) {
			return false
		} else if element&(1<<22) == (1 << 22) {
			return false
		} else if uint64(element)&0xffff > sequence&0xffff {
			return false
		}
	}
	return true
}
