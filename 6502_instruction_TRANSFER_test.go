package main

import "testing"

func TestTAX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = LDA_OPCODE_ABS
	sim.Memory[1] = 10
	sim.Memory[2] = 00
	sim.Memory[10] = 50
	sim.Memory[3] = TAX_OPCODE

	sim.Run(2)

	if sim.REGISTER_X != 50 {
		t.FailNow()
	}
	if sim.REGISTER_PC != 4 {
		t.FailNow()
	}
}
func TestTXA(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = LDX_OPCODE_ABS
	sim.Memory[1] = 10
	sim.Memory[2] = 00
	sim.Memory[10] = 50
	sim.Memory[3] = TXA_OPCODE

	sim.Run(2)

	if sim.REGISTER_A != 50 {
		t.FailNow()
	}
	if sim.REGISTER_PC != 4 {
		t.FailNow()
	}
}
func TestTAY(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = LDA_OPCODE_ABS
	sim.Memory[1] = 10
	sim.Memory[2] = 00
	sim.Memory[10] = 50
	sim.Memory[3] = TAY_OPCODE

	sim.Run(2)

	if sim.REGISTER_Y != 50 {
		t.FailNow()
	}
	if sim.REGISTER_PC != 4 {
		t.FailNow()
	}
}
func TestTYA(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = LDY_OPCODE_ABS
	sim.Memory[1] = 10
	sim.Memory[2] = 00
	sim.Memory[10] = 50
	sim.Memory[3] = TYA_OPCODE

	sim.Run(2)

	if sim.REGISTER_A != 50 {
		t.FailNow()
	}
	if sim.REGISTER_PC != 4 {
		t.FailNow()
	}
}
func TestTSX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = LDA_OPCODE_IMM
	sim.Memory[1] = 55
	sim.Memory[2] = PHA_OPCODE
	sim.Memory[3] = PHA_OPCODE
	sim.Memory[4] = PHA_OPCODE
	sim.Memory[5] = PHA_OPCODE
	sim.Memory[6] = TSX_OPCODE

	sim.Run(6)

	if sim.REGISTER_A != 55 {
		t.FailNow()
	}
	//0 - 4 = 252
	if sim.REGISTER_X != 252 {
		t.FailNow()
	}
	if sim.REGISTER_PC != 7 {
		t.FailNow()
	}
}
func TestTXS(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = LDX_OPCODE_IMM
	sim.Memory[1] = 55
	sim.Memory[2] = TXS_OPCODE

	sim.Run(2)

	if sim.REGISTER_STACKPOINTER != 55 {
		t.FailNow()
	}
	if sim.REGISTER_PC != 3 {
		t.FailNow()
	}
}
