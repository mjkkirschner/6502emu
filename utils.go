package main

func btou(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func (sim *Simulator) computeZeroFlag(value uint8) {
	//zero flag if result is 0
	if value == 0 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	}
}

func (sim *Simulator) computeNegativeFlag(value uint8) {
	//set n
	nbit := GetBit(uint(value), 7)
	if nbit {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	}
}

func (sim *Simulator) computeCarryFlag(carryTrue bool) {
	if carryTrue {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
	}
}

func NewSimulatorFromInstructionData() *Simulator {
	var filePath string = "6502ops.csv"
	instructions := GenerateInstructionMap(filePath)
	sim := NewSimulator(instructions)
	//set reset vector to 0000 because all tests use 0 as start.
	sim.Memory[0xFFFC] = 0
	sim.Memory[0xFFFD] = 0
	sim.reset() //set PC to FFFC/FFFD
	sim.setStatusFlagsDefault()
	return sim
}
