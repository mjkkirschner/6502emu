package main

func INSTRUCTION_AND_IMPLEMENTATION(sim *Simulator, operands []interface{}, instruction InstructionData) {
	//calculate binary AND of accumulator and operand and store result in accumulator.
	a := sim.Register_A
	b := (operands[0]).(uint8)
	c := a & b
	sim.Register_A = c
	//TODO factor out into shared util.
	if sim.Register_A == 0 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	}
	//set n
	nbit := sim.GetBit(REGISTER_A, 7)
	if nbit {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	}
}

func INSTRUCTION_ADC_IMPLEMENTATION(sim *Simulator, operands []interface{}, instruction InstructionData) {
	//calculate the result.
	a := sim.Register_A
	b := (operands[0]).(uint8)
	c := btou(sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY))
	sum := sim.Register_A + b + c

	carryCheck := uint16(a) + uint16(b) + uint16(c)
	overFlowCheck := (a ^ sum) & (b ^ sum) & 0x80 //negative bit.

	sim.Register_A = sim.Register_A + b + c
	//if the addition resulted in an overflow carry should be set to 1 - if not carry should be reset to 0.
	if carryCheck > 255 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
	}
	//overflow occurs when signed arithmetic overflows.
	if overFlowCheck == 1 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW)
	}
	//zero flag if result is 0
	if sim.Register_A == 0 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	}
	//set n
	nbit := sim.GetBit(REGISTER_A, 7)
	if nbit {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	}
}

func INSTRUCTION_CLC_IMPLEMENTATION(sim *Simulator, operands []interface{}, instruction InstructionData) {
	sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
}
