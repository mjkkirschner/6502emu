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

func TestAddWithCarryABS(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.Memory[0] = 0x6d
	sim.Memory[1] = 5
	sim.Memory[2] = 0
	sim.Memory[5] = 100
	if sim.Register_A != 0 {
		t.Log(("a not 0 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Register_A != 100 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestAddWithCarryABSX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_X = 20

	sim.Memory[0] = 0x7d
	sim.Memory[1] = 5
	sim.Memory[2] = 0
	sim.Memory[25] = 101
	if sim.Register_A != 0 {
		t.Log(("a not 0 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Register_A != 101 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}
func TestAddWithCarryABSY(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_Y = 20

	sim.Memory[0] = 0x79
	sim.Memory[1] = 5
	sim.Memory[2] = 0
	sim.Memory[25] = 101
	if sim.Register_A != 0 {
		t.Log(("a not 0 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Register_A != 101 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestAddWithCarryZP(t *testing.T) {
	sim := NewSimulatorFromInstructionData()

	sim.Memory[0] = 0x65
	sim.Memory[1] = 25
	sim.Memory[25] = 101
	if sim.Register_A != 0 {
		t.Log(("a not 0 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Register_A != 101 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestAddWithCarryZPX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_X = 20

	sim.Memory[0] = 0x75
	sim.Memory[1] = 25
	sim.Memory[45] = 101
	if sim.Register_A != 0 {
		t.Log(("a not 0 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Register_A != 101 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestAddWithCarryINDX(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_X = 4

	sim.Memory[0] = 0x61
	sim.Memory[1] = 20
	sim.Memory[24] = 101
	sim.Memory[101] = 255
	if sim.Register_A != 0 {
		t.Log(("a not 0 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Register_A != 255 {
		t.Log(("a not correct"))
		t.FailNow()
	}
}

func TestAddWithCarryINDY(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	sim.REGISTER_Y = 10

	sim.Memory[0] = 0x71
	sim.Memory[1] = 86
	sim.Memory[86] = 0x28
	sim.Memory[87] = 0x40
	sim.Memory[16434] = 111
	if sim.Register_A != 0 {
		t.Log(("a not 0 before run"))
		t.FailNow()
	}

	sim.Run(1)

	if sim.Register_A != 111 {
		t.Log("a not correct", a)
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
