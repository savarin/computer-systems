package vm

import (
	"fmt"
)

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//
func compute(memory []byte) {

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	// Keep looping, like a physical computer's clock
	for {
		pc := registers[0]
		op := memory[pc] // fetch the opcode

		// decode and execute
		switch op {
		case Load:
			register, address := memory[pc+1], memory[pc+2]

			value := memory[address]
			registers[register] = value

		case Store:
			register, address := memory[pc+1], memory[pc+2]

			value := registers[register]
			memory[address] = value

		case Add:
			register1, register2 := memory[pc+1], memory[pc+2]
			registers[register1] += registers[register2]

		case Sub:
			register1, register2 := memory[pc+1], memory[pc+2]
			registers[register1] -= registers[register2]

		case Addi:
			register, i := memory[pc+1], memory[pc+2]
			registers[register] += i

		case Subi:
			register, i := memory[pc+1], memory[pc+2]
			registers[register] -= i

		case Jump:
			address := memory[pc+1]
			registers[0] = address
			continue

		case Beqz:
			register, offset := memory[pc+1], memory[pc+2]

			if registers[register] == 0 {
				registers[0] += offset + 3
				continue
			}

		case Halt:
			return

		default:
			panic(fmt.Sprintf("unknown opcode %x", op))
		}

		registers[0] += 3
	}
}
