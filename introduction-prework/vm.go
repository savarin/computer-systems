package vm

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
	var register, address, value byte

	// Keep looping, like a physical computer's clock
loop:
	for {
		pc := registers[0]
		op := memory[pc] // fetch the opcode

		// decode and execute
		switch op {
		case Load:
			register = memory[pc+1]
			address = memory[pc+2]

			value = memory[address]
			registers[register] = value

			registers[0] += 3

		case Store:
			register = memory[pc+1]
			address = memory[pc+2]

			value = registers[register]
			memory[address] = value

			registers[0] += 3

		case Halt:
			break loop
		}
	}
}
