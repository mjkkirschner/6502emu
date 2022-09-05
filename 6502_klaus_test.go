package main

import "testing"

func TestRunKlausTestBinary(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	//read binary
	sim.loadMemoryFromBinary("6502_functional_test.bin")
	//set pc to start offset.
	sim.REGISTER_PC = 0x400
	//run 'forever'...
	for i := 0; i < 46000; i++ {
		sim.Run(1)
	}
}
