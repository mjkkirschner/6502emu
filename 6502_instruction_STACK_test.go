package main

import "testing"

func TestPHA_PLA(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = LDA_OPCODE_IMM
	sim.Memory[1] = 222
	sim.Memory[2] = PHA_OPCODE
	sim.Memory[3] = LDA_OPCODE_IMM
	sim.Memory[4] = 15         // set acc to some random number
	sim.Memory[5] = PLA_OPCODE // pull 222 from stack to acc
	sim.Run(6)

	//acc should have 222 in it
	if sim.REGISTER_A != 222 {
		t.Fail()
	}
}

func TestPHP_PLP(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = LDA_OPCODE_IMM
	sim.Memory[1] = 251 //-5
	//status should be negative now.
	//check status is not negative.
	//push status to stack
	sim.Memory[2] = PHP_OPCODE
	sim.Memory[3] = LDA_OPCODE_IMM
	sim.Memory[4] = 15 // set acc to a positive number
	//check status is not negative.
	sim.Memory[5] = PLP_OPCODE // pull stack to status reg.

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE) != false {
		t.Fail()
	}

	sim.Run(1)

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE) != true {
		t.Fail()
	}
	sim.Run(2)
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE) != false {
		t.Fail()
	}
	sim.Run(1)
	//status should negative bit high
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE) != true {
		t.Fail()
	}
	if sim.REGISTER_PC != 6 {
		t.Fail()
	}
}
