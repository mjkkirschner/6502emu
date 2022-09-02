package main

import "testing"

func TestBRK(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	///0102 = 258 for our break handler
	sim.Memory[0xFFFF] = 1
	sim.Memory[0xFFFE] = 2
	sim.Memory[258] = ADDWITHCARRY_OPCODE_IMM
	sim.Memory[259] = 111
	sim.Memory[0] = BRK_OPCODE

	sim.Run(2)

	if sim.REGISTER_A != 111 {
		t.FailNow()
	}
}
