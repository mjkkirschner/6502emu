package main

import "testing"

func TestBCC(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = BCC_OPCODE
	sim.Memory[1] = 5

	if sim.REGISTER_PC != 0 {
		t.FailNow()
	}

	sim.Run(1)

	if sim.REGISTER_PC != 5 {
		t.FailNow()
	}
}

func TestBCS(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = SEC_OPCODE
	sim.Memory[1] = BCS_OPCODE
	sim.Memory[2] = 50

	if sim.REGISTER_PC != 0 {
		t.FailNow()
	}

	sim.Run(2)
	//1 + 50
	if sim.REGISTER_PC != 51 {
		t.FailNow()
	}
}

func TestBEQ(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = ADDWITHCARRY_OPCODE_IMM
	sim.Memory[1] = 0
	sim.Memory[2] = BEQ_OPCODE
	sim.Memory[3] = 50

	if sim.REGISTER_PC != 0 {
		t.FailNow()
	}

	sim.Run(2)
	if sim.REGISTER_PC != 52 {
		t.FailNow()
	}
}
func TestBNE(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = ADDWITHCARRY_OPCODE_IMM
	sim.Memory[1] = 1
	sim.Memory[2] = BNE_OPCODE
	sim.Memory[3] = 50

	if sim.REGISTER_PC != 0 {
		t.FailNow()
	}

	sim.Run(2)
	if sim.REGISTER_PC != 52 {
		t.FailNow()
	}
}
func TestBMI(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = ADDWITHCARRY_OPCODE_IMM
	//-2
	sim.Memory[1] = 254
	sim.Memory[2] = BMI_OPCODE
	sim.Memory[3] = 50

	if sim.REGISTER_PC != 0 {
		t.FailNow()
	}

	sim.Run(2)
	if sim.REGISTER_PC != 52 {
		t.FailNow()
	}
}
func TestBPL(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = ADDWITHCARRY_OPCODE_IMM
	sim.Memory[1] = 1
	sim.Memory[2] = BPL_OPCODE
	sim.Memory[3] = 50

	if sim.REGISTER_PC != 0 {
		t.FailNow()
	}

	sim.Run(2)
	if sim.REGISTER_PC != 52 {
		t.FailNow()
	}
}
func TestBVC(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = ADDWITHCARRY_OPCODE_IMM
	sim.Memory[1] = 1
	sim.Memory[2] = BVC_OPCODE
	sim.Memory[3] = 50

	if sim.REGISTER_PC != 0 {
		t.FailNow()
	}

	sim.Run(2)
	if sim.REGISTER_PC != 52 {
		t.FailNow()
	}
}
func TestBVS(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Register_A = 131 //-125
	sim.Memory[0] = ADDWITHCARRY_OPCODE_IMM
	sim.Memory[1] = 131 //-125
	//-125 + -125 = -250 (signed overflow, max of -127)
	sim.Memory[2] = BVS_OPCODE
	sim.Memory[3] = 50

	if sim.REGISTER_PC != 0 {
		t.FailNow()
	}

	sim.Run(2)
	if sim.REGISTER_PC != 52 {
		t.FailNow()
	}
}
