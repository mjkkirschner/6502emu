package main

import (
	"testing"
)

func TestANDImmediate(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = 0x29
	sim.Memory[1] = 4
	sim.REGISTER_A = 5
	if sim.REGISTER_A != 5 {
		t.Log(("a correct before run"))
		t.FailNow()
	}

	sim.Run(1)
	//0100 //4
	//0101 //5 &
	//0100 // result in 4

	if sim.REGISTER_A != 4 {
		t.Log(("a not correct after run"))
		t.FailNow()
	}
}

func TestANDABS(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = 0x2d
	sim.Memory[1] = 5
	sim.Memory[2] = 0
	sim.Memory[5] = 4
	sim.REGISTER_A = 5
	if sim.REGISTER_A != 5 {
		t.Log(("a not 5 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.REGISTER_A != 4 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestANDABSX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_X = 20
	sim.REGISTER_A = 5

	sim.Memory[0] = AND_OPCODE_ABSX
	sim.Memory[1] = 5
	sim.Memory[2] = 0
	sim.Memory[25] = 4
	if sim.REGISTER_A != 5 {
		t.Log(("a not 5 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.REGISTER_A != 4 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}
func TestANDABSY(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_Y = 20
	sim.REGISTER_A = 5

	sim.Memory[0] = AND_OPCODE_ABSY
	sim.Memory[1] = 5
	sim.Memory[2] = 0
	sim.Memory[25] = 4
	if sim.REGISTER_A != 5 {
		t.Log(("a not 5 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.REGISTER_A != 4 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestANDZP(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.REGISTER_A = 5
	sim.Memory[0] = AND_OPCODE_ZP
	sim.Memory[1] = 25
	sim.Memory[25] = 4
	if sim.REGISTER_A != 5 {
		t.Log(("a not 5 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.REGISTER_A != 4 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestANDZPX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_X = 20
	sim.REGISTER_A = 5
	sim.Memory[0] = AND_OPCODE_ZPX
	sim.Memory[1] = 25
	sim.Memory[45] = 4
	if sim.REGISTER_A != 5 {
		t.Log(("a not 5 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.REGISTER_A != 4 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestANDINDX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_X = 4
	sim.REGISTER_A = 5
	sim.Memory[0] = AND_OPCODE_INDX
	sim.Memory[1] = 20
	sim.Memory[24] = 101
	sim.Memory[101] = 4
	if sim.REGISTER_A != 5 {
		t.Log(("a not 5 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.REGISTER_A != 4 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestANDINDY(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_Y = 10
	sim.REGISTER_A = 5
	sim.Memory[0] = AND_OPCODE_INDY
	sim.Memory[1] = 86
	sim.Memory[86] = 0x28
	sim.Memory[87] = 0x40
	sim.Memory[16434] = 4
	if sim.REGISTER_A != 5 {
		t.Log(("a not 5 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.REGISTER_A != 4 {
		t.Log("a not correct", sim.REGISTER_A)
		t.FailNow()
	}
}

func TestANDImmediateFlags(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_A = 255

	sim.Memory[0] = AND_OPCODE_IMM
	sim.Memory[1] = 255
	sim.Memory[2] = AND_OPCODE_IMM
	sim.Memory[3] = 0

	sim.Run(1)
	if sim.REGISTER_A != 255 {
		t.FailNow()
	}
	//A should be 255 and negative should be high
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE) != true {
		t.FailNow()
	}

	sim.Run(1)
	if sim.REGISTER_A != 0 {
		t.FailNow()
	}
	//Z should be high
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO) != true {
		t.FailNow()
	}
}
