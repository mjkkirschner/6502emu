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
