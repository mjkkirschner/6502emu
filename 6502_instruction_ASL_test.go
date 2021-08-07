package main

import (
	"testing"
)

func TestASLABS(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = ASL_OPCODE_ABS
	sim.Memory[1] = 5
	sim.Memory[2] = 0
	sim.Memory[5] = 1

	if sim.Memory[5] != 1 {
		t.Log(("not correct before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Memory[5] != 2 {
		t.Log(("not correct after run"))
		t.FailNow()
	}
}

func TestASLABSX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_X = 20
	sim.Memory[0] = ASL_OPCODE_ABSX
	sim.Memory[1] = 5
	sim.Memory[2] = 0
	sim.Memory[3] = ASL_OPCODE_ABSX
	sim.Memory[4] = 5
	sim.Memory[5] = 0
	sim.Memory[25] = 1

	if sim.Memory[25] != 1 {
		t.Log(("not correct before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Memory[25] != 2 {
		t.Log(("not correct after run"))
		t.FailNow()
	}
	sim.Run(1)

	//shifting again should increment to 4
	if sim.Memory[25] != 4 {
		t.Log(("not correct after run"))
		t.FailNow()
	}
}

func TestASLZP(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = ASL_OPCODE_ZP
	sim.Memory[1] = 25
	sim.Memory[25] = 4
	if sim.Memory[25] != 4 {
		t.Log(("not correct before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Memory[25] != 8 {
		t.Log(("not correct after run"))
		t.FailNow()
	}
}

func TestASLZPX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_X = 10
	sim.Memory[0] = ASL_OPCODE_ZPX
	sim.Memory[1] = 25
	sim.Memory[35] = 4
	if sim.Memory[35] != 4 {
		t.Log(("not correct before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Memory[35] != 8 {
		t.Log(("not correct after run"))
		t.FailNow()
	}
}

func TestASLACC(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Register_A = 8
	sim.Memory[0] = ASL_OPCODE_ACC

	if sim.Register_A != 8 {
		t.Log(("not correct before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Register_A != 16 {
		t.Log(("not correct after run"))
		t.FailNow()
	}
}

func TestASLFlags(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Register_A = 255
	sim.Memory[0] = ASL_OPCODE_ACC
	sim.Memory[1] = ASL_OPCODE_ACC
	sim.Memory[2] = ASL_OPCODE_ACC
	sim.Memory[3] = ASL_OPCODE_ACC
	sim.Memory[4] = ASL_OPCODE_ACC
	sim.Memory[5] = ASL_OPCODE_ACC
	sim.Memory[6] = ASL_OPCODE_ACC
	sim.Memory[7] = ASL_OPCODE_ACC

	sim.Run(8)
	if sim.Register_A != 0 {
		t.Log("a not 0")
		t.FailNow()
	}
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY) != true {
		t.Log("carry not 1")
		t.FailNow()
	}
	//Z should be high
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO) != true {
		t.Log("z not 1")
		t.FailNow()
	}
}
