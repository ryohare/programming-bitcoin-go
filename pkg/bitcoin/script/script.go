package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/script/opcodes"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
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

// conditions (ifs) require the commands array, but
// the stack its a layer below so we need repackage
// the commands as the raw byte arrays to make it
// something the dump stack layer understands
func makeBinaryCommands(cmds []Command) *opcodes.Stack {
	// b := make([]byte, len(cmds))
	// for _, c := range cmds {
	// 	b = append(b, c.Bytes...)
	// }
	// return b
	stack := &opcodes.Stack{}
	for _, c := range cmds {
		stack.Push(c.Bytes)
	}
	return stack
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

func (s *Script) Evaluate(z *big.Int, locktime, sequence, version uint64) bool {

	// Commands list will change so we need to make a local copy
	cmds := s.Commands

	// executable stack to be created
	var stack opcodes.Stack
	var altStack opcodes.Stack
	result := true

	// range over the commands and process each command peice by peice
	for i, c := range cmds {

		if !result {
			break
		}

		// before we hit the main swith block, we need to inspect the stack for
		// specific signatures indicating p2sh, segwit, etc

		// check if the command is an opcode
		if c.OpCode {
			switch opCode := binary.BigEndian.Uint32(c.Bytes); opCode {
			case opcodes.OP_0:
				result = stack.Op0()
			case opcodes.OP_PUSHDATA1:
			case opcodes.OP_PUSHDATA2:
			case opcodes.OP_PUSHDATA4:
				break
			case opcodes.OP_1NEGATE:
				result = stack.Op1Negate()
			case opcodes.OP_1:
				result = stack.Op1()
			case opcodes.OP_2:
				result = stack.Op2()
			case opcodes.OP_3:
				result = stack.Op3()
			case opcodes.OP_4:
				result = stack.Op4()
			case opcodes.OP_5:
				result = stack.Op5()
			case opcodes.OP_6:
				result = stack.Op6()
			case opcodes.OP_7:
				result = stack.Op7()
			case opcodes.OP_8:
				result = stack.Op8()
			case opcodes.OP_9:
				result = stack.Op9()
			case opcodes.OP_10:
				result = stack.Op10()
			case opcodes.OP_11:
				result = stack.Op11()
			case opcodes.OP_12:
				result = stack.Op12()
			case opcodes.OP_13:
				result = stack.Op13()
			case opcodes.OP_14:
				result = stack.Op14()
			case opcodes.OP_15:
				result = stack.Op15()
			case opcodes.OP_16:
				result = stack.Op16()
			case opcodes.OP_NOP:
				result = stack.OpNop()
			case opcodes.OP_IF:
				// conditions (ifs) require the commands array, but
				// the stack its a layer below so we need repackage
				// the commands as the raw byte arrays to make it
				// something the dump stack layer understands
				result = stack.OpIf(makeBinaryCommands(cmds[i : len(cmds)-1]))
			case opcodes.OP_NOTIF:
				// same as above
				result = stack.OpNotIf(makeBinaryCommands(cmds[i : len(cmds)-1]))
			case opcodes.OP_ELSE:
				// this is handled during and op if we find this out here, something
				// fucked up and we need to get out
				result = false
			case opcodes.OP_ENDIF:
				// same as above
				result = false
			case opcodes.OP_VERIFY:
				result = stack.OpVerify()
			case opcodes.OP_RETURN:
				result = stack.OpReturn()
			case opcodes.OP_TOALTSTACK:
				result = stack.OpToAltStack(altStack)
			case opcodes.OP_FROMALTSTACK:
				result = stack.OpFromAltStack(altStack)
			case opcodes.OP_2DROP:
				result = stack.Op2Drop()
			case opcodes.OP_2DUP:
				result = stack.Op2Dup()
			case opcodes.OP_3DUP:
				result = stack.Op3Dup()
			case opcodes.OP_2OVER:
				result = stack.Op2Over()
			case opcodes.OP_2ROT:
				result = stack.Op2Rot()
			case opcodes.OP_2SWAP:
				result = stack.Op2Swap()
			case opcodes.OP_IFDUP:
				result = stack.OpIfDup()
			case opcodes.OP_DEPTH:
				result = stack.OpDepth()
			case opcodes.OP_DROP:
				result = stack.OpDrop()
			case opcodes.OP_DUP:
				result = stack.OpDup()
			case opcodes.OP_NIP:
				result = stack.OpNip()
			case opcodes.OP_OVER:
				result = stack.OpOver()
			case opcodes.OP_PICK:
				result = stack.OpPick()
			case opcodes.OP_ROLL:
				result = stack.OpRoll()
			case opcodes.OP_ROT:
				result = stack.OpRot()
			case opcodes.OP_SWAP:
				result = stack.OpSwap()
			case opcodes.OP_TUCK:
				result = stack.OpTuck()
			case opcodes.OP_SIZE:
				result = stack.OpSize()
			case opcodes.OP_EQUAL:
				result = stack.OpEqual()
			case opcodes.OP_EQUALVERIFY:
				result = stack.OpEqualVerify()
			case opcodes.OP_1ADD:
				result = stack.Op1Add()
			case opcodes.OP_1SUB:
				result = stack.Op1Sub()
			case opcodes.OP_NEGATE:
				result = stack.OpNegate()
			case opcodes.OP_ABS:
				result = stack.OpAbs()
			case opcodes.OP_NOT:
				result = stack.OpNot()
			case opcodes.OP_0NOTEQUAL:
				result = stack.Op0NotEqual()
			case opcodes.OP_ADD:
				result = stack.OpAdd()
			case opcodes.OP_SUB:
				result = stack.OpSub()
			case opcodes.OP_BOOLAND:
				result = stack.OpBoolAnd()
			case opcodes.OP_BOOLOR:
				result = stack.OpBoolOr()
			case opcodes.OP_NUMEQUAL:
				result = stack.OpNumEqual()
			case opcodes.OP_NUMEQUALVERIFY:
				result = stack.OpNumEqualVerify()
			case opcodes.OP_NUMNOTEQUAL:
				result = stack.OpNumNotEqual()
			case opcodes.OP_LESSTHAN:
				result = stack.OpLessThan()
			case opcodes.OP_GREATERTHAN:
				result = stack.OpGreaterThan()
			case opcodes.OP_LESSTHANOREQUAL:
				result = stack.OpLessOrEqualThan()
			case opcodes.OP_GREATERTHANOREQUAL:
				result = stack.OpGreaterOrEqualThan()
			case opcodes.OP_MIN:
				result = stack.OpMin()
			case opcodes.OP_MAX:
				result = stack.OpMax()
			case opcodes.OP_WITHIN:
				result = stack.OpWithIn()
			case opcodes.OP_RIPEMD160:
				result = stack.OpRipeMd160()
			case opcodes.OP_SHA1:
				result = stack.OpSha1()
			case opcodes.OP_SHA256:
				result = stack.OpSha256()
			case opcodes.OP_HASH160:
				result = stack.OpHash160()
			case opcodes.OP_HASH256:
				result = stack.OpHash256()
			case opcodes.OP_CODESEPARATOR:
				break
			case opcodes.OP_CHECKSIG:
				result = stack.OpCheckSig(z)
			case opcodes.OP_CHECKSIGVERIFY:
				result = stack.OpCheckSigVerify(z)
			case opcodes.OP_CHECKMULTISIG:
				result = stack.OpCheckMultisig(z)
			case opcodes.OP_CHECKMULTISIGVERIFY:
				result = stack.OpCheckMultisigVerify(z)
			case opcodes.OP_CHECKLOCKTIMEVERIFY:
				result = stack.OpCheckTimeLockVerify(locktime, sequence)
			case opcodes.OP_CHECKSEQUENCEVERIFY:
				result = stack.OpCheckSequenceVerify(version, sequence)
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
			// is not an op code and is thus a data element
			stack.Push(c.Bytes)

			// check for p2sh signature pattern. It will occur just after we push a the redeem script
			// onto the stack, this branch, and we want to read ahead the list of commands and
			// see if the signature occurs
			if len(cmds[i:]) == 3 &&
				uint32(cmds[0].Bytes[0]) == opcodes.OP_HASH160 &&
				cmds[1].OpCode != true &&
				len(cmds[1].Bytes) == 20 &&
				uint32(cmds[2].Bytes[0]) == opcodes.OP_EQUAL {

				// indicates we have a p2sh signature left in the command list
				// dont care about HASH_160 opcode, care about the hash value
				// and dont are about the OP_EQUAL
				h160 := cmds[1].Bytes

				// have the 20 byte hash160 that was supplied, hash the redeem scrip twhich was pushed
				// into the stack just a second ago
				if !stack.OpHash160() {
					return false
				}

				// push the h160 onto the stack and evaulate it
				stack.Push(h160)

				// verify that the hashes matched with the op verify command
				if !stack.OpEqual() {
					return false
				}

				// no execute the op verify to verify the signature
				if !stack.OpVerify() {
					fmt.Printf("Failed to verify p2sh h160")
					return false
				}

				// insert into the stack the varint length of the the redeem script and insert it at
				// the front of the stack for parsing
				length, err := utils.EncodeUVarInt(uint64(len(c.Bytes)))
				if err != nil {
					fmt.Printf("failed to encode redeem script length because %s\n", err.Error())
					return false
				}
				redeemScript := append(length, c.Bytes...)

				// redeem script is its own script now in byte serialized. Need to parse it into
				// a script object so we can process it fully
				script, err := Parse(bytes.NewReader(redeemScript))
				if err != nil {
					fmt.Printf("failed to parse redeem script because %s\n", err.Error())
					return false
				}

				// extend the redeem script into the commands array so this evaulate will continue
				// with more elements in the cmds array
				cmds = append(cmds, script.Commands...)
			}
		}
	}
	return result
}

// Checks if the pubkey for the script is a P2PKH
func (s Script) IsP2pkhScriptPubkey() bool {
	// this will check if the list of commands held by the script
	// match those for a P2PKH. This should be exactly 5 instructions
	// OP_DUP
	// OP_HASH160
	// Script Pubkey
	// OP_EQUALVERIFY
	// OP_CHECKSIG

	if len(s.Commands) != 5 {
		return false
	}
	if s.Commands[0].Bytes[0] != byte(opcodes.OP_DUP) {
		return false
	}
	if s.Commands[1].Bytes[0] != byte(opcodes.OP_HASH160) {
		return false
	}
	if len(s.Commands[2].Bytes) != 20 {
		return false
	}
	if s.Commands[3].Bytes[0] != byte(opcodes.OP_EQUALVERIFY) {
		return false
	}
	if s.Commands[4].Bytes[0] != byte(opcodes.OP_CHECKSIG) {
		return false
	}
	return true
}

// Checks if the pubkey for the script is P2SH
func (s Script) IsP2shScriptPubkey() bool {
	// This will check if the list of commands helpd by the script
	// match those for a P2SH. Thias should be exactly 3 commands
	// OP_HASH160
	// Script Pubkey
	// OP_EQAL

	if len(s.Commands) != 3 {
		return false
	}
	if s.Commands[0].Bytes[0] != byte(opcodes.OP_HASH160) {
		return false
	}
	if len(s.Commands[1].Bytes) != 20 {
		return false
	}
	if s.Commands[2].Bytes[0] != byte(opcodes.OP_EQUAL) {
		return false
	}
	return true
}

func MakeP2pkh(h160 []byte) *Script {
	cmds := []Command{}

	cmds = append(cmds, Command{
		Bytes:  []byte{byte(opcodes.OP_DUP)},
		OpCode: true,
	})
	cmds = append(cmds, Command{
		Bytes:  []byte{byte(opcodes.OP_HASH160)},
		OpCode: true,
	})

	cmds = append(cmds, Command{
		Bytes:  h160,
		OpCode: false,
	})

	cmds = append(cmds, Command{
		Bytes:  []byte{byte(opcodes.OP_EQUALVERIFY)},
		OpCode: true,
	})

	cmds = append(cmds, Command{
		Bytes:  []byte{byte(opcodes.OP_CHECKSIG)},
		OpCode: true,
	})

	return &Script{Commands: cmds}
}
