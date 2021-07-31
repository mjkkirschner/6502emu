package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

type Simulator struct {
	Register_A            uint8
	REGISTER_X            uint8
	REGISTER_Y            uint8
	REGISTER_PC           uint16
	REGISTER_STATUS_P     uint8
	REGISTER_STACKPOINTER uint8
	Instructions          map[OPCODE]InstructionData
	Memory                []uint8
}

func NewSimulator(instructions map[OPCODE]InstructionData) *Simulator {
	return &Simulator{Instructions: instructions, Memory: make([]uint8, 65536)}
}

func (sim *Simulator) executeInstruction(instr InstructionData) {
	//decode operands based on address mode type.
	operands := sim.decodeOperands(instr)
	//lookup opcode execution.
	opFunc := InstructionFunctionMap[OPCODE(instr.opcode)]

	//execute
	opFunc(sim, operands, instr)

}

//we leave out PC register deliberately since its 16 bits.
//also unlikely we'll ever bit set PC.
func (sim *Simulator) SetBit(reg REGISTER, bit uint) {

	switch reg {
	case REGISTER_A:
		sim.Register_A |= (1 << bit)

	case REGISTER_STACKPOINTER:
		sim.REGISTER_STACKPOINTER |= (1 << bit)

	case REGISTER_X:
		sim.REGISTER_X |= (1 << bit)

	case REGISTER_Y:
		sim.REGISTER_Y |= (1 << bit)

	case REGISTER_STATUS:
		sim.REGISTER_STATUS_P |= (1 << bit)
	}
}
func (sim *Simulator) ClearBit(reg REGISTER, bit uint) {

	switch reg {
	case REGISTER_A:
		sim.Register_A &= ^(1 << bit)

	case REGISTER_STACKPOINTER:
		sim.REGISTER_STACKPOINTER &= ^(1 << bit)

	case REGISTER_X:
		sim.REGISTER_X &= ^(1 << bit)

	case REGISTER_Y:
		sim.REGISTER_Y &= ^(1 << bit)

	case REGISTER_STATUS:
		sim.REGISTER_STATUS_P &= ^(1 << bit)
	}
}

func (sim *Simulator) GetBit(reg REGISTER, bit uint) uint8 {

	switch reg {
	case REGISTER_A:
		return sim.Register_A & (1 << bit)

	case REGISTER_STACKPOINTER:
		return sim.REGISTER_STACKPOINTER & (1 << bit)

	case REGISTER_X:
		return sim.REGISTER_X & (1 << bit)

	case REGISTER_Y:
		return sim.REGISTER_Y & (1 << bit)

	case REGISTER_STATUS:
		return sim.REGISTER_STATUS_P & (1 << bit)
	}
	log.Fatal("unhandled register in get bit")
	return 255
}

func (sim *Simulator) incrementPC(inc uint8) {
	sim.REGISTER_PC = sim.REGISTER_PC + uint16(inc)
}

//get operands based on address type
func (sim *Simulator) decodeOperands(instr InstructionData) []interface{} {

	//since these are memory locations negatives usually don't make sense.
	var a uint8
	var b uint8
	var longaddr uint16
	outputOperands := make([]interface{}, 0)

	switch instr.addressMode {
	//load 8bit constants into memory.
	//not sure there will ever be a valid b operand.
	case IMMEDIATE:
		a = uint8(sim.Memory[sim.REGISTER_PC+1])
		b = uint8(sim.Memory[sim.REGISTER_PC+2])
		outputOperands = append(outputOperands, a, b)

	case ZEROPAGE:
		a = uint8(sim.Memory[sim.Memory[sim.REGISTER_PC+1]])
		outputOperands = append(outputOperands, a)

	case ZEROPAGE_INDEXEDX:
		a = uint8(sim.Memory[sim.Memory[sim.REGISTER_PC+1]+sim.REGISTER_X])
		outputOperands = append(outputOperands, a+b)

	//address at absolute 16bit address
	case ABSOLUTE:
		//a will be LSB, b will be MSB since 6502 is little endian
		a = uint8(sim.Memory[sim.REGISTER_PC+1])
		b = uint8(sim.Memory[sim.REGISTER_PC+2])
		//shift msb up 8 then or with a (and with 255 clears any upper  bits...)
		longaddr = uint16(b)<<8 | (uint16(a) & 0xff)
		output := sim.Memory[longaddr]
		outputOperands = append(outputOperands, output)

	case ABSOLUTE_INDEXEDX:

		//a will be LSB, b will be MSB since 6502 is little endian
		a = uint8(sim.Memory[sim.REGISTER_PC+1])
		b = uint8(sim.Memory[sim.REGISTER_PC+2])
		//shift msb up 8 then or with a (and with 255 clears any upper  bits...)
		longaddr = uint16(b)<<8 | (uint16(a) & 0xff)
		b = sim.REGISTER_X
		longaddr = longaddr + uint16(b)
		output := sim.Memory[longaddr]
		outputOperands = append(outputOperands, output)

	case ABSOLUTE_INDEXEDY:
		//a will be LSB, b will be MSB since 6502 is little endian
		a = uint8(sim.Memory[sim.REGISTER_PC+1])
		b = uint8(sim.Memory[sim.REGISTER_PC+2])
		//shift msb up 8 then or with a (and with 255 clears any upper  bits...)
		longaddr = uint16(b)<<8 | (uint16(a) & 0xff)
		b = sim.REGISTER_Y
		longaddr = longaddr + uint16(b)
		output := sim.Memory[longaddr]
		outputOperands = append(outputOperands, output)

	case INDEXED_INDIRECT_X:
		//zp address
		a = uint8(sim.Memory[sim.REGISTER_PC+1])
		//then offset by x
		addr := a + sim.REGISTER_X

		//get address at a+x
		lowbyte := sim.Memory[addr]
		highByte := sim.Memory[addr+1]
		//now combine bytes highLow and return that as the final address for the jump
		longaddr = uint16(highByte)<<8 | (uint16(lowbyte) & 0xff)
		//now indirect
		output := sim.Memory[longaddr]
		outputOperands = append(outputOperands, output)

	case INDIRECT_INDEXED_Y:
		//zp address indirect
		a = uint8(sim.Memory[sim.REGISTER_PC+1])

		//get address at a+x
		lowbyte := sim.Memory[a]
		highByte := sim.Memory[a+1]
		//now combine bytes highLow and return that as the final address for the jump
		longaddr = uint16(highByte)<<8 | (uint16(lowbyte) & 0xff) + uint16(sim.REGISTER_Y)
		//now indirect
		output := sim.Memory[longaddr]
		outputOperands = append(outputOperands, output)

		//only JMP will use INDIRECT this address mode.
	case INDIRECT:
		//a will be LSB, b will be MSB since 6502 is little endian
		a = uint8(sim.Memory[sim.Memory[sim.REGISTER_PC+1]])
		b = uint8(sim.Memory[sim.Memory[sim.REGISTER_PC+2]])
		//shift msb up 8 then or with a (and with 255 clears any upper  bits...)
		longaddr = uint16(b)<<8 | (uint16(a) & 0xff)

		//now we indirect.
		lowbyte := sim.Memory[longaddr]
		highByte := sim.Memory[longaddr+1]
		//now combine bytes highLow and return that as the final address for the jump
		longaddr = uint16(highByte)<<8 | (uint16(lowbyte) & 0xff)
		outputOperands = append(outputOperands, longaddr)

	case RELATIVE:
	case ACCUMULATOR:
	case IMPLIED:

		//TODO some instructions like branch intructions will need to reinterpert the results
		//as signed offset numbers.
	}
	return outputOperands
}
func (sim *Simulator) SingleStep() {
	//fetch
	//get instruction at program counter
	currentOP := sim.Memory[sim.REGISTER_PC]

	//decode
	instruction := sim.Instructions[OPCODE(currentOP)]
	//execute
	sim.executeInstruction(instruction)

	sim.incrementPC(instruction.bytes)
}

func (sim *Simulator) Run(instructions uint) {
	for i := 0; i < int(instructions); i++ {
		sim.SingleStep()
	}
}

type REGISTER int

const (
	REGISTER_A            REGISTER = 100
	REGISTER_X            REGISTER = 200
	REGISTER_Y            REGISTER = 300
	REGISTER_PC           REGISTER = 400
	REGISTER_STACKPOINTER REGISTER = 500
	REGISTER_STATUS       REGISTER = 600

	BITFLAG_STATUS_CARRY             = 0
	BITFLAG_STATUS_ZERO              = 1
	BITFLAG_STATUS_INTERRUPT_DISABLE = 2
	BITFLAG_STATUS_DECIMAL           = 3
	BITFLAG_STATUS_B_FLAG            = 4
	BITFLAG_STATUS_OVERFLOW          = 5
	BITFLAG_STATUS_NEGATIVE          = 6
)

type OPCODE int

//todo consider using memonic or different name for each opcode with addressing...
//or to try to centralize decode logic of operands. - see trial in decodeOperands() function
const (
	ADDWITHCARRY_OPCODE_IMM  = 105
	ADDWITHCARRY_OPCODE_ZP   = 101
	ADDWITHCARRY_OPCODE_ZPX  = 0x75
	ADDWITHCARRY_OPCODE_ABS  = 0x6d
	ADDWITHCARRY_OPCODE_ABSX = 0x7d
	ADDWITHCARRY_OPCODE_ABSY = 0x79
	ADDWITHCARRY_OPCODE_INDX = 0x61
	ADDWITHCARRY_OPCODE_INDY = 0x71
)

func INSTRUCTION_ADC_IMPLEMENTATION(sim *Simulator, operands []interface{}, instruction InstructionData) {
	//calculate the result.
	a := sim.Register_A
	b := (operands[0]).(uint8)
	c := sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
	sum := sim.Register_A + b + c

	carryCheck := uint16(a) + uint16(b) + uint16(c)
	overFlowCheck := (a ^ sum) & (b ^ sum) & 0x80 //negative bit.

	sim.Register_A = sim.Register_A + b + c
	//if the addition resulted in an overflow carry should be set to 1 - if not carry should be reset to 0.
	if carryCheck > 255 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
	}
	//overflow occurs when signed arithmetic overflows.
	if overFlowCheck == 1 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW)
	}
	//zero flag if result is 0
	if sim.Register_A == 0 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	}
	//set n
	nbit := sim.GetBit(REGISTER_A, 7)
	if nbit == 1 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	}
}

var InstructionFunctionMap = map[OPCODE]func(sim *Simulator, operands []interface{}, instruction InstructionData){
	//TODO the code below should be the same for all ADC commands regardless of address mode I think - share it.
	ADDWITHCARRY_OPCODE_IMM:  INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ZP:   INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ZPX:  INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ABS:  INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ABSX: INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ABSY: INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_INDX: INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_INDY: INSTRUCTION_ADC_IMPLEMENTATION,
}

type ADDRESS_MODE uint8

const (
	IMMEDIATE   ADDRESS_MODE = 0
	ABSOLUTE    ADDRESS_MODE = 1
	ZEROPAGE    ADDRESS_MODE = 2
	ACCUMULATOR ADDRESS_MODE = 3
	IMPLIED     ADDRESS_MODE = 4
	RELATIVE    ADDRESS_MODE = 5
	INDIRECT    ADDRESS_MODE = 6

	ABSOLUTE_INDEXEDX  ADDRESS_MODE = 8
	ABSOLUTE_INDEXEDY  ADDRESS_MODE = 9
	ZEROPAGE_INDEXEDX  ADDRESS_MODE = 10
	INDEXED_INDIRECT_X ADDRESS_MODE = 12
	INDIRECT_INDEXED_Y ADDRESS_MODE = 13
)

var ADDRESS_MODE_NAME_MAP = map[ADDRESS_MODE]string{
	IMMEDIATE:          "IMM",
	ABSOLUTE:           "ABS",
	ZEROPAGE:           "ZP",
	ACCUMULATOR:        "ACC",
	IMPLIED:            "IMP",
	RELATIVE:           "REL",
	INDIRECT:           "IND",
	ABSOLUTE_INDEXEDX:  "ABSX",
	ABSOLUTE_INDEXEDY:  "ABSY",
	ZEROPAGE_INDEXEDX:  "ZPX",
	INDEXED_INDIRECT_X: "INDX",
	INDIRECT_INDEXED_Y: "INDY",
}

type InstructionData struct {
	opcode      OPCODE
	mnemonic    string
	addressMode ADDRESS_MODE
	bytes       uint8
	cycles      uint8
	flags       flagsEffected
}

type flagsEffected struct {
	carry           bool
	zero            bool
	interuptDisable bool
	decimal         bool
	bflag           bool
	overflowV       bool
	negative        bool
}

func newFlagsEffected(str string) *flagsEffected {
	var f flagsEffected
	if unicode.IsUpper(rune(str[0])) {
		f.carry = true
	}
	if unicode.IsUpper(rune(str[1])) {
		f.zero = true
	}
	if unicode.IsUpper(rune(str[2])) {
		f.interuptDisable = true
	}
	if unicode.IsUpper(rune(str[3])) {
		f.decimal = true
	}
	if unicode.IsUpper(rune(str[4])) {
		f.bflag = true
	}
	if unicode.IsUpper(rune(str[5])) {
		f.overflowV = true
	}
	if unicode.IsUpper(rune(str[6])) {
		f.negative = true
	}
	return &f
}

//casting unsigned int 256 to signed int should yield -127.
var a uint = 255

func GenerateInstructionMap(filePath string) map[OPCODE]InstructionData {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	//close file when main exits... useless here.
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	result := map[OPCODE]InstructionData{}

	for _, record := range records[1:] {
		//fmt.Println(record)
		var op, _ = strconv.ParseUint(record[0], 0, 8)
		var mem = record[1]
		var addmode = IMMEDIATE

		//if its a valid address mode use it.
		for key, element := range ADDRESS_MODE_NAME_MAP {
			if record[2] == element {
				addmode = key
				break
			}
		}

		var bytes, _ = strconv.ParseUint(record[3], 0, 8)
		var cycles, _ = strconv.ParseUint(record[4], 0, 8)
		var flags = newFlagsEffected(record[5])

		currentInstructionData := InstructionData{OPCODE(op), mem, addmode, uint8(bytes), uint8(cycles), *flags}
		//fmt.Println(currentInstructionData)

		result[currentInstructionData.opcode] = currentInstructionData
	}
	return result
}

func main() {
	fmt.Println("generate simulator from csv")
	var filePath string = "6502ops.csv"
	instructions := GenerateInstructionMap(filePath)

	fmt.Println("instantiate simulator")
	simulator := NewSimulator(instructions)
	fmt.Println(simulator.Instructions)

}
