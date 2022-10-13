package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/script/opcodes"
)

type Command struct {
	// raw bytes for the command
	Bytes []byte

	// flag to indicate if this is an opcode
	OpCode bool
}

type Script struct {
	// serialized byte array of the full script
	RawScript []byte

	// Each "command" (opcode or element) parsed
	Commands []Command
}

func MakeScript() *Script {
	return &Script{
		RawScript: []byte{0x00},
		Commands:  nil,
	}
}

// Combine a scriptSig with a scriptPubKey resulting in a new script with both elements
func Combine(scriptPubkey, scriptSig Script) *Script {
	// create a new command list
	cmds := scriptSig.Commands
	cmds = append(cmds, scriptPubkey.Commands...)

	return &Script{
		Commands: cmds,
	}
}

// Parses a byte stream into a script
func Parse(reader *bytes.Reader) (*Script, error) {
	// read the length of the script,
	// a script always starts with the overal length of the script
	length, err := binary.ReadUvarint(reader)

	if err != nil {
		return nil, err
	}

	// commands array we will parse everyting into. Its an array of byte arrays
	var commands []Command

	// init the numer of bytes being read to 0
	count := uint64(0)

	// loop for all the bytes in the stream
	for {

		// break condition is once all the bytes have been read
		// count will == length.
		if count >= length {
			break
		}

		// read the first byte which determines if we have an opcode or an element
		current, err := reader.ReadByte()

		if err != nil {
			// break the loop if we run out of data to read regardless if we finished or not
			// in a perfect world, we would have exited this at the first circuit breaker
			// for the for loop, but this handles the case were the supplied data length does
			// not match the actual length of the data.
			if err == io.EOF {
				break
			}
		}

		// increment by 1 bytes that we read
		count += 1

		// convert the current byte to an integer
		// which is the length of the element
		currentByte := uint64(current)

		// Checks if this is an element which is defined as a size of
		// 1 - 75. If it is an element (data), the byte indicates the
		// size of the element
		if currentByte >= 1 && currentByte <= 75 {

			// get the length of the element
			elementLength := currentByte

			// create a stack item with the element read direct from the stream
			lr := io.LimitReader(reader, int64(elementLength))
			cmd, err := ioutil.ReadAll(lr)

			if err != nil {
				// on error dont append any new command
				fmt.Printf("Failed to read because %v\n", err.Error())
			} else {
				// append element to the list of commands (stack items)
				commands = append(
					commands,
					Command{
						Bytes:  cmd,
						OpCode: false,
					},
				)
			}

			count += elementLength

		} else if currentByte == 76 {
			// OP_PUSHDATA1 - next byte is the size to read

			// data length is stored as a little endian byte
			dataLengthLittleEndian, _ := reader.ReadByte()
			dataLength := binary.LittleEndian.Uint32([]byte{dataLengthLittleEndian})

			// create a command with the raw data read direct from the stream
			lr := io.LimitReader(reader, int64(dataLength))
			cmd, err := ioutil.ReadAll(lr)

			if err != nil {

				// on error dont append any new command
				fmt.Printf("Failed to read because %v\n", err.Error())
			} else {

				// append the data element to the commands list
				commands = append(
					commands,
					Command{
						Bytes:  cmd,
						OpCode: false,
					},
				)
			}

			count += uint64(dataLength)

		} else if currentByte == 77 {
			// OP_PUSHDATA2 - Next 2 bytes says how long the data to read is

			// data lenght is stored as a little endian byte
			lr := io.LimitReader(reader, int64(2))
			dataLengthLittleEndian, _ := ioutil.ReadAll(lr)
			dataLength := binary.LittleEndian.Uint64(dataLengthLittleEndian)

			// create a command with the raw data read direct from the stream
			lr = io.LimitReader(reader, int64(dataLength))
			cmd, err := ioutil.ReadAll(lr)

			if err != nil {

				// on error dont append any new command
				fmt.Printf("Failed to read because %v\n", err.Error())
			} else {

				// append the validated command
				commands = append(
					commands,
					Command{
						Bytes:  cmd,
						OpCode: false,
					},
				)
			}
			count += dataLength
		} else {

			// this indicates there is an opcode (single byte) which needs to be stored
			opCode := currentByte
			commands = append(
				commands,
				Command{
					Bytes:  []byte{byte(opCode)},
					OpCode: true,
				},
			)
		}
	}

	// should consume exactly the length of bytes expected otherwise
	// there is an error condition which needs to be handled
	if count != length {
		return nil, fmt.Errorf("failed to parse script because count does not match the read bytes")
	}

	// create a script object with the specified commands
	return &Script{
		Commands: commands,
	}, nil
}

func (s Script) RawSerialize() ([]byte, error) {

	// result byte array which will be the serialized stream
	var result []byte

	// iterate over the commands
	for _, c := range s.Commands {
		v := c.Bytes

		// if the command is a single byte, then its the opcode
		if c.OpCode {
			result = append(result, v[0])
		} else {

			// otherwize its the
			length := len(v)
			if length < 75 {

				// if length is < 75, it gets encoded as a single byte
				result = append(result, byte(length))
			} else if length > 75 && length < 0x100 {

				// check if it is OP_PUSHDATA1 (len=76), if it is then encode the length
				// as a single byte, but need to push in the the next element
				result = append(result, byte(76))
				result = append(result, byte(length))
			} else if length > 0x100 && length <= 510 {

				// OP_PUSHDATA2, the length needs to be 2 bytes encoded little endian
				// then encode the element
				result = append(result, byte(77))

				// doing this the dirty way
				var length2 []byte
				binary.BigEndian.PutUint64(length2, uint64(length))
				result = append(result, length2...)
			} else {
				return nil, fmt.Errorf("element is longer than 520 bytes and cannot be seralized")
			}

			// finally, add the data to the command
			result = append(result, v...)
		}
	}

	// return the result
	return result, nil
}

// Returns a serialized byte array containing the script
func (s Script) Serialize() []byte {

	// get the serialized script
	rawScript, err := s.RawSerialize()

	if err != nil {
		return nil
	}

	// scripts always start with the length of the full script
	fullScript := []byte{byte(len(rawScript))}
	fullScript = append(fullScript, rawScript...)

	return fullScript
}

func (s *Script) Evaluate(z *big.Int, locktime, sequence, version uint64) {

	// Commands list will change so we need to make a local copy
	cmds := s.Commands

	// executable stack to be created
	var stack opcodes.Stack
	var altStack opcodes.Stack

	// range over the commands and process each command peice by peice
	for _, c := range cmds {

		// check if the command is an opcode
		if c.OpCode {
			switch opCode := binary.BigEndian.Uint32(c.Bytes); opCode {
			case opcodes.OP_0:
				stack.Op0()
			case opcodes.OP_PUSHDATA1:
			case opcodes.OP_PUSHDATA2:
			case opcodes.OP_PUSHDATA4:
				break
			case opcodes.OP_1NEGATE:
				stack.Op1Negate()
			case opcodes.OP_1:
				stack.Op1()
			case opcodes.OP_2:
				stack.Op2()
			case opcodes.OP_3:
				stack.Op3()
			case opcodes.OP_4:
				stack.Op4()
			case opcodes.OP_5:
				stack.Op5()
			case opcodes.OP_6:
				stack.Op6()
			case opcodes.OP_7:
				stack.Op7()
			case opcodes.OP_8:
				stack.Op8()
			case opcodes.OP_9:
				stack.Op9()
			case opcodes.OP_10:
				stack.Op10()
			case opcodes.OP_11:
				stack.Op11()
			case opcodes.OP_12:
				stack.Op12()
			case opcodes.OP_13:
				stack.Op13()
			case opcodes.OP_14:
				stack.Op14()
			case opcodes.OP_15:
				stack.Op15()
			case opcodes.OP_16:
				stack.Op16()
			case opcodes.OP_NOP:
				stack.OpNop()
			case opcodes.OP_IF:
				fmt.Println("not implemented")
			case opcodes.OP_NOTIF:
				fmt.Println("not implemented")

			case opcodes.OP_ELSE:
				fmt.Println("not implemented")

			case opcodes.OP_ENDIF:
				fmt.Println("not implemented")

			case opcodes.OP_VERIFY:
				stack.OpVerify()
			case opcodes.OP_RETURN:
				stack.OpReturn()
			case opcodes.OP_TOALTSTACK:
				stack.OpToAltStack(altStack)
			case opcodes.OP_FROMALTSTACK:
				stack.OpFromAltStack(altStack)
			case opcodes.OP_2DROP:
				stack.Op2Drop()
			case opcodes.OP_2DUP:
				stack.Op2Dup()
			case opcodes.OP_3DUP:
				stack.Op3Dup()
			case opcodes.OP_2OVER:
				stack.Op2Over()
			case opcodes.OP_2ROT:
				stack.Op2Rot()
			case opcodes.OP_2SWAP:
				stack.Op2Swap()
			case opcodes.OP_IFDUP:
				stack.OpIfDup()
			case opcodes.OP_DEPTH:
				stack.OpDepth()
			case opcodes.OP_DROP:
				stack.OpDrop()
			case opcodes.OP_DUP:
				stack.OpDup()
			case opcodes.OP_NIP:
				stack.OpNip()
			case opcodes.OP_OVER:
				stack.OpOver()
			case opcodes.OP_PICK:
				stack.OpPick()
			case opcodes.OP_ROLL:
				stack.OpRoll()
			case opcodes.OP_ROT:
				stack.OpRot()
			case opcodes.OP_SWAP:
				stack.OpSwap()
			case opcodes.OP_TUCK:
				stack.OpTuck()
			case opcodes.OP_SIZE:
				stack.OpSize()
			case opcodes.OP_EQUAL:
				stack.OpEqual()
			case opcodes.OP_EQUALVERIFY:
				stack.OpEqualVerify()
			case opcodes.OP_1ADD:
				stack.Op1Add()
			case opcodes.OP_1SUB:
				stack.Op1Sub()
			case opcodes.OP_NEGATE:
				stack.OpNegate()
			case opcodes.OP_ABS:
				stack.OpAbs()
			case opcodes.OP_NOT:
				stack.OpNot()
			case opcodes.OP_0NOTEQUAL:
				stack.Op0NotEqual()
			case opcodes.OP_ADD:
				stack.OpAdd()
			case opcodes.OP_SUB:
				stack.OpSub()
			case opcodes.OP_BOOLAND:
				stack.OpBoolAnd()
			case opcodes.OP_BOOLOR:
				stack.OpBoolOr()
			case opcodes.OP_NUMEQUAL:
				stack.OpNumEqual()
			case opcodes.OP_NUMEQUALVERIFY:
				stack.OpNumEqualVerify()
			case opcodes.OP_NUMNOTEQUAL:
				stack.OpNumNotEqual()
			case opcodes.OP_LESSTHAN:
				stack.OpLessThan()
			case opcodes.OP_GREATERTHAN:
				stack.OpGreaterThan()
			case opcodes.OP_LESSTHANOREQUAL:
				stack.OpLessOrEqualThan()
			case opcodes.OP_GREATERTHANOREQUAL:
				stack.OpGreaterOrEqualThan()
			case opcodes.OP_MIN:
				stack.OpMin()
			case opcodes.OP_MAX:
				stack.OpMax()
			case opcodes.OP_WITHIN:
				stack.OpWithIn()
			case opcodes.OP_RIPEMD160:
				stack.OpRipeMd160()
			case opcodes.OP_SHA1:
				stack.OpSha1()
			case opcodes.OP_SHA256:
				stack.OpSha256()
			case opcodes.OP_HASH160:
				stack.OpHash160()
			case opcodes.OP_HASH256:
				stack.OpHash256()
			case opcodes.OP_CODESEPARATOR:
				break
			case opcodes.OP_CHECKSIG:
				stack.OpCheckSig(z)
			case opcodes.OP_CHECKSIGVERIFY:
				stack.OpCheckSigVerify(z)
			case opcodes.OP_CHECKMULTISIG:
				stack.OpCheckMultisig(z)
			case opcodes.OP_CHECKMULTISIGVERIFY:
				stack.OpCheckMultisigVerify(z)
			case opcodes.OP_CHECKLOCKTIMEVERIFY:
				stack.OpCheckTimeLockVerify(locktime, sequence)
			case opcodes.OP_CHECKSEQUENCEVERIFY:
				stack.OpCheckSequenceVerify(version, sequence)
			case opcodes.OP_NOP1:
			case opcodes.OP_NOP4:
			case opcodes.OP_NOP5:
			case opcodes.OP_NOP6:
			case opcodes.OP_NOP7:
			case opcodes.OP_NOP8:
			case opcodes.OP_NOP9:
			case opcodes.OP_NOP10:
				break
			}
		} else {
			// Tsaaasasahe selement is not an op code and is thus a data element
			stack.Push(c.Bytes)
		}
	}
}
