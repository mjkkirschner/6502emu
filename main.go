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
	X_JUMPING             bool
}

func NewSimulator(instructions map[OPCODE]InstructionData) *Simulator {
	return &Simulator{Instructions: instructions, Memory: make([]uint8, 65536)}
}

func (sim *Simulator) reset() {
	//read from fffc and fffd
	//then transfer control.
	addrlow := sim.Memory[0xFFFC]
	addrhigh := sim.Memory[0xFFFD]
	longaddr := uint16(addrhigh)<<8 | (uint16(addrlow) & 0xff)
	sim.REGISTER_PC = longaddr
}

func (sim *Simulator) executeInstruction(instr InstructionData) {

	//decode operands based on address mode type.
	operands := sim.decodeOperands(instr)
	//lookup opcode execution.
	opFunc := InstructionFunctionMap[OPCODE(instr.opcode)]
	if opFunc == nil {
		log.Fatal("no implementation for ", instr.mnemonic, " ", ADDRESS_MODE_NAME_MAP[instr.addressMode])
	}

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

func (sim *Simulator) GetBit(reg REGISTER, bit uint) bool {

	switch reg {
	case REGISTER_A:
		return sim.Register_A&(1<<bit) > 0

	case REGISTER_STACKPOINTER:
		return sim.REGISTER_STACKPOINTER&(1<<bit) > 0

	case REGISTER_X:
		return sim.REGISTER_X&(1<<bit) > 0

	case REGISTER_Y:
		return sim.REGISTER_Y&(1<<bit) > 0

	case REGISTER_STATUS:
		return sim.REGISTER_STATUS_P&(1<<bit) > 0
	}
	log.Fatal("unhandled register in get bit")
	return false
}

func GetBit(value uint, bit uint) bool {
	return value&(1<<bit) > 0
}

func (sim *Simulator) incrementPC(inc uint8) {
	sim.REGISTER_PC = sim.REGISTER_PC + uint16(inc)
}

type decodeResults struct {
	operands []interface{}
	//for some operations which do not have implicit return locations
	//or registers, this address stores the return address of the operation
	//where computed results go.
	//for example - ASL can shift values
	returnAddress uint16
}

//get operands based on address type
func (sim *Simulator) decodeOperands(instr InstructionData) decodeResults {

	//since these are memory locations negatives usually don't make sense.
	var a uint8
	var b uint8
	var longaddr uint16
	outputOperands := make([]interface{}, 0)
	//return address only makes sense in the context of some instructions...
	//for now we'll just set it to the final address before we get the value.
	var returnAddress uint16 = 0

	switch instr.addressMode {
	//load 8bit constants into memory.
	//not sure there will ever be a valid b operand.
	case IMMEDIATE:
		a = uint8(sim.Memory[sim.REGISTER_PC+1])
		b = uint8(sim.Memory[sim.REGISTER_PC+2])
		outputOperands = append(outputOperands, a, b)

	case ZEROPAGE:
		x0 := sim.Memory[sim.REGISTER_PC+1]
		a = uint8(sim.Memory[x0])
		outputOperands = append(outputOperands, a)
		returnAddress = uint16(x0)

	case ZEROPAGE_INDEXEDX:
		x0 := sim.Memory[sim.REGISTER_PC+1] + sim.REGISTER_X
		a = uint8(sim.Memory[x0])
		outputOperands = append(outputOperands, a+b)
		returnAddress = uint16(x0)

	//address at absolute 16bit address
	case ABSOLUTE:
		//a will be LSB, b will be MSB since 6502 is little endian
		a = uint8(sim.Memory[sim.REGISTER_PC+1])
		b = uint8(sim.Memory[sim.REGISTER_PC+2])
		//shift msb up 8 then or with a (and with 255 clears any upper  bits...)
		longaddr = uint16(b)<<8 | (uint16(a) & 0xff)
		output := sim.Memory[longaddr]
		outputOperands = append(outputOperands, output)
		returnAddress = longaddr

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
		returnAddress = longaddr

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
		//in the case of relative its a signed byte max of 127, min -127 from pc.
		//convert to signed.
		offset := int8(sim.Memory[sim.REGISTER_PC+1])
		jumpAddr := uint16(offset) + uint16(sim.REGISTER_PC)
		outputOperands = append(outputOperands, jumpAddr)

	case ACCUMULATOR:
		outputOperands = append(outputOperands, sim.Register_A)
	case IMPLIED:

		//TODO some instructions like branch intructions will need to reinterpert the results
		//as signed offset numbers.
	}
	return decodeResults{outputOperands, returnAddress}
}
func (sim *Simulator) SingleStep() {
	//fetch
	//get instruction at program counter
	currentOP := sim.Memory[sim.REGISTER_PC]

	//decode
	instruction := sim.Instructions[OPCODE(currentOP)]
	//reset jump flag.
	sim.X_JUMPING = false
	//execute
	sim.executeInstruction(instruction)
	//if the instruciton set the PC directly, then don't increment it
	if !sim.X_JUMPING {
		sim.incrementPC(instruction.bytes)
	}
}

func (sim *Simulator) Run(instructions uint) {
	for i := 0; i < int(instructions); i++ {
		sim.SingleStep()
	}
}

var InstructionFunctionMap = map[OPCODE]func(sim *Simulator, operands decodeResults, instruction InstructionData){
	ADDWITHCARRY_OPCODE_IMM:  INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ZP:   INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ZPX:  INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ABS:  INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ABSX: INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_ABSY: INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_INDX: INSTRUCTION_ADC_IMPLEMENTATION,
	ADDWITHCARRY_OPCODE_INDY: INSTRUCTION_ADC_IMPLEMENTATION,

	AND_OPCODE_IMM:  INSTRUCTION_AND_IMPLEMENTATION,
	AND_OPCODE_ZP:   INSTRUCTION_AND_IMPLEMENTATION,
	AND_OPCODE_ZPX:  INSTRUCTION_AND_IMPLEMENTATION,
	AND_OPCODE_ABS:  INSTRUCTION_AND_IMPLEMENTATION,
	AND_OPCODE_ABSX: INSTRUCTION_AND_IMPLEMENTATION,
	AND_OPCODE_ABSY: INSTRUCTION_AND_IMPLEMENTATION,
	AND_OPCODE_INDX: INSTRUCTION_AND_IMPLEMENTATION,
	AND_OPCODE_INDY: INSTRUCTION_AND_IMPLEMENTATION,

	CLC_OPCODE: INSTRUCTION_CLC_IMPLEMENTATION,
	SEC_OPCODE: INSTRUCTION_SEC_IMPLEMENTATION,
	SED_OPCODE: INSTRUCTION_SED_IMPLEMENTATION,
	SEI_OPCODE: INSTRUCTION_SEI_IMPLEMENTATION,

	ASL_OPCODE_ABS:  INSTRUCTION_ASL_IMPLEMENTATION,
	ASL_OPCODE_ABSX: INSTRUCTION_ASL_IMPLEMENTATION,
	ASL_OPCODE_ZP:   INSTRUCTION_ASL_IMPLEMENTATION,
	ASL_OPCODE_ZPX:  INSTRUCTION_ASL_IMPLEMENTATION,
	ASL_OPCODE_ACC:  INSTRUCTION_ASL_IMPLEMENTATION,

	BCC_OPCODE: INSTRUCTION_BCC_IMPLEMENTATION,
	BCS_OPCODE: INSTRUCTION_BCS_IMPLEMENTATION,
	BEQ_OPCODE: INSTRUCTION_BEQ_IMPLEMENTATION,
	BMI_OPCODE: INSTRUCTION_BMI_IMPLEMENTATION,
	BNE_OPCODE: INSTRUCTION_BNE_IMPLEMENTATION,
	BPL_OPCODE: INSTRUCTION_BPL_IMPLEMENTATION,
	BVC_OPCODE: INSTRUCTION_BVC_IMPLEMENTATION,
	BVS_OPCODE: INSTRUCTION_BVS_IMPLEMENTATION,

	BIT_OPCODE_ZP:  INSTRUCTION_BIT_IMPLEMENTATION,
	BIT_OPCODE_ABS: INSTRUCTION_BIT_IMPLEMENTATION,

	BRK_OPCODE: INSTRUCTION_BRK_IMPLEMENTATION,
	CLD_OPCODE: INSTRUCTION_CLD_IMPLEMENTATION,
	CLI_OPCODE: INSTRUCTION_CLI_IMPLEMENTATION,
	CLV_OPCODE: INSTRUCTION_CLV_IMPLEMENTATION,
	NOP_OPCODE: INSTRUCTION_NOP_IMPLEMENTATION,

	PHA_OPCODE: INSTRUCTION_PHA_IMPLEMENTATION,
	PLA_OPCODE: INSTRUCTION_PLA_IMPLEMENTATION,
	RTS_OPCODE: INSTRUCTION_RTS_IMPLEMENTATION,
	RTI_OPCODE: INSTRUCTION_RTI_IMPLEMENTATION,
	TAX_OPCODE: INSTRUCTION_TAX_IMPLEMENTATION,
	TXA_OPCODE: INSTRUCTION_TXA_IMPLEMENTATION,
	TAY_OPCODE: INSTRUCTION_TAY_IMPLEMENTATION,
	TYA_OPCODE: INSTRUCTION_TYA_IMPLEMENTATION,
	TSX_OPCODE: INSTRUCTION_TSX_IMPLEMENTATION,
	TXS_OPCODE: INSTRUCTION_TXS_IMPLEMENTATION,
	PHP_OPCODE: INSTRUCTION_PHP_IMPLEMENTATION,
	PLP_OPCODE: INSTRUCTION_PLP_IMPLEMENTATION,

	CMP_OPCODE_IMM:  INSTRUCTION_CMP_IMPLEMENTATION,
	CMP_OPCODE_ZP:   INSTRUCTION_CMP_IMPLEMENTATION,
	CMP_OPCODE_ZPX:  INSTRUCTION_CMP_IMPLEMENTATION,
	CMP_OPCODE_ABS:  INSTRUCTION_CMP_IMPLEMENTATION,
	CMP_OPCODE_ABSX: INSTRUCTION_CMP_IMPLEMENTATION,
	CMP_OPCODE_ABSY: INSTRUCTION_CMP_IMPLEMENTATION,
	CMP_OPCODE_INDX: INSTRUCTION_CMP_IMPLEMENTATION,
	CMP_OPCODE_INDY: INSTRUCTION_CMP_IMPLEMENTATION,

	CPX_OPCODE_IMM: INSTRUCTION_CPX_IMPLEMENTATION,
	CPX_OPCODE_ZP:  INSTRUCTION_CPX_IMPLEMENTATION,
	CPX_OPCODE_ABS: INSTRUCTION_CPX_IMPLEMENTATION,

	CPY_OPCODE_IMM: INSTRUCTION_CPY_IMPLEMENTATION,
	CPY_OPCODE_ZP:  INSTRUCTION_CPY_IMPLEMENTATION,
	CPY_OPCODE_ABS: INSTRUCTION_CPY_IMPLEMENTATION,

	DEC_OPCODE_ZP:   INSTRUCTION_DEC_IMPLEMENTATION,
	DEC_OPCODE_ZPX:  INSTRUCTION_DEC_IMPLEMENTATION,
	DEC_OPCODE_ABS:  INSTRUCTION_DEC_IMPLEMENTATION,
	DEC_OPCODE_ABSX: INSTRUCTION_DEC_IMPLEMENTATION,

	DEX_OPCODE_IMM: INSTRUCTION_DEX_IMPLEMENTATION,
	DEY_OPCODE_IMM: INSTRUCTION_DEY_IMPLEMENTATION,
	INX_OPCODE_IMM: INSTRUCTION_INX_IMPLEMENTATION,
	INY_OPCODE_IMM: INSTRUCTION_INY_IMPLEMENTATION,

	EOR_OPCODE_IMM:  INSTRUCTION_EOR_IMPLEMENTATION,
	EOR_OPCODE_ZP:   INSTRUCTION_EOR_IMPLEMENTATION,
	EOR_OPCODE_ZPX:  INSTRUCTION_EOR_IMPLEMENTATION,
	EOR_OPCODE_ABS:  INSTRUCTION_EOR_IMPLEMENTATION,
	EOR_OPCODE_ABSX: INSTRUCTION_EOR_IMPLEMENTATION,
	EOR_OPCODE_ABSY: INSTRUCTION_EOR_IMPLEMENTATION,
	EOR_OPCODE_INDX: INSTRUCTION_EOR_IMPLEMENTATION,
	EOR_OPCODE_INDY: INSTRUCTION_EOR_IMPLEMENTATION,

	INC_OPCODE_ZP:   INSTRUCTION_INC_IMPLEMENTATION,
	INC_OPCODE_ZPX:  INSTRUCTION_INC_IMPLEMENTATION,
	INC_OPCODE_ABS:  INSTRUCTION_INC_IMPLEMENTATION,
	INC_OPCODE_ABSX: INSTRUCTION_INC_IMPLEMENTATION,

	JMP_OPCODE_ABS: INSTRUCTION_JMP_IMPLEMENTATION,
	JMP_OPCODE_IND: INSTRUCTION_JMP_IMPLEMENTATION,

	JSR_OPCODE_ABS: INSTRUCTION_JSR_IMPLEMENTATION,

	LDA_OPCODE_IMM:  INSTRUCTION_LDA_IMPLEMENTATION,
	LDA_OPCODE_ZP:   INSTRUCTION_LDA_IMPLEMENTATION,
	LDA_OPCODE_ZPX:  INSTRUCTION_LDA_IMPLEMENTATION,
	LDA_OPCODE_ABS:  INSTRUCTION_LDA_IMPLEMENTATION,
	LDA_OPCODE_ABSX: INSTRUCTION_LDA_IMPLEMENTATION,
	LDA_OPCODE_ABSY: INSTRUCTION_LDA_IMPLEMENTATION,
	LDA_OPCODE_INDX: INSTRUCTION_LDA_IMPLEMENTATION,
	LDA_OPCODE_INDY: INSTRUCTION_LDA_IMPLEMENTATION,

	LDX_OPCODE_IMM:  INSTRUCTION_LDX_IMPLEMENTATION,
	LDX_OPCODE_ZP:   INSTRUCTION_LDX_IMPLEMENTATION,
	LDX_OPCODE_ZPY:  INSTRUCTION_LDX_IMPLEMENTATION,
	LDX_OPCODE_ABS:  INSTRUCTION_LDX_IMPLEMENTATION,
	LDX_OPCODE_ABSY: INSTRUCTION_LDX_IMPLEMENTATION,

	LDY_OPCODE_IMM:  INSTRUCTION_LDY_IMPLEMENTATION,
	LDY_OPCODE_ZP:   INSTRUCTION_LDY_IMPLEMENTATION,
	LDY_OPCODE_ZPX:  INSTRUCTION_LDY_IMPLEMENTATION,
	LDY_OPCODE_ABS:  INSTRUCTION_LDY_IMPLEMENTATION,
	LDY_OPCODE_ABSX: INSTRUCTION_LDY_IMPLEMENTATION,

	LSR_OPCODE_ACC:  INSTRUCTION_LSR_IMPLEMENTATION,
	LSR_OPCODE_ZP:   INSTRUCTION_LSR_IMPLEMENTATION,
	LSR_OPCODE_ZPX:  INSTRUCTION_LSR_IMPLEMENTATION,
	LSR_OPCODE_ABS:  INSTRUCTION_LSR_IMPLEMENTATION,
	LSR_OPCODE_ABSX: INSTRUCTION_LSR_IMPLEMENTATION,

	ORA_OPCODE_IMM:  INSTRUCTION_ORA_IMPLEMENTATION,
	ORA_OPCODE_ZP:   INSTRUCTION_ORA_IMPLEMENTATION,
	ORA_OPCODE_ZPX:  INSTRUCTION_ORA_IMPLEMENTATION,
	ORA_OPCODE_ABS:  INSTRUCTION_ORA_IMPLEMENTATION,
	ORA_OPCODE_ABSX: INSTRUCTION_ORA_IMPLEMENTATION,
	ORA_OPCODE_ABSY: INSTRUCTION_ORA_IMPLEMENTATION,
	ORA_OPCODE_INDX: INSTRUCTION_ORA_IMPLEMENTATION,
	ORA_OPCODE_INDY: INSTRUCTION_ORA_IMPLEMENTATION,

	ROR_OPCODE_ACC:  INSTRUCTION_ROR_IMPLEMENTATION,
	ROR_OPCODE_ZP:   INSTRUCTION_ROR_IMPLEMENTATION,
	ROR_OPCODE_ZPX:  INSTRUCTION_ROR_IMPLEMENTATION,
	ROR_OPCODE_ABS:  INSTRUCTION_ROR_IMPLEMENTATION,
	ROR_OPCODE_ABSX: INSTRUCTION_ROR_IMPLEMENTATION,

	ROL_OPCODE_ACC:  INSTRUCTION_ROL_IMPLEMENTATION,
	ROL_OPCODE_ZP:   INSTRUCTION_ROL_IMPLEMENTATION,
	ROL_OPCODE_ZPX:  INSTRUCTION_ROL_IMPLEMENTATION,
	ROL_OPCODE_ABS:  INSTRUCTION_ROL_IMPLEMENTATION,
	ROL_OPCODE_ABSX: INSTRUCTION_ROL_IMPLEMENTATION,

	SBC_OPCODE_IMM:  INSTRUCTION_SBC_IMPLEMENTATION,
	SBC_OPCODE_ZP:   INSTRUCTION_SBC_IMPLEMENTATION,
	SBC_OPCODE_ZPX:  INSTRUCTION_SBC_IMPLEMENTATION,
	SBC_OPCODE_ABS:  INSTRUCTION_SBC_IMPLEMENTATION,
	SBC_OPCODE_ABSX: INSTRUCTION_SBC_IMPLEMENTATION,
	SBC_OPCODE_ABSY: INSTRUCTION_SBC_IMPLEMENTATION,
	SBC_OPCODE_INDX: INSTRUCTION_SBC_IMPLEMENTATION,
	SBC_OPCODE_INDY: INSTRUCTION_SBC_IMPLEMENTATION,

	STA_OPCODE_ZP:   INSTRUCTION_STA_IMPLEMENTATION,
	STA_OPCODE_ZPX:  INSTRUCTION_STA_IMPLEMENTATION,
	STA_OPCODE_ABS:  INSTRUCTION_STA_IMPLEMENTATION,
	STA_OPCODE_ABSX: INSTRUCTION_STA_IMPLEMENTATION,
	STA_OPCODE_ABSY: INSTRUCTION_STA_IMPLEMENTATION,
	STA_OPCODE_INDX: INSTRUCTION_STA_IMPLEMENTATION,
	STA_OPCODE_INDY: INSTRUCTION_STA_IMPLEMENTATION,

	STX_OPCODE_ZP:  INSTRUCTION_STX_IMPLEMENTATION,
	STX_OPCODE_ZPY: INSTRUCTION_STX_IMPLEMENTATION,
	STX_OPCODE_ABS: INSTRUCTION_STX_IMPLEMENTATION,

	STY_OPCODE_ZP:  INSTRUCTION_STY_IMPLEMENTATION,
	STY_OPCODE_ZPX: INSTRUCTION_STY_IMPLEMENTATION,
	STY_OPCODE_ABS: INSTRUCTION_STY_IMPLEMENTATION,
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
