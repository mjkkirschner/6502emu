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

const (
	OPCODE_ADD OPCODE = 0
	OPCPDE_MOV OPCODE = 1
)

type ADDRESS_MODE uint8

//TOOD USE same names as in
const (
	IMMEDIATE ADDRESS_MODE = 0
	ABSOLUTE  ADDRESS_MODE = 1
	ZEROPAGE  ADDRESS_MODE = 2
)

var ADDRESS_MODE_NAME_MAP = map[ADDRESS_MODE]string{
	IMMEDIATE: "IMM",
	ABSOLUTE:  "ABS",
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

func main() {
	fmt.Printf("cast int8: %d\n", int8(a))
	var filePath string = "6502ops.csv"
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
	for _, record := range records {
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

	}

	//
	//fmt.Println(len(record))
	//for value := range record {
	//	fmt.Printf("  %v\n", record[value])
	//}
}
