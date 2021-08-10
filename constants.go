package main

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

	AND_OPCODE_IMM  = 0x29
	AND_OPCODE_ZP   = 0x25
	AND_OPCODE_ZPX  = 0x35
	AND_OPCODE_ABS  = 0x2d
	AND_OPCODE_ABSX = 0x3d
	AND_OPCODE_ABSY = 0x39
	AND_OPCODE_INDX = 0x21
	AND_OPCODE_INDY = 0x31

	ASL_OPCODE_ACC  = 0xa
	ASL_OPCODE_ZP   = 0x06
	ASL_OPCODE_ZPX  = 0x016
	ASL_OPCODE_ABS  = 0x0e
	ASL_OPCODE_ABSX = 0x1e

	CLC_OPCODE = 0x18
	SEC_OPCODE = 0x38
	SED_OPCODE = 0xf8
	SEI_OPCODE = 0x78

	BCC_OPCODE = 0x90
	BCS_OPCODE = 0xB0
	BEQ_OPCODE = 0xF0
	BMI_OPCODE = 0x30
	BNE_OPCODE = 0xD0
	BPL_OPCODE = 0x10
	BVC_OPCODE = 0x50
	BVS_OPCODE = 0x70

	BIT_OPCODE_ZP  = 0x24
	BIT_OPCODE_ABS = 0x2c

	BRK_OPCODE = 0x00

	CLD_OPCODE = 0xd8
	CLI_OPCODE = 0x58
	CLV_OPCODE = 0xb8
	NOP_OPCODE = 0xea

	PHA_OPCODE = 0x48
	PLA_OPCODE = 0x68
	RTS_OPCODE = 0x60
	TAX_OPCODE = 0xaa
	TXA_OPCODE = 0x8a
	TAY_OPCODE = 0xa8
	TYA_OPCODE = 0x98
	TSX_OPCODE = 0xba
	TXS_OPCODE = 0x9a

	CMP_OPCODE_IMM  = 0xc9
	CMP_OPCODE_ZP   = 0xc5
	CMP_OPCODE_ZPX  = 0xd5
	CMP_OPCODE_ABS  = 0xcd
	CMP_OPCODE_ABSX = 0xdd
	CMP_OPCODE_ABSY = 0xd9
	CMP_OPCODE_INDX = 0xc1
	CMP_OPCODE_INDY = 0xd1
)

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

type memoryRegion struct {
	start uint16
	end   uint16
}

var memoryMap = map[string]memoryRegion{
	"STACK": {0x100, 0x1FF},
}
