package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

type Script struct {
	RawScript []byte
	Commands  [][]byte
}

func MakeScript() *Script {
	return &Script{
		RawScript: []byte{0x00},
		Commands:  nil,
	}
}

// func (s Script) String() {
// 	var result []string
// 	for _,cmd :=range s.Commands{
// 		// ops have to formats, named or number op codes.
// 	}
// }

// Receives
func Parse(reader *bytes.Reader) (*Script, error) {
	// read the length of the script,
	// a script always starts with the overal length of the script
	length, err := binary.ReadUvarint(reader)

	if err != nil {
		return nil, err
	}

	// commands array we will parse everyting into. Its an array of byte arrays
	var commands [][]byte

	// init the numer of bytes being read to 0
	count := uint64(0)

	// loop for all the bytes in the stream
	for {

		// break condition is once all the bytes have been read
		// count will == length.
		// might need a saftefy valve here incase of a bad script
		// otherwise this might loop for ever
		if count >= length {
			break
		} else {
			// safety valve.
			// check if the stream is empty, if it is, break here
		}

		// read the first byte which determines if we have an opcode or an element
		current, _ := reader.ReadByte()

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
				commands = append(commands, cmd)
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
				commands = append(commands, cmd)
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
				commands = append(commands, cmd)
			}
			count += dataLength
		} else {

			// this indicates there is an opcode (single byte) which needs to be stored
			opCode := currentByte
			commands = append(commands, []byte{byte(opCode)})
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
	for _, v := range s.Commands {

		// if the command is a single byte, then its the opcode
		if len(v) == 1 {
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
