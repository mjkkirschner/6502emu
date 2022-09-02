package main

import "testing"

func TestCLC(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = SEC_OPCODE
	sim.Memory[1] = NOP_OPCODE
	sim.Memory[2] = CLC_OPCODE

	sim.Run(1)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY) != true {
		t.FailNow()
	}
	sim.Run(2)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY) != false {
		t.FailNow()
	}
}
func TestCLD(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = SED_OPCODE
	sim.Memory[1] = NOP_OPCODE
	sim.Memory[2] = CLD_OPCODE

	sim.Run(1)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_DECIMAL) != true {
		t.FailNow()
	}
	sim.Run(2)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_DECIMAL) != false {
		t.FailNow()
	}
}
func TestCLI(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = SEI_OPCODE
	sim.Memory[1] = NOP_OPCODE
	sim.Memory[2] = CLI_OPCODE

	sim.Run(1)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_INTERRUPT_DISABLE) != true {
		t.FailNow()
	}
	sim.Run(2)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_INTERRUPT_DISABLE) != false {
		t.FailNow()
	}
}
func TestCLV(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.REGISTER_A = 131 //-125
	sim.Memory[0] = ADDWITHCARRY_OPCODE_IMM
	sim.Memory[1] = 131 //-125
	//-125 + -125 = -250 (signed overflow, max of -127)
	sim.Memory[2] = CLV_OPCODE //clear overflow
	sim.Run(1)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW) != true {
		t.FailNow()
	}
	sim.Run(1)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW) != false {
		t.FailNow()
	}
}
