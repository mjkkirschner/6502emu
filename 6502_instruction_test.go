package main

import (
	"testing"
)

func TestAddWithCarryImmediate(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = 0x69
	sim.Memory[1] = 5
	if sim.Register_A != 0 {
		t.FailNow()
	}

	sim.Run(1)

	if sim.Register_A != 5 {
		t.FailNow()
	}

}
func TestAddWithCarryImmediateFlags(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = 0x69
	sim.Memory[1] = 5
	sim.Memory[2] = 0x69
	sim.Memory[3] = 255
	sim.Memory[4] = 0x69
	sim.Memory[5] = 0
	if sim.Register_A != 0 {
		t.FailNow()
	}

	sim.Run(2)
	//carry should be high
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY) != 1 {
		t.FailNow()
	}
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW) != 0 {
		t.FailNow()
	}
	//5+255 rolls over to 256 + 4
	if sim.Register_A != 4 {
		t.FailNow()
	}
	sim.Run(1)
	//carry should clear
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY) != 0 {
		t.FailNow()
	}

	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW) != 0 {
		t.FailNow()
	}
}

func NewSimulatorFromInstructionData() *Simulator {
	var filePath string = "6502ops.csv"
	instructions := GenerateInstructionMap(filePath)
	return NewSimulator(instructions)
}
