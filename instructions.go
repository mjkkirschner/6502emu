package main

func INSTRUCTION_AND_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//calculate binary AND of accumulator and operand and store result in accumulator.
	a := sim.Register_A
	b := operands.operands[0].(uint8)
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

func INSTRUCTION_ADC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//calculate the result.
	a := sim.Register_A
	b := operands.operands[0].(uint8)
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
	sim.computeZeroFlag(sim.Register_A)
	sim.computeNegativeFlag(sim.Register_A)
}

func INSTRUCTION_CLC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
}

func INSTRUCTION_SEC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
}

func INSTRUCTION_ASL_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//I think we'll only ever have operand A.
	//to determine where to store the result, check the returnAddress if it exists.
	var result uint8 = 0
	//grab the value we need to shift left.
	a := operands.operands[0].(uint8)
	carrycheck := a&128 > 0
	result = a << 1
	switch instruction.addressMode {
	case ACCUMULATOR:
		sim.Register_A = result
	//should cover all other cases
	default:
		sim.Memory[operands.returnAddress] = result
	}
	sim.computeCarryFlag(carrycheck)
	sim.computeZeroFlag(result)
	sim.computeNegativeFlag(result)

}

func INSTRUCTION_BCC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//if carry is 0 branch.
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY) == false {
		//branch
		sim.REGISTER_PC = operands.operands[0].(uint16)
		sim.X_JUMPING = true
	}
}
func INSTRUCTION_BCS_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//if carry is 1 branch.
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY) == true {
		//branch
		sim.REGISTER_PC = operands.operands[0].(uint16)
		sim.X_JUMPING = true
	}
}
func INSTRUCTION_BEQ_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//if zero is 1 branch
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO) == true {
		//branch
		sim.REGISTER_PC = operands.operands[0].(uint16)
		sim.X_JUMPING = true
	}
}
func INSTRUCTION_BNE_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//if zero is 0 branch
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO) == false {
		//branch
		sim.REGISTER_PC = operands.operands[0].(uint16)
		sim.X_JUMPING = true
	}
}
func INSTRUCTION_BMI_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//if negative is 1 branch
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE) == true {
		//branch
		sim.REGISTER_PC = operands.operands[0].(uint16)
		sim.X_JUMPING = true
	}
}
func INSTRUCTION_BPL_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//if negative is 0 branch
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE) == false {
		//branch
		sim.REGISTER_PC = operands.operands[0].(uint16)
		sim.X_JUMPING = true
	}
}

func INSTRUCTION_BVC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//if overflow is 0 branch
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW) == false {
		//branch
		sim.REGISTER_PC = operands.operands[0].(uint16)
		sim.X_JUMPING = true
	}
}

func INSTRUCTION_BVS_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//if overflow is 1 branch
	if sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW) == true {
		//branch
		sim.REGISTER_PC = operands.operands[0].(uint16)
		sim.X_JUMPING = true
	}
}
