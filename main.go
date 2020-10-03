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
	Register_A            int8
	REGISTER_X            int8
	REGISTER_Y            int8
	REGISTER_PC           uint16
	REGISTER_STATUS_P     uint8
	REGISTER_STACKPOINTER uint8
	Instructions          map[OPCODE]InstructionData
	Memory                []int8
}

func (sim *Simulator) executeInstruction(instr InstructionData) {
	//lookup opcode execution.
	InstructionFunctionMap[OPCODE(instr.opcode)](sim)
}
func (sim *Simulator) Run() {
	//fetch
	//get instruction at program counter
	currentOP := sim.Memory[sim.REGISTER_PC]

	//decode
	instruction := sim.Instructions[OPCODE(currentOP)]
	//execute
	sim.executeInstruction(instruction)
}

type REGISTER int

const (
	REGISTER_A            REGISTER = 0
	REGISTER_X            REGISTER = 1
	REGISTER_Y            REGISTER = 2
	REGISTER_PC           REGISTER = 3
	REGISTER_STATUS_P     REGISTER = 4
	REGISTER_STACKPOINTER REGISTER = 5
)

type OPCODE int

var InstructionFunctionMap = map[OPCODE]func(sim *Simulator){
	OPCODE(105): func(sim *Simulator) {
		sim.Register_A = sim.Register_A + 1
	},
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
	ZEROPAGE_INDEXEDY  ADDRESS_MODE = 11
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
	ZEROPAGE_INDEXEDY:  "ZPY",
	INDEXED_INDIRECT_X: "INDX",
	INDIRECT_INDEXED_Y: "INDY",
}

type InstructionData struct {
	opcode      OPCODE
	memonic     string
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

func generateInstructionMap(filePath string) map[OPCODE]InstructionData {
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
		fmt.Println(record)
		var op, _ = strconv.ParseUint(record[0], 0, 8)
		var mem = record[1]
		var addmode = IMMEDIATE

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
		fmt.Println(currentInstructionData)

		result[currentInstructionData.opcode] = currentInstructionData
	}
	return result
}

func main() {
	fmt.Println("generate simulator from csv")
	var filePath string = "6502ops.csv"
	generateInstructionMap(filePath)

	fmt.Println("instantiate simulator")

}
