package main

import "testing"

func TestBootROM(t *testing.T) {
	sim := NewSimulatorFromInstructionData()
	//read binary
	sim.loadMemoryFromBinaryAtAddress("c64roms/kernal.901227-02.bin", 0xE000)
	sim.loadMemoryFromBinaryAtAddress("c64roms/characters.901225-01.bin", 0xD000)
	sim.loadMemoryFromBinaryAtAddress("c64roms/basic.901226-01.bin", 0xA000)
	sim.reset()
	sim.Verbose = true
	//run 'forever'...
	for i := 0; i < 1000000; i++ {
		sim.Run(1)
	}
	sim.printMemoryAt(0x400, 400)
	//TODO assert commodore is printed in memory.
	t.Fail()
}
