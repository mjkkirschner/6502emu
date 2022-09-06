package main

import "fmt"

func INSTRUCTION_AND_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//calculate binary AND of accumulator and operand and store result in accumulator.
	a := sim.REGISTER_A
	b := operands.operands[0].(uint8)
	c := a & b
	sim.REGISTER_A = c
	//TODO factor out into shared util.
	if sim.REGISTER_A == 0 {
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
	a := sim.REGISTER_A
	b := operands.operands[0].(uint8)
	c := btou(sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY))
	sum := sim.REGISTER_A + b + c

	carryCheck := uint16(a) + uint16(b) + uint16(c)
	overFlowCheck := (a ^ sum) & (b ^ sum) & 0x80 //negative bit.

	sim.REGISTER_A = sim.REGISTER_A + b + c
	//if the addition resulted in an overflow carry should be set to 1 - if not carry should be reset to 0.
	if carryCheck > 255 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
	}
	//overflow occurs when signed arithmetic overflows.
	if overFlowCheck > 0 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW)
	}
	sim.computeZeroFlag(sim.REGISTER_A)
	sim.computeNegativeFlag(sim.REGISTER_A)
}

/*
A,Z,C,N = A-M-(1-C)

This instruction subtracts the contents of a memory location to the accumulator together with the not of the carry bit.
If overflow occurs the carry bit is clear, this enables multiple byte subtraction to be performed.
*/
func INSTRUCTION_SBC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	modops := operands
	//invert the operand before adding to do subtraction.
	modops.operands[0] = modops.operands[0].(uint8) ^ 0xFF
	INSTRUCTION_ADC_IMPLEMENTATION(sim, modops, instruction)
}

func INSTRUCTION_TAX_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//calculate the result.
	a := sim.REGISTER_A
	sim.REGISTER_X = a
	sim.computeZeroFlag(sim.REGISTER_A)
	sim.computeNegativeFlag(sim.REGISTER_A)
}

func INSTRUCTION_TXA_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//calculate the result.
	sim.REGISTER_A = sim.REGISTER_X
	sim.computeZeroFlag(sim.REGISTER_A)
	sim.computeNegativeFlag(sim.REGISTER_A)
}
func INSTRUCTION_TYA_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//calculate the result.
	sim.REGISTER_A = sim.REGISTER_Y
	sim.computeZeroFlag(sim.REGISTER_A)
	sim.computeNegativeFlag(sim.REGISTER_A)
}
func INSTRUCTION_TAY_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//calculate the result.
	sim.REGISTER_Y = sim.REGISTER_A
	sim.computeZeroFlag(sim.REGISTER_A)
	sim.computeNegativeFlag(sim.REGISTER_A)
}
func INSTRUCTION_TSX_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//calculate the result.
	sim.REGISTER_X = sim.REGISTER_STACKPOINTER
	sim.computeZeroFlag(sim.REGISTER_STACKPOINTER)
	sim.computeNegativeFlag(sim.REGISTER_STACKPOINTER)
}
func INSTRUCTION_TXS_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//calculate the result.
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_X)
}

func INSTRUCTION_CLC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
}

func INSTRUCTION_CLD_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_DECIMAL)
}
func INSTRUCTION_CLI_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_INTERRUPT_DISABLE)
}
func INSTRUCTION_CLV_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW)
}
func INSTRUCTION_NOP_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {

}
func INSTRUCTION_SEC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY)
}
func INSTRUCTION_SED_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_DECIMAL)
}
func INSTRUCTION_SEI_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_INTERRUPT_DISABLE)
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
		sim.REGISTER_A = result
	//should cover all other cases
	default:
		sim.Memory[operands.returnAddress] = result
	}
	sim.computeCarryFlag(carrycheck)
	sim.computeZeroFlag(result)
	sim.computeNegativeFlag(result)

}

func INSTRUCTION_LSR_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//I think we'll only ever have operand A.
	//to determine where to store the result, check the returnAddress if it exists.
	var result uint8 = 0
	//grab the value we need to shift left.
	a := operands.operands[0].(uint8)
	carrycheck := a&1 > 0
	result = a >> 1
	switch instruction.addressMode {
	case ACCUMULATOR:
		sim.REGISTER_A = result
	//should cover all other cases
	default:
		sim.Memory[operands.returnAddress] = result
	}
	sim.computeCarryFlag(carrycheck)
	sim.computeZeroFlag(result)
	sim.computeNegativeFlag(result)

}

//Move each of the bits in either A or M one place to the right.
// Bit 7 is filled with the current value of the carry flag whilst the old bit 0 becomes the new carry flag value.

func INSTRUCTION_ROR_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//I think we'll only ever have operand A.
	//to determine where to store the result, check the returnAddress if it exists.
	var result uint8 = 0
	//grab the value we need to shift left.
	a := operands.operands[0].(uint8)
	carrycheck := a&1 > 0
	result = a >> 1
	result = result | 128*btou(sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY))
	switch instruction.addressMode {
	case ACCUMULATOR:
		sim.Set8BitRegister(REGISTER_A, result)
	//should cover all other cases
	default:
		sim.SetMemory(operands.returnAddress, result)
	}
	sim.computeCarryFlag(carrycheck)
	sim.computeZeroFlag(result)
	sim.computeNegativeFlag(result)

}

//Move each of the bits in either A or M one place to the left.
// Bit 0 is filled with the current value of the carry flag whilst the old bit 7 becomes the new carry flag value.

func INSTRUCTION_ROL_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//I think we'll only ever have operand A.
	//to determine where to store the result, check the returnAddress if it exists.
	var result uint8 = 0
	//grab the value we need to shift left.
	a := operands.operands[0].(uint8)
	carrycheck := a&128 > 0
	result = a << 1
	result = result | 1&btou(sim.GetBit(REGISTER_STATUS, BITFLAG_STATUS_CARRY))
	switch instruction.addressMode {
	case ACCUMULATOR:
		sim.REGISTER_A = result
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

func INSTRUCTION_BIT_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {

	m := operands.operands[0].(uint8)
	and := sim.REGISTER_A & m
	if and == 0 {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_ZERO)
	}
	v := GetBit(uint(m), 6)
	n := GetBit(uint(m), 7)
	if v {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_OVERFLOW)
	}
	if n {
		sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	} else {
		sim.ClearBit(REGISTER_STATUS, BITFLAG_STATUS_NEGATIVE)
	}
}

func INSTRUCTION_BRK_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//push program counter to stack - we push high then low - so when reading it off
	//its low high.

	stackRegionStart := memoryMap["STACK"].start
	pchigh := stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)

	sim.Memory[pchigh] = uint8(((sim.REGISTER_PC + 2) >> 8) & 0xff)
	//decrement sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER-1)

	pclow := stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)
	sim.Memory[pclow] = uint8((sim.REGISTER_PC + 2) & 0x00ff)

	//decrement sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER-1)

	//push status reg to stack
	addrForStatus := stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)
	sim.Memory[addrForStatus] = sim.REGISTER_STATUS_P

	//decrement sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER-1)
	//load IRQ vector from FFFE/F to pc

	addrlow := sim.Memory[0xFFFE]
	addrhigh := sim.Memory[0xFFFF]
	longaddr := uint16(addrhigh)<<8 | (uint16(addrlow) & 0xff)
	sim.REGISTER_PC = longaddr
	//set break flag high
	sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_B_FLAG)
	//set interupt disable flag high
	sim.SetBit(REGISTER_STATUS, BITFLAG_STATUS_INTERRUPT_DISABLE)
	//TODO unsure if this should be true. BRK should jump to the interupt request handler -right?
	sim.X_JUMPING = true
}

func INSTRUCTION_PHA_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//push A to stack

	stackRegionStart := memoryMap["STACK"].start
	//TODO fix all these overflows.
	stackaddr := uint16(stackRegionStart) + uint16(sim.REGISTER_STACKPOINTER)

	sim.Memory[stackaddr] = sim.REGISTER_A
	//decrement sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER-1)
}
func INSTRUCTION_PHP_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//push status to stack.

	stackRegionStart := memoryMap["STACK"].start
	stackaddr := uint16(stackRegionStart) + uint16(sim.REGISTER_STACKPOINTER)

	sim.SetMemory(stackaddr, sim.REGISTER_STATUS_P)
	//decrement sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER-1)
}

func INSTRUCTION_PLA_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//pull stack to A. //set status flags. n,z

	stackRegionStart := memoryMap["STACK"].start
	//increment sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER+1)
	stackaddrToPull := uint16(stackRegionStart) + uint16(sim.REGISTER_STACKPOINTER)

	sim.Set8BitRegister(REGISTER_A, sim.Memory[stackaddrToPull])
	sim.computeNegativeFlag(sim.REGISTER_A)
	sim.computeZeroFlag(sim.REGISTER_A)
}

func INSTRUCTION_PLP_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//pull stack to status.
	stackRegionStart := memoryMap["STACK"].start
	//increment sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER+1)
	stackaddrToPull := uint16(stackRegionStart) + uint16(sim.REGISTER_STACKPOINTER)

	sim.Set8BitRegister(REGISTER_STATUS, sim.Memory[stackaddrToPull])
	//b and r flags should not be affected by PLP.
	sim.setStatusFlagsDefault()
}

func INSTRUCTION_RTI_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//pull stack to status.
	//then pull PC from stack.
	stackRegionStart := memoryMap["STACK"].start
	//increment sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER+1)
	stackaddrToPull := stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)

	sim.REGISTER_STATUS_P = sim.Memory[stackaddrToPull]
	//increment sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER+1)
	//now get program counter
	stackaddrToPull = stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)
	addrlow := sim.Memory[stackaddrToPull]
	//increment sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER+1)
	//now get program counter
	stackaddrToPull = stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)
	addrhigh := sim.Memory[stackaddrToPull]
	longaddr := uint16(addrhigh)<<8 | (uint16(addrlow) & 0xff)
	sim.REGISTER_PC = longaddr
	sim.X_JUMPING = true
}

func INSTRUCTION_RTS_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	// pull PC-1 from stack.
	stackRegionStart := memoryMap["STACK"].start
	//increment sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER+1)
	//now get program counter
	stackaddrToPull := stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)
	addrlow := sim.Memory[stackaddrToPull]
	//increment sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER+1)
	//now get program counter
	stackaddrToPull = stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)
	addrhigh := sim.Memory[stackaddrToPull]
	longaddr := uint16(addrhigh)<<8 | (uint16(addrlow) & 0xff)
	sim.REGISTER_PC = longaddr + 1
	sim.X_JUMPING = true
}

func INSTRUCTION_CMP_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	m := operands.operands[0].(uint8)
	b := sim.REGISTER_A - m
	fmt.Println("comparing", m, "to", sim.REGISTER_A)
	sim.computeCarryFlag(sim.REGISTER_A >= m)
	sim.computeNegativeFlag(b)
	sim.computeZeroFlag(b)
}

func INSTRUCTION_CPX_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	m := operands.operands[0].(uint8)
	b := sim.REGISTER_X - m
	sim.computeCarryFlag(sim.REGISTER_X >= m)
	sim.computeNegativeFlag(b)
	sim.computeZeroFlag(b)
}

func INSTRUCTION_CPY_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	m := operands.operands[0].(uint8)
	b := sim.REGISTER_Y - m
	sim.computeCarryFlag(sim.REGISTER_Y >= m)
	sim.computeNegativeFlag(b)
	sim.computeZeroFlag(b)
}

func INSTRUCTION_EOR_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	m := operands.operands[0].(uint8)
	sim.Set8BitRegister(REGISTER_A, sim.REGISTER_A^m)
	sim.computeNegativeFlag(sim.REGISTER_A)
	sim.computeZeroFlag(sim.REGISTER_A)
}

func INSTRUCTION_ORA_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	m := operands.operands[0].(uint8)
	sim.Set8BitRegister(REGISTER_A, sim.REGISTER_A|m)
	sim.computeNegativeFlag(sim.REGISTER_A)
	sim.computeZeroFlag(sim.REGISTER_A)
}

func INSTRUCTION_DEC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	x := operands.operands[0].(uint8) - 1
	sim.SetMemory(operands.returnAddress, x)
	sim.computeNegativeFlag(x)
	sim.computeZeroFlag(x)
}
func INSTRUCTION_DEX_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.Set8BitRegister(REGISTER_X, sim.REGISTER_X-1)
	sim.computeNegativeFlag(sim.REGISTER_X)
	sim.computeZeroFlag(sim.REGISTER_X)
}
func INSTRUCTION_DEY_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.Set8BitRegister(REGISTER_Y, sim.REGISTER_Y-1)
	sim.computeNegativeFlag(sim.REGISTER_Y)
	sim.computeZeroFlag(sim.REGISTER_Y)
}
func INSTRUCTION_INX_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.Set8BitRegister(REGISTER_X, sim.REGISTER_X+1)
	sim.computeNegativeFlag(sim.REGISTER_X)
	sim.computeZeroFlag(sim.REGISTER_X)
}
func INSTRUCTION_INY_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.Set8BitRegister(REGISTER_Y, sim.REGISTER_Y+1)
	sim.computeNegativeFlag(sim.REGISTER_Y)
	sim.computeZeroFlag(sim.REGISTER_Y)
}
func INSTRUCTION_INC_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.SetMemory(operands.returnAddress, operands.operands[0].(uint8)+1)
	sim.computeNegativeFlag(sim.Memory[operands.returnAddress])
	sim.computeZeroFlag(sim.Memory[operands.returnAddress])
}

func INSTRUCTION_JMP_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//now set PC to address to jump to.
	sim.Set16BitRegister(REGISTER_PC, operands.returnAddress)
	sim.X_JUMPING = true
}

func INSTRUCTION_JSR_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	//push program counter to stack - we push high then low - so when reading it off
	//its low high.

	//importantly we want to push the address of the low byte of the JSR operand.
	//so that is PC + 1 because in this emulator we increment the PC AFTER the entire instruction...

	stackRegionStart := memoryMap["STACK"].start
	pchigh := stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)

	sim.Memory[pchigh] = uint8(((sim.REGISTER_PC + 2) >> 8) & 0xff)
	//decrement sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER-1)

	pclow := stackRegionStart + uint16(sim.REGISTER_STACKPOINTER)
	sim.Memory[pclow] = uint8((sim.REGISTER_PC + 2) & 0x00ff)

	//decrement sp
	sim.Set8BitRegister(REGISTER_STACKPOINTER, sim.REGISTER_STACKPOINTER-1)

	//now set PC to SR address to jump to.
	sim.REGISTER_PC = operands.returnAddress
	sim.X_JUMPING = true

}
func INSTRUCTION_LDA_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.Set8BitRegister(REGISTER_A, operands.operands[0].(uint8))
	sim.computeNegativeFlag(sim.REGISTER_A)
	sim.computeZeroFlag(sim.REGISTER_A)
}
func INSTRUCTION_LDX_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.REGISTER_X = operands.operands[0].(uint8)
	sim.computeNegativeFlag(sim.REGISTER_X)
	sim.computeZeroFlag(sim.REGISTER_X)
}

func INSTRUCTION_LDY_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.REGISTER_Y = operands.operands[0].(uint8)
	sim.computeNegativeFlag(sim.REGISTER_Y)
	sim.computeZeroFlag(sim.REGISTER_Y)
}

func INSTRUCTION_STA_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.SetMemory(operands.returnAddress, sim.REGISTER_A)
}
func INSTRUCTION_STX_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.SetMemory(operands.returnAddress, sim.REGISTER_X)
}
func INSTRUCTION_STY_IMPLEMENTATION(sim *Simulator, operands decodeResults, instruction InstructionData) {
	sim.SetMemory(operands.returnAddress, sim.REGISTER_Y)
}
