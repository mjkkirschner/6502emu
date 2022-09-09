package main

import (
	"fmt"
	"os"
)

func (sim *Simulator) printMemoryAt(address uint16, width uint16) {
	start := address - width/2
	end := address + width/2
	for i := start; i <= end; i++ {
		fmt.Println(i, ":", sim.Memory[i])
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (sim *Simulator) resetMemoryFromBinary(filePath string) {
	data, err := os.ReadFile(filePath)
	check(err)
	fmt.Println(len(data))
	//TODO is this safe? Should we instead iterate?
	sim.Memory = data
}

func (sim *Simulator) loadMemoryFromBinaryAtAddress(filePath string, address int) {
	data, err := os.ReadFile(filePath)
	check(err)
	fmt.Println(len(data))
	j := 0
	for i := 0; i < len(data); i++ {
		sim.Memory[address+i] = data[j]
		j++
	}
}
