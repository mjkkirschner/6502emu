package main

import "testing"

func TestBITZP(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = BIT_OPCODE_ZP
	sim.Memory[1] = 5
	sim.Memory[5] = 0

	sim.Run(1)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE) != false {
		t.FailNow()
	}
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO) != true {
		t.FailNow()
	}
}

func TestBITABS(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = BIT_OPCODE_ABS
	sim.Memory[1] = 5
	sim.Memory[2] = 0
	sim.Memory[5] = 128

	sim.Run(1)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE) != true {
		t.FailNow()
	}
	//confusing but - this is set if the AND of A and M is 0 (0 & 128 is 0)
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO) != true {
		t.FailNow()
	}
}
